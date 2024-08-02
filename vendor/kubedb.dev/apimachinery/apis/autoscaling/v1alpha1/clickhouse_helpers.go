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

func (r ClickHouseAutoscaler) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralClickHouseAutoscaler))
}

var _ apis.ResourceInfo = &ClickHouseAutoscaler{}

func (r ClickHouseAutoscaler) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", ResourcePluralClickHouseAutoscaler, autoscaling.GroupName)
}

func (r ClickHouseAutoscaler) ResourceShortCode() string {
	return ResourceCodeClickHouseAutoscaler
}

func (r ClickHouseAutoscaler) ResourceKind() string {
	return ResourceKindClickHouseAutoscaler
}

func (r ClickHouseAutoscaler) ResourceSingular() string {
	return ResourceSingularClickHouseAutoscaler
}

func (r ClickHouseAutoscaler) ResourcePlural() string {
	return ResourcePluralClickHouseAutoscaler
}

func (r ClickHouseAutoscaler) ValidateSpecs() error {
	return nil
}

var _ StatusAccessor = &ClickHouseAutoscaler{}

func (r *ClickHouseAutoscaler) GetStatus() AutoscalerStatus {
	return r.Status
}

func (r *ClickHouseAutoscaler) SetStatus(s AutoscalerStatus) {
	r.Status = s
}
