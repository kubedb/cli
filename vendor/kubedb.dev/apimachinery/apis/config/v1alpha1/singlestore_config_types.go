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
	appcat "kmodules.xyz/custom-resources/apis/appcatalog/v1alpha1"
)

const (
	ResourceKindSinglestoreConfiguration = "SinglestoreConfiguration"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SinglestoreConfiguration defines a Singlestore appBinding configuration.
type SinglestoreConfiguration struct {
	metav1.TypeMeta `json:",inline,omitempty"`

	// ReplicaSets contains the dns of each replicaset of sharding. The DSNs are in key-value pair, where
	// the keys are host-0, host-1 etc, and the values are DSN of each replicaset. If there is no sharding
	// but only one replicaset, then ReplicaSets field contains only one key-value pair where the key is
	// host-0 and the value is dsn of that replicaset.
	ReplicaSets map[string]string `json:"replicaSets,omitempty"`

	MasterAggregator string `json:"masterAggregator,omitempty"`

	// Stash defines backup and restore task definitions.
	// +optional
	Stash appcat.StashAddonSpec `json:"stash,omitempty"`
}
