package validator

import (
	"fmt"

	"github.com/ghodss/yaml"
	tapi "github.com/k8sdb/apimachinery/apis/kubedb/v1alpha1"
	amv "github.com/k8sdb/apimachinery/pkg/validator"
	"github.com/k8sdb/cli/pkg/encoder"
	esv "github.com/k8sdb/elasticsearch/pkg/validator"
	mgv "github.com/k8sdb/mongodb/pkg/validator"
	pgv "github.com/k8sdb/postgres/pkg/validator"
	"k8s.io/client-go/kubernetes"
	"k8s.io/kubernetes/pkg/kubectl/resource"
)

func Validate(client kubernetes.Interface, info *resource.Info) error {
	objByte, err := encoder.Encode(info.Object)
	if err != nil {
		return err
	}

	kind := info.Object.GetObjectKind().GroupVersionKind().Kind
	switch kind {
	case tapi.ResourceKindElasticsearch:
		var elastic *tapi.Elasticsearch
		if err := yaml.Unmarshal(objByte, &elastic); err != nil {
			return err
		}
		return esv.ValidateElastic(client, elastic)
	case tapi.ResourceKindPostgres:
		var postgres *tapi.Postgres
		if err := yaml.Unmarshal(objByte, &postgres); err != nil {
			return err
		}
		return pgv.ValidatePostgres(client, postgres)
	case tapi.ResourceKindMongoDB:
		var mongodb *tapi.MongoDB
		if err := yaml.Unmarshal(objByte, &mongodb); err != nil {
			return err
		}
		return mgv.ValidateMongoDB(client, mongodb)
	case tapi.ResourceKindSnapshot:
		var snapshot *tapi.Snapshot
		if err := yaml.Unmarshal(objByte, &snapshot); err != nil {
			return err
		}
		return amv.ValidateSnapshotSpec(client, snapshot.Spec.SnapshotStorageSpec, info.Namespace)
	}
	return nil
}

func ValidateDeletion(info *resource.Info) error {
	objByte, err := encoder.Encode(info.Object)
	if err != nil {
		return err
	}

	kind := info.Object.GetObjectKind().GroupVersionKind().Kind
	switch kind {
	case tapi.ResourceKindElasticsearch:
		var elastic *tapi.Elasticsearch
		if err := yaml.Unmarshal(objByte, &elastic); err != nil {
			return err
		}
		if elastic.Spec.DoNotPause {
			return fmt.Errorf(`Elasticsearch "%v" can't be paused. To continue delete, unset spec.doNotPause and retry.`, elastic.Name)
		}
	case tapi.ResourceKindPostgres:
		var postgres *tapi.Postgres
		if err := yaml.Unmarshal(objByte, &postgres); err != nil {
			return err
		}
		if postgres.Spec.DoNotPause {
			return fmt.Errorf(`Postgres "%v" can't be paused. To continue delete, unset spec.doNotPause and retry.`, postgres.Name)
		}
	case tapi.ResourceKindMongoDB:
		var mongodb *tapi.MongoDB
		if err := yaml.Unmarshal(objByte, &mongodb); err != nil {
			return err
		}
		if mongodb.Spec.DoNotPause {
			return fmt.Errorf(`MongoDB "%v" can't be paused. To continue delete, unset spec.doNotPause and retry.`, mongodb.Name)
		}

	}
	return nil
}
