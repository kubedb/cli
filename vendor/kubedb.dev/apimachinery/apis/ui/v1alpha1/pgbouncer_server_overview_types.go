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
	ResourceKindPgBouncerServerOverView = "PgBouncerServerOverview"
	ResourcePgBouncerServerOverview     = "pgbouncerserveroverview"
	ResourcePgBouncerServerOverviews    = "pgbouncerserveroverviews"
)

type PgBouncerServerOverviewSpec struct {
	Servers []PgBouncerServer `json:"servers"`
}

type PgBouncerServer struct {
	// PodName represents the pod name
	PodName string `json:"podName"`
	// User which name pgbouncer uses to connect to server.
	User string `json:"user"`
	// Database represents the database name
	Database string `json:"database"`
	// State of the pgbouncer server connection, one of active, idle, used, tested, new, active_cancel, being_canceled.
	State string `json:"state"`
	// IP Address of PostgreSQL server
	Address string `json:"address"`
	// Port of PostgreSQL server
	Port *int32 `json:"port,omitempty"`
	// LocalAddress represents connection start address on local machine
	LocalAddress string `json:"localAddress"`
	// LocalPort represents connection start port on local machine
	LocalPort *int32 `json:"localPort,omitempty"`
	// ConnectTime represents when the connection was made
	ConnectTime *metav1.Time `json:"connectTime,omitempty"`
	// RequestTime represents when last request was issued
	RequestTime *metav1.Time `json:"requestTime,omitempty"`
	// CloseNeeded is 1 if the connection will be closed as soon as possible,
	// because a configuration file reload or DNS update changed the connection information or RECONNECT was issued
	CloseNeeded *int32 `json:"closeNeeded,omitempty"`
	// PTR represents address of internal object for this connection. Used as unique ID
	PTR string `json:"ptr"`
	// Link represents address of client connection the server is paired with.
	Link string `json:"link"`
	// RemotePid represents PID of backend server process.
	RemotePid string `json:"remotePid"`
	// A string with TLS connection information, or empty if not using TLS
	TLS string `json:"tls"`
	// A string containing the application_name set on the linked client connection
	ApplicationName string `json:"applicationName"`
}

// PgBouncerServerOverview is the Schema for the PgBouncerServerOverviews API

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type PgBouncerServerOverview struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec PgBouncerServerOverviewSpec `json:"spec,omitempty"`
}

// PgBouncerServerOverviewList contains a list of PgBouncerServerOverview

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type PgBouncerServerOverviewList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PgBouncerServerOverview `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PgBouncerServerOverview{}, &PgBouncerServerOverviewList{})
}
