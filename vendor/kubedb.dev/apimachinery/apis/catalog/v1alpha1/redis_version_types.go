/*
Copyright The KubeDB Authors.

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
	ResourceCodeRedisVersion     = "rdversion"
	ResourceKindRedisVersion     = "RedisVersion"
	ResourceSingularRedisVersion = "redisversion"
	ResourcePluralRedisVersion   = "redisversions"
)

// RedisVersion defines a Redis database version.

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=updateStatus
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=redisversions,singular=redisversion,scope=Cluster,shortName=rdversion,categories={datastore,kubedb,appscode}
// +kubebuilder:printcolumn:name="Version",type="string",JSONPath=".spec.version"
// +kubebuilder:printcolumn:name="DB_IMAGE",type="string",JSONPath=".spec.db.image"
// +kubebuilder:printcolumn:name="Deprecated",type="boolean",JSONPath=".spec.deprecated"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type RedisVersion struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Spec              RedisVersionSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
}

// RedisVersionSpec is the spec for redis version
type RedisVersionSpec struct {
	// Version
	Version string `json:"version" protobuf:"bytes,1,opt,name=version"`
	// Database Image
	DB RedisVersionDatabase `json:"db" protobuf:"bytes,2,opt,name=db"`
	// Exporter Image
	Exporter RedisVersionExporter `json:"exporter" protobuf:"bytes,3,opt,name=exporter"`
	// Deprecated versions usable but regarded as obsolete and best avoided, typically due to having been superseded.
	// +optional
	Deprecated bool `json:"deprecated,omitempty" protobuf:"varint,4,opt,name=deprecated"`
	// PSP names
	PodSecurityPolicies RedisVersionPodSecurityPolicy `json:"podSecurityPolicies" protobuf:"bytes,5,opt,name=podSecurityPolicies"`
}

// RedisVersionDatabase is the Redis Database image
type RedisVersionDatabase struct {
	Image string `json:"image" protobuf:"bytes,1,opt,name=image"`
}

// RedisVersionExporter is the image for the Redis exporter
type RedisVersionExporter struct {
	Image string `json:"image" protobuf:"bytes,1,opt,name=image"`
}

// RedisVersionPodSecurityPolicy is the Redis pod security policies
type RedisVersionPodSecurityPolicy struct {
	DatabasePolicyName string `json:"databasePolicyName" protobuf:"bytes,1,opt,name=databasePolicyName"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// RedisVersionList is a list of RedisVersions
type RedisVersionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	// Items is a list of RedisVersion CRD objects
	Items []RedisVersion `json:"items,omitempty" protobuf:"bytes,2,rep,name=items"`
}
