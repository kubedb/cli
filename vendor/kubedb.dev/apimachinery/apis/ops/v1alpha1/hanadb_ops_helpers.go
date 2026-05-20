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

func (r *HanaDBOpsRequest) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralHanaDBOpsRequest))
}

var _ apis.ResourceInfo = &HanaDBOpsRequest{}

func (r *HanaDBOpsRequest) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", ResourcePluralHanaDBOpsRequest, ops.GroupName)
}

func (r *HanaDBOpsRequest) ResourceShortCode() string {
	return ResourceCodeHanaDBOpsRequest
}

func (r *HanaDBOpsRequest) ResourceKind() string {
	return ResourceKindHanaDBOpsRequest
}

func (r *HanaDBOpsRequest) ResourceSingular() string {
	return ResourceSingularHanaDBOpsRequest
}

func (r *HanaDBOpsRequest) ResourcePlural() string {
	return ResourcePluralHanaDBOpsRequest
}

var _ Accessor = &HanaDBOpsRequest{}

func (r *HanaDBOpsRequest) GetObjectMeta() metav1.ObjectMeta {
	return r.ObjectMeta
}

func (r *HanaDBOpsRequest) GetDBRefName() string {
	return r.Spec.DatabaseRef.Name
}

func (r *HanaDBOpsRequest) GetRequestType() string {
	return string(r.Spec.Type)
}

func (r *HanaDBOpsRequest) GetStatus() OpsRequestStatus {
	return r.Status
}

func (r *HanaDBOpsRequest) SetStatus(s OpsRequestStatus) {
	r.Status = s
}
