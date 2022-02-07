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
	ResourceKindMongoDBSchemaOverview = "MongoDBSchemaOverview"
	ResourceMongoDBSchemaOverview     = "mongodbschemaoverview"
	ResourceMongoDBSchemaOverviews    = "mongodbschemaoverviews"
)

// MongoDBSchemaOverviewSpec defines the desired state of MongoDBSchemaOverview
type MongoDBSchemaOverviewSpec struct {
	Collections []MongoDBCollectionSpec `json:"collections"`
}

type MongoDBCollectionSpec struct {
	Name string `json:"name"`

	// Slice is used to store shards specific collection size for Sharded MongoDB
	TotalSize []int32 `json:"size"`
}

// MongoDBSchemaOverview is the Schema for the MongoDBSchemaOverviews API

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type MongoDBSchemaOverview struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec MongoDBSchemaOverviewSpec `json:"spec,omitempty"`
}

// MongoDBSchemaOverviewList contains a list of MongoDBSchemaOverview

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type MongoDBSchemaOverviewList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MongoDBSchemaOverview `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MongoDBSchemaOverview{}, &MongoDBSchemaOverviewList{})
}
