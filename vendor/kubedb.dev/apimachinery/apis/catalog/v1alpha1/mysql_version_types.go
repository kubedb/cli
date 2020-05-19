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

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

const (
	ResourceCodeMySQLVersion     = "myversion"
	ResourceKindMySQLVersion     = "MySQLVersion"
	ResourceSingularMySQLVersion = "mysqlversion"
	ResourcePluralMySQLVersion   = "mysqlversions"
)

// MySQLVersion defines a MySQL database version.

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=updateStatus
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=mysqlversions,singular=mysqlversion,scope=Cluster,shortName=myversion,categories={datastore,kubedb,appscode}
// +kubebuilder:printcolumn:name="Version",type="string",JSONPath=".spec.version"
// +kubebuilder:printcolumn:name="DB_IMAGE",type="string",JSONPath=".spec.db.image"
// +kubebuilder:printcolumn:name="Deprecated",type="boolean",JSONPath=".spec.deprecated"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type MySQLVersion struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Spec              MySQLVersionSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
}

// MySQLVersionSpec is the spec for postgres version
type MySQLVersionSpec struct {
	// Version
	Version string `json:"version" protobuf:"bytes,1,opt,name=version"`
	// Database Image
	DB MySQLVersionDatabase `json:"db" protobuf:"bytes,2,opt,name=db"`
	// Exporter Image
	Exporter MySQLVersionExporter `json:"exporter" protobuf:"bytes,3,opt,name=exporter"`
	// Tools Image
	Tools MySQLVersionTools `json:"tools" protobuf:"bytes,4,opt,name=tools"`
	// ReplicationModeDetector Image
	ReplicationModeDetector MySQLVersionReplicationModeDetector `json:"replicationModeDetector" protobuf:"bytes,5,opt,name=replicationModeDetector"`
	// Deprecated versions usable but regarded as obsolete and best avoided, typically due to having been superseded.
	// +optional
	Deprecated bool `json:"deprecated,omitempty" protobuf:"varint,6,opt,name=deprecated"`
	// Init container Image
	InitContainer MySQLVersionInitContainer `json:"initContainer" protobuf:"bytes,7,opt,name=initContainer"`
	// PSP names
	PodSecurityPolicies MySQLVersionPodSecurityPolicy `json:"podSecurityPolicies" protobuf:"bytes,8,opt,name=podSecurityPolicies"`
	//upgrade constraints
	UpgradeConstraints MySQLUpgradeConstraints `json:"upgradeConstraints" protobuf:"bytes,9,opt,name=upgradeConstraints"`
}

// MySQLVersionDatabase is the MySQL Database image
type MySQLVersionDatabase struct {
	Image string `json:"image" protobuf:"bytes,1,opt,name=image"`
}

// MySQLVersionExporter is the image for the MySQL exporter
type MySQLVersionExporter struct {
	Image string `json:"image" protobuf:"bytes,1,opt,name=image"`
}

// MySQLVersionReplicationModeDetector is the image for the MySQL replication mode detector
type MySQLVersionReplicationModeDetector struct {
	Image string `json:"image" protobuf:"bytes,1,opt,name=image"`
}

// MySQLVersionTools is the image for the postgres tools
type MySQLVersionTools struct {
	Image string `json:"image" protobuf:"bytes,1,opt,name=image"`
}

// MySQLVersionInitContainer is the Elasticsearch Container initializer
type MySQLVersionInitContainer struct {
	Image string `json:"image" protobuf:"bytes,1,opt,name=image"`
}

// MySQLVersionPodSecurityPolicy is the MySQL pod security policies
type MySQLVersionPodSecurityPolicy struct {
	DatabasePolicyName string `json:"databasePolicyName" protobuf:"bytes,1,opt,name=databasePolicyName"`
}

type MySQLUpgradeConstraints struct {
	// List of all accepted version for modification request
	Whitelist Whitelist `json:"whitelist,omitempty" protobuf:"bytes,1,opt,name=whitelist"`
	// List of all rejected version for modification request
	Blacklist Blacklist `json:"blacklist,omitempty" protobuf:"bytes,2,opt,name=blacklist"`
}

type Whitelist struct {
	// List of all accepted version for standalone modification request. empty indicates all accepted
	Standalone []string `json:"standalone,omitempty" protobuf:"bytes,1,opt,name=standalone"`
	// List of all accepted version for groupReplication modification request. empty indicates all accepted
	GroupReplication []string `json:"groupReplication,omitempty" protobuf:"bytes,2,opt,name=groupReplication"`
}

type Blacklist struct {
	// List of all rejected version for standalone modification request
	Standalone []string `json:"standalone,omitempty" protobuf:"bytes,1,opt,name=standalone"`
	// List of all rejected version for groupReplication modification request
	GroupReplication []string `json:"groupReplication,omitempty" protobuf:"bytes,2,opt,name=groupReplication"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MySQLVersionList is a list of MySQLVersions
type MySQLVersionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	// Items is a list of MySQLVersion CRD objects
	Items []MySQLVersion `json:"items,omitempty" protobuf:"bytes,2,rep,name=items"`
}
