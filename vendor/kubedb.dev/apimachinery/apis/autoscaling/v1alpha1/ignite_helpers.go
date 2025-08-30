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

func (r IgniteAutoscaler) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralIgniteAutoscaler))
}

var _ apis.ResourceInfo = &IgniteAutoscaler{}

func (p IgniteAutoscaler) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", ResourcePluralIgniteAutoscaler, autoscaling.GroupName)
}

func (p IgniteAutoscaler) ResourceShortCode() string {
	return ResourceCodeIgniteAutoscaler
}

func (p IgniteAutoscaler) ResourceKind() string {
	return ResourceKindIgniteAutoscaler
}

func (p IgniteAutoscaler) ResourceSingular() string {
	return ResourceSingularIgniteAutoscaler
}

func (p IgniteAutoscaler) ResourcePlural() string {
	return ResourcePluralIgniteAutoscaler
}

func (p IgniteAutoscaler) ValidateSpecs() error {
	return nil
}

var _ StatusAccessor = &IgniteAutoscaler{}

func (e *IgniteAutoscaler) GetStatus() AutoscalerStatus {
	return e.Status
}

func (e *IgniteAutoscaler) SetStatus(s AutoscalerStatus) {
	e.Status = s
}
