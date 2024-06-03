/*
Copyright AppsCode Inc. and Contributors.

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

package v1

import (
	"kubeops.dev/petset/pkg/features"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// log is for logging in this package.
var petsetlog = logf.Log.WithName("petset-resource")

// SetupWebhookWithManager will setup the manager to manage the webhooks
func (r *PetSet) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

//+kubebuilder:webhook:path=/mutate-apps-k8s-appscode-com-v1-petset,mutating=true,failurePolicy=fail,sideEffects=None,groups=apps.k8s.appscode.com,resources=petsets,verbs=create;update,versions=v1,name=mpetset.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &PetSet{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (obj *PetSet) Default() {
	if len(obj.Spec.PodManagementPolicy) == 0 {
		obj.Spec.PodManagementPolicy = appsv1.OrderedReadyPodManagement
	}

	if obj.Spec.UpdateStrategy.Type == "" {
		obj.Spec.UpdateStrategy.Type = appsv1.RollingUpdateStatefulSetStrategyType

		if obj.Spec.UpdateStrategy.RollingUpdate == nil {
			// UpdateStrategy.RollingUpdate will take default values below.
			obj.Spec.UpdateStrategy.RollingUpdate = &appsv1.RollingUpdateStatefulSetStrategy{}
		}
	}

	if obj.Spec.UpdateStrategy.Type == appsv1.RollingUpdateStatefulSetStrategyType &&
		obj.Spec.UpdateStrategy.RollingUpdate != nil {

		if obj.Spec.UpdateStrategy.RollingUpdate.Partition == nil {
			obj.Spec.UpdateStrategy.RollingUpdate.Partition = ptr.To[int32](0)
		}
		if features.DefaultFeatureGate.Enabled(features.MaxUnavailablePetSet) {
			if obj.Spec.UpdateStrategy.RollingUpdate.MaxUnavailable == nil {
				obj.Spec.UpdateStrategy.RollingUpdate.MaxUnavailable = ptr.To(intstr.FromInt32(1))
			}
		}
	}

	if features.DefaultFeatureGate.Enabled(features.PetSetAutoDeletePVC) {
		if obj.Spec.PersistentVolumeClaimRetentionPolicy == nil {
			obj.Spec.PersistentVolumeClaimRetentionPolicy = &appsv1.StatefulSetPersistentVolumeClaimRetentionPolicy{}
		}
		if len(obj.Spec.PersistentVolumeClaimRetentionPolicy.WhenDeleted) == 0 {
			obj.Spec.PersistentVolumeClaimRetentionPolicy.WhenDeleted = appsv1.RetainPersistentVolumeClaimRetentionPolicyType
		}
		if len(obj.Spec.PersistentVolumeClaimRetentionPolicy.WhenScaled) == 0 {
			obj.Spec.PersistentVolumeClaimRetentionPolicy.WhenScaled = appsv1.RetainPersistentVolumeClaimRetentionPolicyType
		}
	}

	if obj.Spec.Replicas == nil {
		obj.Spec.Replicas = new(int32)
		*obj.Spec.Replicas = 1
	}
	if obj.Spec.RevisionHistoryLimit == nil {
		obj.Spec.RevisionHistoryLimit = new(int32)
		*obj.Spec.RevisionHistoryLimit = 10
	}
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-apps-k8s-appscode-com-v1-petset,mutating=false,failurePolicy=fail,sideEffects=None,groups=apps.k8s.appscode.com,resources=petsets,verbs=create;update,versions=v1,name=vpetset.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &PetSet{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *PetSet) ValidateCreate() (admission.Warnings, error) {
	petsetlog.Info("validate create", "name", r.Name)

	// TODO(user): fill in your validation logic upon object creation.
	return nil, nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *PetSet) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	petsetlog.Info("validate update", "name", r.Name)

	// TODO(user): fill in your validation logic upon object update.
	return nil, nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *PetSet) ValidateDelete() (admission.Warnings, error) {
	petsetlog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil, nil
}
