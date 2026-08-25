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

func (r EtcdAutoscaler) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralEtcdAutoscaler))
}

var _ apis.ResourceInfo = &EtcdAutoscaler{}

func (r EtcdAutoscaler) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", ResourcePluralEtcdAutoscaler, autoscaling.GroupName)
}

func (r EtcdAutoscaler) ResourceShortCode() string {
	return ResourceCodeEtcdAutoscaler
}

func (r EtcdAutoscaler) ResourceKind() string {
	return ResourceKindEtcdAutoscaler
}

func (r EtcdAutoscaler) ResourceSingular() string {
	return ResourceSingularEtcdAutoscaler
}

func (r EtcdAutoscaler) ResourcePlural() string {
	return ResourcePluralEtcdAutoscaler
}

func (r EtcdAutoscaler) ValidateSpecs() error {
	return nil
}

var _ StatusAccessor = &EtcdAutoscaler{}

func (r *EtcdAutoscaler) GetStatus() AutoscalerStatus {
	return r.Status
}

func (r *EtcdAutoscaler) SetStatus(s AutoscalerStatus) {
	r.Status = s
}
