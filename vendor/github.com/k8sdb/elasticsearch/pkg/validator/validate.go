package validator

import (
	"fmt"

	tapi "github.com/k8sdb/apimachinery/api"
	"github.com/k8sdb/apimachinery/pkg/docker"
	amv "github.com/k8sdb/apimachinery/pkg/validator"
	clientset "k8s.io/client-go/kubernetes"
)

func ValidateElastic(client clientset.Interface, elastic *tapi.Elastic) error {
	if elastic.Spec.Version == "" {
		return fmt.Errorf(`Object 'Version' is missing in '%v'`, elastic.Spec)
	}

	if err := docker.CheckDockerImageVersion(docker.ImageElasticsearch, string(elastic.Spec.Version)); err != nil {
		return fmt.Errorf(`Image %v:%v not found`, docker.ImageElasticsearch, elastic.Spec.Version)
	}

	storage := elastic.Spec.Storage
	if storage != nil {
		var err error
		if _, err = amv.ValidateStorageSpec(client, storage); err != nil {
			return err
		}
	}

	backupScheduleSpec := elastic.Spec.BackupSchedule
	if backupScheduleSpec != nil {
		if err := amv.ValidateBackupSchedule(backupScheduleSpec); err != nil {
			return err
		}

		if err := amv.CheckBucketAccess(client, backupScheduleSpec.SnapshotStorageSpec, elastic.Namespace); err != nil {
			return err
		}
	}

	monitorSpec := elastic.Spec.Monitor
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
