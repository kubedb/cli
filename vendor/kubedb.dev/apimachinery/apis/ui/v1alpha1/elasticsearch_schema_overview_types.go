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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ResourceKindElasticsearchSchemaOverview = "ElasticsearchSchemaOverview"
	ResourceElasticsearchSchemaOverview     = "elasticsearchschemaoverview"
	ResourceElasticsearchSchemaOverviews    = "elasticsearchschemaoverviews"
)

// ElasticsearchSchemaOverviewSpec defines the desired state of ElasticsearchSchemaOverview
type ElasticsearchSchemaOverviewSpec struct {
	Indices []ElasticsearchIndexSpec `json:"indices"`
}

type ElasticsearchIndexSpec struct {
	IndexName             *string `json:"indexName,omitempty"`
	PrimaryStoreSizeBytes *string `json:"primaryStoreSizeBytes,omitempty"`
	TotalStoreSizeBytes   *string `json:"totalStoreSizeBytes,omitempty"`
}

// ElasticsearchSchemaOverview is the Schema for the elasticsearchindices API

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ElasticsearchSchemaOverview struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec ElasticsearchSchemaOverviewSpec `json:"spec,omitempty"`
}

// ElasticsearchSchemaOverviewList contains a list of ElasticsearchSchemaOverview

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ElasticsearchSchemaOverviewList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ElasticsearchSchemaOverview `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ElasticsearchSchemaOverview{}, &ElasticsearchSchemaOverviewList{})
}
