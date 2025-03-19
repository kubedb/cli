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

// Elasticsearch is the Schema for the elasticsearches API.

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:resource:path=elasticsearches,singular=elasticsearch
// +kubebuilder:subresource:status
type Elasticsearch struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   dbv1.ElasticsearchSpec `json:"spec,omitempty"`
	Status ElasticsearchStatus    `json:"status,omitempty"`
}

type ElasticsearchStatus struct {
	dbv1.ElasticsearchStatus `json:",inline"`
	GitOps                   GitOpsStatus `json:"gitops,omitempty"`
}

// ElasticsearchList contains a list of Elasticsearch.

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
type ElasticsearchList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Elasticsearch `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Elasticsearch{}, &ElasticsearchList{})
}
