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
	"k8s.io/utils/ptr"
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

func (*MariaDB) Hub() {}

func (_ MariaDB) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralMariaDB))
}

func (m *MariaDB) AsOwner() *metav1.OwnerReference {
	return metav1.NewControllerRef(m, SchemeGroupVersion.WithKind(ResourceKindMariaDB))
}

var _ apis.ResourceInfo = &MariaDB{}

func (m MariaDB) OffshootName() string {
	return m.Name
}

func (m MariaDB) OffshootSelectors() map[string]string {
	return map[string]string{
		meta_util.NameLabelKey:      m.ResourceFQN(),
		meta_util.InstanceLabelKey:  m.Name,
		meta_util.ManagedByLabelKey: kubedb.GroupName,
	}
}

func (m MariaDB) OffshootLabels() map[string]string {
	return m.offshootLabels(m.OffshootSelectors(), nil)
}

func (m MariaDB) PodLabels() map[string]string {
	return m.offshootLabels(m.OffshootSelectors(), m.Spec.PodTemplate.Labels)
}

func (m MariaDB) PodControllerLabels() map[string]string {
	return m.offshootLabels(m.OffshootSelectors(), m.Spec.PodTemplate.Controller.Labels)
}

func (m MariaDB) SidekickLabels(skName string) map[string]string {
	return meta_util.OverwriteKeys(nil, kubedb.CommonSidekickLabels(), map[string]string{
		meta_util.InstanceLabelKey: skName,
		kubedb.SidekickOwnerName:   m.Name,
		kubedb.SidekickOwnerKind:   m.ResourceFQN(),
	})
}

func (m MariaDB) offshootLabels(selector, override map[string]string) map[string]string {
	selector[meta_util.ComponentLabelKey] = kubedb.ComponentDatabase
	return meta_util.FilterKeys(kubedb.GroupName, selector, meta_util.OverwriteKeys(nil, m.Labels, override))
}

func (m MariaDB) ServiceLabels(alias ServiceAlias, extraLabels ...map[string]string) map[string]string {
	svcTemplate := GetServiceTemplate(m.Spec.ServiceTemplates, alias)
	return m.offshootLabels(meta_util.OverwriteKeys(m.OffshootSelectors(), extraLabels...), svcTemplate.Labels)
}

func (m MariaDB) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", ResourcePluralMariaDB, kubedb.GroupName)
}

func (m MariaDB) ResourceShortCode() string {
	return ResourceCodeMariaDB
}

func (m MariaDB) ResourceKind() string {
	return ResourceKindMariaDB
}

func (m MariaDB) ResourceSingular() string {
	return ResourceSingularMariaDB
}

func (m MariaDB) ResourcePlural() string {
	return ResourcePluralMariaDB
}

func (m MariaDB) ServiceName() string {
	return m.OffshootName()
}

func (m MariaDB) IsCluster() bool {
	return pointer.Int32(m.Spec.Replicas) > 1
}

func (m MariaDB) GoverningServiceName() string {
	return meta_util.NameWithSuffix(m.ServiceName(), "pods")
}

func (m MariaDB) PeerName(idx int) string {
	return fmt.Sprintf("%s-%d.%s.%s", m.OffshootName(), idx, m.GoverningServiceName(), m.Namespace)
}

func (m MariaDB) GetAuthSecretName() string {
	if m.Spec.AuthSecret != nil && m.Spec.AuthSecret.Name != "" {
		return m.Spec.AuthSecret.Name
	}
	return meta_util.NameWithSuffix(m.OffshootName(), "auth")
}

func (m MariaDB) ClusterName() string {
	return m.OffshootName()
}

type mariadbApp struct {
	*MariaDB
}

func (m mariadbApp) Name() string {
	return m.MariaDB.Name
}

func (m mariadbApp) Type() appcat.AppType {
	return appcat.AppType(fmt.Sprintf("%s/%s", kubedb.GroupName, ResourceSingularMariaDB))
}

func (m MariaDB) AppBindingMeta() appcat.AppBindingMeta {
	return &mariadbApp{&m}
}

type mariadbStatsService struct {
	*MariaDB
}

func (m mariadbStatsService) GetNamespace() string {
	return m.MariaDB.GetNamespace()
}

func (m mariadbStatsService) ServiceName() string {
	return m.OffshootName() + "-stats"
}

func (m mariadbStatsService) ServiceMonitorName() string {
	return m.ServiceName()
}

func (m mariadbStatsService) ServiceMonitorAdditionalLabels() map[string]string {
	return m.OffshootLabels()
}

func (m mariadbStatsService) Path() string {
	return kubedb.DefaultStatsPath
}

func (m mariadbStatsService) Scheme() string {
	return ""
}

func (m mariadbStatsService) TLSConfig() *promapi.TLSConfig {
	return nil
}

func (m MariaDB) StatsService() mona.StatsAccessor {
	return &mariadbStatsService{&m}
}

