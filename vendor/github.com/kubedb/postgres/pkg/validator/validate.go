package validator

import (
	"errors"
	"fmt"

	"github.com/appscode/go/types"
	meta_util "github.com/appscode/kutil/meta"
	api "github.com/kubedb/apimachinery/apis/kubedb/v1alpha1"
	cs "github.com/kubedb/apimachinery/client/clientset/versioned/typed/kubedb/v1alpha1"
	"github.com/kubedb/apimachinery/pkg/storage"
	amv "github.com/kubedb/apimachinery/pkg/validator"
	core "k8s.io/api/core/v1"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/kubernetes"
)

var (
	postgresVersions = sets.NewString("9.6", "9.6.7", "10.2")
)

func ValidatePostgres(client kubernetes.Interface, extClient cs.KubedbV1alpha1Interface, postgres *api.Postgres) error {

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

	if err := matchWithDormantDatabase(extClient, postgres); err != nil {
		return err
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

func matchWithDormantDatabase(extClient cs.KubedbV1alpha1Interface, postgres *api.Postgres) error {
	// Check if DormantDatabase exists or not
	dormantDb, err := extClient.DormantDatabases(postgres.Namespace).Get(postgres.Name, metav1.GetOptions{})
	if err != nil {
		if !kerr.IsNotFound(err) {
			return err
		}
		return nil
	}

	// Check DatabaseKind
	if dormantDb.Labels[api.LabelDatabaseKind] != api.ResourceKindPostgres {
		return fmt.Errorf(`invalid Postgres: "%v". Exists DormantDatabase "%v" of different Kind`, postgres.Name, dormantDb.Name)
	}

	// Check Origin Spec
	drmnOriginSpec := dormantDb.Spec.Origin.Spec.Postgres
	originalSpec := postgres.Spec

	if originalSpec.DatabaseSecret == nil {
		originalSpec.DatabaseSecret = &core.SecretVolumeSource{
			SecretName: postgres.OffshootName() + "-auth",
		}
	}

	// Skip checking doNotPause
	drmnOriginSpec.DoNotPause = originalSpec.DoNotPause

	if !meta_util.Equal(drmnOriginSpec, &originalSpec) {
		return errors.New("object spec in Postgres mismatches with OriginSpec in DormantDatabases")
	}

	return nil
}
