package validator

import (
	"errors"
	"fmt"

	"github.com/appscode/go/types"
	api "github.com/kubedb/apimachinery/apis/kubedb/v1alpha1"
	"github.com/kubedb/apimachinery/pkg/storage"
	amv "github.com/kubedb/apimachinery/pkg/validator"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/kubernetes"
)

var (
	postgresVersions = sets.NewString("9.6", "9.6.6", "10.2")
)

func ValidatePostgres(client kubernetes.Interface, postgres *api.Postgres) error {
	if postgres.Spec.Version == "" {
		return fmt.Errorf(`object 'Version' is missing in '%v'`, postgres.Spec)
	}

	// check Postgres version validation
	if !postgresVersions.Has(string(postgres.Spec.Version)) {
		return fmt.Errorf(`KubeDB doesn't support Postgres version: %s`, string(postgres.Spec.Version))
	}

	if postgres.Spec.Replicas != nil {
		replicas := types.Int32(postgres.Spec.Replicas)
		if replicas < 1 {
			return fmt.Errorf(`spec.replicas "%d" invalid`, replicas)
		}
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
