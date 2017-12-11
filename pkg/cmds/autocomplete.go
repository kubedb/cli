package cmds

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
	cmdutil "k8s.io/kubernetes/pkg/kubectl/cmd/util"
)

func NewCmdAutoComplete(out io.Writer, errOut io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "autocomplete",
		Short: "Generate bash autocompletions script",
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(RunAutoComplete(cmd, out))
		},
	}
	return cmd
}

const bashCompletionDir = "/etc/bash_completion.d"

func RunAutoComplete(cmd *cobra.Command, out io.Writer) error {
	fileName := fmt.Sprintf("%v/kubedb.sh", bashCompletionDir)
	return cmd.Root().GenBashCompletionFile(fileName)
}
