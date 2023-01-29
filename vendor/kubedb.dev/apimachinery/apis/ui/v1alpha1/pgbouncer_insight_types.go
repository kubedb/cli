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
	ResourceKindPgBouncerInsight = "PgBouncerInsight"
	ResourcePgBouncerInsight     = "pgbouncerinsight"
	ResourcePgBouncerInsights    = "pgbouncerinsights"
)

type PgBouncerInsightSpec struct {
	Version        string                `json:"version"`
	Status         string                `json:"status"`
	SSLMode        api.SSLMode           `json:"sslMode,omitempty"`
	MaxConnections *int32                `json:"maxConnections,omitempty"`
	PodInsights    []PgBouncerPodInsight `json:"podInsights,omitempty"`
}

type PgBouncerPodInsight struct {
	// PodName represents the name of the pod
	PodName string `json:"podName"`
	// Databases represents the number of databases
	Databases *int32 `json:"databases,omitempty"`
	// Users represents the number of users
	Users *int32 `json:"users,omitempty"`
	// Pools represents the number of pools
	Pools *int32 `json:"pools,omitempty"`
	// FreeClients represents the number of free clients
	FreeClients *int32 `json:"freeClients,omitempty"`
	// UsedClients represents the number of used clients
	UsedClients *int32 `json:"usedClients,omitempty"`
	// LoginClients represents the number of clients in the login state
	LoginClients *int32 `json:"loginClients,omitempty"`
	// FreeServers represents the number of free servers
	FreeServers *int32 `json:"freeServers,omitempty"`
	// UsedServers represents the number of used servers
	UsedServers *int32 `json:"usedServers,omitempty"`
	// TotalQueryCount represents the total number of query counts
	TotalQueryCount *int32 `json:"totalQueryCount,omitempty"`
	// AverageQueryCount represents the average number of query counts
	AverageQueryCount *int32 `json:"averageQueryCount,omitempty"`
	// TotalQueryTime represents the total time spent for the queries
	TotalQueryTimeMS *int32 `json:"totalQueryTimeMS,omitempty"`
	// AverageQueryTime represents the average time spent for a single query
	AverageQueryTimeMS *int32 `json:"averageQueryTimeMS,omitempty"`
}

// PgBouncerInsight is the Schema for the pgbouncerinsights API
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type PgBouncerInsight struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PgBouncerInsightSpec `json:"spec,omitempty"`
	Status api.PgBouncerStatus  `json:"status,omitempty"`
}

// PgBouncerInsightList contains a list of PgBouncerInsight
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type PgBouncerInsightList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PgBouncerInsight `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PgBouncerInsight{}, &PgBouncerInsightList{})
}
