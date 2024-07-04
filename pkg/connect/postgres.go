/*
Copyright AppsCode Inc. and Contributors

Licensed under the AppsCode Community License 1.0.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://github.com/appscode/licenses/raw/1.0.0/AppsCode-Community-1.0.0.md

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package connect

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"kubedb.dev/apimachinery/apis/kubedb"
	dbapi "kubedb.dev/apimachinery/apis/kubedb/v1"
	cs "kubedb.dev/apimachinery/client/clientset/versioned"
	"kubedb.dev/cli/pkg/lib"

	"github.com/spf13/cobra"
	shell "gomodules.xyz/go-sh"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
)

func PostgresConnectCMD(f cmdutil.Factory) *cobra.Command {
	var (
		dbName         string
		postgresDBName string
	)

	pgConnectCmd := &cobra.Command{
		Use: "postgres",
		Aliases: []string{
			"postgresql",
			"pgsql",
			"pg",
		},
		Short: "Connect to a postgres object",
		Long:  `Use this cmd to exec into a postgres object's primary pod.`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				log.Fatal("Enter postgres object's name as an argument")
			}
			dbName = args[0]

			namespace, _, err := f.ToRawKubeConfigLoader().Namespace()
			if err != nil {
				klog.Error(err, "failed to get current namespace")
			}

			opts, err := newPostgresOpts(f, dbName, namespace, postgresDBName)
			if err != nil {
				log.Fatalln(err)
			}

			tunnel, err := lib.TunnelToDBService(opts.config, dbName, namespace, kubedb.PostgresDatabasePort)
			if err != nil {
				log.Fatal("couldn't create tunnel, error: ", err)
			}

			err = opts.connect(tunnel.Local)
			if err != nil {
				log.Fatal(err)
			}

			tunnel.Close()
		},
	}

	return pgConnectCmd
}

func PostgresExecCMD(f cmdutil.Factory) *cobra.Command {
	var (
		dbName         string
		postgresDBName string
		fileName       string
		command        string
	)

	pgExecCmd := &cobra.Command{
		Use: "postgres",
		Aliases: []string{
			"postgresql",
			"pgsql",
			"pg",
		},
		Short: "Execute SQL commands to a postgres resource",
		Long: `Use this cmd to execute postgresql commands to a postgres object's primary pod.

Examples:
  # Execute a script named 'demo.sql' in 'pg-demo' postgres database in 'demo' namespace
  kubectl dba exec pg pg-demo -n demo -f demo.sql

  # Execute a command in 'pg-demo' postgres database in 'demo' namespace
  kubectl dba exec pg pg-demo -c '\l' -d kubedb"
				`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				log.Fatal("Enter postgres object's name as an argument. Your commands will be applied on a database inside it's primary pod")
			}
			dbName = args[0]

			namespace, _, err := f.ToRawKubeConfigLoader().Namespace()
			if err != nil {
				klog.Error(err, "failed to get current namespace")
			}

			opts, err := newPostgresOpts(f, dbName, namespace, postgresDBName)
			if err != nil {
				log.Fatalln(err)
			}

			if fileName == "" && command == "" {
				log.Fatal("use --file or --command to execute supported commands to a postgres object")
			}

			tunnel, err := lib.TunnelToDBService(opts.config, dbName, namespace, kubedb.PostgresDatabasePort)
			if err != nil {
				log.Fatal("couldn't creat tunnel, error: ", err)
			}

			if command != "" {
				err = opts.executeCommand(tunnel.Local, command)
				if err != nil {
					log.Fatal(err)
				}
			}

			if fileName != "" {
				err = opts.executeFile(tunnel.Local, fileName)
				if err != nil {
					log.Fatal(err)
				}
			}

			tunnel.Close()
		},
	}

	pgExecCmd.Flags().StringVarP(&fileName, "file", "f", "", "path to command file")
	pgExecCmd.Flags().StringVarP(&command, "command", "c", "", "command to execute")
	pgExecCmd.Flags().StringVarP(&postgresDBName, "dbName", "d", "", "name of the database inside postgres to execute command")

	return pgExecCmd
}

type postgresOpts struct {
	db       *dbapi.Postgres
	dbImage  string
	config   *rest.Config
	client   *kubernetes.Clientset
	dbClient *cs.Clientset

	postgresDBName string

	username string
	pass     string

	errWriter *bytes.Buffer
}

func newPostgresOpts(f cmdutil.Factory, dbName, namespace, postgresDBName string) (*postgresOpts, error) {
	config, err := f.ToRESTConfig()
	if err != nil {
		return nil, err
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	dbClient, err := cs.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	db, err := dbClient.KubedbV1().Postgreses(namespace).Get(context.TODO(), dbName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	if db.Status.Phase != dbapi.DatabasePhaseReady {
		return nil, fmt.Errorf("postgres %s/%s is not ready", namespace, dbName)
	}

	dbVersion, err := dbClient.CatalogV1alpha1().PostgresVersions().Get(context.TODO(), db.Spec.Version, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	secret, err := client.CoreV1().Secrets(db.Namespace).Get(context.TODO(), db.Spec.AuthSecret.Name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return &postgresOpts{
		db:             db,
		dbImage:        dbVersion.Spec.DB.Image,
		config:         config,
		client:         client,
		dbClient:       dbClient,
		username:       string(secret.Data[corev1.BasicAuthUsernameKey]),
		pass:           string(secret.Data[corev1.BasicAuthPasswordKey]),
		errWriter:      &bytes.Buffer{},
		postgresDBName: postgresDBName,
	}, nil
}

func (opts *postgresOpts) getDockerShellCommand(localPort int, dockerFlags, postgresExtraFlags []interface{}) (*shell.Session, error) {
	sh := shell.NewSession()
	sh.ShowCMD = false
	sh.Stderr = opts.errWriter

	db := opts.db
	dockerCommand := []interface{}{
		"run", "--network=host",
		"-e", fmt.Sprintf("PGPASSWORD=%s", opts.pass),
	}
	dockerCommand = append(dockerCommand, dockerFlags...)

	postgresCommand := []interface{}{
		"psql",
		"--host=127.0.0.1", fmt.Sprintf("--port=%d", localPort),
		fmt.Sprintf("--username=%s", opts.username),
	}

	if db.Spec.TLS != nil {
		secretName := db.CertificateName(dbapi.PostgresClientCert)
		certSecret, err := opts.client.CoreV1().Secrets(db.Namespace).Get(context.TODO(), secretName, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}

		caCrt, ok := certSecret.Data[corev1.ServiceAccountRootCAKey]
		if !ok {
			return nil, fmt.Errorf("missing %s in secret %s/%s", corev1.ServiceAccountRootCAKey, certSecret.Namespace, certSecret.Name)
		}
		err = os.WriteFile(pgCaFile, caCrt, 0o644)
		if err != nil {
			return nil, err
		}

		crt, ok := certSecret.Data[corev1.TLSCertKey]
		if !ok {
			return nil, fmt.Errorf("missing %s in secret %s/%s", corev1.TLSCertKey, certSecret.Namespace, certSecret.Name)
		}
		err = os.WriteFile(pgCertFile, crt, 0o644)
		if err != nil {
			return nil, err
		}

		key, ok := certSecret.Data[corev1.TLSPrivateKeyKey]
		if !ok {
			return nil, fmt.Errorf("missing %s in secret %s/%s", corev1.TLSPrivateKeyKey, certSecret.Namespace, certSecret.Name)
		}
		err = os.WriteFile(pgKeyFile, key, 0o600)
		if err != nil {
			return nil, err
		}

		dockerCommand = append(dockerCommand,
			"-v", fmt.Sprintf("%s:%s", "/tmp/", "/root/.postgresql/"),
		)
	}

	dockerCommand = append(dockerCommand, opts.dbImage)
	finalCommand := append(dockerCommand, postgresCommand...)
	if postgresExtraFlags != nil {
		finalCommand = append(finalCommand, postgresExtraFlags...)
	}
	return sh.Command("docker", finalCommand...).SetStdin(os.Stdin), nil
}

func (opts *postgresOpts) connect(localPort int) error {
	dockerFlag := []interface{}{
		"-it",
	}
	shSession, err := opts.getDockerShellCommand(localPort, dockerFlag, nil)
	if err != nil {
		return err
	}

	err = shSession.Run()
	if err != nil {
		return err
	}

	return nil
}

func (opts *postgresOpts) executeCommand(localPort int, command string) error {
	dbFlag := ""
	if opts.postgresDBName != "" {
		dbFlag = fmt.Sprintf("--dbname=%s", opts.postgresDBName)
	}
	postgresExtraFlags := []interface{}{
		dbFlag,
		fmt.Sprintf("--command=%s", command),
	}
	shSession, err := opts.getDockerShellCommand(localPort, nil, postgresExtraFlags)
	if err != nil {
		return err
	}

	out, err := shSession.Output()
	if err != nil {
		return fmt.Errorf("failed to execute command, error: %s, output: %s\n", err, out)
	}
	output := ""
	if string(out) != "" {
		output = ", output:\n\n" + string(out)
	}

	errOutput := opts.errWriter.String()
	if errOutput != "" {
		return fmt.Errorf("failed to execute command, stderr: %s%s", errOutput, output)
	}
	fmt.Printf("command applied successfully%s", output)

	return nil
}

func (opts *postgresOpts) executeFile(localPort int, fileName string) error {
	dbFlag := ""
	if opts.postgresDBName != "" {
		dbFlag = fmt.Sprintf("--dbname=%s", opts.postgresDBName)
	}
	fileName, err := filepath.Abs(fileName)
	if err != nil {
		return err
	}
	tempFileName := "/tmp/postgres.sql"

	dockerFlag := []interface{}{
		"-v", fmt.Sprintf("%s:%s", fileName, tempFileName),
	}
	postgresExtraFlags := []interface{}{
		dbFlag,
		fmt.Sprintf("--file=%v", tempFileName),
	}
	shSession, err := opts.getDockerShellCommand(localPort, dockerFlag, postgresExtraFlags)
	if err != nil {
		return err
	}

	out, err := shSession.Output()
	if err != nil {
		return fmt.Errorf("failed to execute file, error: %s, output: %s\n", err, out)
	}

	output := ""
	if string(out) != "" {
		output = ", output:\n\n" + string(out)
	}

	errOutput := opts.errWriter.String()
	if errOutput != "" {
		return fmt.Errorf("failed to execute file, stderr: %s%s", errOutput, output)
	}

	fmt.Printf("file %s applied successfully%s", fileName, output)

	return nil
}
