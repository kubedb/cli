package validator

import (
	"errors"
	"fmt"
	"strings"

	api "github.com/kubedb/apimachinery/apis/kubedb/v1alpha1"
	adr "github.com/kubedb/apimachinery/pkg/docker"
	"github.com/kubedb/apimachinery/pkg/storage"
	amv "github.com/kubedb/apimachinery/pkg/validator"
	dr "github.com/kubedb/postgres/pkg/docker"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func ValidatePostgres(client kubernetes.Interface, postgres *api.Postgres, docker *dr.Docker) error {
	if postgres.Spec.Version == "" {
		return fmt.Errorf(`object 'Version' is missing in '%v'`, postgres.Spec)
	}

	if docker != nil {
		if err := adr.CheckDockerImageVersion(docker.GetImage(postgres), string(postgres.Spec.Version)); err != nil {
			return fmt.Errorf(`image %s not found`, docker.GetImageWithTag(postgres))
		}
	}

	if postgres.Spec.Storage != nil {
		var err error
		if err = amv.ValidateStorage(client, postgres.Spec.Storage); err != nil {
			return err
		}
	}
	if postgres.Spec.Standby != "" {
		if strings.ToLower(string(postgres.Spec.Standby)) != "hot" &&
			strings.ToLower(string(postgres.Spec.Standby)) != "warm" {
			return fmt.Errorf(`configuration.Standby "%v" invalid`, postgres.Spec.Standby)
		}
	}
	if postgres.Spec.Streaming != "" {
		// TODO: synchronous Streaming is unavailable due to lack of support
		if /*strings.ToLower(configuration.Streaming) != "synchronous" &&
		 */strings.ToLower(string(postgres.Spec.Streaming)) != "asynchronous" {
			return fmt.Errorf(`configuration.Streaming "%v" invalid`, postgres.Spec.Streaming)
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
	return nil
}
