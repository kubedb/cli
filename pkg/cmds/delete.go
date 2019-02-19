/*
Copyright 2014 The Kubernetes Authors.

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
	"strings"
	"time"
	"github.com/golang/glog"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/genericclioptions/printers"
	"k8s.io/cli-runtime/pkg/genericclioptions/resource"
	"k8s.io/client-go/dynamic"
	"k8s.io/kubernetes/pkg/kubectl/cmd/templates"
	cmdutil "k8s.io/kubernetes/pkg/kubectl/cmd/util"
	kubectlwait "k8s.io/kubernetes/pkg/kubectl/cmd/wait"
	"k8s.io/kubernetes/pkg/kubectl/util/i18n"
)

var (
	delete_long = templates.LongDesc(`
		Delete resources by filenames, stdin, resources and names, or by resources and label selector.
		JSON and YAML formats are accepted.

		Note that the delete command does NOT do resource version checks`)

	delete_example = templates.Examples(`
		# Delete a elasticsearch using the type and name specified in elastic.json.
		kubedb delete -f ./elastic.json

		# Delete a postgres based on the type and name in the JSON passed into stdin.
		cat postgres.json | kubedb delete -f -

		# Delete elasticsearch with label elasticsearch.kubedb.com/name=elasticsearch-demo.
		kubedb delete elasticsearch -l elasticsearch.kubedb.com/name=elasticsearch-demo

		# Delete all mysql objects
		kubedb delete mysql --all`)
)

type DeleteOptions struct {
	resource.FilenameOptions

	LabelSelector   string
	FieldSelector   string
	DeleteAll       bool
	IgnoreNotFound  bool
	Cascade         bool
	DeleteNow       bool
	ForceDeletion   bool
	WaitForDeletion bool

	GracePeriod int
	Timeout     time.Duration

	Output string

	DynamicClient dynamic.Interface
	Mapper        meta.RESTMapper
	Result        *resource.Result

	genericclioptions.IOStreams
}

func NewCmdDelete(f cmdutil.Factory, streams genericclioptions.IOStreams) *cobra.Command {
	deleteFlags := NewDeleteCommandFlags("containing the resource to delete.")

	cmd := &cobra.Command{
		Use:                   "delete ([-f FILENAME] | TYPE [(NAME | -l label | --all)])",
		DisableFlagsInUseLine: true,
		Short:                 i18n.T("Delete resources by filenames, stdin, resources and names, or by resources and label selector"),
		Long:                  delete_long,
		Example:               delete_example,
		Run: func(cmd *cobra.Command, args []string) {
			o := deleteFlags.ToOptions(nil, streams)
			cmdutil.CheckErr(o.Complete(f, args, cmd))
			cmdutil.CheckErr(o.Validate(cmd))
			cmdutil.CheckErr(o.RunDelete())
		},
		SuggestFor: []string{"rm"},
	}

	deleteFlags.AddFlags(cmd)

	cmdutil.AddIncludeUninitializedFlag(cmd)
	return cmd
}

func (o *DeleteOptions) Complete(f cmdutil.Factory, args []string, cmd *cobra.Command) error {
	cmdNamespace, enforceNamespace, err := f.ToRawKubeConfigLoader().Namespace()
	if err != nil {
		return err
	}

	if o.DeleteAll || len(o.LabelSelector) > 0 || len(o.FieldSelector) > 0 {
		if f := cmd.Flags().Lookup("ignore-not-found"); f != nil && !f.Changed {
			// If the user didn't explicitly set the option, default to ignoring NotFound errors when used with --all, -l, or --field-selector
			o.IgnoreNotFound = true
		}
	}
	if o.DeleteNow {
		if o.GracePeriod != -1 {
			return fmt.Errorf("--now and --grace-period cannot be specified together")
		}
		o.GracePeriod = 1
	}
	if o.GracePeriod == 0 && !o.ForceDeletion {
		// To preserve backwards compatibility, but prevent accidental data loss, we convert --grace-period=0
		// into --grace-period=1. Users may provide --force to bypass this conversion.
		o.GracePeriod = 1
	}

	includeUninitialized := cmdutil.ShouldIncludeUninitialized(cmd, false)
	r := f.NewBuilder().
		Unstructured().
		ContinueOnError().
		NamespaceParam(cmdNamespace).DefaultNamespace().
		FilenameParam(enforceNamespace, &o.FilenameOptions).
		LabelSelectorParam(o.LabelSelector).
		FieldSelectorParam(o.FieldSelector).
		IncludeUninitialized(includeUninitialized).
		SelectAllParam(o.DeleteAll).
		ResourceTypeOrNameArgs(false, args...).RequireObject(false).
		Flatten().
		Do()
	err = r.Err()
	if err != nil {
		return err
	}
	o.Result = r

	o.Mapper, err = f.ToRESTMapper()
	if err != nil {
		return err
	}

	o.DynamicClient, err = f.DynamicClient()
	if err != nil {
		return err
	}

	return nil
}

func (o *DeleteOptions) Validate(cmd *cobra.Command) error {
	if o.Output != "" && o.Output != "name" {
		return cmdutil.UsageErrorf(cmd, "Unexpected -o output mode: %v. We only support '-o name'.", o.Output)
	}

	if o.DeleteAll && len(o.LabelSelector) > 0 {
		return fmt.Errorf("cannot set --all and --selector at the same time")
	}
	if o.DeleteAll && len(o.FieldSelector) > 0 {
		return fmt.Errorf("cannot set --all and --field-selector at the same time")
	}

	if o.GracePeriod == 0 && !o.ForceDeletion && !o.WaitForDeletion {
		// With the explicit --wait flag we need extra validation for backward compatibility
		return fmt.Errorf("--grace-period=0 must have either --force specified, or --wait to be set to true")
	}

	switch {
	case o.GracePeriod == 0 && o.ForceDeletion:
		fmt.Fprintf(o.ErrOut, "warning: Immediate deletion does not wait for confirmation that the running resource has been terminated. The resource may continue to run on the cluster indefinitely.\n")
	case o.ForceDeletion:
		fmt.Fprintf(o.ErrOut, "warning: --force is ignored because --grace-period is not 0.\n")
	}
	return nil
}

func (o *DeleteOptions) RunDelete() error {
	return o.DeleteResult(o.Result)
}

func (o *DeleteOptions) DeleteResult(r *resource.Result) error {
	found := 0
	if o.IgnoreNotFound {
		r = r.IgnoreErrors(errors.IsNotFound)
	}
	deletedInfos := []*resource.Info{}
	uidMap := kubectlwait.UIDMap{}
	err := r.Visit(func(info *resource.Info, err error) error {
		if err != nil {
			return err
		}
		deletedInfos = append(deletedInfos, info)
		found++

		options := &metav1.DeleteOptions{}
		if o.GracePeriod >= 0 {
			options = metav1.NewDeleteOptions(int64(o.GracePeriod))
		}
		policy := metav1.DeletePropagationBackground
		if !o.Cascade {
			policy = metav1.DeletePropagationOrphan
		}
		options.PropagationPolicy = &policy

		response, err := o.deleteResource(info, options)
		if err != nil {
			return err
		}
		resourceLocation := kubectlwait.ResourceLocation{
			GroupResource: info.Mapping.Resource.GroupResource(),
			Namespace:     info.Namespace,
			Name:          info.Name,
		}
		if status, ok := response.(*metav1.Status); ok && status.Details != nil {
			uidMap[resourceLocation] = status.Details.UID
			return nil
		}
		responseMetadata, err := meta.Accessor(response)
		if err != nil {
			// we don't have UID, but we didn't fail the delete, next best thing is just skipping the UID
			glog.V(1).Info(err)
			return nil
		}
		uidMap[resourceLocation] = responseMetadata.GetUID()

		return nil
	})
	if err != nil {
		return err
	}
	if found == 0 {
		fmt.Fprintf(o.Out, "No resources found\n")
		return nil
	}
	if !o.WaitForDeletion {
		return nil
	}
	// if we don't have a dynamic client, we don't want to wait.  Eventually when delete is cleaned up, this will likely
	// drop out.
	if o.DynamicClient == nil {
		return nil
	}

	effectiveTimeout := o.Timeout
	if effectiveTimeout == 0 {
		// if we requested to wait forever, set it to a week.
		effectiveTimeout = 168 * time.Hour
	}
	waitOptions := kubectlwait.WaitOptions{
		ResourceFinder: genericclioptions.ResourceFinderForResult(resource.InfoListVisitor(deletedInfos)),
		UIDMap:         uidMap,
		DynamicClient:  o.DynamicClient,
		Timeout:        effectiveTimeout,

		Printer:     printers.NewDiscardingPrinter(),
		ConditionFn: kubectlwait.IsDeleted,
		IOStreams:   o.IOStreams,
	}
	err = waitOptions.RunWait()
	if errors.IsForbidden(err) || errors.IsMethodNotSupported(err) {
		// if we're forbidden from waiting, we shouldn't fail.
		// if the resource doesn't support a verb we need, we shouldn't fail.
		glog.V(1).Info(err)
		return nil
	}
	return err
}

func (o *DeleteOptions) deleteResource(info *resource.Info, deleteOptions *metav1.DeleteOptions) (runtime.Object, error) {
	deleteResponse, err := resource.NewHelper(info.Client, info.Mapping).DeleteWithOptions(info.Namespace, info.Name, deleteOptions)
	if err != nil {
		return nil, cmdutil.AddSourceToErr("deleting", info.Source, err)
	}

	o.PrintObj(info)
	return deleteResponse, nil
}

// deletion printing is special because we do not have an object to print.
// This mirrors name printer behavior
func (o *DeleteOptions) PrintObj(info *resource.Info) {
	operation := "deleted"
	groupKind := info.Mapping.GroupVersionKind
	kindString := fmt.Sprintf("%s.%s", strings.ToLower(groupKind.Kind), groupKind.Group)
	if len(groupKind.Group) == 0 {
		kindString = strings.ToLower(groupKind.Kind)
	}

	if o.GracePeriod == 0 {
		operation = "force deleted"
	}

	if o.Output == "name" {
		// -o name: prints resource/name
		fmt.Fprintf(o.Out, "%s/%s\n", kindString, info.Name)
		return
	}

	// understandable output by default
	fmt.Fprintf(o.Out, "%s \"%s\" %s\n", kindString, info.Name, operation)
}
