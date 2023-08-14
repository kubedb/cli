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
	kmapi "kmodules.xyz/client-go/api/v1"
	cutil "kmodules.xyz/client-go/conditions"
)

func GetPhase(obj Interface) DatabaseSchemaPhase {
	conditions := obj.GetStatus().Conditions

	if !obj.GetDeletionTimestamp().IsZero() {
		return DatabaseSchemaPhaseTerminating
	}
	if CheckIfSecretExpired(conditions) {
		return DatabaseSchemaPhaseExpired
	}
	if cutil.IsConditionTrue(conditions, string(DatabaseSchemaConditionTypeDBCreationUnsuccessful)) {
		return DatabaseSchemaPhaseFailed
	}

	// If Database or vault is not in ready state, Phase is 'Pending'
	if !cutil.IsConditionTrue(conditions, string(DatabaseSchemaConditionTypeDBServerReady)) ||
		!cutil.IsConditionTrue(conditions, string(DatabaseSchemaConditionTypeVaultReady)) {
		return DatabaseSchemaPhasePending
	}

	if cutil.IsConditionTrue(conditions, string(DatabaseSchemaConditionTypeDoubleOptInNotPossible)) {
		return DatabaseSchemaPhaseFailed
	}

	// If SecretEngine or Role is not in ready state, Phase is 'InProgress'
	if !cutil.IsConditionTrue(conditions, string(DatabaseSchemaConditionTypeSecretEngineReady)) ||
		!cutil.IsConditionTrue(conditions, string(DatabaseSchemaConditionTypeRoleReady)) ||
		!cutil.IsConditionTrue(conditions, string(DatabaseSchemaConditionTypeSecretAccessRequestReady)) {
		return DatabaseSchemaPhaseInProgress
	}
	// we are here means, SecretAccessRequest is approved and not expired. Now handle Init-Restore cases.

	if !cutil.IsConditionTrue(conditions, string(DatabaseSchemaConditionTypeAppBindingFound)) {
		return DatabaseSchemaPhaseInProgress
	}

	if cutil.HasCondition(conditions, string(DatabaseSchemaConditionTypeRepositoryFound)) {
		//  ----------------------------- Restore case -----------------------------
		if !cutil.IsConditionTrue(conditions, string(DatabaseSchemaConditionTypeRepositoryFound)) ||
			!cutil.IsConditionTrue(conditions, string(DatabaseSchemaConditionTypeRestoreCompleted)) {
			return DatabaseSchemaPhaseInProgress
		}
		if CheckIfRestoreFailed(conditions) {
			return DatabaseSchemaPhaseFailed
		} else {
			return DatabaseSchemaPhaseCurrent
		}
	} else if cutil.HasCondition(conditions, string(DatabaseSchemaConditionTypeInitScriptCompleted)) {
		//  ----------------------------- Init case -----------------------------
		if !cutil.IsConditionTrue(conditions, string(DatabaseSchemaConditionTypeInitScriptCompleted)) {
			return DatabaseSchemaPhaseInProgress
		}
		if CheckIfInitScriptFailed(conditions) {
			return DatabaseSchemaPhaseFailed
		} else {
			return DatabaseSchemaPhaseCurrent
		}
	}
	return DatabaseSchemaPhaseCurrent
}

func CheckIfInitScriptFailed(conditions []kmapi.Condition) bool {
	_, cond := cutil.GetCondition(conditions, string(DatabaseSchemaConditionTypeInitScriptCompleted))
	return cond.Message == string(DatabaseSchemaMessageInitScriptFailed)
}

func CheckIfRestoreFailed(conditions []kmapi.Condition) bool {
	_, cond := cutil.GetCondition(conditions, string(DatabaseSchemaConditionTypeRestoreCompleted))
	return cond.Message == string(DatabaseSchemaMessageRestoreSessionFailed)
}

func CheckIfSecretExpired(conditions []kmapi.Condition) bool {
	i, cond := cutil.GetCondition(conditions, string(DatabaseSchemaConditionTypeSecretAccessRequestReady))
	if i == -1 {
		return false
	}
	return cond.Message == string(DatabaseSchemaMessageSecretAccessRequestExpired)
}

func GetFinalizerForSchema() string {
	return SchemeGroupVersion.Group
}

func GetSchemaDoubleOptInLabelKey() string {
	return SchemeGroupVersion.Group + "/doubleoptin"
}

func GetSchemaDoubleOptInLabelValue() string {
	return "enabled"
}
