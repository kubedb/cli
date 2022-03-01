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
	"strconv"

	api "kubedb.dev/apimachinery/apis/kubedb/v1alpha2"
	cs "kubedb.dev/apimachinery/client/clientset/versioned"
	"kubedb.dev/cli/pkg/lib"

	"github.com/spf13/cobra"
	shell "gomodules.xyz/go-sh"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
)

func MemcachedConnectCMD(f cmdutil.Factory) *cobra.Command {
	var dbName string

	mcConnectCmd := &cobra.Command{
		Use: "memcached",
		Aliases: []string{
			"mc",
		},
		Short: "Connect to a telnet shell to run command against a memcached database",
		Long:  `Use this cmd to exec into a memcached object's primary pod.`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				log.Fatal("enter memcached object's name as an argument")
			}
			dbName = args[0]

			namespace, _, err := f.ToRawKubeConfigLoader().Namespace()
			if err != nil {
				klog.Error(err, "failed to get current namespace")
			}
			opts, err := newMemcachedOpts(f, dbName, namespace)
			if err != nil {
				log.Fatalln(err)
			}

			tunnel, err := lib.TunnelToDBService(opts.config, dbName, namespace, api.MemcachedDatabasePort)
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

	return mcConnectCmd
}

type memcachedOpts struct {
	db       *api.Memcached
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

	if db.Status.Phase != api.DatabasePhaseReady {
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
	return sh.Command("docker", "run", "--network=host", "-it", "--entrypoint", "telnet",
		busyboxImg, "127.0.0.1", strconv.Itoa(localPort),
	).SetStdin(os.Stdin).Run()
}
