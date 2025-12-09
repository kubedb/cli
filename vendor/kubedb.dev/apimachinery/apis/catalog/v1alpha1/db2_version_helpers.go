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

func (DB2Version) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralDB2Version))
}

var _ apis.ResourceInfo = &DB2Version{}

func (d DB2Version) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", ResourcePluralDB2Version, catalog.GroupName)
}

func (d DB2Version) ResourceShortCode() string {
	return ResourceCodeDB2Version
}

func (d DB2Version) ResourceKind() string {
	return ResourceKindDB2Version
}

func (d DB2Version) ResourceSingular() string {
	return ResourceSingularDB2Version
}

func (d DB2Version) ResourcePlural() string {
	return ResourcePluralDB2Version
}

func (d DB2Version) ValidateSpecs() error {
	if d.Spec.Version == "" || d.Spec.DB.Image == "" || d.Spec.Coordinator.Image == "" {
		return fmt.Errorf(`at least one of the following specs is not set for DB2 "%v":
spec.version,
spec.coordinator.image`, d.Name)
	}
	// TODO: add m.spec.exporter.image check FOR monitoring
	return nil
}
