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

//go:generate go-enum --mustparse --names --values
package v1alpha1

import (
	core "k8s.io/api/core/v1"
	resource "k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ResourceCodeRedisOpsRequest     = "rdops"
	ResourceKindRedisOpsRequest     = "RedisOpsRequest"
	ResourceSingularRedisOpsRequest = "redisopsrequest"
	ResourcePluralRedisOpsRequest   = "redisopsrequests"
)

// RedisOpsRequest defines a Redis DBA operation.

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=redisopsrequests,singular=redisopsrequest,shortName=rdops,categories={ops,kubedb,appscode}
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Type",type="string",JSONPath=".spec.type"
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type RedisOpsRequest struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              RedisOpsRequestSpec `json:"spec,omitempty"`
	Status            OpsRequestStatus    `json:"status,omitempty"`
}

// RedisOpsRequestSpec is the spec for RedisOpsRequest
type RedisOpsRequestSpec struct {
	// Specifies the Redis reference
	DatabaseRef core.LocalObjectReference `json:"databaseRef"`
	// Specifies the ops request type: Upgrade, HorizontalScaling, VerticalScaling etc.
	Type RedisOpsRequestType `json:"type"`
	// Specifies information necessary for upgrading Redis
	UpdateVersion *RedisUpdateVersionSpec `json:"updateVersion,omitempty"`
	// Specifies information necessary for horizontal scaling
	HorizontalScaling *RedisHorizontalScalingSpec `json:"horizontalScaling,omitempty"`
	// Specifies information necessary for vertical scaling
	VerticalScaling *RedisVerticalScalingSpec `json:"verticalScaling,omitempty"`
	// Specifies information necessary for volume expansion
	VolumeExpansion *RedisVolumeExpansionSpec `json:"volumeExpansion,omitempty"`
	// Specifies information necessary for custom configuration of Redis
	Configuration *RedisCustomConfigurationSpec `json:"configuration,omitempty"`
	// Specifies information necessary for configuring TLS
	TLS *RedisTLSSpec `json:"tls,omitempty"`
	// Specifies information necessary for configuring authSecret of the database
	Authentication *AuthSpec `json:"authentication,omitempty"`
	// Specifies information necessary for restarting database
	Restart *RestartSpec `json:"restart,omitempty"`
	// Announce is used to announce the redis cluster endpoints.
	// It is used to set
	// cluster-announce-ip, cluster-announce-port, cluster-announce-bus-port, cluster-announce-tls-port
	Announce *Announce `json:"announce,omitempty"`
	// Specifies information necessary for replacing sentinel instances
	Sentinel *RedisSentinelSpec `json:"sentinel,omitempty"`
	// Timeout for each step of the ops request in second. If a step doesn't finish within the specified timeout, the ops request will result in failure.
	Timeout *metav1.Duration `json:"timeout,omitempty"`
	// ApplyOption is to control the execution of OpsRequest depending on the database state.
	// +kubebuilder:default="IfReady"
	Apply ApplyOption `json:"apply,omitempty"`
}

// +kubebuilder:validation:Enum=UpdateVersion;HorizontalScaling;VerticalScaling;VolumeExpansion;Restart;Reconfigure;ReconfigureTLS;ReplaceSentinel;RotateAuth;Announce
// ENUM(UpdateVersion, HorizontalScaling, VerticalScaling, VolumeExpansion, Restart, Reconfigure, ReconfigureTLS, ReplaceSentinel, RotateAuth, Announce)
type RedisOpsRequestType string

type RedisTLSSpec struct {
	*TLSSpec `json:",inline"`
	// This field is only needed in Redis Sentinel Mode when we add or remove TLS. In Redis Sentinel Mode, both redis instances and
	// sentinel instances either have TLS or don't have TLS. So when want to add TLS to Redis in Sentinel Mode, current sentinel instances don't
	// have TLS enabled, so we need to give a new Sentinel Reference which has TLS enabled and which will monitor the Redis instances when we
	// add TLS to it
	// +optional
	Sentinel *RedisSentinelSpec `json:"sentinel,omitempty"`
}

