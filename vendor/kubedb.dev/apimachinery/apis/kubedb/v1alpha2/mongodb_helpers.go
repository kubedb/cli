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
	"strings"

	"kubedb.dev/apimachinery/apis"
	"kubedb.dev/apimachinery/apis/catalog/v1alpha1"
	"kubedb.dev/apimachinery/apis/kubedb"
	"kubedb.dev/apimachinery/crds"

	"gomodules.xyz/pointer"
	"gomodules.xyz/version"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	appslister "k8s.io/client-go/listers/apps/v1"
	kmapi "kmodules.xyz/client-go/api/v1"
	"kmodules.xyz/client-go/apiextensions"
	core_util "kmodules.xyz/client-go/core/v1"
	v1 "kmodules.xyz/client-go/core/v1"
	meta_util "kmodules.xyz/client-go/meta"
	appcat "kmodules.xyz/custom-resources/apis/appcatalog/v1alpha1"
	mona "kmodules.xyz/monitoring-agent-api/api/v1"
	ofst "kmodules.xyz/offshoot-api/api/v1"
)

func (_ MongoDB) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralMongoDB))
}

var _ apis.ResourceInfo = &MongoDB{}

const (
	TLSCAKeyFileName    = "ca.key"
	TLSCACertFileName   = "ca.crt"
	MongoPemFileName    = "mongo.pem"
	MongoClientFileName = "client.pem"
	MongoCertDirectory  = "/var/run/mongodb/tls"

	MongoDBShardLabelKey  = "mongodb.kubedb.com/node.shard"
	MongoDBConfigLabelKey = "mongodb.kubedb.com/node.config"
	MongoDBMongosLabelKey = "mongodb.kubedb.com/node.mongos"

	MongoDBShardAffinityTemplateVar = "SHARD_INDEX"
)

func (m MongoDB) OffshootName() string {
	return m.Name
}

func (m MongoDB) ShardNodeName(nodeNum int32) string {
	if m.Spec.ShardTopology == nil {
		return ""
	}
	return fmt.Sprintf("%v%v", m.ShardCommonNodeName(), nodeNum)
}

func (m MongoDB) ShardNodeTemplate() string {
	if m.Spec.ShardTopology == nil {
		return ""
	}
	return fmt.Sprintf("%s${%s}", m.ShardCommonNodeName(), MongoDBShardAffinityTemplateVar)
}

func (m MongoDB) ShardCommonNodeName() string {
	if m.Spec.ShardTopology == nil {
		return ""
	}
	shardName := fmt.Sprintf("%v-shard", m.OffshootName())
	return m.Spec.ShardTopology.Shard.Prefix + shardName
}

func (m MongoDB) ConfigSvrNodeName() string {
	if m.Spec.ShardTopology == nil {
		return ""
	}
	configsvrName := fmt.Sprintf("%v-configsvr", m.OffshootName())
	return m.Spec.ShardTopology.ConfigServer.Prefix + configsvrName
}

func (m MongoDB) MongosNodeName() string {
	if m.Spec.ShardTopology == nil {
		return ""
	}
	mongosName := fmt.Sprintf("%v-mongos", m.OffshootName())
	return m.Spec.ShardTopology.Mongos.Prefix + mongosName
}

// RepSetName returns Replicaset name only for spec.replicaset
func (m MongoDB) RepSetName() string {
	if m.Spec.ReplicaSet == nil {
		return ""
	}
	return m.Spec.ReplicaSet.Name
}

func (m MongoDB) ShardRepSetName(nodeNum int32) string {
	repSetName := fmt.Sprintf("shard%v", nodeNum)
	if m.Spec.ShardTopology != nil && m.Spec.ShardTopology.Shard.Prefix != "" {
		repSetName = fmt.Sprintf("%v%v", m.Spec.ShardTopology.Shard.Prefix, nodeNum)
	}
	return repSetName
}

func (m MongoDB) ConfigSvrRepSetName() string {
	repSetName := "cnfRepSet"
	if m.Spec.ShardTopology != nil && m.Spec.ShardTopology.ConfigServer.Prefix != "" {
		repSetName = m.Spec.ShardTopology.ConfigServer.Prefix
	}
	return repSetName
}

func (m MongoDB) OffshootSelectors() map[string]string {
	return map[string]string{
		LabelDatabaseName: m.Name,
		LabelDatabaseKind: ResourceKindMongoDB,
	}
}

