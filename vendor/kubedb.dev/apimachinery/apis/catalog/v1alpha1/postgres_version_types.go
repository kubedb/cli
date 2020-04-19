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
	ResourceCodePostgresVersion     = "pgversion"
	ResourceKindPostgresVersion     = "PostgresVersion"
	ResourceSingularPostgresVersion = "postgresversion"
	ResourcePluralPostgresVersion   = "postgresversions"
)

// PostgresVersion defines a Postgres database version.

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=updateStatus
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=postgresversions,singular=postgresversion,scope=Cluster,shortName=pgversion,categories={datastore,kubedb,appscode}
// +kubebuilder:printcolumn:name="Version",type="string",JSONPath=".spec.version"
// +kubebuilder:printcolumn:name="DB_IMAGE",type="string",JSONPath=".spec.db.image"
// +kubebuilder:printcolumn:name="Deprecated",type="boolean",JSONPath=".spec.deprecated"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type PostgresVersion struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Spec              PostgresVersionSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
}

// PostgresVersionSpec is the spec for postgres version
type PostgresVersionSpec struct {
	// Version
	Version string `json:"version" protobuf:"bytes,1,opt,name=version"`
	// Database Image
	DB PostgresVersionDatabase `json:"db" protobuf:"bytes,2,opt,name=db"`
	// Exporter Image
	Exporter PostgresVersionExporter `json:"exporter" protobuf:"bytes,3,opt,name=exporter"`
	// Tools Image
	Tools PostgresVersionTools `json:"tools" protobuf:"bytes,4,opt,name=tools"`
	// Deprecated versions usable but regarded as obsolete and best avoided, typically due to having been superseded.
	// +optional
	Deprecated bool `json:"deprecated,omitempty" protobuf:"varint,5,opt,name=deprecated"`
	// PSP names
	PodSecurityPolicies PostgresVersionPodSecurityPolicy `json:"podSecurityPolicies" protobuf:"bytes,6,opt,name=podSecurityPolicies"`
}

// PostgresVersionDatabase is the Postgres Database image
type PostgresVersionDatabase struct {
	Image string `json:"image" protobuf:"bytes,1,opt,name=image"`
}

// PostgresVersionExporter is the image for the Postgres exporter
type PostgresVersionExporter struct {
	Image string `json:"image" protobuf:"bytes,1,opt,name=image"`
}

// PostgresVersionTools is the image for the postgres tools
type PostgresVersionTools struct {
	Image string `json:"image" protobuf:"bytes,1,opt,name=image"`
}

// PostgresVersionPodSecurityPolicy is the Postgres pod security policies
type PostgresVersionPodSecurityPolicy struct {
	DatabasePolicyName string `json:"databasePolicyName" protobuf:"bytes,1,opt,name=databasePolicyName"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PostgresVersionList is a list of PostgresVersions
type PostgresVersionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	// Items is a list of PostgresVersion CRD objects
	Items []PostgresVersion `json:"items,omitempty" protobuf:"bytes,2,rep,name=items"`
}
