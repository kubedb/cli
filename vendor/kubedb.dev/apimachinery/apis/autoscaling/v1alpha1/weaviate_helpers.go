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

func (r WeaviateAutoscaler) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralWeaviateAutoscaler))
}

var _ apis.ResourceInfo = &WeaviateAutoscaler{}

func (r WeaviateAutoscaler) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", ResourcePluralWeaviateAutoscaler, autoscaling.GroupName)
}

func (r WeaviateAutoscaler) ResourceShortCode() string {
	return ResourceCodeWeaviateAutoscaler
}

func (r WeaviateAutoscaler) ResourceKind() string {
	return ResourceKindWeaviateAutoscaler
}

func (r WeaviateAutoscaler) ResourceSingular() string {
	return ResourceSingularWeaviateAutoscaler
}

func (r WeaviateAutoscaler) ResourcePlural() string {
	return ResourcePluralWeaviateAutoscaler
}

func (r WeaviateAutoscaler) ValidateSpecs() error {
	return nil
}

var _ StatusAccessor = &WeaviateAutoscaler{}

func (r *WeaviateAutoscaler) GetStatus() AutoscalerStatus {
	return r.Status
}

func (r *WeaviateAutoscaler) SetStatus(s AutoscalerStatus) {
	r.Status = s
}
