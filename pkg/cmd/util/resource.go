package util

import (
	"fmt"
	"strings"

	"github.com/ghodss/yaml"
	tapi "github.com/k8sdb/apimachinery/api"
	"github.com/k8sdb/apimachinery/client/clientset"
	"github.com/k8sdb/cli/pkg/cmd/decoder"
	k8serr "k8s.io/kubernetes/pkg/api/errors"
	cmdutil "k8s.io/kubernetes/pkg/kubectl/cmd/util"
	"k8s.io/kubernetes/pkg/runtime"
	"k8s.io/kubernetes/pkg/util/json"
	"k8s.io/kubernetes/pkg/util/strategicpatch"
)

func GetSupportedResource(resource string) (string, error) {
	switch strings.ToLower(resource) {
	case strings.ToLower(tapi.ResourceKindElastic),
		strings.ToLower(tapi.ResourceTypeElastic),
		strings.ToLower(tapi.ResourceCodeElastic):
		return tapi.ResourceKindElastic + "." + tapi.V1beta1SchemeGroupVersion.Group, nil
	case strings.ToLower(tapi.ResourceKindPostgres),
		strings.ToLower(tapi.ResourceTypePostgres),
		strings.ToLower(tapi.ResourceCodePostgres):
		return tapi.ResourceKindPostgres + "." + tapi.V1beta1SchemeGroupVersion.Group, nil
	case strings.ToLower(tapi.ResourceKindSnapshot),
		strings.ToLower(tapi.ResourceTypeSnapshot),
		strings.ToLower(tapi.ResourceCodeSnapshot):
		return tapi.ResourceKindSnapshot + "." + tapi.V1beta1SchemeGroupVersion.Group, nil
	case strings.ToLower(tapi.ResourceKindDormantDatabase),
		strings.ToLower(tapi.ResourceTypeDormantDatabase),
		strings.ToLower(tapi.ResourceCodeDormantDatabase):
		return tapi.ResourceKindDormantDatabase + "." + tapi.V1beta1SchemeGroupVersion.Group, nil
	default:
		return "", fmt.Errorf(`kubedb doesn't support a resource type "%v"`, resource)
	}
	return resource, nil
}

func CheckSupportedResource(kind string) error {
	switch kind {
	case tapi.ResourceKindElastic,
		tapi.ResourceKindPostgres,
		tapi.ResourceKindSnapshot,
		tapi.ResourceKindDormantDatabase:
		return nil
	default:
		return fmt.Errorf(`kubedb doesn't support a resource type "%v"`, kind)
	}
}

func GetAllSupportedResources(f cmdutil.Factory) ([]string, error) {

	resources := map[string]string{
		tapi.ResourceNameElastic:         tapi.ResourceKindElastic + "." + tapi.V1beta1SchemeGroupVersion.Group,
		tapi.ResourceNamePostgres:        tapi.ResourceKindPostgres + "." + tapi.V1beta1SchemeGroupVersion.Group,
		tapi.ResourceNameSnapshot:        tapi.ResourceKindSnapshot + "." + tapi.V1beta1SchemeGroupVersion.Group,
		tapi.ResourceNameDormantDatabase: tapi.ResourceKindDormantDatabase + "." + tapi.V1beta1SchemeGroupVersion.Group,
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
		return !ok
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
		)
	case tapi.ResourceKindPostgres:
		preconditions = append(
			preconditions,
		)
	case tapi.ResourceKindDormantDatabase:
		preconditions = append(
			preconditions,
			RequireChainKeyUnchanged("spec.origin"),
		)
	}
	return preconditions
}

func GetConditionalPreconditionFunc(kind string) []strategicpatch.PreconditionFunc {
	preconditions := []strategicpatch.PreconditionFunc{}

	switch kind {
	case tapi.ResourceKindElastic:
		preconditions = append(
			preconditions,
			RequireChainKeyUnchanged("spec.version"),
			RequireChainKeyUnchanged("spec.storage"),
			RequireChainKeyUnchanged("spec.nodeSelector"),
			RequireChainKeyUnchanged("spec.init"),
		)
	case tapi.ResourceKindPostgres:
		preconditions = append(
			preconditions,
			RequireChainKeyUnchanged("spec.version"),
			RequireChainKeyUnchanged("spec.storage"),
			RequireChainKeyUnchanged("spec.databaseSecret"),
			RequireChainKeyUnchanged("spec.nodeSelector"),
			RequireChainKeyUnchanged("spec.init"),
		)
	}
	return preconditions
}

func CheckResourceExists(extClient clientset.ExtensionInterface, kind, name, namespace string) (bool, error) {
	var err error
	switch kind {
	case tapi.ResourceKindElastic:
		_, err = extClient.Elastics(namespace).Get(name)
	case tapi.ResourceKindPostgres:
		_, err = extClient.Postgreses(namespace).Get(name)
	}

	if err != nil {
		if k8serr.IsNotFound(err) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func IsPreconditionFailed(err error) bool {
	_, ok := err.(errPreconditionFailed)
	return ok
}

type errPreconditionFailed struct {
	message string
}

func newErrPreconditionFailed(target map[string]interface{}) errPreconditionFailed {
	s := fmt.Sprintf("precondition failed for: %v", target)
	return errPreconditionFailed{s}
}

func (err errPreconditionFailed) Error() string {
	return err.message
}

func CheckConditionalPrecondition(patchData []byte, fns ...strategicpatch.PreconditionFunc) error {
	patch := make(map[string]interface{})
	if err := json.Unmarshal(patchData, &patch); err != nil {
		return err
	}
	for _, fn := range fns {
		if !fn(patch) {
			return newErrPreconditionFailed(patch)
		}
	}
	return nil
}
