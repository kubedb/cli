package validator

import (
	"fmt"

	tapi "github.com/k8sdb/apimachinery/apis/kubedb/v1alpha1"
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

	if postgres.Spec.Storage != nil {
		var err error
		if err = amv.ValidateStorage(client, postgres.Spec.Storage); err != nil {
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