func (m MariaDB) StatsServiceLabels() map[string]string {
	return m.ServiceLabels(StatsServiceAlias, map[string]string{kubedb.LabelRole: kubedb.RoleStats})
}

func (m MariaDB) PrimaryServiceDNS() string {
	return fmt.Sprintf("%s.%s.svc", m.ServiceName(), m.Namespace)
}

func (m *MariaDB) SetDefaults(mdVersion *v1alpha1.MariaDBVersion) {
	if m == nil {
		return
	}

	if m.Spec.StorageType == "" {
		m.Spec.StorageType = StorageTypeDurable
	}
	if m.Spec.DeletionPolicy == "" {
		m.Spec.DeletionPolicy = DeletionPolicyDelete
	}

	if m.Spec.Replicas == nil {
		m.Spec.Replicas = pointer.Int32P(1)
	}

	if m.Spec.PodTemplate.Spec.ServiceAccountName == "" {
		m.Spec.PodTemplate.Spec.ServiceAccountName = m.OffshootName()
	}
	if m.Spec.Init != nil && m.Spec.Init.Archiver != nil && m.Spec.Init.Archiver.ReplicationStrategy == nil {
		m.Spec.Init.Archiver.ReplicationStrategy = ptr.To(ReplicationStrategyNone)
	}
	m.setDefaultContainerSecurityContext(mdVersion, &m.Spec.PodTemplate)
	m.setDefaultContainerResourceLimits(&m.Spec.PodTemplate)
	m.SetTLSDefaults()
	m.SetHealthCheckerDefaults()

	m.Spec.Monitor.SetDefaults()
	if m.Spec.Monitor != nil && m.Spec.Monitor.Prometheus != nil {
		if m.Spec.Monitor.Prometheus.Exporter.SecurityContext.RunAsUser == nil {
			m.Spec.Monitor.Prometheus.Exporter.SecurityContext.RunAsUser = mdVersion.Spec.SecurityContext.RunAsUser
		}
		if m.Spec.Monitor.Prometheus.Exporter.SecurityContext.RunAsGroup == nil {
			m.Spec.Monitor.Prometheus.Exporter.SecurityContext.RunAsGroup = mdVersion.Spec.SecurityContext.RunAsUser
		}
	}
}

func (m *MariaDB) setDefaultContainerSecurityContext(mdVersion *v1alpha1.MariaDBVersion, podTemplate *ofstv2.PodTemplateSpec) {
	if podTemplate == nil {
		return
	}

	if podTemplate.Spec.SecurityContext == nil {
		podTemplate.Spec.SecurityContext = &core.PodSecurityContext{}
	}
	if podTemplate.Spec.SecurityContext.FSGroup == nil {
		podTemplate.Spec.SecurityContext.FSGroup = mdVersion.Spec.SecurityContext.RunAsUser
	}
	dbContainer := core_util.GetContainerByName(podTemplate.Spec.Containers, kubedb.MariaDBContainerName)
	if dbContainer == nil {
		dbContainer = &core.Container{
			Name: kubedb.MariaDBContainerName,
		}
	}
	if dbContainer.SecurityContext == nil {
		dbContainer.SecurityContext = &core.SecurityContext{}
	}
	m.assignDefaultContainerSecurityContext(mdVersion, dbContainer.SecurityContext)
	podTemplate.Spec.Containers = core_util.UpsertContainer(podTemplate.Spec.Containers, *dbContainer)

	initContainer := core_util.GetContainerByName(podTemplate.Spec.InitContainers, kubedb.MariaDBInitContainerName)
	if initContainer == nil {
		initContainer = &core.Container{
			Name: kubedb.MariaDBInitContainerName,
		}
	}
	if initContainer.SecurityContext == nil {
		initContainer.SecurityContext = &core.SecurityContext{}
	}
	m.assignDefaultContainerSecurityContext(mdVersion, initContainer.SecurityContext)
	podTemplate.Spec.InitContainers = core_util.UpsertContainer(podTemplate.Spec.InitContainers, *initContainer)

	if m.IsCluster() {
		coordinatorContainer := core_util.GetContainerByName(podTemplate.Spec.Containers, kubedb.MariaDBCoordinatorContainerName)
		if coordinatorContainer == nil {
			coordinatorContainer = &core.Container{
				Name: kubedb.MariaDBCoordinatorContainerName,
			}
		}
		if coordinatorContainer.SecurityContext == nil {
			coordinatorContainer.SecurityContext = &core.SecurityContext{}
		}
		m.assignDefaultContainerSecurityContext(mdVersion, coordinatorContainer.SecurityContext)
		podTemplate.Spec.Containers = core_util.UpsertContainer(podTemplate.Spec.Containers, *coordinatorContainer)
	}
}

