package validator

import (
	"fmt"

	api "github.com/kubedb/apimachinery/apis/kubedb/v1alpha1"
	adr "github.com/kubedb/apimachinery/pkg/docker"
	amv "github.com/kubedb/apimachinery/pkg/validator"
	dr "github.com/kubedb/mysql/pkg/docker"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func ValidateMySQL(client kubernetes.Interface, mysql *api.MySQL, docker *dr.Docker) error {
	if mysql.Spec.Version == "" {
		return fmt.Errorf(`object 'Version' is missing in '%v'`, mysql.Spec)
	}

	if docker != nil {
		if err := adr.CheckDockerImageVersion(docker.GetImage(mysql), string(mysql.Spec.Version)); err != nil {
			return fmt.Errorf(`image %s not found`, docker.GetImageWithTag(mysql))
		}
	}

	if mysql.Spec.Storage != nil {
		var err error
		if err = amv.ValidateStorage(client, mysql.Spec.Storage); err != nil {
			return err
		}
	}

	databaseSecret := mysql.Spec.DatabaseSecret
	if databaseSecret != nil {
		if _, err := client.CoreV1().Secrets(mysql.Namespace).Get(databaseSecret.SecretName, metav1.GetOptions{}); err != nil {
			return err
		}
	}

	backupScheduleSpec := mysql.Spec.BackupSchedule
	if backupScheduleSpec != nil {
		if err := amv.ValidateBackupSchedule(client, backupScheduleSpec, mysql.Namespace); err != nil {
			return err
		}
	}

	monitorSpec := mysql.Spec.Monitor
	if monitorSpec != nil {
		if err := amv.ValidateMonitorSpec(monitorSpec); err != nil {
			return err
		}

	}
	return nil
}
