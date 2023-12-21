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
	"kubedb.dev/cli/pkg/dashboard"

	"github.com/spf13/cobra"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/util/i18n"
	"k8s.io/kubectl/pkg/util/templates"
)

var dashboardLong = templates.LongDesc(`
		Check availability in prometheus server of metrics used in a grafana dashboard.
    `)

var dashboardExample = templates.Examples(`
		# Check availability of mongodb-summary-dashboard grafana dashboard of mongodb
		kubectl dba dashboard mongodb mongodb-summary-dashboard

 		Valid dashboards include:
    		* elasticsearch
			* mongodb
			* mariadb
			* mysql
			* postgres
			* redis
`)

func NewCmdDashboard(f cmdutil.Factory) *cobra.Command {
	var branch string
	var prom dashboard.PromSvc
	cmd := &cobra.Command{
		Use:   "dashboard",
		Short: i18n.T("Check availability of a grafana dashboard"),
		Long:  dashboardLong,

		Run: func(cmd *cobra.Command, args []string) {
			dashboard.Run(f, args, branch, prom)
		},
		Example:               dashboardExample,
		DisableFlagsInUseLine: true,
		DisableAutoGenTag:     true,
	}
	cmd.Flags().StringVarP(&branch, "branch", "b", "master", "branch name of the github repo")
	cmd.Flags().StringVarP(&prom.Name, "prom-svc-name", "", "", "name of the prometheus service")
	cmd.Flags().StringVarP(&prom.Namespace, "prom-svc-namespace", "", "", "namespace of the prometheus service")
	cmd.Flags().IntVarP(&prom.Port, "prom-svc-port", "", 9090, "port of the prometheus service")
	return cmd
}
