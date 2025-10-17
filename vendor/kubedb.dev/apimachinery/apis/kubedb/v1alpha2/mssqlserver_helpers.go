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
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"kubedb.dev/apimachinery/apis"
	catalog "kubedb.dev/apimachinery/apis/catalog/v1alpha1"
	"kubedb.dev/apimachinery/apis/kubedb"
	"kubedb.dev/apimachinery/crds"

	promapi "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"gomodules.xyz/pointer"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
	"k8s.io/utils/ptr"
	kmapi "kmodules.xyz/client-go/api/v1"
	"kmodules.xyz/client-go/apiextensions"
	coreutil "kmodules.xyz/client-go/core/v1"
	meta_util "kmodules.xyz/client-go/meta"
	metautil "kmodules.xyz/client-go/meta"
	"kmodules.xyz/client-go/policy/secomp"
	appcat "kmodules.xyz/custom-resources/apis/appcatalog/v1alpha1"
	mona "kmodules.xyz/monitoring-agent-api/api/v1"
	ofst "kmodules.xyz/offshoot-api/api/v2"
	pslister "kubeops.dev/petset/client/listers/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type MSSQLServerApp struct {
	*MSSQLServer
}

func (m *MSSQLServer) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralMSSQLServer))
}

func (m MSSQLServerApp) Name() string {
	return m.MSSQLServer.Name
}

func (m MSSQLServerApp) Type() appcat.AppType {
	return appcat.AppType(fmt.Sprintf("%s/%s", kubedb.GroupName, ResourceSingularMSSQLServer))
}

func (m *MSSQLServer) ResourceKind() string {
	return ResourceKindMSSQLServer
}

func (m *MSSQLServer) ResourcePlural() string {
	return ResourcePluralMSSQLServer
}

func (m *MSSQLServer) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", m.ResourcePlural(), kubedb.GroupName)
}

// Owner returns owner reference to resources
func (m *MSSQLServer) Owner() *meta.OwnerReference {
	return meta.NewControllerRef(m, SchemeGroupVersion.WithKind(m.ResourceKind()))
}

func (m *MSSQLServer) OffshootName() string {
	return m.Name
}

func (m *MSSQLServer) ServiceName() string {
	return m.OffshootName()
}

func (m *MSSQLServer) SecondaryServiceName() string {
	return metautil.NameWithPrefix(m.ServiceName(), string(SecondaryServiceAlias))
}

func (m *MSSQLServer) GoverningServiceName() string {
	return metautil.NameWithSuffix(m.ServiceName(), "pods")
}

func (m *MSSQLServer) DefaultUserCredSecretName(username string) string {
	return metautil.NameWithSuffix(m.Name, strings.ReplaceAll(fmt.Sprintf("%s-cred", username), "_", "-"))
}

func (m *MSSQLServer) offshootLabels(selector, override map[string]string) map[string]string {
	selector[metautil.ComponentLabelKey] = kubedb.ComponentDatabase
	return metautil.FilterKeys(kubedb.GroupName, selector, metautil.OverwriteKeys(nil, m.Labels, override))
}

func (m *MSSQLServer) ServiceLabels(alias ServiceAlias, extraLabels ...map[string]string) map[string]string {
	svcTemplate := GetServiceTemplate(m.Spec.ServiceTemplates, alias)
	return m.offshootLabels(metautil.OverwriteKeys(m.OffshootSelectors(), extraLabels...), svcTemplate.Labels)
}

func (m *MSSQLServer) OffshootLabels() map[string]string {
	return m.offshootLabels(m.OffshootSelectors(), nil)
}

func (m *MSSQLServer) OffshootSelectors(extraSelectors ...map[string]string) map[string]string {
	selector := map[string]string{
		metautil.NameLabelKey:      m.ResourceFQN(),
		metautil.InstanceLabelKey:  m.Name,
		metautil.ManagedByLabelKey: kubedb.GroupName,
	}
	return metautil.OverwriteKeys(selector, extraSelectors...)
}

