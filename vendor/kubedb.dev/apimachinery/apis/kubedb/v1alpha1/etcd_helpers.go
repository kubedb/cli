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

	"github.com/appscode/go/types"
	apps "k8s.io/api/apps/v1"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
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
		LabelDatabaseName: e.Name,
		LabelDatabaseKind: ResourceKindEtcd,
	}
}

func (e Etcd) OffshootLabels() map[string]string {
	out := e.OffshootSelectors()
	out[meta_util.NameLabelKey] = ResourceSingularEtcd
	out[meta_util.VersionLabelKey] = string(e.Spec.Version)
	out[meta_util.InstanceLabelKey] = e.Name
	out[meta_util.ComponentLabelKey] = ComponentDatabase
	out[meta_util.ManagedByLabelKey] = GenericKey
	return meta_util.FilterKeys(GenericKey, out, e.Labels)
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
	return fmt.Sprintf("kubedb-%s-%s", e.Namespace, e.Name)
}

func (e etcdStatsService) Path() string {
	return fmt.Sprintf("/metrics")
}

func (e etcdStatsService) Scheme() string {
	return ""
}

func (e Etcd) StatsService() mona.StatsAccessor {
	return &etcdStatsService{&e}
}

func (e Etcd) StatsServiceLabels() map[string]string {
	lbl := meta_util.FilterKeys(GenericKey, e.OffshootSelectors(), e.Labels)
	lbl[LabelRole] = RoleStats
	return lbl
}

func (e *Etcd) GetMonitoringVendor() string {
	if e.Spec.Monitor != nil {
		return e.Spec.Monitor.Agent.Vendor()
	}
	return ""
}

func (e *Etcd) SetDefaults() {
	if e == nil {
		return
	}

	// perform defaulting
	if e.Spec.Replicas == nil {
		e.Spec.Replicas = types.Int32P(1)
	}
	if e.Spec.StorageType == "" {
		e.Spec.StorageType = StorageTypeDurable
	}
	if e.Spec.UpdateStrategy.Type == "" {
		e.Spec.UpdateStrategy.Type = apps.RollingUpdateStatefulSetStrategyType
	}
	if e.Spec.TerminationPolicy == "" {
		e.Spec.TerminationPolicy = TerminationPolicyDelete
	} else if e.Spec.TerminationPolicy == TerminationPolicyPause {
		e.Spec.TerminationPolicy = TerminationPolicyHalt
	}

	e.Spec.Monitor.SetDefaults()
}

func (e *EtcdSpec) GetSecrets() []string {
	return nil
}
