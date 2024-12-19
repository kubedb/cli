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

// nolint:goconst
package v1

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"kubedb.dev/apimachinery/apis"
	catalog "kubedb.dev/apimachinery/apis/catalog/v1alpha1"
	"kubedb.dev/apimachinery/apis/kubedb"
	"kubedb.dev/apimachinery/crds"

	"github.com/Masterminds/semver/v3"
	promapi "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"gomodules.xyz/pointer"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/intstr"
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

const (
	ElasticsearchNodeAffinityTemplateVar = "NODE_ROLE"
)

func (*Elasticsearch) Hub() {}

func (_ Elasticsearch) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralElasticsearch))
}

func (e *Elasticsearch) AsOwner() *metav1.OwnerReference {
	return metav1.NewControllerRef(e, SchemeGroupVersion.WithKind(ResourceKindElasticsearch))
}

var _ apis.ResourceInfo = &Elasticsearch{}

func (e Elasticsearch) OffshootName() string {
	return e.Name
}

func (e Elasticsearch) OffshootSelectors(extraSelectors ...map[string]string) map[string]string {
	selector := map[string]string{
		meta_util.NameLabelKey:      e.ResourceFQN(),
		meta_util.InstanceLabelKey:  e.Name,
		meta_util.ManagedByLabelKey: kubedb.GroupName,
	}
	return meta_util.OverwriteKeys(selector, extraSelectors...)
}

func (e Elasticsearch) NodeRoleSpecificLabelKey(roleType ElasticsearchNodeRoleType) string {
	return kubedb.GroupName + "/role-" + string(roleType)
}

func (e Elasticsearch) MasterSelectors() map[string]string {
	return e.OffshootSelectors(map[string]string{e.NodeRoleSpecificLabelKey(ElasticsearchNodeRoleTypeMaster): kubedb.ElasticsearchNodeRoleSet})
}

func (e Elasticsearch) DataSelectors() map[string]string {
	return e.OffshootSelectors(map[string]string{e.NodeRoleSpecificLabelKey(ElasticsearchNodeRoleTypeData): kubedb.ElasticsearchNodeRoleSet})
}

func (e Elasticsearch) IngestSelectors() map[string]string {
	return e.OffshootSelectors(map[string]string{e.NodeRoleSpecificLabelKey(ElasticsearchNodeRoleTypeIngest): kubedb.ElasticsearchNodeRoleSet})
}

func (e Elasticsearch) NodeRoleSpecificSelectors(roleType ElasticsearchNodeRoleType) map[string]string {
	return e.OffshootSelectors(map[string]string{e.NodeRoleSpecificLabelKey(roleType): kubedb.ElasticsearchNodeRoleSet})
}

func (e Elasticsearch) OffshootLabels() map[string]string {
	return e.offshootLabels(e.OffshootSelectors(), nil)
}

func (e Elasticsearch) PodLabels(extraLabels ...map[string]string) map[string]string {
	return e.offshootLabels(meta_util.OverwriteKeys(e.OffshootSelectors(), extraLabels...), e.Spec.PodTemplate.Labels)
}

func (e Elasticsearch) PodControllerLabels(extraLabels ...map[string]string) map[string]string {
	return e.offshootLabels(meta_util.OverwriteKeys(e.OffshootSelectors(), extraLabels...), e.Spec.PodTemplate.Controller.Labels)
}

func (e Elasticsearch) ServiceLabels(alias ServiceAlias, extraLabels ...map[string]string) map[string]string {
	svcTemplate := GetServiceTemplate(e.Spec.ServiceTemplates, alias)
	return e.offshootLabels(meta_util.OverwriteKeys(e.OffshootSelectors(), extraLabels...), svcTemplate.Labels)
}

