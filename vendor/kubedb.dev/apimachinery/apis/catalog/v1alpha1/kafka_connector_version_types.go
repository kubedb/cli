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
	ResourceCodeKafkaConnectorVersion     = "kcversion"
	ResourceKindKafkaConnectorVersion     = "KafkaConnectorVersion"
	ResourceSingularKafkaConnectorVersion = "kafkaconnectorversion"
	ResourcePluralKafkaConnectorVersion   = "kafkaconnectorversions"
)

// KafkaConnectorVersion defines a Kafka connector plugins version.

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=updateStatus
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=kafkaconnectorversions,singular=kafkaconnectorversion,scope=Cluster,shortName=kcversion,categories={catalog,kubedb,appscode}
// +kubebuilder:printcolumn:name="Version",type="string",JSONPath=".spec.version"
// +kubebuilder:printcolumn:name="Connector_Image",type="string",JSONPath=".spec.connectorPlugin.image"
// +kubebuilder:printcolumn:name="Deprecated",type="boolean",JSONPath=".spec.deprecated"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type KafkaConnectorVersion struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              KafkaConnectorVersionSpec `json:"spec,omitempty"`
}

// KafkaConnectorVersionSpec is the spec for kafka connector plugin version
type KafkaConnectorVersionSpec struct {
	// Type of the connector plugins(ex. mongodb, s3, etc.)
	Type string `json:"type"`
	// Version of the connector plugins
	Version string `json:"version"`
	// Database Image
	ConnectorPlugin ConnectorPlugin `json:"connectorPlugin"`
	// Deprecated versions usable but regarded as obsolete and best avoided, typically due to having been superseded.
	// +optional
	Deprecated bool `json:"deprecated,omitempty"`
	// SecurityContext is for the additional config for the init container
	// +optional
	SecurityContext SecurityContext `json:"securityContext"`
}

// ConnectorPlugin is the Kafka connector plugin image
type ConnectorPlugin struct {
	Image string `json:"image"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KafkaConnectorVersionList is a list of KafkaConnectorVersion
type KafkaConnectorVersionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items is a list of KafkaConnectorVersion CRD objects
	Items []KafkaConnectorVersion `json:"items,omitempty"`
}
