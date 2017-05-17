package decoder

import (
	"fmt"

	"github.com/ghodss/yaml"
	tapi "github.com/k8sdb/apimachinery/api"
	"k8s.io/kubernetes/pkg/runtime"
)

func Decode(kind string, data []byte) (runtime.Object, error) {
	switch kind {
	case tapi.ResourceKindElastic:
		var elastic *tapi.Elastic
		if err := yaml.Unmarshal(data, &elastic); err != nil {
			return nil, err
		}
		return elastic, nil
	case tapi.ResourceKindPostgres:
		var postgres *tapi.Postgres
		if err := yaml.Unmarshal(data, &postgres); err != nil {
			return nil, err
		}
		return postgres, nil

	case tapi.ResourceKindSnapshot:
		var snapshot *tapi.Snapshot
		if err := yaml.Unmarshal(data, &snapshot); err != nil {
			return nil, err
		}
		return snapshot, nil
	case tapi.ResourceKindDeletedDatabase:
		var deletedDb *tapi.DeletedDatabase
		if err := yaml.Unmarshal(data, &deletedDb); err != nil {
			return nil, err
		}
		return deletedDb, nil
	}

	return nil, fmt.Errorf(`Invalid kind: "%v"`, kind)
}
