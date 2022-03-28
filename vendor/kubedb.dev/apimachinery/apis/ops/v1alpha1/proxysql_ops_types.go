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
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kmapi "kmodules.xyz/client-go/api/v1"
)

const (
	ResourceKindProxySQLOpsRequest     = "ProxySQLOpsRequest"
	ResourceSingularProxySQLOpsRequest = "proxysqlopsrequest"
	ResourcePluralProxySQLOpsRequest   = "proxysqlopsrequests"
)

// ProxySQLOpsRequest defines a ProxySQL load-balancer DBA operation.

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=proxysqlopsrequests,singular=proxysqlopsrequest,shortName=proxyops,categories={datastore,kubedb,appscode}
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Type",type="string",JSONPath=".spec.type"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type ProxySQLOpsRequest struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              ProxySQLOpsRequestSpec   `json:"spec,omitempty"`
	Status            ProxySQLOpsRequestStatus `json:"status,omitempty"`
}

// ProxySQLOpsRequestSpec is the spec for ProxySQLOpsRequest
type ProxySQLOpsRequestSpec struct {
	// Specifies the ProxySQL reference
	DatabaseRef core.LocalObjectReference `json:"databaseRef"`
	// Specifies the ops request type: Upgrade, HorizontalScaling, VerticalScaling etc.
	Type OpsRequestType `json:"type"`
	// Specifies information necessary for upgrading ProxySQL
	Upgrade *ProxySQLUpgradeSpec `json:"upgrade,omitempty"`
	// Specifies information necessary for horizontal scaling
	HorizontalScaling *ProxySQLHorizontalScalingSpec `json:"horizontalScaling,omitempty"`
	// Specifies information necessary for vertical scaling
	VerticalScaling *ProxySQLVerticalScalingSpec `json:"verticalScaling,omitempty"`
	// Specifies information necessary for volume expansion
	VolumeExpansion *ProxySQLVolumeExpansionSpec `json:"volumeExpansion,omitempty"`
	// Specifies information necessary for custom configuration of ProxySQL
	Configuration *ProxySQLCustomConfigurationSpec `json:"configuration,omitempty"`
	// Specifies information necessary for configuring TLS
	TLS *TLSSpec `json:"tls,omitempty"`
	// Specifies information necessary for restarting database
	Restart *RestartSpec `json:"restart,omitempty"`
}

// ProxySQLReplicaReadinessCriteria is the criteria for checking readiness of a ProxySQL pod
// after updating, horizontal scaling etc.
type ProxySQLReplicaReadinessCriteria struct{}

type ProxySQLUpgradeSpec struct {
	// Specifies the target version name from catalog
	TargetVersion     string                            `json:"targetVersion,omitempty"`
	ReadinessCriteria *ProxySQLReplicaReadinessCriteria `json:"readinessCriteria,omitempty"`
}

// HorizontalScaling is the spec for ProxySQL horizontal scaling
type ProxySQLHorizontalScalingSpec struct{}

// ProxySQLVerticalScalingSpec is the spec for ProxySQL vertical scaling
type ProxySQLVerticalScalingSpec struct {
	ReadinessCriteria *ProxySQLReplicaReadinessCriteria `json:"readinessCriteria,omitempty"`
}

// ProxySQLVolumeExpansionSpec is the spec for ProxySQL volume expansion
type ProxySQLVolumeExpansionSpec struct{}

type ProxySQLCustomConfigurationSpec struct{}

type ProxySQLCustomConfiguration struct {
	ConfigMap *core.LocalObjectReference `json:"configMap,omitempty"`
	Data      map[string]string          `json:"data,omitempty"`
	Remove    bool                       `json:"remove,omitempty"`
}

// ProxySQLOpsRequestStatus is the status for ProxySQLOpsRequest
type ProxySQLOpsRequestStatus struct {
	Phase OpsRequestPhase `json:"phase,omitempty"`
	// observedGeneration is the most recent generation observed for this resource. It corresponds to the
	// resource's generation, which is updated on mutation by the API Server.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
	// Conditions applied to the request, such as approval or denial.
	// +optional
	Conditions []kmapi.Condition `json:"conditions,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ProxySQLOpsRequestList is a list of ProxySQLOpsRequests
type ProxySQLOpsRequestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items is a list of ProxySQLOpsRequest CRD objects
	Items []ProxySQLOpsRequest `json:"items,omitempty"`
}