func (m MongoDB) ShardSelectors(nodeNum int32) map[string]string {
	return v1.UpsertMap(m.OffshootSelectors(), map[string]string{
		MongoDBShardLabelKey: m.ShardNodeName(nodeNum),
	})
}

func (m MongoDB) ConfigSvrSelectors() map[string]string {
	return v1.UpsertMap(m.OffshootSelectors(), map[string]string{
		MongoDBConfigLabelKey: m.ConfigSvrNodeName(),
	})
}

func (m MongoDB) MongosSelectors() map[string]string {
	return v1.UpsertMap(m.OffshootSelectors(), map[string]string{
		MongoDBMongosLabelKey: m.MongosNodeName(),
	})
}

func (m MongoDB) OffshootLabels() map[string]string {
	out := m.OffshootSelectors()
	out[meta_util.NameLabelKey] = ResourceSingularMongoDB
	out[meta_util.VersionLabelKey] = string(m.Spec.Version)
	out[meta_util.InstanceLabelKey] = m.Name
	out[meta_util.ComponentLabelKey] = ComponentDatabase
	out[meta_util.ManagedByLabelKey] = kubedb.GroupName
	return meta_util.FilterKeys(kubedb.GroupName, out, m.Labels)
}

func (m MongoDB) ShardLabels(nodeNum int32) map[string]string {
	return meta_util.FilterKeys(kubedb.GroupName, m.OffshootLabels(), m.ShardSelectors(nodeNum))
}

func (m MongoDB) ConfigSvrLabels() map[string]string {
	return meta_util.FilterKeys(kubedb.GroupName, m.OffshootLabels(), m.ConfigSvrSelectors())
}

func (m MongoDB) MongosLabels() map[string]string {
	return meta_util.FilterKeys(kubedb.GroupName, m.OffshootLabels(), m.MongosSelectors())
}

func (m MongoDB) ResourceShortCode() string {
	return ResourceCodeMongoDB
}

func (m MongoDB) ResourceKind() string {
	return ResourceKindMongoDB
}

func (m MongoDB) ResourceSingular() string {
	return ResourceSingularMongoDB
}

func (m MongoDB) ResourcePlural() string {
	return ResourcePluralMongoDB
}

func (m MongoDB) ServiceName() string {
	return m.OffshootName()
}

// Governing Service Name. Here, name parameter is either
// OffshootName, ShardNodeName or ConfigSvrNodeName
func (m MongoDB) GoverningServiceName(name string) string {
	if name == "" {
		panic(fmt.Sprintf("StatefulSet name is missing for MongoDB %s/%s", m.Namespace, m.Name))
	}
	return name + "-pods"
}

// HostAddress returns serviceName for standalone mongodb.
// and, for replica set = <replName>/<host1>,<host2>,<host3>
// Here, host1 = <pod-name>.<governing-serviceName>
// Governing service name is used for replica host because,
// we used governing service name as part of host while adding members
// to replicaset.
func (m MongoDB) HostAddress() string {
	if m.Spec.ReplicaSet != nil {
		return fmt.Sprintf("%v/", m.RepSetName()) + strings.Join(m.Hosts(), ",")
	}

	return m.ServiceName()
}

func (m MongoDB) Hosts() []string {
	hosts := []string{fmt.Sprintf("%v-0.%v.%v.svc", m.Name, m.GoverningServiceName(m.OffshootName()), m.Namespace)}
	if m.Spec.ReplicaSet != nil {
		hosts = make([]string, *m.Spec.Replicas)
		for i := 0; i < int(pointer.Int32(m.Spec.Replicas)); i++ {
			hosts[i] = fmt.Sprintf("%v-%d.%v.%v.svc", m.Name, i, m.GoverningServiceName(m.OffshootName()), m.Namespace)
		}
	}
	return hosts
}

// ShardDSN = <shardReplName>/<host1:port>,<host2:port>,<host3:port>
//// Here, host1 = <pod-name>.<governing-serviceName>.svc
func (m MongoDB) ShardDSN(nodeNum int32) string {
	if m.Spec.ShardTopology == nil {
		return ""
	}
	return fmt.Sprintf("%v/", m.ShardRepSetName(nodeNum)) + strings.Join(m.ShardHosts(nodeNum), ",")
}

func (m MongoDB) ShardHosts(nodeNum int32) []string {
	if m.Spec.ShardTopology == nil {
		return []string{}
	}
	hosts := make([]string, m.Spec.ShardTopology.Shard.Replicas)
	for i := 0; i < int(m.Spec.ShardTopology.Shard.Replicas); i++ {
		hosts[i] = fmt.Sprintf("%v-%d.%v.%v.svc:%v", m.ShardNodeName(nodeNum), i, m.GoverningServiceName(m.ShardNodeName(nodeNum)), m.Namespace, MongoDBDatabasePort)
	}
	return hosts
}

