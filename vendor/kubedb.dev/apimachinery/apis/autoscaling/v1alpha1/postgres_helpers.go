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

func (_ PostgresAutoscaler) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralPostgresAutoscaler))
}

var _ apis.ResourceInfo = &PostgresAutoscaler{}

func (p PostgresAutoscaler) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", ResourcePluralPostgresAutoscaler, autoscaling.GroupName)
}

func (p PostgresAutoscaler) ResourceShortCode() string {
	return ResourceCodePostgresAutoscaler
}

func (p PostgresAutoscaler) ResourceKind() string {
	return ResourceKindPostgresAutoscaler
}

func (p PostgresAutoscaler) ResourceSingular() string {
	return ResourceSingularPostgresAutoscaler
}

func (p PostgresAutoscaler) ResourcePlural() string {
	return ResourcePluralPostgresAutoscaler
}

func (p PostgresAutoscaler) ValidateSpecs() error {
	return nil
}

var _ StatusAccessor = &PostgresAutoscaler{}

func (e *PostgresAutoscaler) GetStatus() AutoscalerStatus {
	return e.Status
}

func (e *PostgresAutoscaler) SetStatus(s AutoscalerStatus) {
	e.Status = s
}
