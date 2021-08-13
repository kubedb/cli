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
	"log"
	"os"
	"strconv"

	apiv1alpha2 "kubedb.dev/apimachinery/apis/kubedb/v1alpha2"
	cs "kubedb.dev/apimachinery/client/clientset/versioned"
	"kubedb.dev/cli/pkg/lib"

	shell "github.com/codeskyblue/go-sh"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
)

func NewMemcachedCMD(f cmdutil.Factory) *cobra.Command {
	var (
		dbName    string
		namespace string
	)

	currentNamespace, _, err := f.ToRawKubeConfigLoader().Namespace()
	if err != nil {
		klog.Error(err, "failed to get current namespace")
	}

	var mcCmd = &cobra.Command{
		Use: "memcached",
		Aliases: []string{
			"mc",
		},
		Short: "Use to operate memcached pods",
		Long: `Use this cmd to operate memcached pods. Available sub-commands:
				connect`,
		Run: func(cmd *cobra.Command, args []string) {},
	}

	var mcConnectCmd = &cobra.Command{
		Use:   "connect",
		Short: "Connect to a memcached object's pod",
		Long:  `Use this cmd to exec into a memcached object's primary pod.`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				log.Fatal("enter memcached object's name as an argument")
			}
			dbName = args[0]
			opts, err := newMemcachedOpts(f, dbName, namespace)
			if err != nil {
				log.Fatalln(err)
			}

			tunnel, err := lib.TunnelToDBService(opts.config, dbName, namespace, apiv1alpha2.MemcachedDatabasePort)
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

	mcCmd.AddCommand(mcConnectCmd)
	mcCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", currentNamespace, "namespace of the memcached object to connect to.")

	return mcCmd
}

type memcachedOpts struct {
	db       *apiv1alpha2.Memcached
	config   *rest.Config
	client   *kubernetes.Clientset
	dbClient *cs.Clientset
}

func newMemcachedOpts(f cmdutil.Factory, dbName, namespace string) (*memcachedOpts, error) {
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

	db, err := dbClient.KubedbV1alpha2().Memcacheds(namespace).Get(context.TODO(), dbName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	if db.Status.Phase != apiv1alpha2.DatabasePhaseReady {
		return nil, fmt.Errorf("memcached %s/%s is not ready", namespace, dbName)
	}

	return &memcachedOpts{
		db:       db,
		config:   config,
		client:   client,
		dbClient: dbClient,
	}, nil
}

func (opts *memcachedOpts) connect(localPort int) error {
	sh := shell.NewSession()
	return sh.Command("docker", "run", "--network=host", "-it",
		alpineTelnetImg, "127.0.0.1", strconv.Itoa(localPort),
	).SetStdin(os.Stdin).Run()
}
