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
	"kubedb.dev/cli/pkg/alerts"

	"github.com/spf13/cobra"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/util/i18n"
	"k8s.io/kubectl/pkg/util/templates"
)

var (
	alertLong = templates.LongDesc(`
		Get the prometheus alerts for a specific database in just one command
    `)
	alertExample = templates.Examples(`
	    kubectl dba get-alerts mongodb -n demo sample-mongodb --prom-svc-name=prometheus-kube-prometheus-prometheus --prom-svc-namespace=monitoring
		
 		Valid resource types include:
    		* elasticsearch
			* mongodb
			* mariadb
			* mysql
			* postgres
			* redis
`)
)

func NewCmdAlert(f cmdutil.Factory) *cobra.Command {
	var prom alerts.PromSvc
	cmd := &cobra.Command{
		Use:     "get-alerts",
		Short:   i18n.T("Alerts associated with a database"),
		Long:    alertLong,
		Example: alertExample,
		Run: func(cmd *cobra.Command, args []string) {
			alerts.Run(f, args, prom)
		},
		DisableFlagsInUseLine: true,
		DisableAutoGenTag:     true,
	}
	cmd.Flags().StringVarP(&prom.Name, "prom-svc-name", "", "", "name of the prometheus service")
	cmd.Flags().StringVarP(&prom.Namespace, "prom-svc-namespace", "", "", "namespace of the prometheus service")
	cmd.Flags().IntVarP(&prom.Port, "prom-svc-port", "", 9090, "port of the prometheus service")
	return cmd
}
