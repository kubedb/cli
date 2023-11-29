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

package apis

const (
	ImageRegistry = "imageRegistry"
	ImageRepo     = "imageRepo"
	ImageTag      = "imageTag"

	InvokerKind = "invokerKind"
	InvokerName = "invokerName"

	Snapshot = "snapshot"

	Namespace      = "namespace"
	BackupSession  = "backupSession"
	RestoreSession = "restoreSession"

	// EnableCache is false when TmpDir.DisableCaching is true in backupConfig/restoreSession
	// default is true
	EnableCache    = "enableCache"
	InterimDataDir = "interimDataDir"

	LicenseApiService = "licenseApiService"
)
