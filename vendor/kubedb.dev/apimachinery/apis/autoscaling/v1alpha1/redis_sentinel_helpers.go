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

func (_ RedisSentinelAutoscaler) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralRedisSentinelAutoscaler))
}

var _ apis.ResourceInfo = &RedisSentinelAutoscaler{}

func (r RedisSentinelAutoscaler) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", ResourcePluralRedisSentinelAutoscaler, autoscaling.GroupName)
}

func (r RedisSentinelAutoscaler) ResourceShortCode() string {
	return ResourceCodeRedisSentinelAutoscaler
}

func (r RedisSentinelAutoscaler) ResourceKind() string {
	return ResourceKindRedisSentinelAutoscaler
}

func (r RedisSentinelAutoscaler) ResourceSingular() string {
	return ResourceSingularRedisSentinelAutoscaler
}

func (r RedisSentinelAutoscaler) ResourcePlural() string {
	return ResourcePluralRedisSentinelAutoscaler
}

func (r RedisSentinelAutoscaler) ValidateSpecs() error {
	return nil
}

var _ StatusAccessor = &RedisSentinelAutoscaler{}

func (e *RedisSentinelAutoscaler) GetStatus() AutoscalerStatus {
	return e.Status
}

func (e *RedisSentinelAutoscaler) SetStatus(s AutoscalerStatus) {
	e.Status = s
}
