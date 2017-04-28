package cmd

import (
	"io"

	"github.com/k8sdb/kubedb/pkg/cmd/templates"
	"github.com/spf13/cobra"
)

// NewKubedbCommand creates the `kubedb` command and its nested children.
func NewKubedbCommand(in io.Reader, out, err io.Writer) *cobra.Command {
	cmds := &cobra.Command{
		Use: "kubedb",
	}

	groups := templates.CommandGroups{
		{
			Message:  "Basic Commands (Beginner):",
			Commands: []*cobra.Command{},
		},
		{
			Message: "Basic Commands (Intermediate):",
			Commands: []*cobra.Command{
				NewCmdGet(out, err),
			},
		},
	}

	groups.Add(cmds)
	return cmds
}
