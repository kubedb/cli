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
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	appcat "kmodules.xyz/custom-resources/apis/appcatalog/v1alpha1"
)

const (
	ResourceKindRedisConfiguration = "RedisConfiguration"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// RedisConfiguration defines a Redis appBinding configuration.
type RedisConfiguration struct {
	metav1.TypeMeta `json:",inline,omitempty"`

	// ClientCertSecret is the client secret name which needs to provide when the Redis client need to authenticate with client key and cert.
	// It will be used when `tls-auth-clients` value is set to `required` or `yes`.
	// +optional
	ClientCertSecret *v1.LocalObjectReference `json:"clientCertSecret,omitempty" protobuf:"bytes,1,opt,name=clientCertSecret"`

	// Stash defines backup and restore task definitions.
	// +optional
	Stash appcat.StashAddonSpec `json:"stash,omitempty" protobuf:"bytes,2,opt,name=stash"`
}
