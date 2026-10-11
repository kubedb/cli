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

	"kubedb.dev/apimachinery/apis"
	catalog "kubedb.dev/apimachinery/apis/catalog/v1alpha1"
	"kubedb.dev/apimachinery/apis/kubedb"
	"kubedb.dev/apimachinery/crds"

	promapi "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"gomodules.xyz/pointer"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
	kmapi "kmodules.xyz/client-go/api/v1"
	"kmodules.xyz/client-go/apiextensions"
	coreutil "kmodules.xyz/client-go/core/v1"
	meta_util "kmodules.xyz/client-go/meta"
	"kmodules.xyz/client-go/policy/secomp"
	appcat "kmodules.xyz/custom-resources/apis/appcatalog/v1alpha1"
	mona "kmodules.xyz/monitoring-agent-api/api/v1"
	ofstv2 "kmodules.xyz/offshoot-api/api/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type MilvusApp struct {
	*Milvus
}

func (Milvus) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralMilvus))
}

func (m *Milvus) ResourceKind() string {
	return ResourceKindMilvus
}

func (m *Milvus) ResourceSingular() string {
	return ResourceSingularMilvus
}

func (m *Milvus) ResourcePlural() string {
	return ResourcePluralMilvus
}

func (m *Milvus) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", m.ResourcePlural(), kubedb.GroupName)
}

func (m *Milvus) AppBindingMeta() appcat.AppBindingMeta {
	return &MilvusApp{m}
}

func (r MilvusApp) Name() string {
	return r.Milvus.Name
}

func (m Milvus) Type() appcat.AppType {
	return appcat.AppType(fmt.Sprintf("%s/%s", kubedb.GroupName, m.ResourceSingular()))
}

func (m *Milvus) GetConnectionScheme() string {
	scheme := "http"
	if m.Spec.TLS != nil && m.Spec.TLS.External != nil && m.Spec.TLS.External.Mode != TLSModeDisabled {
		scheme = "https"
	}
	return scheme
}

func (m *Milvus) OffshootName() string {
	return m.Name
}

func (m *Milvus) ServiceName() string {
	return m.OffshootName()
}

func (m *Milvus) PetSetName(nodeRole MilvusNodeRoleType) string {
	if m.IsDistributed() {
		return meta_util.NameWithSuffix(m.OffshootName(), string(nodeRole))
	}
	return m.OffshootName()
}

// PetSetNameForGroup returns the PetSet name for one MilvusNodeGroup within
// a Distributed role -- <db>-<role>-<group>. When groupName is empty, this
// is identical to PetSetName(nodeRole), covering the role's default
// (ungrouped) PetSet.
func (m *Milvus) PetSetNameForGroup(nodeRole MilvusNodeRoleType, groupName string) string {
	if groupName == "" {
		return m.PetSetName(nodeRole)
	}
	return meta_util.NameWithSuffix(m.PetSetName(nodeRole), groupName)
}

// GetNodeGroups returns the MilvusNodeGroup list for the given Distributed
// role, or nil if the role has no groups configured (the common case --
// one PetSet per role, built from the role-level Replicas/PodTemplate/
// Network/GPU fields instead).
func (m *Milvus) GetNodeGroups(nodeType MilvusNodeRoleType) []MilvusNodeGroup {
	nodeSpec, dataNodeSpec := m.GetNodeSpec(nodeType)
	if nodeSpec != nil {
		return nodeSpec.Groups
	}
	if dataNodeSpec != nil {
		return dataNodeSpec.Groups
	}
	return nil
}

func (m *Milvus) GetNodeSpec(nodeType MilvusNodeRoleType) (*MilvusNode, *MilvusDataNode) {
	switch nodeType {
	case MilvusNodeRoleMixCoord:
		return m.Spec.Topology.Distributed.MixCoord, nil
	case MilvusNodeRoleDataNode:
		return m.Spec.Topology.Distributed.DataNode, nil
	case MilvusNodeRoleProxy:
		return m.Spec.Topology.Distributed.Proxy, nil
	case MilvusNodeRoleQueryNode:
		return m.Spec.Topology.Distributed.QueryNode, nil
	case MilvusNodeRoleStreamingNode:
		return nil, m.Spec.Topology.Distributed.StreamingNode
	default:
		klog.Errorf("unknown milvus node role %s\n", nodeType)
		return nil, nil
	}
}

