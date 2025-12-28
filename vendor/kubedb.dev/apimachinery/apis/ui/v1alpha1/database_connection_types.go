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
	dbapi "kubedb.dev/apimachinery/apis/kubedb/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kmapi "kmodules.xyz/client-go/api/v1"
	ofst "kmodules.xyz/offshoot-api/api/v1"
)

const (
	ResourceKindDatabaseConnection = "DatabaseConnection"
	ResourceDatabaseConnection     = "databaseconnection"
	ResourceDatabaseConnections    = "databaseconnections"
)

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type DatabaseConnection struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DatabaseConnectionSpec `json:"spec,omitempty"`
	Status dbapi.MariaDBStatus    `json:"status,omitempty"`
}

// DatabaseConnectionSpec defines the desired state of DatabaseConnection
type DatabaseConnectionSpec struct {
	Public    []PublicInfo        `json:"public,omitempty"`
	InCluster InClusterConnection `json:"inCluster,omitempty"`

	// Databases already present on the referred database server
	Databases []string `json:"databases,omitempty"`
}

type PublicInfo struct {
	Gateway        []GatewayConnection `json:"gateway,omitempty"`
	ConnectOptions map[string]string   `json:"connectOptions,omitempty"`
}

type GatewayConnection struct {
	*ofst.Gateway `json:",inline"`
	SecretRef     *kmapi.ObjectReference `json:"secretRef,omitempty"`
	CACert        []byte                 `json:"caCert,omitempty"`
}

type InClusterConnection struct {
	Host string `json:"host,omitempty"`
	Port int32  `json:"port,omitempty"`
	// Command for exec-ing into the db pod
	// Example: kubectl exec -it -n default service/mongo-test1  -c mongodb -- bash -c '<the actual command>'
	Exec      string                 `json:"exec,omitempty"`
	SecretRef *kmapi.ObjectReference `json:"secretRef,omitempty"`
	CACert    []byte                 `json:"caCert,omitempty"`

	// Parameters: `username = <username>\n
	// password = <password>\n
	// host = <host>\n
	// database = <database>\n
	// sslmode = REQUIRED`
	//
	// URI: `mongodb+srv://<username>:<password>@<host>:<port>/<database>?authSource=<authSource>&tls=true&replicaSet=arnob`
	//
	// Flags: `mongo "mongodb+srv://<username>:<password>@<host>:<port>/<database>?authSource=<authsource>&replicaSet=arnob" --tls`
	//
	// And some language specific template strings. Like: Java, C#, Go, Python, Javascript, Ruby etc.
	ConnectOptions map[string]string `json:"connectOptions,omitempty"`
}

// DatabaseConnectionList contains a list of DatabaseConnection
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type DatabaseConnectionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DatabaseConnection `json:"items"`
}
