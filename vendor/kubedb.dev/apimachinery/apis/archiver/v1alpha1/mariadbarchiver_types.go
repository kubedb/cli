/*
Copyright 2022.

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
	ResourceKindMariaDBArchiver     = "MariaDBArchiver"
	ResourceSingularMariaDBArchiver = "mariadbarchiver"
	ResourcePluralMariaDBArchiver   = "mariadbarchivers"
)

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=mariadbarchivers,singular=mariadbarchiver,shortName=mdarchiver,categories={archiver,kubedb,appscode}
// +kubebuilder:subresource:status
type MariaDBArchiver struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MariaDBArchiverSpec   `json:"spec,omitempty"`
	Status MariaDBArchiverStatus `json:"status,omitempty"`
}

// MariaDBArchiverSpec defines the desired state of MariaDBArchiver
type MariaDBArchiverSpec struct {
	// Databases define which MariaDB databases are allowed to consume this archiver
	Databases *dbapi.AllowedConsumers `json:"databases"`
	// Pause defines if the backup process should be paused or not
	// +optional
	Pause bool `json:"pause,omitempty"`
	// RetentionPolicy refers to a RetentionPolicy CR which defines how to cleanup the old Snapshots
	// +optional
	RetentionPolicy *kmapi.ObjectReference `json:"retentionPolicy"`
	// FullBackup defines the session configuration for the full backup
	// +optional
	FullBackup *FullBackupOptions `json:"fullBackup"`
	// LogBackup defines the sidekick configuration for the log backup
	// +optional
	LogBackup *LogBackupOptions `json:"logBackup"`
	// ManifestBackup defines the session configuration for the manifest backup
	// This options will eventually go to the manifest-backup job's yaml
	// +optional
	ManifestBackup *ManifestBackupOptions `json:"manifestBackup"`
	// EncryptionSecret refers to the Secret containing the encryption key used to encode backed-up data.
	// +optional
	EncryptionSecret *kmapi.ObjectReference `json:"encryptionSecret"`
	// BackupStorage holds the storage information for storing backup data
	// +optional
	BackupStorage *BackupStorage `json:"backupStorage"`
	// DeletionPolicy defines the DeletionPolicy for the backup repository
	// +optional
	DeletionPolicy *storageapi.BackupConfigDeletionPolicy `json:"deletionPolicy"`
}

// MariaDBArchiverStatus defines the observed state of MariaDBArchiver
type MariaDBArchiverStatus struct {
	// Specifies the information of all the databases managed by this archiver
	// +optional
	DatabaseRefs []ArchiverDatabaseRef `json:"databaseRefs,omitempty"`
}

// +kubebuilder:object:root=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type MariaDBArchiverList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MariaDBArchiver `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MariaDBArchiver{}, &MariaDBArchiverList{})
}
