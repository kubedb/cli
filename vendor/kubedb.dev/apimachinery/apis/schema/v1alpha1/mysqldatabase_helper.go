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

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kmodules.xyz/client-go/apiextensions"
	"kmodules.xyz/client-go/meta"
)

const (
	MySQLSuffix string = "mysql"
)

func (in MySQLDatabase) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourceMySQLDatabases))
}

var _ Interface = &MySQLDatabase{}

func (in *MySQLDatabase) GetInit() *InitSpec {
	return in.Spec.Init
}

func (in *MySQLDatabase) GetStatus() DatabaseStatus {
	return in.Status
}

// GetAppBindingMeta returns meta info of the appbinding which has been created by schema manager
func (in *MySQLDatabase) GetAppBindingMeta() metav1.ObjectMeta {
	meta := metav1.ObjectMeta{
		Name:      meta.NameWithSuffix(in.Name, MySQLSuffix+"-apbng"),
		Namespace: in.Namespace,
	}
	return meta
}

// GetVaultSecretEngineMeta returns meta info of the secret engine which has been created by schema manager
func (in *MySQLDatabase) GetVaultSecretEngineMeta() metav1.ObjectMeta {
	meta := metav1.ObjectMeta{
		Name:      meta.NameWithSuffix(in.Name, MySQLSuffix+"-engine"),
		Namespace: in.Namespace,
	}
	return meta
}

// GetMySQLRoleMeta returns meta info of the MySQL role which has been created by schema manager
func (in *MySQLDatabase) GetMySQLRoleMeta() metav1.ObjectMeta {
	meta := metav1.ObjectMeta{
		Name:      meta.NameWithSuffix(in.Name, MySQLSuffix+"-role"),
		Namespace: in.Namespace,
	}
	return meta
}

// GetSecretAccessRequestMeta returns meta info of the secret access request which has been created by schema manager
func (in *MySQLDatabase) GetSecretAccessRequestMeta() metav1.ObjectMeta {
	meta := metav1.ObjectMeta{
		Name:      meta.NameWithSuffix(in.Name, MySQLSuffix+"-req"),
		Namespace: in.Namespace,
	}
	return meta
}

// GetInitJobMeta returns meta info of the init job which has been created by schema manager
func (in *MySQLDatabase) GetInitJobMeta() metav1.ObjectMeta {
	meta := metav1.ObjectMeta{
		Name:      meta.NameWithSuffix(in.Name, MySQLSuffix+"-job"),
		Namespace: in.Namespace,
	}
	return meta
}

// GetMySQLAuthSecretMeta returns meta info of the mysql auth secret
func (in *MySQLDatabase) GetMySQLAuthSecretMeta() metav1.ObjectMeta {
	meta := metav1.ObjectMeta{
		Name:      in.Spec.Database.ServerRef.Name + "-auth",
		Namespace: in.Spec.Database.ServerRef.Namespace,
	}
	return meta
}

// GetRestoreSessionMeta returns meta info of the restore session which has been created by schema manager
func (in *MySQLDatabase) GetRestoreSessionMeta() metav1.ObjectMeta {
	meta := metav1.ObjectMeta{
		Name:      meta.NameWithSuffix(in.Name, MySQLSuffix+"-rs"),
		Namespace: in.Namespace,
	}
	return meta
}

// GetRepositoryMeta returns meta info of the repository which has been created by schema manager
func (in *MySQLDatabase) GetRepositoryMeta() metav1.ObjectMeta {
	meta := metav1.ObjectMeta{
		Name:      meta.NameWithSuffix(in.Name, MySQLSuffix+"-rp"),
		Namespace: in.Namespace,
	}
	return meta
}

// GetRepositorySecretMeta returns meta info of the repository which has been created by schema manager
func (in *MySQLDatabase) GetRepositorySecretMeta() metav1.ObjectMeta {
	meta := metav1.ObjectMeta{
		Name:      meta.NameWithSuffix(in.Name, MySQLSuffix+"-rp-sec"),
		Namespace: in.Namespace,
	}
	return meta
}
