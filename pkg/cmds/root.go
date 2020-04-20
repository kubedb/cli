/*
Copyright The KubeDB Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmds

import (
	"flag"
	"io"

	"kubedb.dev/cli/pkg/cmds/create"
	"kubedb.dev/cli/pkg/cmds/get"

	v "github.com/appscode/go/version"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	cliflag "k8s.io/component-base/cli/flag"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/util/templates"
	"kmodules.xyz/client-go/logs"
	"kmodules.xyz/client-go/tools/cli"
)

// NewKubeDBCommand creates the `kubedb` command and its nested children.
func NewKubeDBCommand(in io.Reader, out, err io.Writer) *cobra.Command {
	cmds := &cobra.Command{
		Use:   "kubectl-dba",
		Short: "kubectl plugin for KubeDB",
		Long: templates.LongDesc(`
      KubeDB by AppsCode - Kubernetes ready production-grade Databases

      Find more information at https://kubedb.com`),
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			cli.SendAnalytics(cmd, v.Version.Version)
		},
		Run: runHelp,
	}

	flags := cmds.PersistentFlags()
	// Normalize all flags that are coming from other packages or pre-configurations
	// a.k.a. change all "_" to "-". e.g. glog package
	flags.SetNormalizeFunc(cliflag.WordSepNormalizeFunc)

	kubeConfigFlags := genericclioptions.NewConfigFlags(true)
	kubeConfigFlags.AddFlags(flags)
	matchVersionKubeConfigFlags := cmdutil.NewMatchVersionFlags(kubeConfigFlags)
	matchVersionKubeConfigFlags.AddFlags(flags)

	flags.AddGoFlagSet(flag.CommandLine)
	logs.ParseFlags()
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
