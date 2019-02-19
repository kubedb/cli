/*
Copyright 2015 The Kubernetes Authors.

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
	"fmt"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/kubernetes/pkg/kubectl/cmd/templates"
	cmdutil "k8s.io/kubernetes/pkg/kubectl/cmd/util"
	"k8s.io/kubernetes/pkg/kubectl/cmd/util/editor"
	"k8s.io/kubernetes/pkg/kubectl/util/i18n"
)

var (
	editLong = templates.LongDesc(`
		Edit a resource from the default editor.

		The edit command allows you to directly edit any API resource you can retrieve via the
		command line tools. It will open the editor defined by your KUBEDB_EDITOR, or EDITOR
		environment variables, or fall back to 'nano'`)

	editExample = templates.Examples(`
		# Edit the elasticsearch named 'elasticsearch-demo':
		kubedb edit es/elasticsearch-demo

		# Use an alternative editor
		KUBEDB_EDITOR="nano" kubedb edit es/elasticsearch-demo`)
)

func NewCmdEdit(f cmdutil.Factory, ioStreams genericclioptions.IOStreams) *cobra.Command {
	o := editor.NewEditOptions(editor.NormalEditMode, ioStreams)
	o.ValidateOptions = cmdutil.ValidateOptions{EnableValidation: true}

	cmd := &cobra.Command{
		Use:                   "edit (RESOURCE/NAME | -f FILENAME)",
		DisableFlagsInUseLine: true,
		Short:                 i18n.T("Edit a resource on the server"),
		Long:                  editLong,
		Example:               fmt.Sprintf(editExample),
		Run: func(cmd *cobra.Command, args []string) {
			if err := o.Complete(f, args, cmd); err != nil {
				cmdutil.CheckErr(err)
			}
			if err := o.Run(); err != nil {
				cmdutil.CheckErr(err)
			}
		},
	}

	// bind flag structs
	o.RecordFlags.AddFlags(cmd)
	o.PrintFlags.AddFlags(cmd)

	usage := "to use to edit the resource"
	cmdutil.AddFilenameOptionFlags(cmd, &o.FilenameOptions, usage)
	cmdutil.AddValidateOptionFlags(cmd, &o.ValidateOptions)
	cmd.Flags().BoolVarP(&o.OutputPatch, "output-patch", "", o.OutputPatch, "Output the patch if the resource is edited.")
	cmd.Flags().BoolVar(&o.WindowsLineEndings, "windows-line-endings", o.WindowsLineEndings,
		"Defaults to the line ending native to your platform.")

	cmdutil.AddApplyAnnotationVarFlags(cmd, &o.ApplyAnnotation)
	cmdutil.AddIncludeUninitializedFlag(cmd)
	return cmd
}
