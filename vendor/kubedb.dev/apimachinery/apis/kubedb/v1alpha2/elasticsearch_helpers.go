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
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"kubedb.dev/apimachinery/apis"
	catalog "kubedb.dev/apimachinery/apis/catalog/v1alpha1"
	"kubedb.dev/apimachinery/apis/kubedb"
	"kubedb.dev/apimachinery/crds"

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

const (
	ElasticsearchNodeAffinityTemplateVar = "NODE_ROLE"
)

func (_ Elasticsearch) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralElasticsearch))
}

var _ apis.ResourceInfo = &Elasticsearch{}

func (e Elasticsearch) OffshootName() string {
	return e.Name
}

func (e Elasticsearch) OffshootSelectors() map[string]string {
	return map[string]string{
		meta_util.NameLabelKey:      e.ResourceFQN(),
		meta_util.InstanceLabelKey:  e.Name,
		meta_util.ManagedByLabelKey: kubedb.GroupName,
	}
}

func (e Elasticsearch) NodeRoleSpecificLabelKey(roleType ElasticsearchNodeRoleType) string {
	return kubedb.GroupName + "/role-" + string(roleType)
}

func (e Elasticsearch) MasterSelectors() map[string]string {
	selectors := e.OffshootSelectors()
	selectors[e.NodeRoleSpecificLabelKey(ElasticsearchNodeRoleTypeMaster)] = ElasticsearchNodeRoleSet
	return selectors
}

func (e Elasticsearch) DataSelectors() map[string]string {
	selectors := e.OffshootSelectors()
	selectors[e.NodeRoleSpecificLabelKey(ElasticsearchNodeRoleTypeData)] = ElasticsearchNodeRoleSet
	return selectors
}

func (e Elasticsearch) IngestSelectors() map[string]string {
	selectors := e.OffshootSelectors()
	selectors[e.NodeRoleSpecificLabelKey(ElasticsearchNodeRoleTypeIngest)] = ElasticsearchNodeRoleSet
	return selectors
}

func (e Elasticsearch) NodeRoleSpecificSelectors(roleType ElasticsearchNodeRoleType) map[string]string {
	selectors := e.OffshootSelectors()
	selectors[e.NodeRoleSpecificLabelKey(roleType)] = ElasticsearchNodeRoleSet
	return selectors
}

func (e Elasticsearch) OffshootLabels() map[string]string {
	out := e.OffshootSelectors()
	out[meta_util.ComponentLabelKey] = ComponentDatabase
	return meta_util.FilterKeys(kubedb.GroupName, out, e.Labels)
}

func (e Elasticsearch) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", ResourcePluralElasticsearch, kubedb.GroupName)
}

func (e Elasticsearch) ResourceShortCode() string {
	return ResourceCodeElasticsearch
}

func (e Elasticsearch) ResourceKind() string {
	return ResourceKindElasticsearch
}

func (e Elasticsearch) ResourceSingular() string {
	return ResourceSingularElasticsearch
}

func (e Elasticsearch) ResourcePlural() string {
	return ResourcePluralElasticsearch
}

func (e Elasticsearch) ServiceName() string {
	return e.OffshootName()
}

func (e *Elasticsearch) MasterDiscoveryServiceName() string {
	return meta_util.NameWithSuffix(e.ServiceName(), "master")
}

func (e Elasticsearch) GoverningServiceName() string {
	return meta_util.NameWithSuffix(e.ServiceName(), "pods")
}

// CertificateName returns the default certificate name and/or certificate secret name for a certificate alias
func (e *Elasticsearch) CertificateName(alias ElasticsearchCertificateAlias) string {
	return meta_util.NameWithSuffix(e.Name, fmt.Sprintf("%s-cert", string(alias)))
}

// GetCertSecretName returns the secret name for a certificate alias if any,
// otherwise returns default certificate secret name for the given alias.
func (e *Elasticsearch) GetCertSecretName(alias ElasticsearchCertificateAlias) string {
	if e.Spec.TLS != nil {
		name, ok := kmapi.GetCertificateSecretName(e.Spec.TLS.Certificates, string(alias))
		if ok {
			return name
		}
	}
	return e.CertificateName(alias)
}

// ClientCertificateCN returns the CN for a client certificate
func (e *Elasticsearch) ClientCertificateCN(alias ElasticsearchCertificateAlias) string {
	return fmt.Sprintf("%s-%s", e.Name, string(alias))
}

