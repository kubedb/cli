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
	"kubedb.dev/apimachinery/apis/catalog"
	"kubedb.dev/apimachinery/crds"

	"kmodules.xyz/client-go/apiextensions"
)

func (_ RedisVersion) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralRedisVersion))
}

var _ apis.ResourceInfo = &RedisVersion{}

func (r RedisVersion) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", ResourcePluralRedisVersion, catalog.GroupName)
}

func (r RedisVersion) ResourceShortCode() string {
	return ResourceCodeRedisVersion
}

func (r RedisVersion) ResourceKind() string {
	return ResourceKindRedisVersion
}

func (r RedisVersion) ResourceSingular() string {
	return ResourceSingularRedisVersion
}

func (r RedisVersion) ResourcePlural() string {
	return ResourcePluralRedisVersion
}

func (r RedisVersion) ValidateSpecs() error {
	if r.Spec.Version == "" ||
		r.Spec.DB.Image == "" ||
		r.Spec.Exporter.Image == "" {
		return fmt.Errorf(`atleast one of the following specs is not set for redisVersion "%v":
spec.version,
spec.db.image,
spec.exporter.image.`, r.Name)
	}
	return nil
}
