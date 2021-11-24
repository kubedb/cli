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
	api "kubedb.dev/apimachinery/apis/kubedb/v1alpha2"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ResourceKindMongoDBOverview = "MongoDBOverview"
	ResourceMongoDBOverview     = "mongodboverview"
	ResourceMongoDBOverviews    = "mongodboverviews"
)

// MongoDBOverviewSpec defines the desired state of MongoDBOverview
type MongoDBOverviewSpec struct {
	Name           string                  `json:"name" protobuf:"bytes,1,opt,name=name"`
	Namespace      string                  `json:"namespace" protobuf:"bytes,2,opt,name=namespace"`
	Version        string                  `json:"version" protobuf:"bytes,3,opt,name=version"`
	Type           MongoDBMode             `json:"type" protobuf:"bytes,4,opt,name=type,casttype=MongoDBMode"`
	Status         api.DatabasePhase       `json:"status" protobuf:"bytes,5,opt,name=status,casttype=DBStatus"`
	ConnectionInfo *DBConnectionInfo       `json:"connectionInfo" protobuf:"bytes,6,opt,name=connectionInfo"`
	Credentials    *DBCredentials          `json:"credentials" protobuf:"bytes,7,opt,name=credentials"`
	Connections    *MongoDBConnectionsInfo `json:"connections" protobuf:"bytes,8,opt,name=connections"`
	DBStats        *MongoDBDatabaseStats   `json:"dbStats" protobuf:"bytes,9,opt,name=dbStats"`
	ShardsInfo     *MongoDBShardsInfo      `json:"shardsInfo,omitempty" protobuf:"bytes,10,opt,name=shardsInfo"`
	ReplicaSetInfo *MongoDBReplicaSetInfo  `json:"replicaSetInfo,omitempty" protobuf:"bytes,11,opt,name=replicaSetInfo"`
	Resources      corev1.ResourceList     `json:"resources" protobuf:"bytes,12,rep,name=resources,casttype=k8s.io/api/core/v1.ResourceList,castkey=k8s.io/api/core/v1.ResourceName"`
	SecurityReport *SecurityReport         `json:"securityReport" protobuf:"bytes,13,opt,name=securityReport"`
}

type MongoDBDatabaseStats struct {
	TotalCollections int32 `json:"totalCollections" protobuf:"varint,1,opt,name=totalCollections"`
	DataSize         int64 `json:"dataSize" protobuf:"varint,2,opt,name=dataSize"`
	TotalIndexes     int32 `json:"totalIndexes" protobuf:"varint,3,opt,name=totalIndexes"`
	IndexSize        int64 `json:"indexSize" protobuf:"varint,4,opt,name=indexSize"`
}

type MongoDBConnectionsInfo struct {
	CurrentConnections   int32 `json:"currentConnections" protobuf:"varint,1,opt,name=currentConnections"`
	TotalConnections     int32 `json:"totalConnections" protobuf:"varint,2,opt,name=totalConnections"`
	AvailableConnections int32 `json:"availableConnections" protobuf:"varint,3,opt,name=availableConnections"`
	ActiveConnections    int32 `json:"activeConnections" protobuf:"varint,4,opt,name=activeConnections"`
}

type MongoDBReplicaSetInfo struct {
	NumberOfReplicas int32 `json:"numberOfReplicas" protobuf:"varint,1,opt,name=numberOfReplicas"`
}

type SecurityReport struct {
	TLSSecured        bool `json:"tlsSecured" protobuf:"varint,1,opt,name=tlsSecured"`
	BackupEnabled     bool `json:"backupEnabled" protobuf:"varint,2,opt,name=backupEnabled"`
	MonitoringEnabled bool `json:"monitoringEnabled" protobuf:"varint,3,opt,name=monitoringEnabled"`
}

type DBConnectionInfo struct {
	URL  string `json:"url,omitempty" protobuf:"bytes,1,opt,name=url"`
	Port int32  `json:"port,omitempty" protobuf:"varint,2,opt,name=port"`
}

type MongoDBShardsInfo struct {
	NumberOfShards    int32 `json:"numberOfShards" protobuf:"varint,1,opt,name=numberOfShards"`
	ReplicasPerShards int32 `json:"replicasPerShards" protobuf:"varint,2,opt,name=replicasPerShards"`
	NumberOfChunks    int32 `json:"numberOfChunks" protobuf:"varint,3,opt,name=numberOfChunks"`
	BalancerEnabled   bool  `json:"balancerEnabled" protobuf:"varint,4,opt,name=balancerEnabled"`
	ChunksBalanced    bool  `json:"chunksBalanced" protobuf:"varint,5,opt,name=chunksBalanced"`
}

type DBCredentials struct {
	Username    string `json:"username,omitempty" protobuf:"bytes,1,opt,name=username"`
	Password    string `json:"password,omitempty" protobuf:"bytes,2,opt,name=password"`
	Certificate string `json:"certificate,omitempty" protobuf:"bytes,3,opt,name=certificate"`
}

// MongoDBOverview is the Schema for the MongoDBOverviews API

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type MongoDBOverview struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec   MongoDBOverviewSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	Status api.MongoDBStatus   `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// MongoDBOverviewList contains a list of MongoDBOverview

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type MongoDBOverviewList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Items           []MongoDBOverview `json:"items" protobuf:"bytes,2,rep,name=items"`
}

func init() {
	SchemeBuilder.Register(&MongoDBOverview{}, &MongoDBOverviewList{})
}
