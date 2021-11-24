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
	ResourceKindRedisQueries = "RedisQueries"
	ResourceRedisQueries     = "redisqueries"
	ResourceRedisQuerieses   = "redisqueries"
)

// RedisQueriesSpec defines the desired state of RedisQueries
type RedisQueriesSpec struct {
	Queries []RedisQuerySpec `json:"queries" protobuf:"bytes,1,rep,name=queries"`
}

type RedisQuerySpec struct {
	QueryId                int64    `json:"queryId" protobuf:"varint,1,opt,name=queryId"`
	QueryTimestamp         int64    `json:"queryTimestamp" protobuf:"varint,2,opt,name=queryTimestamp"`
	ExecTimeInMircoSeconds int64    `json:"execTimeInMircoSeconds" protobuf:"varint,3,opt,name=execTimeInMircoSeconds"`
	Args                   []string `json:"args" protobuf:"bytes,4,rep,name=args"`
}

// RedisQueries is the Schema for the RedisQueries API

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type RedisQueries struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec RedisQueriesSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
}

// RedisQueriesList contains a list of RedisQueries

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type RedisQueriesList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Items           []RedisQueries `json:"items" protobuf:"bytes,2,rep,name=items"`
}

func init() {
	SchemeBuilder.Register(&RedisQueries{}, &RedisQueriesList{})
}
