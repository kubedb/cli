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

func (s SinglestoreVersion) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralSinglestoreVersion))
}

var _ apis.ResourceInfo = &SinglestoreVersion{}

func (s SinglestoreVersion) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", ResourcePluralSinglestoreVersion, catalog.GroupName)
}

func (s SinglestoreVersion) ResourceShortCode() string {
	return ResourceCodeSinglestoreVersion
}

func (s SinglestoreVersion) ResourceKind() string {
	return ResourceKindSinlestoreVersion
}

func (s SinglestoreVersion) ResourceSingular() string {
	return ResourceSingularSinglestoreVersion
}

func (s SinglestoreVersion) ResourcePlural() string {
	return ResourcePluralSinglestoreVersion
}

func (s SinglestoreVersion) ValidateSpecs() error {
	if s.Spec.Version == "" ||
		s.Spec.DB.Image == "" ||
		s.Spec.Coordinator.Image == "" ||
		s.Spec.Standalone.Image == "" {
		return fmt.Errorf(`atleast one of the following specs is not set for singlestoreVersion "%v":
spec.version,
spec.coordinator.image,
spec.standalone.image.`, s.Name)
	}
	return nil
}
