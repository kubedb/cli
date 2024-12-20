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
	cutil "kmodules.xyz/client-go/conditions"
	"kmodules.xyz/client-go/meta"
	"kubestash.dev/apimachinery/apis"
	"kubestash.dev/apimachinery/crds"
	"time"

	"kmodules.xyz/client-go/apiextensions"
	meta_util "kmodules.xyz/client-go/meta"
)

func (_ BackupVerificationSession) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(GroupVersion.WithResource(ResourcePluralBackupVerificationSession))
}

func (b *BackupVerificationSession) IsCompleted() bool {
	phase := b.Status.Phase

	return phase == BackupVerificationSessionSucceeded ||
		phase == BackupVerificationSessionFailed ||
		phase == BackupVerificationSessionSkipped
}

func (b *BackupVerificationSession) CalculatePhase() BackupVerificationSessionPhase {
	if cutil.IsConditionFalse(b.Status.Conditions, TypeVerificationSessionHistoryCleaned) {
		return BackupVerificationSessionFailed
	}

	if cutil.IsConditionTrue(b.Status.Conditions, TypeBackupVerificationSkipped) {
		return BackupVerificationSessionSkipped
	}

	if b.sessionHistoryCleanupSucceeded() &&
		(b.failedToRestoreBackup() ||
			b.failedToVerifyBackup()) {
		return BackupVerificationSessionFailed
	}

	if cutil.IsConditionTrue(b.Status.Conditions, TypeVerificationSessionHistoryCleaned) {
		return BackupVerificationSessionSucceeded
	}

	return BackupVerificationSessionRunning
}

func (b *BackupVerificationSession) sessionHistoryCleanupFailed() bool {
	return cutil.IsConditionFalse(b.Status.Conditions, TypeVerificationSessionHistoryCleaned)
}

func (b *BackupVerificationSession) sessionHistoryCleanupSucceeded() bool {
	return cutil.IsConditionTrue(b.Status.Conditions, TypeVerificationSessionHistoryCleaned)
}

func (b *BackupVerificationSession) failedToRestoreBackup() bool {
	return cutil.IsConditionFalse(b.Status.Conditions, TypeRestoreSucceeded)
}

func (b *BackupVerificationSession) failedToVerifyBackup() bool {
	return cutil.IsConditionFalse(b.Status.Conditions, TypeBackupVerified)
}

func GenerateBackupVerificationSessionName(repoName, sessionName string) string {
	return meta.ValidNameWithPrefixNSuffix(repoName, sessionName, fmt.Sprintf("%d", time.Now().Unix()))
}

func (b *BackupVerificationSession) OffshootLabels() map[string]string {
	newLabels := make(map[string]string)
	newLabels[meta_util.ManagedByLabelKey] = apis.KubeStashKey
	newLabels[apis.KubeStashInvokerName] = b.Name
	newLabels[apis.KubeStashInvokerNamespace] = b.Namespace
	newLabels[apis.KubeStashSessionName] = b.Spec.Session
	newLabels[apis.KubeStashRepoName] = b.Spec.Repository

	return apis.UpsertLabels(b.Labels, newLabels)
}

func (b *BackupVerificationSession) SetBackupVerifiedConditionToFalse(err error) {
	newCond := kmapi.Condition{
		Type:    TypeBackupVerified,
		Status:  metav1.ConditionFalse,
		Reason:  ReasonFailedToVerifyBackup,
		Message: fmt.Sprintf("Failed to verify backup. Reason: %q", err.Error()),
	}
	b.Status.Conditions = cutil.SetCondition(b.Status.Conditions, newCond)
}

func (b *BackupVerificationSession) SetBackupVerifiedConditionToTrue() {
	newCond := kmapi.Condition{
		Type:   TypeBackupVerified,
		Status: metav1.ConditionTrue,
		Reason: ReasonSuccessfullyVerifiedBackup,
	}
	b.Status.Conditions = cutil.SetCondition(b.Status.Conditions, newCond)
}
