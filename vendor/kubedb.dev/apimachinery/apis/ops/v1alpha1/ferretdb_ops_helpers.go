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

func (r *FerretDBOpsRequest) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralFerretDBOpsRequest))
}

var _ apis.ResourceInfo = &FerretDBOpsRequest{}

func (r *FerretDBOpsRequest) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", ResourcePluralFerretDBOpsRequest, ops.GroupName)
}

func (r *FerretDBOpsRequest) ResourceShortCode() string {
	return ResourceCodeFerretDBOpsRequest
}

func (r *FerretDBOpsRequest) ResourceKind() string {
	return ResourceKindFerretDBOpsRequest
}

func (r *FerretDBOpsRequest) ResourceSingular() string {
	return ResourceSingularFerretDBOpsRequest
}

func (r *FerretDBOpsRequest) ResourcePlural() string {
	return ResourcePluralFerretDBOpsRequest
}

var _ Accessor = &FerretDBOpsRequest{}

func (r *FerretDBOpsRequest) GetObjectMeta() metav1.ObjectMeta {
	return r.ObjectMeta
}

func (r *FerretDBOpsRequest) GetDBRefName() string {
	return r.Spec.DatabaseRef.Name
}

func (r *FerretDBOpsRequest) GetRequestType() any {
	return r.Spec.Type
}

func (r *FerretDBOpsRequest) GetStatus() OpsRequestStatus {
	return r.Status
}

func (r *FerretDBOpsRequest) SetStatus(s OpsRequestStatus) {
	r.Status = s
}