type mssqlserverStatsService struct {
	*MSSQLServer
}

func (m mssqlserverStatsService) GetNamespace() string {
	return m.MSSQLServer.GetNamespace()
}

func (m mssqlserverStatsService) ServiceName() string {
	return m.OffshootName() + "-stats"
}

func (m mssqlserverStatsService) ServiceMonitorName() string {
	return m.ServiceName()
}

func (m mssqlserverStatsService) ServiceMonitorAdditionalLabels() map[string]string {
	return m.OffshootLabels()
}

func (m mssqlserverStatsService) Path() string {
	return kubedb.DefaultStatsPath
}

func (m mssqlserverStatsService) Scheme() string {
	return ""
}

func (m mssqlserverStatsService) TLSConfig() *promapi.TLSConfig {
	return nil
}

func (m MSSQLServer) StatsService() mona.StatsAccessor {
	return &mssqlserverStatsService{&m}
}

func (m MSSQLServer) StatsServiceLabels() map[string]string {
	return m.ServiceLabels(StatsServiceAlias, map[string]string{kubedb.LabelRole: kubedb.RoleStats})
}

func (m *MSSQLServer) IsAvailabilityGroup() bool {
	return m.Spec.Topology != nil &&
		m.Spec.Topology.Mode != nil &&
		*m.Spec.Topology.Mode == MSSQLServerModeAvailabilityGroup
}

func (m *MSSQLServer) IsDistributedAG() bool {
	return m.Spec.Topology != nil &&
		m.Spec.Topology.Mode != nil &&
		*m.Spec.Topology.Mode == MSSQLServerModeDistributedAG &&
		m.Spec.Topology.DistributedAG != nil
}

func (m *MSSQLServer) IsCluster() bool {
	return m.IsAvailabilityGroup() || m.IsDistributedAG()
}

func (m *MSSQLServer) DistributedAGName() string {
	return "dag"
}

func (m *MSSQLServer) IsStandalone() bool {
	return m.Spec.Topology == nil
}

func (m *MSSQLServer) PVCName(alias string) string {
	return metautil.NameWithSuffix(m.Name, alias)
}

func (m *MSSQLServer) PodLabels(extraLabels ...map[string]string) map[string]string {
	return m.offshootLabels(metautil.OverwriteKeys(m.OffshootSelectors(), extraLabels...), m.Spec.PodTemplate.Labels)
}

func (m *MSSQLServer) PodLabel(podTemplate *ofst.PodTemplateSpec) map[string]string {
	if podTemplate != nil && podTemplate.Labels != nil {
		return m.offshootLabels(m.OffshootSelectors(), m.Spec.PodTemplate.Labels)
	}
	return m.offshootLabels(m.OffshootSelectors(), nil)
}

func (m *MSSQLServer) ConfigSecretName() string {
	return metautil.NameWithSuffix(m.OffshootName(), "config")
}

func (m *MSSQLServer) PetSetName() string {
	return m.OffshootName()
}

func (m *MSSQLServer) ServiceAccountName() string {
	return m.OffshootName()
}

func (m *MSSQLServer) AvailabilityGroupName() string {
	// Get the database name
	dbName := m.Name

	// Regular expression pattern to match allowed characters (alphanumeric only)
	allowedPattern := regexp.MustCompile(`[a-zA-Z0-9]+`)

	// Extract valid characters from the database name
	validChars := allowedPattern.FindAllString(dbName, -1)

	// Ensure that the availability group name is not empty
	var availabilityGroupName string
	if len(validChars) == 0 {
		klog.Warningf("Database name '%s' contains no valid characters for the availability group name. Setting availability group name to 'DefaultGroupName'.", dbName)
		availabilityGroupName = "DefaultGroupName" // Provide a default name if the database name contains no valid characters
	} else {
		// Concatenate the valid characters to form the availability group name
		availabilityGroupName = ""
		for _, char := range validChars {
			availabilityGroupName += char
		}
	}

	return availabilityGroupName
}

