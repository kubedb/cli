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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kmapi "kmodules.xyz/client-go/api/v1"
	"kubestash.dev/apimachinery/apis"
	"time"

	storageapi "kubestash.dev/apimachinery/apis/storage/v1alpha1"
	"kubestash.dev/apimachinery/crds"

	"kmodules.xyz/client-go/apiextensions"
	cutil "kmodules.xyz/client-go/conditions"
	"kmodules.xyz/client-go/meta"
	meta_util "kmodules.xyz/client-go/meta"
)

func (_ BackupSession) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(GroupVersion.WithResource(ResourcePluralBackupSession))
}

func (b *BackupSession) IsRunning() bool {
	return b.Status.Phase == BackupSessionRunning
}

func (b *BackupSession) IsCompleted() bool {
	phase := b.Status.Phase

	return phase == BackupSessionSucceeded ||
		phase == BackupSessionFailed ||
		phase == BackupSessionSkipped
}

func (b *BackupSession) CalculatePhase() BackupSessionPhase {
	if cutil.IsConditionFalse(b.Status.Conditions, TypeMetricsPushed) {
		return BackupSessionFailed
	}

	if cutil.IsConditionTrue(b.Status.Conditions, TypeBackupSkipped) {
		return BackupSessionSkipped
	}

	if cutil.IsConditionTrue(b.Status.Conditions, TypeMetricsPushed) &&
		(b.failedToEnsurebackupExecutor() ||
			b.failedToEnsureSnapshots() ||
			b.failedToExecutePreBackupHooks() ||
			b.failedToExecutePostBackupHooks() ||
			b.failedToApplyRetentionPolicy() ||
			b.sessionHistoryCleanupFailed() ||
			b.snapshotCleanupIncomplete()) {
		return BackupSessionFailed
	}

	componentsPhase := b.calculateBackupSessionPhaseFromSnapshots()
	if componentsPhase == BackupSessionPending || b.FinalStepExecuted() {
		return componentsPhase
	}

	return BackupSessionRunning
}

func (b *BackupSession) snapshotCleanupIncomplete() bool {
	return cutil.IsConditionTrue(b.Status.Conditions, TypeSnapshotCleanupIncomplete)
}

func (b *BackupSession) sessionHistoryCleanupFailed() bool {
	return cutil.IsConditionFalse(b.Status.Conditions, TypeSessionHistoryCleaned)
}

func (b *BackupSession) failedToEnsureSnapshots() bool {
	return cutil.IsConditionFalse(b.Status.Conditions, TypeSnapshotsEnsured)
}

func (b *BackupSession) failedToEnsurebackupExecutor() bool {
	return cutil.IsConditionFalse(b.Status.Conditions, TypeBackupExecutorEnsured)
}

func (b *BackupSession) FinalStepExecuted() bool {
	return cutil.HasCondition(b.Status.Conditions, TypeMetricsPushed)
}

func (b *BackupSession) failedToExecutePreBackupHooks() bool {
	return cutil.IsConditionFalse(b.Status.Conditions, TypePreBackupHooksExecutionSucceeded)
}

func (b *BackupSession) failedToExecutePostBackupHooks() bool {
	return cutil.IsConditionFalse(b.Status.Conditions, TypePostBackupHooksExecutionSucceeded)
}

func (b *BackupSession) failedToApplyRetentionPolicy() bool {
	for _, status := range b.Status.RetentionPolicies {
		if status.Phase == RetentionPolicyFailedToApply {
			return true
		}
	}

	return false
}

func (b *BackupSession) calculateBackupSessionPhaseFromSnapshots() BackupSessionPhase {
	status := b.Status.Snapshots
	if len(status) == 0 {
		return BackupSessionPending
	}

	pending := 0
	failed := 0
	succeeded := 0

	for _, s := range status {
		if s.Phase == storageapi.SnapshotFailed {
			failed++
		}
		if s.Phase == storageapi.SnapshotPending {
			pending++
		}
		if s.Phase == storageapi.SnapshotSucceeded {
			succeeded++
		}
	}

	if pending == len(status) {
		return BackupSessionPending
	}

	if succeeded+failed != len(status) {
		return BackupSessionRunning
	}

	if failed > 0 {
		return BackupSessionFailed
	}

	return BackupSessionSucceeded
}

func GenerateBackupSessionName(invokerName, sessionName string) string {
	return meta.ValidNameWithPrefixNSuffix(invokerName, sessionName, fmt.Sprintf("%d", time.Now().Unix()))
}

func (b *BackupSession) OffshootLabels() map[string]string {
	newLabels := make(map[string]string)
	newLabels[meta_util.ManagedByLabelKey] = apis.KubeStashKey
	newLabels[apis.KubeStashInvokerName] = b.Name
	newLabels[apis.KubeStashInvokerNamespace] = b.Namespace
	newLabels[apis.KubeStashSessionName] = b.Spec.Session

	return apis.UpsertLabels(b.Labels, newLabels)
}

func (b *BackupSession) GetSummary(targetRef *kmapi.TypedObjectReference) *Summary {
	errMsg := b.getFailureMessage()
	phase := BackupSessionSucceeded
	if errMsg != "" {
		phase = BackupSessionFailed
	}

	return &Summary{
		Name:      b.Name,
		Namespace: b.Namespace,

		Invoker: &kmapi.TypedObjectReference{
			APIGroup:  GroupVersion.Group,
			Kind:      b.Spec.Invoker.Kind,
			Name:      b.Spec.Invoker.Name,
			Namespace: b.Namespace,
		},

		Target: targetRef,

		Status: TargetStatus{
			Phase:    string(phase),
			Duration: b.Status.Duration,
			Error:    errMsg,
		},
	}
}

func (b *BackupSession) getFailureMessage() string {
	failureFound, reason := b.checkFailureInConditions()
	if failureFound {
		return reason
	}
	failureFound, reason = b.checkFailureInSnapshots()
	if failureFound {
		return reason
	}
	failureFound, reason = b.checkFailureInRetentionPolicy()
	if failureFound {
		return reason
	}

	return ""
}

func (b *BackupSession) checkFailureInConditions() (bool, string) {
	for _, condition := range b.Status.Conditions {
		if condition.Status == metav1.ConditionFalse {
			return true, condition.Message
		}
	}

	return false, ""
}

func (b *BackupSession) checkFailureInSnapshots() (bool, string) {
	for _, snapStatus := range b.Status.Snapshots {
		if snapStatus.Phase == storageapi.SnapshotFailed {
			return true, "one or more snapshots are failed"
		}
	}
	return false, ""
}

func (b *BackupSession) checkFailureInRetentionPolicy() (bool, string) {
	for _, retention := range b.Status.RetentionPolicies {
		if retention.Phase == RetentionPolicyFailedToApply {
			return true, "one or more retention policies are failed to apply"
		}
	}
	return false, ""
}

func (b *BackupSession) GetRemainingTimeoutDuration() (*metav1.Duration, error) {
	if b.Spec.BackupTimeout == nil || b.Status.BackupDeadline == nil {
		return nil, nil
	}
	currentTime := metav1.Now()
	if b.Status.BackupDeadline.Before(&currentTime) {
		return nil, fmt.Errorf("deadline exceeded")
	}
	return &metav1.Duration{Duration: b.Status.BackupDeadline.Sub(currentTime.Time)}, nil
}
