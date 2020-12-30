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
	"kubedb.dev/apimachinery/apis/ops"
	"kubedb.dev/apimachinery/crds"

	"kmodules.xyz/client-go/apiextensions"
)

func (_ PostgresOpsRequest) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralPostgresOpsRequest))
}

var _ apis.ResourceInfo = &PostgresOpsRequest{}

func (p PostgresOpsRequest) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", ResourcePluralPostgresOpsRequest, ops.GroupName)
}

func (p PostgresOpsRequest) ResourceShortCode() string {
	return ResourceCodePostgresOpsRequest
}

func (p PostgresOpsRequest) ResourceKind() string {
	return ResourceKindPostgresOpsRequest
}

func (p PostgresOpsRequest) ResourceSingular() string {
	return ResourceSingularPostgresOpsRequest
}

func (p PostgresOpsRequest) ResourcePlural() string {
	return ResourcePluralPostgresOpsRequest
}

func (p PostgresOpsRequest) ValidateSpecs() error {
	return nil
}
