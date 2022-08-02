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

func (_ MariaDBAutoscaler) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralMariaDBAutoscaler))
}

var _ apis.ResourceInfo = &MariaDBAutoscaler{}

func (m MariaDBAutoscaler) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", ResourcePluralMariaDBAutoscaler, autoscaling.GroupName)
}

func (m MariaDBAutoscaler) ResourceShortCode() string {
	return ResourceCodeMariaDBAutoscaler
}

func (m MariaDBAutoscaler) ResourceKind() string {
	return ResourceKindMariaDBAutoscaler
}

func (m MariaDBAutoscaler) ResourceSingular() string {
	return ResourceSingularMariaDBAutoscaler
}

func (m MariaDBAutoscaler) ResourcePlural() string {
	return ResourcePluralMariaDBAutoscaler
}

func (m MariaDBAutoscaler) ValidateSpecs() error {
	return nil
}

var _ StatusAccessor = &MariaDBAutoscaler{}

func (e *MariaDBAutoscaler) GetStatus() AutoscalerStatus {
	return e.Status
}

func (e *MariaDBAutoscaler) SetStatus(s AutoscalerStatus) {
	e.Status = s
}
