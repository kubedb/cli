package api

import (
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/unversioned"
)

const (
	ResourceKindElastic = "Elastic"
	ResourceNameElastic = "elastic"
	ResourceTypeElastic = "elastics"
)

// Elastic defines a Elasticsearch database.
type Elastic struct {
	unversioned.TypeMeta `json:",inline,omitempty"`
	api.ObjectMeta       `json:"metadata,omitempty"`
	Spec                 ElasticSpec   `json:"spec,omitempty"`
	Status               ElasticStatus `json:"status,omitempty"`
}

type ElasticSpec struct {
	// Version of Elasticsearch to be deployed.
	Version string `json:"version,omitempty"`
	// Number of instances to deploy for a Elasticsearch database.
	Replicas int32 `json:"replicas,omitempty"`
	// Storage spec to specify how storage shall be used.
	Storage *StorageSpec `json:"storage,omitempty"`
	// ServiceAccountName is the name of the ServiceAccount to use to run the
	// Prometheus Pods.
	ServiceAccountName string `json:"serviceAccountName,omitempty"`
	// NodeSelector is a selector which must be true for the pod to fit on a node
	// +optional
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
	// Init is used to initialize database
	// +optional
	Init *InitSpec `json:"init,omitempty"`
	// BackupSchedule spec to specify how database backup will be taken
	// +optional
	BackupSchedule *BackupScheduleSpec `json:"backupSchedule,omitempty"`
	// If DoNotDelete is true, controller will prevent to delete this Elastic object.
	// Controller will create same Elastic object and ignore other process.
	// +optional
	DoNotDelete bool `json:"doNotDelete,omitempty"`
}

type ElasticStatus struct {
	CreationTime   *unversioned.Time `json:"creationTime,omitempty"`
	DatabaseStatus `json:",inline,omitempty"`
	Reason         string `json:"reason,omitempty"`
}

type ElasticList struct {
	unversioned.TypeMeta `json:",inline"`
	unversioned.ListMeta `json:"metadata,omitempty"`
	// Items is a list of Elastic TPR objects
	Items []*Elastic `json:"items,omitempty"`
}