// returns the volume name for certificate secret.
// Values will be like: transport-certs, http-certs etc.
func (e *Elasticsearch) CertSecretVolumeName(alias ElasticsearchCertificateAlias) string {
	return string(alias) + "-certs"
}

// returns the mountPath for certificate secrets.
// if configDir is "/usr/share/elasticsearch/config",
// mountPath will be, "/usr/share/elasticsearch/config/certs/<alias>".
func (e *Elasticsearch) CertSecretVolumeMountPath(configDir string, alias ElasticsearchCertificateAlias) string {
	return filepath.Join(configDir, "certs", string(alias))
}

// returns the default secret name for the  user credentials (ie. username, password)
// If username contains underscore (_), it will be replaced by hyphen (‐) for
// the Kubernetes naming convention.
func (e *Elasticsearch) DefaultUserCredSecretName(userName string) string {
	return meta_util.NameWithSuffix(e.Name, strings.ReplaceAll(fmt.Sprintf("%s-cred", userName), "_", "-"))
}

// Return the secret name for the given user.
// Return error, if the secret name is missing.
func (e *Elasticsearch) GetUserCredSecretName(username string) (string, error) {
	userSpec, err := getElasticsearchUser(e.Spec.InternalUsers, username)
	if err != nil {
		return "", err
	}
	if userSpec.SecretName == "" {
		return "", errors.New("secretName cannot be empty")
	}
	return userSpec.SecretName, nil
}

// returns the secret name for the default elasticsearch configuration
func (e *Elasticsearch) ConfigSecretName() string {
	return meta_util.NameWithSuffix(e.Name, "config")
}

func (e *Elasticsearch) GetConnectionScheme() string {
	scheme := "http"
	if e.Spec.EnableSSL {
		scheme = "https"
	}
	return scheme
}

func (e *Elasticsearch) GetConnectionURL() string {
	return fmt.Sprintf("%v://%s.%s:%d", e.GetConnectionScheme(), e.OffshootName(), e.Namespace, ElasticsearchRestPort)
}

func (e *Elasticsearch) CombinedStatefulSetName() string {
	return e.OffshootName()
}

func (e *Elasticsearch) MasterStatefulSetName() string {
	if e.Spec.Topology.Master.Suffix != "" {
		return meta_util.NameWithSuffix(e.OffshootName(), e.Spec.Topology.Master.Suffix)
	}
	return meta_util.NameWithSuffix(e.OffshootName(), string(ElasticsearchNodeRoleTypeMaster))
}

func (e *Elasticsearch) DataStatefulSetName() string {
	if e.Spec.Topology.Data.Suffix != "" {
		return meta_util.NameWithSuffix(e.OffshootName(), e.Spec.Topology.Data.Suffix)
	}
	return meta_util.NameWithSuffix(e.OffshootName(), string(ElasticsearchNodeRoleTypeData))
}

func (e *Elasticsearch) IngestStatefulSetName() string {
	if e.Spec.Topology.Ingest.Suffix != "" {
		return meta_util.NameWithSuffix(e.OffshootName(), e.Spec.Topology.Ingest.Suffix)
	}
	return meta_util.NameWithSuffix(e.OffshootName(), string(ElasticsearchNodeRoleTypeIngest))
}

func (e *Elasticsearch) DataContentStatefulSetName() string {
	if e.Spec.Topology.DataContent.Suffix != "" {
		return meta_util.NameWithSuffix(e.OffshootName(), e.Spec.Topology.DataContent.Suffix)
	}
	return meta_util.NameWithSuffix(e.OffshootName(), string(ElasticsearchNodeRoleTypeDataContent))
}

func (e *Elasticsearch) DataFrozenStatefulSetName() string {
	if e.Spec.Topology.DataFrozen.Suffix != "" {
		return meta_util.NameWithSuffix(e.OffshootName(), e.Spec.Topology.DataFrozen.Suffix)
	}
	return meta_util.NameWithSuffix(e.OffshootName(), string(ElasticsearchNodeRoleTypeDataFrozen))
}

func (e *Elasticsearch) DataColdStatefulSetName() string {
	if e.Spec.Topology.DataCold.Suffix != "" {
		return meta_util.NameWithSuffix(e.OffshootName(), e.Spec.Topology.DataCold.Suffix)
	}
	return meta_util.NameWithSuffix(e.OffshootName(), string(ElasticsearchNodeRoleTypeDataCold))
}

