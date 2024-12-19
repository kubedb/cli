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
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// log is for logging in this package.
var backupverifierlog = logf.Log.WithName("backupverifier-resource")

func (v *BackupVerifier) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(v).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-core-kubestash-com-v1alpha1-backupverifier,mutating=false,failurePolicy=fail,sideEffects=None,groups=core.kubestash.com,resources=backupverifiers,verbs=create;update,versions=v1alpha1,name=vbackupverifier.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &BackupVerifier{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (v *BackupVerifier) ValidateCreate() (admission.Warnings, error) {
	backupverifierlog.Info("validate create", "name", v.Name)

	if err := v.validateVerifier(); err != nil {
		return nil, err
	}

	return nil, nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (v *BackupVerifier) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	backupverifierlog.Info("validate update", "name", v.Name)

	if err := v.validateVerifier(); err != nil {
		return nil, err
	}

	return nil, nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (v *BackupVerifier) ValidateDelete() (admission.Warnings, error) {
	backupverifierlog.Info("validate delete", "name", v.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil, nil
}

func (v *BackupVerifier) validateVerifier() error {
	if v.Spec.RestoreOption == nil {
		return fmt.Errorf("restoreOption for backupVerifier %s/%s cannot be empty", v.Namespace, v.Name)
	}

	if v.Spec.RestoreOption.AddonInfo == nil {
		return fmt.Errorf("addonInfo in restoreOption for backupVerifier %s/%s cannot be empty", v.Namespace, v.Name)
	}

	if v.Spec.Scheduler != nil {
		return fmt.Errorf("scheduler for backupVerifier %s/%s cannot be empty", v.Namespace, v.Name)
	}

	if v.Spec.Type == "" {
		return fmt.Errorf("type of backupVerifier %s/%s cannot be empty", v.Namespace, v.Name)
	}

	if v.Spec.Type == QueryVerificationType {
		if v.Spec.Query == nil {
			return fmt.Errorf("query in backupVerifier %s/%s cannot be empty", v.Namespace, v.Name)
		}
		if v.Spec.Function == "" {
			return fmt.Errorf("function in backupVerifier %s/%s cannot be empty", v.Namespace, v.Name)
		}
	}

	if v.Spec.Type == ScriptVerificationType {
		if v.Spec.Script == nil {
			return fmt.Errorf("script in backupVerifier %s/%s cannot be empty", v.Namespace, v.Name)
		}

		if v.Spec.Script.Location == "" {
			return fmt.Errorf("script location in backupVerifier %s/%s cannot be empty", v.Namespace, v.Name)
		}

		if v.Spec.Function == "" {
			return fmt.Errorf("function in backupVerifier %s/%s cannot be empty", v.Namespace, v.Name)
		}
	}

	return nil
}
