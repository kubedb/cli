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
	api "kubedb.dev/apimachinery/apis/kubedb/v1alpha2"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ResourceKindRedisInsight = "RedisInsight"
	ResourceRedisInsight     = "redisinsight"
	ResourceRedisInsights    = "redisinsights"
)

// RedisInsightSpec defines the desired state of RedisInsight
type RedisInsightSpec struct {
	Version                      string `json:"version"`
	Status                       string `json:"status"`
	Mode                         string `json:"mode"`
	EvictionPolicy               string `json:"evictionPolicy"`
	MaxClients                   *int64 `json:"maxClients,omitempty"`
	ConnectedClients             *int64 `json:"connectedClients,omitempty"`
	BlockedClients               *int64 `json:"blockedClients,omitempty"`
	TotalKeys                    *int64 `json:"totalKeys,omitempty"`
	ExpiredKeys                  *int64 `json:"expiredKeys,omitempty"`
	EvictedKeys                  *int64 `json:"evictedKeys,omitempty"`
	ReceivedConnections          *int64 `json:"receivedConnections,omitempty"`
	RejectedConnections          *int64 `json:"rejectedConnections,omitempty"`
	SlowLogThresholdMicroSeconds *int64 `json:"slowLogThresholdMicroSeconds,omitempty"`
	SlowLogMaxLen                *int64 `json:"slowLogMaxLen,omitempty"`
}

// RedisInsight is the Schema for the redisinsights API

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type RedisInsight struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RedisInsightSpec `json:"spec,omitempty"`
	Status api.RedisStatus  `json:"status,omitempty"`
}

// RedisInsightList contains a list of RedisInsight

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type RedisInsightList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RedisInsight `json:"items"`
}

func init() {
	SchemeBuilder.Register(&RedisInsight{}, &RedisInsightList{})
}