func (e Elasticsearch) offshootLabels(selector, override map[string]string) map[string]string {
	selector[meta_util.ComponentLabelKey] = kubedb.ComponentDatabase
	return meta_util.FilterKeys(kubedb.GroupName, selector, meta_util.OverwriteKeys(nil, e.Labels, override))
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

func (e Elasticsearch) GetAuthSecretName() string {
	if e.Spec.AuthSecret != nil && e.Spec.AuthSecret.Name != "" {
		return e.Spec.AuthSecret.Name
	}
	return meta_util.NameWithSuffix(e.OffshootName(), "auth")
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
	return fmt.Sprintf("%v://%s.%s:%d", e.GetConnectionScheme(), e.OffshootName(), e.Namespace, kubedb.ElasticsearchRestPort)
}

func (e *Elasticsearch) CombinedPetSetName() string {
	return e.OffshootName()
}

func (e *Elasticsearch) MasterPetSetName() string {
	if e.Spec.Topology.Master.Suffix != "" {
		return meta_util.NameWithSuffix(e.OffshootName(), e.Spec.Topology.Master.Suffix)
	}
	return meta_util.NameWithSuffix(e.OffshootName(), string(ElasticsearchNodeRoleTypeMaster))
}

func (e *Elasticsearch) DataPetSetName() string {
	if e.Spec.Topology.Data.Suffix != "" {
		return meta_util.NameWithSuffix(e.OffshootName(), e.Spec.Topology.Data.Suffix)
	}
	return meta_util.NameWithSuffix(e.OffshootName(), string(ElasticsearchNodeRoleTypeData))
}

func (e *Elasticsearch) IngestPetSetName() string {
	if e.Spec.Topology.Ingest.Suffix != "" {
		return meta_util.NameWithSuffix(e.OffshootName(), e.Spec.Topology.Ingest.Suffix)
	}
	return meta_util.NameWithSuffix(e.OffshootName(), string(ElasticsearchNodeRoleTypeIngest))
}

func (e *Elasticsearch) DataContentPetSetName() string {
	if e.Spec.Topology.DataContent.Suffix != "" {
		return meta_util.NameWithSuffix(e.OffshootName(), e.Spec.Topology.DataContent.Suffix)
	}
	return meta_util.NameWithSuffix(e.OffshootName(), string(ElasticsearchNodeRoleTypeDataContent))
}

func (e *Elasticsearch) DataFrozenPetSetName() string {
	if e.Spec.Topology.DataFrozen.Suffix != "" {
		return meta_util.NameWithSuffix(e.OffshootName(), e.Spec.Topology.DataFrozen.Suffix)
	}
	return meta_util.NameWithSuffix(e.OffshootName(), string(ElasticsearchNodeRoleTypeDataFrozen))
}

func (e *Elasticsearch) DataColdPetSetName() string {
	if e.Spec.Topology.DataCold.Suffix != "" {
		return meta_util.NameWithSuffix(e.OffshootName(), e.Spec.Topology.DataCold.Suffix)
	}
	return meta_util.NameWithSuffix(e.OffshootName(), string(ElasticsearchNodeRoleTypeDataCold))
}

func (e *Elasticsearch) DataHotPetSetName() string {
	if e.Spec.Topology.DataHot.Suffix != "" {
		return meta_util.NameWithSuffix(e.OffshootName(), e.Spec.Topology.DataHot.Suffix)
	}
	return meta_util.NameWithSuffix(e.OffshootName(), string(ElasticsearchNodeRoleTypeDataHot))
}

func (e *Elasticsearch) DataWarmPetSetName() string {
	if e.Spec.Topology.DataWarm.Suffix != "" {
		return meta_util.NameWithSuffix(e.OffshootName(), e.Spec.Topology.DataWarm.Suffix)
	}
	return meta_util.NameWithSuffix(e.OffshootName(), string(ElasticsearchNodeRoleTypeDataWarm))
}

func (e *Elasticsearch) MLPetSetName() string {
	if e.Spec.Topology.ML.Suffix != "" {
		return meta_util.NameWithSuffix(e.OffshootName(), e.Spec.Topology.ML.Suffix)
	}
	return meta_util.NameWithSuffix(e.OffshootName(), string(ElasticsearchNodeRoleTypeML))
}

func (e *Elasticsearch) TransformPetSetName() string {
	if e.Spec.Topology.Transform.Suffix != "" {
		return meta_util.NameWithSuffix(e.OffshootName(), e.Spec.Topology.Transform.Suffix)
	}
	return meta_util.NameWithSuffix(e.OffshootName(), string(ElasticsearchNodeRoleTypeTransform))
}

func (e *Elasticsearch) CoordinatingPetSetName() string {
	if e.Spec.Topology.Coordinating.Suffix != "" {
		return meta_util.NameWithSuffix(e.OffshootName(), e.Spec.Topology.Coordinating.Suffix)
	}
	return meta_util.NameWithSuffix(e.OffshootName(), string(ElasticsearchNodeRoleTypeCoordinating))
}

func (e *Elasticsearch) InitialMasterNodes() []string {
	// For combined clusters
	psName := e.CombinedPetSetName()
	replicas := int32(1)
	if e.Spec.Replicas != nil {
		replicas = pointer.Int32(e.Spec.Replicas)
	}

	// For topology cluster, overwrite the values
	if e.Spec.Topology != nil {
		psName = e.MasterPetSetName()
		if e.Spec.Topology.Master.Replicas != nil {
			replicas = pointer.Int32(e.Spec.Topology.Master.Replicas)
		}
	}

	var nodeNames []string
	for i := int32(0); i < replicas; i++ {
		nodeNames = append(nodeNames, fmt.Sprintf("%s-%d", psName, i))
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
	return kubedb.DefaultStatsPath
}

func (e elasticsearchStatsService) Scheme() string {
	return ""
}

func (e elasticsearchStatsService) TLSConfig() *promapi.TLSConfig {
	return nil
}

func (e Elasticsearch) StatsService() mona.StatsAccessor {
	return &elasticsearchStatsService{&e}
}

func (e Elasticsearch) StatsServiceLabels() map[string]string {
	return e.ServiceLabels(StatsServiceAlias, map[string]string{kubedb.LabelRole: kubedb.RoleStats})
}

func (e Elasticsearch) setContainerSecurityContextDefaults(esVersion *catalog.ElasticsearchVersion, podTemplate *ofstv2.PodTemplateSpec, containerName string) {
	if podTemplate == nil {
		return
	}

	if podTemplate.Spec.SecurityContext == nil {
		podTemplate.Spec.SecurityContext = &core.PodSecurityContext{}
	}
	if podTemplate.Spec.SecurityContext.FSGroup == nil {
		podTemplate.Spec.SecurityContext.FSGroup = esVersion.Spec.SecurityContext.RunAsUser
	}
	getContainers := func() []core.Container {
		if containerName == kubedb.ElasticsearchInitConfigMergerContainerName {
			return podTemplate.Spec.InitContainers
		}
		return podTemplate.Spec.Containers
	}
	container := core_util.GetContainerByName(getContainers(), containerName)
	if container == nil {
		container = &core.Container{
			Name: containerName,
		}
	}
	if container.SecurityContext == nil {
		container.SecurityContext = &core.SecurityContext{}
	}
	e.assignDefaultContainerSecurityContext(esVersion, container.SecurityContext)
	if containerName == kubedb.ElasticsearchInitConfigMergerContainerName {
		podTemplate.Spec.InitContainers = core_util.UpsertContainer(podTemplate.Spec.InitContainers, *container)
		return
	}
	podTemplate.Spec.Containers = core_util.UpsertContainer(podTemplate.Spec.Containers, *container)
}

func (e Elasticsearch) assignDefaultContainerSecurityContext(esVersion *catalog.ElasticsearchVersion, sc *core.SecurityContext) {
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
		sc.RunAsUser = esVersion.Spec.SecurityContext.RunAsUser
	}
	if sc.RunAsGroup == nil {
		sc.RunAsGroup = esVersion.Spec.SecurityContext.RunAsUser
	}
	if sc.SeccompProfile == nil {
		sc.SeccompProfile = secomp.DefaultSeccompProfile()
	}
}

func (e *Elasticsearch) SetDefaults(esVersion *catalog.ElasticsearchVersion) {
	if e == nil {
		return
	}

	if e.Spec.StorageType == "" {
		e.Spec.StorageType = StorageTypeDurable
	}

	if e.Spec.DeletionPolicy == "" {
		e.Spec.DeletionPolicy = DeletionPolicyDelete
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
		e.setContainerSecurityContextDefaults(esVersion, &e.Spec.Topology.Ingest.PodTemplate, kubedb.ElasticsearchContainerName)
		e.setContainerSecurityContextDefaults(esVersion, &e.Spec.Topology.Ingest.PodTemplate, kubedb.ElasticsearchInitConfigMergerContainerName)
		dbContainer := core_util.GetContainerByName(e.Spec.Topology.Ingest.PodTemplate.Spec.Containers, kubedb.ElasticsearchContainerName)
		if dbContainer != nil && (dbContainer.Resources.Requests == nil && dbContainer.Resources.Limits == nil) {
			apis.SetDefaultResourceLimits(&dbContainer.Resources, kubedb.DefaultResourcesMemoryIntensive)
		}
		if e.Spec.Topology.Ingest.Replicas == nil {
			e.Spec.Topology.Ingest.Replicas = pointer.Int32P(1)
		}
		if e.Spec.Topology.Ingest.MaxUnavailable == nil && *e.Spec.Topology.Ingest.Replicas > 1 {
			e.Spec.Topology.Ingest.MaxUnavailable = &intstr.IntOrString{IntVal: 1}
		}

		// Required nodes, must exist!
		// Default to "master"
		if e.Spec.Topology.Master.Suffix == "" {
			e.Spec.Topology.Master.Suffix = string(ElasticsearchNodeRoleTypeMaster)
		}
		e.setContainerSecurityContextDefaults(esVersion, &e.Spec.Topology.Master.PodTemplate, kubedb.ElasticsearchContainerName)
		e.setContainerSecurityContextDefaults(esVersion, &e.Spec.Topology.Master.PodTemplate, kubedb.ElasticsearchInitConfigMergerContainerName)
		dbContainer = core_util.GetContainerByName(e.Spec.Topology.Master.PodTemplate.Spec.Containers, kubedb.ElasticsearchContainerName)
		if dbContainer != nil && (dbContainer.Resources.Requests == nil && dbContainer.Resources.Limits == nil) {
			apis.SetDefaultResourceLimits(&dbContainer.Resources, kubedb.DefaultResourcesMemoryIntensive)
		}
		if e.Spec.Topology.Master.Replicas == nil {
			e.Spec.Topology.Master.Replicas = pointer.Int32P(1)
		}
		if e.Spec.Topology.Master.MaxUnavailable == nil && *e.Spec.Topology.Master.Replicas > 1 {
			e.Spec.Topology.Master.MaxUnavailable = &intstr.IntOrString{IntVal: 1}
		}

		// Optional nodes, when other type of data nodes are not empty.
		// Otherwise required nodes.
		if e.Spec.Topology.Data != nil {
			// Default to "data"
			if e.Spec.Topology.Data.Suffix == "" {
				e.Spec.Topology.Data.Suffix = string(ElasticsearchNodeRoleTypeData)
			}
			e.setContainerSecurityContextDefaults(esVersion, &e.Spec.Topology.Data.PodTemplate, kubedb.ElasticsearchContainerName)
			e.setContainerSecurityContextDefaults(esVersion, &e.Spec.Topology.Data.PodTemplate, kubedb.ElasticsearchInitConfigMergerContainerName)
			dbContainer = core_util.GetContainerByName(e.Spec.Topology.Data.PodTemplate.Spec.Containers, kubedb.ElasticsearchContainerName)
			if dbContainer != nil && (dbContainer.Resources.Requests == nil && dbContainer.Resources.Limits == nil) {
				apis.SetDefaultResourceLimits(&dbContainer.Resources, kubedb.DefaultResourcesMemoryIntensive)
			}
			if e.Spec.Topology.Data.Replicas == nil {
				e.Spec.Topology.Data.Replicas = pointer.Int32P(1)
			}
			if e.Spec.Topology.Data.MaxUnavailable == nil && *e.Spec.Topology.Data.Replicas > 1 {
				e.Spec.Topology.Data.MaxUnavailable = &intstr.IntOrString{IntVal: 1}
			}
		}

		// Optional, can be empty
		if e.Spec.Topology.DataHot != nil {
			// Default to "data-hot"
			if e.Spec.Topology.DataHot.Suffix == "" {
				e.Spec.Topology.DataHot.Suffix = string(ElasticsearchNodeRoleTypeDataHot)
			}
			e.setContainerSecurityContextDefaults(esVersion, &e.Spec.Topology.DataHot.PodTemplate, kubedb.ElasticsearchContainerName)
			e.setContainerSecurityContextDefaults(esVersion, &e.Spec.Topology.DataHot.PodTemplate, kubedb.ElasticsearchInitConfigMergerContainerName)
			dbContainer = core_util.GetContainerByName(e.Spec.Topology.DataHot.PodTemplate.Spec.Containers, kubedb.ElasticsearchContainerName)
			if dbContainer != nil && (dbContainer.Resources.Requests == nil && dbContainer.Resources.Limits == nil) {
				apis.SetDefaultResourceLimits(&dbContainer.Resources, kubedb.DefaultResourcesMemoryIntensive)
			}
			if e.Spec.Topology.DataHot.Replicas == nil {
				e.Spec.Topology.DataHot.Replicas = pointer.Int32P(1)
			}
			if e.Spec.Topology.DataHot.MaxUnavailable == nil && *e.Spec.Topology.DataHot.Replicas > 1 {
				e.Spec.Topology.DataHot.MaxUnavailable = &intstr.IntOrString{IntVal: 1}
			}
		}

		// Optional, can be empty
		if e.Spec.Topology.DataWarm != nil {
			// Default to "data-warm"
			if e.Spec.Topology.DataWarm.Suffix == "" {
				e.Spec.Topology.DataWarm.Suffix = string(ElasticsearchNodeRoleTypeDataWarm)
			}
			e.setContainerSecurityContextDefaults(esVersion, &e.Spec.Topology.DataWarm.PodTemplate, kubedb.ElasticsearchContainerName)
			e.setContainerSecurityContextDefaults(esVersion, &e.Spec.Topology.DataWarm.PodTemplate, kubedb.ElasticsearchInitConfigMergerContainerName)
			dbContainer = core_util.GetContainerByName(e.Spec.Topology.DataWarm.PodTemplate.Spec.Containers, kubedb.ElasticsearchContainerName)
			if dbContainer != nil && (dbContainer.Resources.Requests == nil && dbContainer.Resources.Limits == nil) {
				apis.SetDefaultResourceLimits(&dbContainer.Resources, kubedb.DefaultResourcesMemoryIntensive)
			}
			if e.Spec.Topology.DataWarm.Replicas == nil {
				e.Spec.Topology.DataWarm.Replicas = pointer.Int32P(1)
			}
			if e.Spec.Topology.DataWarm.MaxUnavailable == nil && *e.Spec.Topology.DataWarm.Replicas > 1 {
				e.Spec.Topology.DataWarm.MaxUnavailable = &intstr.IntOrString{IntVal: 1}
			}
		}

		// Optional, can be empty
		if e.Spec.Topology.DataCold != nil {
			// Default to "data-warm"
			if e.Spec.Topology.DataCold.Suffix == "" {
				e.Spec.Topology.DataCold.Suffix = string(ElasticsearchNodeRoleTypeDataCold)
			}
			e.setContainerSecurityContextDefaults(esVersion, &e.Spec.Topology.DataCold.PodTemplate, kubedb.ElasticsearchContainerName)
			e.setContainerSecurityContextDefaults(esVersion, &e.Spec.Topology.DataCold.PodTemplate, kubedb.ElasticsearchInitConfigMergerContainerName)
			dbContainer = core_util.GetContainerByName(e.Spec.Topology.DataCold.PodTemplate.Spec.Containers, kubedb.ElasticsearchContainerName)
			if dbContainer != nil && (dbContainer.Resources.Requests == nil && dbContainer.Resources.Limits == nil) {
				apis.SetDefaultResourceLimits(&dbContainer.Resources, kubedb.DefaultResourcesMemoryIntensive)
			}
			if e.Spec.Topology.DataCold.Replicas == nil {
				e.Spec.Topology.DataCold.Replicas = pointer.Int32P(1)
			}
			if e.Spec.Topology.DataCold.MaxUnavailable == nil && *e.Spec.Topology.DataCold.Replicas > 1 {
				e.Spec.Topology.DataCold.MaxUnavailable = &intstr.IntOrString{IntVal: 1}
			}
		}

		// Optional, can be empty
		if e.Spec.Topology.DataFrozen != nil {
			// Default to "data-frozen"
			if e.Spec.Topology.DataFrozen.Suffix == "" {
				e.Spec.Topology.DataFrozen.Suffix = string(ElasticsearchNodeRoleTypeDataFrozen)
			}
			e.setContainerSecurityContextDefaults(esVersion, &e.Spec.Topology.DataFrozen.PodTemplate, kubedb.ElasticsearchContainerName)
			e.setContainerSecurityContextDefaults(esVersion, &e.Spec.Topology.DataFrozen.PodTemplate, kubedb.ElasticsearchInitConfigMergerContainerName)
			dbContainer = core_util.GetContainerByName(e.Spec.Topology.DataFrozen.PodTemplate.Spec.Containers, kubedb.ElasticsearchContainerName)
			if dbContainer != nil && (dbContainer.Resources.Requests == nil && dbContainer.Resources.Limits == nil) {
				apis.SetDefaultResourceLimits(&dbContainer.Resources, kubedb.DefaultResourcesMemoryIntensive)
			}
			if e.Spec.Topology.DataFrozen.Replicas == nil {
				e.Spec.Topology.DataFrozen.Replicas = pointer.Int32P(1)
			}
			if e.Spec.Topology.DataFrozen.MaxUnavailable == nil && *e.Spec.Topology.DataFrozen.Replicas > 1 {
				e.Spec.Topology.DataFrozen.MaxUnavailable = &intstr.IntOrString{IntVal: 1}
			}
		}

		// Optional, can be empty
		if e.Spec.Topology.DataContent != nil {
			// Default to "data-content"
			if e.Spec.Topology.DataContent.Suffix == "" {
				e.Spec.Topology.DataContent.Suffix = string(ElasticsearchNodeRoleTypeDataContent)
			}
			e.setContainerSecurityContextDefaults(esVersion, &e.Spec.Topology.DataContent.PodTemplate, kubedb.ElasticsearchContainerName)
			e.setContainerSecurityContextDefaults(esVersion, &e.Spec.Topology.DataContent.PodTemplate, kubedb.ElasticsearchInitConfigMergerContainerName)
			dbContainer = core_util.GetContainerByName(e.Spec.Topology.DataContent.PodTemplate.Spec.Containers, kubedb.ElasticsearchContainerName)
			if dbContainer != nil && (dbContainer.Resources.Requests == nil && dbContainer.Resources.Limits == nil) {
				apis.SetDefaultResourceLimits(&dbContainer.Resources, kubedb.DefaultResourcesMemoryIntensive)
			}
			if e.Spec.Topology.DataContent.Replicas == nil {
				e.Spec.Topology.DataContent.Replicas = pointer.Int32P(1)
			}
			if e.Spec.Topology.DataContent.MaxUnavailable == nil && *e.Spec.Topology.DataContent.Replicas > 1 {
				e.Spec.Topology.DataContent.MaxUnavailable = &intstr.IntOrString{IntVal: 1}
			}
		}

		// Optional, can be empty
		if e.Spec.Topology.ML != nil {
			// Default to "ml"
			if e.Spec.Topology.ML.Suffix == "" {
				e.Spec.Topology.ML.Suffix = string(ElasticsearchNodeRoleTypeML)
			}
			e.setContainerSecurityContextDefaults(esVersion, &e.Spec.Topology.ML.PodTemplate, kubedb.ElasticsearchContainerName)
			e.setContainerSecurityContextDefaults(esVersion, &e.Spec.Topology.ML.PodTemplate, kubedb.ElasticsearchInitConfigMergerContainerName)
			dbContainer = core_util.GetContainerByName(e.Spec.Topology.ML.PodTemplate.Spec.Containers, kubedb.ElasticsearchContainerName)
			if dbContainer != nil && (dbContainer.Resources.Requests == nil && dbContainer.Resources.Limits == nil) {
				apis.SetDefaultResourceLimits(&dbContainer.Resources, kubedb.DefaultResourcesMemoryIntensive)
			}
			if e.Spec.Topology.ML.Replicas == nil {
				e.Spec.Topology.ML.Replicas = pointer.Int32P(1)
			}
			if e.Spec.Topology.ML.MaxUnavailable == nil && *e.Spec.Topology.ML.Replicas > 1 {
				e.Spec.Topology.ML.MaxUnavailable = &intstr.IntOrString{IntVal: 1}
			}
		}

		// Optional, can be empty
		if e.Spec.Topology.Transform != nil {
			// Default to "transform"
			if e.Spec.Topology.Transform.Suffix == "" {
				e.Spec.Topology.Transform.Suffix = string(ElasticsearchNodeRoleTypeTransform)
			}
			e.setContainerSecurityContextDefaults(esVersion, &e.Spec.Topology.Transform.PodTemplate, kubedb.ElasticsearchContainerName)
			e.setContainerSecurityContextDefaults(esVersion, &e.Spec.Topology.Transform.PodTemplate, kubedb.ElasticsearchInitConfigMergerContainerName)
			dbContainer = core_util.GetContainerByName(e.Spec.Topology.Transform.PodTemplate.Spec.Containers, kubedb.ElasticsearchContainerName)
			if dbContainer != nil && (dbContainer.Resources.Requests == nil && dbContainer.Resources.Limits == nil) {
				apis.SetDefaultResourceLimits(&dbContainer.Resources, kubedb.DefaultResourcesMemoryIntensive)
			}
			if e.Spec.Topology.Transform.Replicas == nil {
				e.Spec.Topology.Transform.Replicas = pointer.Int32P(1)
			}
			if e.Spec.Topology.Transform.MaxUnavailable == nil && *e.Spec.Topology.Transform.Replicas > 1 {
				e.Spec.Topology.Transform.MaxUnavailable = &intstr.IntOrString{IntVal: 1}
			}
		}

	} else {
		e.setContainerSecurityContextDefaults(esVersion, &e.Spec.PodTemplate, kubedb.ElasticsearchContainerName)
		e.setContainerSecurityContextDefaults(esVersion, &e.Spec.PodTemplate, kubedb.ElasticsearchInitConfigMergerContainerName)
		dbContainer := core_util.GetContainerByName(e.Spec.PodTemplate.Spec.Containers, kubedb.ElasticsearchContainerName)
		if dbContainer != nil && (dbContainer.Resources.Requests == nil && dbContainer.Resources.Limits == nil) {
			apis.SetDefaultResourceLimits(&dbContainer.Resources, kubedb.DefaultResourcesMemoryIntensive)
		}
		if e.Spec.Replicas == nil {
			e.Spec.Replicas = pointer.Int32P(1)
		}
		if e.Spec.MaxUnavailable == nil && *e.Spec.Replicas > 1 {
			e.Spec.MaxUnavailable = &intstr.IntOrString{IntVal: 1}
		}
	}

	// set default kernel settings
	// -	Ref: https://www.elastic.co/guide/en/elasticsearch/reference/7.9/vm-max-map-count.html
	// if kernelSettings defaults is enabled systls-init container will be injected with the default vm_map_count settings
	// if not init container will not be injected and default values will not be set
	if e.Spec.KernelSettings == nil {
		e.Spec.KernelSettings = &KernelSettings{
			DisableDefaults: false,
		}
	}
	if !e.Spec.KernelSettings.DisableDefaults {
		e.Spec.KernelSettings.Privileged = true
		vmMapCountNotSet := true
		if len(e.Spec.KernelSettings.Sysctls) != 0 {
			for i := 0; i < len(e.Spec.KernelSettings.Sysctls); i++ {
				if e.Spec.KernelSettings.Sysctls[i].Name == "vm.max_map_count" {
					vmMapCountNotSet = false
					break
				}
			}
		}
		if vmMapCountNotSet {
			e.Spec.KernelSettings.Sysctls = append(e.Spec.KernelSettings.Sysctls, core.Sysctl{
				Name:  "vm.max_map_count",
				Value: "262144",
			})
		}
	}

	e.SetDefaultInternalUsersAndRoleMappings(esVersion)
	e.SetMetricsExporterDefaults(esVersion)
	e.SetTLSDefaults(esVersion)
}

func (e *Elasticsearch) SetMetricsExporterDefaults(esVersion *catalog.ElasticsearchVersion) {
	e.Spec.Monitor.SetDefaults()
	if e.Spec.Monitor != nil && e.Spec.Monitor.Prometheus != nil {
		if e.Spec.Monitor.Prometheus.Exporter.SecurityContext.RunAsUser == nil {
			e.Spec.Monitor.Prometheus.Exporter.SecurityContext.RunAsUser = esVersion.Spec.SecurityContext.RunAsUser
		}
		if e.Spec.Monitor.Prometheus.Exporter.SecurityContext.RunAsGroup == nil {
			e.Spec.Monitor.Prometheus.Exporter.SecurityContext.RunAsGroup = esVersion.Spec.SecurityContext.RunAsUser
		}
	}
}

// Set Default internal users settings
func (e *Elasticsearch) SetDefaultInternalUsersAndRoleMappings(esVersion *catalog.ElasticsearchVersion) {
	// If security is disabled (ie. DisableSecurity: true), ignore.
	if e.Spec.DisableSecurity {
		return
	}

	version, err := semver.NewVersion(esVersion.Spec.Version)
	if err != nil {
		return
	}
	// set missing internal users for Xpack,
	// internal users are supported for version>=7.8.x
	if esVersion.Spec.AuthPlugin == catalog.ElasticsearchAuthPluginXpack &&
		(version.Major() >= 8 || (version.Major() == 7 && version.Minor() >= 8)) {
		inUsers := e.Spec.InternalUsers
		// If not set, create empty map
		if inUsers == nil {
			inUsers = make(map[string]ElasticsearchUserSpec)
		}

		// "elastic" user
		if userSpec, exists := inUsers[string(ElasticsearchInternalUserElastic)]; !exists {
			inUsers[string(ElasticsearchInternalUserElastic)] = ElasticsearchUserSpec{
				BackendRoles: []string{"superuser"},
			}
		} else {
			// upsert "superuser" role, if missing
			// elastic user must have the superuser role
			userSpec.BackendRoles = upsertStringSlice(userSpec.BackendRoles, "superuser")
			inUsers[string(ElasticsearchInternalUserElastic)] = userSpec
		}

		// "Kibana_system", "logstash_system", "beats_system", "apm_system", "remote_monitoring_user" user
		setMissingElasticsearchUser(inUsers, string(ElasticsearchInternalUserKibanaSystem), ElasticsearchUserSpec{
			BackendRoles: []string{"kibana_system"},
		})
		setMissingElasticsearchUser(inUsers, string(ElasticsearchInternalUserBeatsSystem), ElasticsearchUserSpec{
			BackendRoles: []string{"beats_system"},
		})
		setMissingElasticsearchUser(inUsers, string(ElasticsearchInternalUserApmSystem), ElasticsearchUserSpec{
			BackendRoles: []string{"apm_system"},
		})
		setMissingElasticsearchUser(inUsers, string(ElasticsearchInternalUserRemoteMonitoringUser), ElasticsearchUserSpec{
			BackendRoles: []string{"remote_monitoring_collector", "remote_monitoring_agent"},
		})
		setMissingElasticsearchUser(inUsers, string(ElasticsearchInternalUserLogstashSystem), ElasticsearchUserSpec{
			BackendRoles: []string{"logstash_system"},
		})

		e.Spec.InternalUsers = inUsers
	}

	// set missing internal users and roles for OpenSearch
	if esVersion.Spec.AuthPlugin == catalog.ElasticsearchAuthPluginOpenSearch {

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

		// If monitoring is enabled,
		// The "metric_exporter" user needs to have "readall_monitor" role mapped to itself.
		if e.Spec.Monitor != nil {
			rolesMapping := e.Spec.RolesMapping
			if rolesMapping == nil {
				rolesMapping = make(map[string]ElasticsearchRoleMapSpec)
			}
			var monitorRole string
			if esVersion.Spec.AuthPlugin == catalog.ElasticsearchAuthPluginSearchGuard {
				// readall_and_monitor role name varies in ES version
				// 	V7        = "SGS_READALL_AND_MONITOR"
				//	V6        = "sg_readall_and_monitor"
				if strings.HasPrefix(esVersion.Spec.Version, "6.") {
					monitorRole = kubedb.ElasticsearchSearchGuardReadallMonitorRoleV6
					// Delete unsupported role, if any
					delete(rolesMapping, string(kubedb.ElasticsearchSearchGuardReadallMonitorRoleV7))
				} else {
					monitorRole = kubedb.ElasticsearchSearchGuardReadallMonitorRoleV7
					// Delete unsupported role, if any
					// Required during upgrade process, from v6 --> v7
					delete(rolesMapping, string(kubedb.ElasticsearchSearchGuardReadallMonitorRoleV6))
				}
			} else if esVersion.Spec.AuthPlugin == catalog.ElasticsearchAuthPluginOpenDistro {
				monitorRole = kubedb.ElasticsearchOpendistroReadallMonitorRole
			} else {
				monitorRole = kubedb.ElasticsearchOpenSearchReadallMonitorRole
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

	inUsers := e.Spec.InternalUsers
	// Set missing user secret names
	for username, userSpec := range inUsers {
		// For admin user, spec.authSecret.Name must have high precedence over default field
		if username == string(ElasticsearchInternalUserAdmin) || username == string(ElasticsearchInternalUserElastic) {
			if e.Spec.AuthSecret != nil && e.Spec.AuthSecret.Name != "" {
				userSpec.SecretName = e.Spec.AuthSecret.Name
			} else {
				if userSpec.SecretName == "" {
					userSpec.SecretName = e.GetAuthSecretName()
				}
				e.Spec.AuthSecret = &SecretReference{
					LocalObjectReference: core.LocalObjectReference{
						Name: userSpec.SecretName,
					},
				}
			}
		} else if userSpec.SecretName == "" {
			userSpec.SecretName = e.DefaultUserCredSecretName(username)
		}
		inUsers[username] = userSpec
	}
	e.Spec.InternalUsers = inUsers
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

	// Set missing admin certificate spec, if authPlugin is "OpenDistro", "SearchGuard", or "OpenSearch"
	// Create the admin certificate, even if the enable.SSL is false. This is necessary to securityadmin.sh command.
	if esVersion.Spec.AuthPlugin == catalog.ElasticsearchAuthPluginSearchGuard ||
		esVersion.Spec.AuthPlugin == catalog.ElasticsearchAuthPluginOpenDistro ||
		esVersion.Spec.AuthPlugin == catalog.ElasticsearchAuthPluginOpenSearch {
		tlsConfig.Certificates = kmapi.SetMissingSpecForCertificate(tlsConfig.Certificates, kmapi.CertificateSpec{
			Alias:      string(ElasticsearchAdminCert),
			SecretName: e.CertificateName(ElasticsearchAdminCert),
		})
	}

	// If SSL is enabled, set missing certificate spec
	if e.Spec.EnableSSL {
		// http
		tlsConfig.Certificates = kmapi.SetMissingSpecForCertificate(tlsConfig.Certificates, kmapi.CertificateSpec{
			Alias:      string(ElasticsearchHTTPCert),
			SecretName: e.CertificateName(ElasticsearchHTTPCert),
		})

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
			Alias:      string(ElasticsearchClientCert),
			SecretName: e.CertificateName(ElasticsearchClientCert),
		})
	}

	// remove archiverCert from old spec if exists
	tlsConfig.Certificates = kmapi.RemoveCertificate(tlsConfig.Certificates, string(ElasticsearchArchiverCert))

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
				Organizations: []string{kubedb.KubeDBOrganization},
			}
		}
	}

	e.Spec.TLS = tlsConfig
}

func (e *Elasticsearch) SetHealthCheckerDefaults() {
	if e.Spec.HealthChecker.PeriodSeconds == nil {
		e.Spec.HealthChecker.PeriodSeconds = pointer.Int32P(10)
	}
	if e.Spec.HealthChecker.TimeoutSeconds == nil {
		e.Spec.HealthChecker.TimeoutSeconds = pointer.Int32P(10)
	}
	if e.Spec.HealthChecker.FailureThreshold == nil {
		e.Spec.HealthChecker.FailureThreshold = pointer.Int32P(1)
	}
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

func (e *Elasticsearch) ReplicasAreReady(lister pslister.PetSetLister) (bool, string, error) {
	// Desire number of statefulSets
	expectedItems := 1
	if e.Spec.Topology != nil {
		expectedItems = 3
	}
	return checkReplicas(lister.PetSets(e.Namespace), labels.SelectorFromSet(e.OffshootLabels()), expectedItems)
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
