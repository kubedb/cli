package validator

import (
	"fmt"

	"github.com/ghodss/yaml"
	tapi "github.com/kubedb/apimachinery/apis/kubedb/v1alpha1"
	cs "github.com/kubedb/apimachinery/client/clientset/versioned/typed/kubedb/v1alpha1"
	amv "github.com/kubedb/apimachinery/pkg/validator"
	"github.com/kubedb/cli/pkg/encoder"
	esv "github.com/kubedb/elasticsearch/pkg/validator"
	memv "github.com/kubedb/memcached/pkg/admission"
	mgv "github.com/kubedb/mongodb/pkg/admission"
	msv "github.com/kubedb/mysql/pkg/admission"
	pgv "github.com/kubedb/postgres/pkg/admission"
	rdv "github.com/kubedb/redis/pkg/admission"
	"k8s.io/client-go/kubernetes"
	"k8s.io/kubernetes/pkg/kubectl/resource"
)

func Validate(client kubernetes.Interface, extClient cs.KubedbV1alpha1Interface, info *resource.Info) error {
	objByte, err := encoder.Encode(info.Object)
	if err != nil {
		return err
	}

	kind := info.Object.GetObjectKind().GroupVersionKind().Kind
	switch kind {
	case tapi.ResourceKindElasticsearch:
		var elasticsearch *tapi.Elasticsearch
		if err := yaml.Unmarshal(objByte, &elasticsearch); err != nil {
			return err
		}
		return esv.ValidateElasticsearch(client, extClient, elasticsearch)
	case tapi.ResourceKindPostgres:
		var postgres *tapi.Postgres
		if err := yaml.Unmarshal(objByte, &postgres); err != nil {
			return err
		}
		return pgv.ValidatePostgres(client, extClient, postgres)
	case tapi.ResourceKindMySQL:
		var mysql *tapi.MySQL
		if err := yaml.Unmarshal(objByte, &mysql); err != nil {
			return err
		}
		return msv.ValidateMySQL(client, extClient, mysql)
	case tapi.ResourceKindMongoDB:
		var mongodb *tapi.MongoDB
		if err := yaml.Unmarshal(objByte, &mongodb); err != nil {
			return err
		}
		return mgv.ValidateMongoDB(client, extClient, mongodb)
	case tapi.ResourceKindRedis:
		var redis *tapi.Redis
		if err := yaml.Unmarshal(objByte, &redis); err != nil {
			return err
		}
		return rdv.ValidateRedis(client, extClient, redis)
	case tapi.ResourceKindMemcached:
		var memcached *tapi.Memcached
		if err := yaml.Unmarshal(objByte, &memcached); err != nil {
			return err
		}
		return memv.ValidateMemcached(client, extClient, memcached)
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
		var elasticsearch *tapi.Elasticsearch
		if err := yaml.Unmarshal(objByte, &elasticsearch); err != nil {
			return err
		}
		if elasticsearch.Spec.DoNotPause {
			return fmt.Errorf(`elasticsearch "%v" can't be paused. To continue delete, unset spec.doNotPause and retry`, elasticsearch.Name)
		}
	case tapi.ResourceKindPostgres:
		var postgres *tapi.Postgres
		if err := yaml.Unmarshal(objByte, &postgres); err != nil {
			return err
		}
		if postgres.Spec.DoNotPause {
			return fmt.Errorf(`postgres "%v" can't be paused. To continue delete, unset spec.doNotPause and retry`, postgres.Name)
		}
	case tapi.ResourceKindMySQL:
		var mysql *tapi.MySQL
		if err := yaml.Unmarshal(objByte, &mysql); err != nil {
			return err
		}
		if mysql.Spec.DoNotPause {
			return fmt.Errorf(`mysql "%v" can't be paused. To continue delete, unset spec.doNotPause and retry`, mysql.Name)
		}
	case tapi.ResourceKindMongoDB:
		var mongodb *tapi.MongoDB
		if err := yaml.Unmarshal(objByte, &mongodb); err != nil {
			return err
		}
		if mongodb.Spec.DoNotPause {
			return fmt.Errorf(`mongodb "%v" can't be paused. To continue delete, unset spec.doNotPause and retry`, mongodb.Name)
		}
	case tapi.ResourceKindRedis:
		var redis *tapi.Redis
		if err := yaml.Unmarshal(objByte, &redis); err != nil {
			return err
		}
		if redis.Spec.DoNotPause {
			return fmt.Errorf(`redis "%v" can't be paused. To continue delete, unset spec.doNotPause and retry`, redis.Name)
		}
	case tapi.ResourceKindMemcached:
		var memcached *tapi.Memcached
		if err := yaml.Unmarshal(objByte, &memcached); err != nil {
			return err
		}
		if memcached.Spec.DoNotPause {
			return fmt.Errorf(`memcached "%v" can't be paused. To continue delete, unset spec.doNotPause and retry`, memcached.Name)
		}
	}
	return nil
}
