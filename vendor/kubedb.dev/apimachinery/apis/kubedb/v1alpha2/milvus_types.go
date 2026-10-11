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
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kmapi "kmodules.xyz/client-go/api/v1"
	mona "kmodules.xyz/monitoring-agent-api/api/v1"
	ofstv2 "kmodules.xyz/offshoot-api/api/v2"
)

const (
	ResourceCodeMilvus     = "mv"
	ResourceKindMilvus     = "Milvus"
	ResourceSingularMilvus = "milvus"
	ResourcePluralMilvus   = "milvuses"
)

// +kubebuilder:validation:Enum=Standalone;Distributed
type MilvusMode string

// Package v1alpha2 contains API Schema definitions for the  v1alpha2 API group.

// +genclient
// +k8s:openapi-gen=true

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=milvuses,singular=milvus,shortName=mv,categories={datastore,kubedb,appscode,all}
// +kubebuilder:printcolumn:name="Version",type="string",JSONPath=".spec.version"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type Milvus struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MilvusSpec   `json:"spec,omitempty"`
	Status MilvusStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen=true
// MilvusSpec defines the desired state of Milvus
type MilvusSpec struct {
	// Version of Milvus to be deployed
	Version string `json:"version"`

	// Meta contains configuration for etcd meta storage
	MetaStorage *MetaStorageSpec `json:"metaStorage,omitempty"`

	// ObjectStorage contains specification for druid to connect to the object storage
	ObjectStorage *ObjectStorageSpec `json:"objectStorage"`

	// Milvus cluster topology
	// +optional
	Topology *MilvusTopology `json:"topology,omitempty"`

	// PodTemplate is an optional configuration for pods used to expose database
	// +optional
	PodTemplate *ofstv2.PodTemplateSpec `json:"podTemplate,omitempty"`

	// StorageType can be durable (default) or ephemeral
	StorageType StorageType `json:"storageType,omitempty"`

	// Storage to specify how storage shall be used.
	Storage *core.PersistentVolumeClaimSpec `json:"storage,omitempty"`

	// disable security. It disables authentication security of user.
	// If unset, default is false
	// +optional
	DisableSecurity bool `json:"disableSecurity,omitempty"`

	// Database authentication secret
	// +optional
	AuthSecret *SecretReference `json:"authSecret,omitempty"`

	// +optional
	Configuration *ConfigurationSpec `json:"configuration,omitempty"`

	// ServiceTemplates is an optional configuration for services used to expose database
	// +optional
	ServiceTemplates []NamedServiceTemplateSpec `json:"serviceTemplates,omitempty"`

	// DeletionPolicy controls the delete operation for database
	// +optional
	DeletionPolicy DeletionPolicy `json:"deletionPolicy,omitempty"`

	// Indicates that the database is halted and all offshoot Kubernetes resources except PVCs are deleted.
	// +optional
	Halted bool `json:"halted,omitempty"`

	// HealthChecker defines attributes of the health checker
	// +optional
	// +kubebuilder:default={periodSeconds: 10, timeoutSeconds: 10, failureThreshold: 3}
	HealthChecker kmapi.HealthCheckSpec `json:"healthChecker"`

	// Monitor is used monitor database instance
	// +optional
	Monitor *mona.AgentSpec `json:"monitor,omitempty"`

	// TLS contains tls configurations
	// +optional
	TLS *MilvusTLSConfig `json:"tls,omitempty"`

	// Network configures secondary networking (e.g. SR-IOV via Multus) for
	// this component. Only meaningful in Standalone mode; Distributed roles
	// configure this per-role under spec.topology.distributed.<role>.network.
	// +optional
	Network *MilvusNetworkSpec `json:"network,omitempty"`

	// GPU configures GPU device scheduling for this component. Only
	// meaningful in Standalone mode; Distributed roles configure this
	// per-role under spec.topology.distributed.<role>.gpu.
	// +optional
	GPU *MilvusGPUSpec `json:"gpu,omitempty"`

	// Init is used to initialize the database from an archiver: a full backup
	// (and, when the archiver records a change log, the log up to a recovery
	// timestamp) is restored into this new Milvus before it starts.
	// +optional
	Init *InitSpec `json:"init,omitempty"`

	// Archiver controls the continuous archiver (base backups and point-in-time
	// recovery) of this database. It references a MilvusArchiver that also
	// selects this database (double opt-in).
	// +optional
	Archiver *Archiver `json:"archiver,omitempty"`
}

