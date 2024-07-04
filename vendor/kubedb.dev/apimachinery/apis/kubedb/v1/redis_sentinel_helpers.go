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
	"kmodules.xyz/client-go/policy/secomp"
	appcat "kmodules.xyz/custom-resources/apis/appcatalog/v1alpha1"
	mona "kmodules.xyz/monitoring-agent-api/api/v1"
	ofstv2 "kmodules.xyz/offshoot-api/api/v2"
	pslister "kubeops.dev/petset/client/listers/apps/v1"
)

func (*RedisSentinel) Hub() {}

func (rs RedisSentinel) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralRedisSentinel))
}

func (rs *RedisSentinel) AsOwner() *metav1.OwnerReference {
	return metav1.NewControllerRef(rs, SchemeGroupVersion.WithKind(ResourceKindRedisSentinel))
}

var _ apis.ResourceInfo = &RedisSentinel{}

func (rs RedisSentinel) OffshootName() string {
	return rs.Name
}

func (rs RedisSentinel) OffshootSelectors() map[string]string {
	return map[string]string{
		meta_util.NameLabelKey:      rs.ResourceFQN(),
		meta_util.InstanceLabelKey:  rs.Name,
		meta_util.ManagedByLabelKey: kubedb.GroupName,
	}
}

func (rs RedisSentinel) OffshootLabels() map[string]string {
	return rs.offshootLabels(rs.OffshootSelectors(), nil)
}

func (rs RedisSentinel) PodLabels() map[string]string {
	return rs.offshootLabels(rs.OffshootSelectors(), rs.Spec.PodTemplate.Labels)
}

func (rs RedisSentinel) PodControllerLabels() map[string]string {
	return rs.offshootLabels(rs.OffshootSelectors(), rs.Spec.PodTemplate.Controller.Labels)
}

func (rs RedisSentinel) ServiceLabels(alias ServiceAlias, extraLabels ...map[string]string) map[string]string {
	svcTemplate := GetServiceTemplate(rs.Spec.ServiceTemplates, alias)
	return rs.offshootLabels(meta_util.OverwriteKeys(rs.OffshootSelectors(), extraLabels...), svcTemplate.Labels)
}

func (rs RedisSentinel) offshootLabels(selector, override map[string]string) map[string]string {
	selector[meta_util.ComponentLabelKey] = kubedb.ComponentDatabase
	return meta_util.FilterKeys(kubedb.GroupName, selector, meta_util.OverwriteKeys(rs.Labels, override))
}

func (rs RedisSentinel) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", ResourcePluralRedisSentinel, kubedb.GroupName)
}

func (rs RedisSentinel) ResourceShortCode() string {
	return ResourceCodeRedisSentinel
}

func (rs RedisSentinel) ResourceKind() string {
	return ResourceKindRedisSentinel
}

func (rs RedisSentinel) ResourceSingular() string {
	return ResourceSingularRedisSentinel
}

func (rs RedisSentinel) ResourcePlural() string {
	return ResourcePluralRedisSentinel
}

func (rs RedisSentinel) GetAuthSecretName() string {
	if rs.Spec.AuthSecret != nil && rs.Spec.AuthSecret.Name != "" {
		return rs.Spec.AuthSecret.Name
	}
	return meta_util.NameWithSuffix(rs.OffshootName(), "auth")
}

func (rs RedisSentinel) GoverningServiceName() string {
	return meta_util.NameWithSuffix(rs.OffshootName(), "pods")
}

func (r RedisSentinel) Address() string {
	return fmt.Sprintf("%v.%v.svc:%d", r.Name, r.Namespace, kubedb.RedisSentinelPort)
}

func (rs RedisSentinel) ConfigSecretName() string {
	return rs.OffshootName()
}

type redisSentinelApp struct {
	*RedisSentinel
}

func (rs redisSentinelApp) Name() string {
	return rs.RedisSentinel.Name
}

func (rs redisSentinelApp) Type() appcat.AppType {
	return appcat.AppType(fmt.Sprintf("%s/%s", kubedb.GroupName, ResourceSingularRedisSentinel))
}

func (rs RedisSentinel) AppBindingMeta() appcat.AppBindingMeta {
	return &redisSentinelApp{&rs}
}

type redisSentinelStatsService struct {
	*RedisSentinel
}

func (rs redisSentinelStatsService) GetNamespace() string {
	return rs.RedisSentinel.GetNamespace()
}

