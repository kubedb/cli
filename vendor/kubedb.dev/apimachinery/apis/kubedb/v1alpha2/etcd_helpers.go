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

	"gomodules.xyz/pointer"
	"k8s.io/apimachinery/pkg/labels"
	appslister "k8s.io/client-go/listers/apps/v1"
	"kmodules.xyz/client-go/apiextensions"
	meta_util "kmodules.xyz/client-go/meta"
	appcat "kmodules.xyz/custom-resources/apis/appcatalog/v1alpha1"
	mona "kmodules.xyz/monitoring-agent-api/api/v1"
)

func (_ Etcd) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralEtcd))
}

var _ apis.ResourceInfo = &Etcd{}

func (e Etcd) OffshootName() string {
	return e.Name
}

func (e Etcd) OffshootSelectors() map[string]string {
	return map[string]string{
		meta_util.NameLabelKey:      e.ResourceFQN(),
		meta_util.InstanceLabelKey:  e.Name,
		meta_util.ManagedByLabelKey: kubedb.GroupName,
	}
}

func (e Etcd) OffshootLabels() map[string]string {
	return e.offshootLabels(e.OffshootSelectors(), nil)
}

func (e Etcd) PodLabels() map[string]string {
	return e.offshootLabels(e.OffshootSelectors(), e.Spec.PodTemplate.Labels)
}

func (e Etcd) PodControllerLabels() map[string]string {
	return e.offshootLabels(e.OffshootSelectors(), e.Spec.PodTemplate.Controller.Labels)
}

func (e Etcd) ServiceLabels(alias ServiceAlias, extraLabels ...map[string]string) map[string]string {
	svcTemplate := GetServiceTemplate(e.Spec.ServiceTemplates, alias)
	return e.offshootLabels(meta_util.OverwriteKeys(e.OffshootSelectors(), extraLabels...), svcTemplate.Labels)
}

func (e Etcd) offshootLabels(selector, override map[string]string) map[string]string {
	selector[meta_util.ComponentLabelKey] = ComponentDatabase
	return meta_util.FilterKeys(kubedb.GroupName, selector, meta_util.OverwriteKeys(nil, e.Labels, override))
}

func (e Etcd) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", ResourcePluralEtcd, kubedb.GroupName)
}

func (e Etcd) ResourceShortCode() string {
	return ResourceCodeEtcd
}

func (e Etcd) ResourceKind() string {
	return ResourceKindEtcd
}

func (e Etcd) ResourceSingular() string {
	return ResourceSingularEtcd
}

func (e Etcd) ResourcePlural() string {
	return ResourcePluralEtcd
}

func (e Etcd) ClientServiceName() string {
	return e.OffshootName() + "-client"
}

func (e Etcd) PeerServiceName() string {
	return e.OffshootName()
}

type etcdApp struct {
	*Etcd
}

func (r etcdApp) Name() string {
	return r.Etcd.Name
}

func (r etcdApp) Type() appcat.AppType {
	return appcat.AppType(fmt.Sprintf("%s/%s", kubedb.GroupName, ResourceSingularEtcd))
}

func (r Etcd) AppBindingMeta() appcat.AppBindingMeta {
	return &etcdApp{&r}
}

type etcdStatsService struct {
	*Etcd
}

func (e etcdStatsService) GetNamespace() string {
	return e.Etcd.GetNamespace()
}

func (e etcdStatsService) ServiceName() string {
	return e.OffshootName() + "-stats"
}

func (e etcdStatsService) ServiceMonitorName() string {
	return e.ServiceName()
}

func (e etcdStatsService) ServiceMonitorAdditionalLabels() map[string]string {
	return e.OffshootLabels()
}

func (e etcdStatsService) Path() string {
	return "/metrics"
}

func (e etcdStatsService) Scheme() string {
	return ""
}

func (e Etcd) StatsService() mona.StatsAccessor {
	return &etcdStatsService{&e}
}

func (e Etcd) StatsServiceLabels() map[string]string {
	return e.ServiceLabels(StatsServiceAlias, map[string]string{LabelRole: RoleStats})
}

func (e *Etcd) SetDefaults() {
	if e == nil {
		return
	}

	// perform defaulting
	if e.Spec.Replicas == nil {
		e.Spec.Replicas = pointer.Int32P(1)
	}
	if e.Spec.StorageType == "" {
		e.Spec.StorageType = StorageTypeDurable
	}
	if e.Spec.TerminationPolicy == "" {
		e.Spec.TerminationPolicy = TerminationPolicyDelete
	}

	e.Spec.Monitor.SetDefaults()
	SetDefaultResourceLimits(&e.Spec.PodTemplate.Spec.Resources, DefaultResources)
}

func (e *EtcdSpec) GetPersistentSecrets() []string {
	return nil
}

func (e *Etcd) ReplicasAreReady(lister appslister.StatefulSetLister) (bool, string, error) {
	// Desire number of statefulSets
	expectedItems := 1
	return checkReplicas(lister.StatefulSets(e.Namespace), labels.SelectorFromSet(e.OffshootLabels()), expectedItems)
}
