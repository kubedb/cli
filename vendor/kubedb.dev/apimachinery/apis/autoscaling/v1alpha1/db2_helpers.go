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
	"kubedb.dev/apimachinery/apis/autoscaling"
	"kubedb.dev/apimachinery/crds"

	"kmodules.xyz/client-go/apiextensions"
)

func (DB2Autoscaler) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralDB2Autoscaler))
}

var _ apis.ResourceInfo = &DB2Autoscaler{}

func (d DB2Autoscaler) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", ResourcePluralDB2Autoscaler, autoscaling.GroupName)
}

func (d DB2Autoscaler) ResourceShortCode() string {
	return ResourceCodeDB2Autoscaler
}

func (d DB2Autoscaler) ResourceKind() string {
	return ResourceKindDB2Autoscaler
}

func (d DB2Autoscaler) ResourceSingular() string {
	return ResourceSingularDB2Autoscaler
}

func (d DB2Autoscaler) ResourcePlural() string {
	return ResourcePluralDB2Autoscaler
}

func (d DB2Autoscaler) ValidateSpecs() error {
	return nil
}

var _ StatusAccessor = &DB2Autoscaler{}

func (d *DB2Autoscaler) GetStatus() AutoscalerStatus {
	return d.Status
}

func (d *DB2Autoscaler) SetStatus(s AutoscalerStatus) {
	d.Status = s
}