func (m MSSQLServer) SidekickLabels(skName string) map[string]string {
	return meta_util.OverwriteKeys(nil, kubedb.CommonSidekickLabels(), map[string]string{
		meta_util.InstanceLabelKey: skName,
		kubedb.SidekickOwnerName:   m.Name,
		kubedb.SidekickOwnerKind:   m.ResourceFQN(),
	})
}

func (m *MSSQLServer) PodControllerLabels(extraLabels ...map[string]string) map[string]string {
	return m.offshootLabels(metautil.OverwriteKeys(m.OffshootSelectors(), extraLabels...), m.Spec.PodTemplate.Controller.Labels)
}

func (m *MSSQLServer) PodControllerLabel(podTemplate *ofst.PodTemplateSpec) map[string]string {
	if podTemplate != nil && podTemplate.Controller.Labels != nil {
		return m.offshootLabels(m.OffshootSelectors(), podTemplate.Controller.Labels)
	}
	return m.offshootLabels(m.OffshootSelectors(), nil)
}

func (m *MSSQLServer) GetPersistentSecrets() []string {
	var secrets []string
	if m.Spec.AuthSecret != nil {
		secrets = append(secrets, m.Spec.AuthSecret.Name)
	}

	secrets = append(secrets, m.EndpointCertSecretName())
	secrets = append(secrets, m.DbmLoginSecretName())
	secrets = append(secrets, m.MasterKeySecretName())
	secrets = append(secrets, m.ConfigSecretName())

	return secrets
}

func (m *MSSQLServer) AppBindingMeta() appcat.AppBindingMeta {
	return &MSSQLServerApp{m}
}