// +k8s:deepcopy-gen=true
// MilvusNetworkSpec configures secondary networking for a Milvus component.
type MilvusNetworkSpec struct {
	// SRIOV attaches an SR-IOV virtual function to the pod via Multus CNI.
	// +optional
	SRIOV *MilvusSRIOVSpec `json:"sriov,omitempty"`
}

// +k8s:deepcopy-gen=true
// MilvusSRIOVSpec requests a Multus/SR-IOV secondary network attachment.
type MilvusSRIOVSpec struct {
	// AttachmentRef names a cluster-admin-authored NetworkAttachmentDefinition
	// (k8s.cni.cncf.io/v1) in the same namespace as this Milvus.
	// +kubebuilder:validation:MinLength=1
	AttachmentRef string `json:"attachmentRef"`

	// ResourceName must match the SR-IOV device plugin's advertised extended
	// resource for the requested VF (e.g. intel.com/sriov_net_A,
	// nvidia.com/hostdev).
	// +kubebuilder:validation:MinLength=1
	ResourceName string `json:"resourceName"`

	// RDMAResourceName optionally requests an additional RDMA device resource
	// (e.g. rdma/rdma_shared_device_a) alongside ResourceName, for GPUDirect
	// RDMA setups. Most deployments do not need this.
	// +optional
	RDMAResourceName string `json:"rdmaResourceName,omitempty"`

	// Interface names the secondary interface Multus should attach this
	// network as. Defaults to Multus's own convention ("net1") when unset.
	// Setting a non-default value requires the rendered
	// k8s.v1.cni.cncf.io/networks annotation to use the JSON-array form
	// naming this interface explicitly, not just AttachmentRef by itself.
	// +optional
	Interface string `json:"interface,omitempty"`

	// GDS optionally attaches a second, independent SR-IOV network for GPU
	// Direct Storage to object storage -- confirmed as a real requirement
	// on data-plane roles (QueryNode, DataNode, StreamingNode) by a
	// customer reference implementation
	// (pl2/milvus-fixes/sriov-gpu-design.md § 13.4); Proxy/MixCoord only
	// need the RDMA attachment above. When set, the rendered
	// k8s.v1.cni.cncf.io/networks annotation always uses the JSON-array
	// form (two entries), regardless of whether Interface above is at its
	// default.
	// +optional
	GDS *MilvusSRIOVAttachmentSpec `json:"gds,omitempty"`
}

// +k8s:deepcopy-gen=true
// MilvusSRIOVAttachmentSpec requests one additional Multus/SR-IOV secondary
// network attachment, alongside the primary one in MilvusSRIOVSpec.
type MilvusSRIOVAttachmentSpec struct {
	// AttachmentRef names a cluster-admin-authored NetworkAttachmentDefinition
	// (k8s.cni.cncf.io/v1) in the same namespace as this Milvus.
	// +kubebuilder:validation:MinLength=1
	AttachmentRef string `json:"attachmentRef"`

	// ResourceName must match the SR-IOV device plugin's advertised extended
	// resource for the requested VF.
	// +kubebuilder:validation:MinLength=1
	ResourceName string `json:"resourceName"`

	// Interface names the secondary interface Multus should attach this
	// network as. Defaults to "net2" when unset (net1 is MilvusSRIOVSpec's
	// own default interface, above).
	// +optional
	Interface string `json:"interface,omitempty"`
}

