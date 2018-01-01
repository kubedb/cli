package cmds

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	core_util "github.com/appscode/kutil/core/v1"
	"github.com/ghodss/yaml"
	tapi "github.com/kubedb/apimachinery/apis/kubedb/v1alpha1"
	"github.com/kubedb/apimachinery/client/scheme"
	tcs "github.com/kubedb/apimachinery/client/typed/kubedb/v1alpha1"
	"github.com/kubedb/cli/pkg/encoder"
	"github.com/kubedb/cli/pkg/kube"
	"github.com/kubedb/cli/pkg/util"
	"github.com/kubedb/cli/pkg/validator"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/jsonmergepatch"
	yml "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/kubernetes/pkg/kubectl/cmd/templates"
	cmdutil "k8s.io/kubernetes/pkg/kubectl/cmd/util"
	"k8s.io/kubernetes/pkg/kubectl/resource"
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
	return deleteResult(f, r, out, shortOutput, mapper)
}

func deleteResult(f cmdutil.Factory, r *resource.Result, out io.Writer, shortOutput bool, mapper meta.RESTMapper) error {
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
		if err := deleteResource(f, info, out, shortOutput, mapper); err != nil {
			return err
		}
	}

	if found == 0 {
		fmt.Fprintln(out, "No resources found")
	}
	return nil
}

func deleteResource(f cmdutil.Factory, info *resource.Info, out io.Writer, shortOutput bool, mapper meta.RESTMapper) error {
	if forceDeletion {
		if err := forceRemoveFinalizer(f, info); err != nil {
			return cmdutil.AddSourceToErr("deleting forcefully", info.Source, err)
		}
	}
	if err := resource.NewHelper(info.Client, info.Mapping).Delete(info.Namespace, info.Name); err != nil {
		return cmdutil.AddSourceToErr("deleting", info.Source, err)
	}
	cmdutil.PrintSuccess(mapper, shortOutput, out, info.Mapping.Resource, info.Name, false, "deleted")
	return nil
}

func forceRemoveFinalizer(f cmdutil.Factory, info *resource.Info) error {
	restConfig, err := f.ClientConfig()
	if err != nil {
		return err
	}
	extClient := tcs.NewForConfigOrDie(restConfig)
	codec := scheme.Codecs.LegacyCodec(tapi.SchemeGroupVersion)
	currOriginalObj, err := util.GetStructuredObject(info.Object)
	if err != nil {
		return err
	}

	originalSerialization, err := runtime.Encode(codec, currOriginalObj)
	if err != nil {
		return err
	}

	originalJson, err := yml.ToJSON(originalSerialization)
	if err != nil {
		return err
	}

	if modJson, err := patchFinalizer(info); err == nil {
		if patch, err := jsonmergepatch.CreateThreeWayJSONMergePatch(originalJson, modJson, originalJson); err == nil {
			h := resource.NewHelper(extClient.RESTClient(), info.Mapping)
			_, err := extClient.RESTClient().Patch(types.MergePatchType).
				NamespaceIfScoped(info.Namespace, h.NamespaceScoped).
				Resource(h.Resource).
				Name(info.Name).
				Body(patch).
				Do().
				Get()
			return err
		}
		return err
	}
	return err
}

func patchFinalizer(info *resource.Info) ([]byte, error) {
	objByte, err := encoder.Encode(info.Object)
	if err != nil {
		return nil, err
	}
	kind := info.Object.GetObjectKind().GroupVersionKind().Kind
	switch kind {
	case tapi.ResourceKindElasticsearch:
		var elasticSearch *tapi.Elasticsearch
		if err := yaml.Unmarshal(objByte, &elasticSearch); err != nil {
			return nil, err
		}
		if core_util.HasFinalizer(elasticSearch.ObjectMeta, "kubedb.com") {
			elasticSearch.ObjectMeta = core_util.RemoveFinalizer(elasticSearch.ObjectMeta, "kubedb.com")
		}
		return json.Marshal(elasticSearch)
	case tapi.ResourceKindPostgres:
		var postgres *tapi.Postgres
		if err := yaml.Unmarshal(objByte, &postgres); err != nil {
			return nil, err
		}
		if core_util.HasFinalizer(postgres.ObjectMeta, "kubedb.com") {
			postgres.ObjectMeta = core_util.RemoveFinalizer(postgres.ObjectMeta, "kubedb.com")
		}
		return json.Marshal(postgres)
	case tapi.ResourceKindMySQL:
		var mysql *tapi.MySQL
		if err := yaml.Unmarshal(objByte, &mysql); err != nil {
			return nil, err
		}
		if core_util.HasFinalizer(mysql.ObjectMeta, "kubedb.com") {
			mysql.ObjectMeta = core_util.RemoveFinalizer(mysql.ObjectMeta, "kubedb.com")
		}
		return json.Marshal(mysql)
	case tapi.ResourceKindMongoDB:
		var mongodb *tapi.MongoDB
		if err := yaml.Unmarshal(objByte, &mongodb); err != nil {
			return nil, err
		}
		if core_util.HasFinalizer(mongodb.ObjectMeta, "kubedb.com") {
			mongodb.ObjectMeta = core_util.RemoveFinalizer(mongodb.ObjectMeta, "kubedb.com")
		}
		return json.Marshal(mongodb)
	case tapi.ResourceKindRedis:
		var redis *tapi.Redis
		if err := yaml.Unmarshal(objByte, &redis); err != nil {
			return nil, err
		}
		if core_util.HasFinalizer(redis.ObjectMeta, "kubedb.com") {
			redis.ObjectMeta = core_util.RemoveFinalizer(redis.ObjectMeta, "kubedb.com")
		}
		return json.Marshal(redis)
	case tapi.ResourceKindMemcached:
		var memcached *tapi.Memcached
		if err := yaml.Unmarshal(objByte, &memcached); err != nil {
			return nil, err
		}
		if core_util.HasFinalizer(memcached.ObjectMeta, "kubedb.com") {
			memcached.ObjectMeta = core_util.RemoveFinalizer(memcached.ObjectMeta, "kubedb.com")
		}
		return json.Marshal(memcached)
	case tapi.ResourceKindDormantDatabase:
		var dormantdb *tapi.DormantDatabase
		if err := yaml.Unmarshal(objByte, &dormantdb); err != nil {
			return nil, err
		}
		if core_util.HasFinalizer(dormantdb.ObjectMeta, "kubedb.com") {
			dormantdb.ObjectMeta = core_util.RemoveFinalizer(dormantdb.ObjectMeta, "kubedb.com")
		}
		return json.Marshal(dormantdb)
	case tapi.ResourceKindSnapshot:
		var snapshot *tapi.Snapshot
		if err := yaml.Unmarshal(objByte, &snapshot); err != nil {
			return nil, err
		}
		if core_util.HasFinalizer(snapshot.ObjectMeta, "kubedb.com") {
			snapshot.ObjectMeta = core_util.RemoveFinalizer(snapshot.ObjectMeta, "kubedb.com")
		}
		return json.Marshal(snapshot)
	}
	return nil, fmt.Errorf(`invalid kind: "%v"`, kind)
}
