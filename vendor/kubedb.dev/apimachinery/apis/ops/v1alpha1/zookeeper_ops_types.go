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
	ResourceCodeZooKeeperOpsRequest     = "zkops"
	ResourceKindZooKeeperOpsRequest     = "ZooKeeperOpsRequest"
	ResourceSingularZooKeeperOpsRequest = "zookeeperopsrequest"
	ResourcePluralZooKeeperOpsRequest   = "zookeeperopsrequests"
)

// ZooKeeperDBOpsRequest defines a ZooKeeper DBA operation.

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=zookeeperopsrequests,singular=zookeeperopsrequest,shortName=zkops,categories={ops,kubedb,appscode}
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Type",type="string",JSONPath=".spec.type"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type ZooKeeperOpsRequest struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              ZooKeeperOpsRequestSpec `json:"spec,omitempty"`
	Status            OpsRequestStatus        `json:"status,omitempty"`
}

// ZooKeeperOpsRequestSpec is the spec for ZooKeeperOpsRequest
type ZooKeeperOpsRequestSpec struct {
	// Specifies the ZooKeeper reference
	DatabaseRef core.LocalObjectReference `json:"databaseRef"`
	// Specifies the ops request type: UpdateVersion, HorizontalScaling, VerticalScaling etc.
	Type ZooKeeperOpsRequestType `json:"type"`
	// Specifies information necessary for restarting database
	Restart *RestartSpec `json:"restart,omitempty"`
	// Timeout for each step of the ops request in second. If a step doesn't finish within the specified timeout, the ops request will result in failure.
	Timeout *metav1.Duration `json:"timeout,omitempty"`
	// ApplyOption is to control the execution of OpsRequest depending on the database state.
	// +kubebuilder:default="IfReady"
	Apply ApplyOption `json:"apply,omitempty"`
}

// +kubebuilder:validation:Enum=Restart
// ENUM(Restart)
type ZooKeeperOpsRequestType string

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ZooKeeperOpsRequestList is a list of ZooKeeperOpsRequests
type ZooKeeperOpsRequestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items is a list of ZooKeeperOpsRequest CRD objects
	Items []ZooKeeperOpsRequest `json:"items,omitempty"`
}
