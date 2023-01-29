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
	ResourceKindPgBouncerPoolOverview = "PgBouncerPoolOverview"
	ResourcePgBouncerPoolOverview     = "pgbouncerpooloverview"
	ResourcePgBouncerPoolOverviews    = "pgbouncerpooloverviews"
)

type PgBouncerPoolOverviewSpec struct {
	Pools []PgBouncerPool `json:"pools"`
}

type PgBouncerPool struct {
	// PodName represents the corresponding pod for the pool
	PodName string `json:"podName"`
	// Database represents the database of the pool
	Database string `json:"database"`
	// User represents the user of the pool
	User string `json:"user"`
	// ActiveClientConnections represents client connections that are either linked to server connections
	// or are idle with no queries waiting to be processed
	ActiveClientConnections *int32 `json:"activeClientConnections,omitempty"`
	// WaitingClientConnections represents client connections that have sent queries
	// but have not yet got a server connection
	WaitingClientConnections *int32 `json:"waitingClientConnections,omitempty"`
	// ActiveQueryCancellationRequest represents client connections that have forwarded query cancellations to the server
	ActiveQueryCancellationRequest *int32 `json:"activeQueryCancellationRequest,omitempty"`
	// WaitingQueryCancellationRequest represents client connections that are waiting to forwarded query cancellations to the server
	WaitingQueryCancellationRequest *int32 `json:"waitingQueryCancellationRequest,omitempty"`
	// ActiveServerConnections represents server connections that are linked to a client.
	ActiveServerConnections *int32 `json:"activeServerConnections,omitempty"`
	// ActiveServersCancelRequest represents server connections that are currently forwarding a cancel request.
	ActiveServersCancelRequest *int32 `json:"activeServersCancelRequest,omitempty"`
	// ServersBeingCanceled represents servers that normally could become idle but are waiting to do
	// so until all in-flight cancel requests have completed
	ServersBeingCanceled *int32 `json:"serversBeingCanceled,omitempty"`
	// IdleServers represents server connections that are unused and immediately usable for client queries.
	IdleServers *int32 `json:"idleServers,omitempty"`
	// UsedServers represents server connections that have been idle for more than server_check_delay
	UsedServers *int32 `json:"usedServers,omitempty"`
	// ServersTested represents server connections that are currently running either server_reset_query or server_check_query
	TestedServers *int32 `json:"testedServers,omitempty"`
	// ServersInLogin represents server connections currently in the process of logging in
	ServersInLogin *int32 `json:"serversInLogin,omitempty"`
	// MaxWaitMS represents how long the first (oldest) client in the queue has waited
	MaxWaitMS *metav1.Duration `json:"maxWaitMS,omitempty"`
	// Mode represents the pooling mode in use. ex: session, transaction, statement
	Mode string `json:"mode"`
}

// PgBouncerPoolOverview is the Schema for the PgBouncerPoolOverviews API
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type PgBouncerPoolOverview struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec PgBouncerPoolOverviewSpec `json:"spec,omitempty"`
}

// PgBouncerPoolOverviewList contains a list of PgBouncerPoolOverviews
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type PgBouncerPoolOverviewList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PgBouncerPoolOverview `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PgBouncerPoolOverview{}, &PgBouncerPoolOverviewList{})
}
