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

func (r *DB2OpsRequest) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralDB2OpsRequest))
}

var _ apis.ResourceInfo = &DB2OpsRequest{}

func (r *DB2OpsRequest) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", ResourcePluralDB2OpsRequest, ops.GroupName)
}

func (r *DB2OpsRequest) ResourceShortCode() string {
	return ResourceCodeDB2OpsRequest
}

func (r *DB2OpsRequest) ResourceKind() string {
	return ResourceKindDB2OpsRequest
}

func (r *DB2OpsRequest) ResourceSingular() string {
	return ResourceSingularDB2OpsRequest
}

func (r *DB2OpsRequest) ResourcePlural() string {
	return ResourcePluralDB2OpsRequest
}

var _ Accessor = &DB2OpsRequest{}

func (r *DB2OpsRequest) GetObjectMeta() metav1.ObjectMeta {
	return r.ObjectMeta
}

func (r *DB2OpsRequest) GetDBRefName() string {
	return r.Spec.DatabaseRef.Name
}

func (r *DB2OpsRequest) GetRequestType() string {
	return string(r.Spec.Type)
}

func (r *DB2OpsRequest) GetStatus() OpsRequestStatus {
	return r.Status
}

func (r *DB2OpsRequest) SetStatus(s OpsRequestStatus) {
	r.Status = s
}