func (e *Elasticsearch) DataHotStatefulSetName() string {
	if e.Spec.Topology.DataHot.Suffix != "" {
		return meta_util.NameWithSuffix(e.OffshootName(), e.Spec.Topology.DataHot.Suffix)
	}
	return meta_util.NameWithSuffix(e.OffshootName(), string(ElasticsearchNodeRoleTypeDataHot))
}

func (e *Elasticsearch) DataWarmStatefulSetName() string {
	if e.Spec.Topology.DataWarm.Suffix != "" {
		return meta_util.NameWithSuffix(e.OffshootName(), e.Spec.Topology.DataWarm.Suffix)
	}
	return meta_util.NameWithSuffix(e.OffshootName(), string(ElasticsearchNodeRoleTypeDataWarm))
}

func (e *Elasticsearch) MLStatefulSetName() string {
	if e.Spec.Topology.ML.Suffix != "" {
		return meta_util.NameWithSuffix(e.OffshootName(), e.Spec.Topology.ML.Suffix)
	}
	return meta_util.NameWithSuffix(e.OffshootName(), string(ElasticsearchNodeRoleTypeML))
}

func (e *Elasticsearch) TransformStatefulSetName() string {
	if e.Spec.Topology.Transform.Suffix != "" {
		return meta_util.NameWithSuffix(e.OffshootName(), e.Spec.Topology.Transform.Suffix)
	}
	return meta_util.NameWithSuffix(e.OffshootName(), string(ElasticsearchNodeRoleTypeTransform))
}

func (e *Elasticsearch) CoordinatingStatefulSetName() string {
	if e.Spec.Topology.Coordinating.Suffix != "" {
		return meta_util.NameWithSuffix(e.OffshootName(), e.Spec.Topology.Coordinating.Suffix)
	}
	return meta_util.NameWithSuffix(e.OffshootName(), string(ElasticsearchNodeRoleTypeCoordinating))
}

func (e *Elasticsearch) InitialMasterNodes() []string {
	// For combined clusters
	stsName := e.CombinedStatefulSetName()
	replicas := int32(1)
	if e.Spec.Replicas != nil {
		replicas = pointer.Int32(e.Spec.Replicas)
	}

	// For topology cluster, overwrite the values
	if e.Spec.Topology != nil {
		stsName = e.MasterStatefulSetName()
		if e.Spec.Topology.Master.Replicas != nil {
			replicas = pointer.Int32(e.Spec.Topology.Master.Replicas)
		}
	}

	var nodeNames []string
	for i := int32(0); i < replicas; i++ {
		nodeNames = append(nodeNames, fmt.Sprintf("%s-%d", stsName, i))
	}

	return nodeNames
}

type elasticsearchApp struct {
	*Elasticsearch
}

func (r elasticsearchApp) Name() string {
	return r.Elasticsearch.Name
}

func (r elasticsearchApp) Type() appcat.AppType {
	return appcat.AppType(fmt.Sprintf("%s/%s", kubedb.GroupName, ResourceSingularElasticsearch))
}

func (e Elasticsearch) AppBindingMeta() appcat.AppBindingMeta {
	return &elasticsearchApp{&e}
}

type elasticsearchStatsService struct {
	*Elasticsearch
}

func (e elasticsearchStatsService) GetNamespace() string {
	return e.Elasticsearch.GetNamespace()
}

func (e elasticsearchStatsService) ServiceName() string {
	return e.OffshootName() + "-stats"
}

func (e elasticsearchStatsService) ServiceMonitorName() string {
	return e.ServiceName()
}

func (e elasticsearchStatsService) ServiceMonitorAdditionalLabels() map[string]string {
	return e.OffshootLabels()
}

func (e elasticsearchStatsService) Path() string {
	return DefaultStatsPath
}

func (e elasticsearchStatsService) Scheme() string {
	return ""
}

func (e Elasticsearch) StatsService() mona.StatsAccessor {
	return &elasticsearchStatsService{&e}
}

func (e Elasticsearch) StatsServiceLabels() map[string]string {
	lbl := meta_util.FilterKeys(kubedb.GroupName, e.OffshootSelectors(), e.Labels)
	lbl[LabelRole] = RoleStats
	return lbl
}

