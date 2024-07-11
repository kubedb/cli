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

package v1

import (
	"fmt"

	"kubedb.dev/apimachinery/apis"
	catalog "kubedb.dev/apimachinery/apis/catalog/v1alpha1"
	"kubedb.dev/apimachinery/apis/kubedb"
	"kubedb.dev/apimachinery/crds"

	promapi "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"gomodules.xyz/pointer"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	kmapi "kmodules.xyz/client-go/api/v1"
	"kmodules.xyz/client-go/apiextensions"
	core_util "kmodules.xyz/client-go/core/v1"
	meta_util "kmodules.xyz/client-go/meta"
	appcat "kmodules.xyz/custom-resources/apis/appcatalog/v1alpha1"
	mona "kmodules.xyz/monitoring-agent-api/api/v1"
	ofstv2 "kmodules.xyz/offshoot-api/api/v2"
	ofst_util "kmodules.xyz/offshoot-api/util"
	pslister "kubeops.dev/petset/client/listers/apps/v1"
)

func (*PgBouncer) Hub() {}

func (p PgBouncer) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralPgBouncer))
}

func (p *PgBouncer) AsOwner() *metav1.OwnerReference {
	return metav1.NewControllerRef(p, SchemeGroupVersion.WithKind(ResourceKindPgBouncer))
}

var _ apis.ResourceInfo = &PgBouncer{}

func (p PgBouncer) OffshootName() string {
	return p.Name
}

func (p PgBouncer) OffshootSelectors() map[string]string {
	return map[string]string{
		meta_util.NameLabelKey:      p.ResourceFQN(),
		meta_util.InstanceLabelKey:  p.Name,
		meta_util.ManagedByLabelKey: kubedb.GroupName,
	}
}

func (p PgBouncer) OffshootLabels() map[string]string {
	return p.offshootLabels(p.OffshootSelectors(), nil)
}

func (p PgBouncer) PodLabels(backendSecretRV string) map[string]string {
	podLabels := p.offshootLabels(p.OffshootSelectors(), p.Spec.PodTemplate.Labels)
	podLabels[kubedb.BackendSecretResourceVersion] = backendSecretRV
	return podLabels
}

func (p PgBouncer) PodControllerLabels() map[string]string {
	return p.offshootLabels(p.OffshootSelectors(), p.Spec.PodTemplate.Controller.Labels)
}

func (p PgBouncer) ServiceLabels(alias ServiceAlias, extraLabels ...map[string]string) map[string]string {
	svcTemplate := GetServiceTemplate(p.Spec.ServiceTemplates, alias)
	return p.offshootLabels(meta_util.OverwriteKeys(p.OffshootSelectors(), extraLabels...), svcTemplate.Labels)
}

func (p PgBouncer) offshootLabels(selector, override map[string]string) map[string]string {
	selector[meta_util.ComponentLabelKey] = kubedb.ComponentConnectionPooler
	return meta_util.FilterKeys(kubedb.GroupName, selector, meta_util.OverwriteKeys(nil, p.Labels, override))
}

func (p PgBouncer) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", ResourcePluralPgBouncer, kubedb.GroupName)
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

func (p PgBouncer) GoverningServiceName() string {
	return meta_util.NameWithSuffix(p.ServiceName(), "pods")
}

func (p PgBouncer) GetAuthSecretName() string {
	if p.Spec.AuthSecret != nil && p.Spec.AuthSecret.Name != "" {
		return p.Spec.AuthSecret.Name
	}
	return meta_util.NameWithSuffix(p.OffshootName(), "auth")
}

func (p PgBouncer) GetBackendSecretName() string {
	return meta_util.NameWithSuffix(p.OffshootName(), "backend")
}

func (p PgBouncer) ConfigSecretName() string {
	return meta_util.NameWithSuffix(p.ServiceName(), "config")
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
	return p.ServiceName()
}

func (p pgbouncerStatsService) ServiceMonitorAdditionalLabels() map[string]string {
	return p.OffshootLabels()
}

