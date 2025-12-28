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

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

const (
	ResourceCodeFerretDBVersion     = "frversion"
	ResourceKindFerretDBVersion     = "FerretDBVersion"
	ResourceSingularFerretDBVersion = "ferretdbversion"
	ResourcePluralFerretDBVersion   = "ferretdbversions"
)

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=updateStatus
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=ferretdbversions,singular=ferretdbversion,scope=Cluster,shortName=frversion,categories={catalog,kubedb,appscode}
// +kubebuilder:printcolumn:name="Version",type="string",JSONPath=".spec.version"
// +kubebuilder:printcolumn:name="DB_IMAGE",type="string",JSONPath=".spec.db.image"
// +kubebuilder:printcolumn:name="Deprecated",type="boolean",JSONPath=".spec.deprecated"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type FerretDBVersion struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              FerretDBVersionSpec `json:"spec,omitempty"`
}

// FerretDBVersionSpec defines the desired state of FerretDBVersion
type FerretDBVersionSpec struct {
	// Version
	Version string `json:"version"`

	// EndOfLife refers if this version reached into its end of the life or not, based on https://endoflife.date/
	// +optional
	EndOfLife bool `json:"endOfLife"`

	// Database Image
	DB FerretDBVersionDatabase `json:"db"`

	// Postgres version
	Postgres FerretDBVersionPostgres `json:"postgres"`

	// Deprecated versions usable but regarded as obsolete and best avoided, typically due to having been superseded.
	// +optional
	Deprecated bool `json:"deprecated,omitempty"`

	// update constraints
	// +optional
	UpdateConstraints UpdateConstraints `json:"updateConstraints,omitempty"`

	// SecurityContext is for the additional security information for the FerretDB container
	// +optional
	SecurityContext SecurityContext `json:"securityContext"`

	// +optional
	UI []ChartInfo `json:"ui,omitempty"`
}

// FerretDBVersionDatabase is the FerretDB Database image
type FerretDBVersionDatabase struct {
	Image string `json:"image"`
}

type FerretDBVersionPostgres struct {
	// Which versions pg will be used as backend of ferretdb. default 13.13 when backend internally managed
	Version string `json:"version"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// FerretDBVersionList contains a list of FerretDBVersion
type FerretDBVersionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []FerretDBVersion `json:"items,omitempty"`
}