func (e *Elasticsearch) SetDefaults(esVersion *catalog.ElasticsearchVersion, topology *core_util.Topology) {
	if e == nil {
		return
	}

	if e.Spec.StorageType == "" {
		e.Spec.StorageType = StorageTypeDurable
	}

	if e.Spec.TerminationPolicy == "" {
		e.Spec.TerminationPolicy = TerminationPolicyDelete
	}

	if e.Spec.PodTemplate.Spec.ServiceAccountName == "" {
		e.Spec.PodTemplate.Spec.ServiceAccountName = e.OffshootName()
	}

	// set default elasticsearch node name prefix
	if e.Spec.Topology != nil {
		// Required nodes, must exist!
		// Default to "ingest"
		if e.Spec.Topology.Ingest.Suffix == "" {
			e.Spec.Topology.Ingest.Suffix = string(ElasticsearchNodeRoleTypeIngest)
		}
		SetDefaultResourceLimits(&e.Spec.Topology.Ingest.Resources, DefaultResources)
		if e.Spec.Topology.Ingest.Replicas == nil {
			e.Spec.Topology.Ingest.Replicas = pointer.Int32P(1)
		}

		// Required nodes, must exist!
		// Default to "master"
		if e.Spec.Topology.Master.Suffix == "" {
			e.Spec.Topology.Master.Suffix = string(ElasticsearchNodeRoleTypeMaster)
		}
		SetDefaultResourceLimits(&e.Spec.Topology.Master.Resources, DefaultResources)
		if e.Spec.Topology.Master.Replicas == nil {
			e.Spec.Topology.Master.Replicas = pointer.Int32P(1)
		}

		// Optional nodes, when other type of data nodes are not empty.
		// Otherwise required nodes.
		if e.Spec.Topology.Data != nil {
			// Default to "data"
			if e.Spec.Topology.Data.Suffix == "" {
				e.Spec.Topology.Data.Suffix = string(ElasticsearchNodeRoleTypeData)
			}
			SetDefaultResourceLimits(&e.Spec.Topology.Data.Resources, DefaultResources)
			if e.Spec.Topology.Data.Replicas == nil {
				e.Spec.Topology.Data.Replicas = pointer.Int32P(1)
			}
		}

		// Optional, can be empty
		if e.Spec.Topology.DataHot != nil {
			// Default to "data-hot"
			if e.Spec.Topology.DataHot.Suffix == "" {
				e.Spec.Topology.DataHot.Suffix = string(ElasticsearchNodeRoleTypeDataHot)
			}
			SetDefaultResourceLimits(&e.Spec.Topology.DataHot.Resources, DefaultResources)
			if e.Spec.Topology.DataHot.Replicas == nil {
				e.Spec.Topology.DataHot.Replicas = pointer.Int32P(1)
			}
		}

		// Optional, can be empty
		if e.Spec.Topology.DataWarm != nil {
			// Default to "data-warm"
			if e.Spec.Topology.DataWarm.Suffix == "" {
				e.Spec.Topology.DataWarm.Suffix = string(ElasticsearchNodeRoleTypeDataWarm)
			}
			SetDefaultResourceLimits(&e.Spec.Topology.DataWarm.Resources, DefaultResources)
			if e.Spec.Topology.DataWarm.Replicas == nil {
				e.Spec.Topology.DataWarm.Replicas = pointer.Int32P(1)
			}
		}

		// Optional, can be empty
		if e.Spec.Topology.DataCold != nil {
			// Default to "data-warm"
			if e.Spec.Topology.DataCold.Suffix == "" {
				e.Spec.Topology.DataCold.Suffix = string(ElasticsearchNodeRoleTypeDataCold)
			}
			SetDefaultResourceLimits(&e.Spec.Topology.DataCold.Resources, DefaultResources)
			if e.Spec.Topology.DataCold.Replicas == nil {
				e.Spec.Topology.DataCold.Replicas = pointer.Int32P(1)
			}
		}

		// Optional, can be empty
		if e.Spec.Topology.DataFrozen != nil {
			// Default to "data-frozen"
			if e.Spec.Topology.DataFrozen.Suffix == "" {
				e.Spec.Topology.DataFrozen.Suffix = string(ElasticsearchNodeRoleTypeDataFrozen)
			}
			SetDefaultResourceLimits(&e.Spec.Topology.DataFrozen.Resources, DefaultResources)
			if e.Spec.Topology.DataFrozen.Replicas == nil {
				e.Spec.Topology.DataFrozen.Replicas = pointer.Int32P(1)
			}
		}

		// Optional, can be empty
		if e.Spec.Topology.DataContent != nil {
			// Default to "data-content"
			if e.Spec.Topology.DataContent.Suffix == "" {
				e.Spec.Topology.DataContent.Suffix = string(ElasticsearchNodeRoleTypeDataContent)
			}
			SetDefaultResourceLimits(&e.Spec.Topology.DataContent.Resources, DefaultResources)
			if e.Spec.Topology.DataContent.Replicas == nil {
				e.Spec.Topology.DataContent.Replicas = pointer.Int32P(1)
			}
		}

		// Optional, can be empty
		if e.Spec.Topology.ML != nil {
			// Default to "ml"
			if e.Spec.Topology.ML.Suffix == "" {
				e.Spec.Topology.ML.Suffix = string(ElasticsearchNodeRoleTypeML)
			}
			SetDefaultResourceLimits(&e.Spec.Topology.ML.Resources, DefaultResources)
			if e.Spec.Topology.ML.Replicas == nil {
				e.Spec.Topology.ML.Replicas = pointer.Int32P(1)
			}
		}

		// Optional, can be empty
		if e.Spec.Topology.Transform != nil {
			// Default to "transform"
			if e.Spec.Topology.Transform.Suffix == "" {
				e.Spec.Topology.Transform.Suffix = string(ElasticsearchNodeRoleTypeTransform)
			}
			SetDefaultResourceLimits(&e.Spec.Topology.Transform.Resources, DefaultResources)
			if e.Spec.Topology.Transform.Replicas == nil {
				e.Spec.Topology.Transform.Replicas = pointer.Int32P(1)
			}
		}

	} else {
		SetDefaultResourceLimits(&e.Spec.PodTemplate.Spec.Resources, DefaultResources)
		if e.Spec.Replicas == nil {
			e.Spec.Replicas = pointer.Int32P(1)
		}
	}

	// set default kernel settings
	// -	Ref: https://www.elastic.co/guide/en/elasticsearch/reference/7.9/vm-max-map-count.html
	if e.Spec.KernelSettings == nil {
		e.Spec.KernelSettings = &KernelSettings{
			Privileged: true,
			Sysctls: []core.Sysctl{
				{
					Name:  "vm.max_map_count",
					Value: "262144",
				},
			},
		}
	}

	if e.Spec.PodTemplate.Spec.ContainerSecurityContext == nil {
		e.Spec.PodTemplate.Spec.ContainerSecurityContext = &core.SecurityContext{
			Privileged: pointer.BoolP(false),
			Capabilities: &core.Capabilities{
				Add: []core.Capability{"IPC_LOCK", "SYS_RESOURCE"},
			},
		}
	}

	// Add default Elasticsearch UID
	if e.Spec.PodTemplate.Spec.ContainerSecurityContext.RunAsUser == nil &&
		esVersion.Spec.SecurityContext.RunAsUser != nil {
		e.Spec.PodTemplate.Spec.ContainerSecurityContext.RunAsUser = esVersion.Spec.SecurityContext.RunAsUser
	}

	e.setDefaultAffinity(&e.Spec.PodTemplate, e.OffshootSelectors(), topology)
	e.SetTLSDefaults(esVersion)
	e.setDefaultInternalUsersAndRoleMappings(esVersion)
	e.Spec.Monitor.SetDefaults()
}

