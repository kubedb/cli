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
	"kubedb.dev/cli/pkg/remote_replica"

	"github.com/spf13/cobra"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
)

var (
	desLong = "generate appbinding , secrets for remote replica"
	example = "kubectl dba remote-config mysql -n <ns> -u <user_name> -p$<password> -d<dns_name>  <db_name>"
)

func NewCmdGenApb(f cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "remote-config",
		Short:   desLong,
		Long:    desLong,
		Example: example,
		Run: func(cmd *cobra.Command, args []string) {
		},
		DisableAutoGenTag:     false,
		DisableFlagsInUseLine: false,
	}
	cmd.AddCommand(remote_replica.MysqlAPP(f))
	cmd.AddCommand(remote_replica.PostgreSQlAPP(f))
	return cmd
}
