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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ResourceKindSinglestoreQueries = "SinglestoreQueries"
	ResourceSinglestoreQueries     = "singlestorequeries"
	ResourceSinglestoreQuerieses   = "singlestorequeries"
)

// SinglestoreQueriesSpec defines the desired state of SinglestoreQueries.
type SinglestoreQueriesSpec struct {
	// Queries contains cumulative normalized plan-cache entries, not individual query executions.
	Queries []SinglestoreQuerySpec `json:"queries"`
}

// SinglestoreQuerySpec contains cumulative statistics for a normalized plan-cache entry.
type SinglestoreQuerySpec struct {
	DatabaseName                     string       `json:"databaseName,omitempty"`
	QueryText                        string       `json:"queryText,omitempty"`
	PlanID                           *int64       `json:"planId,omitempty"`
	PlanType                         string       `json:"planType,omitempty"`
	Commits                          *int64       `json:"commits,omitempty"`
	Rollbacks                        *int64       `json:"rollbacks,omitempty"`
	RowCount                         *int64       `json:"rowCount,omitempty"`
	ExecutionTimeMilliSeconds        *float64     `json:"executionTimeMilliSeconds,omitempty"`
	AverageExecutionTimeMilliSeconds *float64     `json:"averageExecutionTimeMilliSeconds,omitempty"`
	LastExecuted                     *metav1.Time `json:"lastExecuted,omitempty"`
	AverageMemoryUseBytes            *int64       `json:"averageMemoryUseBytes,omitempty"`
	CPUTimeMilliSeconds              *float64     `json:"cpuTimeMilliSeconds,omitempty"`
}

// SinglestoreQueries is the Schema for the singlestorequeries API.

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type SinglestoreQueries struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec SinglestoreQueriesSpec `json:"spec,omitempty"`
}

// SinglestoreQueriesList contains a list of SinglestoreQueries.

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type SinglestoreQueriesList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SinglestoreQueries `json:"items"`
}
