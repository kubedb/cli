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
	"fmt"

	kdm "kubedb.dev/apimachinery/apis/kubedb/v1alpha2"
	"kubedb.dev/apimachinery/crds"

	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"kmodules.xyz/client-go/apiextensions"
	kmeta "kmodules.xyz/client-go/meta"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (_ PostgresDatabase) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePostgresDatabases))
}

var _ Interface = &PostgresDatabase{}

func (in *PostgresDatabase) GetInit() *InitSpec {
	return in.Spec.Init
}

func (in *PostgresDatabase) GetStatus() DatabaseStatus {
	return in.Status
}

const (
	EnvPGPassword            string = "PGPASSWORD"
	EnvPGUser                string = "PGUSER"
	PostgresSchemaKubeSystem string = "kube_system"
)

func GetPostgresSchemaFinalizerString() string {
	return SchemeGroupVersion.Group
}

func GetPostgresInitVolumeNameForPod(pgSchema *PostgresDatabase) string {
	return kmeta.NameWithSuffix(pgSchema.Name, kdm.ResourceSingularPostgres+"-vol")
}

func GetPostgresInitJobContainerName(pgSchema *PostgresDatabase) string {
	return kmeta.NameWithSuffix(pgSchema.Name, kdm.ResourceSingularPostgres)
}

func GetPostgresSchemaJobName(pgSchema *PostgresDatabase) string {
	return kmeta.NameWithSuffix(pgSchema.Name, kdm.ResourceSingularPostgres+"-job")
}

func GetPostgresHostName(db *kdm.Postgres) string {
	return fmt.Sprintf("%v.%v.svc", db.ServiceName(), db.Namespace)
}

func GetPostgresSchemaSecretEngineName(pgSchema *PostgresDatabase) string {
	return kmeta.NameWithSuffix(pgSchema.Name, kdm.ResourceSingularPostgres+"-engine")
}

func GetPostgresSchemaRoleName(pgSchema *PostgresDatabase) string {
	return kmeta.NameWithSuffix(pgSchema.Name, kdm.ResourceSingularPostgres+"-role")
}

func GetPostgresSchemaCreationStatements(pgSchema *PostgresDatabase) []string {
	createRole := "CREATE ROLE \"{{name}}\" WITH LOGIN PASSWORD '{{password}}';"
	grantRole := fmt.Sprintf("GRANT All PRIVILEGES ON DATABASE %s TO \"{{name}}\";", pgSchema.Spec.Database.Config.Name)
	return []string{createRole, grantRole}
}

func GetPostgresSchemaRevocationStatements(pgSchema *PostgresDatabase) []string {
	switchDatabase := fmt.Sprintf(`\c %s`, pgSchema.Spec.Database.Config.Name)
	revokeRole := fmt.Sprintf("REVOKE All PRIVILEGES ON DATABASE %s FROM \"{{name}}\";", pgSchema.Spec.Database.Config.Name)
	reassignOwned := "REASSIGN OWNED BY \"{{name}}\" TO POSTGRES;"
	dropOwned := "DROP OWNED BY \"{{name}}\";"
	dropRole := "DROP ROLE IF EXISTS \"{{name}}\";"
	return []string{switchDatabase, reassignOwned, revokeRole, dropOwned, dropRole}
}

func GetPostgresSchemaRoleSecretAccessName(pgSchema *PostgresDatabase) string {
	return kmeta.NameWithSuffix(pgSchema.Name, kdm.ResourceSingularPostgres+"-req")
}

func GetPostgresSchemaAppBinding(pgSchema *PostgresDatabase) string {
	return kmeta.NameWithSuffix(pgSchema.Name, kdm.ResourceSingularPostgres+"-appbdng")
}

func GetPostgresSchemaRestoreSessionName(pgSchema *PostgresDatabase) string {
	return kmeta.NameWithSuffix(pgSchema.Name, kdm.ResourceSingularPostgres+"-restore")
}

func GetPostgresSchemaSecretName(pgSchema *PostgresDatabase) string {
	return kmeta.NameWithSuffix(pgSchema.Name, kdm.ResourceSingularPostgres+"-secret")
}

func (in *PostgresDatabase) CheckDoubleOptIn(ctx context.Context, client client.Client) (bool, error) {
	// Get updated PostgresDatabase object
	var schema PostgresDatabase
	err := client.Get(ctx, types.NamespacedName{
		Namespace: in.GetNamespace(),
		Name:      in.GetName(),
	}, &schema)
	if err != nil {
		return false, err
	}

	// Get the database server
	var pg kdm.Postgres
	err = client.Get(ctx, types.NamespacedName{
		Namespace: schema.Spec.Database.ServerRef.Namespace,
		Name:      schema.Spec.Database.ServerRef.Name,
	}, &pg)
	if err != nil {
		return false, err
	}

	if pg.Spec.AllowedSchemas == nil {
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

	possible, err := CheckIfDoubleOptInPossible(schema.ObjectMeta, nsSchema.ObjectMeta, nsDB.ObjectMeta, pg.Spec.AllowedSchemas)
	if err != nil {
		return false, err
	}

	return possible, nil
}
