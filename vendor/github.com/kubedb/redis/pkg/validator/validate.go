package validator

import (
	"fmt"

	api "github.com/kubedb/apimachinery/apis/kubedb/v1alpha1"
	amv "github.com/kubedb/apimachinery/pkg/validator"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/kubernetes"
)

var (
	redisVersions = sets.NewString("4", "4.0", "4.0.6")
)

func ValidateRedis(client kubernetes.Interface, redis *api.Redis) error {
	if redis.Spec.Version == "" {
		return fmt.Errorf(`object 'Version' is missing in '%v'`, redis.Spec)
	}

	// check Redis version validation
	if !redisVersions.Has(string(redis.Spec.Version)) {
		return fmt.Errorf(`KubeDB doesn't support Redis version: %s`, string(redis.Spec.Version))
	}

	if redis.Spec.Storage != nil {
		var err error
		if err = amv.ValidateStorage(client, redis.Spec.Storage); err != nil {
			return err
		}
	}

	monitorSpec := redis.Spec.Monitor
	if monitorSpec != nil {
		if err := amv.ValidateMonitorSpec(monitorSpec); err != nil {
			return err
		}

	}
	return nil
}
