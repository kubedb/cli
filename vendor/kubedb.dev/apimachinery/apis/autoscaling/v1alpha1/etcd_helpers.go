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
	"kubedb.dev/apimachinery/apis"
	"kubedb.dev/apimachinery/crds"

	"kmodules.xyz/client-go/apiextensions"
)

func (_ EtcdAutoscaler) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralEtcdAutoscaler))
}

var _ apis.ResourceInfo = &EtcdAutoscaler{}

func (e EtcdAutoscaler) ResourceShortCode() string {
	return ResourceCodeEtcdAutoscaler
}

func (e EtcdAutoscaler) ResourceKind() string {
	return ResourceKindEtcdAutoscaler
}

func (e EtcdAutoscaler) ResourceSingular() string {
	return ResourceSingularEtcdAutoscaler
}

func (e EtcdAutoscaler) ResourcePlural() string {
	return ResourcePluralEtcdAutoscaler
}

func (e EtcdAutoscaler) ValidateSpecs() error {
	return nil
}
