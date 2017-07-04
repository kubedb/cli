package validator

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/appscode/log"
	tapi "github.com/k8sdb/apimachinery/api"
	"github.com/k8sdb/apimachinery/pkg/storage"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
	apiv1 "k8s.io/client-go/pkg/api/v1"
)

func ValidateStorageSpec(client clientset.Interface, spec *tapi.StorageSpec) (*tapi.StorageSpec, error) {
	if spec == nil {
		return nil, nil
	}

	if spec.Class == "" {
		return nil, fmt.Errorf(`Object 'Class' is missing in '%v'`, *spec)
	}

	if _, err := client.StorageV1().StorageClasses().Get(spec.Class, metav1.GetOptions{}); err != nil {
		if kerr.IsNotFound(err) {
			return nil, fmt.Errorf(`Spec.Storage.Class "%v" not found`, spec.Class)
		}
		return nil, err
	}

	if len(spec.AccessModes) == 0 {
		spec.AccessModes = []apiv1.PersistentVolumeAccessMode{
			apiv1.ReadWriteOnce,
		}
		log.Infof(`Using "%v" as AccessModes in "%v"`, apiv1.ReadWriteOnce, *spec)
	}

	if val, found := spec.Resources.Requests[apiv1.ResourceStorage]; found {
		if val.Value() <= 0 {
			return nil, errors.New("Invalid ResourceStorage request")
		}
	} else {
		return nil, errors.New("Missing ResourceStorage request")
	}

	return spec, nil
}

func ValidateBackupSchedule(spec *tapi.BackupScheduleSpec) error {
	if spec == nil {
		return nil
	}
	// CronExpression can't be empty
	if spec.CronExpression == "" {
		return errors.New("Invalid cron expression")
	}

	return ValidateSnapshotSpec(spec.SnapshotStorageSpec)
}

func ValidateSnapshotSpec(spec tapi.SnapshotStorageSpec) error {
	// BucketName can't be empty
	if spec.S3 == nil && spec.GCS == nil && spec.Azure == nil && spec.Swift == nil && spec.Local == nil {
		return errors.New("No storage provider is configured.")
	}

	// Need to provide Storage credential secret
	if spec.StorageSecretName == "" {
		return fmt.Errorf(`Object 'SecretName' is missing in '%v'`, spec.StorageSecretName)
	}
	return nil
}

func ValidateMonitorSpec(monitorSpec *tapi.MonitorSpec) error {
	specData, err := json.Marshal(monitorSpec)
	if err != nil {
		return err
	}

	if monitorSpec.Agent == "" {
		return fmt.Errorf(`Object 'Agent' is missing in '%v'`, string(specData))
	}
	if monitorSpec.Prometheus != nil {
		if monitorSpec.Agent != tapi.AgentCoreosPrometheus {
			return fmt.Errorf(`Invalid 'Agent' in '%v'`, string(specData))
		}
	}

	return nil
}

func ValidateSnapshot(client clientset.Interface, snapshot *tapi.Snapshot) error {
	snapshotSpec := snapshot.Spec.SnapshotStorageSpec
	if err := ValidateSnapshotSpec(snapshotSpec); err != nil {
		return err
	}

	if err := storage.CheckBucketAccess(client, snapshot.Spec.SnapshotStorageSpec, snapshot.Namespace); err != nil {
		return err
	}
	return nil
}
