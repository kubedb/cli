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
)

func GetFinalizer() string {
	return SchemeGroupVersion.Group
}

// Migration is the duck type projected from the per-engine {DB}Migration
// underlying types; it is not served as its own CRD, so it has no
// CustomResourceDefinition helper. The concrete resources users create are the
// per-engine kinds below.

func (PostgresMigration) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralPostgresMigrations))
}

func (MySQLMigration) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralMySQLMigrations))
}

func (MariaDBMigration) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralMariaDBMigrations))
}

func (MongoDBMigration) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralMongoDBMigrations))
}

func (MSSQLServerMigration) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralMSSQLServerMigrations))
}

func (Branch) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralBranches))
}

func (BranchWork) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(SchemeGroupVersion.WithResource(ResourcePluralBranchWorks))
}

func (m Migration) GetDBKindAndCommand() (string, string) {
	switch {
	case m.Spec.Source.PostgresSource != nil && m.Spec.Target.PostgresTarget != nil:
		return "Postgres", "postgres"
	case m.Spec.Source.MongoDBSource != nil && m.Spec.Target.MongoDBTarget != nil:
		return "MongoDB", "mongodb"
	case m.Spec.Source.MySQLSource != nil && m.Spec.Target.MySQLTarget != nil:
		return "MySQL", "mysql"
	case m.Spec.Source.MariaDBSource != nil && m.Spec.Target.MariaDBTarget != nil:
		return "MariaDB", "mariadb"
	case m.Spec.Source.MSSQLServerSource != nil && m.Spec.Target.MSSQLServerTarget != nil:
		return "MSSQLServer", "mssqlserver"
	}

	return "", ""
}

func (m Migration) GetConnectionInfos() (*ConnectionInfo, *ConnectionInfo) {
	switch {
	case m.Spec.Source.PostgresSource != nil && m.Spec.Target.PostgresTarget != nil:
		return &m.Spec.Source.PostgresSource.ConnectionInfo, &m.Spec.Target.PostgresTarget.ConnectionInfo
	case m.Spec.Source.MongoDBSource != nil && m.Spec.Target.MongoDBTarget != nil:
		return &m.Spec.Source.MongoDBSource.ConnectionInfo, &m.Spec.Target.MongoDBTarget.ConnectionInfo
	case m.Spec.Source.MySQLSource != nil && m.Spec.Target.MySQLTarget != nil:
		return m.Spec.Source.MySQLSource.ConnectionInfo, m.Spec.Target.MySQLTarget.ConnectionInfo
	case m.Spec.Source.MariaDBSource != nil && m.Spec.Target.MariaDBTarget != nil:
		return m.Spec.Source.MariaDBSource.ConnectionInfo, m.Spec.Target.MariaDBTarget.ConnectionInfo
	case m.Spec.Source.MSSQLServerSource != nil && m.Spec.Target.MSSQLServerTarget != nil:
		src := &ConnectionInfo{
			AppBinding: m.Spec.Source.MSSQLServerSource.ConnectionInfo.AppBinding,
			DBName:     m.Spec.Source.MSSQLServerSource.ConnectionInfo.Database,
		}
		tgt := &ConnectionInfo{
			AppBinding: m.Spec.Target.MSSQLServerTarget.ConnectionInfo.AppBinding,
			DBName:     m.Spec.Target.MSSQLServerTarget.ConnectionInfo.Database,
		}
		return src, tgt
	}
	return nil, nil
}
