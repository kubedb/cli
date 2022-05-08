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
	ResourceKindElasticsearchInsight = "ElasticsearchInsight"
	ResourceElasticsearchInsight     = "elasticsearchinsight"
	ResourceElasticsearchInsights    = "elasticsearchinsights"
)

// ElasticsearchInsightSpec defines the desired state of ElasticsearchInsight
type ElasticsearchInsightSpec struct {
	Version                    string `json:"version"`
	Status                     string `json:"status"`
	Mode                       string `json:"mode"`
	ElasticsearchClusterHealth `json:",inline"`
}

type ElasticsearchClusterHealth struct {
	ActivePrimaryShards               *float64 `json:"activePrimaryShards,omitempty"`
	ActiveShards                      *float64 `json:"activeShards,omitempty"`
	ActiveShardsPercentAsNumber       *float64 `json:"activeShardsPercentAsNumber,omitempty"`
	ClusterName                       string   `json:"clusterName,omitempty"`
	DelayedUnassignedShards           *float64 `json:"delayedUnassignedShards,omitempty"`
	InitializingShards                *float64 `json:"initializingShards,omitempty"`
	NumberOfDataNodes                 *float64 `json:"numberOfDataNodes,omitempty"`
	NumberOfInFlightFetch             *float64 `json:"numberOfInFlightFetch,omitempty"`
	NumberOfNodes                     *float64 `json:"numberOfNodes,omitempty"`
	NumberOfPendingTasks              *float64 `json:"numberOfPendingTasks,omitempty"`
	RelocatingShards                  *float64 `json:"relocatingShards,omitempty"`
	ClusterStatus                     string   `json:"clusterStatus,omitempty"`
	UnassignedShards                  *float64 `json:"unassignedShards,omitempty"`
	TaskMaxWaitingInQueueMilliSeconds *float64 `json:"taskMaxWaitingInQueueMilliSeconds,omitempty"`
}

// ElasticsearchInsight is the Schema for the elasticsearchinsights API

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ElasticsearchInsight struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ElasticsearchInsightSpec `json:"spec,omitempty"`
	Status api.ElasticsearchStatus  `json:"status,omitempty"`
}

// ElasticsearchInsightList contains a list of ElasticsearchInsight

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ElasticsearchInsightList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ElasticsearchInsight `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ElasticsearchInsight{}, &ElasticsearchInsightList{})
}
