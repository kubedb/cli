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

func (m MongoDBOpsRequest) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralMongoDBOpsRequest))
}

var _ apis.ResourceInfo = &MongoDBOpsRequest{}

func (m MongoDBOpsRequest) ResourceFQN() string {
	return fmt.Sprintf("%s.%s", ResourcePluralMongoDBOpsRequest, ops.GroupName)
}

func (m MongoDBOpsRequest) ResourceShortCode() string {
	return ResourceCodeMongoDBOpsRequest
}

func (m MongoDBOpsRequest) ResourceKind() string {
	return ResourceKindMongoDBOpsRequest
}

func (m MongoDBOpsRequest) ResourceSingular() string {
	return ResourceSingularMongoDBOpsRequest
}

func (m MongoDBOpsRequest) ResourcePlural() string {
	return ResourcePluralMongoDBOpsRequest
}

func (m MongoDBOpsRequest) ValidateSpecs() error {
	return nil
}

var _ Accessor = &MongoDBOpsRequest{}

func (m *MongoDBOpsRequest) GetObjectMeta() metav1.ObjectMeta {
	return m.ObjectMeta
}

func (m MongoDBOpsRequest) GetRequestType() any {
	switch m.Spec.Type {
	case MongoDBOpsRequestTypeUpgrade:
		return MongoDBOpsRequestTypeUpdateVersion
	}
	return m.Spec.Type
}

func (m MongoDBOpsRequest) GetUpdateVersionSpec() *MongoDBUpdateVersionSpec {
	if m.Spec.UpdateVersion != nil {
		return m.Spec.UpdateVersion
	}
	return m.Spec.Upgrade
}

func (m *MongoDBOpsRequest) GetDBRefName() string {
	return m.Spec.DatabaseRef.Name
}

func (m *MongoDBOpsRequest) GetStatus() OpsRequestStatus {
	return m.Status
}

func (m *MongoDBOpsRequest) SetStatus(s OpsRequestStatus) {
	m.Status = s
}
