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

func (*MySQL) Hub() {}

func (_ MySQL) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralMySQL))
}

func (m *MySQL) AsOwner() *metav1.OwnerReference {
	return metav1.NewControllerRef(m, SchemeGroupVersion.WithKind(ResourceKindMySQL))
}

var _ apis.ResourceInfo = &MySQL{}

func (m MySQL) OffshootName() string {
	return m.Name
}

func (m MySQL) OffshootSelectors() map[string]string {
	return m.offshootSelectors(kubedb.MySQLComponentDB)
}

func (m MySQL) RouterOffshootSelectors() map[string]string {
	return m.offshootSelectors(kubedb.MySQLComponentRouter)
}

func (m MySQL) offshootSelectors(component string) map[string]string {
	selectors := map[string]string{
		meta_util.NameLabelKey:      m.ResourceFQN(),
		meta_util.InstanceLabelKey:  m.Name,
		meta_util.ManagedByLabelKey: kubedb.GroupName,
	}
	if m.IsInnoDBCluster() {
		selectors[kubedb.MySQLComponentKey] = component
	}
	return selectors
}

func (m MySQL) OffshootLabels() map[string]string {
	return m.offshootLabels(m.OffshootSelectors(), nil)
}

func (m MySQL) PodLabels() map[string]string {
	return m.offshootLabels(m.OffshootSelectors(), m.Spec.PodTemplate.Labels)
}

func (m MySQL) PodControllerLabels() map[string]string {
	return m.offshootLabels(m.OffshootSelectors(), m.Spec.PodTemplate.Controller.Labels)
}

func (m MySQL) SidekickLabels(skName string) map[string]string {
	return meta_util.OverwriteKeys(nil, kubedb.CommonSidekickLabels(), map[string]string{
		meta_util.InstanceLabelKey: skName,
		kubedb.SidekickOwnerName:   m.Name,
		kubedb.SidekickOwnerKind:   m.ResourceFQN(),
	})
}

func (m MySQL) RouterOffshootLabels() map[string]string {
	return m.offshootLabels(m.RouterOffshootSelectors(), nil)
}

func (m MySQL) RouterPodLabels() map[string]string {
	return m.offshootLabels(m.RouterOffshootLabels(), m.Spec.Topology.InnoDBCluster.Router.PodTemplate.Labels)
}

func (m MySQL) RouterPodControllerLabels() map[string]string {
	return m.offshootLabels(m.RouterOffshootLabels(), m.Spec.Topology.InnoDBCluster.Router.PodTemplate.Controller.Labels)
}

func (m MySQL) ServiceLabels(alias ServiceAlias, extraLabels ...map[string]string) map[string]string {
	svcTemplate := GetServiceTemplate(m.Spec.ServiceTemplates, alias)
	return m.offshootLabels(meta_util.OverwriteKeys(m.OffshootSelectors(), extraLabels...), svcTemplate.Labels)
}

func (m MySQL) offshootLabels(selector, override map[string]string) map[string]string {
	selector[meta_util.ComponentLabelKey] = kubedb.ComponentDatabase
	return meta_util.FilterKeys(kubedb.GroupName, selector, meta_util.OverwriteKeys(nil, m.Labels, override))
}

func (m MySQL) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", ResourcePluralMySQL, kubedb.GroupName)
}

func (m MySQL) ResourceShortCode() string {
	return ResourceCodeMySQL
}

func (m MySQL) ResourceKind() string {
	return ResourceKindMySQL
}

func (m MySQL) ResourceSingular() string {
	return ResourceSingularMySQL
}

func (m MySQL) ResourcePlural() string {
	return ResourcePluralMySQL
}

func (m MySQL) ServiceName() string {
	return m.OffshootName()
}

func (m MySQL) StandbyServiceName() string {
	return meta_util.NameWithPrefix(m.ServiceName(), "standby")
}

func (m MySQL) GoverningServiceName() string {
	return meta_util.NameWithSuffix(m.ServiceName(), "pods")
}

func (m MySQL) PrimaryServiceDNS() string {
	return fmt.Sprintf("%s.%s.svc", m.ServiceName(), m.Namespace)
}

func (m MySQL) StandbyServiceDNS() string {
	return fmt.Sprintf("%s.%s.svc", m.StandbyServiceName(), m.Namespace)
}

func (m MySQL) Hosts() []string {
	replicas := 1
	if m.Spec.Replicas != nil {
		replicas = int(*m.Spec.Replicas)
	}
	hosts := make([]string, replicas)
	for i := 0; i < replicas; i++ {
		hosts[i] = fmt.Sprintf("%v-%d.%v.%v.svc", m.Name, i, m.GoverningServiceName(), m.Namespace)
	}
	return hosts
}

