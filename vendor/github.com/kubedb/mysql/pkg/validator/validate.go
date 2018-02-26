package validator

import (
	"fmt"

	"github.com/appscode/go/types"
	meta_util "github.com/appscode/kutil/meta"
	api "github.com/kubedb/apimachinery/apis/kubedb/v1alpha1"
	cs "github.com/kubedb/apimachinery/client/clientset/versioned/typed/kubedb/v1alpha1"
	amv "github.com/kubedb/apimachinery/pkg/validator"
	"github.com/pkg/errors"
	core "k8s.io/api/core/v1"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/kubernetes"
)

var (
	mysqlVersions = sets.NewString("8.0", "8")
)

func ValidateMySQL(client kubernetes.Interface, extClient cs.KubedbV1alpha1Interface, mysql *api.MySQL) error {
	if mysql.Spec.Version == "" {
		return fmt.Errorf(`object 'Version' is missing in '%v'`, mysql.Spec)
	}

	// check MySQL version validation
	if !mysqlVersions.Has(string(mysql.Spec.Version)) {
		return fmt.Errorf(`KubeDB doesn't support MySQL version: %s`, string(mysql.Spec.Version))
	}

	if mysql.Spec.Replicas != nil {
		replicas := types.Int32(mysql.Spec.Replicas)
		if replicas != 1 {
			return fmt.Errorf(`spec.replicas "%d" invalid. Value must be one`, replicas)
		}
	}

	if err := matchWithDormantDatabase(extClient, mysql); err != nil {
		return err
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

func matchWithDormantDatabase(extClient cs.KubedbV1alpha1Interface, mysql *api.MySQL) error {
	// Check if DormantDatabase exists or not
	dormantDb, err := extClient.DormantDatabases(mysql.Namespace).Get(mysql.Name, metav1.GetOptions{})
	if err != nil {
		if !kerr.IsNotFound(err) {
			return err
		}
		return nil
	}

	// Check DatabaseKind
	if dormantDb.Labels[api.LabelDatabaseKind] != api.ResourceKindMySQL {
		return fmt.Errorf(`invalid MySQL: "%v". Exists DormantDatabase "%v" of different Kind`, mysql.Name, dormantDb.Name)
	}

	// Check Origin Spec
	drmnOriginSpec := dormantDb.Spec.Origin.Spec.MySQL
	originalSpec := mysql.Spec

	if originalSpec.DatabaseSecret == nil {
		originalSpec.DatabaseSecret = &core.SecretVolumeSource{
			SecretName: mysql.Name + "-auth",
		}
	}

	// Skip checking doNotPause
	drmnOriginSpec.DoNotPause = originalSpec.DoNotPause

	if !meta_util.Equal(drmnOriginSpec, &originalSpec) {
		return errors.New("mysql spec mismatches with OriginSpec in DormantDatabases")
	}

	return nil
}
