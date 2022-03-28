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
	ResourceKindMariaDBQueries = "MariaDBQueries"
	ResourceMariaDBQueries     = "mariadbqueries"
	ResourceMariaDBQuerieses   = "mariadbqueries"
)

// MariaDBQueriesSpec defines the desired state of MariaDBQueries
type MariaDBQueriesSpec struct {
	Queries []MariaDBQuerySpec `json:"queries"`
}

type MariaDBQuerySpec struct {
	StartTime             *metav1.Time `json:"startTime"`
	UserHost              string       `json:"userHost"`
	QueryTimeMilliSeconds string       `json:"queryTimeMilliSeconds"`
	LockTimeMilliSeconds  string       `json:"lockTimeMilliSeconds"`
	RowsSent              *int64       `json:"rowsSent,omitempty"`
	RowsExamined          *int64       `json:"rowsExamined,omitempty"`
	DB                    string       `json:"db"`
	LastInsertId          *int64       `json:"lastInsertId,omitempty"`
	InsertId              *int64       `json:"insertId,omitempty"`
	ServerId              *int64       `json:"serverId,omitempty"`
	SQLText               string       `json:"sqlText,omitempty"`
	ThreadId              *int64       `json:"threadId,omitempty"`
	RowsAffected          *int64       `json:"rowsAffected,omitempty"`
}

// MariaDBQueries is the Schema for the mariadbslowqueries API

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type MariaDBQueries struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec MariaDBQueriesSpec `json:"spec,omitempty"`
}

// MariaDBQueriesList contains a list of MariaDBQueries

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type MariaDBQueriesList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MariaDBQueries `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MariaDBQueries{}, &MariaDBQueriesList{})
}
