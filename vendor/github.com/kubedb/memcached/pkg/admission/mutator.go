package admission

import (
	"fmt"
	"sync"

	"github.com/appscode/go/log"
	"github.com/appscode/go/types"
	mon_api "github.com/appscode/kube-mon/api"
	hookapi "github.com/appscode/kubernetes-webhook-util/admission/v1beta1"
	meta_util "github.com/appscode/kutil/meta"
	api "github.com/kubedb/apimachinery/apis/kubedb/v1alpha1"
	cs "github.com/kubedb/apimachinery/client/clientset/versioned"
	"github.com/pkg/errors"
	admission "k8s.io/api/admission/v1beta1"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type MemcachedMutator struct {
	client      kubernetes.Interface
	extClient   cs.Interface
	lock        sync.RWMutex
	initialized bool
}

var _ hookapi.AdmissionHook = &MemcachedMutator{}

func (a *MemcachedMutator) Resource() (plural schema.GroupVersionResource, singular string) {
	return schema.GroupVersionResource{
			Group:    "admission.kubedb.com",
			Version:  "v1alpha1",
			Resource: "memcachedmutationreviews",
		},
		"memcachedmutationreview"
}

func (a *MemcachedMutator) Initialize(config *rest.Config, stopCh <-chan struct{}) error {
	a.lock.Lock()
	defer a.lock.Unlock()

	a.initialized = true

	var err error
	if a.client, err = kubernetes.NewForConfig(config); err != nil {
		return err
	}
	if a.extClient, err = cs.NewForConfig(config); err != nil {
		return err
	}
	return err
}

func (a *MemcachedMutator) Admit(req *admission.AdmissionRequest) *admission.AdmissionResponse {
	status := &admission.AdmissionResponse{}

	// N.B.: No Mutating for delete
	if (req.Operation != admission.Create && req.Operation != admission.Update) ||
		len(req.SubResource) != 0 ||
		req.Kind.Group != api.SchemeGroupVersion.Group ||
		req.Kind.Kind != api.ResourceKindMemcached {
		status.Allowed = true
		return status
	}

	a.lock.RLock()
	defer a.lock.RUnlock()
	if !a.initialized {
		return hookapi.StatusUninitialized()
	}
	obj, err := meta_util.UnmarshalFromJSON(req.Object.Raw, api.SchemeGroupVersion)
	if err != nil {
		return hookapi.StatusBadRequest(err)
	}
	memcachedMod, err := setDefaultValues(a.client, a.extClient, obj.(*api.Memcached).DeepCopy())
	if err != nil {
		return hookapi.StatusForbidden(err)
	} else if memcachedMod != nil {
		patch, err := meta_util.CreateJSONPatch(obj, memcachedMod)
		if err != nil {
			return hookapi.StatusInternalServerError(err)
		}
		status.Patch = patch
		patchType := admission.PatchTypeJSONPatch
		status.PatchType = &patchType
	}

	status.Allowed = true
	return status
}

// setDefaultValues provides the defaulting that is performed in mutating stage of creating/updating a Memcached database
func setDefaultValues(client kubernetes.Interface, extClient cs.Interface, memcached *api.Memcached) (runtime.Object, error) {
	if memcached.Spec.Version == "" {
		return nil, fmt.Errorf(`object 'Version' is missing in '%v'`, memcached.Spec)
	}

	if memcached.Spec.Replicas == nil {
		memcached.Spec.Replicas = types.Int32P(1)
	}

	if err := setDefaultsFromDormantDB(extClient, memcached); err != nil {
		return nil, err
	}

	// If monitoring spec is given without port,
	// set default Listening port
	setMonitoringPort(memcached)

	return memcached, nil
}

// setDefaultsFromDormantDB takes values from Similar Dormant Database
func setDefaultsFromDormantDB(extClient cs.Interface, memcached *api.Memcached) error {
	// Check if DormantDatabase exists or not
	dormantDb, err := extClient.KubedbV1alpha1().DormantDatabases(memcached.Namespace).Get(memcached.Name, metav1.GetOptions{})
	if err != nil {
		if !kerr.IsNotFound(err) {
			return err
		}
		return nil
	}

	// Check DatabaseKind
	if value, _ := meta_util.GetStringValue(dormantDb.Labels, api.LabelDatabaseKind); value != api.ResourceKindMemcached {
		return errors.New(fmt.Sprintf(`invalid Memcached: "%v". Exists DormantDatabase "%v" of different Kind`, memcached.Name, dormantDb.Name))
	}

	// Check Origin Spec
	ddbOriginSpec := dormantDb.Spec.Origin.Spec.Memcached

	// If Monitoring Spec of new object is not given,
	// Take Monitoring Settings from Dormant
	if memcached.Spec.Monitor == nil {
		memcached.Spec.Monitor = ddbOriginSpec.Monitor
	} else {
		ddbOriginSpec.Monitor = memcached.Spec.Monitor
	}

	// Skip checking DoNotPause
	ddbOriginSpec.DoNotPause = memcached.Spec.DoNotPause

	if !meta_util.Equal(ddbOriginSpec, &memcached.Spec) {
		diff := meta_util.Diff(ddbOriginSpec, &memcached.Spec)
		log.Errorf("memcached spec mismatches with OriginSpec in DormantDatabases. Diff: %v", diff)
		return errors.New(fmt.Sprintf("memcached spec mismatches with OriginSpec in DormantDatabases. Diff: %v", diff))
	}

	// Delete  Matching dormantDatabase in Controller

	return nil
}

// Assign Default Monitoring Port if MonitoringSpec Exists
// and the AgentVendor is Prometheus.
func setMonitoringPort(memcached *api.Memcached) {
	if memcached.Spec.Monitor != nil &&
		memcached.GetMonitoringVendor() == mon_api.VendorPrometheus {
		if memcached.Spec.Monitor.Prometheus == nil {
			memcached.Spec.Monitor.Prometheus = &mon_api.PrometheusSpec{}
		}
		if memcached.Spec.Monitor.Prometheus.Port == 0 {
			memcached.Spec.Monitor.Prometheus.Port = api.PrometheusExporterPortNumber
		}
	}
}