// setDefaultAffinity
func (e *Elasticsearch) setDefaultAffinity(podTemplate *ofst.PodTemplateSpec, labels map[string]string, topology *core_util.Topology) {
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
						Namespaces: []string{e.Namespace},
						LabelSelector: &metav1.LabelSelector{
							MatchLabels:      labels,
							MatchExpressions: e.GetMatchExpressions(),
						},

						TopologyKey: core.LabelHostname,
					},
				},
				// Prefer to not schedule multiple pods on the node with same zone
				{
					Weight: 50,
					PodAffinityTerm: core.PodAffinityTerm{
						Namespaces: []string{e.Namespace},
						LabelSelector: &metav1.LabelSelector{
							MatchLabels:      labels,
							MatchExpressions: e.GetMatchExpressions(),
						},
						TopologyKey: topology.LabelZone,
					},
				},
			},
		},
	}
}

// Set Default internal users settings
func (e *Elasticsearch) setDefaultInternalUsersAndRoleMappings(esVersion *catalog.ElasticsearchVersion) {
	// If security is disabled (ie. DisableSecurity: true), ignore.
	if e.Spec.DisableSecurity {
		return
	}

	// The internalUsers feature only works with searchGuard and openDistro
	if esVersion.Spec.Distribution == catalog.ElasticsearchDistroOpenDistro ||
		esVersion.Spec.Distribution == catalog.ElasticsearchDistroSearchGuard {

		inUsers := e.Spec.InternalUsers
		// If not set, create empty map
		if inUsers == nil {
			inUsers = make(map[string]ElasticsearchUserSpec)
		}

		// "Admin" user
		if userSpec, exists := inUsers[string(ElasticsearchInternalUserAdmin)]; !exists {
			inUsers[string(ElasticsearchInternalUserAdmin)] = ElasticsearchUserSpec{
				Reserved:     true,
				BackendRoles: []string{"admin"},
			}
		} else {
			// upsert "admin" role, if missing
			// Admin user must have the admin role
			userSpec.BackendRoles = upsertStringSlice(userSpec.BackendRoles, "admin")
			inUsers[string(ElasticsearchInternalUserAdmin)] = userSpec
		}

		// "Kibanaserver", "Kibanaro", "Logstash", "Readall", "Snapshotrestore"
		setMissingElasticsearchUser(inUsers, string(ElasticsearchInternalUserKibanaserver), ElasticsearchUserSpec{Reserved: true})
		setMissingElasticsearchUser(inUsers, string(ElasticsearchInternalUserKibanaro), ElasticsearchUserSpec{})
		setMissingElasticsearchUser(inUsers, string(ElasticsearchInternalUserLogstash), ElasticsearchUserSpec{})
		setMissingElasticsearchUser(inUsers, string(ElasticsearchInternalUserReadall), ElasticsearchUserSpec{})
		setMissingElasticsearchUser(inUsers, string(ElasticsearchInternalUserSnapshotrestore), ElasticsearchUserSpec{})
		// "MetricsExporter", Only if the monitoring is enabled.
		if e.Spec.Monitor != nil {
			setMissingElasticsearchUser(inUsers, string(ElasticsearchInternalUserMetricsExporter), ElasticsearchUserSpec{})
		}

		// Set missing user secret names
		for username, userSpec := range inUsers {
			// For admin user, spec.authSecret.Name must have high precedence over default field
			if username == string(ElasticsearchInternalUserAdmin) {
				if e.Spec.AuthSecret != nil && e.Spec.AuthSecret.Name != "" {
					userSpec.SecretName = e.Spec.AuthSecret.Name
				} else {
					if userSpec.SecretName == "" {
						userSpec.SecretName = e.DefaultUserCredSecretName(username)
					}
					e.Spec.AuthSecret = &core.LocalObjectReference{
						Name: userSpec.SecretName,
					}
				}
			} else if userSpec.SecretName == "" {
				userSpec.SecretName = e.DefaultUserCredSecretName(username)
			}
			inUsers[username] = userSpec
		}

		// If monitoring is enabled,
		// The "metric_exporter" user needs to have "readall_monitor" role mapped to itself.
		if e.Spec.Monitor != nil {
			rolesMapping := e.Spec.RolesMapping
			if rolesMapping == nil {
				rolesMapping = make(map[string]ElasticsearchRoleMapSpec)
			}
			var monitorRole string
			if esVersion.Spec.Distribution == catalog.ElasticsearchDistroSearchGuard {
				// readall_and_monitor role name varies in ES version
				// 	V7        = "SGS_READALL_AND_MONITOR"
				//	V6        = "sg_readall_and_monitor"
				if strings.HasPrefix(esVersion.Spec.Version, "6.") {
					monitorRole = ElasticsearchSearchGuardReadallMonitorRoleV6
					// Delete unsupported role, if any
					delete(rolesMapping, string(ElasticsearchSearchGuardReadallMonitorRoleV7))
				} else {
					monitorRole = ElasticsearchSearchGuardReadallMonitorRoleV7
					// Delete unsupported role, if any
					// Required during upgrade process, from v6 --> v7
					delete(rolesMapping, string(ElasticsearchSearchGuardReadallMonitorRoleV6))
				}
			} else {
				monitorRole = ElasticsearchOpendistroReadallMonitorRole
			}

			// Create rolesMapping if not exists.
			if value, exist := rolesMapping[monitorRole]; exist {
				value.Users = upsertStringSlice(value.Users, string(ElasticsearchInternalUserMetricsExporter))
				rolesMapping[monitorRole] = value
			} else {
				rolesMapping[monitorRole] = ElasticsearchRoleMapSpec{
					Users: []string{string(ElasticsearchInternalUserMetricsExporter)},
				}
			}
			e.Spec.RolesMapping = rolesMapping
		}
		e.Spec.InternalUsers = inUsers
	}
}

