package cmd

import (
	"fmt"
	"io"
	"strings"

	"github.com/k8sdb/kubedb/pkg/cmd/util"
	"github.com/k8sdb/kubedb/pkg/kube"
	"github.com/spf13/cobra"
	"k8s.io/kubernetes/pkg/api/meta"
	"k8s.io/kubernetes/pkg/kubectl/cmd/templates"
	cmdutil "k8s.io/kubernetes/pkg/kubectl/cmd/util"
	"k8s.io/kubernetes/pkg/kubectl/resource"
	"k8s.io/kubernetes/pkg/runtime"
)

// ref: k8s.io/kubernetes/pkg/kubectl/cmd/delete.go

var (
	delete_long = templates.LongDesc(`
		Delete resources by filenames, stdin, resources and names, or by resources and label selector.
		JSON and YAML formats are accepted.

		Note that the delete command does NOT do resource version checks`)

	delete_example = templates.Examples(`
		# Delete a elastic using the type and name specified in elastic.json.
		kubedb delete -f ./elastic.json

		# Delete a postgres based on the type and name in the JSON passed into stdin.
		cat postgres.json | kubedb delete -f -

		# Delete elastic with label elastic.k8sdb.com/name=elasticsearch-demo.
		kubedb delete elastic -l elastic.k8sdb.com/name=elasticsearch-demo`)
)

func NewCmdDelete(out, errOut io.Writer) *cobra.Command {
	options := &resource.FilenameOptions{}

	cmd := &cobra.Command{
		Use:     "delete ([-f FILENAME] | TYPE [(NAME | -l label)])",
		Short:   "Delete resources by filenames, stdin, resources and names, or by resources and label selector",
		Long:    delete_long,
		Example: delete_example,
		Run: func(cmd *cobra.Command, args []string) {
			f := kube.NewKubeFactory(cmd)
			cmdutil.CheckErr(RunDelete(f, cmd, out, args, options))
		},
	}

	util.AddFilenameOptionFlags(cmd, options)
	util.AddDeleteFlags(cmd)
	return cmd
}

func RunDelete(f cmdutil.Factory, cmd *cobra.Command, out io.Writer, args []string, options *resource.FilenameOptions) error {
	cmdNamespace, enforceNamespace, err := f.DefaultNamespace()
	if err != nil {
		return err
	}
	mapper, typer, err := f.UnstructuredObject()
	if err != nil {
		return err
	}

	if len(args) > 0 {
		resources := strings.Split(args[0], ",")
		for i, r := range resources {
			items := strings.Split(r, "/")
			kind, err := util.GetSupportedResourceKind(items[0])
			if err != nil {
				return err
			}
			items[0] = kind
			resources[i] = strings.Join(items, "/")
		}
		args[0] = strings.Join(resources, ",")
	}

	r := resource.NewBuilder(mapper, typer, resource.ClientMapperFunc(f.UnstructuredClientForMapping), runtime.UnstructuredJSONScheme).
		ContinueOnError().
		NamespaceParam(cmdNamespace).DefaultNamespace().
		FilenameParam(enforceNamespace, options).
		SelectorParam(cmdutil.GetFlagString(cmd, "selector")).
		ResourceTypeOrNameArgs(true, args...).RequireObject(true).
		Flatten().
		Do()
	err = r.Err()
	if err != nil {
		return err
	}

	shortOutput := cmdutil.GetFlagString(cmd, "output") == "name"
	return deleteResult(r, out, shortOutput, mapper)
}

func deleteResult(r *resource.Result, out io.Writer, shortOutput bool, mapper meta.RESTMapper) error {
	infoList := make([]*resource.Info, 0)
	err := r.Visit(func(info *resource.Info, err error) error {
		if err != nil {
			return err
		}
		kind := info.GetObjectKind().GroupVersionKind().Kind
		if err := util.CheckSupportedResource(kind); err != nil {
			return err
		}

		infoList = append(infoList, info)
		return nil

	})
	if err != nil {
		return err
	}

	found := 0
	for _, info := range infoList {
		found++
		if err := deleteResource(info, out, shortOutput, mapper); err != nil {
			return err
		}
	}

	if found == 0 {
		fmt.Fprintf(out, "No resources found\n")
	}
	return nil
}

func deleteResource(info *resource.Info, out io.Writer, shortOutput bool, mapper meta.RESTMapper) error {
	if err := resource.NewHelper(info.Client, info.Mapping).Delete(info.Namespace, info.Name); err != nil {
		return cmdutil.AddSourceToErr("deleting", info.Source, err)
	}
	cmdutil.PrintSuccess(mapper, shortOutput, out, info.Mapping.Resource, info.Name, false, "deleted")
	return nil
}
