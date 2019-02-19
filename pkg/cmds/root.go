package cmds

import (
	"flag"
	"io"
	v "github.com/appscode/go/version"
	"github.com/appscode/kutil/tools/cli"
	"github.com/kubedb/cli/pkg/cmds/create"
	"github.com/kubedb/cli/pkg/cmds/get"
	"github.com/spf13/cobra"
	utilflag "k8s.io/apiserver/pkg/util/flag"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/kubernetes/pkg/kubectl/cmd/templates"
	cmdutil "k8s.io/kubernetes/pkg/kubectl/cmd/util"
)

// NewKubedbCommand creates the `kubedb` command and its nested children.
func NewKubedbCommand(in io.Reader, out, err io.Writer) *cobra.Command {
	cmds := &cobra.Command{
		Use:   "kubedb",
		Short: "Command line interface for KubeDB",
		Long: templates.LongDesc(`
      KubeDB by AppsCode - Kubernetes ready production-grade Databases

      Find more information at https://github.com/kubedb/cli.`),
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			cli.SendAnalytics(cmd, v.Version.Version)
		},
		Run: runHelp,
	}

	flags := cmds.PersistentFlags()
	// Normalize all flags that are coming from other packages or pre-configurations
	// a.k.a. change all "_" to "-". e.g. glog package
	flags.SetNormalizeFunc(utilflag.WordSepNormalizeFunc)

	kubeConfigFlags := genericclioptions.NewConfigFlags()
	kubeConfigFlags.AddFlags(flags)
	matchVersionKubeConfigFlags := cmdutil.NewMatchVersionFlags(kubeConfigFlags)
	matchVersionKubeConfigFlags.AddFlags(flags)

	flags.AddGoFlagSet(flag.CommandLine)
	// ref: https://github.com/kubernetes/kubernetes/issues/17162#issuecomment-225596212
	flag.CommandLine.Parse([]string{})
	flags.BoolVar(&cli.EnableAnalytics, "enable-analytics", cli.EnableAnalytics, "Send analytical events to Google Analytics")

	f := cmdutil.NewFactory(matchVersionKubeConfigFlags)

	ioStreams := genericclioptions.IOStreams{In: in, Out: out, ErrOut: err}

	groups := templates.CommandGroups{
		{
			Message: "Basic Commands (Beginner):",
			Commands: []*cobra.Command{
				create.NewCmdCreate(f, ioStreams),
			},
		},
		{
			Message: "Basic Commands (Intermediate):",
			Commands: []*cobra.Command{
				get.NewCmdGet("kubedb", f, ioStreams),
				NewCmdEdit(f, ioStreams),
				NewCmdDelete(f, ioStreams),
			},
		},
		{
			Message: "Troubleshooting and Debugging Commands:",
			Commands: []*cobra.Command{
				NewCmdDescribe("kubedb", f, ioStreams),
				NewCmdApiResources(f, ioStreams),
				v.NewCmdVersion(),
			},
		},
	}
	groups.Add(cmds)
	templates.ActsAsRootCommand(cmds, nil, groups...)

	return cmds
}

func runHelp(cmd *cobra.Command, args []string) {
	cmd.Help()
}
