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

func (_ RedisAutoscaler) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralRedisAutoscaler))
}

var _ apis.ResourceInfo = &RedisAutoscaler{}

func (r RedisAutoscaler) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", ResourcePluralRedisAutoscaler, autoscaling.GroupName)
}

func (r RedisAutoscaler) ResourceShortCode() string {
	return ResourceCodeRedisAutoscaler
}

func (r RedisAutoscaler) ResourceKind() string {
	return ResourceKindRedisAutoscaler
}

func (r RedisAutoscaler) ResourceSingular() string {
	return ResourceSingularRedisAutoscaler
}

func (r RedisAutoscaler) ResourcePlural() string {
	return ResourcePluralRedisAutoscaler
}

func (r RedisAutoscaler) ValidateSpecs() error {
	return nil
}

var _ StatusAccessor = &RedisAutoscaler{}

func (e *RedisAutoscaler) GetStatus() AutoscalerStatus {
	return e.Status
}

func (e *RedisAutoscaler) SetStatus(s AutoscalerStatus) {
	e.Status = s
}