func (m MySQL) PeerName(idx int) string {
	return fmt.Sprintf("%s-%d.%s.%s", m.OffshootName(), idx, m.GoverningServiceName(), m.Namespace)
}

func (m MySQL) GetAuthSecretName() string {
	if m.Spec.AuthSecret != nil && m.Spec.AuthSecret.Name != "" {
		return m.Spec.AuthSecret.Name
	}
	return meta_util.NameWithSuffix(m.OffshootName(), "auth")
}

type mysqlApp struct {
	*MySQL
}

func (r mysqlApp) Name() string {
	return r.MySQL.Name
}

func (r mysqlApp) Type() appcat.AppType {
	return appcat.AppType(fmt.Sprintf("%s/%s", kubedb.GroupName, ResourceSingularMySQL))
}

func (m MySQL) AppBindingMeta() appcat.AppBindingMeta {
	return &mysqlApp{&m}
}

type mysqlStatsService struct {
	*MySQL
}

func (m mysqlStatsService) GetNamespace() string {
	return m.MySQL.GetNamespace()
}

func (m MySQL) GetNameSpacedName() string {
	return m.Namespace + "/" + m.Name
}

func (m mysqlStatsService) ServiceName() string {
	return m.OffshootName() + "-stats"
}

func (m mysqlStatsService) ServiceMonitorName() string {
	return m.ServiceName()
}

func (m mysqlStatsService) ServiceMonitorAdditionalLabels() map[string]string {
	return m.OffshootLabels()
}

func (m mysqlStatsService) Path() string {
	return kubedb.DefaultStatsPath
}

func (m mysqlStatsService) Scheme() string {
	return ""
}

func (m mysqlStatsService) TLSConfig() *promapi.TLSConfig {
	return nil
}

func (m MySQL) StatsService() mona.StatsAccessor {
	return &mysqlStatsService{&m}
}

func (m MySQL) StatsServiceLabels() map[string]string {
	return m.ServiceLabels(StatsServiceAlias, map[string]string{kubedb.LabelRole: kubedb.RoleStats})
}

func (m *MySQL) UsesGroupReplication() bool {
	return m.Spec.Topology != nil &&
		m.Spec.Topology.Mode != nil &&
		*m.Spec.Topology.Mode == MySQLModeGroupReplication
}

func (m *MySQL) IsInnoDBCluster() bool {
	return m.Spec.Topology != nil &&
		m.Spec.Topology.Mode != nil &&
		*m.Spec.Topology.Mode == MySQLModeInnoDBCluster
}

func (m *MySQL) IsRemoteReplica() bool {
	return m.Spec.Topology != nil &&
		m.Spec.Topology.Mode != nil &&
		*m.Spec.Topology.Mode == MySQLModeRemoteReplica
}

func (m *MySQL) IsSemiSync() bool {
	return m.Spec.Topology != nil &&
		m.Spec.Topology.Mode != nil &&
		*m.Spec.Topology.Mode == MySQLModeSemiSync
}

func (m *MySQL) SetDefaults(myVersion *v1alpha1.MySQLVersion) {
	if m == nil {
		return
	}
	if m.Spec.StorageType == "" {
		m.Spec.StorageType = StorageTypeDurable
	}
	if m.Spec.DeletionPolicy == "" {
		m.Spec.DeletionPolicy = DeletionPolicyDelete
	}

	if m.UsesGroupReplication() || m.IsInnoDBCluster() || m.IsSemiSync() {
		if m.Spec.Replicas == nil {
			m.Spec.Replicas = pointer.Int32P(kubedb.MySQLDefaultGroupSize)
		}
	} else {
		if m.Spec.Replicas == nil {
			m.Spec.Replicas = pointer.Int32P(1)
		}
	}

	m.setDefaultContainerSecurityContext(myVersion, &m.Spec.PodTemplate)

	if m.IsInnoDBCluster() {
		if m.Spec.Topology != nil && m.Spec.Topology.InnoDBCluster != nil {
			m.setDefaultInnoDBContainerSecurityContext(myVersion, m.Spec.Topology.InnoDBCluster.Router.PodTemplate)

			routerContainer := core_util.GetContainerByName(m.Spec.Topology.InnoDBCluster.Router.PodTemplate.Spec.Containers, kubedb.MySQLRouterContainerName)
			if routerContainer != nil && (routerContainer.Resources.Requests == nil && routerContainer.Resources.Limits == nil) {
				apis.SetDefaultResourceLimits(&routerContainer.Resources, kubedb.CoordinatorDefaultResources)
			}
		}
	}

	if m.Spec.PodTemplate.Spec.ServiceAccountName == "" {
		m.Spec.PodTemplate.Spec.ServiceAccountName = m.OffshootName()
	}

	m.setDefaultContainerResourceLimits(&m.Spec.PodTemplate)
	m.SetTLSDefaults()
	m.SetHealthCheckerDefaults()

	m.Spec.Monitor.SetDefaults()
	if m.Spec.Monitor != nil && m.Spec.Monitor.Prometheus != nil {
		if m.Spec.Monitor.Prometheus.Exporter.SecurityContext.RunAsUser == nil {
			m.Spec.Monitor.Prometheus.Exporter.SecurityContext.RunAsUser = myVersion.Spec.SecurityContext.RunAsUser
		}
		if m.Spec.Monitor.Prometheus.Exporter.SecurityContext.RunAsGroup == nil {
			m.Spec.Monitor.Prometheus.Exporter.SecurityContext.RunAsGroup = myVersion.Spec.SecurityContext.RunAsUser
		}
	}
	if m.Spec.Init != nil && m.Spec.Init.Archiver != nil && m.Spec.Init.Archiver.ReplicationStrategy == nil {
		m.Spec.Init.Archiver.ReplicationStrategy = ptr.To(ReplicationStrategyNone)
	}
}

