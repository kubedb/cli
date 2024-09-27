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

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type DatabaseConnection struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DatabaseConnectionSpec `json:"spec,omitempty"`
	Status dbapi.MariaDBStatus    `json:"status,omitempty"`
}

// TODO: Need to pass the Type information in the ObjectMeta. For example: MongoDB, MySQL etc.

// DatabaseConnectionSpec defines the desired state of DatabaseConnection
type DatabaseConnectionSpec struct {
	// Public Connections are exposed via Gateway
	Public []PublicConnection `json:"public,omitempty"`

	// Private Connections are in-cluster. Accessible from another pod in the same cluster.
	Private []PrivateConnection `json:"private,omitempty"`

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

//type ConnectOption struct {
//	// username = <username>
//	// password = <password>
//	// host = <host>
//	// database = <database>
//	// sslmode = REQUIRED
//	Parameters []string `json:"parameters,omitempty"`
//
//	// Actual: mongodb+srv://doadmin:show-password@arnob-a013a268.mongo.ondigitalocean.com/admin?authSource=admin&tls=true&replicaSet=arnob
//	// Template: `mongodb+srv://<username>:<password>@<host>:<port>/<database>?authSource=<authSource>&tls=true&replicaSet=arnob`
//	ConnectionString string `json:"connectionString,omitempty"`
//
//	// Actual: mongo "mongodb+srv://doadmin:show-password@private-arnob-aa409eb4.mongo.ondigitalocean.com/admin?authSource=admin&replicaSet=arnob" --tls
//	// Template: `mongo "mongodb+srv://<username>:<password>@<host>:<port>/<database>?authSource=<authsource>&replicaSet=arnob" --tls`
//	Flags string `json:"flags,omitempty"`
//}

type PublicConnection struct {
	*ofst.Gateway `json:",inline"`
	SecretRef     *kmapi.ObjectReference `json:"secretRef,omitempty"`
}

type PrivateConnection struct {
	Host      string                 `json:"host,omitempty"`
	Port      int32                  `json:"port,omitempty"`
	SecretRef *kmapi.ObjectReference `json:"secretRef,omitempty"`
}

// DatabaseConnectionList contains a list of DatabaseConnection
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type DatabaseConnectionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DatabaseConnection `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DatabaseConnection{}, &DatabaseConnectionList{})
}
