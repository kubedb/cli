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
	ResourceKindMySQLOverview = "MySQLOverview"
	ResourceMySQLOverview     = "mysqloverview"
	ResourceMySQLOverviews    = "mysqloverviews"
)

// MySQLOverviewSpec defines the desired state of MySQLOverview
type MySQLOverviewSpec struct {
	Version                string           `json:"version" protobuf:"bytes,1,opt,name=version"`
	Status                 string           `json:"status" protobuf:"bytes,2,opt,name=status"`
	Mode                   string           `json:"mode" protobuf:"bytes,3,opt,name=mode"`
	ConnectionInfo         DBConnectionInfo `json:"connectionsInfo,omitempty" protobuf:"bytes,4,opt,name=connectionsInfo"`
	Credentials            DBCredentials    `json:"credentials,omitempty" protobuf:"bytes,5,opt,name=credentials"`
	MaxConnections         *int32           `json:"maxConnections,omitempty" protobuf:"varint,6,opt,name=maxConnections"`
	MaxUsedConnections     *int32           `json:"maxUsedConnections,omitempty" protobuf:"varint,7,opt,name=maxUsedConnections"`
	Questions              *int32           `json:"questions" protobuf:"varint,8,opt,name=questions"`
	LongQueryTimeThreshold *float64         `json:"longQueryTimeThreshold,omitempty" protobuf:"fixed64,9,opt,name=longQueryTimeThreshold"`
	NumberOfSlowQueries    *int32           `json:"numberOfSlowQueries,omitempty" protobuf:"varint,10,opt,name=numberOfSlowQueries"`
	AbortedClients         *int32           `json:"abortedClients,omitempty" protobuf:"varint,11,opt,name=abortedClients"`
	AbortedConnections     *int32           `json:"abortedConnections,omitempty" protobuf:"varint,12,opt,name=abortedConnections"`
	ThreadsCached          *int32           `json:"threadsCached,omitempty" protobuf:"varint,13,opt,name=threadsCached"`
	ThreadsConnected       *int32           `json:"threadsConnected,omitempty" protobuf:"varint,14,opt,name=threadsConnected"`
	ThreadsCreated         *int32           `json:"threadsCreated,omitempty" protobuf:"varint,15,opt,name=threadsCreated"`
	ThreadsRunning         *int32           `json:"threadsRunning,omitempty" protobuf:"varint,16,opt,name=threadsRunning"`
}

// MySQLOverview is the Schema for the mysqloverviews API

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type MySQLOverview struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec   MySQLOverviewSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	Status api.MySQLStatus   `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// MySQLOverviewList contains a list of MySQLOverview

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type MySQLOverviewList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Items           []MySQLOverview `json:"items" protobuf:"bytes,2,rep,name=items"`
}

func init() {
	SchemeBuilder.Register(&MySQLOverview{}, &MySQLOverviewList{})
}
