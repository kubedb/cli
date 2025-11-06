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

package cmds

import (
	"kubedb.dev/cli/pkg/debug"

	"github.com/spf13/cobra"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/util/i18n"
	"k8s.io/kubectl/pkg/util/templates"
)

var (
	debugLong = templates.LongDesc(`
		Collect all the logs and yamls of a specific database in just one command
    `)
	debugExample = templates.Examples(`
	    kubectl dba debug mongodb -n demo sample-mongodb --operator-namespace kubedb
		
 		Valid resource types include:
    		* elasticsearch
			* mongodb
			* mariadb
			* mysql
			* postgres
			* redis
`)
)

func NewCmdDebug(f cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "debug",
		Short:                 i18n.T("Debug any Database issue"),
		Long:                  debugLong,
		Example:               debugExample,
		Run:                   func(cmd *cobra.Command, args []string) {},
		DisableFlagsInUseLine: true,
		DisableAutoGenTag:     true,
	}

	cmd.AddCommand(debug.ElasticsearchDebugCMD(f))
	cmd.AddCommand(debug.MariaDebugCMD(f))
	cmd.AddCommand(debug.MongoDBDebugCMD(f))
	cmd.AddCommand(debug.MySQLDebugCMD(f))
	cmd.AddCommand(debug.PostgresDebugCMD(f))
	cmd.AddCommand(debug.RedisDebugCMD(f))

	return cmd
}
