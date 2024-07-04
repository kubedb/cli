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
	"path/filepath"

	"kubedb.dev/apimachinery/apis"
	"kubedb.dev/apimachinery/apis/catalog/v1alpha1"
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

func (*PerconaXtraDB) Hub() {}

func (_ PerconaXtraDB) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralPerconaXtraDB))
}

func (p *PerconaXtraDB) AsOwner() *metav1.OwnerReference {
	return metav1.NewControllerRef(p, SchemeGroupVersion.WithKind(ResourceKindPerconaXtraDB))
}

var _ apis.ResourceInfo = &PerconaXtraDB{}

func (p PerconaXtraDB) OffshootName() string {
	return p.Name
}

func (p PerconaXtraDB) OffshootSelectors() map[string]string {
	return map[string]string{
		meta_util.NameLabelKey:      p.ResourceFQN(),
		meta_util.InstanceLabelKey:  p.Name,
		meta_util.ManagedByLabelKey: kubedb.GroupName,
	}
}

func (p PerconaXtraDB) OffshootLabels() map[string]string {
	return p.offshootLabels(p.OffshootSelectors(), nil)
}

func (p PerconaXtraDB) PodLabels() map[string]string {
	return p.offshootLabels(p.OffshootSelectors(), p.Spec.PodTemplate.Labels)
}

func (p PerconaXtraDB) PodControllerLabels() map[string]string {
	return p.offshootLabels(p.OffshootSelectors(), p.Spec.PodTemplate.Controller.Labels)
}

func (p PerconaXtraDB) ServiceLabels(alias ServiceAlias, extraLabels ...map[string]string) map[string]string {
	svcTemplate := GetServiceTemplate(p.Spec.ServiceTemplates, alias)
	return p.offshootLabels(meta_util.OverwriteKeys(p.OffshootSelectors(), extraLabels...), svcTemplate.Labels)
}

func (p PerconaXtraDB) offshootLabels(selector, override map[string]string) map[string]string {
	selector[meta_util.ComponentLabelKey] = kubedb.ComponentDatabase
	return meta_util.FilterKeys(kubedb.GroupName, selector, meta_util.OverwriteKeys(nil, p.Labels, override))
}

func (p PerconaXtraDB) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", ResourcePluralPerconaXtraDB, kubedb.GroupName)
}

func (p PerconaXtraDB) ResourceShortCode() string {
	return ResourceCodePerconaXtraDB
}

func (p PerconaXtraDB) ResourceKind() string {
	return ResourceKindPerconaXtraDB
}

func (p PerconaXtraDB) ResourceSingular() string {
	return ResourceSingularPerconaXtraDB
}

func (p PerconaXtraDB) ResourcePlural() string {
	return ResourcePluralPerconaXtraDB
}

func (p PerconaXtraDB) ServiceName() string {
	return p.OffshootName()
}

func (p PerconaXtraDB) GoverningServiceName() string {
	return meta_util.NameWithSuffix(p.ServiceName(), "pods")
}

func (p PerconaXtraDB) PeerName(idx int) string {
	return fmt.Sprintf("%s-%d.%s.%s", p.OffshootName(), idx, p.GoverningServiceName(), p.Namespace)
}

func (p PerconaXtraDB) GetAuthSecretName() string {
	if p.Spec.AuthSecret != nil && p.Spec.AuthSecret.Name != "" {
		return p.Spec.AuthSecret.Name
	}
	return meta_util.NameWithSuffix(p.OffshootName(), "auth")
}

func (p PerconaXtraDB) GetReplicationSecretName() string {
	if p.Spec.SystemUserSecrets != nil &&
		p.Spec.SystemUserSecrets.ReplicationUserSecret != nil &&
		p.Spec.SystemUserSecrets.ReplicationUserSecret.Name != "" {
		return p.Spec.SystemUserSecrets.ReplicationUserSecret.Name
	}
	return meta_util.NameWithSuffix(p.OffshootName(), "replication")
}

