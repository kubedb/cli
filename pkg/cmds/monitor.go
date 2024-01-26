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
	"kubedb.dev/cli/pkg/connection"
	"kubedb.dev/cli/pkg/dashboard"

	"github.com/spf13/cobra"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/util/i18n"
	"k8s.io/kubectl/pkg/util/templates"
)

var monitorLong = templates.LongDesc(`
		All monitoring related commands from AppsCode.
    `)

var monitorExample = templates.Examples(`

		# Check triggered alerts for a specific database
		kubectl dba monitor get-alerts mongodb -n demo sample-mongodb --prom-svc-name=prometheus-kube-prometheus-prometheus --prom-svc-namespace=monitoring

		# Check availability of mongodb-summary-dashboard grafana dashboard of mongodb
		kubectl dba monitor dashboard mongodb mongodb-summary-dashboard

		# Check connection status of target with prometheus server
		kubectl dba monitor check-connection mongodb

 		Valid sub command include:
    		* get-alerts
			* dashboard
			* check-connection
`)

func NewCmdMonitor(f cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "monitor",
		Short:                 i18n.T("get-alerts,grafana dashboard check and target connection check for monitoring"),
		Long:                  monitorLong,
		Example:               monitorExample,
		Run:                   func(cmd *cobra.Command, args []string) {},
		DisableFlagsInUseLine: true,
		DisableAutoGenTag:     true,
	}

	cmd.AddCommand(DashboardCMD(f))
	cmd.AddCommand(AlertCMD(f))
	cmd.AddCommand(ConnectionCMD(f))

	return cmd
}

// alert
var alertLong = templates.LongDesc(`
		Get the prometheus alerts for a specific database in just one command
    `)
var alertExample = templates.Examples(`
	    kubectl dba monitor get-alerts mongodb -n demo sample-mongodb --prom-svc-name=prometheus-kube-prometheus-prometheus --prom-svc-namespace=monitoring
		
 		Valid resource types include:
    		* elasticsearch
			* mongodb
			* mariadb
			* mysql
			* postgres
			* redis
`)

func AlertCMD(f cmdutil.Factory) *cobra.Command {
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

// dashboard
var dashboardLong = templates.LongDesc(`
		Check availability of metrics in prometheus server used in a grafana dashboard.
    `)

var dashboardExample = templates.Examples(`
		# Check availability of mongodb-summary-dashboard grafana dashboard of mongodb
		kubectl dba monitor dashboard mongodb mongodb-summary-dashboard

 		Valid dashboards include:
    		* elasticsearch
			* mongodb
			* mariadb
			* mysql
			* postgres
			* redis
`)

func DashboardCMD(f cmdutil.Factory) *cobra.Command {
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

// check-connection
// TODO
var connectionLong = templates.LongDesc(`
		
    `)

// TODO
var connectionExample = templates.Examples(`
# kubectl dba monitor check-connection mongodb -n demo sample_mg -> all connection check and report
`)

func ConnectionCMD(f cmdutil.Factory) *cobra.Command {
	var prom connection.PromSvc
	cmd := &cobra.Command{
		Use:     "check-connection",
		Short:   i18n.T("Check connection status of prometheus target with server"),
		Long:    connectionLong,
		Example: connectionExample,
		Run: func(cmd *cobra.Command, args []string) {
			connection.Run(f, args, prom)
		},
		DisableFlagsInUseLine: true,
		DisableAutoGenTag:     true,
	}
	cmd.Flags().StringVarP(&prom.Name, "prom-svc-name", "", "", "name of the prometheus service")
	cmd.Flags().StringVarP(&prom.Namespace, "prom-svc-namespace", "", "", "namespace of the prometheus service")
	cmd.Flags().IntVarP(&prom.Port, "prom-svc-port", "", 9090, "port of the prometheus service")
	return cmd
}
