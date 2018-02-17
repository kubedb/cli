package validator

import (
	"fmt"

	api "github.com/kubedb/apimachinery/apis/kubedb/v1alpha1"
	amv "github.com/kubedb/apimachinery/pkg/validator"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/kubernetes"
)

var (
	mysqlVersions = sets.NewString("8.0", "8")
)

func ValidateMySQL(client kubernetes.Interface, mysql *api.MySQL) error {
	if mysql.Spec.Version == "" {
		return fmt.Errorf(`object 'Version' is missing in '%v'`, mysql.Spec)
	}

	// check MySQL version validation
	if !mysqlVersions.Has(string(mysql.Spec.Version)) {
		return fmt.Errorf(`KubeDB doesn't support MySQL version: %s`, string(mysql.Spec.Version))
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
