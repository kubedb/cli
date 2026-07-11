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
	ofst "kmodules.xyz/offshoot-api/api/v1"
)

const (
	ResourceKindMariaDBMigration     = "MariaDBMigration"
	ResourceSingularMariaDBMigration = "mariadbmigration"
	ResourcePluralMariaDBMigrations  = "mariadbmigrations"
)

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=mariadbmigrations,singular=mariadbmigration,shortName=mrmig,categories={kubedb,appscode,all}
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced
// +kubebuilder:printcolumn:name="Phase",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Stage",type="string",JSONPath=".status.progress.info.Stage"
// +kubebuilder:printcolumn:name="Lag",type="string",JSONPath=".status.progress.info.Lag"
// +kubebuilder:printcolumn:name="Progress",type="string",JSONPath=".status.progress.info.Progress"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type MariaDBMigration struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitzero"`

	// spec defines the desired state of MariaDBMigration
	// +required
	Spec MariaDBMigrationSpec `json:"spec"`

	// status defines the observed state of MariaDBMigration.
	// It reuses the shared MigrationStatus so that the Migration duck type can
	// project it and the operator's status patches replay onto it unchanged.
	// +optional
	Status MigrationStatus `json:"status,omitzero"`
}

// MariaDBMigrationSpec defines the desired state of MariaDBMigration
type MariaDBMigrationSpec struct {
	// Source defines the source MariaDB database configuration
	Source MariaDBSource `json:"source"`

	// Target defines the target MariaDB database configuration
	Target MariaDBTarget `json:"target"`

	// JobDefaults specifies default settings for migration jobs
	// +optional
	JobDefaults *JobDefaults `json:"jobDefaults,omitempty"`

	// JobTemplate specifies runtime configurations for the migration Job
	// +optional
	JobTemplate *ofst.PodTemplateSpec `json:"jobTemplate,omitempty"`
}

// MariaDBMigrationList contains a list of MariaDBMigration

// +kubebuilder:object:root=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type MariaDBMigrationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitzero"`
	Items           []MariaDBMigration `json:"items"`
}

type MariaDBSource struct {
	// ConnectionInfo refers to the source MariaDB database connection information.
	ConnectionInfo *ConnectionInfo `yaml:"connectionInfo" json:"connectionInfo"`
	Schema         *MySQLSchema    `yaml:"schema" json:"schema,omitempty"`
	Snapshot       *MySQLSnapshot  `yaml:"snapshot" json:"snapshot,omitempty"`
	Streaming      *MySQLStreaming `yaml:"streaming" json:"streaming,omitempty"`
}

type MariaDBTarget struct {
	// ConnectionInfo refers to the target MariaDB database connection information.
	ConnectionInfo *ConnectionInfo `yaml:"connectionInfo" json:"connectionInfo"`
}

func init() {
	SchemeBuilder.Register(&MariaDBMigration{}, &MariaDBMigrationList{})
}