func (p pgbouncerStatsService) Path() string {
	return kubedb.DefaultStatsPath
}

func (p pgbouncerStatsService) Scheme() string {
	return ""
}

func (p pgbouncerStatsService) TLSConfig() *promapi.TLSConfig {
	return nil
}

func (p PgBouncer) StatsService() mona.StatsAccessor {
	return &pgbouncerStatsService{&p}
}

func (p PgBouncer) StatsServiceLabels() map[string]string {
	return p.ServiceLabels(StatsServiceAlias, map[string]string{kubedb.LabelRole: kubedb.RoleStats})
}

func (p PgBouncer) ReplicasServiceName() string {
	return fmt.Sprintf("%v-replicas", p.Name)
}

func (p *PgBouncer) SetDefaults(pgBouncerVersion *catalog.PgBouncerVersion, usesAcme bool) {
	if p == nil {
		return
	}

	if p.Spec.DeletionPolicy == "" {
		p.Spec.DeletionPolicy = PgBouncerDeletionPolicyDelete
	}

	p.setConnectionPoolConfigDefaults()

	if p.Spec.TLS != nil {
		if p.Spec.SSLMode == "" {
			p.Spec.SSLMode = PgBouncerSSLModeVerifyFull
		}
	} else {
		if p.Spec.SSLMode == "" {
			p.Spec.SSLMode = PgBouncerSSLModeDisable
		}
	}

	p.setPgBouncerContainerDefaults(&p.Spec.PodTemplate)

	p.SetSecurityContext(pgBouncerVersion)
	if p.Spec.TLS != nil {
		p.SetTLSDefaults(usesAcme)
	}

	p.Spec.Monitor.SetDefaults()
	if p.Spec.Monitor != nil && p.Spec.Monitor.Prometheus != nil {
		if p.Spec.Monitor.Prometheus.Exporter.SecurityContext.RunAsUser == nil {
			p.Spec.Monitor.Prometheus.Exporter.SecurityContext.RunAsUser = pgBouncerVersion.Spec.SecurityContext.RunAsUser
		}
		if p.Spec.Monitor.Prometheus.Exporter.SecurityContext.RunAsGroup == nil {
			p.Spec.Monitor.Prometheus.Exporter.SecurityContext.RunAsGroup = pgBouncerVersion.Spec.SecurityContext.RunAsUser
		}
	}
	dbContainer := core_util.GetContainerByName(p.Spec.PodTemplate.Spec.Containers, ResourceSingularPgBouncer)
	if dbContainer != nil && (dbContainer.Resources.Requests == nil && dbContainer.Resources.Limits == nil) {
		apis.SetDefaultResourceLimits(&dbContainer.Resources, kubedb.DefaultResourcesMemoryIntensive)
	}
}

func (p *PgBouncer) setPgBouncerContainerDefaults(podTemplate *ofstv2.PodTemplateSpec) {
	if podTemplate == nil {
		return
	}
	container := ofst_util.EnsureContainerExists(podTemplate, kubedb.PgBouncerContainerName)
	p.setContainerDefaultResources(container, *kubedb.DefaultResources.DeepCopy())
}

func (p *PgBouncer) setContainerDefaultResources(container *core.Container, defaultResources core.ResourceRequirements) {
	if container.Resources.Requests == nil && container.Resources.Limits == nil {
		apis.SetDefaultResourceLimits(&container.Resources, defaultResources)
	}
}

func (p *PgBouncer) SetTLSDefaults(usesAcme bool) {
	if p.Spec.TLS == nil || p.Spec.TLS.IssuerRef == nil {
		return
	}

	p.Spec.TLS.Certificates = kmapi.SetMissingSecretNameForCertificate(p.Spec.TLS.Certificates, string(PgBouncerServerCert), p.CertificateName(PgBouncerServerCert))
	if !usesAcme {
		p.Spec.TLS.Certificates = kmapi.SetMissingSecretNameForCertificate(p.Spec.TLS.Certificates, string(PgBouncerClientCert), p.CertificateName(PgBouncerClientCert))
		p.Spec.TLS.Certificates = kmapi.SetMissingSecretNameForCertificate(p.Spec.TLS.Certificates, string(PgBouncerMetricsExporterCert), p.CertificateName(PgBouncerMetricsExporterCert))
	}
}