// GetNodeGPUSpec returns the effective *MilvusGPUSpec for the given role in
// Distributed mode, or nil if unset or the role has no such spec.
func (m *Milvus) GetNodeGPUSpec(nodeType MilvusNodeRoleType) *MilvusGPUSpec {
	nodeSpec, dataNodeSpec := m.GetNodeSpec(nodeType)
	if nodeSpec != nil {
		return nodeSpec.GPU
	}
	if dataNodeSpec != nil {
		return dataNodeSpec.GPU
	}
	return nil
}

// GetNodeNetworkSpec returns the effective *MilvusNetworkSpec for the given
// role in Distributed mode, or nil if unset or the role has no such spec.
func (m *Milvus) GetNodeNetworkSpec(nodeType MilvusNodeRoleType) *MilvusNetworkSpec {
	nodeSpec, dataNodeSpec := m.GetNodeSpec(nodeType)
	if nodeSpec != nil {
		return nodeSpec.Network
	}
	if dataNodeSpec != nil {
		return dataNodeSpec.Network
	}
	return nil
}

// podTemplateRequestsGPU reports whether any container in podTemplate has a
// resources.requests or resources.limits entry for MilvusGPUDefaultResourceName.
// This catches a user hand-writing "nvidia.com/gpu" directly under
// podTemplate.spec.containers[].resources, independent of whether the typed
// GPU field is also set.
func podTemplateRequestsGPU(podTemplate *ofstv2.PodTemplateSpec) bool {
	if podTemplate == nil {
		return false
	}
	containsGPU := func(rl core.ResourceList) bool {
		if rl == nil {
			return false
		}
		_, ok := rl[MilvusGPUDefaultResourceName]
		return ok
	}
	for _, c := range podTemplate.Spec.Containers {
		if containsGPU(c.Resources.Requests) || containsGPU(c.Resources.Limits) {
			return true
		}
	}
	// Kubernetes permits GPU resources on init containers too (they're
	// counted against the node's allocatable the same as any other
	// container's); skipping them here would let an init-container GPU
	// request slip past milvusValidateGPU undetected.
	for _, c := range podTemplate.Spec.InitContainers {
		if containsGPU(c.Resources.Requests) || containsGPU(c.Resources.Limits) {
			return true
		}
	}
	return false
}

// RequestsGPU reports whether this Milvus, in whichever mode it's configured
// for, requests a GPU anywhere -- via the typed GPU field (Standalone's
// spec.gpu, or any Distributed role's .gpu), or via a hand-written
// nvidia.com/gpu resource request/limit on a podTemplate container. Used by
// the admission webhook to decide whether spec.version's MilvusVersion must
// declare GPU support.
func (m *Milvus) RequestsGPU() bool {
	if m.Spec.GPU != nil || podTemplateRequestsGPU(m.Spec.PodTemplate) {
		return true
	}
	if !m.IsDistributed() || m.Spec.Topology.Distributed == nil {
		return false
	}
	for _, nodeType := range []MilvusNodeRoleType{
		MilvusNodeRoleMixCoord, MilvusNodeRoleDataNode, MilvusNodeRoleProxy,
		MilvusNodeRoleQueryNode, MilvusNodeRoleStreamingNode,
	} {
		// When Groups is set, the role-level GPU/PodTemplate fields are
		// ignored for pod-building (nodes.go's ensureNodeOrGroups never
		// reads them) -- checking them here too would let a stale
		// role-level GPU setting force a CPU-only MilvusVersion to be
		// rejected even though every rendered group is CPU-only.
		groups := m.GetNodeGroups(nodeType)
		if len(groups) > 0 {
			for _, group := range groups {
				if group.GPU != nil || podTemplateRequestsGPU(group.PodTemplate) {
					return true
				}
			}
			continue
		}
		nodeSpec, dataNodeSpec := m.GetNodeSpec(nodeType)
		switch {
		case nodeSpec != nil:
			if nodeSpec.GPU != nil || podTemplateRequestsGPU(nodeSpec.PodTemplate) {
				return true
			}
		case dataNodeSpec != nil:
			if dataNodeSpec.GPU != nil || podTemplateRequestsGPU(dataNodeSpec.PodTemplate) {
				return true
			}
		}
	}
	return false
}

