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

func (_ ElasticsearchAutoscaler) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralElasticsearchAutoscaler))
}

var _ apis.ResourceInfo = &ElasticsearchAutoscaler{}

func (e ElasticsearchAutoscaler) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", ResourcePluralElasticsearchAutoscaler, autoscaling.GroupName)
}

func (e ElasticsearchAutoscaler) ResourceShortCode() string {
	return ResourceCodeElasticsearchAutoscaler
}

func (e ElasticsearchAutoscaler) ResourceKind() string {
	return ResourceKindElasticsearchAutoscaler
}

func (e ElasticsearchAutoscaler) ResourceSingular() string {
	return ResourceSingularElasticsearchAutoscaler
}

func (e ElasticsearchAutoscaler) ResourcePlural() string {
	return ResourcePluralElasticsearchAutoscaler
}

func (e ElasticsearchAutoscaler) ValidateSpecs() error {
	return nil
}

var _ StatusAccessor = &ElasticsearchAutoscaler{}

func (e *ElasticsearchAutoscaler) GetStatus() AutoscalerStatus {
	return e.Status
}

func (e *ElasticsearchAutoscaler) SetStatus(s AutoscalerStatus) {
	e.Status = s
}
