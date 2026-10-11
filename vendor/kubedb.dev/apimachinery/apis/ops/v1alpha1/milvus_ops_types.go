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

//go:generate go-enum --mustparse --names --values
package v1alpha1

import (
	api "kubedb.dev/apimachinery/apis/kubedb/v1alpha2"

	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ResourceCodeMilvusOpsRequest     = "mvops"
	ResourceKindMilvusOpsRequest     = "MilvusOpsRequest"
	ResourceSingularMilvusOpsRequest = "milvusopsrequest"
	ResourcePluralMilvusOpsRequest   = "milvusopsrequests"
)

// MilvusDBOpsRequest defines a Milvus DBA operation.

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=milvusopsrequests,singular=milvusopsrequest,shortName=mvops,categories={ops,kubedb,appscode}
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Type",type="string",JSONPath=".spec.type"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type MilvusOpsRequest struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              MilvusOpsRequestSpec `json:"spec,omitempty"`
	Status            OpsRequestStatus     `json:"status,omitempty"`
}

// MilvusOpsRequestSpec is the spec for MilvusOpsRequest
type MilvusOpsRequestSpec struct {
	// Specifies the Milvus reference
	DatabaseRef core.LocalObjectReference `json:"databaseRef"`
	// Specifies the ops request type: UpdateVersion, HorizontalScaling, VerticalScaling etc.
	Type MilvusOpsRequestType `json:"type"`
	// Specifies information necessary for upgrading milvus
	UpdateVersion *MilvusUpdateVersionSpec `json:"updateVersion,omitempty"`
	// Specifies information necessary for horizontal scaling
	HorizontalScaling *MilvusHorizontalScalingSpec `json:"horizontalScaling,omitempty"`
	// Specifies information necessary for vertical scaling
	VerticalScaling *MilvusVerticalScalingSpec `json:"verticalScaling,omitempty"`
	// Specifies information necessary for volume expansion
	VolumeExpansion *MilvusVolumeExpansionSpec `json:"volumeExpansion,omitempty"`
	// Specifies information necessary for custom configuration of milvus
	Configuration *ReconfigurationSpec `json:"configuration,omitempty"`
	// Specifies information necessary for configuring TLS
	TLS *MilvusTLSSpec `json:"tls,omitempty"`
	// Specifies information necessary for configuring authSecret of the database
	Authentication *AuthSpec `json:"authentication,omitempty"`
	// Specifies information necessary for restarting database
	Restart *RestartSpec `json:"restart,omitempty"`
	// Specifies information necessary for migrating storageClass or data
	Migration *StorageMigrationSpec `json:"migration,omitempty"`
	// Timeout for each step of the ops request in second. If a step doesn't finish within the specified timeout, the ops request will result in failure.
	Timeout *metav1.Duration `json:"timeout,omitempty"`
	// ApplyOption is to control the execution of OpsRequest depending on the database state.
	// +kubebuilder:default="IfReady"
	Apply ApplyOption `json:"apply,omitempty"`
	// +kubebuilder:default=1
	MaxRetries int32 `json:"maxRetries,omitempty"`
}

// +kubebuilder:validation:Enum=UpdateVersion;HorizontalScaling;VerticalScaling;VolumeExpansion;Restart;Reconfigure;ReconfigureTLS;RotateAuth;StorageMigration
// ENUM(UpdateVersion, HorizontalScaling, VerticalScaling, VolumeExpansion, Restart, Reconfigure, ReconfigureTLS, RotateAuth, StorageMigration)
type MilvusOpsRequestType string

// MilvusUpdateVersionSpec contains the update version information of a milvus cluster
type MilvusUpdateVersionSpec struct {
	// Specifies the target version name from catalog
	TargetVersion string `json:"targetVersion,omitempty"`
}

// MilvusReplicaReadinessCriteria is the criteria for checking readiness of a Milvus pod
// after updating, horizontal scaling etc.
type MilvusReplicaReadinessCriteria struct{}

// MilvusHorizontalScalingSpec contains the horizontal scaling information of a Milvus cluster
type MilvusHorizontalScalingSpec struct {
	// Node topology specification
	Topology *MilvusHorizontalScalingTopologySpec `json:"topology,omitempty"`
}