// ConfigSvrDSN = <configSvrReplName>/<host1:port>,<host2:port>,<host3:port>
//// Here, host1 = <pod-name>.<governing-serviceName>.svc
func (m MongoDB) ConfigSvrDSN() string {
	if m.Spec.ShardTopology == nil {
		return ""
	}

	return fmt.Sprintf("%v/", m.ConfigSvrRepSetName()) + strings.Join(m.ConfigSvrHosts(), ",")
}

func (m MongoDB) ConfigSvrHosts() []string {
	if m.Spec.ShardTopology == nil {
		return []string{}
	}

	hosts := make([]string, m.Spec.ShardTopology.ConfigServer.Replicas)
	for i := 0; i < int(m.Spec.ShardTopology.ConfigServer.Replicas); i++ {
		hosts[i] = fmt.Sprintf("%v-%d.%v.%v.svc:%v", m.ConfigSvrNodeName(), i, m.GoverningServiceName(m.ConfigSvrNodeName()), m.Namespace, MongoDBDatabasePort)
	}
	return hosts
}

func (m MongoDB) MongosHosts() []string {
	if m.Spec.ShardTopology == nil {
		return []string{}
	}

	hosts := make([]string, m.Spec.ShardTopology.Mongos.Replicas)
	for i := 0; i < int(m.Spec.ShardTopology.Mongos.Replicas); i++ {
		hosts[i] = fmt.Sprintf("%v-%d.%v.%v.svc:%v", m.MongosNodeName(), i, m.GoverningServiceName(m.MongosNodeName()), m.Namespace, MongoDBDatabasePort)
	}
	return hosts
}

type mongoDBApp struct {
	*MongoDB
}

func (r mongoDBApp) Name() string {
	return r.MongoDB.Name
}

func (r mongoDBApp) Type() appcat.AppType {
	return appcat.AppType(fmt.Sprintf("%s/%s", kubedb.GroupName, ResourceSingularMongoDB))
}

func (m MongoDB) AppBindingMeta() appcat.AppBindingMeta {
	return &mongoDBApp{&m}
}

type mongoDBStatsService struct {
	*MongoDB
}

func (m mongoDBStatsService) GetNamespace() string {
	return m.MongoDB.GetNamespace()
}

func (m mongoDBStatsService) ServiceName() string {
	return m.OffshootName() + "-stats"
}

func (m mongoDBStatsService) ServiceMonitorName() string {
	return m.ServiceName()
}

func (m mongoDBStatsService) ServiceMonitorAdditionalLabels() map[string]string {
	return m.OffshootLabels()
}

func (m mongoDBStatsService) Path() string {
	return DefaultStatsPath
}

func (m mongoDBStatsService) Scheme() string {
	return ""
}

func (m MongoDB) StatsService() mona.StatsAccessor {
	return &mongoDBStatsService{&m}
}

func (m MongoDB) StatsServiceLabels() map[string]string {
	lbl := meta_util.FilterKeys(kubedb.GroupName, m.OffshootSelectors(), m.Labels)
	lbl[LabelRole] = RoleStats
	return lbl
}