// DistributedNodeRolesWithSRIOV returns the Distributed roles whose
// .network.sriov is set. Used by the admission webhook's topology-coupling
// check (warn if some but not all five roles have it).
func (m *Milvus) DistributedNodeRolesWithSRIOV() []MilvusNodeRoleType {
	if !m.IsDistributed() || m.Spec.Topology.Distributed == nil {
		return nil
	}
	var withSRIOV []MilvusNodeRoleType
	for _, nodeType := range []MilvusNodeRoleType{
		MilvusNodeRoleMixCoord, MilvusNodeRoleDataNode, MilvusNodeRoleProxy,
		MilvusNodeRoleQueryNode, MilvusNodeRoleStreamingNode,
	} {
		has := false
		// When Groups is set, the role-level Network field is ignored for
		// pod-building -- evaluating it here too could mark the role as
		// SR-IOV-enabled (and trigger the partial-topology warning below)
		// for a stale setting no rendered pod actually uses. A role counts
		// as "has SR-IOV" if any single group does -- this feeds the
		// webhook's topology-coupling warning (§ 10), which is necessarily
		// approximate once groups exist: it can't know whether an
		// ungrouped peer role needs to reach every group or just some.
		groups := m.GetNodeGroups(nodeType)
		if len(groups) > 0 {
			for _, group := range groups {
				if group.Network != nil && group.Network.SRIOV != nil {
					has = true
					break
				}
			}
		} else if net := m.GetNodeNetworkSpec(nodeType); net != nil && net.SRIOV != nil {
			has = true
		}
		if has {
			withSRIOV = append(withSRIOV, nodeType)
		}
	}
	return withSRIOV
}

func (m *Milvus) PodControllerLabels(nodeType MilvusNodeRoleType, extraLabels ...map[string]string) map[string]string {
	nodeSpec, dataNodeSpec := m.GetNodeSpec(nodeType)
	var labels map[string]string
	if nodeSpec != nil {
		labels = nodeSpec.PodTemplate.Controller.Labels
	} else {
		labels = dataNodeSpec.PodTemplate.Controller.Labels
	}
	return m.OffshootLabel(meta_util.OverwriteKeys(m.OffshootSelectors(), extraLabels...), labels)
}

// PodControllerLabelsForGroup is PodControllerLabels, plus a
// kubedb.com/node-group label distinguishing this group's pods from every
// other group's (and the ungrouped default's) pods of the same role -- each
// group's PetSet needs a selector that matches only its own pods, or two
// PetSets end up adopting each other's. groupTemplate is the specific
// MilvusNodeGroup's own PodTemplate (its Controller.Labels are honored the
// same way the role-level PodTemplate's are above); pass a zero value if the
// group has none.
func (m *Milvus) PodControllerLabelsForGroup(nodeType MilvusNodeRoleType, groupName string, groupTemplate *ofstv2.PodTemplateSpec, extraLabels ...map[string]string) map[string]string {
	if groupName == "" {
		return m.PodControllerLabels(nodeType, extraLabels...)
	}
	var templateLabels map[string]string
	if groupTemplate != nil {
		templateLabels = groupTemplate.Controller.Labels
	}
	groupLabel := map[string]string{kubedb.LabelNodeGroup: groupName}
	selectors := meta_util.OverwriteKeys(m.OffshootSelectors(), append([]map[string]string{groupLabel}, extraLabels...)...)
	return m.OffshootLabel(selectors, templateLabels)
}

// GroupSelectors is the label selector for one MilvusNodeGroup's PetSet --
// OffshootSelectors() (role-agnostic base) plus the role and group-name
// labels, matching exactly the labels PodControllerLabelsForGroup stamps.
func (m *Milvus) GroupSelectors(nodeType MilvusNodeRoleType, groupName string) map[string]string {
	roleLabels := map[string]string{kubedb.LabelRole: string(nodeType)}
	if groupName == "" {
		return m.OffshootSelectors(roleLabels)
	}
	return m.OffshootSelectors(roleLabels, map[string]string{kubedb.LabelNodeGroup: groupName})
}

func (m *Milvus) ServiceAccountName() string {
	return m.OffshootName()
}

func (m *Milvus) MilvusNodeContainerPort(nodeRole MilvusNodeRoleType) int32 {
	switch nodeRole {
	case MilvusNodeRoleMixCoord:
		return kubedb.MilvusMetricsPort
	case MilvusNodeRoleDataNode:
		return kubedb.MilvusPortDataNode
	case MilvusNodeRoleProxy:
		return kubedb.MilvusGrpcPort
	case MilvusNodeRoleQueryNode:
		return kubedb.MilvusPortQueryNode
	case MilvusNodeRoleStreamingNode:
		return kubedb.MilvusPortStreamingNode
	default:
		klog.Errorf("unknown Milvus node role %s\n", nodeRole)
		return -1
	}
}

