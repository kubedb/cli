package admission

import (
	"fmt"
	"sync"

	"github.com/appscode/go/log"
	"github.com/appscode/go/types"
	mon_api "github.com/appscode/kube-mon/api"
	hookapi "github.com/appscode/kubernetes-webhook-util/admission/v1beta1"
	"github.com/appscode/kutil"
	core_util "github.com/appscode/kutil/core/v1"
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

type PostgresMutator struct {
	client      kubernetes.Interface
	extClient   cs.Interface
	lock        sync.RWMutex
	initialized bool
}

var _ hookapi.AdmissionHook = &PostgresMutator{}

func (a *PostgresMutator) Resource() (plural schema.GroupVersionResource, singular string) {
	return schema.GroupVersionResource{
			Group:    "admission.kubedb.com",
			Version:  "v1alpha1",
			Resource: "postgresmutationreviews",
		},
		"postgresmutationreview"
}

func (a *PostgresMutator) Initialize(config *rest.Config, stopCh <-chan struct{}) error {
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

func (a *PostgresMutator) Admit(req *admission.AdmissionRequest) *admission.AdmissionResponse {
	status := &admission.AdmissionResponse{}

	// N.B.: No Mutating for delete
	if (req.Operation != admission.Create && req.Operation != admission.Update) ||
		len(req.SubResource) != 0 ||
		req.Kind.Group != api.SchemeGroupVersion.Group ||
		req.Kind.Kind != api.ResourceKindPostgres {
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
	dbMod, err := setDefaultValues(a.client, a.extClient, obj.(*api.Postgres).DeepCopy())
	if err != nil {
		return hookapi.StatusForbidden(err)
	} else if dbMod != nil {
		patch, err := meta_util.CreateJSONPatch(obj, dbMod)
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

// setDefaultValues provides the defaulting that is performed in mutating stage of creating/updating a Postgres database
func setDefaultValues(client kubernetes.Interface, extClient cs.Interface, postgres *api.Postgres) (runtime.Object, error) {
	if postgres.Spec.Version == "" {
		return nil, fmt.Errorf(`object 'Version' is missing in '%v'`, postgres.Spec)
	}

	if postgres.Spec.Replicas == nil {
		postgres.Spec.Replicas = types.Int32P(1)
	}

	if err := setDefaultsFromDormantDB(extClient, postgres); err != nil {
		return nil, err
	}

	// If monitoring spec is given without port,
	// set default Listening port
	setMonitoringPort(postgres)

	return postgres, nil
}

// setDefaultsFromDormantDB takes values from Similar Dormant Database
func setDefaultsFromDormantDB(extClient cs.Interface, postgres *api.Postgres) error {
	// Check if DormantDatabase exists or not
	dormantDb, err := extClient.KubedbV1alpha1().DormantDatabases(postgres.Namespace).Get(postgres.Name, metav1.GetOptions{})
	if err != nil {
		if !kerr.IsNotFound(err) {
			return err
		}
		return nil
	}

	// Check DatabaseKind
	if value, _ := meta_util.GetStringValue(dormantDb.Labels, api.LabelDatabaseKind); value != api.ResourceKindPostgres {
		return errors.New(fmt.Sprintf(`invalid Postgres: "%v". Exists DormantDatabase "%v" of different Kind`, postgres.Name, dormantDb.Name))
	}

	// Check Origin Spec
	ddbOriginSpec := dormantDb.Spec.Origin.Spec.Postgres

	// If DatabaseSecret of new object is not given,
	// Take dormantDatabaseSecretName
	if postgres.Spec.DatabaseSecret == nil {
		postgres.Spec.DatabaseSecret = ddbOriginSpec.DatabaseSecret
	}

	// If Monitoring Spec of new object is not given,
	// Take Monitoring Settings from Dormant
	if postgres.Spec.Monitor == nil {
		postgres.Spec.Monitor = ddbOriginSpec.Monitor
	} else {
		ddbOriginSpec.Monitor = postgres.Spec.Monitor
	}

	// If Backup Scheduler of new object is not given,
	// Take Backup Scheduler Settings from Dormant
	if postgres.Spec.BackupSchedule == nil {
		postgres.Spec.BackupSchedule = ddbOriginSpec.BackupSchedule
	} else {
		ddbOriginSpec.BackupSchedule = postgres.Spec.BackupSchedule
	}

	// Skip checking DoNotPause
	ddbOriginSpec.DoNotPause = postgres.Spec.DoNotPause

	if !meta_util.Equal(ddbOriginSpec, &postgres.Spec) {
		diff := meta_util.Diff(ddbOriginSpec, &postgres.Spec)
		log.Errorf("postgres spec mismatches with OriginSpec in DormantDatabases. Diff: %v", diff)
		return errors.New(fmt.Sprintf("postgres spec mismatches with OriginSpec in DormantDatabases. Diff: %v", diff))
	}

	if _, err := meta_util.GetString(postgres.Annotations, api.AnnotationInitialized); err == kutil.ErrNotFound &&
		postgres.Spec.Init != nil &&
		postgres.Spec.Init.SnapshotSource != nil {
		postgres.Annotations = core_util.UpsertMap(postgres.Annotations, map[string]string{
			api.AnnotationInitialized: "",
		})
	}

	// Delete  Matching dormantDatabase in Controller

	return nil
}

// Assign Default Monitoring Port if MonitoringSpec Exists
// and the AgentVendor is Prometheus.
func setMonitoringPort(postgres *api.Postgres) {
	if postgres.Spec.Monitor != nil &&
		postgres.GetMonitoringVendor() == mon_api.VendorPrometheus {
		if postgres.Spec.Monitor.Prometheus == nil {
			postgres.Spec.Monitor.Prometheus = &mon_api.PrometheusSpec{}
		}
		if postgres.Spec.Monitor.Prometheus.Port == 0 {
			postgres.Spec.Monitor.Prometheus.Port = api.PrometheusExporterPortNumber
		}
	}
}
