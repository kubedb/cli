package api

import (
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/unversioned"
)

const (
	ResourceCodeDeletedDatabase = "ddb"
	ResourceKindDeletedDatabase = "DeletedDatabase"
	ResourceNameDeletedDatabase = "deleted-database"
	ResourceTypeDeletedDatabase = "deleteddatabases"
)

type DeletedDatabase struct {
	unversioned.TypeMeta `json:",inline,omitempty"`
	api.ObjectMeta       `json:"metadata,omitempty"`
	Spec                 DeletedDatabaseSpec   `json:"spec,omitempty"`
	Status               DeletedDatabaseStatus `json:"status,omitempty"`
}

type DeletedDatabaseSpec struct {
	// If true, invoke wipe out operation
	// +optional
	WipeOut bool `json:"wipeOut,omitempty"`
	// If true, invoke recover operation
	// +optional
	Recover bool `json:"recover,omitempty"`
	// Origin to store original database information
	Origin Origin `json:"origin,omitempty"`
}

type Origin struct {
	api.ObjectMeta `json:"metadata,omitempty"`
	// Origin Spec to store original database Spec
	Spec OriginSpec `json:"spec,omitempty"`
}

type OriginSpec struct {
	// Elastic Spec
	// +optional
	Elastic *ElasticSpec `json:"elastic,omitempty"`
	// Postgres Spec
	// +optional
	Postgres *PostgresSpec `json:"postgres,omitempty"`
}

type DeletedDatabasePhase string

const (
	// used for Databases that are deleted
	DeletedDatabasePhaseDeleted DeletedDatabasePhase = "Deleted"
	// used for Databases that are currently deleting
	DeletedDatabasePhaseDeleting DeletedDatabasePhase = "Deleting"
	// used for Databases that are wiped out
	DeletedDatabasePhaseWipedOut DeletedDatabasePhase = "WipedOut"
	// used for Databases that are currently wiping out
	DeletedDatabasePhaseWipingOut DeletedDatabasePhase = "WipingOut"
	// used for Databases that are currently recovering
	DeletedDatabasePhaseRecovering DeletedDatabasePhase = "Recovering"
)

type DeletedDatabaseStatus struct {
	CreationTime *unversioned.Time    `json:"creationTime,omitempty"`
	DeletionTime *unversioned.Time    `json:"deletionTime,omitempty"`
	WipeOutTime  *unversioned.Time    `json:"wipeOutTime,omitempty"`
	Phase        DeletedDatabasePhase `json:"phase,omitempty"`
	Reason       string               `json:"reason,omitempty"`
}

type DeletedDatabaseList struct {
	unversioned.TypeMeta `json:",inline"`
	unversioned.ListMeta `json:"metadata,omitempty"`
	// Items is a list of DeletedDatabase TPR objects
	Items []DeletedDatabase `json:"items,omitempty"`
}
