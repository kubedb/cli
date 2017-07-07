package validator

import (
	"github.com/ghodss/yaml"
	tapi "github.com/k8sdb/apimachinery/api"
	amv "github.com/k8sdb/apimachinery/pkg/validator"
	"github.com/k8sdb/cli/pkg/encoder"
	esv "github.com/k8sdb/elasticsearch/pkg/validator"
	pgv "github.com/k8sdb/postgres/pkg/validator"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/kubernetes/pkg/kubectl/resource"
)

func Validate(client clientset.Interface, info *resource.Info) error {
	objByte, err := encoder.Encode(info.Object)
	if err != nil {
		return err
	}

	kind := info.Object.GetObjectKind().GroupVersionKind().Kind
	switch kind {
	case tapi.ResourceKindElastic:
		var elastic *tapi.Elastic
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
	case tapi.ResourceKindSnapshot:
		var snapshot *tapi.Snapshot
		if err := yaml.Unmarshal(objByte, &snapshot); err != nil {
			return err
		}
		return amv.ValidateSnapshotSpec(client, snapshot.Spec.SnapshotStorageSpec, info.Namespace)
	}
	return nil
}
