package dashboard

import (
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"log"
)

type DashboardOpts struct {
	Branch    string
	Database  string
	Dashboard string
}

func Run(f cmdutil.Factory, args []string, branch string) {
	if len(args) < 2 {
		log.Fatal("Enter database and grafana dashboard name as argument")
	}
	dashboardOpts := DashboardOpts{
		Branch:    branch,
		Database:  args[0],
		Dashboard: args[1],
	}
}
