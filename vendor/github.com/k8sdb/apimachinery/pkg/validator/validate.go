package validator

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/appscode/log"
	"github.com/ghodss/yaml"
	"github.com/graymeta/stow"
	tapi "github.com/k8sdb/apimachinery/api"
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
	bucketName := spec.BucketName
	if bucketName == "" {
		return fmt.Errorf(`Object 'BucketName' is missing in '%v'`, spec)
	}

	// Need to provide Storage credential secret
	storageSecret := spec.StorageSecret
	if storageSecret == nil {
		return fmt.Errorf(`Object 'StorageSecret' is missing in '%v'`, spec)
	}

	// Credential SecretName  can't be empty
	storageSecretName := storageSecret.SecretName
	if storageSecretName == "" {
		return fmt.Errorf(`Object 'SecretName' is missing in '%v'`, *spec.StorageSecret)
	}
	return nil
}

const (
	KeyProvider = "provider"
	KeyConfig   = "config"
)

func CheckBucketAccess(client clientset.Interface, snapshotSpec tapi.SnapshotStorageSpec, namespace string) error {
	secret, err := client.CoreV1().Secrets(namespace).Get(snapshotSpec.StorageSecret.SecretName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	provider := secret.Data[KeyProvider]
	if provider == nil {
		return errors.New("Missing provider key")
	}
	configData := secret.Data[KeyConfig]
	if configData == nil {
		return errors.New("Missing config key")
	}

	var config stow.ConfigMap
	if err := yaml.Unmarshal(configData, &config); err != nil {
		return err
	}

	loc, err := stow.Dial(string(provider), config)
	if err != nil {
		return err
	}

	container, err := loc.Container(snapshotSpec.BucketName)
	if err != nil {
		return err
	}

	r := bytes.NewReader([]byte("CheckBucketAccess"))
	item, err := container.Put(".kubedb", r, r.Size(), nil)
	if err != nil {
		return err
	}

	if err := container.RemoveItem(item.ID()); err != nil {
		return err
	}
	return nil
}

func ValidateSnapshot(client clientset.Interface, snapshot *tapi.Snapshot) error {
	snapshotSpec := snapshot.Spec.SnapshotStorageSpec
	if err := ValidateSnapshotSpec(snapshotSpec); err != nil {
		return err
	}

	if err := CheckBucketAccess(client, snapshot.Spec.SnapshotStorageSpec, snapshot.Namespace); err != nil {
		return err
	}
	return nil
}
