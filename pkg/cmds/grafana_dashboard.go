package cmds

import (
	"github.com/spf13/cobra"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/util/i18n"
	"k8s.io/kubectl/pkg/util/templates"
	"kubedb.dev/cli/pkg/dashboard"
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

func NewCmdDashboardCMD(f cmdutil.Factory) *cobra.Command {
	var branch string
	cmd := &cobra.Command{
		Use:   "dashboard",
		Short: i18n.T("Check availability of a grafana dashboard"),
		Long:  dashboardLong,

		Run: func(cmd *cobra.Command, args []string) {
			dashboard.Run(f, args, branch)
		},
		Example:               dashboardExample,
		DisableFlagsInUseLine: true,
		DisableAutoGenTag:     true,
	}
	cmd.Flags().StringVarP(&branch, "branch", "b", "master", "branch name of the github repo")
	return cmd
}