func (rs redisSentinelStatsService) ServiceName() string {
	return rs.OffshootName() + "-stats"
}

func (rs redisSentinelStatsService) ServiceMonitorName() string {
	return rs.ServiceName()
}

func (rs redisSentinelStatsService) ServiceMonitorAdditionalLabels() map[string]string {
	return rs.OffshootLabels()
}

func (rs redisSentinelStatsService) Path() string {
	return kubedb.DefaultStatsPath
}

func (r redisSentinelStatsService) Scheme() string {
	return ""
}

func (r redisSentinelStatsService) TLSConfig() *promapi.TLSConfig {
	return nil
}

func (rs RedisSentinel) StatsService() mona.StatsAccessor {
	return &redisSentinelStatsService{&rs}
}

func (rs RedisSentinel) StatsServiceLabels() map[string]string {
	return rs.ServiceLabels(StatsServiceAlias, map[string]string{kubedb.LabelRole: kubedb.RoleStats})
}

func (rs *RedisSentinel) SetDefaults(rdVersion *catalog.RedisVersion) {
	if rs == nil {
		return
	}

	if rs.Spec.StorageType == "" {
		rs.Spec.StorageType = StorageTypeDurable
	}
	if rs.Spec.DeletionPolicy == "" {
		rs.Spec.DeletionPolicy = DeletionPolicyDelete
	}

	rs.setDefaultContainerSecurityContext(rdVersion, &rs.Spec.PodTemplate)
	rs.setDefaultContainerResourceLimits(&rs.Spec.PodTemplate)

	if rs.Spec.PodTemplate.Spec.ServiceAccountName == "" {
		rs.Spec.PodTemplate.Spec.ServiceAccountName = rs.OffshootName()
	}

	rs.SetTLSDefaults()
	rs.SetHealthCheckerDefaults()
	rs.Spec.Monitor.SetDefaults()
	if rs.Spec.Monitor != nil && rs.Spec.Monitor.Prometheus != nil {
		if rs.Spec.Monitor.Prometheus.Exporter.SecurityContext.RunAsUser == nil {
			rs.Spec.Monitor.Prometheus.Exporter.SecurityContext.RunAsUser = rdVersion.Spec.SecurityContext.RunAsUser
		}
		if rs.Spec.Monitor.Prometheus.Exporter.SecurityContext.RunAsGroup == nil {
			rs.Spec.Monitor.Prometheus.Exporter.SecurityContext.RunAsGroup = rdVersion.Spec.SecurityContext.RunAsUser
		}
	}
}

func (rs *RedisSentinel) SetHealthCheckerDefaults() {
	if rs.Spec.HealthChecker.PeriodSeconds == nil {
		rs.Spec.HealthChecker.PeriodSeconds = pointer.Int32P(10)
	}
	if rs.Spec.HealthChecker.TimeoutSeconds == nil {
		rs.Spec.HealthChecker.TimeoutSeconds = pointer.Int32P(10)
	}
	if rs.Spec.HealthChecker.FailureThreshold == nil {
		rs.Spec.HealthChecker.FailureThreshold = pointer.Int32P(1)
	}
}

func (rs *RedisSentinel) SetTLSDefaults() {
	if rs.Spec.TLS == nil || rs.Spec.TLS.IssuerRef == nil {
		return
	}
	rs.Spec.TLS.Certificates = kmapi.SetMissingSecretNameForCertificate(rs.Spec.TLS.Certificates, string(RedisServerCert), rs.CertificateName(RedisServerCert))
	rs.Spec.TLS.Certificates = kmapi.SetMissingSecretNameForCertificate(rs.Spec.TLS.Certificates, string(RedisClientCert), rs.CertificateName(RedisClientCert))
	rs.Spec.TLS.Certificates = kmapi.SetMissingSecretNameForCertificate(rs.Spec.TLS.Certificates, string(RedisMetricsExporterCert), rs.CertificateName(RedisMetricsExporterCert))
}

func (rs *RedisSentinel) GetPersistentSecrets() []string {
	if rs == nil {
		return nil
	}

	var secrets []string
	if rs.Spec.AuthSecret != nil {
		secrets = append(secrets, rs.Spec.AuthSecret.Name)
	}
	return secrets
}