func (p PerconaXtraDB) GetMonitorSecretName() string {
	if p.Spec.SystemUserSecrets != nil &&
		p.Spec.SystemUserSecrets.MonitorUserSecret != nil &&
		p.Spec.SystemUserSecrets.MonitorUserSecret.Name != "" {
		return p.Spec.SystemUserSecrets.MonitorUserSecret.Name
	}
	return meta_util.NameWithSuffix(p.OffshootName(), "monitor")
}

func (p PerconaXtraDB) ClusterName() string {
	return p.OffshootName()
}

type perconaXtraDBApp struct {
	*PerconaXtraDB
}

func (p perconaXtraDBApp) Name() string {
	return p.PerconaXtraDB.Name
}

func (p perconaXtraDBApp) Type() appcat.AppType {
	return appcat.AppType(fmt.Sprintf("%s/%s", kubedb.GroupName, ResourceSingularPerconaXtraDB))
}

func (p PerconaXtraDB) AppBindingMeta() appcat.AppBindingMeta {
	return &perconaXtraDBApp{&p}
}

type perconaXtraDBStatsService struct {
	*PerconaXtraDB
}

func (p perconaXtraDBStatsService) GetNamespace() string {
	return p.PerconaXtraDB.GetNamespace()
}

func (p perconaXtraDBStatsService) ServiceName() string {
	return p.OffshootName() + "-stats"
}

func (p perconaXtraDBStatsService) ServiceMonitorName() string {
	return p.ServiceName()
}

func (p perconaXtraDBStatsService) ServiceMonitorAdditionalLabels() map[string]string {
	return p.OffshootLabels()
}

func (p perconaXtraDBStatsService) Path() string {
	return kubedb.DefaultStatsPath
}

func (p perconaXtraDBStatsService) Scheme() string {
	return ""
}

func (p perconaXtraDBStatsService) TLSConfig() *promapi.TLSConfig {
	return nil
}

func (p PerconaXtraDB) StatsService() mona.StatsAccessor {
	return &perconaXtraDBStatsService{&p}
}

func (p PerconaXtraDB) StatsServiceLabels() map[string]string {
	return p.ServiceLabels(StatsServiceAlias, map[string]string{kubedb.LabelRole: kubedb.RoleStats})
}

func (p PerconaXtraDB) PrimaryServiceDNS() string {
	return fmt.Sprintf("%s.%s.svc", p.ServiceName(), p.Namespace)
}

func (p *PerconaXtraDB) SetDefaults(pVersion *v1alpha1.PerconaXtraDBVersion) {
	if p == nil {
		return
	}

	if p.Spec.Replicas == nil {
		p.Spec.Replicas = pointer.Int32P(kubedb.PerconaXtraDBDefaultClusterSize)
	}

	if p.Spec.StorageType == "" {
		p.Spec.StorageType = StorageTypeDurable
	}
	if p.Spec.DeletionPolicy == "" {
		p.Spec.DeletionPolicy = DeletionPolicyDelete
	}

	if p.Spec.PodTemplate.Spec.ServiceAccountName == "" {
		p.Spec.PodTemplate.Spec.ServiceAccountName = p.OffshootName()
	}

	// Need to set FSGroup equal to  p.Spec.PodTemplate.Spec.ContainerSecurityContext.RunAsGroup.
	// So that /var/pv directory have the group permission for the RunAsGroup user GID.
	// Otherwise, We will get write permission denied.
	p.setDefaultContainerSecurityContext(pVersion, &p.Spec.PodTemplate)
	p.setDefaultContainerResourceLimits(&p.Spec.PodTemplate)
	p.SetTLSDefaults()
	p.Spec.Monitor.SetDefaults()
	if p.Spec.Monitor != nil && p.Spec.Monitor.Prometheus != nil {
		if p.Spec.Monitor.Prometheus.Exporter.SecurityContext.RunAsUser == nil {
			p.Spec.Monitor.Prometheus.Exporter.SecurityContext.RunAsUser = pVersion.Spec.SecurityContext.RunAsUser
		}
		if p.Spec.Monitor.Prometheus.Exporter.SecurityContext.RunAsGroup == nil {
			p.Spec.Monitor.Prometheus.Exporter.SecurityContext.RunAsGroup = pVersion.Spec.SecurityContext.RunAsUser
		}
	}
}

