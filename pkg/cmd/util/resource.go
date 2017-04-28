package util

import (
	"fmt"
	"strings"

	tapi "github.com/k8sdb/apimachinery/api"
	"github.com/k8sdb/kubedb/pkg/kube"
	k8serr "k8s.io/kubernetes/pkg/api/errors"
)

func CheckSupportedResources(args []string) error {
	if len(args) > 0 {
		resources := strings.Split(args[0], ",")
		for _, r := range resources {
			switch strings.ToLower(r) {
			case strings.ToLower(tapi.ResourceKindElastic):
			case strings.ToLower(tapi.ResourceTypeElastic):
			case strings.ToLower(tapi.ResourceKindPostgres):
			case strings.ToLower(tapi.ResourceTypePostgres):
			case strings.ToLower(tapi.ResourceKindDatabaseSnapshot):
			case strings.ToLower(tapi.ResourceTypeDatabaseSnapshot):
			case strings.ToLower(tapi.ResourceKindDeletedDatabase):
			case strings.ToLower(tapi.ResourceTypeDeletedDatabase):
				continue
			case "all":
				continue
			default:
				return fmt.Errorf(`kubedb doesn't support a resource type "%v"`, r)
			}
		}
	}
	return nil
}

func GetAllSupportedResources(client *kube.Client) (string, error) {

	resources := map[string]string{
		tapi.ResourceNameElastic:          tapi.ResourceTypeElastic,
		tapi.ResourceNamePostgres:         tapi.ResourceTypePostgres,
		tapi.ResourceNameDatabaseSnapshot: tapi.ResourceTypeDatabaseSnapshot,
		tapi.ResourceNameDeletedDatabase:  tapi.ResourceTypeDeletedDatabase,
	}

	clientset, err := client.ClientSet()
	if err != nil {
		return "", err
	}

	availableResources := make([]string, 0)
	for key, val := range resources {
		_, err := clientset.ThirdPartyResources().Get(key + "." + tapi.V1beta1SchemeGroupVersion.Group)
		if err != nil {
			if k8serr.IsNotFound(err) {
				continue
			}
			return "", err
		}
		availableResources = append(availableResources, val)
	}

	return strings.Join(availableResources, ","), nil
}

var ShortForms = map[string]string{
	"es": tapi.ResourceTypeElastic,
	"pg": tapi.ResourceTypePostgres,
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
