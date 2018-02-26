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
	memcachedVersions = sets.NewString("1.5", "1.5.4")
)

func ValidateMemcached(client kubernetes.Interface, extClient cs.KubedbV1alpha1Interface, memcached *api.Memcached) error {
	if memcached.Spec.Version == "" {
		return fmt.Errorf(`object 'Version' is missing in '%v'`, memcached.Spec)
	}

	// check Memcached version validation
	if !memcachedVersions.Has(string(memcached.Spec.Version)) {
		return fmt.Errorf(`KubeDB doesn't support Memcached version: %s`, string(memcached.Spec.Version))
	}

	if memcached.Spec.Replicas != nil {
		replicas := types.Int32(memcached.Spec.Replicas)
		if replicas < 1 {
			return fmt.Errorf(`spec.replicas "%d" invalid`, replicas)
		}
	}

	if err := matchWithDormantDatabase(extClient, memcached); err != nil {
		return err
	}

	monitorSpec := memcached.Spec.Monitor
	if monitorSpec != nil {
		if err := amv.ValidateMonitorSpec(monitorSpec); err != nil {
			return err
		}
	}
	return nil
}

func matchWithDormantDatabase(extClient cs.KubedbV1alpha1Interface, memcached *api.Memcached) error {
	// Check if DormantDatabase exists or not
	dormantDb, err := extClient.DormantDatabases(memcached.Namespace).Get(memcached.Name, metav1.GetOptions{})
	if err != nil {
		if !kerr.IsNotFound(err) {
			return err
		}
		return nil
	}

	// Check DatabaseKind
	if dormantDb.Labels[api.LabelDatabaseKind] != api.ResourceKindMemcached {
		return fmt.Errorf(`invalid Memcached: "%v". Exists DormantDatabase "%v" of different Kind`, memcached.Name, dormantDb.Name)
	}

	// Check Origin Spec
	drmnOriginSpec := dormantDb.Spec.Origin.Spec.Memcached
	originalSpec := memcached.Spec

	// Skip checking doNotPause
	drmnOriginSpec.DoNotPause = originalSpec.DoNotPause

	if !meta_util.Equal(drmnOriginSpec, &originalSpec) {
		return errors.New("memcached spec mismatches with OriginSpec in DormantDatabases")
	}

	return nil
}
