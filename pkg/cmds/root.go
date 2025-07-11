/*
Copyright AppsCode Inc. and Contributors

Licensed under the AppsCode Community License 1.0.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://github.com/appscode/licenses/raw/1.0.0/AppsCode-Community-1.0.0.md

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmds

import (
	"flag"
	"fmt"
	"io"

	"github.com/spf13/cobra"
	v "gomodules.xyz/x/version"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/genericiooptions"
	cliflag "k8s.io/component-base/cli/flag"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/util/templates"
)

// NewKubeDBCommand creates the `kubedb` command and its nested children.
func NewKubeDBCommand(in io.Reader, out, err io.Writer) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "kubectl-dba",
		Short: "kubectl plugin for KubeDB",
		Long: templates.LongDesc(`
      kubectl plugin for KubeDB by AppsCode - Kubernetes ready production-grade Databases

      Find more information at https://kubedb.com`),
		Run:               runHelp,
		DisableAutoGenTag: true,
	}

	flags := rootCmd.PersistentFlags()
	flags.SetNormalizeFunc(cliflag.WarnWordSepNormalizeFunc) // Warn for "_" flags

	// Normalize all flags that are coming from other packages or pre-configurations
	// a.k.a. change all "_" to "-". e.g. glog package
	flags.SetNormalizeFunc(cliflag.WordSepNormalizeFunc)

	kubeConfigFlags := genericclioptions.NewConfigFlags(true).WithDeprecatedPasswordFlag()
	kubeConfigFlags.AddFlags(flags)
	matchVersionKubeConfigFlags := cmdutil.NewMatchVersionFlags(kubeConfigFlags)
	matchVersionKubeConfigFlags.AddFlags(rootCmd.PersistentFlags())
	rootCmd.PersistentFlags().AddGoFlagSet(flag.CommandLine)

	f := cmdutil.NewFactory(matchVersionKubeConfigFlags)

	ioStreams := genericiooptions.IOStreams{In: in, Out: out, ErrOut: err}

	groups := templates.CommandGroups{
		{
			Message: "Troubleshooting and Debugging Commands:",
			Commands: []*cobra.Command{
				NewCmdDescribe("kubedb", f, ioStreams),
				NewCmdCompletion(),
				v.NewCmdVersion(),
				NewCmdShowCredentials("kubedb", f, ioStreams),
			},
		},
		{
			Message: "Database Ops Commands:",
			Commands: []*cobra.Command{
				NewCmdRestart("kubedb", f, ioStreams),
			},
		},
		{
			Message: "Pause and Resume Commands:",
			Commands: []*cobra.Command{
				NewCmdPause("kubedb", f, ioStreams),
				NewCmdResume("kubedb", f, ioStreams),
			},
		},
		{
			Message: "Database Connection Commands",
			Commands: []*cobra.Command{
				NewCmdConnect(f),
				NewCmdExec(f),
			},
		},
		{
			Message: "Insert and Verify data in Database",
			Commands: []*cobra.Command{
				NewCmdData(f),
			},
		},
		{
			Message: "Debug any Database issue",
			Commands: []*cobra.Command{
				NewCmdDebug(f),
			},
		},
		{
			Message: "Generate appbinding and secrets for remote Replica",
			Commands: []*cobra.Command{
				NewCmdGenApb(f),
			},
		},
		{
			Message: "MSSQLServer specific commands",
			Commands: []*cobra.Command{
				NewCmdMSSQL(f),
			},
		},
		{
			Message: "Metric related CMDs",
			Commands: []*cobra.Command{
				NewCmdMonitor(f),
			},
		},
	}

	filters := []string{"options"}
	groups.Add(rootCmd)
	templates.ActsAsRootCommand(rootCmd, filters, groups...)

	rootCmd.AddCommand(NewCmdOptions(ioStreams.Out, ioStreams.ErrOut))

	return rootCmd
}

func runHelp(cmd *cobra.Command, args []string) {
	err := cmd.Help()
	if err != nil {
		fmt.Println("Failed to execute 'help' command. Reason:", err)
	}
}
