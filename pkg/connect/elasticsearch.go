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

	"kubedb.dev/apimachinery/apis/kubedb"
	api "kubedb.dev/apimachinery/apis/kubedb/v1"
	cs "kubedb.dev/apimachinery/client/clientset/versioned"
	"kubedb.dev/cli/pkg/lib"

	"github.com/spf13/cobra"
	shell "gomodules.xyz/go-sh"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
)

func ElasticSearchConnectCMD(f cmdutil.Factory) *cobra.Command {
	var dbName string

	esConnectCmd := &cobra.Command{
		Use: "elasticsearch",
		Aliases: []string{
			"es",
		},
		Short: "Connect to a shell to run elasticsearch api calls",
		Long: `Use this cmd to run api calls to your elasticsearch database. 

This command connects you to a shell to run curl commands. 

It exports the following environment variables to run api calls to your database:
  $USERNAME
  $PASSWORD
  $ADDRESS
  $CACERT
  $CERT
  $KEY

Example connect command:
  # connect to a shell with curl access to the database of name es-demo in demo namespace
  kubectl dba connect es es-demo -n demo

Example curl commands:
  # curl command to run on the connected elasticsearch database:
  curl -u $USERNAME:$PASSWORD $ADDRESS/_cluster/health?pretty

  # curl command to run on the connected tls secured elasticsearch database:
  curl --cacert $CACERT --cert $CERT --key $KEY  -u $USERNAME:$PASSWORD $ADDRESS/_cluster/health?pretty`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				log.Fatal("enter elasticsearch object's name as an argument")
			}
			dbName = args[0]

			namespace, _, err := f.ToRawKubeConfigLoader().Namespace()
			if err != nil {
				klog.Error(err, "failed to get current namespace")
			}
			opts, err := newElasticsearchOpts(f, dbName, namespace)
			if err != nil {
				log.Fatalln(err)
			}

			tunnel, err := lib.TunnelToDBService(opts.config, dbName, namespace, kubedb.ElasticsearchRestPort)
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

	return esConnectCmd
}

type elasticsearchOpts struct {
	db       *api.Elasticsearch
	config   *rest.Config
	client   *kubernetes.Clientset
	dbClient *cs.Clientset

	username string
	pass     string
}

func newElasticsearchOpts(f cmdutil.Factory, dbName, namespace string) (*elasticsearchOpts, error) {
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

	db, err := dbClient.KubedbV1().Elasticsearches(namespace).Get(context.TODO(), dbName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	if db.Status.Phase != api.DatabasePhaseReady {
		return nil, fmt.Errorf("elasticsearch %s/%s is not ready", namespace, dbName)
	}

	secret, err := client.CoreV1().Secrets(db.Namespace).Get(context.TODO(), db.Spec.AuthSecret.Name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return &elasticsearchOpts{
		db:       db,
		config:   config,
		client:   client,
		dbClient: dbClient,
		username: string(secret.Data[v1.BasicAuthUsernameKey]),
		pass:     string(secret.Data[v1.BasicAuthPasswordKey]),
	}, nil
}

func (opts *elasticsearchOpts) getDockerShellCommand(localPort int, dockerFlags, esCommands []interface{}) (*shell.Session, error) {
	sh := shell.NewSession()
	sh.ShowCMD = false

	db := opts.db
	dockerCommand := []interface{}{
		"run", "--network=host", "-it",
		"-e", fmt.Sprintf("USERNAME=%s", opts.username),
		"-e", fmt.Sprintf("PASSWORD=%s", opts.pass),
	}
	dockerCommand = append(dockerCommand, dockerFlags...)

	if db.Spec.EnableSSL {
		secretName := db.CertificateName(api.ElasticsearchAdminCert)
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
			"-e", fmt.Sprintf("ADDRESS=https://localhost:%d", localPort),
			"-e", fmt.Sprintf("CACERT=%s", caFile),
			"-e", fmt.Sprintf("CERT=%s", certFile),
			"-e", fmt.Sprintf("KEY=%s", keyFile),
			"-v", fmt.Sprintf("%s:%s", caFile, caFile),
			"-v", fmt.Sprintf("%s:%s", certFile, certFile),
			"-v", fmt.Sprintf("%s:%s", keyFile, keyFile),
		)
	} else {
		dockerCommand = append(dockerCommand,
			"-e", fmt.Sprintf("ADDRESS=http://localhost:%d", localPort))
	}

	dockerCommand = append(dockerCommand, alpineCurlImg)
	finalCommand := dockerCommand
	if esCommands != nil {
		finalCommand = append(finalCommand, esCommands...)
	}
	return sh.Command("docker", finalCommand...).SetStdin(os.Stdin), nil
}

func (opts *elasticsearchOpts) connect(localPort int) error {
	shellCmd, err := opts.getDockerShellCommand(localPort, []interface{}{"-it"}, []interface{}{"sh"})
	if err != nil {
		return err
	}
	return shellCmd.Run()
}
