package v1alpha1

import (
	"fmt"

	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	crdutils "kmodules.xyz/client-go/apiextensions/v1beta1"
	meta_util "kmodules.xyz/client-go/meta"
	appcat "kmodules.xyz/custom-resources/apis/appcatalog/v1alpha1"
	mona "kmodules.xyz/monitoring-agent-api/api/v1"
	"kubedb.dev/apimachinery/apis"
	"kubedb.dev/apimachinery/apis/kubedb"
)

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

// Snapshot service account name.
func (p PgBouncer) SnapshotSAName() string {
	return fmt.Sprintf("%v-snapshot", p.OffshootName())
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

func (p PgBouncer) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crdutils.NewCustomResourceDefinition(crdutils.Config{
		Group:         SchemeGroupVersion.Group,
		Plural:        ResourcePluralPgBouncer,
		Singular:      ResourceSingularPgBouncer,
		Kind:          ResourceKindPgBouncer,
		ShortNames:    []string{ResourceCodePgBouncer},
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
		SpecDefinitionName:      "kubedb.dev/apimachinery/apis/kubedb/v1alpha1.PgBouncer",
		EnableValidation:        true,
		GetOpenAPIDefinitions:   GetOpenAPIDefinitions,
		EnableStatusSubresource: apis.EnableStatusSubresource,
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

func (p *PgBouncer) SetDefaults() {
	if p == nil {
		return
	}
	p.Spec.SetDefaults()
}

func (p *PgBouncerSpec) SetDefaults() {
	if p == nil {
		return
	}
}

//func (e *PgBouncerSpec) GetSecrets() []string {
//	if e == nil {
//		return nil
//	}
//
//	var secrets []string
//	if e.DatabaseSecret != nil {
//		secrets = append(secrets, e.DatabaseSecret.SecretName)
//	}
//	return secrets
//}
