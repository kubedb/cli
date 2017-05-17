package util

import (
	"fmt"
	"strings"

	tapi "github.com/k8sdb/apimachinery/api"
	k8serr "k8s.io/kubernetes/pkg/api/errors"
	cmdutil "k8s.io/kubernetes/pkg/kubectl/cmd/util"
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
	case strings.ToLower(tapi.ResourceKindDeletedDatabase):
	case strings.ToLower(tapi.ResourceTypeDeletedDatabase):
	case strings.ToLower(tapi.ResourceCodeDeletedDatabase):
		return tapi.ResourceKindDeletedDatabase, nil
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
	case tapi.ResourceKindDeletedDatabase:
		return nil
	default:
		return fmt.Errorf(`kubedb doesn't support a resource type "%v"`, kind)
	}
	return nil
}

func GetAllSupportedResources(f cmdutil.Factory) ([]string, error) {

	resources := map[string]string{
		tapi.ResourceNameElastic:          tapi.ResourceTypeElastic,
		tapi.ResourceNamePostgres:         tapi.ResourceTypePostgres,
		tapi.ResourceNameSnapshot: tapi.ResourceTypeSnapshot,
		tapi.ResourceNameDeletedDatabase:  tapi.ResourceTypeDeletedDatabase,
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
	tapi.ResourceCodeElastic:          tapi.ResourceTypeElastic,
	tapi.ResourceCodePostgres:         tapi.ResourceTypePostgres,
	tapi.ResourceCodeSnapshot: tapi.ResourceTypeSnapshot,
	tapi.ResourceCodeDeletedDatabase:  tapi.ResourceTypeDeletedDatabase,
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
