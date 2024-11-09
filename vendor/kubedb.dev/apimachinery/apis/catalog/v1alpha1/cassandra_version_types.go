/*
Copyright 2024.

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
)

const (
	ResourceCodeCassandraVersion     = "casversion"
	ResourceKindCassandraVersion     = "CassandraVersion"
	ResourceSingularCassandraVersion = "cassandraversion"
	ResourcePluralCassandraVersion   = "cassandraversions"
)

// CassandraVersion defines a Cassandra database version.

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=updateStatus
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=cassandraversions,singular=cassandraversion,scope=Cluster,shortName=casversion,categories={catalog,kubedb,appscode}
// +kubebuilder:printcolumn:name="Version",type="string",JSONPath=".spec.version"
// +kubebuilder:printcolumn:name="DB_IMAGE",type="string",JSONPath=".spec.db.image"
// +kubebuilder:printcolumn:name="Deprecated",type="boolean",JSONPath=".spec.deprecated"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type CassandraVersion struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CassandraVersionSpec   `json:"spec,omitempty"`
	Status CassandraVersionStatus `json:"status,omitempty"`
}

// CassandraVersionSpec defines the desired state of CassandraVersion
type CassandraVersionSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Version
	Version string `json:"version"`

	// Database Image
	DB CassandraVersionDatabase `json:"db"`

	// Exporter Image
	Exporter CassandraVersionExporter `json:"exporter"`

	// Database Image
	InitContainer CassandraInitContainer `json:"initContainer"`

	// SecurityContext is for the additional config for the DB container
	// +optional
	SecurityContext SecurityContext `json:"securityContext"`

	// +optional
	UI []ChartInfo `json:"ui,omitempty"`
}

// CassandraVersionExporter is the image for the Cassandra exporter
type CassandraVersionExporter struct {
	Image string `json:"image"`
}

// CassandraVersionDatabase is the Cassandra Database image
type CassandraVersionDatabase struct {
	Image string `json:"image"`
}

// CassandraInitContainer is the Cassandra init Container image
type CassandraInitContainer struct {
	Image string `json:"image"`
}

// CassandraVersionStatus defines the observed state of CassandraVersion
type CassandraVersionStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CassandraVersionList contains a list of CassandraVersion
type CassandraVersionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CassandraVersion `json:"items"`
}