func (m *MongoDB) SetDefaults(mgVersion *v1alpha1.MongoDBVersion, topology *core_util.Topology) {
	if m == nil {
		return
	}

	if m.Spec.StorageType == "" {
		m.Spec.StorageType = StorageTypeDurable
	}
	if m.Spec.StorageEngine == "" {
		m.Spec.StorageEngine = StorageEngineWiredTiger
	}
	if m.Spec.TerminationPolicy == "" {
		m.Spec.TerminationPolicy = TerminationPolicyDelete
	}

	if m.Spec.SSLMode == "" {
		m.Spec.SSLMode = SSLModeDisabled
	}

	if (m.Spec.ReplicaSet != nil || m.Spec.ShardTopology != nil) && m.Spec.ClusterAuthMode == "" {
		if m.Spec.SSLMode == SSLModeDisabled || m.Spec.SSLMode == SSLModeAllowSSL {
			m.Spec.ClusterAuthMode = ClusterAuthModeKeyFile
		} else {
			m.Spec.ClusterAuthMode = ClusterAuthModeX509
		}
	}

	if m.Spec.ShardTopology != nil {
		if m.Spec.ShardTopology.Mongos.PodTemplate.Spec.Lifecycle == nil {
			m.Spec.ShardTopology.Mongos.PodTemplate.Spec.Lifecycle = new(core.Lifecycle)
		}

		m.Spec.ShardTopology.Mongos.PodTemplate.Spec.Lifecycle.PreStop = &core.Handler{
			Exec: &core.ExecAction{
				Command: []string{
					"bash",
					"-c",
					"mongo admin --username=$MONGO_INITDB_ROOT_USERNAME --password=$MONGO_INITDB_ROOT_PASSWORD --quiet --eval \"db.adminCommand({ shutdown: 1 })\" || true",
				},
			},
		}

		if m.Spec.ShardTopology.ConfigServer.PodTemplate.Spec.ServiceAccountName == "" {
			m.Spec.ShardTopology.ConfigServer.PodTemplate.Spec.ServiceAccountName = m.OffshootName()
		}
		if m.Spec.ShardTopology.Mongos.PodTemplate.Spec.ServiceAccountName == "" {
			m.Spec.ShardTopology.Mongos.PodTemplate.Spec.ServiceAccountName = m.OffshootName()
		}
		if m.Spec.ShardTopology.Shard.PodTemplate.Spec.ServiceAccountName == "" {
			m.Spec.ShardTopology.Shard.PodTemplate.Spec.ServiceAccountName = m.OffshootName()
		}

		// set default probes
		m.setDefaultProbes(&m.Spec.ShardTopology.Shard.PodTemplate, mgVersion)
		m.setDefaultProbes(&m.Spec.ShardTopology.ConfigServer.PodTemplate, mgVersion)
		m.setDefaultProbes(&m.Spec.ShardTopology.Mongos.PodTemplate, mgVersion)

		// set default affinity (PodAntiAffinity)
		shardLabels := m.OffshootSelectors()
		shardLabels[MongoDBShardLabelKey] = m.ShardNodeTemplate()
		m.setDefaultAffinity(&m.Spec.ShardTopology.Shard.PodTemplate, shardLabels, topology)

		configServerLabels := m.OffshootSelectors()
		configServerLabels[MongoDBConfigLabelKey] = m.ConfigSvrNodeName()
		m.setDefaultAffinity(&m.Spec.ShardTopology.ConfigServer.PodTemplate, configServerLabels, topology)

		mongosLabels := m.OffshootSelectors()
		mongosLabels[MongoDBMongosLabelKey] = m.MongosNodeName()
		m.setDefaultAffinity(&m.Spec.ShardTopology.Mongos.PodTemplate, mongosLabels, topology)
	} else {
		if m.Spec.Replicas == nil {
			m.Spec.Replicas = pointer.Int32P(1)
		}

		if m.Spec.PodTemplate == nil {
			m.Spec.PodTemplate = new(ofst.PodTemplateSpec)
		}
		if m.Spec.PodTemplate.Spec.ServiceAccountName == "" {
			m.Spec.PodTemplate.Spec.ServiceAccountName = m.OffshootName()
		}

		// set default probes
		m.setDefaultProbes(m.Spec.PodTemplate, mgVersion)
		// set default affinity (PodAntiAffinity)
		m.setDefaultAffinity(m.Spec.PodTemplate, m.OffshootSelectors(), topology)
	}

	m.SetTLSDefaults()
}

