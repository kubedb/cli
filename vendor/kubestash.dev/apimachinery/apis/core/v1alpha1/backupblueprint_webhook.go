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
	"context"
	"fmt"
	core "k8s.io/api/core/v1"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	kmapi "kmodules.xyz/client-go/api/v1"
	"kubestash.dev/apimachinery/apis"
	storageapi "kubestash.dev/apimachinery/apis/storage/v1alpha1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// log is for logging in this package.
var backupblueprintlog = logf.Log.WithName("backupblueprint-resource")

func (r *BackupBlueprint) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

//+kubebuilder:webhook:path=/mutate-core-kubestash-com-v1alpha1-backupblueprint,mutating=true,failurePolicy=fail,sideEffects=None,groups=core.kubestash.com,resources=backupblueprints,verbs=create;update,versions=v1alpha1,name=mbackupblueprint.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &BackupBlueprint{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *BackupBlueprint) Default() {
	backupblueprintlog.Info("default", "name", r.Name)

	if r.Spec.UsagePolicy == nil {
		r.setDefaultUsagePolicy()
	}
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-core-kubestash-com-v1alpha1-backupblueprint,mutating=false,failurePolicy=fail,sideEffects=None,groups=core.kubestash.com,resources=backupblueprints,verbs=create;update,versions=v1alpha1,name=vbackupblueprint.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &BackupBlueprint{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *BackupBlueprint) ValidateCreate() (admission.Warnings, error) {
	backupblueprintlog.Info("validate create", "name", r.Name)

	if err := r.validateUsagePolicy(); err != nil {
		return nil, err
	}

	return nil, r.validateBackendsAgainstUsagePolicy(context.Background(), apis.GetRuntimeClient())
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *BackupBlueprint) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	backupblueprintlog.Info("validate update", "name", r.Name)

	if err := r.validateUsagePolicy(); err != nil {
		return nil, err
	}

	return nil, r.validateBackendsAgainstUsagePolicy(context.Background(), apis.GetRuntimeClient())
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *BackupBlueprint) ValidateDelete() (admission.Warnings, error) {
	backupblueprintlog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil, nil
}

func (r *BackupBlueprint) setDefaultUsagePolicy() {
	fromSameNamespace := apis.NamespacesFromSame
	r.Spec.UsagePolicy = &apis.UsagePolicy{
		AllowedNamespaces: apis.AllowedNamespaces{
			From: &fromSameNamespace,
		},
	}
}

func (r *BackupBlueprint) validateBackendsAgainstUsagePolicy(ctx context.Context, c client.Client) error {
	if r.Spec.BackupConfigurationTemplate == nil {
		return fmt.Errorf("backupConfigurationTemplate can not be empty")
	}

	for _, backend := range r.Spec.BackupConfigurationTemplate.Backends {
		bs, err := r.getBackupStorage(ctx, c, backend.StorageRef)
		if err != nil {
			if kerr.IsNotFound(err) {
				continue
			}
			return err
		}

		ns := &core.Namespace{ObjectMeta: v1.ObjectMeta{Name: r.Namespace}}
		if err := c.Get(ctx, client.ObjectKeyFromObject(ns), ns); err != nil {
			return err
		}

		if !bs.UsageAllowed(ns) {
			return fmt.Errorf("namespace %q is not allowed to refer BackupStorage %s/%s. Please, check the `usagePolicy` of the BackupStorage", r.Namespace, bs.Name, bs.Namespace)
		}
	}
	return nil
}

func (r *BackupBlueprint) getBackupStorage(ctx context.Context, c client.Client, ref *kmapi.ObjectReference) (*storageapi.BackupStorage, error) {
	bs := &storageapi.BackupStorage{
		ObjectMeta: v1.ObjectMeta{
			Name:      ref.Name,
			Namespace: ref.Namespace,
		},
	}

	if bs.Namespace == "" {
		bs.Namespace = r.Namespace
	}

	if err := c.Get(ctx, client.ObjectKeyFromObject(bs), bs); err != nil {
		return nil, err
	}
	return bs, nil
}

func (r *BackupBlueprint) validateUsagePolicy() error {
	if *r.Spec.UsagePolicy.AllowedNamespaces.From == apis.NamespacesFromSelector &&
		r.Spec.UsagePolicy.AllowedNamespaces.Selector == nil {
		return fmt.Errorf("selector cannot be empty for usage policy of type %q", apis.NamespacesFromSelector)
	}
	return nil
}
