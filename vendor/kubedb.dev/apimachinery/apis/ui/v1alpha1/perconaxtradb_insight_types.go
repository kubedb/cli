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
	dbapi "kubedb.dev/apimachinery/apis/kubedb/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ResourceKindPerconaXtraDBInsight = "PerconaXtraDBInsight"
	ResourcePerconaXtraDBInsight     = "perconaxtradbinsight"
	ResourcePerconaXtraDBInsights    = "perconaxtradbinsights"
)

// PerconaXtraDBInsightSpec defines the desired state of PerconaXtraDBInsight.
type PerconaXtraDBInsightSpec struct {
	Version                       string   `json:"version"`
	Status                        string   `json:"status"`
	Mode                          string   `json:"mode"`
	MaxConnections                *int64   `json:"maxConnections,omitempty"`
	MaxUsedConnections            *int64   `json:"maxUsedConnections,omitempty"`
	Questions                     *int64   `json:"questions,omitempty"`
	LongQueryTimeThresholdSeconds *float64 `json:"longQueryTimeThresholdSeconds,omitempty"`
	ThreadsConnected              *int64   `json:"threadsConnected,omitempty"`
	ThreadsRunning                *int64   `json:"threadsRunning,omitempty"`
	AbortedClients                *int64   `json:"abortedClients,omitempty"`
	AbortedConnections            *int64   `json:"abortedConnections,omitempty"`
	WSRepClusterSize              *int64   `json:"wsrepClusterSize,omitempty"`
	WSRepClusterStatus            string   `json:"wsrepClusterStatus,omitempty"`
	WSRepConnected                string   `json:"wsrepConnected,omitempty"`
	WSRepReady                    string   `json:"wsrepReady,omitempty"`
	WSRepLocalState               *int64   `json:"wsrepLocalState,omitempty"`
	WSRepLocalStateComment        string   `json:"wsrepLocalStateComment,omitempty"`
	WSRepFlowControlPaused        *float64 `json:"wsrepFlowControlPaused,omitempty"`
}

// PerconaXtraDBInsight is the Schema for the perconaxtradbinsights API.

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type PerconaXtraDBInsight struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PerconaXtraDBInsightSpec  `json:"spec,omitempty"`
	Status dbapi.PerconaXtraDBStatus `json:"status,omitempty"`
}

// PerconaXtraDBInsightList contains a list of PerconaXtraDBInsight.

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type PerconaXtraDBInsightList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PerconaXtraDBInsight `json:"items"`
}
