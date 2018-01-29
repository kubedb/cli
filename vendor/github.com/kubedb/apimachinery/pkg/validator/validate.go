package validator

import (
	"encoding/json"
	"errors"
	"fmt"

	mona "github.com/appscode/kube-mon/api"
	api "github.com/kubedb/apimachinery/apis/kubedb/v1alpha1"
	"github.com/kubedb/apimachinery/pkg/storage"
	core "k8s.io/api/core/v1"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func ValidateStorage(client kubernetes.Interface, spec *core.PersistentVolumeClaimSpec) error {
	if spec == nil {
		return nil
	}

	if spec.StorageClassName != nil {
		if _, err := client.StorageV1beta1().StorageClasses().Get(*spec.StorageClassName, metav1.GetOptions{}); err != nil {
			if kerr.IsNotFound(err) {
				return fmt.Errorf(`spec.storage.storageClassName "%v" not found`, *spec.StorageClassName)
			}
			return err
		}
	}

	if val, found := spec.Resources.Requests[core.ResourceStorage]; found {
		if val.Value() <= 0 {
			return errors.New("invalid ResourceStorage request")
		}
	} else {
		return errors.New("missing ResourceStorage request")
	}

	return nil
}

func ValidateBackupSchedule(client kubernetes.Interface, spec *api.BackupScheduleSpec, namespace string) error {
	if spec == nil {
		return nil
	}
	// CronExpression can't be empty
	if spec.CronExpression == "" {
		return errors.New("invalid cron expression")
	}

	return ValidateSnapshotSpec(client, spec.SnapshotStorageSpec, namespace)
}

func ValidateSnapshotSpec(client kubernetes.Interface, spec api.SnapshotStorageSpec, namespace string) error {
	// BucketName can't be empty
	if spec.S3 == nil && spec.GCS == nil && spec.Azure == nil && spec.Swift == nil && spec.Local == nil {
		return errors.New("no storage provider is configured")
	}

	if spec.Local != nil {
		return nil
	}

	// Need to provide Storage credential secret
	if spec.StorageSecretName == "" {
		return fmt.Errorf(`object 'SecretName' is missing in '%v'`, spec)
	}

	if err := storage.CheckBucketAccess(client, spec, namespace); err != nil {
		return err
	}

	return nil
}

func ValidateMonitorSpec(monitorSpec *mona.AgentSpec) error {
	specData, err := json.Marshal(monitorSpec)
	if err != nil {
		return err
	}

	if monitorSpec.Agent == "" {
		return fmt.Errorf(`object 'Agent' is missing in '%v'`, string(specData))
	}

	if monitorSpec.Agent.Vendor() == mona.VendorPrometheus {
		if monitorSpec.Agent == mona.AgentPrometheusBuiltin ||
			(monitorSpec.Agent == mona.AgentCoreOSPrometheus && monitorSpec.Prometheus != nil) {
			return nil
		}
	}

	return fmt.Errorf(`invalid 'Agent' in '%v'`, string(specData))
}
