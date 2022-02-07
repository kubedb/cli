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
	ResourceKindProxySQLVersion     = "ProxySQLVersion"
	ResourceSingularProxySQLVersion = "proxysqlversion"
	ResourcePluralProxySQLVersion   = "proxysqlversions"
)

// ProxySQLVersion defines a ProxySQL load-balancer version.

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=updateStatus
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=proxysqlversions,singular=proxysqlversion,scope=Cluster,categories={datastore,kubedb,appscode}
// +kubebuilder:printcolumn:name="Version",type="string",JSONPath=".spec.version"
// +kubebuilder:printcolumn:name="ProxySQL_IMAGE",type="string",JSONPath=".spec.db.image"
// +kubebuilder:printcolumn:name="Deprecated",type="boolean",JSONPath=".spec.deprecated"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type ProxySQLVersion struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              ProxySQLVersionSpec `json:"spec,omitempty"`
}

// ProxySQLVersionSpec is the spec for ProxySQL version
type ProxySQLVersionSpec struct {
	// Version
	Version string `json:"version"`
	// Proxysql Image
	Proxysql ProxySQLVersionProxysql `json:"proxysql"`
	// Exporter Image
	Exporter ProxySQLVersionExporter `json:"exporter"`
	// Deprecated versions usable but regarded as obsolete and best avoided, typically due to having been superseded.
	// +optional
	Deprecated bool `json:"deprecated,omitempty"`
	// PSP names
	PodSecurityPolicies ProxySQLVersionPodSecurityPolicy `json:"podSecurityPolicies"`
}

// ProxySQLVersionProxysql is the proxysql image
type ProxySQLVersionProxysql struct {
	Image string `json:"image"`
}

// ProxySQLVersionExporter is the image for the ProxySQL exporter
type ProxySQLVersionExporter struct {
	Image string `json:"image"`
}

// ProxySQLVersionPodSecurityPolicy is the ProxySQL pod security policies
type ProxySQLVersionPodSecurityPolicy struct {
	DatabasePolicyName string `json:"databasePolicyName"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ProxySQLVersionList is a list of ProxySQLVersions
type ProxySQLVersionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items is a list of ProxySQLVersion CRD objects
	Items []ProxySQLVersion `json:"items,omitempty"`
}