// +k8s:deepcopy-gen=true
// MilvusGPUSpec configures GPU device scheduling for a Milvus component.
type MilvusGPUSpec struct {
	// ResourceName is the extended resource name to request (default
	// nvidia.com/gpu). Override only for vGPU/MIG resource names.
	// +optional
	ResourceName string `json:"resourceName,omitempty"`

	// Count is the number of GPU devices to request. Must be a positive
	// integer if set.
	// +optional
	Count int64 `json:"count,omitempty"`

	// NodeSelector is merged into the pod's node selector.
	// +optional
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`

	// Tolerations is merged into the pod's tolerations.
	// +optional
	Tolerations []core.Toleration `json:"tolerations,omitempty"`
}

type MilvusTLSConfig struct {
	kmapi.TLSConfig `json:",inline"`

	// External controls TLS for client-facing traffic (gRPC + REST).
	// +optional
	External *ProtocolTLSConfig `json:"external,omitempty"`

	// Internal enables TLS for inter-component communication (one-way only)
	// +optional
	Internal *ProtocolTLSConfig `json:"internal,omitempty"`
}

// +kubebuilder:validation:Enum=server;client
type MilvusCertificateType string

const (
	MilvusCertificateTypeServer MilvusCertificateType = "server"
	MilvusCertificateTypeClient MilvusCertificateType = "client"
)

type MilvusTopology struct {
	// If set to -
	// "Standalone", Milvus will start a Standalone Mode
	// "Distributed", Milvus will start a Distributed Mode
	Mode *MilvusMode `json:"mode,omitempty"`

	// Distributed contains information of the Distributed configuration.
	// Used when Mode is "Distributed".
	// +optional
	Distributed *MilvusDistributedSpec `json:"distributed,omitempty"`
}

type MilvusDistributedSpec struct {
	// +optional
	DataNode *MilvusNode `json:"datanode,omitempty"`

	// +optional
	MixCoord *MilvusNode `json:"mixcoord,omitempty"`

	// +optional
	QueryNode *MilvusNode `json:"querynode,omitempty"`

	// +optional
	StreamingNode *MilvusDataNode `json:"streamingnode,omitempty"`

	// +optional
	Proxy *MilvusNode `json:"proxy,omitempty"`
}

type MilvusNode struct {
	// Replicas represents number of replicas for the specific type of node
	// +optional
	Replicas *int32 `json:"replicas,omitempty"`

	// PodTemplate is an optional configuration for pods used to expose database
	// +optional
	PodTemplate *ofstv2.PodTemplateSpec `json:"podTemplate,omitempty"`

	// Network configures secondary networking (e.g. SR-IOV via Multus) for
	// this component. Every Distributed role dials, or is dialed by, at
	// least one other role directly by its advertised address, so setting
	// this on some roles but not others is very likely a connectivity bug,
	// not a deliberate choice; the admission webhook warns accordingly.
	// +optional
	Network *MilvusNetworkSpec `json:"network,omitempty"`

	// GPU configures GPU device scheduling for this component. Meaningful
	// on QueryNode (search) and DataNode (index build); harmless if set on
	// other roles.
	// +optional
	GPU *MilvusGPUSpec `json:"gpu,omitempty"`

	// Groups optionally splits this role into multiple named,
	// independently scheduled sub-pools -- e.g. two GPU classes of
	// QueryNode in one cluster, each with its own network/gpu/podTemplate
	// (pl2/milvus-fixes/sriov-gpu-design.md § 13.5). When set, the
	// operator creates one PetSet per group, named <db>-<role>-<group
	// name>, instead of one PetSet for the whole role, and the
	// Replicas/PodTemplate/Network/GPU fields above are ignored for
	// pod-building purposes (each group carries its own). When unset (the
	// default), behavior is unchanged: one PetSet named <db>-<role>, built
	// from the fields above.
	//
	// Group names must be unique within this role.
	// +optional
	Groups []MilvusNodeGroup `json:"groups,omitempty"`
}

// +k8s:deepcopy-gen=true
// MilvusNodeGroup is one named sub-pool within a Distributed role's
// MilvusNode.Groups.
type MilvusNodeGroup struct {
	// Name distinguishes this group from others in the same role. Must be
	// a valid DNS label segment; used in the PetSet/pod names as
	// <db>-<role>-<name>.
	Name string `json:"name"`

	// Replicas represents number of replicas for this group.
	// +optional
	Replicas *int32 `json:"replicas,omitempty"`

	// PodTemplate is an optional configuration for pods in this group.
	// +optional
	PodTemplate *ofstv2.PodTemplateSpec `json:"podTemplate,omitempty"`

	// Network configures secondary networking for this group.
	// +optional
	Network *MilvusNetworkSpec `json:"network,omitempty"`

	// GPU configures GPU device scheduling for this group.
	// +optional
	GPU *MilvusGPUSpec `json:"gpu,omitempty"`

	// StorageType and Storage optionally override this group's own storage.
	// Meaningful only for StreamingNode groups (spec.topology.distributed.
	// streamingnode.groups[]) -- the only Distributed role with storage at
	// all; harmless no-ops on every other role's groups. When unset, a
	// StreamingNode group falls back to MilvusDataNode's top-level
	// StorageType/Storage (the pre-Groups default).
	// +optional
	StorageType StorageType `json:"storageType,omitempty"`
	// +optional
	Storage *core.PersistentVolumeClaimSpec `json:"storage,omitempty"`
}

type MilvusDataNode struct {
	// MilvusDataNode has all the characteristics of MilvusNode
	MilvusNode `json:",inline"`

	// StorageType specifies if the storage
	// of this node is durable (default) or ephemeral.
	StorageType StorageType `json:"storageType,omitempty"`

	// Storage to specify how storage shall be used.
	Storage *core.PersistentVolumeClaimSpec `json:"storage,omitempty"`
}

// +k8s:deepcopy-gen=true
// Meta Storage defines the configuration for etcd meta storage
type MetaStorageSpec struct {
	// ExternallyManaged indicates whether etcd is managed outside this operator.
	// If true, only endpoints are used. If false, an EtcdCluster CR is created.
	// +optional
	ExternallyManaged bool `json:"externallyManaged,omitempty"`

	// Endpoints are the client endpoints of etcd (e.g., ["http://etcd-svc:2379"]).
	// Required when ExternallyManaged=true.
	// +kubebuilder:validation:MinItems=1
	// +optional
	Endpoints []string `json:"endpoints,omitempty"`

	// Size is the expected size of the cluster.
	// Required when ExternallyManaged=false. Ignored otherwise.
	// +kubebuilder:validation:Minimum=1
	// +optional
	Size int `json:"size,omitempty"`

	// StorageType can be durable (default) or ephemeral
	// +optional
	StorageType StorageType `json:"storageType,omitempty"`

	// Storage to specify how storage shall be used.
	// +optional
	Storage *core.PersistentVolumeClaimSpec `json:"storage,omitempty"`

	// TLS configures certificates issued from spec.metaStorage.tls.issuerRef for
	// the internally-managed meta etcd's client, server and peer traffic.
	// Only used when ExternallyManaged is false.
	// +optional
	TLS *kmapi.TLSConfig `json:"tls,omitempty"`

	// AuthSecret is the root credential for the internally-managed meta etcd.
	// If omitted, the etcd operator auto-generates one. Only used when
	// ExternallyManaged is false.
	// +optional
	AuthSecret *SecretReference `json:"authSecret,omitempty"`
}

// +k8s:deepcopy-gen=true
// ObjectStorageStorageSpec defines the configuration for MinIO or S3 object storage
type ObjectStorageSpec struct {
	// ConfigSecret should contain the necessary data to connect to external MinIO
	// +optional
	ConfigSecret *core.LocalObjectReference `json:"configSecret,omitempty"`
}

// +k8s:deepcopy-gen=true
// MilvusStatus defines the observed state of Milvus
type MilvusStatus struct {
	// Specifies the current phase of the database
	// +optional
	Phase DatabasePhase `json:"phase,omitempty"`

	// ObservedGeneration is the most recent generation observed for this resource
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// Conditions applied to the database, such as approval or denial.
	// +optional
	Conditions []kmapi.Condition `json:"conditions,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// MilvusList contains a list of Milvus
type MilvusList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Milvus `json:"items"`
}
type MilvusNodeRoleType string

