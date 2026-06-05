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

func (w *WeaviateOpsRequest) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralWeaviateOpsRequest))
}

var _ apis.ResourceInfo = &WeaviateOpsRequest{}

func (w *WeaviateOpsRequest) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", ResourcePluralWeaviateOpsRequest, ops.GroupName)
}

func (w *WeaviateOpsRequest) ResourceShortCode() string {
	return ResourceCodeWeaviateOpsRequest
}

func (w *WeaviateOpsRequest) ResourceKind() string {
	return ResourceKindWeaviateOpsRequest
}

func (w *WeaviateOpsRequest) ResourceSingular() string {
	return ResourceSingularWeaviateOpsRequest
}

func (w *WeaviateOpsRequest) ResourcePlural() string {
	return ResourcePluralWeaviateOpsRequest
}

var _ Accessor = &WeaviateOpsRequest{}

func (w *WeaviateOpsRequest) GetObjectMeta() metav1.ObjectMeta {
	return w.ObjectMeta
}

func (w *WeaviateOpsRequest) GetDBRefName() string {
	return w.Spec.DatabaseRef.Name
}

func (w *WeaviateOpsRequest) GetRequestType() string {
	return string(w.Spec.Type)
}

func (w *WeaviateOpsRequest) GetStatus() OpsRequestStatus {
	return w.Status
}

func (w *WeaviateOpsRequest) SetStatus(s OpsRequestStatus) {
	w.Status = s
}
