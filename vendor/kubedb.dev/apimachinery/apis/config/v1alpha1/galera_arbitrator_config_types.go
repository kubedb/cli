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
	// Resource Kind for GaleraArbitratorConfiguration
	ResourceKindGaleraArbitratorConfiguration = "GaleraArbitratorConfiguration"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// GaleraArbitratorConfiguration defines Galera ARBitrator Daemon (garbd) configuration.
// Ref: https://galeracluster.com/library/documentation/arbitrator.html
// 		https://galeracluster.com/library/documentation/backup-cluster.html
type GaleraArbitratorConfiguration struct {
	metav1.TypeMeta `json:",inline,omitempty"`

	// Address denotes the logical name of the galera cluster. It is
	// used as the value of the variable named "wsrep_cluster_name"
	// in the replication configuration for galera
	// Ref: https://galeracluster.com/library/documentation/mysql-wsrep-options.html#wsrep-cluster-name
	Address string `json:"address,omitempty" protobuf:"bytes,1,opt,name=address"`

	// Group denotes the collection of cluster members by IP address
	// or resolvable domain name. This address is used as the value of the
	// variable named "wsrep_cluster_address" in the replication configuration
	// for galera. It must be in galera format.
	// Ref: https://galeracluster.com/library/documentation/mysql-wsrep-options.html#wsrep-cluster-address
	Group string `json:"group,omitempty" protobuf:"bytes,2,opt,name=group"`

	// SSTMethod denotes the method or script the node uses during a State Snapshot Transfer.
	// This method is needed to form the SST request string that contains SST request to
	// trigger state snapshot dump (state backup) on one of the other nodes.
	// Ref: https://galeracluster.com/library/documentation/mysql-wsrep-options.html#wsrep-sst-method
	SSTMethod string `json:"sstMethod,omitempty" protobuf:"bytes,3,opt,name=sstMethod"`

	// Stash defines backup and restore task definitions.
	// +optional
	Stash appcat.StashAddonSpec `json:"stash,omitempty" protobuf:"bytes,4,opt,name=stash"`
}
