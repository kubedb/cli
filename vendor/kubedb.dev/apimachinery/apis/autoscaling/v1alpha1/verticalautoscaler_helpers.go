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
	"kubedb.dev/apimachinery/apis"
	"kubedb.dev/apimachinery/crds"

	"kmodules.xyz/client-go/apiextensions"
)

func (_ VerticalAutoscaler) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralVerticalAutoscaler))
}

var _ apis.ResourceInfo = &VerticalAutoscaler{}

func (r VerticalAutoscaler) ResourceShortCode() string {
	return ResourceCodeVerticalAutoscaler
}

func (r VerticalAutoscaler) ResourceKind() string {
	return ResourceKindVerticalAutoscaler
}

func (r VerticalAutoscaler) ResourceSingular() string {
	return ResourceSingularVerticalAutoscaler
}

func (r VerticalAutoscaler) ResourcePlural() string {
	return ResourcePluralVerticalAutoscaler
}

func (r VerticalAutoscaler) ValidateSpecs() error {
	return nil
}

func (_ VerticalAutoscalerCheckpoint) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource("verticalautoscalercheckpoints"))
}
