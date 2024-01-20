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

func (d DruidVersion) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralDruidVersion))
}

var _ apis.ResourceInfo = &DruidVersion{}

func (d DruidVersion) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", ResourcePluralDruidVersion, catalog.GroupName)
}

func (d DruidVersion) ResourceShortCode() string {
	return ResourceCodeDruidVersion
}

func (d DruidVersion) ResourceKind() string {
	return ResourceKindDruidVersion
}

func (d DruidVersion) ResourceSingular() string {
	return ResourceSingularDruidVersion
}

func (d DruidVersion) ResourcePlural() string {
	return ResourcePluralDruidVersion
}

func (d DruidVersion) ValidateSpecs() error {
	if d.Spec.Version == "" ||
		d.Spec.DB.Image == "" {
		return fmt.Errorf(`atleast one of the following specs is not set for druidVersion "%v":
spec.version,
spec.db.image`, d.Name)
	}
	return nil
}
