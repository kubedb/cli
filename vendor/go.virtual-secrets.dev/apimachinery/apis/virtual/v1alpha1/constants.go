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

const (
	AWSAccessKeyID     = "AWS_ACCESS_KEY_ID"
	AWSSecretAccessKey = "AWS_SECRET_ACCESS_KEY"
	AWSSessionToken    = "AWS_SESSION_TOKEN"
	AWSRegion          = "AWS_REGION"

	AnnotationDataLenKey           = "virtual-secrets.dev/data-len"
	AnnotationCreationTimestampKey = "virtual-secrets.dev/creation-timestamp"

	AzureTenantID     = "AZURE_TENANT_ID"
	AzureClientID     = "AZURE_CLIENT_ID"
	AzureClientSecret = "AZURE_CLIENT_SECRET"

	AccessModeWorkloadIdentity = "WorkloadIdentity"
	AccessModeServicePrincipal = "ServicePrincipal"

	GCPServiceAccountJSONKey = "credentials.json"
)
