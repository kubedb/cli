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
	"fmt"

	"kubedb.dev/apimachinery/apis"
	"kubedb.dev/apimachinery/apis/kubedb"
	"kubedb.dev/apimachinery/crds"

	promapi "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"gomodules.xyz/pointer"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	appslister "k8s.io/client-go/listers/apps/v1"
	kmapi "kmodules.xyz/client-go/api/v1"
	"kmodules.xyz/client-go/apiextensions"
	core_util "kmodules.xyz/client-go/core/v1"
	meta_util "kmodules.xyz/client-go/meta"
	appcat "kmodules.xyz/custom-resources/apis/appcatalog/v1alpha1"
	mona "kmodules.xyz/monitoring-agent-api/api/v1"
	ofst "kmodules.xyz/offshoot-api/api/v1"
)

func (_ MySQL) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralMySQL))
}

var _ apis.ResourceInfo = &MySQL{}

func (m MySQL) OffshootName() string {
	return m.Name
}

func (m MySQL) OffshootSelectors() map[string]string {
	return m.offshootSelectors(MySQLComponentDB)
}

func (m MySQL) RouterOffshootSelectors() map[string]string {
	return m.offshootSelectors(MySQLComponentRouter)
}

func (m MySQL) offshootSelectors(component string) map[string]string {
	selectors := map[string]string{
		meta_util.NameLabelKey:      m.ResourceFQN(),
		meta_util.InstanceLabelKey:  m.Name,
		meta_util.ManagedByLabelKey: kubedb.GroupName,
	}
	if m.IsInnoDBCluster() {
		selectors[MySQLComponentKey] = component
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
	selector[meta_util.ComponentLabelKey] = ComponentDatabase
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
	return DefaultStatsPath
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
	return m.ServiceLabels(StatsServiceAlias, map[string]string{LabelRole: RoleStats})
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

func (m *MySQL) IsReadReplica() bool {
	return m.Spec.Topology != nil &&
		m.Spec.Topology.Mode != nil &&
		*m.Spec.Topology.Mode == MySQLModeReadReplica
}

func (m *MySQL) IsSemiSync() bool {
	return m.Spec.Topology != nil &&
		m.Spec.Topology.Mode != nil &&
		*m.Spec.Topology.Mode == MySQLModeSemiSync
}

func (m *MySQL) SetDefaults(topology *core_util.Topology) {
	if m == nil {
		return
	}
	if m.Spec.StorageType == "" {
		m.Spec.StorageType = StorageTypeDurable
	}
	if m.Spec.TerminationPolicy == "" {
		m.Spec.TerminationPolicy = TerminationPolicyDelete
	}

	if m.UsesGroupReplication() {
		if m.Spec.Replicas == nil {
			m.Spec.Replicas = pointer.Int32P(MySQLDefaultGroupSize)
		}
	} else {
		if m.Spec.Replicas == nil {
			m.Spec.Replicas = pointer.Int32P(1)
		}
	}

	if m.Spec.PodTemplate.Spec.ServiceAccountName == "" {
		m.Spec.PodTemplate.Spec.ServiceAccountName = m.OffshootName()
	}

	m.Spec.Monitor.SetDefaults()
	m.setDefaultAffinity(&m.Spec.PodTemplate, m.OffshootSelectors(), topology)
	m.SetTLSDefaults()
	m.SetHealthCheckerDefaults()
	apis.SetDefaultResourceLimits(&m.Spec.PodTemplate.Spec.Resources, DefaultResources)
}

// setDefaultAffinity
func (m *MySQL) setDefaultAffinity(podTemplate *ofst.PodTemplateSpec, labels map[string]string, topology *core_util.Topology) {
	if podTemplate == nil {
		return
	} else if podTemplate.Spec.Affinity != nil {
		// Update topologyKey fields according to Kubernetes version
		topology.ConvertAffinity(podTemplate.Spec.Affinity)
		return
	}

	podTemplate.Spec.Affinity = &core.Affinity{
		PodAntiAffinity: &core.PodAntiAffinity{
			PreferredDuringSchedulingIgnoredDuringExecution: []core.WeightedPodAffinityTerm{
				// Prefer to not schedule multiple pods on the same node
				{
					Weight: 100,
					PodAffinityTerm: core.PodAffinityTerm{
						Namespaces: []string{m.Namespace},
						LabelSelector: &metav1.LabelSelector{
							MatchLabels: labels,
						},
						TopologyKey: core.LabelHostname,
					},
				},
				// Prefer to not schedule multiple pods on the node with same zone
				{
					Weight: 50,
					PodAffinityTerm: core.PodAffinityTerm{
						Namespaces: []string{m.Namespace},
						LabelSelector: &metav1.LabelSelector{
							MatchLabels: labels,
						},
						TopologyKey: topology.LabelZone,
					},
				},
			},
		},
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

func (m *MySQL) ReplicasAreReady(lister appslister.StatefulSetLister) (bool, string, error) {
	// Desire number of statefulSets
	expectedItems := 1
	return checkReplicas(lister.StatefulSets(m.Namespace), labels.SelectorFromSet(m.OffshootLabels()), expectedItems)
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
