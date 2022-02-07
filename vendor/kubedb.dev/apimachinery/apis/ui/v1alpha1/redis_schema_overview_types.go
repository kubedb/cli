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
	Databases []RedisDatabaseSpec `json:"databases"`
}

type RedisDatabaseSpec struct {
	DBId               string       `json:"dbId"`
	Keys               string       `json:"keys"`
	Expires            *metav1.Time `json:"expires,omitempty"`
	AvgTTLMilliSeconds string       `json:"avgTTLMilliSeconds,omitempty"`
}

// RedisSchemaOverview is the Schema for the RedisSchemaOverviews API

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type RedisSchemaOverview struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec RedisSchemaOverviewSpec `json:"spec,omitempty"`
}

// RedisSchemaOverviewList contains a list of RedisSchemaOverview

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type RedisSchemaOverviewList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RedisSchemaOverview `json:"items"`
}

func init() {
	SchemeBuilder.Register(&RedisSchemaOverview{}, &RedisSchemaOverviewList{})
}