// set default tls configuration (ie. alias, secretName)
func (e *Elasticsearch) SetTLSDefaults(esVersion *catalog.ElasticsearchVersion) {
	// If security is disabled (ie. DisableSecurity: true), ignore.
	if e.Spec.DisableSecurity {
		return
	}

	tlsConfig := e.Spec.TLS
	if tlsConfig == nil {
		tlsConfig = &kmapi.TLSConfig{}
	}

	// If the issuerRef is nil, the operator will create the CA certificate.
	// It is required even if the spec.EnableSSL is false. Because, the transport
	// layer is always secured with certificates. Unless you turned off all the security
	// by setting spec.DisableSecurity to true.
	if tlsConfig.IssuerRef == nil {
		tlsConfig.Certificates = kmapi.SetMissingSpecForCertificate(tlsConfig.Certificates, kmapi.CertificateSpec{
			Alias:      string(ElasticsearchCACert),
			SecretName: e.CertificateName(ElasticsearchCACert),
		})
	}

	// transport layer is always secured with certificate
	tlsConfig.Certificates = kmapi.SetMissingSpecForCertificate(tlsConfig.Certificates, kmapi.CertificateSpec{
		Alias:      string(ElasticsearchTransportCert),
		SecretName: e.CertificateName(ElasticsearchTransportCert),
	})

	// If SSL is enabled, set missing certificate spec
	if e.Spec.EnableSSL {
		// http
		tlsConfig.Certificates = kmapi.SetMissingSpecForCertificate(tlsConfig.Certificates, kmapi.CertificateSpec{
			Alias:      string(ElasticsearchHTTPCert),
			SecretName: e.CertificateName(ElasticsearchHTTPCert),
		})

		// Set missing admin certificate spec, if authPlugin is either "OpenDistro" or "SearchGuard"
		if esVersion.Spec.Distribution == catalog.ElasticsearchDistroSearchGuard ||
			esVersion.Spec.Distribution == catalog.ElasticsearchDistroOpenDistro {
			tlsConfig.Certificates = kmapi.SetMissingSpecForCertificate(tlsConfig.Certificates, kmapi.CertificateSpec{
				Alias:      string(ElasticsearchAdminCert),
				SecretName: e.CertificateName(ElasticsearchAdminCert),
			})
		}

		// Set missing metrics-exporter certificate, if monitoring is enabled.
		if e.Spec.Monitor != nil {
			// matrics-exporter
			tlsConfig.Certificates = kmapi.SetMissingSpecForCertificate(tlsConfig.Certificates, kmapi.CertificateSpec{
				Alias:      string(ElasticsearchMetricsExporterCert),
				SecretName: e.CertificateName(ElasticsearchMetricsExporterCert),
			})
		}

		// archiver
		tlsConfig.Certificates = kmapi.SetMissingSpecForCertificate(tlsConfig.Certificates, kmapi.CertificateSpec{
			Alias:      string(ElasticsearchArchiverCert),
			SecretName: e.CertificateName(ElasticsearchArchiverCert),
		})
	}

	for id := range tlsConfig.Certificates {
		// Force overwrite the private key encoding type to PKCS#8
		tlsConfig.Certificates[id].PrivateKey = &kmapi.CertificatePrivateKey{
			Encoding: kmapi.PKCS8,
		}
		// Set default subject to O:KubeDB, if missing.
		// It isn't set from SetMissingSpecForCertificate(),
		// Because the default organization(ie. kubedb) gets merged, even if
		// the organizations[] isn't empty.
		if tlsConfig.Certificates[id].Subject == nil {
			tlsConfig.Certificates[id].Subject = &kmapi.X509Subject{
				Organizations: []string{KubeDBOrganization},
			}
		}
	}

	e.Spec.TLS = tlsConfig
}

