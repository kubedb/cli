package v1alpha1

import (
	"fmt"

	"kubedb.dev/apimachinery/apis"
	"kubedb.dev/apimachinery/apis/kubedb"

	"github.com/appscode/go/types"
	apps "k8s.io/api/apps/v1"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	crdutils "kmodules.xyz/client-go/apiextensions/v1beta1"
	meta_util "kmodules.xyz/client-go/meta"
	appcat "kmodules.xyz/custom-resources/apis/appcatalog/v1alpha1"
	mona "kmodules.xyz/monitoring-agent-api/api/v1"
)

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

func (p ProxySQL) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crdutils.NewCustomResourceDefinition(crdutils.Config{
		Group:         SchemeGroupVersion.Group,
		Plural:        ResourcePluralProxySQL,
		Singular:      ResourceSingularProxySQL,
		Kind:          ResourceKindProxySQL,
		Categories:    []string{"datastore", "kubedb", "appscode", "all"},
		ResourceScope: string(apiextensions.NamespaceScoped),
		Versions: []apiextensions.CustomResourceDefinitionVersion{
			{
				Name:    SchemeGroupVersion.Version,
				Served:  true,
				Storage: true,
			},
		},
		Labels: crdutils.Labels{
			LabelsMap: map[string]string{"app": "kubedb"},
		},
		SpecDefinitionName:      "kubedb.dev/apimachinery/apis/kubedb/v1alpha1.ProxySQL",
		EnableValidation:        true,
		GetOpenAPIDefinitions:   GetOpenAPIDefinitions,
		EnableStatusSubresource: true,
		AdditionalPrinterColumns: []apiextensions.CustomResourceColumnDefinition{
			{
				Name:     "Version",
				Type:     "string",
				JSONPath: ".spec.version",
			},
			{
				Name:     "Status",
				Type:     "string",
				JSONPath: ".status.phase",
			},
			{
				Name:     "Age",
				Type:     "date",
				JSONPath: ".metadata.creationTimestamp",
			},
		},
	}, apis.SetNameSchema)
}

func (p *ProxySQL) SetDefaults() {
	if p == nil {
		return
	}
	p.Spec.SetDefaults()
}

func (p *ProxySQLSpec) SetDefaults() {
	if p == nil || p.Mode == nil || p.Backend == nil {
		return
	}

	if p.Replicas == nil {
		p.Replicas = types.Int32P(1)
	}

	if p.UpdateStrategy.Type == "" {
		p.UpdateStrategy.Type = apps.RollingUpdateStatefulSetStrategyType
	}
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