func (m *Milvus) GetAuthSecretName() string {
	if m.Spec.AuthSecret != nil && m.Spec.AuthSecret.Name != "" {
		return m.Spec.AuthSecret.Name
	}
	return meta_util.NameWithSuffix(m.OffshootName(), "auth")
}

func (m *Milvus) ConfigSecretName() string {
	uid := string(m.UID)
	return meta_util.NameWithSuffix(m.OffshootName(), uid[len(uid)-6:])
}

func (m *Milvus) GetPersistentSecrets() []string {
	var secrets []string
	if !IsVirtualAuthSecretReferred(m.Spec.AuthSecret) && m.Spec.AuthSecret != nil && m.Spec.AuthSecret.Name != "" {
		secrets = append(secrets, m.GetAuthSecretName())
	}
	secrets = append(secrets, m.ConfigSecretName())
	return secrets
}

func (m *Milvus) OffshootLabel(selector, override map[string]string) map[string]string {
	selector[meta_util.ComponentLabelKey] = kubedb.ComponentDatabase
	return meta_util.FilterKeys(kubedb.GroupName, selector, meta_util.OverwriteKeys(nil, m.Labels, override))
}

func (m *Milvus) OffshootLabels() map[string]string {
	return m.OffshootLabel(m.OffshootSelectors(), nil)
}

func (m *Milvus) OffshootSelectors(extraSelectors ...map[string]string) map[string]string {
	selector := map[string]string{
		meta_util.NameLabelKey:      m.ResourceFQN(),
		meta_util.InstanceLabelKey:  m.Name,
		meta_util.ManagedByLabelKey: kubedb.GroupName,
	}
	return meta_util.OverwriteKeys(selector, extraSelectors...)
}

func (m *Milvus) PodLabel(podTemplate *ofstv2.PodTemplateSpec) map[string]string {
	if podTemplate != nil && podTemplate.Labels != nil {
		return m.OffshootLabel(m.OffshootSelectors(), m.Spec.PodTemplate.Labels)
	}
	return m.OffshootLabel(m.OffshootSelectors(), nil)
}

func (m *Milvus) ServiceDNS() string {
	return fmt.Sprintf("%s.%s.svc.cluster.local:%d", m.ServiceName(), m.Namespace, kubedb.MilvusGrpcPort)
}

func (m *Milvus) SetHealthCheckerDefaults() {
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

func (m *Milvus) EtcdServiceName() string {
	return fmt.Sprintf("%s-%s", m.Name, kubedb.EtcdName)
}

// internalMetaEtcd returns a throwaway Etcd value describing the internally
// managed meta-storage etcd cluster this Milvus creates (see
// pkg/controller/dependency.go in the milvus operator). Milvus and Etcd live in
// the same v1alpha2 package, so its own naming/URL helpers are reused here
// instead of duplicating DNS-building logic.
func (m *Milvus) internalMetaEtcd() *Etcd {
	e := &Etcd{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.EtcdServiceName(),
			Namespace: m.Namespace,
		},
	}
	if m.Spec.MetaStorage != nil {
		e.Spec.TLS = m.Spec.MetaStorage.TLS
		e.Spec.AuthSecret = m.Spec.MetaStorage.AuthSecret
	}
	return e
}

// MetaStorageTLSEnabled reports whether the internally-managed meta etcd has
// TLS configured. Always false when meta storage is externally managed.
func (m *Milvus) MetaStorageTLSEnabled() bool {
	return m.Spec.MetaStorage != nil && !m.Spec.MetaStorage.ExternallyManaged && m.Spec.MetaStorage.TLS != nil
}

// MetaStorageClientCertSecretName returns the internally-managed meta etcd's
// client certificate secret name, or "" if TLS is disabled.
func (m *Milvus) MetaStorageClientCertSecretName() string {
	return m.internalMetaEtcd().GetCertSecretName(EtcdClientCert)
}

// MetaStorageAuthSecretName returns the internally-managed meta etcd's root
// auth secret name (BYO name if spec.metaStorage.authSecret.name is set, else
// the conventional "<etcd-name>-auth").
func (m *Milvus) MetaStorageAuthSecretName() string {
	return m.internalMetaEtcd().GetAuthSecretName()
}

