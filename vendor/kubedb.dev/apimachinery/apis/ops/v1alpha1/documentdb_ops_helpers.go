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

func (d DocumentDBOpsRequest) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralDocumentDBOpsRequest))
}

var _ apis.ResourceInfo = &DocumentDBOpsRequest{}

func (d DocumentDBOpsRequest) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", ResourcePluralDocumentDBOpsRequest, ops.GroupName)
}

func (d DocumentDBOpsRequest) ResourceShortCode() string {
	return ResourceCodeDocumentDBOpsRequest
}

func (d DocumentDBOpsRequest) ResourceKind() string {
	return ResourceKindDocumentDBOpsRequest
}

func (d DocumentDBOpsRequest) ResourceSingular() string {
	return ResourceSingularDocumentDBOpsRequest
}

func (d DocumentDBOpsRequest) ResourcePlural() string {
	return ResourcePluralDocumentDBOpsRequest
}

func (d DocumentDBOpsRequest) ValidateSpecs() error {
	return nil
}

var _ Accessor = &DocumentDBOpsRequest{}

func (d *DocumentDBOpsRequest) GetObjectMeta() metav1.ObjectMeta {
	return d.ObjectMeta
}

func (d *DocumentDBOpsRequest) GetDBRefName() string {
	return d.Spec.DatabaseRef.Name
}

func (d *DocumentDBOpsRequest) GetRequestType() string {
	return string(d.Spec.Type)
}

func (d *DocumentDBOpsRequest) GetStatus() OpsRequestStatus {
	return d.Status
}

func (d *DocumentDBOpsRequest) SetStatus(s OpsRequestStatus) {
	d.Status = s
}
