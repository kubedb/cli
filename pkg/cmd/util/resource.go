package util

import (
	"fmt"
	"strings"

	"github.com/ghodss/yaml"
	tapi "github.com/k8sdb/apimachinery/api"
	"github.com/k8sdb/kubedb/pkg/cmd/decoder"
	k8serr "k8s.io/kubernetes/pkg/api/errors"
	cmdutil "k8s.io/kubernetes/pkg/kubectl/cmd/util"
	"k8s.io/kubernetes/pkg/runtime"
	"k8s.io/kubernetes/pkg/util/strategicpatch"
)

func GetSupportedResourceKind(resource string) (string, error) {
	switch strings.ToLower(resource) {
	case strings.ToLower(tapi.ResourceKindElastic):
	case strings.ToLower(tapi.ResourceTypeElastic):
	case strings.ToLower(tapi.ResourceCodeElastic):
		return tapi.ResourceKindElastic, nil
	case strings.ToLower(tapi.ResourceKindPostgres):
	case strings.ToLower(tapi.ResourceTypePostgres):
	case strings.ToLower(tapi.ResourceCodePostgres):
		return tapi.ResourceKindPostgres, nil
	case strings.ToLower(tapi.ResourceKindSnapshot):
	case strings.ToLower(tapi.ResourceTypeSnapshot):
	case strings.ToLower(tapi.ResourceCodeSnapshot):
		return tapi.ResourceKindSnapshot, nil
	case strings.ToLower(tapi.ResourceKindDormantDatabase):
	case strings.ToLower(tapi.ResourceTypeDormantDatabase):
	case strings.ToLower(tapi.ResourceCodeDormantDatabase):
		return tapi.ResourceKindDormantDatabase, nil
	default:
		return "", fmt.Errorf(`kubedb doesn't support a resource type "%v"`, resource)
	}
	return resource, nil
}

func CheckSupportedResource(kind string) error {
	switch kind {
	case tapi.ResourceKindElastic:
	case tapi.ResourceKindPostgres:
	case tapi.ResourceKindSnapshot:
	case tapi.ResourceKindDormantDatabase:
		return nil
	default:
		return fmt.Errorf(`kubedb doesn't support a resource type "%v"`, kind)
	}
	return nil
}

func GetAllSupportedResources(f cmdutil.Factory) ([]string, error) {

	resources := map[string]string{
		tapi.ResourceNameElastic:         tapi.ResourceTypeElastic,
		tapi.ResourceNamePostgres:        tapi.ResourceTypePostgres,
		tapi.ResourceNameSnapshot:        tapi.ResourceTypeSnapshot,
		tapi.ResourceNameDormantDatabase: tapi.ResourceTypeDormantDatabase,
	}

	clientset, err := f.ClientSet()
	if err != nil {
		return nil, err
	}

	availableResources := make([]string, 0)
	for key, val := range resources {
		_, err := clientset.ThirdPartyResources().Get(key + "." + tapi.V1beta1SchemeGroupVersion.Group)
		if err != nil {
			if k8serr.IsNotFound(err) {
				continue
			}
			return nil, err
		}
		availableResources = append(availableResources, val)
	}

	return availableResources, nil
}

var ShortForms = map[string]string{
	tapi.ResourceCodeElastic:         tapi.ResourceTypeElastic,
	tapi.ResourceCodePostgres:        tapi.ResourceTypePostgres,
	tapi.ResourceCodeSnapshot:        tapi.ResourceTypeSnapshot,
	tapi.ResourceCodeDormantDatabase: tapi.ResourceTypeDormantDatabase,
}

func ResourceShortFormFor(resource string) (string, bool) {
	var alias string
	exists := false
	for k, val := range ShortForms {
		if val == resource {
			alias = k
			exists = true
			break
		}
	}
	return alias, exists
}

func GetObjectData(obj runtime.Object) ([]byte, error) {
	return yaml.Marshal(obj)
}

func GetStructuredObject(obj runtime.Object) (runtime.Object, error) {
	kind := obj.GetObjectKind().GroupVersionKind().Kind
	data, err := GetObjectData(obj)
	if err != nil {
		return obj, err
	}
	return decoder.Decode(kind, data)
}

func checkChainKeyUnchanged(key string, mapData map[string]interface{}) bool {
	keys := strings.Split(key, ".")
	val, ok := mapData[keys[0]]
	if !ok || len(keys) == 1 {
		return false
	}

	newKey := strings.Join(keys[1:], ".")
	return checkChainKeyUnchanged(newKey, val.(map[string]interface{}))
}

func RequireChainKeyUnchanged(key string) strategicpatch.PreconditionFunc {
	return func(patch interface{}) bool {
		patchMap, ok := patch.(map[string]interface{})
		if !ok {
			fmt.Println("Invalid data")
			return true
		}
		check := checkChainKeyUnchanged(key, patchMap)
		if !check {
			fmt.Println(key, "was changed")
		}
		return check
	}
}

func GetPreconditionFunc(kind string) []strategicpatch.PreconditionFunc {
	preconditions := []strategicpatch.PreconditionFunc{
		strategicpatch.RequireKeyUnchanged("apiVersion"),
		strategicpatch.RequireKeyUnchanged("kind"),
		strategicpatch.RequireMetadataKeyUnchanged("name"),
		strategicpatch.RequireMetadataKeyUnchanged("namespace"),
		strategicpatch.RequireKeyUnchanged("status"),
	}

	switch kind {
	case tapi.ResourceKindElastic:
		preconditions = append(
			preconditions,
			RequireChainKeyUnchanged("spec.version"),
		)
	case tapi.ResourceKindPostgres:
	case tapi.ResourceKindSnapshot:
	case tapi.ResourceKindDormantDatabase:
	}
	return preconditions
}
