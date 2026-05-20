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
	"kubedb.dev/apimachinery/apis/autoscaling"
	"kubedb.dev/apimachinery/crds"

	"kmodules.xyz/client-go/apiextensions"
)

func (r Neo4jAutoscaler) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralNeo4jAutoscaler))
}

var _ apis.ResourceInfo = &Neo4jAutoscaler{}

func (r Neo4jAutoscaler) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", ResourcePluralNeo4jAutoscaler, autoscaling.GroupName)
}

func (r Neo4jAutoscaler) ResourceShortCode() string {
	return ResourceCodeNeo4jAutoscaler
}

func (r Neo4jAutoscaler) ResourceKind() string {
	return ResourceKindNeo4jAutoscaler
}

func (r Neo4jAutoscaler) ResourceSingular() string {
	return ResourceSingularNeo4jAutoscaler
}

func (r Neo4jAutoscaler) ResourcePlural() string {
	return ResourcePluralNeo4jAutoscaler
}

func (r Neo4jAutoscaler) ValidateSpecs() error {
	return nil
}

var _ StatusAccessor = &Neo4jAutoscaler{}

func (r *Neo4jAutoscaler) GetStatus() AutoscalerStatus {
	return r.Status
}

func (r *Neo4jAutoscaler) SetStatus(s AutoscalerStatus) {
	r.Status = s
}
