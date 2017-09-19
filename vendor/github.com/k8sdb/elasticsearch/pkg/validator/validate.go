package validator

import (
	"fmt"

	tapi "github.com/k8sdb/apimachinery/apis/kubedb/v1alpha1"
	"github.com/k8sdb/apimachinery/pkg/docker"
	amv "github.com/k8sdb/apimachinery/pkg/validator"
	clientset "k8s.io/client-go/kubernetes"
)

func ValidateElastic(client clientset.Interface, elastic *tapi.Elasticsearch) error {
	if elastic.Spec.Version == "" {
		return fmt.Errorf(`Object 'Version' is missing in '%v'`, elastic.Spec)
	}

	if err := docker.CheckDockerImageVersion(docker.ImageElasticsearch, string(elastic.Spec.Version)); err != nil {
		return fmt.Errorf(`Image %v:%v not found`, docker.ImageElasticsearch, elastic.Spec.Version)
	}

	if elastic.Spec.Storage != nil {
		var err error
		if err = amv.ValidateStorage(client, elastic.Spec.Storage); err != nil {
			return err
		}
	}

	backupScheduleSpec := elastic.Spec.BackupSchedule
	if backupScheduleSpec != nil {
		if err := amv.ValidateBackupSchedule(client, backupScheduleSpec, elastic.Namespace); err != nil {
			return err
		}
	}

	monitorSpec := elastic.Spec.Monitor
	if monitorSpec != nil {
		if err := amv.ValidateMonitorSpec(monitorSpec); err != nil {
			return err
		}

	}
	return nil
}