func (p *PerconaXtraDB) setDefaultContainerSecurityContext(pVersion *v1alpha1.PerconaXtraDBVersion, podTemplate *ofstv2.PodTemplateSpec) {
	if podTemplate == nil {
		return
	}

	if podTemplate.Spec.SecurityContext == nil {
		podTemplate.Spec.SecurityContext = &core.PodSecurityContext{}
	}
	if podTemplate.Spec.SecurityContext.FSGroup == nil {
		podTemplate.Spec.SecurityContext.FSGroup = pVersion.Spec.SecurityContext.RunAsUser
	}
	dbContainer := core_util.GetContainerByName(podTemplate.Spec.Containers, kubedb.PerconaXtraDBContainerName)
	if dbContainer == nil {
		dbContainer = &core.Container{
			Name: kubedb.PerconaXtraDBContainerName,
		}
	}
	if dbContainer.SecurityContext == nil {
		dbContainer.SecurityContext = &core.SecurityContext{}
	}
	p.assignDefaultContainerSecurityContext(pVersion, dbContainer.SecurityContext)
	podTemplate.Spec.Containers = core_util.UpsertContainer(podTemplate.Spec.Containers, *dbContainer)

	initContainer := core_util.GetContainerByName(podTemplate.Spec.InitContainers, kubedb.PerconaXtraDBInitContainerName)
	if initContainer == nil {
		initContainer = &core.Container{
			Name: kubedb.PerconaXtraDBInitContainerName,
		}
	}
	if initContainer.SecurityContext == nil {
		initContainer.SecurityContext = &core.SecurityContext{}
	}
	p.assignDefaultContainerSecurityContext(pVersion, initContainer.SecurityContext)
	podTemplate.Spec.InitContainers = core_util.UpsertContainer(podTemplate.Spec.InitContainers, *initContainer)

	coordinatorContainer := core_util.GetContainerByName(podTemplate.Spec.Containers, kubedb.PerconaXtraDBCoordinatorContainerName)
	if coordinatorContainer == nil {
		coordinatorContainer = &core.Container{
			Name: kubedb.PerconaXtraDBCoordinatorContainerName,
		}
	}
	if coordinatorContainer.SecurityContext == nil {
		coordinatorContainer.SecurityContext = &core.SecurityContext{}
	}
	p.assignDefaultContainerSecurityContext(pVersion, coordinatorContainer.SecurityContext)
	podTemplate.Spec.Containers = core_util.UpsertContainer(podTemplate.Spec.Containers, *coordinatorContainer)
}

func (p *PerconaXtraDB) assignDefaultContainerSecurityContext(pVersion *v1alpha1.PerconaXtraDBVersion, sc *core.SecurityContext) {
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
		sc.RunAsUser = pVersion.Spec.SecurityContext.RunAsUser
	}
	if sc.RunAsGroup == nil {
		sc.RunAsGroup = pVersion.Spec.SecurityContext.RunAsUser
	}
	if sc.SeccompProfile == nil {
		sc.SeccompProfile = secomp.DefaultSeccompProfile()
	}
}

