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
	"context"
	"log"

	olddbapi "kubedb.dev/apimachinery/apis/kubedb/v1alpha2"

	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/klog/v2"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
)

func DruidDebugCMD(f cmdutil.Factory) *cobra.Command {
	var (
		dbName            string
		operatorNamespace string
	)

	mdDebugCmd := &cobra.Command{
		Use: "druid",
		Aliases: []string{
			"dr",
			"druids",
		},
		Short:   "Debug helper for Druid database",
		Example: `kubectl dba debug druid -n demo sample-druid --operator-namespace kubedb`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				log.Fatal("Enter druid object's name as an argument")
			}
			dbName = args[0]

			namespace, _, err := f.ToRawKubeConfigLoader().Namespace()
			if err != nil {
				klog.Error(err, "failed to get current namespace")
			}

			gvk := func() schema.GroupVersionKind {
				kind := olddbapi.ResourceKindDruid
				return schema.GroupVersionKind{
					Group:   olddbapi.SchemeGroupVersion.Group,
					Version: olddbapi.SchemeGroupVersion.Version,
					Kind:    kind,
				}
			}()
			opts, err := newDBOpts(f, gvk, dbName, namespace, operatorNamespace)
			if err != nil {
				log.Fatalln(err)
			}

			var db olddbapi.Druid
			err = opts.kc.Get(context.TODO(), types.NamespacedName{Name: dbName, Namespace: namespace}, &db)
			if err != nil {
				log.Fatalln(err)
			}
			opts.db.OwnerReferences = db.OwnerReferences

			err = writeYaml(&db, getDir(db.GetName()))
			if err != nil {
				return
			}
			opts.selectors = db.OffshootSelectors()
			klog.Infof("db selectors: %v;\nDebug info has been generated in '%v' folder", opts.selectors, dbName)
			err = opts.collectALl()
			if err != nil {
				log.Fatalln(err)
			}
		},
	}
	mdDebugCmd.Flags().StringVarP(&operatorNamespace, "operator-namespace", "o", "kubedb", "the namespace where the kubedb operator is installed")

	return mdDebugCmd
}
