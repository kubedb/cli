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

func (q QdrantAutoscaler) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralQdrantAutoscaler))
}

var _ apis.ResourceInfo = &QdrantAutoscaler{}

func (q QdrantAutoscaler) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", ResourcePluralQdrantAutoscaler, autoscaling.GroupName)
}

func (q QdrantAutoscaler) ResourceShortCode() string {
	return ResourceCodeQdrantAutoscaler
}

func (q QdrantAutoscaler) ResourceKind() string {
	return ResourceKindQdrantAutoscaler
}

func (q QdrantAutoscaler) ResourceSingular() string {
	return ResourceSingularQdrantAutoscaler
}

func (q QdrantAutoscaler) ResourcePlural() string {
	return ResourcePluralQdrantAutoscaler
}

func (q QdrantAutoscaler) ValidateSpecs() error {
	return nil
}

var _ StatusAccessor = &QdrantAutoscaler{}

func (q *QdrantAutoscaler) GetStatus() AutoscalerStatus {
	return q.Status
}

func (q *QdrantAutoscaler) SetStatus(s AutoscalerStatus) {
	q.Status = s
}
