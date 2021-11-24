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
	ResourceKindRedisSchemaOverview = "RedisSchemaOverview"
	ResourceRedisSchemaOverview     = "redisschemaoverview"
	ResourceRedisSchemaOverviews    = "redisschemaoverviews"
)

// RedisSchemaOverviewSpec defines the desired state of RedisSchemaOverview
type RedisSchemaOverviewSpec struct {
	Databases []RedisDatabaseSpec `json:"databases" protobuf:"bytes,1,rep,name=databases"`
}

type RedisDatabaseSpec struct {
	DBId    string `json:"dbId,omitempty" protobuf:"bytes,1,opt,name=dbId"`
	Keys    string `json:"keys,omitempty" protobuf:"bytes,2,opt,name=keys"`
	Expires string `json:"expires,omitempty" protobuf:"bytes,3,opt,name=expires"`
	AvgTTL  string `json:"avgTTL,omitempty" protobuf:"bytes,4,opt,name=avgTTL"`
}

// RedisSchemaOverview is the Schema for the RedisSchemaOverviews API

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type RedisSchemaOverview struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec RedisSchemaOverviewSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
}

// RedisSchemaOverviewList contains a list of RedisSchemaOverview

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type RedisSchemaOverviewList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Items           []RedisSchemaOverview `json:"items" protobuf:"bytes,2,rep,name=items"`
}

func init() {
	SchemeBuilder.Register(&RedisSchemaOverview{}, &RedisSchemaOverviewList{})
}
