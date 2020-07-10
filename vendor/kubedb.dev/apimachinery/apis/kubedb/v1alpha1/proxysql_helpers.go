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

package v1alpha1

import (
	"fmt"

	"kubedb.dev/apimachinery/apis"
	"kubedb.dev/apimachinery/apis/kubedb"
	"kubedb.dev/apimachinery/crds"

	"github.com/appscode/go/types"
	apps "k8s.io/api/apps/v1"
	"kmodules.xyz/client-go/apiextensions"
	meta_util "kmodules.xyz/client-go/meta"
	appcat "kmodules.xyz/custom-resources/apis/appcatalog/v1alpha1"
	mona "kmodules.xyz/monitoring-agent-api/api/v1"
)

func (_ ProxySQL) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralProxySQL))
}

var _ apis.ResourceInfo = &ProxySQL{}

func (p ProxySQL) OffshootName() string {
	return p.Name
}

func (p ProxySQL) OffshootSelectors() map[string]string {
	return map[string]string{
		LabelProxySQLName:        p.Name,
		LabelProxySQLLoadBalance: string(*p.Spec.Mode),
		LabelDatabaseKind:        ResourceKindProxySQL,
	}
}

func (p ProxySQL) OffshootLabels() map[string]string {
	out := p.OffshootSelectors()
	out[meta_util.NameLabelKey] = ResourceSingularProxySQL
	out[meta_util.VersionLabelKey] = string(p.Spec.Version)
	out[meta_util.InstanceLabelKey] = p.Name
	out[meta_util.ComponentLabelKey] = ComponentDatabase
	out[meta_util.ManagedByLabelKey] = GenericKey
	return meta_util.FilterKeys(GenericKey, out, p.Labels)
}

func (p ProxySQL) ResourceShortCode() string {
	return ""
}

func (p ProxySQL) ResourceKind() string {
	return ResourceKindProxySQL
}

func (p ProxySQL) ResourceSingular() string {
	return ResourceSingularProxySQL
}

func (p ProxySQL) ResourcePlural() string {
	return ResourcePluralProxySQL
}

func (p ProxySQL) ServiceName() string {
	return p.OffshootName()
}

type proxysqlApp struct {
	*ProxySQL
}

func (p proxysqlApp) Name() string {
	return p.ProxySQL.Name
}

func (p proxysqlApp) Type() appcat.AppType {
	return appcat.AppType(fmt.Sprintf("%s/%s", kubedb.GroupName, ResourceSingularProxySQL))
}

func (p ProxySQL) AppBindingMeta() appcat.AppBindingMeta {
	return &proxysqlApp{&p}
}

type proxysqlStatsService struct {
	*ProxySQL
}

func (p proxysqlStatsService) GetNamespace() string {
	return p.ProxySQL.GetNamespace()
}

func (p proxysqlStatsService) ServiceName() string {
	return p.OffshootName() + "-stats"
}

func (p proxysqlStatsService) ServiceMonitorName() string {
	return fmt.Sprintf("kubedb-%s-%s", p.Namespace, p.Name)
}

func (p proxysqlStatsService) Path() string {
	return DefaultStatsPath
}

func (p proxysqlStatsService) Scheme() string {
	return ""
}

func (p ProxySQL) StatsService() mona.StatsAccessor {
	return &proxysqlStatsService{&p}
}

func (p ProxySQL) StatsServiceLabels() map[string]string {
	lbl := meta_util.FilterKeys(GenericKey, p.OffshootSelectors(), p.Labels)
	lbl[LabelRole] = RoleStats
	return lbl
}

func (p *ProxySQL) GetMonitoringVendor() string {
	if p.Spec.Monitor != nil {
		return p.Spec.Monitor.Agent.Vendor()
	}
	return ""
}

func (p *ProxySQL) SetDefaults() {
	if p == nil {
		return
	}

	if p == nil || p.Spec.Mode == nil || p.Spec.Backend == nil {
		return
	}

	if p.Spec.Replicas == nil {
		p.Spec.Replicas = types.Int32P(1)
	}

	if p.Spec.UpdateStrategy.Type == "" {
		p.Spec.UpdateStrategy.Type = apps.RollingUpdateStatefulSetStrategyType
	}

	p.Spec.Monitor.SetDefaults()
}

func (p *ProxySQLSpec) GetSecrets() []string {
	if p == nil {
		return nil
	}

	var secrets []string
	if p.ProxySQLSecret != nil {
		secrets = append(secrets, p.ProxySQLSecret.SecretName)
	}
	return secrets
}
