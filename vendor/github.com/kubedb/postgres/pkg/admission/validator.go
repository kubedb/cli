package admission

import (
	"fmt"
	"strings"
	"sync"

	"github.com/appscode/go/log"
	hookapi "github.com/appscode/kubernetes-webhook-util/admission/v1beta1"
	meta_util "github.com/appscode/kutil/meta"
	api "github.com/kubedb/apimachinery/apis/kubedb/v1alpha1"
	cs "github.com/kubedb/apimachinery/client/clientset/versioned"
	kubedbv1alpha1 "github.com/kubedb/apimachinery/client/clientset/versioned/typed/kubedb/v1alpha1"
	"github.com/kubedb/apimachinery/pkg/storage"
	amv "github.com/kubedb/apimachinery/pkg/validator"
	"github.com/pkg/errors"
	admission "k8s.io/api/admission/v1beta1"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/mergepatch"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type PostgresValidator struct {
	client      kubernetes.Interface
	extClient   cs.Interface
	lock        sync.RWMutex
	initialized bool
}

var _ hookapi.AdmissionHook = &PostgresValidator{}

func (a *PostgresValidator) Resource() (plural schema.GroupVersionResource, singular string) {
	return schema.GroupVersionResource{
			Group:    "admission.kubedb.com",
			Version:  "v1alpha1",
			Resource: "postgresvalidationreviews",
		},
		"postgresvalidationreview"
}

