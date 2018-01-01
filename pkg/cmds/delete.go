package cmds

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/kubedb/cli/pkg/kube"
	"github.com/kubedb/cli/pkg/util"
	"github.com/kubedb/cli/pkg/validator"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/kubernetes/pkg/kubectl/cmd/templates"
	"k8s.io/kubernetes/pkg/kubectl/resource"

	core_util "github.com/appscode/kutil/core/v1"
	"github.com/ghodss/yaml"
	tapi "github.com/kubedb/apimachinery/apis/kubedb/v1alpha1"
	"github.com/kubedb/cli/pkg/encoder"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/jsonmergepatch"
	cmdutil "k8s.io/kubernetes/pkg/kubectl/cmd/util"
)

// ref: k8s.io/kubernetes/pkg/kubectl/cmd/delete.go

var (
	forceDeletion = false

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
		kubedb delete elasticsearch -l elasticsearch.kubedb.com/name=elasticsearch-demo`)
)

func NewCmdDelete(out, errOut io.Writer) *cobra.Command {
	options := &resource.FilenameOptions{}

	cmd := &cobra.Command{
		Use:     "delete ([-f FILENAME] | TYPE [(NAME | -l label)])",
		Short:   "Delete resources by filenames, stdin, resources and names, or by resources and label selector",
		Long:    deleteLong,
		Example: deleteExample,
		Run: func(cmd *cobra.Command, args []string) {
			f := kube.NewKubeFactory(cmd)
			cmdutil.CheckErr(RunDelete(f, cmd, out, args, options))
		},
	}

	util.AddDeleteFlags(cmd, options)
	cmd.Flags().BoolVar(&forceDeletion, "force", false, "Immediate deletion of some resources may result in inconsistency or data loss.")
	return cmd
}

func RunDelete(f cmdutil.Factory, cmd *cobra.Command, out io.Writer, args []string, options *resource.FilenameOptions) error {
	selector := cmdutil.GetFlagString(cmd, "selector")
	cmdNamespace, enforceNamespace := util.GetNamespace(cmd)
	categoryExpander := f.CategoryExpander()
	mapper, typer, err := f.UnstructuredObject()
	if err != nil {
		return err
	}

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

	r := resource.NewBuilder(mapper, categoryExpander, typer, resource.ClientMapperFunc(f.UnstructuredClientForMapping), unstructured.UnstructuredJSONScheme).
		ContinueOnError().
		NamespaceParam(cmdNamespace).DefaultNamespace().
		FilenameParam(enforceNamespace, options).
		SelectorParam(selector).
		ResourceTypeOrNameArgs(false, args...).RequireObject(true).
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

	found := 0
	for _, info := range infoList {
		found++
		if err := deleteResource(info, out, shortOutput, mapper); err != nil {
			return err
		}
	}

	if found == 0 {
		fmt.Fprintln(out, "No resources found")
	}
	return nil
}

func deleteResource(info *resource.Info, out io.Writer, shortOutput bool, mapper meta.RESTMapper) error {
	if forceDeletion {
		err := forceRemoveFinalizer(info)
		cmdutil.AddSourceToErr("deleting forcefully", info.Source, err)
	}
	if err := resource.NewHelper(info.Client, info.Mapping).Delete(info.Namespace, info.Name); err != nil {
		return cmdutil.AddSourceToErr("deleting", info.Source, err)
	}
	cmdutil.PrintSuccess(mapper, shortOutput, out, info.Mapping.Resource, info.Name, false, "deleted")
	return nil
}

func forceRemoveFinalizer(info *resource.Info) error {
	objMeta, err := getObjectMeta(info)
	if err == nil {
		if core_util.HasFinalizer(objMeta, "kubedb.com") {
			editedMeta := core_util.RemoveFinalizer(objMeta, "kubedb.com")
			curJson, err := json.Marshal(objMeta)
			if err != nil {
				return err
			}
			modJson, err := json.Marshal(editedMeta)
			if err != nil {
				return err
			}
			patch, err := jsonmergepatch.CreateThreeWayJSONMergePatch(curJson, modJson, curJson)
			patched, err := resource.NewHelper(info.Client, info.Mapping).Patch(info.Namespace, info.Name, types.StrategicMergePatchType, patch)
			return err
			fmt.Println("patch:", patch)
			fmt.Println("patched:", patched)
		}
	}
	return err
}

func getObjectMeta(info *resource.Info) (metav1.ObjectMeta, error) {
	objByte, err := encoder.Encode(info.Object)
	if err != nil {
		return metav1.ObjectMeta{}, err
	}

	kind := info.Object.GetObjectKind().GroupVersionKind().Kind
	switch kind {
	case tapi.ResourceKindElasticsearch:
		var elasticsearch *tapi.Elasticsearch
		if err := yaml.Unmarshal(objByte, &elasticsearch); err != nil {
			return metav1.ObjectMeta{}, err
		}
		return elasticsearch.ObjectMeta, nil
	case tapi.ResourceKindPostgres:
		var postgres *tapi.Postgres
		if err := yaml.Unmarshal(objByte, &postgres); err != nil {
			return metav1.ObjectMeta{}, err
		}
		return postgres.ObjectMeta, nil
	case tapi.ResourceKindMySQL:
		var mysql *tapi.MySQL
		if err := yaml.Unmarshal(objByte, &mysql); err != nil {
			return metav1.ObjectMeta{}, err
		}
		return mysql.ObjectMeta, nil
	case tapi.ResourceKindMongoDB:
		var mongodb *tapi.MongoDB
		if err := yaml.Unmarshal(objByte, &mongodb); err != nil {
			return metav1.ObjectMeta{}, err
		}
		return mongodb.ObjectMeta, nil

	case tapi.ResourceKindRedis:
		var redis *tapi.Redis
		if err := yaml.Unmarshal(objByte, &redis); err != nil {
			return metav1.ObjectMeta{}, err
		}
		return redis.ObjectMeta, nil
	case tapi.ResourceKindMemcached:
		var memcached *tapi.Memcached
		if err := yaml.Unmarshal(objByte, &memcached); err != nil {
			return metav1.ObjectMeta{}, err
		}
		return memcached.ObjectMeta, nil
	}
	return metav1.ObjectMeta{}, nil
}