func (rs *RedisSentinel) setDefaultContainerSecurityContext(rdVersion *catalog.RedisVersion, podTemplate *ofstv2.PodTemplateSpec) {
	if podTemplate == nil {
		podTemplate = &ofstv2.PodTemplateSpec{}
	}

	if podTemplate.Spec.SecurityContext == nil {
		podTemplate.Spec.SecurityContext = &core.PodSecurityContext{}
	}
	if podTemplate.Spec.SecurityContext.FSGroup == nil {
		podTemplate.Spec.SecurityContext.FSGroup = rdVersion.Spec.SecurityContext.RunAsUser
	}
	dbContainer := core_util.GetContainerByName(podTemplate.Spec.Containers, kubedb.RedisSentinelContainerName)
	if dbContainer == nil {
		dbContainer = &core.Container{
			Name: kubedb.RedisSentinelContainerName,
		}
	}
	if dbContainer.SecurityContext == nil {
		dbContainer.SecurityContext = &core.SecurityContext{}
	}
	rs.assignDefaultContainerSecurityContext(rdVersion, dbContainer.SecurityContext)
	podTemplate.Spec.Containers = core_util.UpsertContainer(podTemplate.Spec.Containers, *dbContainer)

	initContainer := core_util.GetContainerByName(podTemplate.Spec.InitContainers, kubedb.RedisSentinelInitContainerName)
	if initContainer == nil {
		initContainer = &core.Container{
			Name: kubedb.RedisSentinelInitContainerName,
		}
	}
	if initContainer.SecurityContext == nil {
		initContainer.SecurityContext = &core.SecurityContext{}
	}
	rs.assignDefaultContainerSecurityContext(rdVersion, initContainer.SecurityContext)
	podTemplate.Spec.InitContainers = core_util.UpsertContainer(podTemplate.Spec.InitContainers, *initContainer)
}

func (rs *RedisSentinel) assignDefaultContainerSecurityContext(rdVersion *catalog.RedisVersion, sc *core.SecurityContext) {
	if sc.AllowPrivilegeEscalation == nil {
		sc.AllowPrivilegeEscalation = pointer.BoolP(false)
	}
	if sc.Capabilities == nil {
		sc.Capabilities = &core.Capabilities{
			Drop: []core.Capability{"ALL"},
		}
	}
	if sc.RunAsNonRoot == nil {
		sc.RunAsNonRoot = pointer.BoolP(true)
	}
	if sc.RunAsUser == nil {
		sc.RunAsUser = rdVersion.Spec.SecurityContext.RunAsUser
	}
	if sc.RunAsGroup == nil {
		sc.RunAsGroup = rdVersion.Spec.SecurityContext.RunAsUser
	}
	if sc.SeccompProfile == nil {
		sc.SeccompProfile = secomp.DefaultSeccompProfile()
	}
}

func (rs *RedisSentinel) setDefaultContainerResourceLimits(podTemplate *ofstv2.PodTemplateSpec) {
	dbContainer := core_util.GetContainerByName(podTemplate.Spec.Containers, kubedb.RedisSentinelContainerName)
	if dbContainer != nil && (dbContainer.Resources.Requests == nil && dbContainer.Resources.Limits == nil) {
		apis.SetDefaultResourceLimits(&dbContainer.Resources, kubedb.DefaultResources)
	}

	initContainer := core_util.GetContainerByName(podTemplate.Spec.InitContainers, kubedb.RedisSentinelInitContainerName)
	if initContainer != nil && (initContainer.Resources.Requests == nil && initContainer.Resources.Limits == nil) {
		apis.SetDefaultResourceLimits(&initContainer.Resources, kubedb.DefaultInitContainerResource)
	}
}

// CertificateName returns the default certificate name and/or certificate secret name for a certificate alias
func (rs *RedisSentinel) CertificateName(alias RedisCertificateAlias) string {
	return meta_util.NameWithSuffix(rs.Name, fmt.Sprintf("%s-cert", string(alias)))
}

// GetCertSecretName returns the secret name for a certificate alias if any provide,
// otherwise returns default certificate secret name for the given alias.
func (rs *RedisSentinel) GetCertSecretName(alias RedisCertificateAlias) string {
	if rs.Spec.TLS != nil {
		name, ok := kmapi.GetCertificateSecretName(rs.Spec.TLS.Certificates, string(alias))
		if ok {
			return name
		}
	}
	return rs.CertificateName(alias)
}

func (rs *RedisSentinel) ReplicasAreReady(lister pslister.PetSetLister) (bool, string, error) {
	// Desire number of statefulSets
	expectedItems := 1
	return checkReplicas(lister.PetSets(rs.Namespace), labels.SelectorFromSet(rs.OffshootLabels()), expectedItems)
}
