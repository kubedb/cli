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
	ResourceKindProxySQLInsight = "ProxySQLInsight"
	ResourceProxySQLInsight     = "proxysqlinsight"
	ResourceProxySQLInsights    = "proxysqlinsights"
)

// RuleExecution gives us insight about the types of executed queries, based on the mysql_query_rules tables.
// To see the digest for corresponding ruleId , please refer to proxysqlsettings api.
type RuleExecution struct {
	RuleId *int64 `json:"ruleId"`
	Hits   *int64 `json:"hits"`
}

// PodInsight gives us insight about the connection status and query traffics of individual pods
type PodInsight struct {
	PodName                    string          `json:"podName"`
	Questions                  *int32          `json:"questions,omitempty"`
	SlowQueries                *int32          `json:"slowQueries,omitempty"`
	AbortedClientConnections   *int32          `json:"abortedClientConnections,omitempty"`
	ConnectedClientConnections *int32          `json:"connectedClientConnections,omitempty"`
	CreatedClientConnections   *int32          `json:"createdClientConnections,omitempty"`
	AbortedServerConnections   *int32          `json:"abortedServerConnections,omitempty"`
	ConnectedServerConnections *int32          `json:"connectedServerConnections,omitempty"`
	CreatedServerConnections   *int32          `json:"createdServerConnections,omitempty"`
	QueryTypeInsight           []RuleExecution `json:"queryInsight,omitempty"`
}

// ProxySQLInsightSpec defines the desired state of ProxySQLInsight
type ProxySQLInsightSpec struct {
	Version                string           `json:"version"`
	Status                 string           `json:"status"`
	MaxConnections         *int32           `json:"maxConnections,omitempty"`
	LongQueryTimeThreshold *metav1.Duration `json:"longQueryTimeThreshold,omitempty"`
	PodInsights            []PodInsight     `json:"podInsights,omitempty"`
}

// ProxySQLInsight is the Schema for the proxysqlinsights API
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ProxySQLInsight struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ProxySQLInsightSpec `json:"spec,omitempty"`
	Status api.ProxySQLStatus  `json:"status,omitempty"`
}

// ProxySQLInsightList contains a list of ProxySQLInsight
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ProxySQLInsightList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ProxySQLInsight `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ProxySQLInsight{}, &ProxySQLInsightList{})
}