const (
	MilvusNodeRoleDataNode      MilvusNodeRoleType = "datanode"
	MilvusNodeRoleMixCoord      MilvusNodeRoleType = "mixcoord"
	MilvusNodeRoleQueryNode     MilvusNodeRoleType = "querynode"
	MilvusNodeRoleStreamingNode MilvusNodeRoleType = "streamingnode"
	MilvusNodeRoleProxy         MilvusNodeRoleType = "proxy"
)

// MilvusGPUDefaultResourceName is the extended resource requested when
// MilvusGPUSpec.ResourceName is left unset.
const MilvusGPUDefaultResourceName core.ResourceName = "nvidia.com/gpu"

// MilvusSRIOVDefaultInterface is the secondary interface name Multus attaches
// as when MilvusSRIOVSpec.Interface is left unset.
const MilvusSRIOVDefaultInterface = "net1"

var _ Accessor = &Milvus{}

func (m *Milvus) GetObjectMeta() metav1.ObjectMeta {
	return m.ObjectMeta
}

func (m *Milvus) GetConditions() []kmapi.Condition {
	return m.Status.Conditions
}

func (m *Milvus) SetCondition(cond kmapi.Condition) {
	m.Status.Conditions = setCondition(m.Status.Conditions, cond)
}

func (m *Milvus) RemoveCondition(typ string) {
	m.Status.Conditions = removeCondition(m.Status.Conditions, typ)
}
