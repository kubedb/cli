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
	"kubedb.dev/apimachinery/crds"

	"kmodules.xyz/client-go/apiextensions"
	"kmodules.xyz/client-go/meta"
)

const (
	InitScriptName              string = "init.js"
	MongoInitScriptPath         string = "init-scripts"
	MongoPrefix                 string = "MongoDB"
	MongoSuffix                 string = "mongo"
	MongoDatabaseNameForEntry   string = "kubedb-system"
	MongoCollectionNameForEntry string = "databases"
)

func (in MongoDBDatabase) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourceMongoDBDatabases))
}

var _ Interface = &MongoDBDatabase{}

func (in *MongoDBDatabase) GetInit() *InitSpec {
	return in.Spec.Init
}

func (in *MongoDBDatabase) GetStatus() DatabaseStatus {
	return in.Status
}

func (in *MongoDBDatabase) GetMongoInitVolumeNameForPod() string {
	return meta.NameWithSuffix(in.GetName(), MongoSuffix+"-vol")
}

func (in *MongoDBDatabase) GetMongoInitJobName() string {
	return meta.NameWithSuffix(in.GetName(), MongoSuffix+"-job")
}

func (in *MongoDBDatabase) GetMongoInitScriptContainerName() string {
	return meta.NameWithSuffix(in.GetName(), MongoSuffix)
}

func (in *MongoDBDatabase) GetMongoRestoreSessionName() string {
	return meta.NameWithSuffix(in.GetName(), MongoSuffix+"-rs")
}

func (in *MongoDBDatabase) GetMongoAdminRoleName() string {
	return meta.NameWithSuffix(in.GetName(), MongoSuffix+"-role")
}

func (in *MongoDBDatabase) GetMongoAdminSecretAccessRequestName() string {
	return meta.NameWithSuffix(in.GetName(), MongoSuffix+"-req")
}

func (in *MongoDBDatabase) GetMongoAdminServiceAccountName() string {
	return meta.NameWithSuffix(in.GetName(), MongoSuffix+"-sa")
}

func (in *MongoDBDatabase) GetMongoSecretEngineName() string {
	return meta.NameWithSuffix(in.GetName(), MongoSuffix+"-engine")
}

func (in *MongoDBDatabase) GetMongoAppBindingName() string {
	return meta.NameWithSuffix(in.GetName(), MongoSuffix+"-apbng")
}

func (in *MongoDBDatabase) GetAuthSecretName(dbServerName string) string {
	return meta.NameWithSuffix(dbServerName, "auth")
}