func (m *MongoDB) SetTLSDefaults() {
	if m.Spec.TLS == nil || m.Spec.TLS.IssuerRef == nil {
		return
	}

	if m.Spec.ShardTopology != nil {
		m.Spec.TLS.Certificates = kmapi.SetMissingSpecForCertificate(m.Spec.TLS.Certificates, kmapi.CertificateSpec{
			Alias:      string(MongoDBServerCert),
			SecretName: "",
			Subject: &kmapi.X509Subject{
				Organizations:       []string{KubeDBOrganization},
				OrganizationalUnits: []string{string(MongoDBServerCert)},
			},
		})
		// reset secret name to empty string, since multiple secrets will be created for each StatefulSet.
		m.Spec.TLS.Certificates = kmapi.SetSecretNameForCertificate(m.Spec.TLS.Certificates, string(MongoDBServerCert), "")
	} else {
		m.Spec.TLS.Certificates = kmapi.SetMissingSpecForCertificate(m.Spec.TLS.Certificates, kmapi.CertificateSpec{
			Alias:      string(MongoDBServerCert),
			SecretName: m.CertificateName(MongoDBServerCert, ""),
			Subject: &kmapi.X509Subject{
				Organizations:       []string{KubeDBOrganization},
				OrganizationalUnits: []string{string(MongoDBServerCert)},
			},
		})
	}
	m.Spec.TLS.Certificates = kmapi.SetMissingSpecForCertificate(m.Spec.TLS.Certificates, kmapi.CertificateSpec{
		Alias:      string(MongoDBClientCert),
		SecretName: m.CertificateName(MongoDBClientCert, ""),
		Subject: &kmapi.X509Subject{
			Organizations:       []string{KubeDBOrganization},
			OrganizationalUnits: []string{string(MongoDBClientCert)},
		},
	})
	m.Spec.TLS.Certificates = kmapi.SetMissingSpecForCertificate(m.Spec.TLS.Certificates, kmapi.CertificateSpec{
		Alias:      string(MongoDBMetricsExporterCert),
		SecretName: m.CertificateName(MongoDBMetricsExporterCert, ""),
		Subject: &kmapi.X509Subject{
			Organizations:       []string{KubeDBOrganization},
			OrganizationalUnits: []string{string(MongoDBMetricsExporterCert)},
		},
	})
}
func (m *MongoDB) getCmdForProbes(mgVersion *v1alpha1.MongoDBVersion) []string {
	var sslArgs string
	if m.Spec.SSLMode == SSLModeRequireSSL {
		sslArgs = fmt.Sprintf("--tls --tlsCAFile=%v/%v --tlsCertificateKeyFile=%v/%v",
			MongoCertDirectory, TLSCACertFileName, MongoCertDirectory, MongoClientFileName)

		breakingVer, _ := version.NewVersion("4.1")
		exceptionVer, _ := version.NewVersion("4.1.4")
		currentVer, err := version.NewVersion(mgVersion.Spec.Version)
		if err != nil {
			panic(fmt.Errorf("MongoDB %s/%s: unable to parse version. reason: %s", m.Namespace, m.Name, err.Error()))
		}
		if currentVer.Equal(exceptionVer) {
			sslArgs = fmt.Sprintf("--tls --tlsCAFile=%v/%v --tlsPEMKeyFile=%v/%v", MongoCertDirectory, TLSCACertFileName, MongoCertDirectory, MongoClientFileName)
		} else if currentVer.LessThan(breakingVer) {
			sslArgs = fmt.Sprintf("--ssl --sslCAFile=%v/%v --sslPEMKeyFile=%v/%v", MongoCertDirectory, TLSCACertFileName, MongoCertDirectory, MongoClientFileName)
		}
	}

	return []string{
		"bash",
		"-c",
		fmt.Sprintf(`set -x; if [[ $(mongo admin --host=localhost %v --username=$MONGO_INITDB_ROOT_USERNAME --password=$MONGO_INITDB_ROOT_PASSWORD --authenticationDatabase=admin --quiet --eval "db.adminCommand('ping').ok" ) -eq "1" ]]; then 
          exit 0
        fi
        exit 1`, sslArgs),
	}
}

func (m *MongoDB) GetDefaultLivenessProbeSpec(mgVersion *v1alpha1.MongoDBVersion) *core.Probe {
	return &core.Probe{
		Handler: core.Handler{
			Exec: &core.ExecAction{
				Command: m.getCmdForProbes(mgVersion),
			},
		},
		FailureThreshold: 3,
		PeriodSeconds:    10,
		SuccessThreshold: 1,
		TimeoutSeconds:   5,
	}
}

func (m *MongoDB) GetDefaultReadinessProbeSpec(mgVersion *v1alpha1.MongoDBVersion) *core.Probe {
	return &core.Probe{
		Handler: core.Handler{
			Exec: &core.ExecAction{
				Command: m.getCmdForProbes(mgVersion),
			},
		},
		FailureThreshold: 3,
		PeriodSeconds:    10,
		SuccessThreshold: 1,
		TimeoutSeconds:   1,
	}
}

// setDefaultProbes sets defaults only when probe fields are nil.
// In operator, check if the value of probe fields is "{}".
// For "{}", ignore readinessprobe or livenessprobe in statefulset.
// ref: https://github.com/helm/charts/blob/345ba987722350ffde56ec34d2928c0b383940aa/stable/mongodb/templates/deployment-standalone.yaml#L93
func (m *MongoDB) setDefaultProbes(podTemplate *ofst.PodTemplateSpec, mgVersion *v1alpha1.MongoDBVersion) {
	if podTemplate == nil {
		return
	}

	if podTemplate.Spec.LivenessProbe == nil {
		podTemplate.Spec.LivenessProbe = m.GetDefaultLivenessProbeSpec(mgVersion)
	}
	if podTemplate.Spec.ReadinessProbe == nil {
		podTemplate.Spec.ReadinessProbe = m.GetDefaultReadinessProbeSpec(mgVersion)
	}
}