func (m *MySQL) SetTLSDefaults() {
	if m.Spec.TLS == nil || m.Spec.TLS.IssuerRef == nil {
		return
	}
	m.Spec.TLS.Certificates = kmapi.SetMissingSecretNameForCertificate(m.Spec.TLS.Certificates, string(MySQLServerCert), m.CertificateName(MySQLServerCert))
	m.Spec.TLS.Certificates = kmapi.SetMissingSecretNameForCertificate(m.Spec.TLS.Certificates, string(MySQLClientCert), m.CertificateName(MySQLClientCert))
	m.Spec.TLS.Certificates = kmapi.SetMissingSecretNameForCertificate(m.Spec.TLS.Certificates, string(MySQLMetricsExporterCert), m.CertificateName(MySQLMetricsExporterCert))
}

func (m *MySQL) SetHealthCheckerDefaults() {
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

func (m *MySQLSpec) GetPersistentSecrets() []string {
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
func (m *MySQL) CertificateName(alias MySQLCertificateAlias) string {
	return meta_util.NameWithSuffix(m.Name, fmt.Sprintf("%s-cert", string(alias)))
}

// GetCertSecretName returns the secret name for a certificate alias if any
// otherwise returns default certificate secret name for the given alias.
func (m *MySQL) GetCertSecretName(alias MySQLCertificateAlias) string {
	if m.Spec.TLS != nil {
		name, ok := kmapi.GetCertificateSecretName(m.Spec.TLS.Certificates, string(alias))
		if ok {
			return name
		}
	}
	return m.CertificateName(alias)
}

func (m *MySQL) ReplicasAreReady(lister pslister.PetSetLister) (bool, string, error) {
	// Desire number of statefulSets
	expectedItems := 1
	return checkReplicas(lister.PetSets(m.Namespace), labels.SelectorFromSet(m.OffshootLabels()), expectedItems)
}

func MySQLRequireSSLArg() string {
	return "--require-secure-transport=ON"
}

func MySQLExporterTLSArg() string {
	return "--config.my-cnf=/etc/mysql/certs/exporter.cnf"
}

func (m *MySQL) MySQLTLSArgs() []string {
	tlsArgs := []string{
		"--ssl-capath=/etc/mysql/certs",
		"--ssl-ca=/etc/mysql/certs/ca.crt",
		"--ssl-cert=/etc/mysql/certs/server.crt",
		"--ssl-key=/etc/mysql/certs/server.key",
	}
	if m.Spec.RequireSSL {
		tlsArgs = append(tlsArgs, MySQLRequireSSLArg())
	}
	return tlsArgs
}

func (m *MySQL) GetRouterName() string {
	return fmt.Sprintf("%s-router", m.Name)
}

func (m *MySQL) setDefaultContainerSecurityContext(myVersion *v1alpha1.MySQLVersion, podTemplate *ofstv2.PodTemplateSpec) {
	if podTemplate == nil {
		podTemplate = &ofstv2.PodTemplateSpec{}
	}

	if podTemplate.Spec.SecurityContext == nil {
		podTemplate.Spec.SecurityContext = &core.PodSecurityContext{}
	}
	if podTemplate.Spec.SecurityContext.FSGroup == nil {
		podTemplate.Spec.SecurityContext.FSGroup = myVersion.Spec.SecurityContext.RunAsUser
	}

	dbContainer := core_util.GetContainerByName(podTemplate.Spec.Containers, kubedb.MySQLContainerName)
	if dbContainer == nil {
		dbContainer = &core.Container{
			Name: kubedb.MySQLContainerName,
		}
	}
	if dbContainer.SecurityContext == nil {
		dbContainer.SecurityContext = &core.SecurityContext{}
	}
	m.assignDefaultContainerSecurityContext(myVersion, dbContainer.SecurityContext)
	podTemplate.Spec.Containers = core_util.UpsertContainer(podTemplate.Spec.Containers, *dbContainer)

	initContainer := core_util.GetContainerByName(podTemplate.Spec.InitContainers, kubedb.MySQLInitContainerName)
	if initContainer == nil {
		initContainer = &core.Container{
			Name: kubedb.MySQLInitContainerName,
		}
	}
	if initContainer.SecurityContext == nil {
		initContainer.SecurityContext = &core.SecurityContext{}
	}
	m.assignDefaultContainerSecurityContext(myVersion, initContainer.SecurityContext)
	podTemplate.Spec.InitContainers = core_util.UpsertContainer(podTemplate.Spec.InitContainers, *initContainer)

	if m.IsInnoDBCluster() || m.IsSemiSync() || m.UsesGroupReplication() {
		coordinatorContainer := core_util.GetContainerByName(podTemplate.Spec.Containers, kubedb.MySQLCoordinatorContainerName)
		if coordinatorContainer == nil {
			coordinatorContainer = &core.Container{
				Name: kubedb.MySQLCoordinatorContainerName,
			}
		}
		if coordinatorContainer.SecurityContext == nil {
			coordinatorContainer.SecurityContext = &core.SecurityContext{}
		}
		m.assignDefaultContainerSecurityContext(myVersion, coordinatorContainer.SecurityContext)
		podTemplate.Spec.Containers = core_util.UpsertContainer(podTemplate.Spec.Containers, *coordinatorContainer)
	}
}

func (m *MySQL) setDefaultInnoDBContainerSecurityContext(myVersion *v1alpha1.MySQLVersion, podTemplate *ofstv2.PodTemplateSpec) {
	if podTemplate == nil {
		podTemplate = &ofstv2.PodTemplateSpec{}
	}

	if podTemplate.Spec.SecurityContext == nil {
		podTemplate.Spec.SecurityContext = &core.PodSecurityContext{}
	}
	if podTemplate.Spec.SecurityContext.FSGroup == nil {
		podTemplate.Spec.SecurityContext.FSGroup = myVersion.Spec.SecurityContext.RunAsUser
	}

	routerContainer := core_util.GetContainerByName(podTemplate.Spec.Containers, kubedb.MySQLRouterContainerName)
	if routerContainer == nil {
		routerContainer = &core.Container{
			Name: kubedb.MySQLRouterContainerName,
		}
	}
	if routerContainer.SecurityContext == nil {
		routerContainer.SecurityContext = &core.SecurityContext{}
	}
	m.assignDefaultContainerSecurityContext(myVersion, routerContainer.SecurityContext)
	podTemplate.Spec.Containers = core_util.UpsertContainer(podTemplate.Spec.Containers, *routerContainer)
}

func (m *MySQL) assignDefaultContainerSecurityContext(myVersion *v1alpha1.MySQLVersion, sc *core.SecurityContext) {
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
		sc.RunAsUser = myVersion.Spec.SecurityContext.RunAsUser
	}
	if sc.RunAsGroup == nil {
		sc.RunAsGroup = myVersion.Spec.SecurityContext.RunAsUser
	}
	if sc.SeccompProfile == nil {
		sc.SeccompProfile = secomp.DefaultSeccompProfile()
	}
}

func (m *MySQL) setDefaultContainerResourceLimits(podTemplate *ofstv2.PodTemplateSpec) {
	dbContainer := core_util.GetContainerByName(podTemplate.Spec.Containers, kubedb.MySQLContainerName)
	if dbContainer != nil && (dbContainer.Resources.Requests == nil && dbContainer.Resources.Limits == nil) {
		apis.SetDefaultResourceLimits(&dbContainer.Resources, kubedb.DefaultResources)
	}

	initContainer := core_util.GetContainerByName(podTemplate.Spec.InitContainers, kubedb.MySQLInitContainerName)
	if initContainer != nil && (initContainer.Resources.Requests == nil && initContainer.Resources.Limits == nil) {
		apis.SetDefaultResourceLimits(&initContainer.Resources, kubedb.DefaultInitContainerResource)
	}

	if m.IsInnoDBCluster() || m.IsSemiSync() || m.UsesGroupReplication() {
		coordinatorContainer := core_util.GetContainerByName(podTemplate.Spec.Containers, kubedb.MySQLCoordinatorContainerName)
		if coordinatorContainer != nil && (coordinatorContainer.Resources.Requests == nil && coordinatorContainer.Resources.Limits == nil) {
			apis.SetDefaultResourceLimits(&coordinatorContainer.Resources, kubedb.CoordinatorDefaultResources)
		}
	}
}
