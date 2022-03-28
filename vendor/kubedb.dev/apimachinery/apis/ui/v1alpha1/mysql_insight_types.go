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
	ResourceKindMySQLInsight = "MySQLInsight"
	ResourceMySQLInsight     = "mysqlinsight"
	ResourceMySQLInsights    = "mysqlinsights"
)

// MySQLInsightSpec defines the desired state of MySQLInsight
type MySQLInsightSpec struct {
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

// MySQLInsight is the Schema for the mysqlinsights API

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type MySQLInsight struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MySQLInsightSpec `json:"spec,omitempty"`
	Status api.MySQLStatus  `json:"status,omitempty"`
}

// MySQLInsightList contains a list of MySQLInsight

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type MySQLInsightList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MySQLInsight `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MySQLInsight{}, &MySQLInsightList{})
}