// setDefaultAffinity
func (m *MongoDB) setDefaultAffinity(podTemplate *ofst.PodTemplateSpec, labels map[string]string, topology *core_util.Topology) {
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

// setSecurityContext will set default PodSecurityContext.
// These values will be applied only to newly created objects.
// These defaultings should not be applied to DBs or dormantDBs,
// that is managed by previous operators,
func (m *MongoDBSpec) SetSecurityContext(podTemplate *ofst.PodTemplateSpec) {
	if podTemplate == nil {
		return
	}
	if podTemplate.Spec.SecurityContext == nil {
		podTemplate.Spec.SecurityContext = new(core.PodSecurityContext)
	}
	if podTemplate.Spec.SecurityContext.FSGroup == nil {
		podTemplate.Spec.SecurityContext.FSGroup = pointer.Int64P(999)
	}
	if podTemplate.Spec.SecurityContext.RunAsNonRoot == nil {
		podTemplate.Spec.SecurityContext.RunAsNonRoot = pointer.BoolP(true)
	}
	if podTemplate.Spec.SecurityContext.RunAsUser == nil {
		podTemplate.Spec.SecurityContext.RunAsUser = pointer.Int64P(999)
	}
}

func (m *MongoDBSpec) GetPersistentSecrets() []string {
	if m == nil {
		return nil
	}

	var secrets []string
	if m.AuthSecret != nil {
		secrets = append(secrets, m.AuthSecret.Name)
	}
	if m.KeyFileSecret != nil {
		secrets = append(secrets, m.KeyFileSecret.Name)
	}
	return secrets
}

func (m *MongoDB) KeyFileRequired() bool {
	if m == nil {
		return false
	}
	return m.Spec.ClusterAuthMode == ClusterAuthModeKeyFile ||
		m.Spec.ClusterAuthMode == ClusterAuthModeSendKeyFile ||
		m.Spec.ClusterAuthMode == ClusterAuthModeSendX509
}

// CertificateName returns the default certificate name and/or certificate secret name for a certificate alias
func (m *MongoDB) CertificateName(alias MongoDBCertificateAlias, stsName string) string {
	if m.Spec.ShardTopology != nil && alias == MongoDBServerCert {
		if stsName == "" {
			panic(fmt.Sprintf("StatefulSet name required to compute %s certificate name for MongoDB %s/%s", alias, m.Namespace, m.Name))
		}
		return meta_util.NameWithSuffix(stsName, fmt.Sprintf("%s-cert", string(alias)))
	}
	return meta_util.NameWithSuffix(m.Name, fmt.Sprintf("%s-cert", string(alias)))
}

// MustCertSecretName returns the secret name for a certificate alias
func (m *MongoDB) MustCertSecretName(alias MongoDBCertificateAlias, stsName string) string {
	if m == nil {
		panic("missing MongoDB database")
	} else if m.Spec.TLS == nil {
		panic(fmt.Errorf("MongoDB %s/%s is missing tls spec", m.Namespace, m.Name))
	}
	if m.Spec.ShardTopology != nil && alias == MongoDBServerCert {
		if stsName == "" {
			panic(fmt.Sprintf("StatefulSet name required to compute %s certificate name for MongoDB %s/%s", alias, m.Namespace, m.Name))
		}
		return m.CertificateName(alias, stsName)
	}
	name, ok := kmapi.GetCertificateSecretName(m.Spec.TLS.Certificates, string(alias))
	if !ok {
		panic(fmt.Errorf("MongoDB %s/%s is missing secret name for %s certificate", m.Namespace, m.Name, alias))
	}
	return name
}

func (m *MongoDB) ReplicasAreReady(lister appslister.StatefulSetLister) (bool, string, error) {
	// Desire number of statefulSets
	expectedItems := 1
	if m.Spec.ShardTopology != nil {
		expectedItems = 2 + int(m.Spec.ShardTopology.Shard.Shards)
	}
	return checkReplicas(lister.StatefulSets(m.Namespace), labels.SelectorFromSet(m.OffshootLabels()), expectedItems)
}
