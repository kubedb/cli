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
	dboldapi "kubedb.dev/apimachinery/apis/kubedb/v1alpha2"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ResourceKindSinglestoreInsight = "SinglestoreInsight"
	ResourceSinglestoreInsight     = "singlestoreinsight"
	ResourceSinglestoreInsights    = "singlestoreinsights"
)

// SinglestoreInsightSpec defines the desired state of SinglestoreInsight.
type SinglestoreInsightSpec struct {
	Version              string `json:"version"`
	Status               string `json:"status"`
	Mode                 string `json:"mode"`
	MaxConnections       *int64 `json:"maxConnections,omitempty"`
	MaxConnectionThreads *int64 `json:"maxConnectionThreads,omitempty"`
	MaxUsedConnections   *int64 `json:"maxUsedConnections,omitempty"`
	Connections          *int64 `json:"connections,omitempty"`
	Questions            *int64 `json:"questions,omitempty"`
	ThreadsCached        *int64 `json:"threadsCached,omitempty"`
	ThreadsConnected     *int64 `json:"threadsConnected,omitempty"`
	ThreadsCreated       *int64 `json:"threadsCreated,omitempty"`
	ThreadsRunning       *int64 `json:"threadsRunning,omitempty"`
	AbortedClients       *int64 `json:"abortedClients,omitempty"`
	AbortedConnections   *int64 `json:"abortedConnections,omitempty"`
}

// SinglestoreInsight is the Schema for the singlestoreinsights API.

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type SinglestoreInsight struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SinglestoreInsightSpec     `json:"spec,omitempty"`
	Status dboldapi.SinglestoreStatus `json:"status,omitempty"`
}

// SinglestoreInsightList contains a list of SinglestoreInsight.

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type SinglestoreInsightList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SinglestoreInsight `json:"items"`
}
