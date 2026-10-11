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
	dbapi "kubedb.dev/apimachinery/apis/kubedb/v1alpha2"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ResourceKindDruidInsight = "DruidInsight"
	ResourceDruidInsight     = "druidinsight"
	ResourceDruidInsights    = "druidinsights"
)

type DruidInsightSpec struct {
	Version     string                  `json:"version"`
	Services    []DruidServiceStatus    `json:"services,omitempty"`
	Datasources []DruidDatasourceStatus `json:"datasources,omitempty"`
}

type DruidServiceStatus struct {
	Type          string `json:"type"`
	ReplicaCount  int32  `json:"replicaCount"`
	Version       string `json:"version,omitempty"`
	LeaderAddress string `json:"leaderAddress,omitempty"`
}

type DruidDatasourceStatus struct {
	Name                string       `json:"name"`
	SegmentCount        *int64       `json:"segmentCount,omitempty"`
	AvailableSegments   *int64       `json:"availableSegments,omitempty"`
	UnavailableSegments *int64       `json:"unavailableSegments,omitempty"`
	RealtimeSegments    *int64       `json:"realtimeSegments,omitempty"`
	TotalRows           *int64       `json:"totalRows,omitempty"`
	TotalSizeBytes      *int64       `json:"totalSizeBytes,omitempty"`
	MinTime             *metav1.Time `json:"minTime,omitempty"`
	MaxTime             *metav1.Time `json:"maxTime,omitempty"`
}

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type DruidInsight struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DruidInsightSpec  `json:"spec,omitempty"`
	Status dbapi.DruidStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type DruidInsightList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DruidInsight `json:"items"`
}
