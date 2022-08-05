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
	ResourceKindMariaDBInsight = "MariaDBInsight"
	ResourceMariaDBInsight     = "mariadbinsight"
	ResourceMariaDBInsights    = "mariadbinsights"
)

// MariaDBInsightSpec defines the desired state of MariaDBInsight
type MariaDBInsightSpec struct {
	Version                       string   `json:"version"`
	Status                        string   `json:"status"`
	Mode                          string   `json:"mode"`
	MaxConnections                *int32   `json:"maxConnections,omitempty"`
	MaxUsedConnections            *int32   `json:"maxUsedConnections,omitempty"`
	Questions                     *int32   `json:"questions,omitempty"`
	LongQueryTimeThresholdSeconds *float64 `json:"longQueryTimeThresholdSeconds,omitempty"`
	NumberOfSlowQueries           *int32   `json:"numberOfSlowQueries,omitempty"`
	AbortedClients                *int32   `json:"abortedClients,omitempty"`
	AbortedConnections            *int32   `json:"abortedConnections,omitempty"`
	ThreadsCached                 *int32   `json:"threadsCached,omitempty"`
	ThreadsConnected              *int32   `json:"threadsConnected,omitempty"`
	ThreadsCreated                *int32   `json:"threadsCreated,omitempty"`
	ThreadsRunning                *int32   `json:"threadsRunning,omitempty"`
}

// MariaDBInsight is the Schema for the mariaDBinsights API

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type MariaDBInsight struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MariaDBInsightSpec `json:"spec,omitempty"`
	Status api.MariaDBStatus  `json:"status,omitempty"`
}

// MariaDBInsightList contains a list of MariaDBInsight

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type MariaDBInsightList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MariaDBInsight `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MariaDBInsight{}, &MariaDBInsightList{})
}
