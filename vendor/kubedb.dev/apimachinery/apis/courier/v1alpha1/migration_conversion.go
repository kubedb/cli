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

	"k8s.io/apimachinery/pkg/runtime"
)

// Duckify projects one of the per-engine {DB}Migration underlying types into
// the Migration duck shape. It implements the duck.Object contract from
// kmodules.xyz/client-go/client/duck so the operator can reconcile every
// engine-specific CRD through a single, engine-agnostic Migration view.
//
// Only Get/List (raw -> duck projection) and Patch/Status().Patch (patch bytes
// computed against the duck, replayed onto the raw) are used by the duck
// client; there is intentionally no reverse (duck -> raw) conversion. Because
// ObjectMeta and Status are copied through unchanged, status and finalizer
// patches land on identical JSON paths in the underlying object.
func (m *Migration) Duckify(srcRaw runtime.Object) error {
	// APIVersion is shared by all underlying kinds; Kind is set per case below so
	// the projected Migration retains the provenance of the concrete CRD it was
	// duckified from.
	m.APIVersion = SchemeGroupVersion.String()

	switch t := srcRaw.(type) {
	case *PostgresMigration:
		m.Kind = ResourceKindPostgresMigration
		m.ObjectMeta = t.ObjectMeta
		m.Status = t.Status
		m.Spec = MigrationSpec{
			Source:      &Source{PostgresSource: &t.Spec.Source},
			Target:      &Target{PostgresTarget: &t.Spec.Target},
			JobDefaults: t.Spec.JobDefaults,
			JobTemplate: t.Spec.JobTemplate,
		}
	case *MySQLMigration:
		m.Kind = ResourceKindMySQLMigration
		m.ObjectMeta = t.ObjectMeta
		m.Status = t.Status
		m.Spec = MigrationSpec{
			Source:      &Source{MySQLSource: &t.Spec.Source},
			Target:      &Target{MySQLTarget: &t.Spec.Target},
			JobDefaults: t.Spec.JobDefaults,
			JobTemplate: t.Spec.JobTemplate,
		}
	case *MariaDBMigration:
		m.Kind = ResourceKindMariaDBMigration
		m.ObjectMeta = t.ObjectMeta
		m.Status = t.Status
		m.Spec = MigrationSpec{
			Source:      &Source{MariaDBSource: &t.Spec.Source},
			Target:      &Target{MariaDBTarget: &t.Spec.Target},
			JobDefaults: t.Spec.JobDefaults,
			JobTemplate: t.Spec.JobTemplate,
		}
	case *MongoDBMigration:
		m.Kind = ResourceKindMongoDBMigration
		m.ObjectMeta = t.ObjectMeta
		m.Status = t.Status
		m.Spec = MigrationSpec{
			Source:      &Source{MongoDBSource: &t.Spec.Source},
			Target:      &Target{MongoDBTarget: &t.Spec.Target},
			JobDefaults: t.Spec.JobDefaults,
			JobTemplate: t.Spec.JobTemplate,
		}
	case *MSSQLServerMigration:
		m.Kind = ResourceKindMSSQLServerMigration
		m.ObjectMeta = t.ObjectMeta
		m.Status = t.Status
		m.Spec = MigrationSpec{
			Source:      &Source{MSSQLServerSource: &t.Spec.Source},
			Target:      &Target{MSSQLServerTarget: &t.Spec.Target},
			JobDefaults: t.Spec.JobDefaults,
			JobTemplate: t.Spec.JobTemplate,
		}
	default:
		return fmt.Errorf("courier: cannot Duckify %T", srcRaw)
	}

	return nil
}