func (m *Milvus) MetaStorageEndpoints() []string {
	if m.Spec.MetaStorage.ExternallyManaged {
		if len(m.Spec.MetaStorage.Endpoints) == 0 {
			klog.Errorf("metadata storage is externally managed but no endpoints were provided")
			return []string{}
		}
		return m.Spec.MetaStorage.Endpoints
	}

	e := m.internalMetaEtcd()
	size := m.Spec.MetaStorage.Size

	endpoints := make([]string, size)
	for i := range size {
		endpoints[i] = e.ClientURL(e.PodName(i))
	}

	return endpoints
}

func (m *Milvus) SetDefaultStorage() *core.PersistentVolumeClaimSpec {
	return &core.PersistentVolumeClaimSpec{
		AccessModes: []core.PersistentVolumeAccessMode{
			core.ReadWriteOnce,
		},
		Resources: core.VolumeResourceRequirements{
			Requests: core.ResourceList{
				core.ResourceStorage: resource.MustParse("1Gi"),
			},
		},
	}
}

func (m *Milvus) setDistributedDefaults(kc client.Client) {
	if m.Spec.Topology.Distributed == nil {
		m.Spec.Topology.Distributed = &MilvusDistributedSpec{}
	}

	var mvVersion catalog.MilvusVersion
	err := kc.Get(context.TODO(), types.NamespacedName{
		Name: m.Spec.Version,
	}, &mvVersion)
	if err != nil {
		return
	}
	m.setComponentDefaults(&mvVersion, &m.Spec.Topology.Distributed.MixCoord)
	m.setComponentDefaults(&mvVersion, &m.Spec.Topology.Distributed.DataNode)
	m.setComponentDefaults(&mvVersion, &m.Spec.Topology.Distributed.Proxy)
	m.setComponentDefaults(&mvVersion, &m.Spec.Topology.Distributed.QueryNode)
	m.setComponentDefaults(&mvVersion, &m.Spec.Topology.Distributed.StreamingNode)
}

func (m *Milvus) setComponentDefaults(mvVersion *catalog.MilvusVersion, node any) {
	var replicas **int32
	var podTemplate **ofstv2.PodTemplateSpec
	var gpu **MilvusGPUSpec
	var network **MilvusNetworkSpec

	switch n := node.(type) {
	case **MilvusNode:
		if *n == nil {
			*n = &MilvusNode{}
		}
		replicas = &(*n).Replicas
		podTemplate = &(*n).PodTemplate
		gpu = &(*n).GPU
		network = &(*n).Network

	case **MilvusDataNode:
		if *n == nil {
			*n = &MilvusDataNode{}
			(*n).Storage = m.SetDefaultStorage()
		}
		replicas = &(*n).Replicas
		podTemplate = &(*n).PodTemplate
		gpu = &(*n).GPU
		network = &(*n).Network
		if (*n).StorageType == "" {
			(*n).StorageType = StorageTypeDurable
		}
	}

	if replicas != nil && *replicas == nil {
		*replicas = pointer.Int32P(1)
	}
	if podTemplate != nil && *podTemplate == nil {
		*podTemplate = &ofstv2.PodTemplateSpec{}
	}

	m.setDefaultContainerSecurityContext(mvVersion, *podTemplate)
	m.setDefaultContainerResourceLimits(*podTemplate)
	apis.SetDefaultResizePolicy((*podTemplate).Spec.Containers, (*podTemplate).Spec.InitContainers)
	setDefaultGPUAndNetwork(*gpu, *network)

	// Groups (if any) get the same replicas/podTemplate/GPU/network
	// defaulting as the role-level fields above, independently per group.
	var groups *[]MilvusNodeGroup
	switch n := node.(type) {
	case **MilvusNode:
		groups = &(*n).Groups
	case **MilvusDataNode:
		groups = &(*n).Groups
	}
	if groups != nil {
		for i := range *groups {
			g := &(*groups)[i]
			if g.Replicas == nil {
				g.Replicas = pointer.Int32P(1)
			}
			if g.PodTemplate == nil {
				g.PodTemplate = &ofstv2.PodTemplateSpec{}
			}
			m.setDefaultContainerSecurityContext(mvVersion, g.PodTemplate)
			m.setDefaultContainerResourceLimits(g.PodTemplate)
			apis.SetDefaultResizePolicy(g.PodTemplate.Spec.Containers, g.PodTemplate.Spec.InitContainers)
			setDefaultGPUAndNetwork(g.GPU, g.Network)
		}
	}
}

