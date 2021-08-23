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
	"kubedb.dev/cli/pkg/connect"

	"github.com/spf13/cobra"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/util/i18n"
	"k8s.io/kubectl/pkg/util/templates"
)

var (
	connectLong = templates.LongDesc(`
		Connect to a database. 
    `)
)

func NewCmdConnect(f cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "connect",
		Short:                 i18n.T("Connect to a database."),
		Long:                  connectLong,
		Run:                   func(cmd *cobra.Command, args []string) {},
		DisableFlagsInUseLine: true,
		DisableAutoGenTag:     true,
	}

	cmd.AddCommand(connect.MongoDBConnectCMD(f))
	cmd.AddCommand(connect.ElasticSearchConnectCMD(f))
	cmd.AddCommand(connect.RedisConnectCMD(f))
	cmd.AddCommand(connect.PostgresConnectCMD(f))
	cmd.AddCommand(connect.MariadbConnectCMD(f))
	cmd.AddCommand(connect.MySQLConnectCMD(f))
	cmd.AddCommand(connect.MemcachedConnectCMD(f))

	return cmd
}
