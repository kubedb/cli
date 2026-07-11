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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kmapi "kmodules.xyz/client-go/api/v1"
	cutil "kmodules.xyz/client-go/conditions"
)

func SetMigratorJobTriggeredConditionToTrue(migrator *Migration) {
	newCond := kmapi.Condition{
		Type:    MigratorJobTriggered,
		Status:  metav1.ConditionTrue,
		Message: "Migration job has been triggered.",
	}
	migrator.Status.Conditions = cutil.SetCondition(migrator.Status.Conditions, newCond)
}

func SetMigratorJobTriggeredConditionToFalse(migrator *Migration, err error) {
	newCond := kmapi.Condition{
		Type:    MigratorJobTriggered,
		Status:  metav1.ConditionFalse,
		Message: err.Error(),
	}
	migrator.Status.Conditions = cutil.SetCondition(migrator.Status.Conditions, newCond)
}

func (m *Migration) CalculatePhase() MigrationPhase {
	if cutil.IsConditionTrue(m.Status.Conditions, MigrationSucceeded) {
		return MigrationPhaseSucceeded
	}
	if cutil.IsConditionTrue(m.Status.Conditions, MigrationFailed) {
		return MigrationPhaseFailed
	}
	if cutil.IsConditionTrue(m.Status.Conditions, MigrationRunning) {
		return MigrationPhaseRunning
	}
	return MigrationPhasePending
}

// SetMigrationRunningCondition sets the condition indicating migration is in progress
func SetMigrationRunningCondition(migrator *Migration) {
	newCond := kmapi.Condition{
		Type:    MigrationRunning,
		Status:  metav1.ConditionTrue,
		Reason:  ReasonMigrationRunning,
		Message: "Migration is currently in progress.",
	}
	migrator.Status.Conditions = cutil.SetCondition(migrator.Status.Conditions, newCond)
}

// SetMigrationSucceededCondition sets the condition indicating migration completed successfully
func SetMigrationSucceededCondition(migrator *Migration) {
	newCond := kmapi.Condition{
		Type:    MigrationSucceeded,
		Status:  metav1.ConditionTrue,
		Reason:  ReasonMigrationSucceeded,
		Message: "Migration completed successfully.",
	}
	migrator.Status.Conditions = cutil.SetCondition(migrator.Status.Conditions, newCond)

	// Clear running condition
	clearCond := kmapi.Condition{
		Type:    MigrationRunning,
		Status:  metav1.ConditionFalse,
		Reason:  ReasonMigrationSucceeded,
		Message: "Migration completed.",
	}
	migrator.Status.Conditions = cutil.SetCondition(migrator.Status.Conditions, clearCond)
}

// SetMigrationFailedCondition sets the condition indicating migration failed
func SetMigrationFailedCondition(migrator *Migration, err error) {
	newCond := kmapi.Condition{
		Type:    MigrationFailed,
		Status:  metav1.ConditionTrue,
		Reason:  ReasonMigrationFailed,
		Message: err.Error(),
	}
	migrator.Status.Conditions = cutil.SetCondition(migrator.Status.Conditions, newCond)

	// Clear running condition
	clearCond := kmapi.Condition{
		Type:    MigrationRunning,
		Status:  metav1.ConditionFalse,
		Reason:  ReasonMigrationFailed,
		Message: "Migration failed.",
	}
	migrator.Status.Conditions = cutil.SetCondition(migrator.Status.Conditions, clearCond)
}