func (e *Elasticsearch) GetMatchExpressions() []metav1.LabelSelectorRequirement {
	if e.Spec.Topology == nil {
		return nil
	}

	return []metav1.LabelSelectorRequirement{
		{
			Key:      fmt.Sprintf("${%s}", ElasticsearchNodeAffinityTemplateVar),
			Operator: metav1.LabelSelectorOpExists,
		},
	}
}

func (e *Elasticsearch) GetPersistentSecrets() []string {
	if e == nil {
		return nil
	}

	var secrets []string
	// Add Admin/Elastic user secret name
	if e.Spec.AuthSecret != nil {
		secrets = append(secrets, e.Spec.AuthSecret.Name)
	}

	// Skip for Admin/Elastic User.
	// Add other user cred secret names.
	if e.Spec.InternalUsers != nil {
		for user := range e.Spec.InternalUsers {
			if user == string(ElasticsearchInternalUserAdmin) || user == string(ElasticsearchInternalUserElastic) {
				continue
			}
			secretName, _ := e.GetUserCredSecretName(user)
			secrets = append(secrets, secretName)
		}
	}
	return secrets
}

func (e *Elasticsearch) ReplicasAreReady(lister appslister.StatefulSetLister) (bool, string, error) {
	// Desire number of statefulSets
	expectedItems := 1
	if e.Spec.Topology != nil {
		expectedItems = 3
	}
	return checkReplicas(lister.StatefulSets(e.Namespace), labels.SelectorFromSet(e.OffshootLabels()), expectedItems)
}

