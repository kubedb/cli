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

package v1alpha1

import (
	"fmt"
	"strconv"

	"kubedb.dev/apimachinery/apis"
	"kubedb.dev/apimachinery/apis/catalog/v1alpha1"
	"kubedb.dev/apimachinery/apis/kubedb"
	"kubedb.dev/apimachinery/crds"

	"github.com/appscode/go/types"
	"gomodules.xyz/version"
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	out[meta_util.ManagedByLabelKey] = GenericKey
	return meta_util.FilterKeys(GenericKey, out, m.Labels)
}

func (m MongoDB) ShardLabels(nodeNum int32) map[string]string {
	return meta_util.FilterKeys(GenericKey, m.OffshootLabels(), m.ShardSelectors(nodeNum))
}

func (m MongoDB) ConfigSvrLabels() map[string]string {
	return meta_util.FilterKeys(GenericKey, m.OffshootLabels(), m.ConfigSvrSelectors())
}

func (m MongoDB) MongosLabels() map[string]string {
	return meta_util.FilterKeys(GenericKey, m.OffshootLabels(), m.MongosSelectors())
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
func (m MongoDB) GvrSvcName(name string) string {
	return name + "-gvr"
}

// HostAddress returns serviceName for standalone mongodb.
// and, for replica set = <replName>/<host1>,<host2>,<host3>
// Here, host1 = <pod-name>.<governing-serviceName>
// Governing service name is used for replica host because,
// we used governing service name as part of host while adding members
// to replicaset.
func (m MongoDB) HostAddress() string {
	host := m.ServiceName()
	if m.Spec.ReplicaSet != nil {
		//host = m.Spec.ReplicaSet.Name + "/" + m.Name + "-0." + m.GvrSvcName(m.OffshootName()) + "." + m.Namespace + ".svc"
		host = fmt.Sprintf("%v/", m.RepSetName())
		for i := 0; i < int(types.Int32(m.Spec.Replicas)); i++ {
			if i != 0 {
				host += ","
			}
			host += fmt.Sprintf("%v-%v.%v.%v.svc", m.Name, strconv.Itoa(i), m.GvrSvcName(m.OffshootName()), m.Namespace)
		}
	}
	return host
}

// ShardDSN = <shardReplName>/<host1:port>,<host2:port>,<host3:port>
//// Here, host1 = <pod-name>.<governing-serviceName>.svc
func (m MongoDB) ShardDSN(nodeNum int32) string {
	if m.Spec.ShardTopology == nil {
		return ""
	}
	host := fmt.Sprintf("%v/", m.ShardRepSetName(nodeNum))
	for i := 0; i < int(m.Spec.ShardTopology.Shard.Replicas); i++ {
		//host += "," + m.ShardNodeName(nodeNum) + "-" + strconv.Itoa(i) + "." + m.GvrSvcName(m.ShardNodeName(nodeNum)) + "." + m.Namespace + ".svc"

		if i != 0 {
			host += ","
		}
		host += fmt.Sprintf("%v-%v.%v.%v.svc:%v", m.ShardNodeName(nodeNum), strconv.Itoa(i), m.GvrSvcName(m.ShardNodeName(nodeNum)), m.Namespace, MongoDBShardPort)
	}
	return host
}

// ConfigSvrDSN = <configSvrReplName>/<host1:port>,<host2:port>,<host3:port>
//// Here, host1 = <pod-name>.<governing-serviceName>.svc
func (m MongoDB) ConfigSvrDSN() string {
	if m.Spec.ShardTopology == nil {
		return ""
	}
	//	host := m.ConfigSvrRepSetName() + "/" + m.ConfigSvrNodeName() + "-0." + m.GvrSvcName(m.ConfigSvrNodeName()) + "." + m.Namespace + ".svc"
	host := fmt.Sprintf("%v/", m.ConfigSvrRepSetName())
	for i := 0; i < int(m.Spec.ShardTopology.ConfigServer.Replicas); i++ {
		if i != 0 {
			host += ","
		}
		host += fmt.Sprintf("%v-%v.%v.%v.svc:%v", m.ConfigSvrNodeName(), strconv.Itoa(i), m.GvrSvcName(m.ConfigSvrNodeName()), m.Namespace, MongoDBShardPort)
	}
	return host
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
	return fmt.Sprintf("kubedb-%s-%s", m.Namespace, m.Name)
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
	lbl := meta_util.FilterKeys(GenericKey, m.OffshootSelectors(), m.Labels)
	lbl[LabelRole] = RoleStats
	return lbl
}

func (m *MongoDB) GetMonitoringVendor() string {
	if m.Spec.Monitor != nil {
		return m.Spec.Monitor.Agent.Vendor()
	}
	return ""
}

func (m *MongoDB) SetDefaults(mgVersion *v1alpha1.MongoDBVersion, topology *core_util.Topology) {
	if m == nil {
		return
	}

	if m.Spec.StorageType == "" {
		m.Spec.StorageType = StorageTypeDurable
	}
	if m.Spec.UpdateStrategy.Type == "" {
		m.Spec.UpdateStrategy.Type = apps.RollingUpdateStatefulSetStrategyType
	}
	if m.Spec.TerminationPolicy == "" {
		m.Spec.TerminationPolicy = TerminationPolicyDelete
	} else if m.Spec.TerminationPolicy == TerminationPolicyPause {
		m.Spec.TerminationPolicy = TerminationPolicyHalt
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

		if m.Spec.ShardTopology.Mongos.Strategy.Type == "" {
			m.Spec.ShardTopology.Mongos.Strategy.Type = apps.RollingUpdateDeploymentStrategyType
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
			m.Spec.Replicas = types.Int32P(1)
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
}

// setDefaultProbes sets defaults only when probe fields are nil.
// In operator, check if the value of probe fields is "{}".
// For "{}", ignore readinessprobe or livenessprobe in statefulset.
// ref: https://github.com/helm/charts/blob/345ba987722350ffde56ec34d2928c0b383940aa/stable/mongodb/templates/deployment-standalone.yaml#L93
func (m *MongoDB) setDefaultProbes(podTemplate *ofst.PodTemplateSpec, mgVersion *v1alpha1.MongoDBVersion) {
	if podTemplate == nil {
		return
	}
	var sslArgs string
	if m.Spec.SSLMode == SSLModeRequireSSL {
		sslArgs = fmt.Sprintf("--tls --tlsCAFile=%v/%v --tlsCertificateKeyFile=%v/%v",
			MongoCertDirectory, TLSCACertFileName, MongoCertDirectory, MongoClientFileName)

		breakingVer, err := version.NewVersion("4.1")
		if err != nil {
			return
		}
		exceptionVer, err := version.NewVersion("4.1.4")
		if err != nil {
			return
		}
		currentVer, err := version.NewVersion(mgVersion.Spec.Version)
		if err != nil {
			return
		}
		if currentVer.Equal(exceptionVer) {
			sslArgs = fmt.Sprintf("--tls --tlsCAFile=%v/%v --tlsPEMKeyFile=%v/%v", MongoCertDirectory, TLSCACertFileName, MongoCertDirectory, MongoClientFileName)
		} else if currentVer.LessThan(breakingVer) {
			sslArgs = fmt.Sprintf("--ssl --sslCAFile=%v/%v --sslPEMKeyFile=%v/%v", MongoCertDirectory, TLSCACertFileName, MongoCertDirectory, MongoClientFileName)
		}
	}

	cmd := []string{
		"bash",
		"-c",
		fmt.Sprintf(`set -x; if [[ $(mongo admin --host=localhost %v --username=$MONGO_INITDB_ROOT_USERNAME --password=$MONGO_INITDB_ROOT_PASSWORD --authenticationDatabase=admin --quiet --eval "db.adminCommand('ping').ok" ) -eq "1" ]]; then 
          exit 0
        fi
        exit 1`, sslArgs),
	}

	if podTemplate.Spec.LivenessProbe == nil {
		podTemplate.Spec.LivenessProbe = &core.Probe{
			Handler: core.Handler{
				Exec: &core.ExecAction{
					Command: cmd,
				},
			},
			FailureThreshold: 3,
			PeriodSeconds:    10,
			SuccessThreshold: 1,
			TimeoutSeconds:   5,
		}
	}
	if podTemplate.Spec.ReadinessProbe == nil {
		podTemplate.Spec.ReadinessProbe = &core.Probe{
			Handler: core.Handler{
				Exec: &core.ExecAction{
					Command: cmd,
				},
			},
			FailureThreshold: 3,
			PeriodSeconds:    10,
			SuccessThreshold: 1,
			TimeoutSeconds:   1,
		}
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
		podTemplate.Spec.SecurityContext.FSGroup = types.Int64P(999)
	}
	if podTemplate.Spec.SecurityContext.RunAsNonRoot == nil {
		podTemplate.Spec.SecurityContext.RunAsNonRoot = types.BoolP(true)
	}
	if podTemplate.Spec.SecurityContext.RunAsUser == nil {
		podTemplate.Spec.SecurityContext.RunAsUser = types.Int64P(999)
	}
}

func (m *MongoDBSpec) GetSecrets() []string {
	if m == nil {
		return nil
	}

	var secrets []string
	if m.DatabaseSecret != nil {
		secrets = append(secrets, m.DatabaseSecret.SecretName)
	}
	if m.KeyFile != nil {
		secrets = append(secrets, m.KeyFile.SecretName)
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
