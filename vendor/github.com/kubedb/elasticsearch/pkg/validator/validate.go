package validator

import (
	"errors"
	"fmt"

	api "github.com/kubedb/apimachinery/apis/kubedb/v1alpha1"
	"github.com/kubedb/apimachinery/pkg/docker"
	amv "github.com/kubedb/apimachinery/pkg/validator"
	"k8s.io/client-go/kubernetes"
)

func ValidateElasticsearch(client kubernetes.Interface, elasticsearch *api.Elasticsearch) error {
	if elasticsearch.Spec.Version == "" {
		return fmt.Errorf(`object 'Version' is missing in '%v'`, elasticsearch.Spec)
	}

	if err := docker.CheckDockerImageVersion(docker.ImageElasticsearch, string(elasticsearch.Spec.Version)); err != nil {
		return fmt.Errorf(`image %v:%v not found`, docker.ImageElasticsearch, elasticsearch.Spec.Version)
	}

	topology := elasticsearch.Spec.Topology
	if topology != nil {
		if topology.Client.Prefix == topology.Master.Prefix {
			return errors.New("client & master node should not have same prefix")
		}
		if topology.Client.Prefix == topology.Data.Prefix {
			return errors.New("client & data node should not have same prefix")
		}
		if topology.Master.Prefix == topology.Data.Prefix {
			return errors.New("master & data node should not have same prefix")
		}
	}

	if elasticsearch.Spec.Storage != nil {
		if err := amv.ValidateStorage(client, elasticsearch.Spec.Storage); err != nil {
			return err
		}
	}

	backupScheduleSpec := elasticsearch.Spec.BackupSchedule
	if backupScheduleSpec != nil {
		if err := amv.ValidateBackupSchedule(client, backupScheduleSpec, elasticsearch.Namespace); err != nil {
			return err
		}
	}

	monitorSpec := elasticsearch.Spec.Monitor
	if monitorSpec != nil {
		if err := amv.ValidateMonitorSpec(monitorSpec); err != nil {
			return err
		}

	}
	return nil
}
