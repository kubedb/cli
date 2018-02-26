package validator

import (
	"fmt"

	"github.com/appscode/go/types"
	meta_util "github.com/appscode/kutil/meta"
	api "github.com/kubedb/apimachinery/apis/kubedb/v1alpha1"
	cs "github.com/kubedb/apimachinery/client/clientset/versioned/typed/kubedb/v1alpha1"
	amv "github.com/kubedb/apimachinery/pkg/validator"
	"github.com/pkg/errors"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/kubernetes"
)

var (
	redisVersions = sets.NewString("4", "4.0", "4.0.6")
)

func ValidateRedis(client kubernetes.Interface, extClient cs.KubedbV1alpha1Interface, redis *api.Redis) error {
	if redis.Spec.Version == "" {
		return fmt.Errorf(`object 'Version' is missing in '%v'`, redis.Spec)
	}

	// check Redis version validation
	if !redisVersions.Has(string(redis.Spec.Version)) {
		return fmt.Errorf(`KubeDB doesn't support Redis version: %s`, string(redis.Spec.Version))
	}

	if redis.Spec.Replicas != nil {
		replicas := types.Int32(redis.Spec.Replicas)
		if replicas != 1 {
			return fmt.Errorf(`spec.replicas "%d" invalid. Value must be one`, replicas)
		}
	}

	if err := matchWithDormantDatabase(extClient, redis); err != nil {
		return err
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

func matchWithDormantDatabase(extClient cs.KubedbV1alpha1Interface, redis *api.Redis) error {
	// Check if DormantDatabase exists or not
	dormantDb, err := extClient.DormantDatabases(redis.Namespace).Get(redis.Name, metav1.GetOptions{})
	if err != nil {
		if !kerr.IsNotFound(err) {
			return err
		}
		return nil
	}

	// Check DatabaseKind
	if dormantDb.Labels[api.LabelDatabaseKind] != api.ResourceKindRedis {
		return fmt.Errorf(`invalid Redis: "%v". Exists DormantDatabase "%v" of different Kind`, redis.Name, dormantDb.Name)
	}

	// Check Origin Spec
	drmnOriginSpec := dormantDb.Spec.Origin.Spec.Redis
	originalSpec := redis.Spec

	// Skip checking doNotPause
	drmnOriginSpec.DoNotPause = originalSpec.DoNotPause

	if !meta_util.Equal(drmnOriginSpec, &originalSpec) {
		return errors.New("redis spec mismatches with OriginSpec in DormantDatabases")
	}

	return nil
}
