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

func (_ *SchemaRegistryVersion) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralSchemaRegistryVersion))
}

var _ apis.ResourceInfo = &SchemaRegistryVersion{}

func (r *SchemaRegistryVersion) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", ResourcePluralSchemaRegistryVersion, catalog.GroupName)
}

func (r *SchemaRegistryVersion) ResourceShortCode() string {
	return ResourceCodeSchemaRegistryVersion
}

func (r *SchemaRegistryVersion) ResourceKind() string {
	return ResourceKindSchemaRegistryVersion
}

func (r *SchemaRegistryVersion) ResourceSingular() string {
	return ResourceSingularSchemaRegistryVersion
}

func (r *SchemaRegistryVersion) ResourcePlural() string {
	return ResourcePluralSchemaRegistryVersion
}

func (r *SchemaRegistryVersion) ValidateSpecs() error {
	if r.Spec.Version == "" ||
		r.Spec.Registry.Image == "" ||
		r.Spec.InMemory.Image == "" {
		return fmt.Errorf(`atleast one of the following specs is not set for schemaRegistryVersion "%v":
							spec.version,
							spec.registry.image, r.inMemory.image`, r.Name)
	}
	return nil
}
