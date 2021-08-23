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

func MongoDBConnectCMD(f cmdutil.Factory) *cobra.Command {
	var (
		dbName    string
		namespace string
	)

	currentNamespace, _, err := f.ToRawKubeConfigLoader().Namespace()
	if err != nil {
		klog.Error(err, "failed to get current namespace")
	}

	var mgConnectCmd = &cobra.Command{
		Use: "mongodb",
		Aliases: []string{
			"mg",
			"mongo",
		},
		Short: "Connect to a mongodb object",
		Long: `Use this cmd to exec into a mongodb object's primary pod.
example:

kubectl dba connect mg <db-name> -n <db-namespace>`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				log.Fatal("Enter mongodb object's name as an argument")
			}
			dbName = args[0]
			opts, err := newMongodbOpts(f, dbName, namespace)
			if err != nil {
				log.Fatalln(err)
			}

			tunnel, err := lib.TunnelToDBService(opts.config, dbName, namespace, api.MongoDBDatabasePort)
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

	mgConnectCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", currentNamespace, "namespace of the mongodb object to connect to.")

	return mgConnectCmd
}

func MongoDBExecCMD(f cmdutil.Factory) *cobra.Command {
	var (
		dbName    string
		namespace string
		fileName  string
		command   string
	)

	currentNamespace, _, err := f.ToRawKubeConfigLoader().Namespace()
	if err != nil {
		klog.Error(err, "failed to get current namespace")
	}

	var mgExecCmd = &cobra.Command{
		Use: "mongodb",
		Aliases: []string{
			"mg",
			"mongo",
		},
		Short: "Execute commands to a mongodb resource",
		Long: `Use this cmd to execute mongodb commands to a mongodb object's primary pod.

Examples:
  # Execute a script named 'mongo.js' in 'mg-rs' mongodb database in 'demo' namespace
  kubectl dba exec mg mg-rs -n demo -f mongo.js

  # Execute a command in 'mg-rs' mongodb database in 'demo' namespace
  kubectl dba exec mg mg-rs -n demo -c "printjson(db.getCollectionNames())"
				`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				log.Fatal("Enter mongodb object's name as an argument. Your commands will be applied on a database inside it's primary pod")
			}
			dbName = args[0]

			opts, err := newMongodbOpts(f, dbName, namespace)
			if err != nil {
				log.Fatalln(err)
			}

			if fileName == "" && command == "" {
				log.Fatal("use --file or --command to execute supported commands to a mongodb object")
			}

			tunnel, err := lib.TunnelToDBService(opts.config, dbName, namespace, api.MongoDBDatabasePort)
			if err != nil {
				log.Fatal("couldn't create tunnel, error: ", err)
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

	mgExecCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", currentNamespace, "namespace of the mongodb object to connect to.")
	mgExecCmd.Flags().StringVarP(&fileName, "file", "f", "", "path of the script file")
	mgExecCmd.Flags().StringVarP(&command, "command", "c", "", "command to execute")

	return mgExecCmd
}

type mongodbOpts struct {
	db       *api.MongoDB
	dbImage  string
	config   *rest.Config
	client   *kubernetes.Clientset
	dbClient *cs.Clientset

	username string
	pass     string
}

func newMongodbOpts(f cmdutil.Factory, dbName, namespace string) (*mongodbOpts, error) {
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

	db, err := dbClient.KubedbV1alpha2().MongoDBs(namespace).Get(context.TODO(), dbName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	if db.Status.Phase != api.DatabasePhaseReady {
		return nil, fmt.Errorf("mongodb %s/%s is not ready", namespace, dbName)
	}

	dbVersion, err := dbClient.CatalogV1alpha1().MongoDBVersions().Get(context.TODO(), db.Spec.Version, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	secret, err := client.CoreV1().Secrets(db.Namespace).Get(context.TODO(), db.Spec.AuthSecret.Name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return &mongodbOpts{
		db:       db,
		dbImage:  dbVersion.Spec.DB.Image,
		config:   config,
		client:   client,
		dbClient: dbClient,
		username: string(secret.Data[corev1.BasicAuthUsernameKey]),
		pass:     string(secret.Data[corev1.BasicAuthPasswordKey]),
	}, nil
}

func (opts *mongodbOpts) getDockerShellCommand(localPort int, dockerFlags, mongoExtraFlags []interface{}) (*shell.Session, error) {
	sh := shell.NewSession()
	sh.ShowCMD = false

	db := opts.db
	dockerCommand := []interface{}{
		"run", "--network=host",
	}
	dockerCommand = append(dockerCommand, dockerFlags...)

	mongoCommand := []interface{}{
		"mongo", "admin",
		"--host=127.0.0.1", fmt.Sprintf("--port=%d", localPort), "--quiet",
		fmt.Sprintf("--username=%s", opts.username),
		fmt.Sprintf("--password=%s", opts.pass),
	}

	if db.Spec.TLS != nil {
		secretName := db.GetCertSecretName(api.MongoDBClientCert, "")
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

		key, ok := certSecret.Data[corev1.TLSPrivateKeyKey]
		if !ok {
			return nil, fmt.Errorf("missing %s in secret %s/%s", corev1.TLSPrivateKeyKey, certSecret.Namespace, certSecret.Name)
		}

		pem := append(crt[:], []byte("\n")...)
		pem = append(pem, key...)
		err = ioutil.WriteFile(pemFile, pem, 0644)
		if err != nil {
			return nil, err
		}

		dockerCommand = append(dockerCommand,
			"-v", fmt.Sprintf("%s:%s", caFile, caFile),
			"-v", fmt.Sprintf("%s:%s", pemFile, pemFile),
		)
		mongoCommand = append(mongoCommand,
			"--tls",
			fmt.Sprintf("--tlsCAFile=%v", caFile),
			fmt.Sprintf("--tlsCertificateKeyFile=%v", pemFile))
	}

	dockerCommand = append(dockerCommand, opts.dbImage)
	finalCommand := append(dockerCommand, mongoCommand...)
	if mongoExtraFlags != nil {
		finalCommand = append(finalCommand, mongoExtraFlags...)
	}

	return sh.Command("docker", finalCommand...).SetStdin(os.Stdin), nil
}

func (opts *mongodbOpts) connect(localPort int) error {
	dockerFlag := []interface{}{
		"-it",
	}
	shSession, err := opts.getDockerShellCommand(localPort, dockerFlag, nil)
	if err != nil {
		return err
	}

	return shSession.Run()
}

func (opts *mongodbOpts) executeCommand(localPort int, command string) error {
	mongoExtraFlags := []interface{}{
		"--eval", command,
	}
	shSession, err := opts.getDockerShellCommand(localPort, nil, mongoExtraFlags)
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

func (opts *mongodbOpts) executeFile(localPort int, fileName string) error {
	fileName, err := filepath.Abs(fileName)
	if err != nil {
		return err
	}
	tempFileName := "/home/mongo.js"

	dockerFlag := []interface{}{
		"-v", fmt.Sprintf("%s:%s", fileName, tempFileName),
	}
	mongoExtraFlags := []interface{}{
		tempFileName,
	}
	shSession, err := opts.getDockerShellCommand(localPort, dockerFlag, mongoExtraFlags)
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
