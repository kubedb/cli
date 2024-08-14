/*
Copyright 2023.

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
	ResourceCodeRabbitMQVersion     = "rmversion"
	ResourceKindRabbitMQVersion     = "RabbitMQVersion"
	ResourceSingularRabbitMQVersion = "rabbitmqversion"
	ResourcePluralRabbitMQVersion   = "rabbitmqversions"
)

// RabbitMQVersion defines a RabbitMQ database version.

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=updateStatus
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=rabbitmqversions,singular=rabbitmqversion,scope=Cluster,shortName=rmversion,categories={catalog,kubedb,appscode}
// +kubebuilder:printcolumn:name="Version",type="string",JSONPath=".spec.version"
// +kubebuilder:printcolumn:name="DB_IMAGE",type="string",JSONPath=".spec.db.image"
// +kubebuilder:printcolumn:name="Deprecated",type="boolean",JSONPath=".spec.deprecated"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type RabbitMQVersion struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              RabbitMQVersionSpec `json:"spec,omitempty"`
}

// RabbitMQVersionSpec is the spec for RabbitMQ version
type RabbitMQVersionSpec struct {
	// Version
	Version string `json:"version"`
	// Database Image
	DB RabbitMQVersionDatabase `json:"db"`
	// Database Image
	InitContainer RabbitMQInitContainer `json:"initContainer"`
	// Deprecated versions usable but regarded as obsolete and best avoided, typically due to having been superseded.
	// +optional
	Deprecated bool `json:"deprecated,omitempty"`
	// SecurityContext is for the additional config for the DB container
	// +optional
	SecurityContext SecurityContext `json:"securityContext"`
	// +optional
	UI []ChartInfo `json:"ui,omitempty"`
}

// RabbitMQVersionDatabase is the RabbitMQ Database image
type RabbitMQVersionDatabase struct {
	Image string `json:"image"`
}

// RabbitMQInitContainer is the RabbitMQ init Container image
type RabbitMQInitContainer struct {
	Image string `json:"image"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// RabbitMQVersionList is a list of RabbitmqVersions
type RabbitMQVersionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items is a list of RedisVersion CRD objects
	Items []RabbitMQVersion `json:"items,omitempty"`
}