func (p *PerconaXtraDB) setDefaultContainerResourceLimits(podTemplate *ofstv2.PodTemplateSpec) {
	dbContainer := core_util.GetContainerByName(podTemplate.Spec.Containers, kubedb.PerconaXtraDBContainerName)
	if dbContainer != nil && (dbContainer.Resources.Requests == nil && dbContainer.Resources.Limits == nil) {
		apis.SetDefaultResourceLimits(&dbContainer.Resources, kubedb.DefaultResources)
	}

	initContainer := core_util.GetContainerByName(podTemplate.Spec.InitContainers, kubedb.PerconaXtraDBInitContainerName)
	if initContainer != nil && (initContainer.Resources.Requests == nil && initContainer.Resources.Limits == nil) {
		apis.SetDefaultResourceLimits(&initContainer.Resources, kubedb.DefaultInitContainerResource)
	}

	coordinatorContainer := core_util.GetContainerByName(podTemplate.Spec.Containers, kubedb.PerconaXtraDBCoordinatorContainerName)
	if coordinatorContainer != nil && (coordinatorContainer.Resources.Requests == nil && coordinatorContainer.Resources.Limits == nil) {
		apis.SetDefaultResourceLimits(&coordinatorContainer.Resources, kubedb.CoordinatorDefaultResources)
	}
}

func (p *PerconaXtraDB) SetHealthCheckerDefaults() {
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

func (p *PerconaXtraDB) SetTLSDefaults() {
	if p.Spec.TLS == nil || p.Spec.TLS.IssuerRef == nil {
		return
	}
	p.Spec.TLS.Certificates = kmapi.SetMissingSecretNameForCertificate(p.Spec.TLS.Certificates, string(PerconaXtraDBServerCert), p.CertificateName(PerconaXtraDBServerCert))
	p.Spec.TLS.Certificates = kmapi.SetMissingSecretNameForCertificate(p.Spec.TLS.Certificates, string(PerconaXtraDBClientCert), p.CertificateName(PerconaXtraDBClientCert))
	p.Spec.TLS.Certificates = kmapi.SetMissingSecretNameForCertificate(p.Spec.TLS.Certificates, string(PerconaXtraDBExporterCert), p.CertificateName(PerconaXtraDBExporterCert))
}

// CertificateName returns the default certificate name and/or certificate secret name for a certificate alias
func (p *PerconaXtraDBSpec) GetPersistentSecrets() []string {
	if p == nil {
		return nil
	}

	var secrets []string
	if p.AuthSecret != nil {
		secrets = append(secrets, p.AuthSecret.Name)
	}
	if p.SystemUserSecrets != nil && p.SystemUserSecrets.ReplicationUserSecret != nil {
		secrets = append(secrets, p.SystemUserSecrets.ReplicationUserSecret.Name)
	}
	if p.SystemUserSecrets != nil && p.SystemUserSecrets.MonitorUserSecret != nil {
		secrets = append(secrets, p.SystemUserSecrets.MonitorUserSecret.Name)
	}
	return secrets
}

// CertificateName returns the default certificate name and/or certificate secret name for a certificate alias
func (p *PerconaXtraDB) CertificateName(alias PerconaXtraDBCertificateAlias) string {
	return meta_util.NameWithSuffix(p.Name, fmt.Sprintf("%s-cert", string(alias)))
}

// GetCertSecretName returns the secret name for a certificate alias if any,
// otherwise returns default certificate secret name for the given alias.
func (p *PerconaXtraDB) GetCertSecretName(alias PerconaXtraDBCertificateAlias) string {
	if p.Spec.TLS != nil {
		name, ok := kmapi.GetCertificateSecretName(p.Spec.TLS.Certificates, string(alias))
		if ok {
			return name
		}
	}
	return p.CertificateName(alias)
}

func (p *PerconaXtraDB) ReplicasAreReady(lister pslister.PetSetLister) (bool, string, error) {
	// Desire number of statefulSets
	expectedItems := 1
	return checkReplicas(lister.PetSets(p.Namespace), labels.SelectorFromSet(p.OffshootLabels()), expectedItems)
}

func (p *PerconaXtraDB) CertMountPath(alias PerconaXtraDBCertificateAlias) string {
	return filepath.Join(kubedb.PerconaXtraDBCertMountPath, string(alias))
}

func (p *PerconaXtraDB) CertFilePath(certAlias PerconaXtraDBCertificateAlias, certFileName string) string {
	return filepath.Join(p.CertMountPath(certAlias), certFileName)
}
