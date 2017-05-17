package api

import (
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/unversioned"
)

const (
	ResourceCodePostgres = "pg"
	ResourceKindPostgres = "Postgres"
	ResourceNamePostgres = "postgres"
	ResourceTypePostgres = "postgreses"
)

// Postgres defines a Postgres database.
type Postgres struct {
	unversioned.TypeMeta `json:",inline,omitempty"`
	api.ObjectMeta       `json:"metadata,omitempty"`
	Spec                 PostgresSpec   `json:"spec,omitempty"`
	Status               PostgresStatus `json:"status,omitempty"`
}

type PostgresSpec struct {
	// Version of Postgres to be deployed.
	Version string `json:"version,omitempty"`
	// Storage spec to specify how storage shall be used.
	Storage *StorageSpec `json:"storage,omitempty"`
	// Database authentication secret
	DatabaseSecret *api.SecretVolumeSource `json:"databaseSecret,omitempty"`
	// NodeSelector is a selector which must be true for the pod to fit on a node
	// +optional
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
	// Init is used to initialize database
	// +optional
	Init *InitSpec `json:"init,omitempty"`
	// BackupSchedule spec to specify how database backup will be taken
	// +optional
	BackupSchedule *BackupScheduleSpec `json:"backupSchedule,omitempty"`
	// If DoNotDelete is true, controller will prevent to delete this Postgres object.
	// Controller will create same Postgres object and ignore other process.
	// +optional
	DoNotDelete bool `json:"doNotDelete,omitempty"`
}

type PostgresStatus struct {
	CreationTime *unversioned.Time `json:"creationTime,omitempty"`
	Phase        DatabasePhase     `json:"phase,omitempty"`
	Reason       string            `json:"reason,omitempty"`
}

type PostgresList struct {
	unversioned.TypeMeta `json:",inline"`
	unversioned.ListMeta `json:"metadata,omitempty"`
	// Items is a list of Postgres TPR objects
	Items []Postgres `json:"items,omitempty"`
}