type RedisSentinelSpec struct {
	// Sentinel Ref for new Sentinel which will replace the old sentinel
	Ref *RedisSentinelRef `json:"ref"`
	// +optional
	RemoveUnusedSentinel bool `json:"removeUnusedSentinel,omitempty"`
}

type RedisSentinelRef struct {
	// Name of the refereed sentinel
	Name string `json:"name,omitempty"`

	// Namespace where refereed sentinel has been deployed
	Namespace string `json:"namespace,omitempty"`
}

// RedisReplicaReadinessCriteria is the criteria for checking readiness of a Redis pod
// after updating, horizontal scaling etc.
type RedisReplicaReadinessCriteria struct{}

type RedisUpdateVersionSpec struct {
	// Specifies the target version name from catalog
	TargetVersion     string                         `json:"targetVersion,omitempty"`
	ReadinessCriteria *RedisReplicaReadinessCriteria `json:"readinessCriteria,omitempty"`
}

type RedisHorizontalScalingSpec struct {
	// Number of shards in the cluster
	Shards *int32 `json:"shards,omitempty"`
	// specifies the number of replica of the shards
	Replicas *int32 `json:"replicas,omitempty"`

	// Announce is used to announce the redis cluster endpoints.
	// It is used to set
	// cluster-announce-ip, cluster-announce-port, cluster-announce-bus-port, cluster-announce-tls-port
	// While scaling up shard or replica just provide the missing announces.
	Announce *Announce `json:"announce,omitempty"`
}

// RedisVerticalScalingSpec is the spec for Redis vertical scaling
type RedisVerticalScalingSpec struct {
	Redis       *PodResources       `json:"redis,omitempty"`
	Exporter    *ContainerResources `json:"exporter,omitempty"`
	Coordinator *ContainerResources `json:"coordinator,omitempty"`
}

type RedisAclSpec struct {
	// SecretRef holds the password against which ACLs will be created if syncACL is given.
	// +optional
	SecretRef *core.LocalObjectReference `json:"secretRef,omitempty"`

	// SyncACL specifies the list of users whose ACLs should be synchronized with the new authentication secret.
	// If provided, the system will update the ACLs for these users to ensure they are in sync with the new authentication settings.
	SyncACL []string `json:"syncACL,omitempty"`

	// DeleteUsers specifies the list of users that should be deleted from the database.
	// If provided, the system will remove these users from the database to enhance security or manage
	DeleteUsers []string `json:"deleteUsers,omitempty"`
}

// RedisVolumeExpansionSpec is the spec for Redis volume expansion
type RedisVolumeExpansionSpec struct {
	Mode  VolumeExpansionMode `json:"mode"`
	Redis *resource.Quantity  `json:"redis,omitempty"`
}

type RedisCustomConfigurationSpec struct {
	ConfigSecret       *core.LocalObjectReference `json:"configSecret,omitempty"`
	ApplyConfig        map[string]string          `json:"applyConfig,omitempty"`
	RemoveCustomConfig bool                       `json:"removeCustomConfig,omitempty"`
	Auth               *RedisAclSpec              `json:"auth,omitempty"`
}

// +kubebuilder:validation:Enum=ip;hostname
type PreferredEndpointType string

const (
	PreferredEndpointTypeIP       PreferredEndpointType = "ip"
	PreferredEndpointTypeHostname PreferredEndpointType = "hostname"
)

type Announce struct {
	// +kubebuilder:default=hostname
	Type PreferredEndpointType `json:"type,omitempty"`
	// This field is used to set cluster-announce information for redis cluster of each shard.
	Shards []Shards `json:"shards,omitempty"`
}

type Shards struct {
	// Endpoints contains the cluster-announce information for all the replicas in a shard.
	// This will be used to set cluster-announce-ip/hostname, cluster-announce-port/cluster-announce-tls-port
	// and cluster-announce-bus-port
	// format cluster-announce (host:port@busport)
	Endpoints []string `json:"endpoints,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// RedisOpsRequestList is a list of RedisOpsRequests
type RedisOpsRequestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items is a list of RedisOpsRequest CRD objects
	Items []RedisOpsRequest `json:"items,omitempty"`
}