// returns true if the user exists.
// otherwise false.
func hasElasticsearchUser(userList map[string]ElasticsearchUserSpec, username string) bool {
	if _, exist := userList[username]; exist {
		return true
	}
	return false
}

// Set user if missing
func setMissingElasticsearchUser(userList map[string]ElasticsearchUserSpec, username string, userSpec ElasticsearchUserSpec) {
	if hasElasticsearchUser(userList, username) {
		return
	}
	userList[username] = userSpec
}

// Returns userSpec if exists
func getElasticsearchUser(userList map[string]ElasticsearchUserSpec, username string) (*ElasticsearchUserSpec, error) {
	if !hasElasticsearchUser(userList, username) {
		return nil, errors.New("user is missing")
	}
	userSpec := userList[username]
	return &userSpec, nil
}

// ToMap returns ClusterTopology in a Map
func (esTopology *ElasticsearchClusterTopology) ToMap() map[ElasticsearchNodeRoleType]ElasticsearchNode {
	topology := make(map[ElasticsearchNodeRoleType]ElasticsearchNode)
	topology[ElasticsearchNodeRoleTypeMaster] = esTopology.Master
	topology[ElasticsearchNodeRoleTypeIngest] = esTopology.Ingest
	if esTopology.Data != nil {
		topology[ElasticsearchNodeRoleTypeData] = *esTopology.Data
	}
	if esTopology.DataHot != nil {
		topology[ElasticsearchNodeRoleTypeDataHot] = *esTopology.DataHot
	}
	if esTopology.DataWarm != nil {
		topology[ElasticsearchNodeRoleTypeDataWarm] = *esTopology.DataWarm
	}
	if esTopology.DataCold != nil {
		topology[ElasticsearchNodeRoleTypeDataCold] = *esTopology.DataCold
	}
	if esTopology.DataFrozen != nil {
		topology[ElasticsearchNodeRoleTypeDataFrozen] = *esTopology.DataFrozen
	}
	if esTopology.DataContent != nil {
		topology[ElasticsearchNodeRoleTypeDataContent] = *esTopology.DataContent
	}
	if esTopology.ML != nil {
		topology[ElasticsearchNodeRoleTypeML] = *esTopology.ML
	}
	if esTopology.Transform != nil {
		topology[ElasticsearchNodeRoleTypeTransform] = *esTopology.Transform
	}
	if esTopology.Coordinating != nil {
		topology[ElasticsearchNodeRoleTypeCoordinating] = *esTopology.Coordinating
	}
	return topology
}
