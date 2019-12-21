/*
Copyright The KubeDB Authors.

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
	ResourceCodeDormantDatabase     = "drmn"
	ResourceKindDormantDatabase     = "DormantDatabase"
	ResourceSingularDormantDatabase = "dormantdatabase"
	ResourcePluralDormantDatabase   = "dormantdatabases"
)

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=dormantdatabases,singular=dormantdatabase,shortName=drmn,categories={datastore,kubedb,appscode,all}
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type DormantDatabase struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Spec              DormantDatabaseSpec   `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	Status            DormantDatabaseStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

type DormantDatabaseSpec struct {
	// If true, invoke wipe out operation
	// +optional
	WipeOut bool `json:"wipeOut,omitempty" protobuf:"varint,1,opt,name=wipeOut"`
	// Origin to store original database information
	Origin Origin `json:"origin" protobuf:"bytes,2,opt,name=origin"`
}

type Origin struct {
	ofst.PartialObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	// Origin Spec to store original database Spec
	Spec OriginSpec `json:"spec" protobuf:"bytes,2,opt,name=spec"`
}

type OriginSpec struct {
	// Elasticsearch Spec
	// +optional
	Elasticsearch *ElasticsearchSpec `json:"elasticsearch,omitempty" protobuf:"bytes,1,opt,name=elasticsearch"`
	// Postgres Spec
	// +optional
	Postgres *PostgresSpec `json:"postgres,omitempty" protobuf:"bytes,2,opt,name=postgres"`
	// MySQL Spec
	// +optional
	MySQL *MySQLSpec `json:"mysql,omitempty" protobuf:"bytes,3,opt,name=mysql"`
	// PerconaXtraDB Spec
	// +optional
	PerconaXtraDB *PerconaXtraDBSpec `json:"perconaxtradb,omitempty" protobuf:"bytes,4,opt,name=perconaxtradb"`
	// MariaDB Spec
	// +optional
	MariaDB *MariaDBSpec `json:"mariadb,omitempty" protobuf:"bytes,5,opt,name=mariadb"`
	// MongoDB Spec
	// +optional
	MongoDB *MongoDBSpec `json:"mongodb,omitempty" protobuf:"bytes,6,opt,name=mongodb"`
	// Redis Spec
	// +optional
	Redis *RedisSpec `json:"redis,omitempty" protobuf:"bytes,7,opt,name=redis"`
	// Memcached Spec
	// +optional
	Memcached *MemcachedSpec `json:"memcached,omitempty" protobuf:"bytes,8,opt,name=memcached"`
	// Etcd Spec
	// +optional
	Etcd *EtcdSpec `json:"etcd,omitempty" protobuf:"bytes,9,opt,name=etcd"`
}

type DormantDatabasePhase string

const (
	// used for Databases that are paused
	DormantDatabasePhasePaused DormantDatabasePhase = "Paused"
	// used for Databases that are currently pausing
	DormantDatabasePhasePausing DormantDatabasePhase = "Pausing"
	// used for Databases that are wiped out
	DormantDatabasePhaseWipedOut DormantDatabasePhase = "WipedOut"
	// used for Databases that are currently wiping out
	DormantDatabasePhaseWipingOut DormantDatabasePhase = "WipingOut"
	// used for Databases that are currently recovering
	DormantDatabasePhaseResuming DormantDatabasePhase = "Resuming"
)

type DormantDatabaseStatus struct {
	PausingTime *metav1.Time         `json:"pausingTime,omitempty" protobuf:"bytes,1,opt,name=pausingTime"`
	WipeOutTime *metav1.Time         `json:"wipeOutTime,omitempty" protobuf:"bytes,2,opt,name=wipeOutTime"`
	Phase       DormantDatabasePhase `json:"phase,omitempty" protobuf:"bytes,3,opt,name=phase,casttype=DormantDatabasePhase"`
	Reason      string               `json:"reason,omitempty" protobuf:"bytes,4,opt,name=reason"`
	// observedGeneration is the most recent generation observed for this resource. It corresponds to the
	// resource's generation, which is updated on mutation by the API Server.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty" protobuf:"varint,5,opt,name=observedGeneration"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type DormantDatabaseList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	// Items is a list of DormantDatabase CRD objects
	Items []DormantDatabase `json:"items,omitempty" protobuf:"bytes,2,rep,name=items"`
}
