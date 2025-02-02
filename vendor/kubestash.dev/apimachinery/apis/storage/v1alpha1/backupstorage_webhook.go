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
	"k8s.io/apimachinery/pkg/runtime"
	"kubestash.dev/apimachinery/apis"
	"reflect"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	"strings"
)

// log is for logging in this package.
var backupstoragelog = logf.Log.WithName("backupstorage-resource")

func (r *BackupStorage) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

//+kubebuilder:webhook:path=/mutate-storage-kubestash-com-v1alpha1-backupstorage,mutating=true,failurePolicy=fail,sideEffects=None,groups=storage.kubestash.com,resources=backupstorages,verbs=create;update,versions=v1alpha1,name=mbackupstorage.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &BackupStorage{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *BackupStorage) Default() {
	backupstoragelog.Info("default", "name", r.Name)

	if r.Spec.UsagePolicy == nil {
		r.setDefaultUsagePolicy()
	}
	r.removeTrailingSlash()
}

func (r *BackupStorage) removeTrailingSlash() {
	if r.Spec.Storage.S3 != nil {
		r.Spec.Storage.S3.Bucket = strings.TrimSuffix(r.Spec.Storage.S3.Bucket, "/")
		r.Spec.Storage.S3.Endpoint = strings.TrimSuffix(r.Spec.Storage.S3.Endpoint, "/")
		r.Spec.Storage.S3.Prefix = strings.TrimSuffix(r.Spec.Storage.S3.Prefix, "/")
	}
	if r.Spec.Storage.GCS != nil {
		r.Spec.Storage.GCS.Bucket = strings.TrimSuffix(r.Spec.Storage.GCS.Bucket, "/")
		r.Spec.Storage.GCS.Prefix = strings.TrimSuffix(r.Spec.Storage.GCS.Prefix, "/")
	}
	if r.Spec.Storage.Azure != nil {
		r.Spec.Storage.Azure.Prefix = strings.TrimSuffix(r.Spec.Storage.Azure.Prefix, "/")
	}
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-storage-kubestash-com-v1alpha1-backupstorage,mutating=false,failurePolicy=fail,sideEffects=None,groups=storage.kubestash.com,resources=backupstorages,verbs=create;update,versions=v1alpha1,name=vbackupstorage.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &BackupStorage{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *BackupStorage) ValidateCreate() (admission.Warnings, error) {
	backupstoragelog.Info("validate create", "name", r.Name)

	c := apis.GetRuntimeClient()

	if r.Spec.Default {
		if err := r.validateSingleDefaultBackupStorageInSameNamespace(context.Background(), c); err != nil {
			return nil, err
		}
	}

	if err := r.validateUsagePolicy(); err != nil {
		return nil, err
	}

	return nil, r.validateUniqueDirectory(context.Background(), c)
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *BackupStorage) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	backupstoragelog.Info("validate update", "name", r.Name)

	c := apis.GetRuntimeClient()

	if r.Spec.Default {
		if err := r.validateSingleDefaultBackupStorageInSameNamespace(context.Background(), c); err != nil {
			return nil, err
		}
	}

	if err := r.validateUsagePolicy(); err != nil {
		return nil, err
	}

	if err := r.validateUpdateStorage(old.(*BackupStorage)); err != nil {
		return nil, err
	}

	return nil, r.validateUniqueDirectory(context.Background(), c)
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *BackupStorage) ValidateDelete() (admission.Warnings, error) {
	backupstoragelog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil, nil
}

func (r *BackupStorage) setDefaultUsagePolicy() {
	fromSameNamespace := apis.NamespacesFromSame
	r.Spec.UsagePolicy = &apis.UsagePolicy{
		AllowedNamespaces: apis.AllowedNamespaces{
			From: &fromSameNamespace,
		},
	}
}

func (r *BackupStorage) validateSingleDefaultBackupStorageInSameNamespace(ctx context.Context, c client.Client) error {
	bsList := BackupStorageList{}
	if err := c.List(ctx, &bsList, client.InNamespace(r.Namespace)); err != nil {
		return err
	}

	for _, bs := range bsList.Items {
		if !r.isSameBackupStorage(bs) &&
			bs.Spec.Default {
			return fmt.Errorf("multiple default BackupStorages are not allowed within the same namespace")
		}
	}

	return nil
}

func (r *BackupStorage) validateUsagePolicy() error {
	if *r.Spec.UsagePolicy.AllowedNamespaces.From == apis.NamespacesFromSelector &&
		r.Spec.UsagePolicy.AllowedNamespaces.Selector == nil {
		return fmt.Errorf("selector cannot be empty for usage policy of type %q", apis.NamespacesFromSelector)
	}
	return nil
}

func (r *BackupStorage) isSameBackupStorage(bs BackupStorage) bool {
	if r.Namespace == bs.Namespace &&
		r.Name == bs.Name {
		return true
	}
	return false
}

func (r *BackupStorage) validateUpdateStorage(old *BackupStorage) error {
	if !reflect.DeepEqual(old.Spec.Storage, r.Spec.Storage) &&
		len(r.Status.Repositories) != 0 {
		return fmt.Errorf("BackupStorage is currently in use and cannot be modified")
	}
	return nil
}

func (r *BackupStorage) validateUniqueDirectory(ctx context.Context, c client.Client) error {
	bsList := BackupStorageList{}
	if err := c.List(ctx, &bsList); err != nil {
		return err
	}

	for _, bs := range bsList.Items {
		if !r.isSameBackupStorage(bs) &&
			r.isPointToSameDir(bs) {
			return fmt.Errorf("no two BackupStorage should point to the same directory of the same bucket")
		}
	}

	return nil
}

func (r *BackupStorage) isPointToSameDir(bs BackupStorage) bool {
	if r.Spec.Storage.Provider != bs.Spec.Storage.Provider {
		return false
	}

	switch r.Spec.Storage.Provider {
	case ProviderS3:
		if r.Spec.Storage.S3.Bucket == bs.Spec.Storage.S3.Bucket &&
			r.Spec.Storage.S3.Region == bs.Spec.Storage.S3.Region &&
			r.Spec.Storage.S3.Prefix == bs.Spec.Storage.S3.Prefix {
			return true
		}
		return false
	case ProviderGCS:
		if r.Spec.Storage.GCS.Bucket == bs.Spec.Storage.GCS.Bucket &&
			r.Spec.Storage.GCS.Prefix == bs.Spec.Storage.GCS.Prefix {
			return true
		}
		return false
	case ProviderAzure:
		if r.Spec.Storage.Azure.StorageAccount == bs.Spec.Storage.Azure.StorageAccount &&
			r.Spec.Storage.Azure.Container == bs.Spec.Storage.Azure.Container &&
			r.Spec.Storage.Azure.Prefix == bs.Spec.Storage.Azure.Prefix {
			return true
		}
		return false
	//case ProviderB2:
	//	if r.Spec.Storage.B2.Bucket == bs.Spec.Storage.B2.Bucket &&
	//		r.Spec.Storage.B2.Prefix == bs.Spec.Storage.B2.Prefix {
	//		return true
	//	}
	//	return false
	//case ProviderSwift:
	//	// TODO: check for account
	//	if r.Spec.Storage.Swift.Container == bs.Spec.Storage.Swift.Container &&
	//		r.Spec.Storage.Swift.Prefix == bs.Spec.Storage.Swift.Prefix {
	//		return true
	//	}
	//	return false
	default:
		return false
	}
}