// setDefaultGPUAndNetwork fills in ResourceName/Interface defaults on an
// already-set GPU/Network spec. A nil gpu/network (the field wasn't set at
// all) is left nil -- this only defaults sub-fields of a spec the user (or
// an earlier defaulting pass) already opted into.
func setDefaultGPUAndNetwork(gpu *MilvusGPUSpec, network *MilvusNetworkSpec) {
	if gpu != nil {
		if gpu.ResourceName == "" {
			gpu.ResourceName = string(MilvusGPUDefaultResourceName)
		}
		// Count is a plain int64 (not a pointer), so an explicit 0 and "never
		// set" are indistinguishable on the wire; either way, requesting a
		// GPU with a count of 0 is meaningless, so default it to 1 rather
		// than trying to reject it at admission.
		if gpu.Count <= 0 {
			gpu.Count = 1
		}
	}
	if network != nil && network.SRIOV != nil && network.SRIOV.Interface == "" {
		network.SRIOV.Interface = MilvusSRIOVDefaultInterface
	}
}

// IsLegacyRocksMQ reports whether this is a Standalone Milvus created before the Woodpecker WAL.
func (m *Milvus) IsLegacyRocksMQ() bool {
	if m.IsDistributed() {
		return false
	}
	return m.GetAnnotations()[kubedb.MilvusWALAnnotation] == kubedb.MilvusWALRocksMQ
}

func (m *Milvus) SidekickName() string {
	return m.Name + "-" + kubedb.MilvusSidekickSuffix
}

func (m *Milvus) SidekickLabels(skName string) map[string]string {
	return meta_util.OverwriteKeys(nil, kubedb.CommonSidekickLabels(), map[string]string{
		meta_util.ComponentLabelKey: skName,
		kubedb.SidekickOwnerName:    m.Name,
		kubedb.SidekickOwnerKind:    m.ResourceFQN(),
	})
}

func (m *Milvus) ArchiverBackupConfigName() string {
	return m.Name + "-" + kubedb.MilvusArchiverBackupConfigSuffix
}

func (m *Milvus) IncrementalSnapshotName() string {
	return m.Name + "-" + kubedb.MilvusIncrementalSnapshotSuffix
}

func (m *Milvus) SetDefaults(kc client.Client) {
	if m.Spec.DeletionPolicy == "" {
		m.Spec.DeletionPolicy = DeletionPolicyDelete
	}

	var mvVersion catalog.MilvusVersion
	err := kc.Get(context.TODO(), types.NamespacedName{
		Name: m.Spec.Version,
	}, &mvVersion)
	if err != nil {
		return
	}

	if m.Spec.DeletionPolicy == "" {
		m.Spec.DeletionPolicy = DeletionPolicyDelete
	}

	if m.Spec.StorageType == "" {
		m.Spec.StorageType = StorageTypeDurable
	}

	if m.Spec.AuthSecret == nil {
		m.Spec.AuthSecret = &SecretReference{}
	}

	if m.Spec.AuthSecret.Kind == "" {
		m.Spec.AuthSecret.Kind = kubedb.ResourceKindSecret
	}

	if m.Spec.PodTemplate == nil {
		m.Spec.PodTemplate = &ofstv2.PodTemplateSpec{}
	}

	if m.IsDistributed() {
		m.setDistributedDefaults(kc)
	} else {
		m.setDefaultContainerSecurityContext(&mvVersion, m.Spec.PodTemplate)
		m.setDefaultContainerResourceLimits(m.Spec.PodTemplate)
		apis.SetDefaultResizePolicy(m.Spec.PodTemplate.Spec.Containers, m.Spec.PodTemplate.Spec.InitContainers)
		setDefaultGPUAndNetwork(m.Spec.GPU, m.Spec.Network)
	}

	m.setMetaStorageDefaults()

	m.setArchiverDefaults()

	m.SetHealthCheckerDefaults()

	m.SetTLSDefaults()

	if m.Spec.Monitor != nil {
		if m.Spec.Monitor.Prometheus == nil {
			m.Spec.Monitor.Prometheus = &mona.PrometheusSpec{}
		}
		if m.Spec.Monitor.Prometheus.Exporter.Port == 0 {
			m.Spec.Monitor.Prometheus.Exporter.Port = kubedb.MilvusMetricsPort
		}
		m.Spec.Monitor.SetDefaults()
		if m.Spec.Monitor.Prometheus.Exporter.SecurityContext.RunAsUser == nil {
			m.Spec.Monitor.Prometheus.Exporter.SecurityContext.RunAsUser = mvVersion.Spec.SecurityContext.RunAsUser
		}
		if m.Spec.Monitor.Prometheus.Exporter.SecurityContext.RunAsGroup == nil {
			m.Spec.Monitor.Prometheus.Exporter.SecurityContext.RunAsGroup = mvVersion.Spec.SecurityContext.RunAsUser
		}
	}
}

