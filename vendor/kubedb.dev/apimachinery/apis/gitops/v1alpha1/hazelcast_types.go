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
	dbv1a2 "kubedb.dev/apimachinery/apis/kubedb/v1alpha2"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Hazelcast is the Schema for the hazelcasts API.

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
type Hazelcast struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   dbv1a2.HazelcastSpec `json:"spec,omitempty"`
	Status HazelcastStatus      `json:"status,omitempty"`
}

type HazelcastStatus struct {
	dbv1a2.HazelcastStatus `json:",inline"`
	GitOps                 GitOpsStatus `json:"gitops,omitempty"`
}

// HazelcastList contains a list of Hazelcast.

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
type HazelcastList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Hazelcast `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Hazelcast{}, &HazelcastList{})
}