func (a *PostgresValidator) Initialize(config *rest.Config, stopCh <-chan struct{}) error {
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

func (a *PostgresValidator) Admit(req *admission.AdmissionRequest) *admission.AdmissionResponse {
	status := &admission.AdmissionResponse{}

	if (req.Operation != admission.Create && req.Operation != admission.Update && req.Operation != admission.Delete) ||
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

	switch req.Operation {
	case admission.Delete:
		// req.Object.Raw = nil, so read from kubernetes
		obj, err := a.extClient.KubedbV1alpha1().Postgreses(req.Namespace).Get(req.Name, metav1.GetOptions{})
		if err != nil && !kerr.IsNotFound(err) {
			return hookapi.StatusInternalServerError(err)
		} else if err == nil && obj.Spec.DoNotPause {
			return hookapi.StatusBadRequest(fmt.Errorf(`postgres "%s" can't be paused. To continue delete, unset spec.doNotPause and retry`, req.Name))
		}
	default:
		obj, err := meta_util.UnmarshalFromJSON(req.Object.Raw, api.SchemeGroupVersion)
		if err != nil {
			return hookapi.StatusBadRequest(err)
		}
		if req.Operation == admission.Update {
			// validate changes made by user
			oldObject, err := meta_util.UnmarshalFromJSON(req.OldObject.Raw, api.SchemeGroupVersion)
			if err != nil {
				return hookapi.StatusBadRequest(err)
			}

			postgres := obj.(*api.Postgres).DeepCopy()
			oldPostgres := oldObject.(*api.Postgres).DeepCopy()
			// Allow changing Database Secret only if there was no secret have set up yet.
			if oldPostgres.Spec.DatabaseSecret == nil {
				oldPostgres.Spec.DatabaseSecret = postgres.Spec.DatabaseSecret
			}

			if err := validateUpdate(postgres, oldPostgres, req.Kind.Kind); err != nil {
				return hookapi.StatusBadRequest(fmt.Errorf("%v", err))
			}
		}
		// validate database specs
		if err = ValidatePostgres(a.client, a.extClient.KubedbV1alpha1(), obj.(*api.Postgres)); err != nil {
			return hookapi.StatusForbidden(err)
		}
	}
	status.Allowed = true
	return status
}

var (
	postgresVersions = sets.NewString("9.6", "9.6.7", "10.2")
)

// ValidatePostgres checks if the object satisfies all the requirements.
// It is not method of Interface, because it is referenced from controller package too.
func ValidatePostgres(client kubernetes.Interface, extClient kubedbv1alpha1.KubedbV1alpha1Interface, postgres *api.Postgres) error {
	if postgres.Spec.Version == "" {
		return fmt.Errorf(`object 'Version' is missing in '%v'`, postgres.Spec)
	}

	// Check Postgres version validation
	if !postgresVersions.Has(string(postgres.Spec.Version)) {
		return fmt.Errorf(`KubeDB doesn't support Postgres version: %s`, string(postgres.Spec.Version))
	}

	if postgres.Spec.Replicas == nil || *postgres.Spec.Replicas < 1 {
		return fmt.Errorf(`spec.replicas "%v" invalid. Value must be greater than zero`, postgres.Spec.Replicas)
	}

	if postgres.Spec.Storage != nil {
		var err error
		if err = amv.ValidateStorage(client, postgres.Spec.Storage); err != nil {
			return err
		}
	}

	if postgres.Spec.StandbyMode != nil {
		standByMode := *postgres.Spec.StandbyMode
		if standByMode != api.HotStandby && standByMode != api.WarmStandby {
			return fmt.Errorf(`spec.standbyMode "%s" invalid`, standByMode)
		}
	}

	if postgres.Spec.StreamingMode != nil {
		streamingMode := *postgres.Spec.StreamingMode
		// TODO: synchronous Streaming is unavailable due to lack of support
		if streamingMode != api.AsynchronousStreaming {
			return fmt.Errorf(`spec.streamingMode "%s" invalid`, streamingMode)
		}
	}

	if postgres.Spec.Archiver != nil {
		archiverStorage := postgres.Spec.Archiver.Storage
		if archiverStorage != nil {
			if archiverStorage.StorageSecretName == "" {
				return fmt.Errorf(`object 'StorageSecretName' is missing in '%v'`, archiverStorage)
			}
			if archiverStorage.S3 == nil {
				return errors.New("no storage provider is configured")
			}
			if !(archiverStorage.GCS == nil && archiverStorage.Azure == nil && archiverStorage.Swift == nil && archiverStorage.Local == nil) {
				return errors.New("invalid storage provider is configured")
			}

			if err := storage.CheckBucketAccess(client, *archiverStorage, postgres.Namespace); err != nil {
				return err
			}
		}
	}

	databaseSecret := postgres.Spec.DatabaseSecret
	if databaseSecret != nil {
		if _, err := client.CoreV1().Secrets(postgres.Namespace).Get(databaseSecret.SecretName, metav1.GetOptions{}); err != nil {
			return err
		}
	}

	if postgres.Spec.Init != nil &&
		postgres.Spec.Init.SnapshotSource != nil &&
		databaseSecret == nil {
		return fmt.Errorf("in Snapshot init, 'spec.databaseSecret.secretName' of %v needs to be similar to older database of snapshot %v",
			postgres.Name, postgres.Spec.Init.SnapshotSource.Name)
	}

	if postgres.Spec.Init != nil && postgres.Spec.Init.PostgresWAL != nil {
		wal := postgres.Spec.Init.PostgresWAL
		if wal.StorageSecretName == "" {
			return fmt.Errorf(`object 'StorageSecretName' is missing in '%v'`, wal)
		}
		if wal.S3 == nil {
			return errors.New("no storage provider is configured")
		}
		if !(wal.GCS == nil && wal.Azure == nil && wal.Swift == nil && wal.Local == nil) {
			return errors.New("invalid storage provider is configured")
		}

		if err := storage.CheckBucketAccess(client, wal.SnapshotStorageSpec, postgres.Namespace); err != nil {
			return err
		}
	}

	backupScheduleSpec := postgres.Spec.BackupSchedule
	if backupScheduleSpec != nil {
		if err := amv.ValidateBackupSchedule(client, backupScheduleSpec, postgres.Namespace); err != nil {
			return err
		}
	}

	monitorSpec := postgres.Spec.Monitor
	if monitorSpec != nil {
		if err := amv.ValidateMonitorSpec(monitorSpec); err != nil {
			return err
		}
	}

	if err := matchWithDormantDatabase(extClient, postgres); err != nil {
		return err
	}
	return nil
}

func matchWithDormantDatabase(extClient kubedbv1alpha1.KubedbV1alpha1Interface, postgres *api.Postgres) error {
	// Check if DormantDatabase exists or not
	dormantDb, err := extClient.DormantDatabases(postgres.Namespace).Get(postgres.Name, metav1.GetOptions{})
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
	drmnOriginSpec := dormantDb.Spec.Origin.Spec.Postgres
	originalSpec := postgres.Spec

	// Skip checking doNotPause
	drmnOriginSpec.DoNotPause = originalSpec.DoNotPause

	// Skip checking Monitoring
	drmnOriginSpec.Monitor = originalSpec.Monitor

	// Skip Checking BackUP Scheduler
	drmnOriginSpec.BackupSchedule = originalSpec.BackupSchedule

	if !meta_util.Equal(drmnOriginSpec, &originalSpec) {
		diff := meta_util.Diff(drmnOriginSpec, &originalSpec)
		log.Errorf("postgres spec mismatches with OriginSpec in DormantDatabases. Diff: %v", diff)
		return errors.New(fmt.Sprintf("postgres spec mismatches with OriginSpec in DormantDatabases. Diff: %v", diff))
	}

	return nil
}

func validateUpdate(obj, oldObj runtime.Object, kind string) error {
	preconditions := getPreconditionFunc()
	_, err := meta_util.CreateStrategicPatch(oldObj, obj, preconditions...)
	if err != nil {
		if mergepatch.IsPreconditionFailed(err) {
			return fmt.Errorf("%v.%v", err, preconditionFailedError(kind))
		}
		return err
	}
	return nil
}

func getPreconditionFunc() []mergepatch.PreconditionFunc {
	preconditions := []mergepatch.PreconditionFunc{
		mergepatch.RequireKeyUnchanged("apiVersion"),
		mergepatch.RequireKeyUnchanged("kind"),
		mergepatch.RequireMetadataKeyUnchanged("name"),
		mergepatch.RequireMetadataKeyUnchanged("namespace"),
	}

	for _, field := range preconditionSpecFields {
		preconditions = append(preconditions,
			meta_util.RequireChainKeyUnchanged(field),
		)
	}
	return preconditions
}

var preconditionSpecFields = []string{
	"spec.version",
	"spec.standby",
	"spec.streaming",
	"spec.archiver",
	"spec.databaseSecret",
	"spec.storage",
	"spec.nodeSelector",
	"spec.init",
}

func preconditionFailedError(kind string) error {
	str := preconditionSpecFields
	strList := strings.Join(str, "\n\t")
	return fmt.Errorf(strings.Join([]string{`At least one of the following was changed:
	apiVersion
	kind
	name
	namespace`, strList}, "\n\t"))
}
