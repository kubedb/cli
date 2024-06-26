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

func (_ ClickHouseVersion) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralClickHouseVersion))
}

var _ apis.ResourceInfo = &ClickHouseVersion{}

func (r ClickHouseVersion) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", ResourcePluralClickHouseVersion, catalog.GroupName)
}

func (r ClickHouseVersion) ResourceShortCode() string {
	return ResourceCodeClickHouseVersion
}

func (r ClickHouseVersion) ResourceKind() string {
	return ResourceKindClickHouseVersion
}

func (r ClickHouseVersion) ResourceSingular() string {
	return ResourceSingularClickHouseVersion
}

func (r ClickHouseVersion) ResourcePlural() string {
	return ResourcePluralClickHouseVersion
}

func (r ClickHouseVersion) ValidateSpecs() error {
	if r.Spec.Version == "" ||
		r.Spec.DB.Image == "" {
		return fmt.Errorf(`atleast one of the following specs is not set for ClickHouseVersion "%v":
							spec.version,
							spec.db.image`, r.Name)
	}
	return nil
}
