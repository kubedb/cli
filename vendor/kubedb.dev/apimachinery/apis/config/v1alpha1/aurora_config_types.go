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
	// Resource Kind for AuroraConfiguration
	ResourceKindAuroraConfiguration = "AuroraConfiguration"

	// AppTypeAWSAurora marks an AppBinding as representing an external AWS Aurora
	// (MySQL-compatible) cluster rather than a KubeDB-managed engine. Consumers such as the
	// ProxySQL operator set this on the AppBinding's spec.type to detect an Aurora backend.
	AppTypeAWSAurora appcat.AppType = "kubedb.com/aws-aurora"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AuroraConfiguration carries innate, connectivity-level facts about an AWS Aurora cluster that
// are intrinsic to the cluster itself, not a per-consumer routing/tuning policy. Consumers read
// this from the Aurora AppBinding's spec.parameters.
type AuroraConfiguration struct {
	metav1.TypeMeta `json:",inline"`

	// DomainName is the FQDN suffix (starting with a dot, e.g.
	// ".xxxxx.<region>.rds.amazonaws.com") ProxySQL's native Aurora monitor appends to each
	// discovered instance identifier to build its hostname — see mysql_aws_aurora_hostgroups.
	// If empty, consumers derive it from the AppBinding's own clientConfig host (the
	// writer/cluster endpoint) by stripping the leading "<cluster-id>.cluster-" label. Set this
	// explicitly for Aurora Global Database secondary regions or custom endpoints that don't
	// follow that convention.
	// +optional
	DomainName string `json:"domainName,omitempty"`
}