func (m *Milvus) SetTLSDefaults() {
	if m.Spec.TLS == nil || m.Spec.TLS.IssuerRef == nil {
		return
	}

	if m.Spec.TLS.External == nil {
		m.Spec.TLS.External = &ProtocolTLSConfig{
			Mode: TLSModeDisabled,
		}
	}

	if m.Spec.TLS.Internal == nil {
		m.Spec.TLS.Internal = &ProtocolTLSConfig{
			Mode: TLSModeDisabled,
		}
	}

	m.Spec.TLS.Certificates = kmapi.SetMissingSecretNameForCertificate(m.Spec.TLS.Certificates, string(MilvusCertificateTypeServer), m.CertificateName(MilvusCertificateTypeServer))
	m.Spec.TLS.Certificates = kmapi.SetMissingSecretNameForCertificate(m.Spec.TLS.Certificates, string(MilvusCertificateTypeClient), m.CertificateName(MilvusCertificateTypeClient))
}

func (m *Milvus) setMetaStorageDefaults() {
	if m.Spec.MetaStorage == nil {
		m.Spec.MetaStorage = &MetaStorageSpec{}
	}

	if m.Spec.MetaStorage.StorageType == "" {
		m.Spec.MetaStorage.StorageType = StorageTypeDurable
	}

	if !m.Spec.MetaStorage.ExternallyManaged {
		if m.Spec.MetaStorage.Size == 0 {
			m.Spec.MetaStorage.Size = 3
		}

		if m.Spec.MetaStorage.Storage == nil {
			m.Spec.MetaStorage.Storage = &core.PersistentVolumeClaimSpec{}
		}

		if len(m.Spec.MetaStorage.Storage.AccessModes) == 0 {
			m.Spec.MetaStorage.Storage.AccessModes = []core.PersistentVolumeAccessMode{
				core.ReadWriteOnce,
			}
		}

		if m.Spec.MetaStorage.Storage.Resources.Requests == nil {
			m.Spec.MetaStorage.Storage.Resources.Requests = core.ResourceList{
				core.ResourceStorage: resource.MustParse("1Gi"),
			}
		}
	}
}

func (m *Milvus) setDefaultContainerSecurityContext(mvVersion *catalog.MilvusVersion, podTemplate *ofstv2.PodTemplateSpec) {
	if podTemplate == nil {
		return
	}
	if podTemplate.Spec.SecurityContext == nil {
		podTemplate.Spec.SecurityContext = &core.PodSecurityContext{}
	}
	if podTemplate.Spec.SecurityContext.FSGroup == nil {
		podTemplate.Spec.SecurityContext.FSGroup = mvVersion.Spec.SecurityContext.RunAsUser
	}

	container := coreutil.GetContainerByName(podTemplate.Spec.Containers, kubedb.MilvusContainerName)
	if container == nil {
		container = &core.Container{
			Name: kubedb.MilvusContainerName,
		}
	}
	if container.SecurityContext == nil {
		container.SecurityContext = &core.SecurityContext{}
	}
	m.AssignDefaultContainerSecurityContext(mvVersion, container.SecurityContext)
	podTemplate.Spec.Containers = coreutil.UpsertContainer(podTemplate.Spec.Containers, *container)
}

func (m *Milvus) AssignDefaultContainerSecurityContext(mvVersion *catalog.MilvusVersion, rc *core.SecurityContext) {
	if rc.AllowPrivilegeEscalation == nil {
		rc.AllowPrivilegeEscalation = pointer.BoolP(false)
	}
	if rc.Capabilities == nil {
		rc.Capabilities = &core.Capabilities{
			Drop: []core.Capability{"ALL"},
		}
	}
	if rc.RunAsNonRoot == nil {
		rc.RunAsNonRoot = pointer.BoolP(true)
	}
	if rc.RunAsUser == nil {
		rc.RunAsUser = mvVersion.Spec.SecurityContext.RunAsUser
	}
	if rc.RunAsGroup == nil {
		rc.RunAsGroup = mvVersion.Spec.SecurityContext.RunAsUser
	}
	if rc.SeccompProfile == nil {
		rc.SeccompProfile = secomp.DefaultSeccompProfile()
	}
}