func (m *MariaDB) assignDefaultContainerSecurityContext(mdVersion *v1alpha1.MariaDBVersion, sc *core.SecurityContext) {
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
		sc.RunAsUser = mdVersion.Spec.SecurityContext.RunAsUser
	}
	if sc.RunAsGroup == nil {
		sc.RunAsGroup = mdVersion.Spec.SecurityContext.RunAsUser
	}
	if sc.SeccompProfile == nil {
		sc.SeccompProfile = secomp.DefaultSeccompProfile()
	}
}

func (m *MariaDB) setDefaultContainerResourceLimits(podTemplate *ofstv2.PodTemplateSpec) {
	dbContainer := core_util.GetContainerByName(podTemplate.Spec.Containers, kubedb.MariaDBContainerName)
	if dbContainer != nil && (dbContainer.Resources.Requests == nil && dbContainer.Resources.Limits == nil) {
		apis.SetDefaultResourceLimits(&dbContainer.Resources, kubedb.DefaultResources)
	}

	initContainer := core_util.GetContainerByName(podTemplate.Spec.InitContainers, kubedb.MariaDBInitContainerName)
	if initContainer != nil && (initContainer.Resources.Requests == nil && initContainer.Resources.Limits == nil) {
		apis.SetDefaultResourceLimits(&initContainer.Resources, kubedb.DefaultInitContainerResource)
	}
	if m.IsCluster() {
		coordinatorContainer := core_util.GetContainerByName(podTemplate.Spec.Containers, kubedb.MariaDBCoordinatorContainerName)
		if coordinatorContainer != nil && (coordinatorContainer.Resources.Requests == nil && coordinatorContainer.Resources.Limits == nil) {
			apis.SetDefaultResourceLimits(&coordinatorContainer.Resources, kubedb.CoordinatorDefaultResources)
		}
	}
}

func (m *MariaDB) SetHealthCheckerDefaults() {
	if m.Spec.HealthChecker.PeriodSeconds == nil {
		m.Spec.HealthChecker.PeriodSeconds = pointer.Int32P(10)
	}
	if m.Spec.HealthChecker.TimeoutSeconds == nil {
		m.Spec.HealthChecker.TimeoutSeconds = pointer.Int32P(10)
	}
	if m.Spec.HealthChecker.FailureThreshold == nil {
		m.Spec.HealthChecker.FailureThreshold = pointer.Int32P(1)
	}
}

func (m *MariaDB) SetTLSDefaults() {
	if m.Spec.TLS == nil || m.Spec.TLS.IssuerRef == nil {
		return
	}
	m.Spec.TLS.Certificates = kmapi.SetMissingSecretNameForCertificate(m.Spec.TLS.Certificates, string(MariaDBServerCert), m.CertificateName(MariaDBServerCert))
	m.Spec.TLS.Certificates = kmapi.SetMissingSecretNameForCertificate(m.Spec.TLS.Certificates, string(MariaDBClientCert), m.CertificateName(MariaDBClientCert))
	m.Spec.TLS.Certificates = kmapi.SetMissingSecretNameForCertificate(m.Spec.TLS.Certificates, string(MariaDBExporterCert), m.CertificateName(MariaDBExporterCert))
}

func (m *MariaDBSpec) GetPersistentSecrets() []string {
	if m == nil {
		return nil
	}

	var secrets []string
	if m.AuthSecret != nil {
		secrets = append(secrets, m.AuthSecret.Name)
	}
	return secrets
}

// CertificateName returns the default certificate name and/or certificate secret name for a certificate alias
func (m *MariaDB) CertificateName(alias MariaDBCertificateAlias) string {
	return meta_util.NameWithSuffix(m.Name, fmt.Sprintf("%s-cert", string(alias)))
}

// GetCertSecretName returns the secret name for a certificate alias if any,
// otherwise returns default certificate secret name for the given alias.
func (m *MariaDB) GetCertSecretName(alias MariaDBCertificateAlias) string {
	if m.Spec.TLS != nil {
		name, ok := kmapi.GetCertificateSecretName(m.Spec.TLS.Certificates, string(alias))
		if ok {
			return name
		}
	}
	return m.CertificateName(alias)
}

func (m *MariaDB) ReplicasAreReady(lister pslister.PetSetLister) (bool, string, error) {
	// Desire number of statefulSets
	expectedItems := 1
	return checkReplicas(lister.PetSets(m.Namespace), labels.SelectorFromSet(m.OffshootLabels()), expectedItems)
}

func (m *MariaDB) InlineConfigSecretName() string {
	return meta_util.NameWithSuffix(m.Name, "inline")
}

func (m *MariaDB) CertMountPath(alias MariaDBCertificateAlias) string {
	return filepath.Join(kubedb.PerconaXtraDBCertMountPath, string(alias))
}

func (m *MariaDB) CertFilePath(certAlias MariaDBCertificateAlias, certFileName string) string {
	return filepath.Join(m.CertMountPath(certAlias), certFileName)
}
