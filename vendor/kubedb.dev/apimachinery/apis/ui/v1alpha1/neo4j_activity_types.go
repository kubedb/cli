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

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

const (
	ResourceKindNeo4jActivity = "Neo4jActivity"
	ResourceNeo4jActivity     = "neo4jactivity"
	ResourceNeo4jActivities   = "neo4jactivities"
)

// Neo4jActivitySpec contains a bounded snapshot of live Neo4j activity.
type Neo4jActivitySpec struct {
	CollectedAt           *metav1.Time         `json:"collectedAt,omitempty"`
	Summary               Neo4jActivitySummary `json:"summary"`
	Transactions          []Neo4jTransaction   `json:"transactions,omitempty"`
	TransactionsTruncated bool                 `json:"transactionsTruncated,omitempty"`
	Incomplete            bool                 `json:"incomplete,omitempty"`
	UnavailableServers    []string             `json:"unavailableServers,omitempty"`
}

type Neo4jActivitySummary struct {
	TotalTransactions          *int32 `json:"totalTransactions,omitempty"`
	RunningQueries             *int32 `json:"runningQueries,omitempty"`
	WaitingQueries             *int32 `json:"waitingQueries,omitempty"`
	BlockedTransactions        *int32 `json:"blockedTransactions,omitempty"`
	LongRunningTransactions    *int32 `json:"longRunningTransactions,omitempty"`
	LongRunningThresholdMillis *int64 `json:"longRunningThresholdMillis,omitempty"`
}

type Neo4jTransaction struct {
	ServerName         string       `json:"serverName"`
	Database           string       `json:"database"`
	TransactionID      string       `json:"transactionID,omitempty"`
	CurrentQueryID     string       `json:"currentQueryID,omitempty"`
	Username           string       `json:"username,omitempty"`
	Query              string       `json:"query,omitempty"`
	QueryTruncated     bool         `json:"queryTruncated,omitempty"`
	Status             string       `json:"status,omitempty"`
	CurrentQueryStatus string       `json:"currentQueryStatus,omitempty"`
	StartTime          *metav1.Time `json:"startTime,omitempty"`
	ElapsedTimeMillis  *int64       `json:"elapsedTimeMillis,omitempty"`
	WaitTimeMillis     *int64       `json:"waitTimeMillis,omitempty"`
	ActiveLockCount    *int64       `json:"activeLockCount,omitempty"`
	ConnectionID       string       `json:"connectionID,omitempty"`
	ClientAddress      string       `json:"clientAddress,omitempty"`
}

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type Neo4jActivity struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec Neo4jActivitySpec `json:"spec,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type Neo4jActivityList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Neo4jActivity `json:"items"`
}