func (m *Milvus) setDefaultContainerResourceLimits(podTemplate *ofstv2.PodTemplateSpec) {
	dbContainer := coreutil.GetContainerByName(podTemplate.Spec.Containers, kubedb.MilvusContainerName)
	if dbContainer != nil {
		apis.SetDefaultResourceLimits(&dbContainer.Resources, kubedb.DefaultResources)
	}
}

func (m *Milvus) IsDistributed() bool {
	return m != nil &&
		m.Spec.Topology != nil &&
		m.Spec.Topology.Mode != nil &&
		*m.Spec.Topology.Mode == "Distributed"
}

func (m *Milvus) ServiceLabels(alias ServiceAlias, extraLabels ...map[string]string) map[string]string {
	svcTemplate := GetServiceTemplate(m.Spec.ServiceTemplates, alias)
	return m.OffshootLabel(meta_util.OverwriteKeys(m.OffshootSelectors(), extraLabels...), svcTemplate.Labels)
}

type milvusStatsService struct {
	*Milvus
}

func (m milvusStatsService) GetNamespace() string {
	return m.Milvus.GetNamespace()
}

func (m milvusStatsService) ServiceName() string {
	return m.OffshootName() + "-stats"
}

func (m milvusStatsService) ServiceMonitorName() string {
	return m.ServiceName()
}

func (m milvusStatsService) ServiceMonitorAdditionalLabels() map[string]string {
	return m.OffshootLabels()
}

func (m milvusStatsService) Path() string {
	return kubedb.DefaultStatsPath
}

func (m milvusStatsService) Scheme() string {
	return "http"
}

func (m Milvus) StatsServiceLabels() map[string]string {
	return m.ServiceLabels(StatsServiceAlias, map[string]string{kubedb.LabelRole: kubedb.RoleStats})
}

func (m milvusStatsService) TLSConfig() *promapi.TLSConfig {
	return nil
}

func (m Milvus) StatsService() mona.StatsAccessor {
	return &milvusStatsService{&m}
}

func (m *Milvus) GetCertSecretName(alias MilvusCertificateType) string {
	if m.Spec.TLS != nil {
		name, ok := kmapi.GetCertificateSecretName(m.Spec.TLS.Certificates, string(alias))
		if ok {
			return name
		}
	}
	return m.CertificateName(alias)
}

func (m *Milvus) CertificateName(alias MilvusCertificateType) string {
	return meta_util.NameWithSuffix(m.Name, fmt.Sprintf("%s-cert", string(alias)))
}

func (m *Milvus) GetStorageClassName() string {
	return *m.Spec.Storage.StorageClassName
}

type MilvusBind struct {
	*Milvus
}

var _ DBBindInterface = &MilvusBind{}

func (m *MilvusBind) ServiceNames() (string, string) {
	return m.ServiceName(), ""
}

func (m *MilvusBind) Ports() (int, int) {
	return int(kubedb.MilvusHttpPort), 0
}

func (m *MilvusBind) SecretName() string {
	return m.GetAuthSecretName()
}

func (m *MilvusBind) CertSecretName() string {
	return m.GetCertSecretName(MilvusCertificateTypeClient)
}

func (m *Milvus) GetDeletionPolicy() string {
	return string(m.Spec.DeletionPolicy)
}

func (m *Milvus) AsOwner() *metav1.OwnerReference {
	return metav1.NewControllerRef(m, SchemeGroupVersion.WithKind(m.ResourceKind()))
}

func (m *Milvus) setArchiverDefaults() {
	if m.Spec.Init != nil && m.Spec.Init.Archiver != nil {
		a := m.Spec.Init.Archiver
		if a.EncryptionSecret != nil && a.EncryptionSecret.Namespace == "" {
			a.EncryptionSecret.Namespace = m.GetNamespace()
		}
		if a.FullDBRepository != nil && a.FullDBRepository.Namespace == "" {
			a.FullDBRepository.Namespace = m.GetNamespace()
		}
		if a.ManifestRepository != nil && a.ManifestRepository.Namespace == "" {
			a.ManifestRepository.Namespace = m.GetNamespace()
		}
	}
}
