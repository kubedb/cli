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
	ResourceCodeAerospikeVersion     = "arversion"
	ResourceKindAerospikeVersion     = "AerospikeVersion"
	ResourceSingularAerospikeVersion = "aerospikeversion"
	ResourcePluralAerospikeVersion   = "aerospikeversions"
)

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=updateStatus
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=aerospikeversions,singular=aerospikeversion,scope=Cluster,shortName=arversion,categories={catalog,kubedb,appscode}
// +kubebuilder:printcolumn:name="Version",type="string",JSONPath=".spec.version"
// +kubebuilder:printcolumn:name="AEROSPIKE_IMAGE",type="string",JSONPath=".spec.aerospike.image"
// +kubebuilder:printcolumn:name="Deprecated",type="boolean",JSONPath=".spec.deprecated"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type AerospikeVersion struct {
	meta.TypeMeta   `json:",inline"`
	meta.ObjectMeta `json:"metadata,omitempty"`
	Spec            AerospikeVersionSpec `json:"spec,omitempty"`
}

// AerospikeVersionSpec defines the desired state of AerospikeVersion
type AerospikeVersionSpec struct {
	// Version
	Version string `json:"version"`

	// EndOfLife refers if this version reached into its end of the life or not, based on https://endoflife.date/
	// +optional
	EndOfLife bool `json:"endOfLife"`

	// Aerospike Image
	Aerospike AerospikeVersionAerospike `json:"aerospike"`

	// +optional
	Deprecated bool `json:"deprecated,omitempty"`

	// +optional
	GitSyncer GitSyncer `json:"gitSyncer,omitempty"`

	// Exporter Image
	Exporter AerospikeVersionExporter `json:"exporter,omitempty"`

	// update constraints
	UpdateConstraints UpdateConstraints `json:"updateConstraints,omitempty"`

	// SecurityContext is for the additional config for aerospike DB container
	// +optional
	SecurityContext AerospikeSecurityContext `json:"securityContext"`

	// +optional
	UI []ChartInfo `json:"ui,omitempty"`
}

// AerospikeVersionPodSecurityPolicy is the Aerospike pod security policies
type AerospikeVersionPodSecurityPolicy struct {
	DatabasePolicyName string `json:"databasePolicyName"`
}

// AerospikeVersionExporter is the image for the Aerospike exporter
type AerospikeVersionExporter struct {
	Image string `json:"image"`
}

// AerospikeVersionDatabase is the Aerospike Database image
type AerospikeVersionAerospike struct {
	Image string `json:"image"`
}

// AerospikeSecurityContext is the additional features for the Aerospike
type AerospikeSecurityContext struct {
	// RunAsUser is default UID for the DB container. It is by default 70 for postgres user.
	RunAsUser *int64 `json:"runAsUser,omitempty"`

	// RunAsAnyNonRoot will be true if user can change the default db container user to other than postgres user.
	RunAsAnyNonRoot bool `json:"runAsAnyNonRoot,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AerospikeVersionList contains a list of AerospikeVersion
type AerospikeVersionList struct {
	meta.TypeMeta `json:",inline"`
	meta.ListMeta `json:"metadata,omitempty"`
	Items         []AerospikeVersion `json:"items"`
}
