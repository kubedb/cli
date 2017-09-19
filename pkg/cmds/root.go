package cmds

import (
	"io"

	v "github.com/appscode/go/version"
	"github.com/spf13/cobra"
	_ "github.com/spf13/cobra/doc"
	"k8s.io/kubernetes/pkg/kubectl/cmd/templates"
	"github.com/jpillora/go-ogle-analytics"
	"strings"
)

const (
	gaTrackingCode = "UA-62096468-20"
)

// NewKubedbCommand creates the `kubedb` command and its nested children.
func NewKubedbCommand(in io.Reader, out, err io.Writer, version string) *cobra.Command {
	enableAnalytics := true
	cmds := &cobra.Command{
		Use:   "kubedb",
		Short: "Command line interface for KubeDB",
		Long: templates.LongDesc(`
      KubeDB by AppsCode - Kubernetes ready production-grade Databases

      Find more information at https://github.com/k8sdb/cli.`),
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if enableAnalytics && gaTrackingCode != "" {
				if client, err := ga.NewClient(gaTrackingCode); err == nil {
					parts := strings.Split(cmd.CommandPath(), " ")
					client.Send(ga.NewEvent(parts[0], strings.Join(parts[1:], "/")).Label(version))
				}
			}
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
				NewCmdSummarize(out, err),
				NewCmdCompare(out, err),
				v.NewCmdVersion(),
			},
		},
	}

	groups.Add(cmds)
	templates.ActsAsRootCommand(cmds, nil, groups...)

	cmds.PersistentFlags().String("kube-context", "", "name of the kubeconfig context to use")
	cmds.PersistentFlags().BoolVar(&enableAnalytics, "analytics", enableAnalytics, "Send analytical events to Google Analytics")

	return cmds
}

func runHelp(cmd *cobra.Command, args []string) {
	cmd.Help()
}