// MilvusHorizontalScalingTopologySpec contains the horizontal scaling information in cluster topology mode
type MilvusHorizontalScalingTopologySpec struct {
	// Number of Proxy nodes
	Proxy *int32 `json:"proxy,omitempty"`
	// Number of MixCoord nodes
	MixCoord *int32 `json:"mixcoord,omitempty"`
	// Number of QueryNode nodes
	QueryNode *int32 `json:"querynode,omitempty"`
	// Number of StreamingNode nodes
	StreamingNode *int32 `json:"streamingnode,omitempty"`
	// Number of DataNode nodes
	DataNode *int32 `json:"dataNode,omitempty"`

	// Groups optionally horizontally scales one or more MilvusNodeGroup
	// sub-pools (pl2/milvus-fixes/sriov-gpu-design.md § 13.5) instead of a
	// whole role above. A role whose spec.topology.distributed.<role>.groups
	// is set is only addressable through this field -- the flat per-role
	// fields above target that role's own (ungrouped) PetSet, which doesn't
	// exist once Groups is in use.
	// +optional
	Groups []MilvusScalingGroupTarget `json:"groups,omitempty"`
}

// MilvusScalingGroupTarget identifies one MilvusNodeGroup (by role + group
// name) and its desired replica count, for horizontal-scaling a group
// individually.
type MilvusScalingGroupTarget struct {
	// NodeType is the Distributed role the target group belongs to.
	NodeType api.MilvusNodeRoleType `json:"nodeType"`
	// Group is the MilvusNodeGroup name (spec.topology.distributed.<nodeType>.groups[].name).
	Group string `json:"group"`
	// Replicas is the desired replica count for this group's PetSet.
	Replicas int32 `json:"replicas"`
}

// MilvusVerticalScalingSpec contains the vertical scaling information of a Milvus cluster
type MilvusVerticalScalingSpec struct {
	// Used when Milvus runs in Standalone mode
	Node *PodResources `json:"node,omitempty"`
	// Used when Milvus runs in Distributed mode
	Proxy         *PodResources `json:"proxy,omitempty"`
	MixCoord      *PodResources `json:"mixcoord,omitempty"`
	DataNode      *PodResources `json:"datanode,omitempty"`
	QueryNode     *PodResources `json:"querynode,omitempty"`
	StreamingNode *PodResources `json:"streamingnode,omitempty"`
	// Mode selects how the vertical scaling is actuated. Defaults to Restart.
	// +optional
	// +kubebuilder:default=Restart
	Mode VerticalScalingMode `json:"mode,omitempty"`

	// Groups optionally vertically scales one or more MilvusNodeGroup
	// sub-pools (pl2/milvus-fixes/sriov-gpu-design.md § 13.5) instead of a
	// whole role above. A role whose spec.topology.distributed.<role>.groups
	// is set is only addressable through this field -- the flat per-role
	// fields above target that role's own (ungrouped) PetSet, which doesn't
	// exist once Groups is in use.
	// +optional
	Groups []MilvusVerticalScalingGroupTarget `json:"groups,omitempty"`
}

// MilvusVerticalScalingGroupTarget identifies one MilvusNodeGroup (by role +
// group name) and its desired resources, for vertical-scaling a group
// individually.
type MilvusVerticalScalingGroupTarget struct {
	// NodeType is the Distributed role the target group belongs to.
	NodeType api.MilvusNodeRoleType `json:"nodeType"`
	// Group is the MilvusNodeGroup name (spec.topology.distributed.<nodeType>.groups[].name).
	Group string `json:"group"`
	// Resources is the desired resource/scheduling spec for this group.
	Resources PodResources `json:"resources"`
}

// MilvusVolumeExpansionSpec is the spec for Milvus volume expansion
type MilvusVolumeExpansionSpec struct {
	Mode VolumeExpansionMode `json:"mode"`
	// volume specification for standalone node
	Node *resource.Quantity `json:"node,omitempty"`
	// volume specification for stremingnode
	StreamingNode *resource.Quantity `json:"streamingnode,omitempty"`
}

type MilvusTLSSpec struct {
	// +optional
	api.MilvusTLSConfig `json:",inline,omitempty"`

	// RotateCertificates tells operator to initiate certificate rotation
	// +optional
	RotateCertificates bool `json:"rotateCertificates,omitempty"`

	// Remove tells operator to remove TLS configuration
	// +optional
	Remove bool `json:"remove,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MilvusOpsRequestList is a list of MilvusOpsRequests
type MilvusOpsRequestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items is a list of MilvusOpsRequest CRD objects
	Items []MilvusOpsRequest `json:"items,omitempty"`
}
