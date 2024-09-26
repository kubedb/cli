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
	"kubedb.dev/cli/pkg/monitor"
	"kubedb.dev/cli/pkg/monitor/alerts"
	"kubedb.dev/cli/pkg/monitor/connection"
	"kubedb.dev/cli/pkg/monitor/dashboard"

	"github.com/spf13/cobra"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/util/i18n"
	"k8s.io/kubectl/pkg/util/templates"
)

var prom monitor.PromSvc

var monitorLong = templates.LongDesc(`
		All monitoring related commands from AppsCode.
    `)

var monitorExample = templates.Examples(`

		# Check triggered alerts for a specific database
		kubectl dba monitor get-alerts [DATABASE] [DATABASE_NAME] -n [NAMESPACE] 

		# Check availability of grafana dashboard of a database
		kubectl dba monitor dashboard [DATABASE] [DASHBOARD_NAME] 

		# Check connection status of target with prometheus server for a specific database
		kubectl dba monitor check-connection [DATABASE] [DATABASE_NAME] -n [NAMESPACE] 

		# Common Flags
		--prom-svc-name : name of the prometheus service
		--prom-svc-namespace : namespace of the prometheus service
		--prom-svc-port : port of the prometheus service

`)

func NewCmdMonitor(f cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "monitor",
		Short:   i18n.T("Monitoring related commands for a database"),
		Long:    monitorLong,
		Example: monitorExample,
		Run: func(cmd *cobra.Command, args []string) {
		},
		DisableFlagsInUseLine: true,
		DisableAutoGenTag:     true,
	}

	cmd.PersistentFlags().StringVarP(&prom.Name, "prom-svc-name", "", "", "name of the prometheus service")
	cmd.PersistentFlags().StringVarP(&prom.Namespace, "prom-svc-namespace", "", "", "namespace of the prometheus service")
	cmd.PersistentFlags().IntVarP(&prom.Port, "prom-svc-port", "", 9090, "port of the prometheus service")

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
		kubectl dba monitor get-alerts [DATABASE] [DATABASE_NAME] -n [NAMESPACE] \
		--prom-svc-name=[PROM_SVC_NAME] --prom-svc-namespace=[PROM_SVC_NS] --prom-svc-port=[PROM_SVC_PORT]

		# Get triggered alert for a specific mongodb
	    kubectl dba monitor get-alerts mongodb sample-mongodb -n demo \
 		--prom-svc-name=prometheus-kube-prometheus-prometheus --prom-svc-namespace=monitoring --prom-svc-port=9090
		
 		Valid resource types include:
			* connectcluster
			* druid
    		* elasticsearch
			* kafka
			* mariadb
			* mongodb
			* mssql
			* mysql
			* perconaxtradb
			* pgpool
			* postgres
			* proxysql
			* rabbitmq
			* redis
			* singlestore
			* solr
			* zookeeper
`)

func AlertCMD(f cmdutil.Factory) *cobra.Command {
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
	return cmd
}

// dashboard
var dashboardLong = templates.LongDesc(`
		Check availability of metrics in prometheus server used in a grafana dashboard.
    `)

var dashboardExample = templates.Examples(`
		kubectl dba monitor dashboard [DATABASE] [DATABASE_NAME] -n [NAMESPACE] \
		[DASHBOARD_NAME] --file=[FILE_CONTAINING_DASHBOARD_JSON] --file=[FILE_CONTAINING_REMOTE_URL] \  <- these are ORed
		--prom-svc-name=[PROM_SVC_NAME] --prom-svc-namespace=[PROM_SVC_NS] --prom-svc-port=[PROM_SVC_PORT]

		# Check availability of a postgres grafana dashboard
		kubectl-dba monitor dashboard postgres pg15 -n demo \
		--file=/home/arnob/yamls/summary.json \
		--prom-svc-name=prometheus-kube-prometheus-prometheus --prom-svc-namespace=monitoring --prom-svc-port=9090

 		Valid dashboards include:
			* connectcluster
			* druid
    		* elasticsearch
			* kafka
			* mariadb
			* mongodb
			* mssql
			* mysql
			* perconaxtradb
			* pgpool
			* postgres
			* proxysql
			* rabbitmq
			* redis
			* singlestore
			* solr
			* zookeeper
		
		If --file is given, that is the local file. absolute or relative path both accepted.
		If --url is given, that is the remote file. You have to specify the full raw url.
		If just the dashboard name is given, then that will be searched in our dashboard repo. To be exact if mongodb-summary-dashboard specified only, The cli will look for the json in
		https://raw.githubusercontent.com/appscode/grafana-dashboards/master/mongodb/mongodb-summary-dashboard.json
`)

func DashboardCMD(f cmdutil.Factory) *cobra.Command {
	var (
		branch string
		file   string
		url    string
		isDB   bool
	)
	cmd := &cobra.Command{
		Use:   "dashboard",
		Short: i18n.T("Check availability of a grafana dashboard"),
		Long:  dashboardLong,

		Run: func(cmd *cobra.Command, args []string) {
			dashboard.Run(f, args, branch, file, url, prom, isDB)
		},
		Example:               dashboardExample,
		DisableFlagsInUseLine: true,
		DisableAutoGenTag:     true,
	}
	cmd.Flags().StringVarP(&branch, "branch", "b", "master", "branch name of the github repo")
	cmd.Flags().StringVarP(&file, "file", "f", "", "absolute or relative path of the file containing dashboard")
	cmd.Flags().StringVarP(&url, "url", "u", "", "url of the raw file containing dashboard. "+
		"For example: https://raw.githubusercontent.com/appscode/grafana-dashboards/master/mongodb/mongodb-summary-dashboard.json")
	cmd.Flags().BoolVarP(&isDB, "isdb", "d", true, "for non db object's. just provide the url")
	return cmd
}

// check-connection
var connectionLong = templates.LongDesc(`
		Check connection status for different targets with prometheus server for specific DB.
`)

var connectionExample = templates.Examples(`
		kubectl dba monitor check-connection [DATABASE] [DATABASE_NAME] -n [NAMESPACE] \
		--prom-svc-name=[PROM_SVC_NAME] --prom-svc-namespace=[PROM_SVC_NS] --prom-svc-port=[PROM_SVC_PORT]

		# Check connection status for different targets with prometheus server for a specific postgres database 
		kubectl dba monitor check-connection mongodb sample_mg -n demo \
		--prom-svc-name=prometheus-kube-prometheus-prometheus --prom-svc-namespace=monitoring --prom-svc-port=9090

 		Valid resource types include:
			* connectcluster
			* druid
    		* elasticsearch
			* kafka
			* mariadb
			* mongodb
			* mssql
			* mysql
			* perconaxtradb
			* pgpool
			* postgres
			* proxysql
			* rabbitmq
			* redis
			* singlestore
			* solr
			* zookeeper
`)

func ConnectionCMD(f cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "check-connection",
		Short:   i18n.T("Check connection status of prometheus targets with server"),
		Long:    connectionLong,
		Example: connectionExample,
		Run: func(cmd *cobra.Command, args []string) {
			connection.Run(f, args, prom)
		},
		DisableFlagsInUseLine: true,
		DisableAutoGenTag:     true,
	}
	return cmd
}
