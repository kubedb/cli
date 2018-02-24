package validator

import (
	"fmt"

	"github.com/appscode/go/types"
	api "github.com/kubedb/apimachinery/apis/kubedb/v1alpha1"
	amv "github.com/kubedb/apimachinery/pkg/validator"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/kubernetes"
)

var (
	elasticVersions = sets.NewString("5.6", "5.6.4")
)

func ValidateElasticsearch(client kubernetes.Interface, elasticsearch *api.Elasticsearch) error {
	if elasticsearch.Spec.Version == "" {
		return fmt.Errorf(`object 'Version' is missing in '%v'`, elasticsearch.Spec)
	}

	// check Elasticsearch version validation
	if !elasticVersions.Has(string(elasticsearch.Spec.Version)) {
		return fmt.Errorf(`KubeDB doesn't support Elasticsearch version: %s`, string(elasticsearch.Spec.Version))
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

		if topology.Client.Replicas != nil {
			replicas := topology.Client.Replicas
			if types.Int32(replicas) < 1 {
				return fmt.Errorf(`topology.client.replicas "%d" invalid. Must be greater than zero`, replicas)
			}
		}

		if topology.Master.Replicas != nil {
			replicas := topology.Master.Replicas
			if types.Int32(replicas) < 1 {
				return fmt.Errorf(`topology.master.replicas "%d" invalid. Must be greater than zero`, replicas)
			}
		}

		if topology.Data.Replicas != nil {
			replicas := topology.Data.Replicas
			if types.Int32(replicas) < 1 {
				return fmt.Errorf(`topology.data.replicas "%d" invalid. Must be greater than zero`, replicas)
			}
		}
	} else {
		if elasticsearch.Spec.Replicas != nil {
			replicas := types.Int32(elasticsearch.Spec.Replicas)
			if replicas < 1 {
				return fmt.Errorf(`spec.replicas "%d" invalid. Must be greater than zero`, replicas)
			}
		}
	}

	if elasticsearch.Spec.Storage != nil {
		if err := amv.ValidateStorage(client, elasticsearch.Spec.Storage); err != nil {
			return err
		}
	}

	databaseSecret := elasticsearch.Spec.DatabaseSecret
	if databaseSecret != nil {
		if _, err := client.CoreV1().Secrets(elasticsearch.Namespace).Get(databaseSecret.SecretName, metav1.GetOptions{}); err != nil {
			return err
		}
	}

	certificateSecret := elasticsearch.Spec.CertificateSecret
	if certificateSecret != nil {
		if _, err := client.CoreV1().Secrets(elasticsearch.Namespace).Get(certificateSecret.SecretName, metav1.GetOptions{}); err != nil {
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
