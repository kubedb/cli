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

func MySQLConnectCMD(f cmdutil.Factory) *cobra.Command {
	var dbName string

	myConnectCmd := &cobra.Command{
		Use: "mysql",
		Aliases: []string{
			"my",
		},
		Short: "Connect to a mysql object",
		Long:  `Use this cmd to exec into a mysql object's primary pod.`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				log.Fatal("Enter mysql object's name as an argument")
			}
			dbName = args[0]

			namespace, _, err := f.ToRawKubeConfigLoader().Namespace()
			if err != nil {
				klog.Error(err, "failed to get current namespace")
			}

			opts, err := newmysqlOpts(f, dbName, namespace)
			if err != nil {
				log.Fatalln(err)
			}

			tunnel, err := lib.TunnelToDBService(opts.config, dbName, namespace, kubedb.MySQLDatabasePort)
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

	return myConnectCmd
}

func MySQLExecCMD(f cmdutil.Factory) *cobra.Command {
	var (
		dbName      string
		mysqlDBName string
		fileName    string
		command     string
	)

	myExecCmd := &cobra.Command{
		Use: "mysql",
		Aliases: []string{
			"my",
		},
		Short: "Execute SQL commands to a mysql resource",
		Long: `Use this cmd to execute sql commands to a mysql object's primary pod.

Examples:
  # Execute a script named 'demo.sql' in 'my-demo' mysql database in 'demo' namespace
  kubectl dba exec my my-demo -n demo -f demo.sql

  # Execute a command in 'my-demo' mysql database in 'demo' namespace
  kubectl dba exec my my-demo -c 'describe pet' -d kubedb"
				`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				log.Fatal("Enter mysql object's name as an argument. Your commands will be applied on a database inside it's primary pod")
			}
			dbName = args[0]

			namespace, _, err := f.ToRawKubeConfigLoader().Namespace()
			if err != nil {
				klog.Error(err, "failed to get current namespace")
			}

			opts, err := newmysqlOpts(f, dbName, namespace)
			if err != nil {
				log.Fatalln(err)
			}

			if fileName == "" && command == "" {
				log.Fatal("use --file or --command to execute supported commands to a mysql object")
			}

			tunnel, err := lib.TunnelToDBService(opts.config, dbName, namespace, kubedb.MySQLDatabasePort)
			if err != nil {
				log.Fatal("couldn't creat tunnel, error: ", err)
			}

			if command != "" {
				err = opts.executeCommand(tunnel.Local, command, mysqlDBName)
				if err != nil {
					log.Fatal(err)
				}
			}

			if fileName != "" {
				err = opts.executeFile(tunnel.Local, fileName, mysqlDBName)
				if err != nil {
					log.Fatal(err)
				}
			}

			tunnel.Close()
		},
	}

	myExecCmd.Flags().StringVarP(&fileName, "file", "f", "", "path to command file")
	myExecCmd.Flags().StringVarP(&command, "command", "c", "", "command to execute")
	myExecCmd.Flags().StringVarP(&mysqlDBName, "dbName", "d", "mysql", "name of the database inside mysql to execute command")

	return myExecCmd
}

type mysqlOpts struct {
	db       *dbapi.MySQL
	dbImage  string
	config   *rest.Config
	client   *kubernetes.Clientset
	dbClient *cs.Clientset

	username string
	pass     string
}

