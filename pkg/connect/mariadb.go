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
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	api "kubedb.dev/apimachinery/apis/kubedb/v1alpha2"
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

func NewMariadbCMD(f cmdutil.Factory) *cobra.Command {
	var (
		dbName        string
		mariadbDBName string
		namespace     string
		fileName      string
		command       string
	)

	currentNamespace, _, err := f.ToRawKubeConfigLoader().Namespace()
	if err != nil {
		klog.Error(err, "failed to get current namespace")
	}

	var mdCmd = &cobra.Command{
		Use: "mariadb",
		Aliases: []string{
			"md",
		},
		Short: "Use to operate mariadb pods",
		Long: `Use this cmd to operate mariadb pods. Available sub-commands:
				apply
				connect`,
		Run: func(cmd *cobra.Command, args []string) {},
	}

	var mdConnectCmd = &cobra.Command{
		Use:   "connect",
		Short: "Connect to a mariadb object's pod",
		Long:  `Use this cmd to exec into a mariadb object's primary pod.`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				log.Fatal("Enter mariadb object's name as an argument")
			}
			dbName = args[0]
			opts, err := newmariadbOpts(f, dbName, namespace)
			if err != nil {
				log.Fatalln(err)
			}

			tunnel, err := lib.TunnelToDBService(opts.config, dbName, namespace, api.MySQLDatabasePort)
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

	var mdApplyCmd = &cobra.Command{
		Use:   "apply",
		Short: "Apply SQL commands to a mariadb resource",
		Long: `Use this cmd to apply SQL commands from a file to a mariadb object's primary pod.
				Syntax: $ kubectl dba mariadb apply <mariadb-object-name> -n <namespace> -f <fileName>`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				log.Fatal("Enter mariadb object's name as an argument. Your commands will be applied on a database inside it's primary pod")
			}
			dbName = args[0]

			opts, err := newmariadbOpts(f, dbName, namespace)
			if err != nil {
				log.Fatalln(err)
			}

			if fileName == "" && command == "" {
				log.Fatal("use --file or --command to apply supported commands to a mariadb object's pods")
			}

			tunnel, err := lib.TunnelToDBService(opts.config, dbName, namespace, api.MySQLDatabasePort)
			if err != nil {
				log.Fatal("couldn't creat tunnel, error: ", err)
			}

			if command != "" {
				err = opts.applyCommand(tunnel.Local, command, mariadbDBName)
				if err != nil {
					log.Fatal(err)
				}
			}

			if fileName != "" {
				err = opts.applyFile(tunnel.Local, fileName, mariadbDBName)
				if err != nil {
					log.Fatal(err)
				}
			}

			tunnel.Close()
		},
	}

	mdCmd.AddCommand(mdConnectCmd)
	mdCmd.AddCommand(mdApplyCmd)
	mdCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", currentNamespace, "namespace of the mariadb object to connect to.")

	mdApplyCmd.Flags().StringVarP(&fileName, "file", "f", "", "path to command file")
	mdApplyCmd.Flags().StringVarP(&command, "command", "c", "", "command to execute")
	mdApplyCmd.Flags().StringVarP(&mariadbDBName, "dbName", "d", "mysql", "name of the database inside mariadb to execute command")

	return mdCmd
}

type mariadbOpts struct {
	db       *api.MariaDB
	dbImage  string
	config   *rest.Config
	client   *kubernetes.Clientset
	dbClient *cs.Clientset

	username string
	pass     string
}

func newmariadbOpts(f cmdutil.Factory, dbName, namespace string) (*mariadbOpts, error) {
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

	db, err := dbClient.KubedbV1alpha2().MariaDBs(namespace).Get(context.TODO(), dbName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	if db.Status.Phase != api.DatabasePhaseReady {
		return nil, fmt.Errorf("mariadb %s/%s is not ready", namespace, dbName)
	}

	dbVersion, err := dbClient.CatalogV1alpha1().MariaDBVersions().Get(context.TODO(), db.Spec.Version, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	secret, err := client.CoreV1().Secrets(db.Namespace).Get(context.TODO(), db.Spec.AuthSecret.Name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return &mariadbOpts{
		db:       db,
		dbImage:  dbVersion.Spec.DB.Image,
		config:   config,
		client:   client,
		dbClient: dbClient,
		username: string(secret.Data[corev1.BasicAuthUsernameKey]),
		pass:     string(secret.Data[corev1.BasicAuthPasswordKey]),
	}, nil
}

func (opts *mariadbOpts) getDockerShellCommand(localPort int, dockerFlags, mysqlExtraFlags []interface{}) (*shell.Session, error) {
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
		secretName := db.CertificateName(api.MariaDBArchiverCert)
		certSecret, err := opts.client.CoreV1().Secrets(db.Namespace).Get(context.TODO(), secretName, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}

		caCrt, ok := certSecret.Data[corev1.ServiceAccountRootCAKey]
		if !ok {
			return nil, fmt.Errorf("missing %s in secret %s/%s", corev1.ServiceAccountRootCAKey, certSecret.Namespace, certSecret.Name)
		}
		err = ioutil.WriteFile(caFile, caCrt, 0644)
		if err != nil {
			return nil, err
		}

		crt, ok := certSecret.Data[corev1.TLSCertKey]
		if !ok {
			return nil, fmt.Errorf("missing %s in secret %s/%s", corev1.TLSCertKey, certSecret.Namespace, certSecret.Name)
		}
		err = ioutil.WriteFile(certFile, crt, 0644)
		if err != nil {
			return nil, err
		}

		key, ok := certSecret.Data[corev1.TLSPrivateKeyKey]
		if !ok {
			return nil, fmt.Errorf("missing %s in secret %s/%s", corev1.TLSPrivateKeyKey, certSecret.Namespace, certSecret.Name)
		}
		err = ioutil.WriteFile(keyFile, key, 0644)
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

func (opts *mariadbOpts) connect(localPort int) error {
	dockerFlag := []interface{}{
		"-it",
	}
	shSession, err := opts.getDockerShellCommand(localPort, dockerFlag, nil)
	if err != nil {
		return err
	}

	return shSession.Run()
}

func (opts *mariadbOpts) applyCommand(localPort int, command, mysqlDBName string) error {
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
		return fmt.Errorf("failed to apply command, error: %s, output: %s\n", err, out)
	}
	output := ""
	if string(out) != "" {
		output = ", output:\n\n" + string(out)
	}
	fmt.Printf("command applied successfully%s", output)

	return nil
}

func (opts *mariadbOpts) applyFile(localPort int, fileName, mysqlDBName string) error {
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
		return fmt.Errorf("failed to apply file, error: %s, output: %s\n", err, out)
	}

	output := ""
	if string(out) != "" {
		output = ", output:\n\n" + string(out)
	}
	fmt.Printf("file %s applied successfully%s", fileName, output)

	return nil
}
