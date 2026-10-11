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
	ResourceKindPerconaXtraDBQueries = "PerconaXtraDBQueries"
	ResourcePerconaXtraDBQueries     = "perconaxtradbqueries"
	ResourcePerconaXtraDBQuerieses   = "perconaxtradbqueries"
)

// PerconaXtraDBQueriesSpec defines the desired state of PerconaXtraDBQueries.
type PerconaXtraDBQueriesSpec struct {
	Queries []PerconaXtraDBQuerySpec `json:"queries"`
}

// PerconaXtraDBQuerySpec describes cumulative statistics for one normalized statement digest.
type PerconaXtraDBQuerySpec struct {
	SchemaName                       string       `json:"schemaName,omitempty"`
	Digest                           string       `json:"digest,omitempty"`
	DigestText                       string       `json:"digestText,omitempty"`
	Count                            *int64       `json:"count,omitempty"`
	TotalExecutionTimeMilliSeconds   *float64     `json:"totalExecutionTimeMilliSeconds,omitempty"`
	AverageExecutionTimeMilliSeconds *float64     `json:"averageExecutionTimeMilliSeconds,omitempty"`
	MinimumExecutionTimeMilliSeconds *float64     `json:"minimumExecutionTimeMilliSeconds,omitempty"`
	MaximumExecutionTimeMilliSeconds *float64     `json:"maximumExecutionTimeMilliSeconds,omitempty"`
	TotalLockTimeMilliSeconds        *float64     `json:"totalLockTimeMilliSeconds,omitempty"`
	RowsAffected                     *int64       `json:"rowsAffected,omitempty"`
	RowsSent                         *int64       `json:"rowsSent,omitempty"`
	RowsExamined                     *int64       `json:"rowsExamined,omitempty"`
	FirstSeen                        *metav1.Time `json:"firstSeen,omitempty"`
	LastSeen                         *metav1.Time `json:"lastSeen,omitempty"`
}

// PerconaXtraDBQueries is the Schema for the perconaxtradbqueries API.

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type PerconaXtraDBQueries struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec PerconaXtraDBQueriesSpec `json:"spec,omitempty"`
}

// PerconaXtraDBQueriesList contains a list of PerconaXtraDBQueries.

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type PerconaXtraDBQueriesList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PerconaXtraDBQueries `json:"items"`
}
