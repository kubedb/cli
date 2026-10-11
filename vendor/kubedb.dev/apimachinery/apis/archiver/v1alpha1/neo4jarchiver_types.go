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
	storageapi "kubestash.dev/apimachinery/apis/storage/v1alpha1"
)

const (
	ResourceKindNeo4jArchiver     = "Neo4jArchiver"
	ResourceSingularNeo4jArchiver = "neo4jarchiver"
	ResourcePluralNeo4jArchiver   = "neo4jarchivers"
)

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=neo4jarchivers,singular=neo4jarchiver,shortName=neojarchiver,categories={archiver,kubedb,appscode}
// +kubebuilder:subresource:status
type Neo4jArchiver struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   Neo4jArchiverSpec   `json:"spec,omitempty"`
	Status Neo4jArchiverStatus `json:"status,omitempty"`
}

// Neo4jArchiverSpec defines the desired state of Neo4jArchiver
type Neo4jArchiverSpec struct {
	// Databases define which Neo4j databases are allowed to consume this archiver
	Databases *dbapi.AllowedConsumers `json:"databases"`
	// Pause defines if the backup process should be paused or not
	// +optional
	Pause bool `json:"pause,omitempty"`
	// RetentionPolicy field is the RetentionPolicy of the backupConfiguration's backend
	// +optional
	RetentionPolicy *kmapi.ObjectReference `json:"retentionPolicy"`
	// FullBackup defines the sessionConfig of the fullBackup
	// This options will eventually go to the full-backup job's yaml
	// +optional
	FullBackup *Neo4jFullBackupOptions `json:"fullBackup"`
	// DifferentialBackup defines the configuration for differential backup
	// +optional
	DifferentialBackup *Neo4jDifferentialBackupOptions `json:"differentialBackup"`
	// ManifestBackup defines the sessionConfig of the manifestBackup
	// This options will eventually go to the manifest-backup job's yaml
	// +optional
	ManifestBackup *ManifestBackupOptions `json:"manifestBackup"`
	// +optional
	EncryptionSecret *kmapi.ObjectReference `json:"encryptionSecret"`
	// BackupStorage is the backend storageRef of the BackupConfiguration
	// +optional
	BackupStorage *BackupStorage `json:"backupStorage"`
	// DeletionPolicy defines the created repository's deletionPolicy
	// +optional
	DeletionPolicy *storageapi.BackupConfigDeletionPolicy `json:"deletionPolicy"`
}

// Neo4jArchiverStatus defines the observed state of Neo4jArchiver
type Neo4jArchiverStatus struct {
	// Specifies the information of all the databases managed by this archiver
	// +optional
	DatabaseRefs []ArchiverDatabaseRef `json:"databaseRefs,omitempty"`
}

// +kubebuilder:object:root=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type Neo4jArchiverList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Neo4jArchiver `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Neo4jArchiver{}, &Neo4jArchiverList{})
}
