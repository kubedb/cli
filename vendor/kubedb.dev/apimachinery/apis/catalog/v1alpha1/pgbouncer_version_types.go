package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

const (
	ResourceCodePgBouncerVersion     = "pbversion"
	ResourceKindPgBouncerVersion     = "PgBouncerVersion"
	ResourceSingularPgBouncerVersion = "pgbouncerversion"
	ResourcePluralPgBouncerVersion   = "pgbouncerversions"
)

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=updateStatus
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PgBouncerVersion defines a PgBouncer database version.
type PgBouncerVersion struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              PgBouncerVersionSpec `json:"spec,omitempty"`
}

// PgBouncerVersionSpec is the spec for pgbouncer version
type PgBouncerVersionSpec struct {
	// Version
	Version string `json:"version"`
	// Database Image
	Server PgBouncerServerVersion `json:"server"`
	// Exporter Image
	Exporter PgBouncerVersionExporter `json:"exporter"`
	// Deprecated versions usable but regarded as obsolete and best avoided, typically due to having been superseded.
	// +optional
	Deprecated bool `json:"deprecated,omitempty"`
}

// PgBouncerServerVersion is the PgBouncer Database image
type PgBouncerServerVersion struct {
	Image string `json:"image"`
}

// PostgresVersionExporter is the image for the Postgres exporter
type PgBouncerVersionExporter struct {
	Image string `json:"image"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PgBouncerVersionList is a list of PgBouncerVersions
type PgBouncerVersionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items is a list of PgBouncerVersion CRD objects
	Items []PgBouncerVersion `json:"items,omitempty"`
}
