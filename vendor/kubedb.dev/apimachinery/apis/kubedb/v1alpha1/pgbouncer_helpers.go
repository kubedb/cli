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
	"kubedb.dev/apimachinery/apis/kubedb"

	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	meta_util "kmodules.xyz/client-go/meta"
	appcat "kmodules.xyz/custom-resources/apis/appcatalog/v1alpha1"
	mona "kmodules.xyz/monitoring-agent-api/api/v1"
)

func (_ PgBouncer) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralPgBouncer))
}

var _ apis.ResourceInfo = &PgBouncer{}

func (p PgBouncer) OffshootName() string {
	return p.Name
}

func (p PgBouncer) OffshootSelectors() map[string]string {
	return map[string]string{
		LabelDatabaseName: p.Name,
		LabelDatabaseKind: ResourceKindPgBouncer,
	}
}

func (p PgBouncer) OffshootLabels() map[string]string {
	out := p.OffshootSelectors()
	out[meta_util.NameLabelKey] = ResourceSingularPgBouncer
	out[meta_util.InstanceLabelKey] = p.Name
	out[meta_util.ComponentLabelKey] = "connection-pooler"
	out[meta_util.VersionLabelKey] = string(p.Spec.Version)
	out[meta_util.ManagedByLabelKey] = GenericKey
	return meta_util.FilterKeys(GenericKey, out, p.Labels)
}

func (p PgBouncer) ResourceShortCode() string {
	return ResourceCodePgBouncer
}

func (p PgBouncer) ResourceKind() string {
	return ResourceKindPgBouncer
}

func (p PgBouncer) ResourceSingular() string {
	return ResourceSingularPgBouncer
}

func (p PgBouncer) ResourcePlural() string {
	return ResourcePluralPgBouncer
}

func (p PgBouncer) ServiceName() string {
	return p.OffshootName()
}

type pgbouncerApp struct {
	*PgBouncer
}

func (r pgbouncerApp) Name() string {
	return r.PgBouncer.Name
}

func (r pgbouncerApp) Type() appcat.AppType {
	return appcat.AppType(fmt.Sprintf("%s/%s", kubedb.GroupName, ResourceSingularPgBouncer))
}

func (p PgBouncer) AppBindingMeta() appcat.AppBindingMeta {
	return &pgbouncerApp{&p}
}

type pgbouncerStatsService struct {
	*PgBouncer
}

func (p pgbouncerStatsService) GetNamespace() string {
	return p.PgBouncer.GetNamespace()
}

func (p pgbouncerStatsService) ServiceName() string {
	return p.OffshootName() + "-stats"
}

func (p pgbouncerStatsService) ServiceMonitorName() string {
	return fmt.Sprintf("kubedb-%s-%s", p.Namespace, p.Name)
}

func (p pgbouncerStatsService) Path() string {
	return DefaultStatsPath
}

func (p pgbouncerStatsService) Scheme() string {
	return ""
}

func (p PgBouncer) StatsService() mona.StatsAccessor {
	return &pgbouncerStatsService{&p}
}

func (p PgBouncer) StatsServiceLabels() map[string]string {
	lbl := meta_util.FilterKeys(GenericKey, p.OffshootSelectors(), p.Labels)
	lbl[LabelRole] = RoleStats
	return lbl
}

func (p *PgBouncer) GetMonitoringVendor() string {
	if p.Spec.Monitor != nil {
		return p.Spec.Monitor.Agent.Vendor()
	}
	return ""
}

func (p PgBouncer) ReplicasServiceName() string {
	return fmt.Sprintf("%v-replicas", p.Name)
}

func (p *PgBouncer) SetDefaults() {
	if p == nil {
		return
	}
	p.Spec.Monitor.SetDefaults()
}
