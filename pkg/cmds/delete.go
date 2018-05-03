package cmds

import (
	"fmt"
	"io"
	"strings"

	"github.com/kubedb/cli/pkg/kube"
	"github.com/kubedb/cli/pkg/util"
	"github.com/kubedb/cli/pkg/validator"
	"github.com/spf13/cobra"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/kubernetes/pkg/kubectl/cmd/templates"
	cmdutil "k8s.io/kubernetes/pkg/kubectl/cmd/util"
	"k8s.io/kubernetes/pkg/kubectl/resource"
)

// ref: k8s.io/kubernetes/pkg/kubectl/cmd/delete.go

var (
	deleteLong = templates.LongDesc(`
		Delete resources by filenames, stdin, resources and names, or by resources and label selector.
		JSON and YAML formats are accepted.

		Note that the delete command does NOT do resource version checks`)

	deleteExample = templates.Examples(`
		# Delete a elasticsearch using the type and name specified in elastic.json.
		kubedb delete -f ./elastic.json

		# Delete a postgres based on the type and name in the JSON passed into stdin.
		cat postgres.json | kubedb delete -f -

		# Delete elasticsearch with label elasticsearch.kubedb.com/name=elasticsearch-demo.
		kubedb delete elasticsearch -l elasticsearch.kubedb.com/name=elasticsearch-demo

		# Force delete a mysql object
		kubedb delete mysql ms-demo --force

		# Delete all mysql objects
		kubedb delete mysql --all`)
)

func NewCmdDelete(out, errOut io.Writer) *cobra.Command {
	options := &resource.FilenameOptions{}

	cmd := &cobra.Command{
		Use:     "delete ([-f FILENAME] | TYPE [(NAME | -l label | --all)])",
		Short:   "Delete resources by filenames, stdin, resources and names, or by resources and label selector",
		Long:    deleteLong,
		Example: deleteExample,
		Run: func(cmd *cobra.Command, args []string) {
			f := kube.NewKubeFactory(cmd)
			cmdutil.CheckErr(RunDelete(f, cmd, out, args, options))
		},
	}

	util.AddDeleteFlags(cmd, options)
	return cmd
}

func RunDelete(f cmdutil.Factory, cmd *cobra.Command, out io.Writer, args []string, options *resource.FilenameOptions) error {
	selector := cmdutil.GetFlagString(cmd, "selector")
	deleteAll := cmdutil.GetFlagBool(cmd, "all")
	cmdNamespace, enforceNamespace := util.GetNamespace(cmd)

	mapper, _ := f.Object()

	if len(args) > 0 {
		resources := strings.Split(args[0], ",")
		for i, r := range resources {
			items := strings.Split(r, "/")
			kind, err := util.GetSupportedResource(items[0])
			if err != nil {
				return err
			}
			items[0] = kind
			resources[i] = strings.Join(items, "/")
		}
		args[0] = strings.Join(resources, ",")
	}

	r := f.NewBuilder().Unstructured().
		NamespaceParam(cmdNamespace).DefaultNamespace().
		FilenameParam(enforceNamespace, options).
		LabelSelectorParam(selector).
		SelectAllParam(deleteAll).
		ResourceTypeOrNameArgs(false, args...).RequireObject(true).
		Flatten().
		Do()
	err := r.Err()
	if err != nil {
		return err
	}

	return deleteResult(f, cmd, r, out, mapper)
}

func deleteResult(f cmdutil.Factory, cmd *cobra.Command, r *resource.Result, out io.Writer, mapper meta.RESTMapper) error {
	shortOutput := cmdutil.GetFlagString(cmd, "output") == "name"

	infoList := make([]*resource.Info, 0)
	err := r.Visit(func(info *resource.Info, err error) error {
		if err != nil {
			return err
		}
		kind := info.Object.GetObjectKind().GroupVersionKind().Kind
		if err := util.CheckSupportedResource(kind); err != nil {
			return err
		}

		if err := validator.ValidateDeletion(info); err != nil {
			return cmdutil.AddSourceToErr("validating", info.Source, err)
		}

		infoList = append(infoList, info)
		return nil

	})
	if err != nil {
		return err
	}

	if len(infoList) == 0 {
		fmt.Fprintln(out, "No resources found")
		return nil
	}

	for _, info := range infoList {
		if err := deleteResource(info, out, shortOutput); err != nil {
			return err
		}
	}

	return nil
}

func deleteResource(info *resource.Info, out io.Writer, shortOutput bool) error {
	if err := resource.NewHelper(info.Client, info.Mapping).Delete(info.Namespace, info.Name); err != nil && !kerr.IsNotFound(err) {
		return cmdutil.AddSourceToErr("deleting", info.Source, err)
	}
	cmdutil.PrintSuccess(shortOutput, out, info.Object, false, "deleted")
	return nil
}