func newmysqlOpts(f cmdutil.Factory, dbName, namespace string) (*mysqlOpts, error) {
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

	db, err := dbClient.KubedbV1().MySQLs(namespace).Get(context.TODO(), dbName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	if db.Status.Phase != dbapi.DatabasePhaseReady {
		return nil, fmt.Errorf("mysql %s/%s is not ready", namespace, dbName)
	}

	dbVersion, err := dbClient.CatalogV1alpha1().MySQLVersions().Get(context.TODO(), db.Spec.Version, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	secret, err := client.CoreV1().Secrets(db.Namespace).Get(context.TODO(), db.Spec.AuthSecret.Name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return &mysqlOpts{
		db:       db,
		dbImage:  dbVersion.Spec.DB.Image,
		config:   config,
		client:   client,
		dbClient: dbClient,
		username: string(secret.Data[corev1.BasicAuthUsernameKey]),
		pass:     string(secret.Data[corev1.BasicAuthPasswordKey]),
	}, nil
}

func (opts *mysqlOpts) getDockerShellCommand(localPort int, dockerFlags, mysqlExtraFlags []interface{}) (*shell.Session, error) {
	sh := shell.NewSession()
	sh.ShowCMD = false

	db := opts.db
	dockerCommand := []interface{}{
		"run", "--network=host",
		"-e", fmt.Sprintf("MYSQL_PWD=%s", opts.pass),
	}
	dockerCommand = append(dockerCommand, dockerFlags...)

	mysqlCommand := []interface{}{
		"mysql",
		"--host=127.0.0.1", fmt.Sprintf("--port=%d", localPort),
		fmt.Sprintf("--user=%s", opts.username),
	}

	if db.Spec.TLS != nil {
		secretName := db.CertificateName(dbapi.MySQLClientCert)
		certSecret, err := opts.client.CoreV1().Secrets(db.Namespace).Get(context.TODO(), secretName, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}

		caCrt, ok := certSecret.Data[corev1.ServiceAccountRootCAKey]
		if !ok {
			return nil, fmt.Errorf("missing %s in secret %s/%s", corev1.ServiceAccountRootCAKey, certSecret.Namespace, certSecret.Name)
		}
		err = os.WriteFile(caFile, caCrt, 0o644)
		if err != nil {
			return nil, err
		}

		crt, ok := certSecret.Data[corev1.TLSCertKey]
		if !ok {
			return nil, fmt.Errorf("missing %s in secret %s/%s", corev1.TLSCertKey, certSecret.Namespace, certSecret.Name)
		}
		err = os.WriteFile(certFile, crt, 0o644)
		if err != nil {
			return nil, err
		}

		key, ok := certSecret.Data[corev1.TLSPrivateKeyKey]
		if !ok {
			return nil, fmt.Errorf("missing %s in secret %s/%s", corev1.TLSPrivateKeyKey, certSecret.Namespace, certSecret.Name)
		}
		err = os.WriteFile(keyFile, key, 0o644)
		if err != nil {
			return nil, err
		}

		dockerCommand = append(dockerCommand,
			"-v", fmt.Sprintf("%s:%s", caFile, caFile),
			"-v", fmt.Sprintf("%s:%s", certFile, certFile),
			"-v", fmt.Sprintf("%s:%s", keyFile, keyFile),
		)
		mysqlCommand = append(mysqlCommand,
			fmt.Sprintf("--ssl-ca=%v", caFile),
			fmt.Sprintf("--ssl-cert=%v", certFile),
			fmt.Sprintf("--ssl-key=%v", keyFile),
		)
	}

	dockerCommand = append(dockerCommand, opts.dbImage)
	finalCommand := append(dockerCommand, mysqlCommand...)
	if mysqlExtraFlags != nil {
		finalCommand = append(finalCommand, mysqlExtraFlags...)
	}
	return sh.Command("docker", finalCommand...).SetStdin(os.Stdin), nil
}

func (opts *mysqlOpts) connect(localPort int) error {
	dockerFlag := []interface{}{
		"-it",
	}
	shSession, err := opts.getDockerShellCommand(localPort, dockerFlag, nil)
	if err != nil {
		return err
	}

	return shSession.Run()
}

func (opts *mysqlOpts) executeCommand(localPort int, command, mysqlDBName string) error {
	mysqlExtraFlags := []interface{}{
		mysqlDBName,
		"-e", command,
	}
	shSession, err := opts.getDockerShellCommand(localPort, nil, mysqlExtraFlags)
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
	fmt.Printf("command applied successfully%s", output)

	return nil
}

func (opts *mysqlOpts) executeFile(localPort int, fileName, mysqlDBName string) error {
	fileName, err := filepath.Abs(fileName)
	if err != nil {
		return err
	}
	tempFileName := "/tmp/my.sql"

	dockerFlag := []interface{}{
		"-v", fmt.Sprintf("%s:%s", fileName, tempFileName),
	}
	mysqlExtraFlags := []interface{}{
		mysqlDBName,
		"-e", fmt.Sprintf("source %s", tempFileName),
	}
	shSession, err := opts.getDockerShellCommand(localPort, dockerFlag, mysqlExtraFlags)
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
	fmt.Printf("file %s applied successfully%s", fileName, output)

	return nil
}
