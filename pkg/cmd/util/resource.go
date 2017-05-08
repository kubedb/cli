package util

import (
	"fmt"
	"strings"

	tapi "github.com/k8sdb/apimachinery/api"
	k8serr "k8s.io/kubernetes/pkg/api/errors"
	cmdutil "k8s.io/kubernetes/pkg/kubectl/cmd/util"
)

const (
	ShortResourceTypeElastic          = "es"
	ShortResourceTypePostgres         = "pg"
	ShortResourceTypeDatabaseSnapshot = "dbs"
	ShortResourceTypeDeletedDatabase  = "ddb"
)

func GetSupportedResourceKind(resource string) (string, error) {
	switch strings.ToLower(resource) {
	case strings.ToLower(tapi.ResourceKindElastic):
	case strings.ToLower(tapi.ResourceTypeElastic):
	case strings.ToLower(ShortResourceTypeElastic):
		return tapi.ResourceKindElastic, nil
	case strings.ToLower(tapi.ResourceKindPostgres):
	case strings.ToLower(tapi.ResourceTypePostgres):
	case strings.ToLower(ShortResourceTypePostgres):
		return tapi.ResourceKindPostgres, nil
	case strings.ToLower(tapi.ResourceKindDatabaseSnapshot):
	case strings.ToLower(tapi.ResourceTypeDatabaseSnapshot):
	case strings.ToLower(ShortResourceTypeDatabaseSnapshot):
		return tapi.ResourceKindDatabaseSnapshot, nil
	case strings.ToLower(tapi.ResourceKindDeletedDatabase):
	case strings.ToLower(tapi.ResourceTypeDeletedDatabase):
	case strings.ToLower(ShortResourceTypeDeletedDatabase):
		return tapi.ResourceKindDeletedDatabase, nil
	default:
		return "", fmt.Errorf(`kubedb doesn't support a resource type "%v"`, resource)
	}
	return resource, nil
}

func GetAllSupportedResources(f cmdutil.Factory) ([]string, error) {

	resources := map[string]string{
		tapi.ResourceNameElastic:          tapi.ResourceTypeElastic,
		tapi.ResourceNamePostgres:         tapi.ResourceTypePostgres,
		tapi.ResourceNameDatabaseSnapshot: tapi.ResourceTypeDatabaseSnapshot,
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
	ShortResourceTypeElastic:          tapi.ResourceTypeElastic,
	ShortResourceTypePostgres:         tapi.ResourceTypePostgres,
	ShortResourceTypeDatabaseSnapshot: tapi.ResourceTypeDatabaseSnapshot,
	ShortResourceTypeDeletedDatabase:  tapi.ResourceTypeDeletedDatabase,
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