// CertificateName returns the default certificate name and/or certificate secret name for a certificate alias
func (p *PgBouncer) CertificateName(alias PgBouncerCertificateAlias) string {
	return meta_util.NameWithSuffix(p.Name, fmt.Sprintf("%s-cert", string(alias)))
}

// GetPersistentSecrets returns auth secret and config secret of a pgbouncer object
func (p *PgBouncer) GetPersistentSecrets() []string {
	if p == nil {
		return nil
	}
	var secrets []string
	secrets = append(secrets, p.GetAuthSecretName())
	secrets = append(secrets, p.GetBackendSecretName())
	secrets = append(secrets, p.ConfigSecretName())

	return secrets
}

// GetCertSecretName returns the secret name for a certificate alias if any provide,
// otherwise returns default certificate secret name for the given alias.
func (p *PgBouncer) GetCertSecretName(alias PgBouncerCertificateAlias) string {
	if p.Spec.TLS != nil {
		name, ok := kmapi.GetCertificateSecretName(p.Spec.TLS.Certificates, string(alias))
		if ok {
			return name
		}
	}
	return p.CertificateName(alias)
}

func (p *PgBouncer) ReplicasAreReady(lister pslister.PetSetLister) (bool, string, error) {
	// Desire number of petSets
	expectedItems := 1
	return checkReplicas(lister.PetSets(p.Namespace), labels.SelectorFromSet(p.OffshootLabels()), expectedItems)
}

func (p *PgBouncer) SetHealthCheckerDefaults() {
	if p.Spec.HealthChecker.PeriodSeconds == nil {
		p.Spec.HealthChecker.PeriodSeconds = pointer.Int32P(10)
	}
	if p.Spec.HealthChecker.TimeoutSeconds == nil {
		p.Spec.HealthChecker.TimeoutSeconds = pointer.Int32P(10)
	}
	if p.Spec.HealthChecker.FailureThreshold == nil {
		p.Spec.HealthChecker.FailureThreshold = pointer.Int32P(1)
	}
}

func (p *PgBouncer) setConnectionPoolConfigDefaults() {
	if p.Spec.ConnectionPool == nil {
		p.Spec.ConnectionPool = &ConnectionPoolConfig{}
	}
	if p.Spec.ConnectionPool.Port == nil {
		p.Spec.ConnectionPool.Port = pointer.Int32P(5432)
	}
	if p.Spec.ConnectionPool.PoolMode == "" {
		p.Spec.ConnectionPool.PoolMode = kubedb.PgBouncerDefaultPoolMode
	}
	if p.Spec.ConnectionPool.MaxClientConnections == nil {
		p.Spec.ConnectionPool.MaxClientConnections = pointer.Int64P(100)
	}
	if p.Spec.ConnectionPool.DefaultPoolSize == nil {
		p.Spec.ConnectionPool.DefaultPoolSize = pointer.Int64P(20)
	}
	if p.Spec.ConnectionPool.MinPoolSize == nil {
		p.Spec.ConnectionPool.MinPoolSize = pointer.Int64P(0)
	}
	if p.Spec.ConnectionPool.ReservePoolSize == nil {
		p.Spec.ConnectionPool.ReservePoolSize = pointer.Int64P(0)
	}
	if p.Spec.ConnectionPool.ReservePoolTimeoutSeconds == nil {
		p.Spec.ConnectionPool.ReservePoolTimeoutSeconds = pointer.Int64P(5)
	}
	if p.Spec.ConnectionPool.MaxDBConnections == nil {
		p.Spec.ConnectionPool.MaxDBConnections = pointer.Int64P(0)
	}
	if p.Spec.ConnectionPool.MaxUserConnections == nil {
		p.Spec.ConnectionPool.MaxUserConnections = pointer.Int64P(0)
	}
	if p.Spec.ConnectionPool.StatsPeriodSeconds == nil {
		p.Spec.ConnectionPool.StatsPeriodSeconds = pointer.Int64P(60)
	}
	if p.Spec.ConnectionPool.AuthType == "" {
		p.Spec.ConnectionPool.AuthType = PgBouncerClientAuthModeMD5
	}
	if p.Spec.ConnectionPool.IgnoreStartupParameters == "" {
		p.Spec.ConnectionPool.IgnoreStartupParameters = kubedb.PgBouncerDefaultIgnoreStartupParameters
	}
}

