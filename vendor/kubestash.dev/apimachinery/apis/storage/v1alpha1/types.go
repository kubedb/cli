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
	ofst "kmodules.xyz/offshoot-api/api/v1"
)

// DeletionPolicy specifies what to do if a resource is deleted
// +kubebuilder:validation:Enum=Delete;WipeOut
type DeletionPolicy string

const (
	DeletionPolicyDelete  DeletionPolicy = "Delete"
	DeletionPolicyWipeOut DeletionPolicy = "WipeOut"
)

// +kubebuilder:validation:Enum=Delete;WipeOut;Retain
type BackupConfigDeletionPolicy string

const (
	BackupConfigDeletionPolicyDelete  BackupConfigDeletionPolicy = "Delete"
	BackupConfigDeletionPolicyWipeOut BackupConfigDeletionPolicy = "WipeOut"
	BackupConfigDeletionPolicyRetain  BackupConfigDeletionPolicy = "Retain"
)

type StorageProvider string

const (
	ProviderLocal StorageProvider = "local"
	ProviderS3    StorageProvider = "s3"
	ProviderGCS   StorageProvider = "gcs"
	ProviderAzure StorageProvider = "azure"
	//ProviderSwift StorageProvider = "swift"
	//ProviderB2    StorageProvider = "b2"
	//ProviderRest  StorageProvider = "rest"
)

type Backend struct {
	// Provider specifies the provider of the storage
	Provider StorageProvider `json:"provider,omitempty"`

	// Local specifies the storage information for local provider
	// +optional
	Local *LocalSpec `json:"local,omitempty"`

	// S3 specifies the storage information for AWS S3 and S3 compatible storage.
	// +optional
	S3 *S3Spec `json:"s3,omitempty"`

	// GCS specifies the storage information for GCS bucket
	// +optional
	GCS *GCSSpec `json:"gcs,omitempty"`

	// Azure specifies the storage information for Azure Blob container
	// +optional
	Azure *AzureSpec `json:"azure,omitempty"`

	/*
		// Swift specifies the storage information for Swift container
		// +optional
		Swift *SwiftSpec `json:"swift,omitempty"`

		// B2 specifies the storage information for B2 bucket
		// +optional
		B2 *B2Spec `json:"b2,omitempty"`

		// Rest specifies the storage information for rest storage server
		// +optional
		Rest *RestServerSpec `json:"rest,omitempty"`
	*/
}

type LocalSpec struct {
	// Represents the source of a volume to mount. Only one of its members may be specified.
	// Make sure the volume exist before using the volume as backend.
	ofst.VolumeSource `json:",inline"`

	// MountPath specifies the directory where this volume will be mounted
	MountPath string `json:"mountPath,omitempty"`

	// Path within the volume from which the container's volume should be mounted.
	// Defaults to "" (volume's root).
	// +optional
	SubPath string `json:"subPath,omitempty"`
}

type S3Spec struct {
	// Endpoint specifies the URL of the S3 or S3 compatible storage bucket.
	Endpoint string `json:"endpoint,omitempty"`

	// Bucket specifies the name of the bucket that will be used as storage backend.
	Bucket string `json:"bucket,omitempty"`

	// Prefix specifies a directory inside the bucket/container where the data for this backend will be stored.
	Prefix string `json:"prefix,omitempty"`

	// Region specifies the region where the bucket is located
	// +optional
	Region string `json:"region,omitempty"`

	// SecretName specifies the name of the Secret that contains the access credential for this storage.
	// +optional
	SecretName string `json:"secretName,omitempty"`

	// InsecureTLS controls whether a client should skip TLS certificate verification.
	// Setting this field to true disables verification, which might be necessary in cases
	// where the server uses self-signed certificates or certificates from an untrusted CA.
	// Use this option with caution, as it can expose the client to man-in-the-middle attacks
	// and other security risks. Only use it when absolutely necessary.
	// +optional
	InsecureTLS bool `json:"insecureTLS,omitempty"`
}

type GCSSpec struct {
	// Bucket specifies the name of the bucket that will be used as storage backend.
	Bucket string `json:"bucket,omitempty"`

	// Prefix specifies a directory inside the bucket/container where the data for this backend will be stored.
	Prefix string `json:"prefix,omitempty"`

	// MaxConnections specifies the maximum number of concurrent connections to use to upload/download data to this backend.
	// +optional
	MaxConnections int64 `json:"maxConnections,omitempty"`

	// SecretName specifies the name of the Secret that contains the access credential for this storage.
	// +optional
	SecretName string `json:"secretName,omitempty"`
}

type AzureSpec struct {
	// StorageAccount specifies the name of the Azure Storage Account
	StorageAccount string `json:"storageAccount,omitempty"`

	// Container specifies the name of the Azure Blob container that will be used as storage backend.
	Container string `json:"container,omitempty"`

	// Prefix specifies a directory inside the bucket/container where the data for this backend will be stored.
	Prefix string `json:"prefix,omitempty"`

	// MaxConnections specifies the maximum number of concurrent connections to use to upload/download data to this backend.
	// +optional
	MaxConnections int64 `json:"maxConnections,omitempty"`

	// SecretName specifies the name of the Secret that contains the access credential for this storage.
	// +optional
	SecretName string `json:"secretName,omitempty"`
}

/*
type SwiftSpec struct {
	// Container specifies the name of the Swift container that will be used as storage backend.
	Container string `json:"container,omitempty"`

	// Prefix specifies a directory inside the bucket/container where the data for this backend will be stored.
	Prefix string `json:"prefix,omitempty"`

	// Secret specifies the name of the Secret that contains the access credential for this storage.
	// +optional
	SecretName string `json:"secretName,omitempty"`
}

type B2Spec struct {
	// Bucket specifies the name of the bucket that will be used as storage backend.
	Bucket string `json:"bucket,omitempty"`

	// Prefix specifies a directory inside the bucket/container where the data for this backend will be stored.
	Prefix string `json:"prefix,omitempty"`

	// MaxConnections specifies the maximum number of concurrent connections to use to upload/download data to this backend.
	// +optional
	MaxConnections int64 `json:"maxConnections,omitempty"`

	// Secret specifies the name of the Secret that contains the access credential for this storage.
	// +optional
	SecretName string `json:"secretName,omitempty"`
}

type RestServerSpec struct {
	// URL specifies the URL of the REST storage server
	URL string `json:"url,omitempty"`

	// Secret specifies the name of the Secret that contains the access credential for this storage.
	// +optional
	SecretName string `json:"secretName,omitempty"`
}
*/
