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
	Collections []MongoDBCollectionSpec `json:"collections" protobuf:"bytes,1,rep,name=collections"`
}

type MongoDBCollectionSpec struct {
	Name      string  `json:"name" protobuf:"bytes,1,opt,name=name"`
	TotalSize []int32 `json:"size" protobuf:"varint,2,rep,name=size"`
}

// MongoDBSchemaOverview is the Schema for the MongoDBSchemaOverviews API

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type MongoDBSchemaOverview struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec MongoDBSchemaOverviewSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
}

// MongoDBSchemaOverviewList contains a list of MongoDBSchemaOverview

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type MongoDBSchemaOverviewList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Items           []MongoDBSchemaOverview `json:"items" protobuf:"bytes,2,rep,name=items"`
}

func init() {
	SchemeBuilder.Register(&MongoDBSchemaOverview{}, &MongoDBSchemaOverviewList{})
}
