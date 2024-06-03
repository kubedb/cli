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

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kmodules.xyz/client-go/apiextensions"
)

func (d *DruidOpsRequest) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralDruidOpsRequest))
}

var _ apis.ResourceInfo = &DruidOpsRequest{}

func (d *DruidOpsRequest) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", ResourcePluralDruidOpsRequest, ops.GroupName)
}

func (d *DruidOpsRequest) ResourceShortCode() string {
	return ResourceCodeDruidOpsRequest
}

func (d *DruidOpsRequest) ResourceKind() string {
	return ResourceKindDruidOpsRequest
}

func (d *DruidOpsRequest) ResourceSingular() string {
	return ResourceSingularDruidOpsRequest
}

func (d *DruidOpsRequest) ResourcePlural() string {
	return ResourcePluralDruidOpsRequest
}

var _ Accessor = &DruidOpsRequest{}

func (d *DruidOpsRequest) GetObjectMeta() metav1.ObjectMeta {
	return d.ObjectMeta
}

func (d *DruidOpsRequest) GetDBRefName() string {
	return d.Spec.DatabaseRef.Name
}

func (d *DruidOpsRequest) GetRequestType() any {
	return d.Spec.Type
}

func (d *DruidOpsRequest) GetStatus() OpsRequestStatus {
	return d.Status
}

func (d *DruidOpsRequest) SetStatus(s OpsRequestStatus) {
	d.Status = s
}
