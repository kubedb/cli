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
	ResourceKindNeo4jInsight = "Neo4jInsight"
	ResourceNeo4jInsight     = "neo4jinsight"
	ResourceNeo4jInsights    = "neo4jinsights"
)

// Neo4jInsightSpec defines the observed database information displayed by the UI.
type Neo4jInsightSpec struct {
	Version       string                `json:"version"`
	Edition       string                `json:"edition,omitempty"`
	Mode          string                `json:"mode"`
	ClusterHealth Neo4jClusterHealth    `json:"clusterHealth"`
	Servers       []Neo4jServerStatus   `json:"servers,omitempty"`
	Databases     []Neo4jDatabaseStatus `json:"databases,omitempty"`
}

type Neo4jClusterHealth struct {
	TotalServers       *int32 `json:"totalServers,omitempty"`
	AvailableServers   *int32 `json:"availableServers,omitempty"`
	UnavailableServers *int32 `json:"unavailableServers,omitempty"`
	TotalDatabases     *int32 `json:"totalDatabases,omitempty"`
	OnlineDatabases    *int32 `json:"onlineDatabases,omitempty"`
	OfflineDatabases   *int32 `json:"offlineDatabases,omitempty"`
}

type Neo4jServerStatus struct {
	Name           string   `json:"name"`
	Address        string   `json:"address,omitempty"`
	State          string   `json:"state"`
	Health         string   `json:"health"`
	Version        string   `json:"version,omitempty"`
	ModeConstraint string   `json:"modeConstraint,omitempty"`
	Hosting        []string `json:"hosting,omitempty"`
}

type Neo4jDatabaseStatus struct {
	Name                      string   `json:"name"`
	Type                      string   `json:"type,omitempty"`
	Access                    string   `json:"access,omitempty"`
	RequestedStatus           string   `json:"requestedStatus,omitempty"`
	CurrentStatus             string   `json:"currentStatus,omitempty"`
	StatusMessage             string   `json:"statusMessage,omitempty"`
	Default                   bool     `json:"default,omitempty"`
	WriterAddress             string   `json:"writerAddress,omitempty"`
	Constituents              []string `json:"constituents,omitempty"`
	CurrentPrimariesCount     *int32   `json:"currentPrimariesCount,omitempty"`
	CurrentSecondariesCount   *int32   `json:"currentSecondariesCount,omitempty"`
	RequestedPrimariesCount   *int32   `json:"requestedPrimariesCount,omitempty"`
	RequestedSecondariesCount *int32   `json:"requestedSecondariesCount,omitempty"`
	NodeCount                 *int64   `json:"nodeCount,omitempty"`
	RelationshipCount         *int64   `json:"relationshipCount,omitempty"`
}

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type Neo4jInsight struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   Neo4jInsightSpec  `json:"spec,omitempty"`
	Status dbapi.Neo4jStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type Neo4jInsightList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Neo4jInsight `json:"items"`
}
