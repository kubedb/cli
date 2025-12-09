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

func (*DruidAutoscaler) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralDruidAutoscaler))
}

var _ apis.ResourceInfo = &DruidAutoscaler{}

func (d *DruidAutoscaler) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", ResourcePluralDruidAutoscaler, autoscaling.GroupName)
}

func (d *DruidAutoscaler) ResourceShortCode() string {
	return ResourceCodeDruidAutoscaler
}

func (d *DruidAutoscaler) ResourceKind() string {
	return ResourceKindDruidAutoscaler
}

func (d *DruidAutoscaler) ResourceSingular() string {
	return ResourceSingularDruidAutoscaler
}

func (d *DruidAutoscaler) ResourcePlural() string {
	return ResourcePluralDruidAutoscaler
}

func (d *DruidAutoscaler) ValidateSpecs() error {
	return nil
}

var _ StatusAccessor = &DruidAutoscaler{}

func (d *DruidAutoscaler) GetStatus() AutoscalerStatus {
	return d.Status
}

func (d *DruidAutoscaler) SetStatus(s AutoscalerStatus) {
	d.Status = s
}