func (p *PgBouncer) SetSecurityContext(pgBouncerVersion *catalog.PgBouncerVersion) {
	container := core_util.GetContainerByName(p.Spec.PodTemplate.Spec.Containers, kubedb.PgBouncerContainerName)
	if container == nil {
		container = &core.Container{
			Name: kubedb.PgBouncerContainerName,
		}
	}
	if container.SecurityContext == nil {
		container.SecurityContext = &core.SecurityContext{
			RunAsUser: func() *int64 {
				if p.Spec.PodTemplate.Spec.SecurityContext == nil || p.Spec.PodTemplate.Spec.SecurityContext.RunAsUser == nil {
					return pgBouncerVersion.Spec.SecurityContext.RunAsUser
				}
				return p.Spec.PodTemplate.Spec.SecurityContext.RunAsUser
			}(),
			RunAsGroup: func() *int64 {
				if p.Spec.PodTemplate.Spec.SecurityContext == nil || p.Spec.PodTemplate.Spec.SecurityContext.RunAsGroup == nil {
					return pgBouncerVersion.Spec.SecurityContext.RunAsUser
				}
				return p.Spec.PodTemplate.Spec.SecurityContext.RunAsGroup
			}(),
			Privileged: pointer.BoolP(false),
		}
	} else {
		if container.SecurityContext.RunAsUser == nil {
			container.SecurityContext.RunAsUser = pgBouncerVersion.Spec.SecurityContext.RunAsUser
		}
		if container.SecurityContext.RunAsGroup == nil {
			container.SecurityContext.RunAsGroup = container.SecurityContext.RunAsUser
		}
	}

	if p.Spec.PodTemplate.Spec.SecurityContext == nil {
		p.Spec.PodTemplate.Spec.SecurityContext = &core.PodSecurityContext{
			RunAsUser:  container.SecurityContext.RunAsUser,
			RunAsGroup: container.SecurityContext.RunAsGroup,
		}
	} else {
		if p.Spec.PodTemplate.Spec.SecurityContext.RunAsUser == nil {
			p.Spec.PodTemplate.Spec.SecurityContext.RunAsUser = container.SecurityContext.RunAsUser
		}
		if p.Spec.PodTemplate.Spec.SecurityContext.RunAsGroup == nil {
			p.Spec.PodTemplate.Spec.SecurityContext.RunAsGroup = container.SecurityContext.RunAsGroup
		}
	}

	// Need to set FSGroup equal to  p.Spec.PodTemplate.Spec.ContainerSecurityContext.RunAsGroup.
	// So that /var/pv directory have the group permission for the RunAsGroup user GID.
	// Otherwise, We will get write permission denied.
	p.Spec.PodTemplate.Spec.SecurityContext.FSGroup = container.SecurityContext.RunAsGroup
	isPgbouncerContainerPresent := core_util.GetContainerByName(p.Spec.PodTemplate.Spec.Containers, kubedb.PgBouncerContainerName)
	if isPgbouncerContainerPresent == nil {
		core_util.UpsertContainer(p.Spec.PodTemplate.Spec.Containers, *container)
	}
}
