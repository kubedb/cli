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
	dbapi "kubedb.dev/apimachinery/apis/kubedb/v1alpha2"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/klog/v2"
	kmapi "kmodules.xyz/client-go/api/v1"
)

func GetPhase(obj Interface) DatabaseSchemaPhase {
	conditions := obj.GetStatus().Conditions

	if !obj.GetDeletionTimestamp().IsZero() {
		return DatabaseSchemaPhaseTerminating
	}
	if CheckIfSecretExpired(conditions) {
		return DatabaseSchemaPhaseExpired
	}
	if kmapi.IsConditionTrue(conditions, string(DatabaseSchemaConditionTypeDBCreationUnsuccessful)) {
		return DatabaseSchemaPhaseFailed
	}

	// If Database or vault is not in ready state, Phase is 'Pending'
	if !kmapi.IsConditionTrue(conditions, string(DatabaseSchemaConditionTypeDBServerReady)) ||
		!kmapi.IsConditionTrue(conditions, string(DatabaseSchemaConditionTypeVaultReady)) {
		return DatabaseSchemaPhasePending
	}

	if kmapi.IsConditionTrue(conditions, string(DatabaseSchemaConditionTypeDoubleOptInNotPossible)) {
		return DatabaseSchemaPhaseFailed
	}

	// If SecretEngine or Role is not in ready state, Phase is 'InProgress'
	if !kmapi.IsConditionTrue(conditions, string(DatabaseSchemaConditionTypeSecretEngineReady)) ||
		!kmapi.IsConditionTrue(conditions, string(DatabaseSchemaConditionTypeRoleReady)) ||
		!kmapi.IsConditionTrue(conditions, string(DatabaseSchemaConditionTypeSecretAccessRequestReady)) {
		return DatabaseSchemaPhaseInProgress
	}
	// we are here means, SecretAccessRequest is approved and not expired. Now handle Init-Restore cases.

	if !kmapi.IsConditionTrue(conditions, string(DatabaseSchemaConditionTypeAppBindingFound)) {
		return DatabaseSchemaPhaseInProgress
	}

	if kmapi.HasCondition(conditions, string(DatabaseSchemaConditionTypeRepositoryFound)) {
		//  ----------------------------- Restore case -----------------------------
		if !kmapi.IsConditionTrue(conditions, string(DatabaseSchemaConditionTypeRepositoryFound)) ||
			!kmapi.IsConditionTrue(conditions, string(DatabaseSchemaConditionTypeRestoreCompleted)) {
			return DatabaseSchemaPhaseInProgress
		}
		if CheckIfRestoreFailed(conditions) {
			return DatabaseSchemaPhaseFailed
		} else {
			return DatabaseSchemaPhaseCurrent
		}
	} else if kmapi.HasCondition(conditions, string(DatabaseSchemaConditionTypeInitScriptCompleted)) {
		//  ----------------------------- Init case -----------------------------
		if !kmapi.IsConditionTrue(conditions, string(DatabaseSchemaConditionTypeInitScriptCompleted)) {
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
	_, cond := kmapi.GetCondition(conditions, string(DatabaseSchemaConditionTypeInitScriptCompleted))
	return cond.Message == string(DatabaseSchemaMessageInitScriptFailed)
}

func CheckIfRestoreFailed(conditions []kmapi.Condition) bool {
	_, cond := kmapi.GetCondition(conditions, string(DatabaseSchemaConditionTypeRestoreCompleted))
	return cond.Message == string(DatabaseSchemaMessageRestoreSessionFailed)
}

func CheckIfSecretExpired(conditions []kmapi.Condition) bool {
	i, cond := kmapi.GetCondition(conditions, string(DatabaseSchemaConditionTypeSecretAccessRequestReady))
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

// CheckIfDoubleOptInPossible is the intended function to be called from operator
// It checks if the namespace, where SchemaDatabase is applied, is allowed.
// It also checks the labels of schemaDatabase, to decide if that is allowed or not.
func CheckIfDoubleOptInPossible(schemaMeta metav1.ObjectMeta, schemaNSMeta metav1.ObjectMeta, dbNSMeta metav1.ObjectMeta, consumers *dbapi.AllowedConsumers) (bool, error) {
	if consumers == nil {
		return false, nil
	}
	matchNamespace, err := IsInAllowedNamespaces(schemaNSMeta, dbNSMeta, consumers)
	if err != nil {
		return false, err
	}
	matchLabels, err := IsMatchByLabels(schemaMeta, consumers)
	if err != nil {
		return false, err
	}
	return matchNamespace && matchLabels, nil
}

func IsInAllowedNamespaces(schemaNSMeta metav1.ObjectMeta, dbNSMeta metav1.ObjectMeta, consumers *dbapi.AllowedConsumers) (bool, error) {
	if consumers.Namespaces == nil || consumers.Namespaces.From == nil {
		return false, nil
	}

	if *consumers.Namespaces.From == dbapi.NamespacesFromAll {
		return true, nil
	}
	if *consumers.Namespaces.From == dbapi.NamespacesFromSame {
		return schemaNSMeta.GetName() == dbNSMeta.GetName(), nil
	}
	if *consumers.Namespaces.From == dbapi.NamespacesFromSelector {
		if consumers.Namespaces.Selector == nil {
			// this says, Select namespace from the Selector, but the Namespace.Selector field is nil. So, no way to select namespace here.
			return false, nil
		}
		ret, err := selectorMatches(consumers.Namespaces.Selector, schemaNSMeta.GetLabels())
		if err != nil {
			return false, err
		}
		return ret, nil
	}
	return false, nil
}

func IsMatchByLabels(schemaMeta metav1.ObjectMeta, consumers *dbapi.AllowedConsumers) (bool, error) {
	if consumers.Selector != nil {
		ret, err := selectorMatches(consumers.Selector, schemaMeta.Labels)
		if err != nil {
			return false, err
		}
		return ret, nil
	}
	// if Selector is not given, all the Schemas are allowed of the selected namespace
	return true, nil
}

func selectorMatches(ls *metav1.LabelSelector, srcLabels map[string]string) (bool, error) {
	selector, err := metav1.LabelSelectorAsSelector(ls)
	if err != nil {
		klog.Infoln("invalid selector: ", ls)
		return false, err
	}
	return selector.Matches(labels.Set(srcLabels)), nil
}
