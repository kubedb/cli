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

	"kubestash.dev/apimachinery/apis"
	"kubestash.dev/apimachinery/crds"

	"gomodules.xyz/restic"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/klog/v2"
	kmapi "kmodules.xyz/client-go/api/v1"
	"kmodules.xyz/client-go/apiextensions"
	cutil "kmodules.xyz/client-go/conditions"
	"kmodules.xyz/client-go/meta"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (BackupStorage) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(GroupVersion.WithResource(ResourcePluralBackupStorage))
}

func (b *BackupStorage) CalculatePhase() BackupStoragePhase {
	if cutil.IsConditionTrue(b.Status.Conditions, TypeBackendInitialized) {
		if !cutil.HasCondition(b.Status.Conditions, TypeBackendSecretFound) {
			return BackupStorageReady
		}
		if cutil.IsConditionTrue(b.Status.Conditions, TypeBackendSecretFound) {
			return BackupStorageReady
		}
	}
	return BackupStorageNotReady
}

func (b *BackupStorage) UsageAllowed(srcNamespace *core.Namespace) bool {
	allowedNamespaces := b.Spec.UsagePolicy.AllowedNamespaces

	if allowedNamespaces.From == nil {
		return false
	}

	if *allowedNamespaces.From == apis.NamespacesFromAll {
		return true
	}

	if *allowedNamespaces.From == apis.NamespacesFromSame {
		return b.Namespace == srcNamespace.Name
	}

	return selectorMatches(allowedNamespaces.Selector, srcNamespace.Labels)
}

func selectorMatches(ls *metav1.LabelSelector, srcLabels map[string]string) bool {
	selector, err := metav1.LabelSelectorAsSelector(ls)
	if err != nil {
		klog.Infoln("invalid label selector: ", ls)
		return false
	}
	return selector.Matches(labels.Set(srcLabels))
}

func (b *BackupStorage) OffshootLabels() map[string]string {
	newLabels := make(map[string]string)
	newLabels[meta.ManagedByLabelKey] = apis.KubeStashKey
	newLabels[apis.KubeStashInvokerKind] = ResourceKindBackupStorage
	newLabels[apis.KubeStashInvokerName] = b.Name
	newLabels[apis.KubeStashInvokerNamespace] = b.Namespace
	return apis.UpsertLabels(b.Labels, newLabels)
}

func (b *BackupStorage) LocalProvider() bool {
	return b.Spec.Storage.Provider == ProviderLocal
}

func (b *BackupStorage) LocalNetworkVolume() bool {
	if b.Spec.Storage.Provider == ProviderLocal &&
		b.Spec.Storage.Local.NFS != nil {
		return true
	}
	return false
}

// NewBackupStorageResolver creates a StorageConfigResolver that resolves storage configuration
// from a BackupStorage custom resource. This is the default resolver for the kubestash project.
func NewBackupStorageResolver(kbClient client.Client, bsRef *kmapi.ObjectReference) restic.StorageConfigResolver {
	return func(backend *restic.Backend) error {
		bs := &BackupStorage{
			ObjectMeta: metav1.ObjectMeta{
				Name:      bsRef.Name,
				Namespace: bsRef.Namespace,
			},
		}

		if err := kbClient.Get(context.Background(), client.ObjectKeyFromObject(bs), bs); err != nil {
			return fmt.Errorf("failed to get BackupStorage %s/%s: %w", bsRef.Namespace, bsRef.Name, err)
		}
		var storageSecretName string
		switch {
		case bs.Spec.Storage.S3 != nil:
			s3 := bs.Spec.Storage.S3
			storageSecretName = s3.SecretName
			backend.StorageConfig = &restic.StorageConfig{
				Provider:       string(ProviderS3),
				Bucket:         s3.Bucket,
				Endpoint:       s3.Endpoint,
				Region:         s3.Region,
				Prefix:         s3.Prefix,
				InsecureTLS:    s3.InsecureTLS,
				MaxConnections: s3.MaxConnections,
			}
		case bs.Spec.Storage.GCS != nil:
			gcs := bs.Spec.Storage.GCS
			storageSecretName = gcs.SecretName
			backend.StorageConfig = &restic.StorageConfig{
				Provider:       string(ProviderGCS),
				Bucket:         gcs.Bucket,
				Prefix:         gcs.Prefix,
				MaxConnections: gcs.MaxConnections,
			}
		case bs.Spec.Storage.Azure != nil:
			azure := bs.Spec.Storage.Azure
			storageSecretName = azure.SecretName
			backend.StorageConfig = &restic.StorageConfig{
				Provider:            string(ProviderAzure),
				Bucket:              azure.Container,
				Prefix:              azure.Prefix,
				AzureStorageAccount: azure.StorageAccount,
				MaxConnections:      azure.MaxConnections,
			}
		case bs.Spec.Storage.Local != nil:
			local := bs.Spec.Storage.Local
			backend.StorageConfig = &restic.StorageConfig{
				Provider:       string(ProviderLocal),
				Bucket:         local.MountPath,
				Prefix:         local.SubPath,
				MaxConnections: local.MaxConnections,
			}
			if backend.MountPath != "" {
				backend.Bucket = backend.MountPath
			}
		default:
			return fmt.Errorf("no storage backend configured in BackupStorage %s/%s", bsRef.Namespace, bsRef.Name)
		}

		if storageSecretName != "" {
			secret := &core.Secret{}
			if err := kbClient.Get(context.Background(), client.ObjectKey{Name: storageSecretName, Namespace: bsRef.Namespace}, secret); err != nil {
				return fmt.Errorf("failed to get storage Secret %s/%s: %w", bsRef.Namespace, storageSecretName, err)
			}
			backend.StorageSecret = secret
		}
		return nil
	}
}
