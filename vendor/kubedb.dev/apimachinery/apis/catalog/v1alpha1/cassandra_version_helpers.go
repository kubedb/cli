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

func (_ CassandraVersion) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralCassandraVersion))
}

var _ apis.ResourceInfo = &CassandraVersion{}

func (r CassandraVersion) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", ResourcePluralCassandraVersion, catalog.GroupName)
}

func (r CassandraVersion) ResourceShortCode() string {
	return ResourceCodeCassandraVersion
}

func (r CassandraVersion) ResourceKind() string {
	return ResourceKindCassandraVersion
}

func (r CassandraVersion) ResourceSingular() string {
	return ResourceSingularCassandraVersion
}

func (r CassandraVersion) ResourcePlural() string {
	return ResourcePluralCassandraVersion
}

func (r CassandraVersion) ValidateSpecs() error {
	if r.Spec.Version == "" ||
		r.Spec.DB.Image == "" {
		return fmt.Errorf(`atleast one of the following specs is not set for CassandraVersion "%v":
							spec.version,
							spec.db.image`, r.Name)
	}
	return nil
}
