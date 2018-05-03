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

type RedisMutator struct {
	client      kubernetes.Interface
	extClient   cs.Interface
	lock        sync.RWMutex
	initialized bool
}

var _ hookapi.AdmissionHook = &RedisMutator{}

func (a *RedisMutator) Resource() (plural schema.GroupVersionResource, singular string) {
	return schema.GroupVersionResource{
			Group:    "admission.kubedb.com",
			Version:  "v1alpha1",
			Resource: "redismutationreviews",
		},
		"redismutationreview"
}

func (a *RedisMutator) Initialize(config *rest.Config, stopCh <-chan struct{}) error {
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

func (a *RedisMutator) Admit(req *admission.AdmissionRequest) *admission.AdmissionResponse {
	status := &admission.AdmissionResponse{}

	// N.B.: No Mutating for delete
	if (req.Operation != admission.Create && req.Operation != admission.Update) ||
		len(req.SubResource) != 0 ||
		req.Kind.Group != api.SchemeGroupVersion.Group ||
		req.Kind.Kind != api.ResourceKindRedis {
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
	redisMod, err := setDefaultValues(a.client, a.extClient, obj.(*api.Redis).DeepCopy())
	if err != nil {
		return hookapi.StatusForbidden(err)
	} else if redisMod != nil {
		patch, err := meta_util.CreateJSONPatch(obj, redisMod)
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

// setDefaultValues provides the defaulting that is performed in mutating stage of creating/updating a Redis database
func setDefaultValues(client kubernetes.Interface, extClient cs.Interface, redis *api.Redis) (runtime.Object, error) {
	if redis.Spec.Version == "" {
		return nil, fmt.Errorf(`object 'Version' is missing in '%v'`, redis.Spec)
	}

	if redis.Spec.Replicas == nil {
		redis.Spec.Replicas = types.Int32P(1)
	}

	if err := setDefaultsFromDormantDB(extClient, redis); err != nil {
		return nil, err
	}

	// If monitoring spec is given without port,
	// set default Listening port
	setMonitoringPort(redis)

	return redis, nil
}

// setDefaultsFromDormantDB takes values from Similar Dormant Database
func setDefaultsFromDormantDB(extClient cs.Interface, redis *api.Redis) error {
	// Check if DormantDatabase exists or not
	dormantDb, err := extClient.KubedbV1alpha1().DormantDatabases(redis.Namespace).Get(redis.Name, metav1.GetOptions{})
	if err != nil {
		if !kerr.IsNotFound(err) {
			return err
		}
		return nil
	}

	// Check DatabaseKind
	if value, _ := meta_util.GetStringValue(dormantDb.Labels, api.LabelDatabaseKind); value != api.ResourceKindRedis {
		return errors.New(fmt.Sprintf(`invalid Redis: "%v". Exists DormantDatabase "%v" of different Kind`, redis.Name, dormantDb.Name))
	}

	// Check Origin Spec
	ddbOriginSpec := dormantDb.Spec.Origin.Spec.Redis

	// If Monitoring Spec of new object is not given,
	// Take Monitoring Settings from Dormant
	if redis.Spec.Monitor == nil {
		redis.Spec.Monitor = ddbOriginSpec.Monitor
	} else {
		ddbOriginSpec.Monitor = redis.Spec.Monitor
	}

	// Skip checking DoNotPause
	ddbOriginSpec.DoNotPause = redis.Spec.DoNotPause

	if !meta_util.Equal(ddbOriginSpec, &redis.Spec) {
		diff := meta_util.Diff(ddbOriginSpec, &redis.Spec)
		log.Errorf("redis spec mismatches with OriginSpec in DormantDatabases. Diff: %v", diff)
		return errors.New(fmt.Sprintf("redis spec mismatches with OriginSpec in DormantDatabases. Diff: %v", diff))
	}

	// Delete  Matching dormantDatabase in Controller

	return nil
}

// Assign Default Monitoring Port if MonitoringSpec Exists
// and the AgentVendor is Prometheus.
func setMonitoringPort(redis *api.Redis) {
	if redis.Spec.Monitor != nil &&
		redis.GetMonitoringVendor() == mon_api.VendorPrometheus {
		if redis.Spec.Monitor.Prometheus == nil {
			redis.Spec.Monitor.Prometheus = &mon_api.PrometheusSpec{}
		}
		if redis.Spec.Monitor.Prometheus.Port == 0 {
			redis.Spec.Monitor.Prometheus.Port = api.PrometheusExporterPortNumber
		}
	}
}
