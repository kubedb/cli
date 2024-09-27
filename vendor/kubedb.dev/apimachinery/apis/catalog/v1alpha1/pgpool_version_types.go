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
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ResourceCodePgpoolVersion     = "ppversion"
	ResourceKindPgpoolVersion     = "PgpoolVersion"
	ResourceSingularPgpoolVersion = "pgpoolversion"
	ResourcePluralPgpoolVersion   = "pgpoolversions"
)

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=updateStatus
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=pgpoolversions,singular=pgpoolversion,scope=Cluster,shortName=ppversion,categories={catalog,kubedb,appscode}
// +kubebuilder:printcolumn:name="Version",type="string",JSONPath=".spec.version"
// +kubebuilder:printcolumn:name="PGPOOL_IMAGE",type="string",JSONPath=".spec.pgpool.image"
// +kubebuilder:printcolumn:name="Deprecated",type="boolean",JSONPath=".spec.deprecated"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type PgpoolVersion struct {
	meta.TypeMeta   `json:",inline"`
	meta.ObjectMeta `json:"metadata,omitempty"`
	Spec            PgpoolVersionSpec `json:"spec,omitempty"`
}

// PgpoolVersionSpec defines the desired state of PgpoolVersion
type PgpoolVersionSpec struct {
	// Version
	Version string `json:"version"`

	// Pgpool Image
	Pgpool PgpoolVersionPgpool `json:"pgpool"`

	// +optional
	Deprecated bool `json:"deprecated,omitempty"`

	// Exporter Image
	Exporter PgpoolVersionExporter `json:"exporter,omitempty"`

	// update constraints
	UpdateConstraints UpdateConstraints `json:"updateConstraints,omitempty"`

	// SecurityContext is for the additional config for pgpool DB container
	// +optional
	SecurityContext PgpoolSecurityContext `json:"securityContext"`

	// +optional
	UI []ChartInfo `json:"ui,omitempty"`
}

// PgpoolVersionPodSecurityPolicy is the Pgpool pod security policies
type PgpoolVersionPodSecurityPolicy struct {
	DatabasePolicyName string `json:"databasePolicyName"`
}

// PgpoolVersionExporter is the image for the Pgpool exporter
type PgpoolVersionExporter struct {
	Image string `json:"image"`
}

// PgpoolVersionDatabase is the Pgpool Database image
type PgpoolVersionPgpool struct {
	Image string `json:"image"`
}

// PgpoolSecurityContext is the additional features for the Pgpool
type PgpoolSecurityContext struct {
	// RunAsUser is default UID for the DB container. It is by default 70 for postgres user.
	RunAsUser *int64 `json:"runAsUser,omitempty"`

	// RunAsAnyNonRoot will be true if user can change the default db container user to other than postgres user.
	RunAsAnyNonRoot bool `json:"runAsAnyNonRoot,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PgpoolVersionList contains a list of PgpoolVersion
type PgpoolVersionList struct {
	meta.TypeMeta `json:",inline"`
	meta.ListMeta `json:"metadata,omitempty"`
	Items         []PgpoolVersion `json:"items"`
}
