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

func (h *HazelcastOpsRequest) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralHazelcastOpsRequest))
}

var _ apis.ResourceInfo = &HazelcastOpsRequest{}

func (h *HazelcastOpsRequest) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", ResourcePluralHazelcastOpsRequest, ops.GroupName)
}

func (h *HazelcastOpsRequest) ResourceShortCode() string {
	return ResourceCodeHazelcastOpsRequest
}

func (h *HazelcastOpsRequest) ResourceKind() string {
	return ResourceKindHazelcastOpsRequest
}

func (h *HazelcastOpsRequest) ResourceSingular() string {
	return ResourceSingularHazelcastOpsRequest
}

func (h *HazelcastOpsRequest) ResourcePlural() string {
	return ResourcePluralHazelcastOpsRequest
}

var _ Accessor = &HazelcastOpsRequest{}

func (h *HazelcastOpsRequest) GetObjectMeta() metav1.ObjectMeta {
	return h.ObjectMeta
}

func (h *HazelcastOpsRequest) GetDBRefName() string {
	return h.Spec.DatabaseRef.Name
}

func (h *HazelcastOpsRequest) GetRequestType() string {
	return string(h.Spec.Type)
}

func (h *HazelcastOpsRequest) GetStatus() OpsRequestStatus {
	return h.Status
}

func (h *HazelcastOpsRequest) SetStatus(st OpsRequestStatus) {
	h.Status = st
}
