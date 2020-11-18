/*
Copyright AppsCode Inc. and Contributors

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

package v1alpha2

import (
	"fmt"

	"kubedb.dev/apimachinery/apis"
	"kubedb.dev/apimachinery/apis/kubedb"
	"kubedb.dev/apimachinery/crds"

	"k8s.io/apimachinery/pkg/labels"
	appslister "k8s.io/client-go/listers/apps/v1"
	"kmodules.xyz/client-go/apiextensions"
	meta_util "kmodules.xyz/client-go/meta"
	appcat "kmodules.xyz/custom-resources/apis/appcatalog/v1alpha1"
	mona "kmodules.xyz/monitoring-agent-api/api/v1"
)

func (_ Memcached) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralMemcached))
}

var _ apis.ResourceInfo = &Memcached{}

func (m Memcached) OffshootName() string {
	return m.Name
}

func (m Memcached) OffshootSelectors() map[string]string {
	return map[string]string{
		LabelDatabaseKind: ResourceKindMemcached,
		LabelDatabaseName: m.Name,
	}
}

func (m Memcached) OffshootLabels() map[string]string {
	out := m.OffshootSelectors()
	out[meta_util.NameLabelKey] = ResourceSingularMemcached
	out[meta_util.InstanceLabelKey] = m.Name
	out[meta_util.ComponentLabelKey] = ComponentDatabase
	out[meta_util.ManagedByLabelKey] = kubedb.GroupName
	return meta_util.FilterKeys(kubedb.GroupName, out, m.Labels)
}

func (m Memcached) ResourceShortCode() string {
	return ResourceCodeMemcached
}

func (m Memcached) ResourceKind() string {
	return ResourceKindMemcached
}

func (m Memcached) ResourceSingular() string {
	return ResourceSingularMemcached
}

func (m Memcached) ResourcePlural() string {
	return ResourcePluralMemcached
}

func (m Memcached) ServiceName() string {
	return m.OffshootName()
}

func (m Memcached) GoverningServiceName() string {
	return meta_util.NameWithSuffix(m.ServiceName(), "pods")
}

type memcachedApp struct {
	*Memcached
}

func (r memcachedApp) Name() string {
	return r.Memcached.Name
}

func (r memcachedApp) Type() appcat.AppType {
	return appcat.AppType(fmt.Sprintf("%s/%s", kubedb.GroupName, ResourceSingularMemcached))
}

func (m Memcached) AppBindingMeta() appcat.AppBindingMeta {
	return &memcachedApp{&m}
}

type memcachedStatsService struct {
	*Memcached
}

func (m memcachedStatsService) GetNamespace() string {
	return m.Memcached.GetNamespace()
}

func (m memcachedStatsService) ServiceName() string {
	return m.OffshootName() + "-stats"
}

func (m memcachedStatsService) ServiceMonitorName() string {
	return m.ServiceName()
}

func (m memcachedStatsService) ServiceMonitorAdditionalLabels() map[string]string {
	return m.OffshootLabels()
}

func (m memcachedStatsService) Path() string {
	return DefaultStatsPath
}

func (m memcachedStatsService) Scheme() string {
	return ""
}

func (m Memcached) StatsService() mona.StatsAccessor {
	return &memcachedStatsService{&m}
}

func (m Memcached) StatsServiceLabels() map[string]string {
	lbl := meta_util.FilterKeys(kubedb.GroupName, m.OffshootSelectors(), m.Labels)
	lbl[LabelRole] = RoleStats
	return lbl
}

func (m *Memcached) SetDefaults() {
	if m == nil {
		return
	}

	// perform defaulting
	if m.Spec.TerminationPolicy == "" {
		m.Spec.TerminationPolicy = TerminationPolicyDelete
	}

	if m.Spec.PodTemplate.Spec.ServiceAccountName == "" {
		m.Spec.PodTemplate.Spec.ServiceAccountName = m.OffshootName()
	}

	m.Spec.Monitor.SetDefaults()
	setDefaultResourceLimits(&m.Spec.PodTemplate.Spec.Resources, defaultResourceLimits, defaultResourceRequests)
}

func (m *MemcachedSpec) GetPersistentSecrets() []string {
	return nil
}

func (m *Memcached) ReplicasAreReady(lister appslister.StatefulSetLister) (bool, string, error) {
	// Desire number of statefulSets
	expectedItems := 1
	return checkReplicas(lister.StatefulSets(m.Namespace), labels.SelectorFromSet(m.OffshootLabels()), expectedItems)
}
