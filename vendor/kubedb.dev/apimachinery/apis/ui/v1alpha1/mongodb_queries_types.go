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
	Queries []MongoDBQuerySpec `json:"queries"`
}

type MongoDBQuerySpec struct {
	Operation                    MongoDBOperation `json:"operation"`
	DatabaseName                 string           `json:"databaseName"`
	CollectionName               string           `json:"collectionName"`
	Command                      string           `json:"command"`
	Count                        *int64           `json:"count,omitempty"`
	AvgExecutionTimeMilliSeconds *int64           `json:"avgExecutionTimeMilliSeconds,omitempty"`
	MinExecutionTimeMilliSeconds *int64           `json:"minExecutionTimeMilliSeconds,omitempty"`
	MaxExecutionTimeMilliSeconds *int64           `json:"maxExecutionTimeMilliSeconds,omitempty"`
}

// MongoDBQueries is the Schema for the MongoDBQueriess API

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type MongoDBQueries struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec MongoDBQueriesSpec `json:"spec,omitempty"`
}

// MongoDBQueriesList contains a list of MongoDBQueries

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type MongoDBQueriesList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MongoDBQueries `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MongoDBQueries{}, &MongoDBQueriesList{})
}
