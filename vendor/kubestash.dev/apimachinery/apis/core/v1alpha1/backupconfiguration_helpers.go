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
	"kubestash.dev/apimachinery/apis"
	"kubestash.dev/apimachinery/crds"

	kmapi "kmodules.xyz/client-go/api/v1"
	"kmodules.xyz/client-go/apiextensions"
	cutil "kmodules.xyz/client-go/conditions"
)

func (BackupConfiguration) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(GroupVersion.WithResource(ResourcePluralBackupConfiguration))
}

func (b *BackupConfiguration) CalculatePhase() BackupInvokerPhase {
	if cutil.IsConditionFalse(b.Status.Conditions, TypeValidationPassed) {
		return BackupInvokerInvalid
	}

	if b.isReady() {
		return BackupInvokerReady
	}

	return BackupInvokerNotReady
}

func (b *BackupConfiguration) isReady() bool {
	if b.Status.TargetFound == nil || !*b.Status.TargetFound {
		return false
	}

	if !b.backendsReady() {
		return false
	}

	if !b.sessionsReady() {
		return false
	}

	return true
}

func (b *BackupConfiguration) sessionsReady() bool {
	if len(b.Status.Sessions) != len(b.Spec.Sessions) {
		return false
	}

	for _, status := range b.Status.Sessions {
		if !cutil.IsConditionTrue(status.Conditions, TypeSchedulerEnsured) {
			return false
		}
	}

	return true
}

func (b *BackupConfiguration) backendsReady() bool {
	if len(b.Status.Backends) != len(b.Spec.Backends) {
		return false
	}

	for _, backend := range b.Status.Backends {
		if !*backend.Ready {
			return false
		}
	}

	return true
}

func (b *BackupConfiguration) GetStorageRef(backend string) *kmapi.ObjectReference {
	for _, b := range b.Spec.Backends {
		if b.Name == backend {
			return b.StorageRef
		}
	}
	return nil
}

func (b *BackupConfiguration) GetTargetRef() *kmapi.TypedObjectReference {
	if b.Spec.Target == nil {
		return &kmapi.TypedObjectReference{
			APIGroup: "na",
			Kind:     apis.KindEmpty,
			Name:     "na",
		}
	}
	return b.Spec.Target
}
