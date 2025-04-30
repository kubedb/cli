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

func (_ IgniteVersion) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralIgniteVersion))
}

var _ apis.ResourceInfo = &IgniteVersion{}

func (m IgniteVersion) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", ResourcePluralIgniteVersion, catalog.GroupName)
}

func (m IgniteVersion) ResourceShortCode() string {
	return ResourceCodeIgniteVersion
}

func (m IgniteVersion) ResourceKind() string {
	return ResourceKindIgniteVersion
}

func (m IgniteVersion) ResourceSingular() string {
	return ResourceSingularIgniteVersion
}

func (m IgniteVersion) ResourcePlural() string {
	return ResourcePluralIgniteVersion
}

func (m IgniteVersion) ValidateSpecs() error {
	if m.Spec.Version == "" ||
		m.Spec.DB.Image == "" {
		return fmt.Errorf(`atleast one of the following specs is not set for IgniteVersion "%v":
spec.version,
spec.db.image,`, m.Name)
	}
	return nil
}
