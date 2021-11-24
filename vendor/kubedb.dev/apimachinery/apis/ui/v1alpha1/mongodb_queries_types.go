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
	ResourceKindMongoDBQueries = "MongoDBQueries"
	ResourceMongoDBQueries     = "mongodbqueries"
	ResourceMongoDBQuerieses   = "mongodbqueries"
)

// MongoDBQueriesSpec defines the desired state of MongoDBQueries
type MongoDBQueriesSpec struct {
	Queries []MongoDBQuerySpec `json:"queries" protobuf:"bytes,1,rep,name=queries"`
}

type MongoDBQuerySpec struct {
	Operation            MongoDBOperation `json:"operation" protobuf:"bytes,1,opt,name=operation,casttype=MongoDBOperation"`
	DatabaseName         string           `json:"databaseName" protobuf:"bytes,2,opt,name=databaseName"`
	CollectionName       string           `json:"collectionName" protobuf:"bytes,3,opt,name=collectionName"`
	Command              string           `json:"command" protobuf:"bytes,4,opt,name=command"`
	Count                int64            `json:"count" protobuf:"varint,5,opt,name=count"`
	AvgExecutionTimeInMS int64            `json:"avgExecutionTimeInMS" protobuf:"varint,6,opt,name=avgExecutionTimeInMS"`
	MinExecutionTimeInMS int64            `json:"minExecutionTimeInMS" protobuf:"varint,7,opt,name=minExecutionTimeInMS"`
	MaxExecutionTimeInMS int64            `json:"maxExecutionTimeInMS" protobuf:"varint,8,opt,name=maxExecutionTimeInMS"`
}

// MongoDBQueries is the Schema for the MongoDBQueriess API

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type MongoDBQueries struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec MongoDBQueriesSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
}

// MongoDBQueriesList contains a list of MongoDBQueries

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type MongoDBQueriesList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Items           []MongoDBQueries `json:"items" protobuf:"bytes,2,rep,name=items"`
}

func init() {
	SchemeBuilder.Register(&MongoDBQueries{}, &MongoDBQueriesList{})
}
