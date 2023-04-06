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
	"kubedb.dev/apimachinery/apis/ops"
	"kubedb.dev/apimachinery/crds"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kmodules.xyz/client-go/apiextensions"
)

func (r RedisOpsRequest) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralRedisOpsRequest))
}

var _ apis.ResourceInfo = &RedisOpsRequest{}

func (r RedisOpsRequest) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", ResourcePluralRedisOpsRequest, ops.GroupName)
}

func (r RedisOpsRequest) ResourceShortCode() string {
	return ResourceCodeRedisOpsRequest
}

func (r RedisOpsRequest) ResourceKind() string {
	return ResourceKindRedisOpsRequest
}

func (r RedisOpsRequest) ResourceSingular() string {
	return ResourceSingularRedisOpsRequest
}

func (r RedisOpsRequest) ResourcePlural() string {
	return ResourcePluralRedisOpsRequest
}

func (r RedisOpsRequest) ValidateSpecs() error {
	return nil
}

var _ Accessor = &RedisOpsRequest{}

func (r *RedisOpsRequest) GetObjectMeta() metav1.ObjectMeta {
	return r.ObjectMeta
}

func (r RedisOpsRequest) GetRequestType() any {
	return r.Spec.Type
}

func (r RedisOpsRequest) GetUpdateVersionSpec() *RedisUpdateVersionSpec {
	return r.Spec.UpdateVersion
}

func (r *RedisOpsRequest) GetDBRefName() string {
	return r.Spec.DatabaseRef.Name
}

func (r *RedisOpsRequest) GetStatus() OpsRequestStatus {
	return r.Status
}

func (r *RedisOpsRequest) SetStatus(s OpsRequestStatus) {
	r.Status = s
}
