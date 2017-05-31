package main

import (
	"io"

	v "github.com/appscode/go/version"
	"github.com/k8sdb/apimachinery/pkg/analytics"
	cmd "github.com/k8sdb/kubedb/pkg/cmd"
	"github.com/spf13/cobra"
	"k8s.io/kubernetes/pkg/kubectl/cmd/templates"
)

// NewKubedbCommand creates the `kubedb` command and its nested children.
func NewKubedbCommand(in io.Reader, out, err io.Writer) *cobra.Command {
	var enableAnalytics bool
	cmds := &cobra.Command{
		Use:   "kubedb",
		Short: "Controls kubedb objects",
		Long: templates.LongDesc(`
      kubedb CLI controls kubedb ThirdPartyResource objects.

      Find more information at https://github.com/k8sdb/kubedb.`),
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if enableAnalytics {
				analytics.Enable()
			}
			analytics.SendEvent("kubedb/cli", cmd.CommandPath(), Version)
		},
		Run: runHelp,
	}

	groups := templates.CommandGroups{
		{
			Message: "Basic Commands (Beginner):",
			Commands: []*cobra.Command{
				cmd.NewCmdCreate(out, err),
				cmd.NewCmdInit(out, err),
			},
		},
		{
			Message: "Basic Commands (Intermediate):",
			Commands: []*cobra.Command{
				cmd.NewCmdGet(out, err),
				cmd.NewCmdEdit(out, err),
				cmd.NewCmdDelete(out, err),
			},
		},
		{
			Message: "Troubleshooting and Debugging Commands:",
			Commands: []*cobra.Command{
				cmd.NewCmdDescribe(out, err),
				v.NewCmdVersion(),
			},
		},
	}

	groups.Add(cmds)
	templates.ActsAsRootCommand(cmds, nil, groups...)

	cmds.PersistentFlags().String("kube-context", "", "name of the kubeconfig context to use")
	cmds.PersistentFlags().BoolVar(&enableAnalytics, "analytics", true, "Send events to Google Analytics")
	return cmds
}

func runHelp(cmd *cobra.Command, args []string) {
	cmd.Help()
}
