package validator

import (
	"fmt"

	tapi "github.com/k8sdb/apimachinery/api"
	"github.com/k8sdb/apimachinery/pkg/docker"
	amv "github.com/k8sdb/apimachinery/pkg/validator"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
)

func ValidatePostgres(client clientset.Interface, postgres *tapi.Postgres) error {
	if postgres.Spec.Version == "" {
		return fmt.Errorf(`Object 'Version' is missing in '%v'`, postgres.Spec)
	}

	version := fmt.Sprintf("%v-db", postgres.Spec.Version)
	if err := docker.CheckDockerImageVersion(docker.ImagePostgres, version); err != nil {
		return fmt.Errorf(`Image %v:%v not found`, docker.ImagePostgres, version)
	}

	storage := postgres.Spec.Storage
	if storage != nil {
		var err error
		if _, err = amv.ValidateStorageSpec(client, storage); err != nil {
			return err
		}
	}

	databaseSecret := postgres.Spec.DatabaseSecret
	if databaseSecret != nil {
		if _, err := client.CoreV1().Secrets(postgres.Namespace).Get(databaseSecret.SecretName, metav1.GetOptions{}); err != nil {
			return err
		}
	}

	backupScheduleSpec := postgres.Spec.BackupSchedule
	if backupScheduleSpec != nil {
		if err := amv.ValidateBackupSchedule(backupScheduleSpec); err != nil {
			return err
		}

		if err := amv.CheckBucketAccess(client, backupScheduleSpec.SnapshotStorageSpec, postgres.Namespace); err != nil {
			return err
		}
	}

	monitorSpec := postgres.Spec.Monitor
	if monitorSpec != nil {
		if monitorSpec.Agent == "" {
			return fmt.Errorf(`Object 'Agent' is missing in '%v'`, monitorSpec)
		}
		if monitorSpec.Prometheus != nil {
			if monitorSpec.Agent != tapi.AgentCoreosPrometheus {
				return fmt.Errorf(`Invalid 'Agent' in '%v'`, monitorSpec)
			}
		}
	}
	return nil
}
