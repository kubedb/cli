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
	"context"

	dbapi "kubedb.dev/apimachinery/apis/kubedb/v1alpha2"
	"kubedb.dev/apimachinery/crds"

	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"kmodules.xyz/client-go/apiextensions"
	"kmodules.xyz/client-go/meta"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	MariaDBSuffix string = dbapi.ResourceSingularMariaDB
)

func (_ MariaDBDatabase) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourceMariaDBDatabases))
}

var _ Interface = &MariaDBDatabase{}

func (in *MariaDBDatabase) GetInit() *InitSpec {
	return in.Spec.Init
}

func (in *MariaDBDatabase) GetStatus() DatabaseStatus {
	return in.Status
}

// GetAppBindingMeta returns meta info of the appbinding which has been created by schema manager
func (in *MariaDBDatabase) GetAppBindingMeta() metav1.ObjectMeta {
	meta := metav1.ObjectMeta{
		Name:      meta.NameWithSuffix(in.Name, MariaDBSuffix+"-apbng"),
		Namespace: in.Namespace,
	}
	return meta
}

// GetVaultSecretEngineMeta returns meta info of the secret engine which has been created by schema manager
func (in *MariaDBDatabase) GetVaultSecretEngineMeta() metav1.ObjectMeta {
	meta := metav1.ObjectMeta{
		Name:      meta.NameWithSuffix(in.Name, MariaDBSuffix+"-engine"),
		Namespace: in.Namespace,
	}
	return meta
}

// GetMySQLRoleMeta returns meta info of the MySQL role which has been created by schema manager
func (in *MariaDBDatabase) GetMySQLRoleMeta() metav1.ObjectMeta {
	meta := metav1.ObjectMeta{
		Name:      meta.NameWithSuffix(in.Name, MariaDBSuffix+"-role"),
		Namespace: in.Namespace,
	}
	return meta
}

// GetSecretAccessRequestMeta returns meta info of the secret access request which has been created by schema manager
func (in *MariaDBDatabase) GetSecretAccessRequestMeta() metav1.ObjectMeta {
	meta := metav1.ObjectMeta{
		Name:      meta.NameWithSuffix(in.Name, MariaDBSuffix+"-req"),
		Namespace: in.Namespace,
	}
	return meta
}

// GetInitJobMeta returns meta info of the init job which has been created by schema manager
func (in *MariaDBDatabase) GetInitJobMeta() metav1.ObjectMeta {
	meta := metav1.ObjectMeta{
		Name:      meta.NameWithSuffix(in.Name, MariaDBSuffix+"-job"),
		Namespace: in.Namespace,
	}
	return meta
}

// GetMySQLAuthSecretMeta returns meta info of the mysql auth secret
func (in *MariaDBDatabase) GetMySQLAuthSecretMeta() metav1.ObjectMeta {
	meta := metav1.ObjectMeta{
		Name:      in.Spec.Database.ServerRef.Name + "-auth",
		Namespace: in.Spec.Database.ServerRef.Namespace,
	}
	return meta
}

// GetRestoreSessionMeta returns meta info of the restore session which has been created by schema manager
func (in *MariaDBDatabase) GetRestoreSessionMeta() metav1.ObjectMeta {
	meta := metav1.ObjectMeta{
		Name:      meta.NameWithSuffix(in.Name, MariaDBSuffix+"-rs"),
		Namespace: in.Namespace,
	}
	return meta
}

// GetRepositoryMeta returns meta info of the repository which has been created by schema manager
func (in *MariaDBDatabase) GetRepositoryMeta() metav1.ObjectMeta {
	meta := metav1.ObjectMeta{
		Name:      meta.NameWithSuffix(in.Name, MariaDBSuffix+"-rp"),
		Namespace: in.Namespace,
	}
	return meta
}

// GetRepositorySecretMeta returns meta info of the repository which has been created by schema manager
func (in *MariaDBDatabase) GetRepositorySecretMeta() metav1.ObjectMeta {
	meta := metav1.ObjectMeta{
		Name:      meta.NameWithSuffix(in.Name, MariaDBSuffix+"-rp-sec"),
		Namespace: in.Namespace,
	}
	return meta
}

func (in *MariaDBDatabase) CheckDoubleOptIn(ctx context.Context, client client.Client) (bool, error) {
	// Get updated MariaDBDatabase object
	var schema MariaDBDatabase
	err := client.Get(ctx, types.NamespacedName{
		Namespace: in.GetNamespace(),
		Name:      in.GetName(),
	}, &schema)
	if err != nil {
		return false, err
	}

	// Get the database server
	var maria dbapi.MariaDB
	err = client.Get(ctx, types.NamespacedName{
		Namespace: schema.Spec.Database.ServerRef.Namespace,
		Name:      schema.Spec.Database.ServerRef.Name,
	}, &maria)
	if err != nil {
		return false, err
	}

	if maria.Spec.AllowedSchemas == nil {
		return false, nil
	}

	// Get namespace object of the schema
	var nsSchema core.Namespace
	err = client.Get(ctx, types.NamespacedName{
		Name: schema.GetNamespace(),
	}, &nsSchema)
	if err != nil {
		return false, err
	}

	// Get namespace object of the Database server
	var nsDB core.Namespace
	err = client.Get(ctx, types.NamespacedName{
		Name: schema.Spec.Database.ServerRef.Namespace,
	}, &nsDB)
	if err != nil {
		return false, err
	}

	possible, err := CheckIfDoubleOptInPossible(schema.ObjectMeta, nsSchema.ObjectMeta, nsDB.ObjectMeta, maria.Spec.AllowedSchemas)
	if err != nil {
		return false, err
	}

	return possible, nil
}
