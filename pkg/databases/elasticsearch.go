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

package databases

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	apiv1alpha2 "kubedb.dev/apimachinery/apis/kubedb/v1alpha2"
	cs "kubedb.dev/apimachinery/client/clientset/versioned"
	"kubedb.dev/cli/pkg/lib"

	shell "github.com/codeskyblue/go-sh"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
)

func NewElasticSearchCMD(f cmdutil.Factory) *cobra.Command {
	var (
		dbName    string
		namespace string
	)

	currentNamespace, _, err := f.ToRawKubeConfigLoader().Namespace()
	if err != nil {
		klog.Error(err, "failed to get current namespace")
	}

	var esCmd = &cobra.Command{
		Use: "elasticsearch",
		Aliases: []string{
			"es",
		},
		Short: "Use to operate elasticsearch pods",
		Long: `Use this cmd to operate elasticsearch pods. Available sub-commands:
				connect`,
		Run: func(cmd *cobra.Command, args []string) {},
	}

	var esConnectCmd = &cobra.Command{
		Use:   "connect",
		Short: "Connect to a elasticsearch object's pod",
		Long:  `Use this cmd to exec into a elasticsearch object's primary pod.`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				log.Fatal("enter elasticsearch object's name as an argument")
			}
			dbName = args[0]
			opts, err := newElasticsearchOpts(f, dbName, namespace)
			if err != nil {
				log.Fatalln(err)
			}

			tunnel, err := lib.TunnelToDBService(opts.config, dbName, namespace, apiv1alpha2.ElasticsearchRestPort)
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

	esCmd.AddCommand(esConnectCmd)
	esCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", currentNamespace, "namespace of the elasticsearch object to connect to.")

	return esCmd
}

type elasticsearchOpts struct {
	db       *apiv1alpha2.Elasticsearch
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

	db, err := dbClient.KubedbV1alpha2().Elasticsearches(namespace).Get(context.TODO(), dbName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	if db.Status.Phase != apiv1alpha2.DatabasePhaseReady {
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
		secretName := db.CertificateName(apiv1alpha2.ElasticsearchAdminCert)
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
