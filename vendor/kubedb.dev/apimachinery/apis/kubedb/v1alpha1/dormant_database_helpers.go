/*
Copyright The KubeDB Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	"kubedb.dev/apimachinery/api/crds"
	"kubedb.dev/apimachinery/apis"
	"kubedb.dev/apimachinery/apis/catalog/v1alpha1"

	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	meta_util "kmodules.xyz/client-go/meta"
)

func (_ DormantDatabase) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralDormantDatabase))
}

var _ apis.ResourceInfo = &DormantDatabase{}

func (d DormantDatabase) OffshootSelectors() map[string]string {
	selector := map[string]string{
		LabelDatabaseName: d.Name,
	}
	switch {
	case d.Spec.Origin.Spec.Etcd != nil:
		selector[LabelDatabaseKind] = ResourceKindEtcd
	case d.Spec.Origin.Spec.Elasticsearch != nil:
		selector[LabelDatabaseKind] = ResourceKindElasticsearch
	case d.Spec.Origin.Spec.Memcached != nil:
		selector[LabelDatabaseKind] = ResourceKindMemcached
	case d.Spec.Origin.Spec.MongoDB != nil:
		selector[LabelDatabaseKind] = ResourceKindMongoDB
	case d.Spec.Origin.Spec.MySQL != nil:
		selector[LabelDatabaseKind] = ResourceKindMySQL
	case d.Spec.Origin.Spec.PerconaXtraDB != nil:
		selector[LabelDatabaseKind] = ResourceKindPerconaXtraDB
	case d.Spec.Origin.Spec.MariaDB != nil:
		selector[LabelDatabaseKind] = ResourceKindMariaDB
	case d.Spec.Origin.Spec.Postgres != nil:
		selector[LabelDatabaseKind] = ResourceKindPostgres
	case d.Spec.Origin.Spec.Redis != nil:
		selector[LabelDatabaseKind] = ResourceKindRedis
	}
	return selector
}

func (d DormantDatabase) OffshootLabels() map[string]string {
	return meta_util.FilterKeys(GenericKey, d.OffshootSelectors(), d.Spec.Origin.Labels)
}

func (d DormantDatabase) OffshootName() string {
	return d.Name
}

func (d DormantDatabase) ResourceShortCode() string {
	return ResourceCodeDormantDatabase
}

func (d DormantDatabase) ResourceKind() string {
	return ResourceKindDormantDatabase
}

func (d DormantDatabase) ResourceSingular() string {
	return ResourceSingularDormantDatabase
}

func (d DormantDatabase) ResourcePlural() string {
	return ResourcePluralDormantDatabase
}

func (d *DormantDatabase) SetDefaults() {
	if d == nil {
		return
	}
	d.Spec.Origin.Spec.Elasticsearch.SetDefaults()
	d.Spec.Origin.Spec.Postgres.SetDefaults()
	d.Spec.Origin.Spec.MySQL.SetDefaults()
	d.Spec.Origin.Spec.PerconaXtraDB.SetDefaults()
	d.Spec.Origin.Spec.MariaDB.SetDefaults()
	d.Spec.Origin.Spec.MongoDB.SetDefaults(&v1alpha1.MongoDBVersion{})
	d.Spec.Origin.Spec.Redis.SetDefaults()
	d.Spec.Origin.Spec.Memcached.SetDefaults()
	d.Spec.Origin.Spec.Etcd.SetDefaults()
}

func (d *DormantDatabase) GetDatabaseSecrets() []string {
	if d == nil {
		return nil
	}

	var secrets []string
	secrets = append(secrets, d.Spec.Origin.Spec.Elasticsearch.GetSecrets()...)
	secrets = append(secrets, d.Spec.Origin.Spec.Postgres.GetSecrets()...)
	secrets = append(secrets, d.Spec.Origin.Spec.MySQL.GetSecrets()...)
	secrets = append(secrets, d.Spec.Origin.Spec.PerconaXtraDB.GetSecrets()...)
	secrets = append(secrets, d.Spec.Origin.Spec.MariaDB.GetSecrets()...)
	secrets = append(secrets, d.Spec.Origin.Spec.MongoDB.GetSecrets()...)
	secrets = append(secrets, d.Spec.Origin.Spec.Redis.GetSecrets()...)
	secrets = append(secrets, d.Spec.Origin.Spec.Memcached.GetSecrets()...)
	secrets = append(secrets, d.Spec.Origin.Spec.Etcd.GetSecrets()...)
	return secrets
}
