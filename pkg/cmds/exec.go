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
	execLong = templates.LongDesc(`
		Execute commands or scripts to a database.
    `)
)

func NewCmdExec(f cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "exec",
		Short:                 i18n.T("Execute script or command to a database."),
		Long:                  execLong,
		Run:                   func(cmd *cobra.Command, args []string) {},
		DisableFlagsInUseLine: true,
		DisableAutoGenTag:     true,
	}

	cmd.AddCommand(connect.MongoDBExecCMD(f))
	cmd.AddCommand(connect.RedisExecCMD(f))
	cmd.AddCommand(connect.PostgresExecCMD(f))
	cmd.AddCommand(connect.MariadbExecCMD(f))
	cmd.AddCommand(connect.MySQLExecCMD(f))

	return cmd
}