func (m MSSQLServer) SetHealthCheckerDefaults() {
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

func (m MSSQLServer) GetAuthSecretName() string {
	if m.Spec.AuthSecret != nil && m.Spec.AuthSecret.Name != "" {
		return m.Spec.AuthSecret.Name
	}
	return metautil.NameWithSuffix(m.OffshootName(), "auth")
}

func (m *MSSQLServer) CAProviderClassName() string {
	return metautil.NameWithSuffix(m.OffshootName(), "ca-provider")
}

func (m *MSSQLServer) DbmLoginSecretName() string {
	if m.Spec.Topology != nil && m.Spec.Topology.AvailabilityGroup != nil && m.Spec.Topology.AvailabilityGroup.LoginSecretName != "" {
		return m.Spec.Topology.AvailabilityGroup.LoginSecretName
	}
	return metautil.NameWithSuffix(m.OffshootName(), "dbm-login")
}

func (m *MSSQLServer) MasterKeySecretName() string {
	if m.Spec.Topology != nil && m.Spec.Topology.AvailabilityGroup != nil && m.Spec.Topology.AvailabilityGroup.MasterKeySecretName != "" {
		return m.Spec.Topology.AvailabilityGroup.MasterKeySecretName
	}
	return metautil.NameWithSuffix(m.OffshootName(), "master-key")
}

func (m *MSSQLServer) EndpointCertSecretName() string {
	if m.Spec.Topology != nil && m.Spec.Topology.AvailabilityGroup != nil && m.Spec.Topology.AvailabilityGroup.EndpointCertSecretName != "" {
		return m.Spec.Topology.AvailabilityGroup.EndpointCertSecretName
	}
	return metautil.NameWithSuffix(m.OffshootName(), "endpoint-cert")
}

// CertificateName returns the default certificate name and/or certificate secret name for a certificate alias
func (m *MSSQLServer) CertificateName(alias MSSQLServerCertificateAlias) string {
	return metautil.NameWithSuffix(m.Name, fmt.Sprintf("%s-cert", string(alias)))
}

func (m *MSSQLServer) SecretName(alias MSSQLServerCertificateAlias) string {
	return metautil.NameWithSuffix(m.Name, string(alias))
}

// GetCertSecretName returns the secret name for a certificate alias if any
// otherwise returns default certificate secret name for the given alias.
func (m *MSSQLServer) GetCertSecretName(alias MSSQLServerCertificateAlias) string {
	if m.Spec.TLS != nil {
		name, ok := kmapi.GetCertificateSecretName(m.Spec.TLS.Certificates, string(alias))
		if ok {
			return name
		}
	}
	return m.CertificateName(alias)
}

func (m *MSSQLServer) GetNameSpacedName() string {
	return m.Namespace + "/" + m.Name
}

func (m *MSSQLServer) PrimaryServiceDNS() string {
	return fmt.Sprintf("%s.%s.svc", m.ServiceName(), m.Namespace)
}

func (m *MSSQLServer) SetDefaults(kc client.Client) {
	if m == nil {
		return
	}
	if m.Spec.StorageType == "" {
		m.Spec.StorageType = StorageTypeDurable
	}
	if m.Spec.DeletionPolicy == "" {
		m.Spec.DeletionPolicy = DeletionPolicyDelete
	}

	if m.Spec.AuthSecret == nil {
		m.Spec.AuthSecret = &SecretReference{}
	}
	if m.Spec.AuthSecret.Kind == "" {
		m.Spec.AuthSecret.Kind = kubedb.ResourceKindSecret
	}

	if m.IsStandalone() {
		if m.Spec.Replicas == nil {
			m.Spec.Replicas = pointer.Int32P(1)
		}
	} else if m.IsCluster() {
		if m.Spec.Topology.AvailabilityGroup == nil {
			m.Spec.Topology.AvailabilityGroup = &MSSQLServerAvailabilityGroupSpec{}
		}
		if m.Spec.Topology.AvailabilityGroup.LeaderElection == nil {
			m.Spec.Topology.AvailabilityGroup.LeaderElection = &MSSQLServerLeaderElectionConfig{
				// The upper limit of election timeout is 50000ms (50s), which should only be used when deploying a
				// globally-distributed etcd cluster. A reasonable round-trip time for the continental United States is around 130-150ms,
				// and the time between US and Japan is around 350-400ms. If the network has uneven performance or regular packet
				// delays/loss then it is possible that a couple of retries may be necessary to successfully send a packet.
				// So 5s is a safe upper limit of global round-trip time. As the election timeout should be an order of magnitude
				// bigger than broadcast time, in the case of ~5s for a globally distributed cluster, then 50 seconds becomes
				// a reasonable maximum.
				Period: meta.Duration{Duration: 300 * time.Millisecond},
				// the amount of HeartbeatTick can be missed before the failOver
				ElectionTick: 10,
				// this value should be one.
				HeartbeatTick: 1,
			}
		}
		if m.Spec.Topology.AvailabilityGroup.LeaderElection.TransferLeadershipInterval == nil {
			m.Spec.Topology.AvailabilityGroup.LeaderElection.TransferLeadershipInterval = &meta.Duration{Duration: 1 * time.Second}
		}
		if m.Spec.Topology.AvailabilityGroup.LeaderElection.TransferLeadershipTimeout == nil {
			m.Spec.Topology.AvailabilityGroup.LeaderElection.TransferLeadershipTimeout = &meta.Duration{Duration: 60 * time.Second}
		}
	}

	if m.Spec.PodTemplate == nil {
		m.Spec.PodTemplate = &ofst.PodTemplateSpec{}
	}

	var mssqlVersion catalog.MSSQLServerVersion
	err := kc.Get(context.TODO(), types.NamespacedName{
		Name: m.Spec.Version,
	}, &mssqlVersion)
	if err != nil {
		klog.Errorf("can't get the MSSQLServer version object %s for %s \n", m.Spec.Version, err.Error())
		return
	}

	m.SetArbiterDefault()

	m.setDefaultContainerSecurityContext(&mssqlVersion, m.Spec.PodTemplate)

	m.SetTLSDefaults()

	m.SetHealthCheckerDefaults()

	m.setDefaultContainerResourceLimits(m.Spec.PodTemplate)

	if m.Spec.Monitor != nil {
		if m.Spec.Monitor.Prometheus == nil {
			m.Spec.Monitor.Prometheus = &mona.PrometheusSpec{}
		}
		if m.Spec.Monitor.Prometheus.Exporter.Port == 0 {
			m.Spec.Monitor.Prometheus.Exporter.Port = kubedb.MSSQLMonitoringDefaultServicePort
		}
		m.Spec.Monitor.SetDefaults()
		if m.Spec.Monitor.Prometheus.Exporter.SecurityContext.RunAsUser == nil {
			m.Spec.Monitor.Prometheus.Exporter.SecurityContext.RunAsUser = mssqlVersion.Spec.SecurityContext.RunAsUser
		}
		if m.Spec.Monitor.Prometheus.Exporter.SecurityContext.RunAsGroup == nil {
			m.Spec.Monitor.Prometheus.Exporter.SecurityContext.RunAsGroup = mssqlVersion.Spec.SecurityContext.RunAsUser
		}
	}

	if m.Spec.Init != nil && m.Spec.Init.Archiver != nil {
		if m.Spec.Init.Archiver.EncryptionSecret != nil && m.Spec.Init.Archiver.EncryptionSecret.Namespace == "" {
			m.Spec.Init.Archiver.EncryptionSecret.Namespace = m.GetNamespace()
		}
		if m.Spec.Init.Archiver.FullDBRepository != nil && m.Spec.Init.Archiver.FullDBRepository.Namespace == "" {
			m.Spec.Init.Archiver.FullDBRepository.Namespace = m.GetNamespace()
		}
		if m.Spec.Init.Archiver.ManifestRepository != nil && m.Spec.Init.Archiver.ManifestRepository.Namespace == "" {
			m.Spec.Init.Archiver.ManifestRepository.Namespace = m.GetNamespace()
		}
	}
}

func (m *MSSQLServer) setDefaultContainerSecurityContext(mssqlVersion *catalog.MSSQLServerVersion, podTemplate *ofst.PodTemplateSpec) {
	if podTemplate == nil {
		return
	}
	if podTemplate.Spec.SecurityContext == nil {
		podTemplate.Spec.SecurityContext = &core.PodSecurityContext{}
	}
	if podTemplate.Spec.SecurityContext.FSGroup == nil {
		podTemplate.Spec.SecurityContext.FSGroup = mssqlVersion.Spec.SecurityContext.RunAsUser
	}

	container := coreutil.GetContainerByName(podTemplate.Spec.Containers, kubedb.MSSQLContainerName)
	if container == nil {
		container = &core.Container{
			Name: kubedb.MSSQLContainerName,
		}
	}
	if container.SecurityContext == nil {
		container.SecurityContext = &core.SecurityContext{}
	}

	m.assignDefaultContainerSecurityContext(mssqlVersion, container.SecurityContext, true)

	podTemplate.Spec.Containers = coreutil.UpsertContainer(podTemplate.Spec.Containers, *container)

	initContainer := coreutil.GetContainerByName(podTemplate.Spec.InitContainers, kubedb.MSSQLInitContainerName)
	if initContainer == nil {
		initContainer = &core.Container{
			Name: kubedb.MSSQLInitContainerName,
		}
	}
	if initContainer.SecurityContext == nil {
		initContainer.SecurityContext = &core.SecurityContext{}
	}
	m.assignDefaultContainerSecurityContext(mssqlVersion, initContainer.SecurityContext, false)
	podTemplate.Spec.InitContainers = coreutil.UpsertContainer(podTemplate.Spec.InitContainers, *initContainer)

	if m.IsCluster() {
		coordinatorContainer := coreutil.GetContainerByName(podTemplate.Spec.Containers, kubedb.MSSQLCoordinatorContainerName)
		if coordinatorContainer == nil {
			coordinatorContainer = &core.Container{
				Name: kubedb.MSSQLCoordinatorContainerName,
			}
		}
		if coordinatorContainer.SecurityContext == nil {
			coordinatorContainer.SecurityContext = &core.SecurityContext{}
		}
		m.assignDefaultContainerSecurityContext(mssqlVersion, coordinatorContainer.SecurityContext, false)
		podTemplate.Spec.Containers = coreutil.UpsertContainer(podTemplate.Spec.Containers, *coordinatorContainer)
	}
}

func (m *MSSQLServer) assignDefaultContainerSecurityContext(mssqlVersion *catalog.MSSQLServerVersion, sc *core.SecurityContext, isMainContainer bool) {
	if sc.AllowPrivilegeEscalation == nil {
		sc.AllowPrivilegeEscalation = pointer.BoolP(false)
	}
	if sc.Capabilities == nil {
		if isMainContainer {
			sc.Capabilities = &core.Capabilities{
				Drop: []core.Capability{"ALL"},
				Add:  []core.Capability{"NET_BIND_SERVICE"},
			}
		} else {
			sc.Capabilities = &core.Capabilities{
				Drop: []core.Capability{"ALL"},
			}
		}
	}
	if sc.RunAsNonRoot == nil {
		sc.RunAsNonRoot = pointer.BoolP(true)
	}
	if sc.RunAsUser == nil {
		sc.RunAsUser = mssqlVersion.Spec.SecurityContext.RunAsUser
	}
	if sc.RunAsGroup == nil {
		sc.RunAsGroup = mssqlVersion.Spec.SecurityContext.RunAsUser
	}
	if sc.SeccompProfile == nil {
		sc.SeccompProfile = secomp.DefaultSeccompProfile()
	}
}

func (m *MSSQLServer) setDefaultContainerResourceLimits(podTemplate *ofst.PodTemplateSpec) {
	dbContainer := coreutil.GetContainerByName(podTemplate.Spec.Containers, kubedb.MSSQLContainerName)
	if dbContainer != nil && (dbContainer.Resources.Requests == nil && dbContainer.Resources.Limits == nil) {
		apis.SetDefaultResourceLimits(&dbContainer.Resources, kubedb.DefaultResourcesMemoryIntensiveMSSQLServer)
	}

	initContainer := coreutil.GetContainerByName(podTemplate.Spec.InitContainers, kubedb.MSSQLInitContainerName)
	if initContainer != nil && (initContainer.Resources.Requests == nil && initContainer.Resources.Limits == nil) {
		apis.SetDefaultResourceLimits(&initContainer.Resources, kubedb.DefaultInitContainerResource)
	}

	if m.IsAvailabilityGroup() {
		coordinatorContainer := coreutil.GetContainerByName(podTemplate.Spec.Containers, kubedb.MSSQLCoordinatorContainerName)
		if coordinatorContainer != nil && (coordinatorContainer.Resources.Requests == nil && coordinatorContainer.Resources.Limits == nil) {
			apis.SetDefaultResourceLimits(&coordinatorContainer.Resources, kubedb.CoordinatorDefaultResources)
		}
	}
}

func (m *MSSQLServer) SetArbiterDefault() {
	if m.IsAvailabilityGroup() && ptr.Deref(m.Spec.Replicas, 0)%2 == 0 && m.Spec.Arbiter == nil {
		m.Spec.Arbiter = &ArbiterSpec{
			Resources: core.ResourceRequirements{},
		}
		apis.SetDefaultResourceLimits(&m.Spec.Arbiter.Resources, kubedb.DefaultArbiter(false))
	}
}

func (m *MSSQLServer) SetTLSDefaults() {
	if m.Spec.TLS == nil || m.Spec.TLS.IssuerRef == nil {
		return
	}

	if m.Spec.TLS.ClientTLS == nil {
		defaultValue := false
		m.Spec.TLS.ClientTLS = &defaultValue
	}

	// Server-cert
	defaultServerOrg := []string{kubedb.KubeDBOrganization}
	defaultServerOrgUnit := []string{string(MSSQLServerServerCert)}
	_, cert := kmapi.GetCertificate(m.Spec.TLS.Certificates, string(MSSQLServerServerCert))
	if cert != nil && cert.Subject != nil {
		if cert.Subject.Organizations != nil {
			defaultServerOrg = cert.Subject.Organizations
		}
		if cert.Subject.OrganizationalUnits != nil {
			defaultServerOrgUnit = cert.Subject.OrganizationalUnits
		}
	}

	m.Spec.TLS.Certificates = kmapi.SetMissingSpecForCertificate(m.Spec.TLS.Certificates, kmapi.CertificateSpec{
		Alias:      string(MSSQLServerServerCert),
		SecretName: m.GetCertSecretName(MSSQLServerServerCert),
		Subject: &kmapi.X509Subject{
			Organizations:       defaultServerOrg,
			OrganizationalUnits: defaultServerOrgUnit,
		},
	})

	// Client-cert
	defaultClientOrg := []string{kubedb.KubeDBOrganization}
	defaultClientOrgUnit := []string{string(MSSQLServerClientCert)}
	_, cert = kmapi.GetCertificate(m.Spec.TLS.Certificates, string(MSSQLServerClientCert))
	if cert != nil && cert.Subject != nil {
		if cert.Subject.Organizations != nil {
			defaultClientOrg = cert.Subject.Organizations
		}
		if cert.Subject.OrganizationalUnits != nil {
			defaultClientOrgUnit = cert.Subject.OrganizationalUnits
		}
	}
	m.Spec.TLS.Certificates = kmapi.SetMissingSpecForCertificate(m.Spec.TLS.Certificates, kmapi.CertificateSpec{
		Alias:      string(MSSQLServerClientCert),
		SecretName: m.GetCertSecretName(MSSQLServerClientCert),
		Subject: &kmapi.X509Subject{
			Organizations:       defaultClientOrg,
			OrganizationalUnits: defaultClientOrgUnit,
		},
	})

	if m.IsAvailabilityGroup() {
		// Endpoint-cert
		defaultEndpointOrg := []string{kubedb.KubeDBOrganization}
		defaultEndpointOrgUnit := []string{string(MSSQLServerEndpointCert)}
		_, cert = kmapi.GetCertificate(m.Spec.TLS.Certificates, string(MSSQLServerEndpointCert))
		if cert != nil && cert.Subject != nil {
			if cert.Subject.Organizations != nil {
				defaultEndpointOrg = cert.Subject.Organizations
			}
			if cert.Subject.OrganizationalUnits != nil {
				defaultEndpointOrgUnit = cert.Subject.OrganizationalUnits
			}
		}

		m.Spec.TLS.Certificates = kmapi.SetMissingSpecForCertificate(m.Spec.TLS.Certificates, kmapi.CertificateSpec{
			Alias:      string(MSSQLServerEndpointCert),
			SecretName: m.GetCertSecretName(MSSQLServerEndpointCert),
			Subject: &kmapi.X509Subject{
				Organizations:       defaultEndpointOrg,
				OrganizationalUnits: defaultEndpointOrgUnit,
			},
		})
	}
}

func (m *MSSQLServer) ReplicasAreReady(lister pslister.PetSetLister) (bool, string, error) {
	// Desire number of petSets
	expectedItems := 1
	return checkReplicasOfPetSet(lister.PetSets(m.Namespace), labels.SelectorFromSet(m.OffshootLabels()), expectedItems)
}

// Map SecondaryAccessMode to SQL string values
func SecondaryAccessSQL(mode SecondaryAccessMode) string {
	switch mode {
	case SecondaryAccessModePassive:
		return "NO"
	case SecondaryAccessModeReadOnly:
		return "READ_ONLY"
	case SecondaryAccessModeAll:
		return "ALL"
	default:
		// Fallback to NO if unset or unknown
		return "NO"
	}
}
