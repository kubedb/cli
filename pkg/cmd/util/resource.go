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

func CheckSupportedResources(args []string) error {
	if len(args) > 0 {
		resources := strings.Split(args[0], ",")
		for _, r := range resources {
			switch strings.ToLower(r) {
			case strings.ToLower(tapi.ResourceKindElastic):
			case strings.ToLower(tapi.ResourceTypeElastic):
			case strings.ToLower(ShortResourceTypeElastic):
			case strings.ToLower(tapi.ResourceKindPostgres):
			case strings.ToLower(tapi.ResourceTypePostgres):
			case strings.ToLower(ShortResourceTypePostgres):
			case strings.ToLower(tapi.ResourceKindDatabaseSnapshot):
			case strings.ToLower(tapi.ResourceTypeDatabaseSnapshot):
			case strings.ToLower(ShortResourceTypeDatabaseSnapshot):
			case strings.ToLower(tapi.ResourceKindDeletedDatabase):
			case strings.ToLower(tapi.ResourceTypeDeletedDatabase):
			case strings.ToLower(ShortResourceTypeDeletedDatabase):
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

func ReplaceAliases(args string) string {
	resources := strings.Split(args, ",")
	typeList := make([]string, 0)
	for _, r := range resources {
		switch strings.ToLower(r) {
		case strings.ToLower(ShortResourceTypeElastic):
			typeList = append(typeList, tapi.ResourceTypeElastic)
		case strings.ToLower(ShortResourceTypePostgres):
			typeList = append(typeList, tapi.ResourceTypePostgres)
		case strings.ToLower(ShortResourceTypeDatabaseSnapshot):
			typeList = append(typeList, tapi.ResourceTypeDatabaseSnapshot)
		case strings.ToLower(ShortResourceTypeDeletedDatabase):
			typeList = append(typeList, tapi.ResourceTypeDeletedDatabase)
		default:
			typeList = append(typeList, r)
		}
	}

	return strings.Join(typeList, ",")
}

func GetAllSupportedResources(f cmdutil.Factory) (string, error) {

	resources := map[string]string{
		tapi.ResourceNameElastic:          tapi.ResourceTypeElastic,
		tapi.ResourceNamePostgres:         tapi.ResourceTypePostgres,
		tapi.ResourceNameDatabaseSnapshot: tapi.ResourceTypeDatabaseSnapshot,
		tapi.ResourceNameDeletedDatabase:  tapi.ResourceTypeDeletedDatabase,
	}

	clientset, err := f.ClientSet()
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
