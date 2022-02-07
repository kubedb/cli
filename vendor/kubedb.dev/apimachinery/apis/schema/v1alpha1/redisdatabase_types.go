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
	kmapi "kmodules.xyz/client-go/api/v1"
)

const (
	ResourceKindRedisDatabase = "RedisDatabase"
	ResourceRedisDatabase     = "redisdatabase"
	ResourceRedisDatabases    = "redisdatabases"
)

// RedisDatabaseSpec defines the desired state of RedisDatabase
type RedisDatabaseSpec struct {
	// Database defines various configuration options for a database
	Database RedisDatabaseInfo `json:"database"`

	// VaultRef refers to a KubeVault managed vault server
	VaultRef kmapi.ObjectReference `json:"vaultRef"`

	// AccessPolicy contains the serviceAccount details and TTL values of the vault-created secret
	AccessPolicy VaultSecretEngineRole `json:"accessPolicy"`

	// Init contains info about the init script or snapshot info
	// +optional
	Init *InitSpec `json:"init,omitempty"`

	// DeletionPolicy controls the delete operation for database
	// +optional
	// +kubebuilder:default:="Delete"
	DeletionPolicy DeletionPolicy `json:"deletionPolicy,omitempty"`
}

type RedisDatabaseInfo struct {
	// ServerRef refers to a KubeDB managed database instance
	ServerRef kmapi.ObjectReference `json:"serverRef"`

	// DatabaseConfig defines various configuration options for a database
	Config RedisDatabaseConfiguration `json:"config"`
}

type RedisDatabaseConfiguration struct {
}

// RedisDatabase is the Schema for the redisdatabases API

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="DB_SERVER",type="string",JSONPath=".spec.database.serverRef.name"
// +kubebuilder:printcolumn:name="DB_NAME",type="string",JSONPath=".spec.database.config.name"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type RedisDatabase struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RedisDatabaseSpec `json:"spec,omitempty"`
	Status DatabaseStatus    `json:"status,omitempty"`
}

// RedisDatabaseList contains a list of RedisDatabase

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
type RedisDatabaseList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RedisDatabase `json:"items"`
}

func init() {
	SchemeBuilder.Register(&RedisDatabase{}, &RedisDatabaseList{})
}
