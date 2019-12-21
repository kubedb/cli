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
	"fmt"

	"kubedb.dev/apimachinery/api/crds"
	"kubedb.dev/apimachinery/apis"
	"kubedb.dev/apimachinery/apis/catalog/v1alpha1"
	"kubedb.dev/apimachinery/apis/kubedb"

	apps "k8s.io/api/apps/v1"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	meta_util "kmodules.xyz/client-go/meta"
	appcat "kmodules.xyz/custom-resources/apis/appcatalog/v1alpha1"
	mona "kmodules.xyz/monitoring-agent-api/api/v1"
)

func (_ Elasticsearch) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralElasticsearch))
}

var _ apis.ResourceInfo = &Elasticsearch{}

func (e Elasticsearch) OffshootName() string {
	return e.Name
}

func (e Elasticsearch) OffshootSelectors() map[string]string {
	return map[string]string{
		LabelDatabaseKind: ResourceKindElasticsearch,
		LabelDatabaseName: e.Name,
	}
}

func (e Elasticsearch) OffshootLabels() map[string]string {
	out := e.OffshootSelectors()
	out[meta_util.NameLabelKey] = ResourceSingularElasticsearch
	out[meta_util.VersionLabelKey] = string(e.Spec.Version)
	out[meta_util.InstanceLabelKey] = e.Name
	out[meta_util.ComponentLabelKey] = ComponentDatabase
	out[meta_util.ManagedByLabelKey] = GenericKey
	return meta_util.FilterKeys(GenericKey, out, e.Labels)
}

func (e Elasticsearch) ResourceShortCode() string {
	return ResourceCodeElasticsearch
}

func (e Elasticsearch) ResourceKind() string {
	return ResourceKindElasticsearch
}

func (e Elasticsearch) ResourceSingular() string {
	return ResourceSingularElasticsearch
}

func (e Elasticsearch) ResourcePlural() string {
	return ResourcePluralElasticsearch
}

func (e Elasticsearch) ServiceName() string {
	return e.OffshootName()
}

func (e *Elasticsearch) MasterServiceName() string {
	return fmt.Sprintf("%v-master", e.ServiceName())
}

// Snapshot service account name.
func (e Elasticsearch) SnapshotSAName() string {
	return fmt.Sprintf("%v-snapshot", e.OffshootName())
}

func (e *Elasticsearch) GetConnectionScheme() string {
	scheme := "http"
	if e.Spec.EnableSSL {
		scheme = "https"
	}
	return scheme
}

func (e *Elasticsearch) GetConnectionURL() string {
	return fmt.Sprintf("%v://%s.%s:%d", e.GetConnectionScheme(), e.OffshootName(), e.Namespace, ElasticsearchRestPort)
}

type elasticsearchApp struct {
	*Elasticsearch
}

func (r elasticsearchApp) Name() string {
	return r.Elasticsearch.Name
}

func (r elasticsearchApp) Type() appcat.AppType {
	return appcat.AppType(fmt.Sprintf("%s/%s", kubedb.GroupName, ResourceSingularElasticsearch))
}

func (r Elasticsearch) AppBindingMeta() appcat.AppBindingMeta {
	return &elasticsearchApp{&r}
}

type elasticsearchStatsService struct {
	*Elasticsearch
}

func (e elasticsearchStatsService) GetNamespace() string {
	return e.Elasticsearch.GetNamespace()
}

func (e elasticsearchStatsService) ServiceName() string {
	return e.OffshootName() + "-stats"
}

func (e elasticsearchStatsService) ServiceMonitorName() string {
	return fmt.Sprintf("kubedb-%s-%s", e.Namespace, e.Name)
}

func (e elasticsearchStatsService) Path() string {
	return DefaultStatsPath
}

func (e elasticsearchStatsService) Scheme() string {
	return ""
}

func (e Elasticsearch) StatsService() mona.StatsAccessor {
	return &elasticsearchStatsService{&e}
}

func (e Elasticsearch) StatsServiceLabels() map[string]string {
	lbl := meta_util.FilterKeys(GenericKey, e.OffshootSelectors(), e.Labels)
	lbl[LabelRole] = RoleStats
	return lbl
}

func (e *Elasticsearch) GetMonitoringVendor() string {
	if e.Spec.Monitor != nil {
		return e.Spec.Monitor.Agent.Vendor()
	}
	return ""
}

func (e *Elasticsearch) SetDefaults() {
	if e == nil {
		return
	}
	e.Spec.SetDefaults()

	if e.Spec.PodTemplate.Spec.ServiceAccountName == "" {
		e.Spec.PodTemplate.Spec.ServiceAccountName = e.OffshootName()
	}
}

func (e *ElasticsearchSpec) SetDefaults() {
	if e == nil {
		return
	}

	// perform defaulting
	e.BackupSchedule.SetDefaults()

	if !e.DisableSecurity && e.AuthPlugin == v1alpha1.ElasticsearchAuthPluginNone {
		e.DisableSecurity = true
	}
	e.AuthPlugin = ""
	if e.StorageType == "" {
		e.StorageType = StorageTypeDurable
	}
	if e.UpdateStrategy.Type == "" {
		e.UpdateStrategy.Type = apps.RollingUpdateStatefulSetStrategyType
	}
	if e.TerminationPolicy == "" {
		e.TerminationPolicy = TerminationPolicyDelete
	}
}

func (e *ElasticsearchSpec) GetSecrets() []string {
	if e == nil {
		return nil
	}

	var secrets []string
	if e.DatabaseSecret != nil {
		secrets = append(secrets, e.DatabaseSecret.SecretName)
	}
	if e.CertificateSecret != nil {
		secrets = append(secrets, e.CertificateSecret.SecretName)
	}
	return secrets
}
