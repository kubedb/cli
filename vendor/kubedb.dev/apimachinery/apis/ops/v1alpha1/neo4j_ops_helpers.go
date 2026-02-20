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
	"fmt"

	"kubedb.dev/apimachinery/apis"
	"kubedb.dev/apimachinery/apis/ops"
	"kubedb.dev/apimachinery/crds"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kmodules.xyz/client-go/apiextensions"
)

func (n *Neo4jOpsRequest) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralNeo4jOpsRequest))
}

var _ apis.ResourceInfo = &Neo4jOpsRequest{}

func (n *Neo4jOpsRequest) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", ResourcePluralNeo4jOpsRequest, ops.GroupName)
}

func (n *Neo4jOpsRequest) ResourceShortCode() string {
	return ResourceCodeNeo4jOpsRequest
}

func (n *Neo4jOpsRequest) ResourceKind() string {
	return ResourceKindNeo4jOpsRequest
}

func (n *Neo4jOpsRequest) ResourceSingular() string {
	return ResourceSingularNeo4jOpsRequest
}

func (n *Neo4jOpsRequest) ResourcePlural() string {
	return ResourcePluralNeo4jOpsRequest
}

var _ Accessor = &Neo4jOpsRequest{}

func (n *Neo4jOpsRequest) GetObjectMeta() metav1.ObjectMeta {
	return n.ObjectMeta
}

func (n *Neo4jOpsRequest) GetDBRefName() string {
	return n.Spec.DatabaseRef.Name
}

func (n *Neo4jOpsRequest) GetRequestType() string {
	return string(n.Spec.Type)
}

func (n *Neo4jOpsRequest) GetStatus() OpsRequestStatus {
	return n.Status
}

func (n *Neo4jOpsRequest) SetStatus(s OpsRequestStatus) {
	n.Status = s
}
