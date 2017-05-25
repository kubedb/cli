package cmd

import (
	"io"

	"github.com/spf13/cobra"
	"k8s.io/kubernetes/pkg/kubectl/cmd/templates"
)

// NewKubedbCommand creates the `kubedb` command and its nested children.
func NewKubedbCommand(in io.Reader, out, err io.Writer) *cobra.Command {

	cmds := &cobra.Command{
		Use:   "kubedb",
		Short: "Controls kubedb objects",
		Long: templates.LongDesc(`
      kubedb CLI controls kubedb ThirdPartyResource objects.

      Find more information at https://github.com/k8sdb/kubedb.`),
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
			},
		},
	}

	groups.Add(cmds)
	templates.ActsAsRootCommand(cmds, nil, groups...)

	cmds.PersistentFlags().String("kube-context", "", "name of the kubeconfig context to use")
	return cmds
}

func runHelp(cmd *cobra.Command, args []string) {
	cmd.Help()
}
