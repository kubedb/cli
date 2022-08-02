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

func (m MariaDBVersion) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralMariaDBVersion))
}

var _ apis.ResourceInfo = &MariaDBVersion{}

func (m MariaDBVersion) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", ResourcePluralMariaDBVersion, catalog.GroupName)
}

func (m MariaDBVersion) ResourceShortCode() string {
	return ResourceCodeMariaDBVersion
}

func (m MariaDBVersion) ResourceKind() string {
	return ResourceKindMariaDBVersion
}

func (m MariaDBVersion) ResourceSingular() string {
	return ResourceSingularMariaDBVersion
}

func (m MariaDBVersion) ResourcePlural() string {
	return ResourcePluralMariaDBVersion
}

func (m MariaDBVersion) ValidateSpecs() error {
	if m.Spec.Version == "" ||
		m.Spec.DB.Image == "" ||
		m.Spec.Exporter.Image == "" ||
		m.Spec.InitContainer.Image == "" ||
		m.Spec.Coordinator.Image == "" {
		return fmt.Errorf(`atleast one of the following specs is not set for mariadbversion "%v":
spec.version,
spec.db.image,
spec.exporter.image,
spec.initContainer.image,
spec.coordinator.image.`, m.Name)
	}
	return nil
}
