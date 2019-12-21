/*
Copyright The KubeDB Authors.

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

	"kubedb.dev/apimachinery/api/crds"
	"kubedb.dev/apimachinery/apis"

	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
)

func (_ PerconaXtraDBVersion) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralPerconaXtraDBVersion))
}

var _ apis.ResourceInfo = &PerconaXtraDBVersion{}

func (p PerconaXtraDBVersion) ResourceShortCode() string {
	return ResourceCodePerconaXtraDBVersion
}

func (p PerconaXtraDBVersion) ResourceKind() string {
	return ResourceKindPerconaXtraDBVersion
}

func (p PerconaXtraDBVersion) ResourceSingular() string {
	return ResourceSingularPerconaXtraDBVersion
}

func (p PerconaXtraDBVersion) ResourcePlural() string {
	return ResourcePluralPerconaXtraDBVersion
}

func (p PerconaXtraDBVersion) ValidateSpecs() error {
	if p.Spec.Version == "" ||
		p.Spec.DB.Image == "" ||
		p.Spec.Exporter.Image == "" ||
		p.Spec.InitContainer.Image == "" {
		return fmt.Errorf(`atleast one of the following specs is not set for perconaxtradbversion "%v":
spec.version,
spec.db.image,
spec.exporter.image,
spec.initContainer.image.`, p.Name)
	}
	return nil
}
