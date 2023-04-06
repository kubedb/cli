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
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ResourceCodeEtcdOpsRequest     = "etcdops"
	ResourceKindEtcdOpsRequest     = "EtcdOpsRequest"
	ResourceSingularEtcdOpsRequest = "etcdopsrequest"
	ResourcePluralEtcdOpsRequest   = "etcdopsrequests"
)

// EtcdOpsRequest defines a Etcd DBA operation.

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=etcdopsrequests,singular=etcdopsrequest,shortName=etcdops,categories={datastore,kubedb,appscode}
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Type",type="string",JSONPath=".spec.type"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type EtcdOpsRequest struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              EtcdOpsRequestSpec `json:"spec,omitempty"`
	Status            OpsRequestStatus   `json:"status,omitempty"`
}

// EtcdOpsRequestSpec is the spec for EtcdOpsRequest
type EtcdOpsRequestSpec struct {
	// Specifies the Etcd reference
	DatabaseRef core.LocalObjectReference `json:"databaseRef"`
	// Specifies the ops request type: Upgrade, HorizontalScaling, VerticalScaling etc.
	Type EtcdOpsRequestType `json:"type"`
	// Specifies information necessary for upgrading Etcd
	UpdateVersion *EtcdUpdateVersionSpec `json:"updateVersion,omitempty"`
	// Specifies information necessary for horizontal scaling
	HorizontalScaling *EtcdHorizontalScalingSpec `json:"horizontalScaling,omitempty"`
	// Specifies information necessary for vertical scaling
	VerticalScaling *EtcdVerticalScalingSpec `json:"verticalScaling,omitempty"`
	// Specifies information necessary for volume expansion
	VolumeExpansion *EtcdVolumeExpansionSpec `json:"volumeExpansion,omitempty"`
	// Specifies information necessary for custom configuration of Etcd
	Configuration *EtcdCustomConfigurationSpec `json:"configuration,omitempty"`
	// Specifies information necessary for configuring TLS
	TLS *TLSSpec `json:"tls,omitempty"`
	// Specifies information necessary for restarting database
	Restart *RestartSpec `json:"restart,omitempty"`
	// ApplyOption is to control the execution of OpsRequest depending on the database state.
	// +kubebuilder:default="IfReady"
	Apply ApplyOption `json:"apply,omitempty"`
}

// +kubebuilder:validation:Enum=UpdateVersion;HorizontalScaling;VerticalScaling;VolumeExpansion;Restart;Reconfigure;ReconfigureTLS
// ENUM(UpdateVersion, HorizontalScaling, VerticalScaling, VolumeExpansion, Restart, Reconfigure, ReconfigureTLS)
type EtcdOpsRequestType string

// EtcdReplicaReadinessCriteria is the criteria for checking readiness of a Etcd pod
// after updating, horizontal scaling etc.
type EtcdReplicaReadinessCriteria struct{}

type EtcdUpdateVersionSpec struct {
	// Specifies the target version name from catalog
	TargetVersion     string                        `json:"targetVersion,omitempty"`
	ReadinessCriteria *EtcdReplicaReadinessCriteria `json:"readinessCriteria,omitempty"`
}

// HorizontalScaling is the spec for Etcd horizontal scaling
type EtcdHorizontalScalingSpec struct{}

// EtcdVerticalScalingSpec is the spec for Etcd vertical scaling
type EtcdVerticalScalingSpec struct {
	ReadinessCriteria *EtcdReplicaReadinessCriteria `json:"readinessCriteria,omitempty"`
}

// EtcdVolumeExpansionSpec is the spec for Etcd volume expansion
type EtcdVolumeExpansionSpec struct{}

type EtcdCustomConfigurationSpec struct{}

type EtcdCustomConfiguration struct {
	ConfigMap *core.LocalObjectReference `json:"configMap,omitempty"`
	Data      map[string]string          `json:"data,omitempty"`
	Remove    bool                       `json:"remove,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// EtcdOpsRequestList is a list of EtcdOpsRequests
type EtcdOpsRequestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items is a list of EtcdOpsRequest CRD objects
	Items []EtcdOpsRequest `json:"items,omitempty"`
}
