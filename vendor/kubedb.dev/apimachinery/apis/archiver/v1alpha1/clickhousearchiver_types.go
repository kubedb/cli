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
	ResourceKindClickHouseArchiver     = "ClickHouseArchiver"
	ResourceSingularClickHouseArchiver = "clickhousearchiver"
	ResourcePluralClickHouseArchiver   = "clickhousearchivers"
)

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=clickhousearchivers,singular=clickhousearchiver,shortName=charchiver,categories={archiver,kubedb,appscode}
// +kubebuilder:subresource:status
type ClickHouseArchiver struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ClickHouseArchiverSpec   `json:"spec,omitempty"`
	Status ClickHouseArchiverStatus `json:"status,omitempty"`
}

// ClickHouseArchiverSpec defines the desired state of ClickHouseArchiver
type ClickHouseArchiverSpec struct {
	// Databases define which ClickHouse databases are allowed to consume this archiver
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
	FullBackup *ClickHouseFullBackupOptions `json:"fullBackup"`
	// IncrementalBackup defines the configuration for incremental backup
	// +optional
	IncrementalBackup *ClickHouseIncrementalBackupOptions `json:"incrementalBackup"`
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

// ClickHouseArchiverStatus defines the observed state of ClickHouseArchiver
type ClickHouseArchiverStatus struct {
	// Specifies the information of all the databases managed by this archiver
	// +optional
	DatabaseRefs []ArchiverDatabaseRef `json:"databaseRefs,omitempty"`
}

// +kubebuilder:object:root=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ClickHouseArchiverList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ClickHouseArchiver `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ClickHouseArchiver{}, &ClickHouseArchiverList{})
}
