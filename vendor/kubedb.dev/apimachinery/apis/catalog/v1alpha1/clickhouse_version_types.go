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
	ResourceCodeClickHouseVersion     = "chversion"
	ResourceKindClickHouseVersion     = "ClickHouseVersion"
	ResourceSingularClickHouseVersion = "clickhouseversion"
	ResourcePluralClickHouseVersion   = "clickhouseversions"
)

// ClickHouseVersion defines a ClickHouse database version.

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=updateStatus
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=clickhouseversions,singular=clickhouseversion,scope=Cluster,shortName=chversion,categories={catalog,kubedb,appscode}
// +kubebuilder:printcolumn:name="Version",type="string",JSONPath=".spec.version"
// +kubebuilder:printcolumn:name="DB_IMAGE",type="string",JSONPath=".spec.db.image"
// +kubebuilder:printcolumn:name="Deprecated",type="boolean",JSONPath=".spec.deprecated"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type ClickHouseVersion struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ClickHouseVersionSpec   `json:"spec,omitempty"`
	Status ClickHouseVersionStatus `json:"status,omitempty"`
}

// ClickHouseVersionSpec defines the desired state of ClickHouseVersion
type ClickHouseVersionSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Version
	Version string `json:"version"`

	// Database Image
	DB ClickHouseVersionDatabase `json:"db"`

	// Database Image
	InitContainer ClickHouseInitContainer `json:"initContainer"`

	// ClickHouse Keeper Image
	ClickHouseKeeper ClickHouseKeeperContainer `json:"clickHouseKeeper"`

	// SecurityContext is for the additional config for the DB container
	// +optional
	SecurityContext SecurityContext `json:"securityContext"`

	// +optional
	UI []ChartInfo `json:"ui,omitempty"`
}

// ClickHouseVersionDatabase is the ClickHouse Database image
type ClickHouseVersionDatabase struct {
	Image string `json:"image"`
}

// ClickHouseInitContainer is the ClickHouse init Container image
type ClickHouseInitContainer struct {
	Image string `json:"image"`
}

// ClickHouseKeeperContainer is the ClickHouse keeper Container image
type ClickHouseKeeperContainer struct {
	Image string `json:"image"`
}

// ClickHouseVersionStatus defines the observed state of ClickHouseVersion
type ClickHouseVersionStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClickHouseVersionList contains a list of ClickHouseVersion
type ClickHouseVersionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ClickHouseVersion `json:"items"`
}
