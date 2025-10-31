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

package debug

import (
	"bytes"
	"context"
	"log"
	"os"
	"path"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/klog/v2"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	kmapi "kmodules.xyz/client-go/api/v1"
	dbapi "kubedb.dev/apimachinery/apis/kubedb/v1"
	kubedbscheme "kubedb.dev/apimachinery/client/clientset/versioned/scheme"
	ps "kubeops.dev/petset/client/clientset/versioned"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var scheme = runtime.NewScheme()

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(kubedbscheme.AddToScheme(scheme))
}

type dbInfo struct {
	resource  string
	name      string
	namespace string
}
type gitOpsOpts struct {
	kc client.Client
	db dbInfo

	operatorNamespace string
	dir               string
	errWriter         *bytes.Buffer
}

func GitOpsDebugCMD(f cmdutil.Factory) *cobra.Command {
	var (
		dbName            string
		operatorNamespace string
	)

	gitOpsDebugCmd := &cobra.Command{
		Use: "gitops",
		Aliases: []string{
			"git",
		},
		Short:   "Debug helper for gitops databases",
		Example: `kubectl dba debug gitops --db-type mysql -n demo sample-mysql --operator-namespace kubedb`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				log.Fatal("Enter mysql object's name as an argument")
			}
			dbName = args[0]

			namespace, _, err := f.ToRawKubeConfigLoader().Namespace()
			if err != nil {
				klog.Error(err, "failed to get current namespace")
			}

			opts, err := newGitOpsOpts(f, dbName, namespace, operatorNamespace)
			if err != nil {
				log.Fatalln(err)
			}

			err = opts.collectGitOpsDatabase()
			if err != nil {
				log.Fatal(err)
			}

			err = opts.collectForAllDBPetSets()
			if err != nil {
				log.Fatal(err)
			}

			err = opts.collectForAllDBPods()
			if err != nil {
				log.Fatal(err)
			}

			err = opts.collectOtherYamls()
			if err != nil {
				log.Fatal(err)
			}
		},
	}
	gitOpsDebugCmd.Flags().StringVarP(&operatorNamespace, "operator-namespace", "o", "kubedb", "the namespace where the kubedb gitops operator is installed")
	gitOpsDebugCmd.Flags().StringVarP(&opts, "operator-namespace", "o", "kubedb", "the namespace where the kubedb gitops operator is installed")

	return gitOpsDebugCmd
}

func newGitOpsOpts(f cmdutil.Factory, dbName, namespace, operatorNS string) (*gitOpsOpts, error) {
	config, err := f.ToRESTConfig()
	if err != nil {
		return nil, err
	}

	kc, err := client.New(config, client.Options{Scheme: scheme})
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}

	pwd, _ := os.Getwd()
	dir := path.Join(pwd, dbName)
	err = os.MkdirAll(path.Join(dir, logsDir), dirPerm)
	if err != nil {
		return nil, err
	}
	err = os.MkdirAll(path.Join(dir, yamlsDir), dirPerm)
	if err != nil {
		return nil, err
	}

	opts := &gitOpsOpts{
		db: dbInfo{
			name:      dbName,
			namespace: namespace,
		},
		kc:                kc,
		operatorNamespace: operatorNS,
		dir:               dir,
		errWriter:         &bytes.Buffer{},
	}
	return opts, nil
}

func (g *gitOpsOpts) collectGitOpsDatabase() error {
	var uns unstructured.Unstructured
	uns.SetGroupVersionKind(dbapi.SchemeGroupVersion.WithKind(g.getKindFromResource(g.db.resource)))
	err := g.kc.Get(context.Background(), types.NamespacedName{
		Namespace: g.db.namespace,
		Name:      g.db.name,
	}, &uns)
	if err != nil {
		log.Fatalf("failed to get database: %v", err)
	}

	return writeYaml(&uns, path.Join(g.dir, yamlsDir))
}
