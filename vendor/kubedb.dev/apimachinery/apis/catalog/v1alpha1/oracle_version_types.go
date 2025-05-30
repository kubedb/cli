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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ResourceCodeOracleVersion     = "oraversion"
	ResourceKindOracleVersion     = "OracleVersion"
	ResourceSingularOracleVersion = "oracleversion"
	ResourcePluralOracleVersion   = "oracleversions"
)

// OracleVersion defines a Oracle database version.

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=updateStatus
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=oracleversions,singular=oracleversion,scope=Cluster,shortName=oraversion,categories={catalog,kubedb,appscode}
// +kubebuilder:printcolumn:name="Version",type="string",JSONPath=".spec.version"
// +kubebuilder:printcolumn:name="Distribution",type="string",JSONPath=".spec.distribution"
// +kubebuilder:printcolumn:name="DB_IMAGE",type="string",JSONPath=".spec.db.image"
// +kubebuilder:printcolumn:name="Deprecated",type="boolean",JSONPath=".spec.deprecated"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type OracleVersion struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              OracleVersionSpec `json:"spec,omitempty"`
}

// OracleVersionSpec is the spec for oracle version
type OracleVersionSpec struct {
	// Version
	Version string `json:"version"`
	// Distribution
	Distribution OracleDistro `json:"distribution,omitempty"`
	// init container image
	InitContainer OracleVersionInitContainer `json:"initContainer,omitempty"`
	// Database Image
	DB OracleVersionDatabase `json:"db"`
	// Exporter Image
	Exporter OracleVersionExporter `json:"exporter"`
	// Coordinator Image
	Coordinator OracleVersionCoordinator `json:"coordinator,omitempty"`
	// DataGuard Images
	DataGuard OracleDataGuard `json:"dataGuard,omitempty"`

	// Deprecated versions usable but regarded as obsolete and best avoided, typically due to having been superseded.
	// +optional
	Deprecated bool `json:"deprecated,omitempty"`

	// SecurityContext is for the additional config for oracle DB container
	// +optional
	SecurityContext OracleSecurityContext `json:"securityContext"`

	// update constraints
	UpdateConstraints UpdateConstraints `json:"updateConstraints,omitempty"`

	// +optional
	UI []ChartInfo `json:"ui,omitempty"`
}

// OracleObserver defines images for observer

type OracleDataGuard struct {
	InitContainer OracleVersionInitContainer `json:"initContainer,omitempty"`
	// Database Image
	Observer OracleVersionDatabase `json:"observer,omitempty"`
}

// OracleVersionInitContainer is the Oracle init container image
type OracleVersionInitContainer struct {
	Image string `json:"image"`
}

// OracleVersionDatabase is the Oracle Database image
type OracleVersionDatabase struct {
	Image  string `json:"image"`
	BaseOS string `json:"baseOS,omitempty"`
}

// OracleVersionCoordinator is the Oracle leader elector image
type OracleVersionCoordinator struct {
	Image string `json:"image"`
}

// OracleVersionExporter is the image for the Oracle exporter
type OracleVersionExporter struct {
	Image string `json:"image"`
}

// OracleVersionPodSecurityPolicy is the Oracle pod security policies
type OracleVersionPodSecurityPolicy struct {
	DatabasePolicyName string `json:"databasePolicyName"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// OracleVersionList is a list of OracleVersions
type OracleVersionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items is a list of OracleVersion CRD objects
	Items []OracleVersion `json:"items,omitempty"`
}

// OracleSecurityContext is the additional features for the Oracle
type OracleSecurityContext struct {
	// RunAsUser is default UID for the DB container. It is by default 999 for debian based image and 70 for alpine based image.
	// oracle UID 999 for debian images https://github.com/docker-library/oracle/blob/14f13e4b399ed1848fa24c2c1f5bd40c25732bdd/13/Dockerfile#L15
	// oracle UID 70  for alpine images https://github.com/docker-library/oracle/blob/14f13e4b399ed1848fa24c2c1f5bd40c25732bdd/13/alpine/Dockerfile#L6
	RunAsUser *int64 `json:"runAsUser,omitempty"`
}

// +kubebuilder:validation:Enum=Official
type OracleDistro string

const (
	OracleDistroOfficial OracleDistro = "Official"
)
