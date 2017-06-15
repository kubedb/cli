package cmd

import (
	"io"
	v "github.com/appscode/go/version"
	"github.com/k8sdb/apimachinery/pkg/analytics"
	"github.com/spf13/cobra"
	_ "github.com/spf13/cobra/doc"
	"k8s.io/kubernetes/pkg/kubectl/cmd/templates"
)

// NewKubedbCommand creates the `kubedb` command and its nested children.
func NewKubedbCommand(in io.Reader, out, err io.Writer, version string) *cobra.Command {
	var enableAnalytics bool
	cmds := &cobra.Command{
		Use:   "kubedb",
		Short: "Controls kubedb objects",
		Long: templates.LongDesc(`
      kubedb CLI controls kubedb ThirdPartyResource objects.

      Find more information at https://github.com/k8sdb/cli.`),
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if enableAnalytics {
				analytics.Enable()
			}
			analytics.SendEvent("kubedb/cli", cmd.CommandPath(), version)
		},
		Run: runHelp,
	}

	groups := templates.CommandGroups{
		{
			Message: "Basic Commands (Beginner):",
			Commands: []*cobra.Command{
				NewCmdCreate(out, err),
				NewCmdInit(out, err),
			},
		},
		{
			Message: "Basic Commands (Intermediate):",
			Commands: []*cobra.Command{
				NewCmdGet(out, err),
				NewCmdEdit(out, err),
				NewCmdDelete(out, err),
			},
		},
		{
			Message: "Troubleshooting and Debugging Commands:",
			Commands: []*cobra.Command{
				NewCmdDescribe(out, err),
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
