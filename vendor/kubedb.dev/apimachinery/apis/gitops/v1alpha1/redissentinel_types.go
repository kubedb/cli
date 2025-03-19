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
	dbv1 "kubedb.dev/apimachinery/apis/kubedb/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// RedisSentinel is the Schema for the redissentinels API.

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
type RedisSentinel struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   dbv1.RedisSentinelSpec `json:"spec,omitempty"`
	Status RedisSentinelStatus    `json:"status,omitempty"`
}

type RedisSentinelStatus struct {
	dbv1.RedisSentinelStatus `json:",inline"`
	GitOps                   GitOpsStatus `json:"gitops,omitempty"`
}

// RedisSentinelList contains a list of RedisSentinel.

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
type RedisSentinelList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RedisSentinel `json:"items"`
}

func init() {
	SchemeBuilder.Register(&RedisSentinel{}, &RedisSentinelList{})
}
