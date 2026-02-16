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
	KindClusterRole           = "ClusterRole"
	KindRole                  = "Role"
	KindPersistentVolumeClaim = "PersistentVolumeClaim"
	KindJob                   = "Job"
)

const (
	MigratorJobClusterRole = "migrator-job"
	MigratorJobPrefix      = "migrator"
	SidecarContainerName   = "status-reporter"
	MigratorGRPCPort       = 50051
	MigratorPVCSuffix      = "pvc"
	PVCVolumeName          = "migrator-data"
	PVCVolumeMountPath     = "/data"
	ConfigVolName          = "migrator-config"
	ConfigVolMountPath     = "/etc/migrator"
	ConfigFileName         = "config.yaml"
	MigratorConfigSuffix   = "config"
	ConfigPath             = ConfigVolMountPath + "/" + ConfigFileName
)

// Conditions Related Constants
const (
	MigratorJobTriggered = "MigratorJobTriggered"

	DestroySignalSend = "DestroySignalSend"

	// MigrationRunning Migration status conditions
	MigrationRunning       = "MigrationRunning"
	ReasonMigrationRunning = "MigrationInProgress"

	MigrationSucceeded       = "MigrationSucceeded"
	ReasonMigrationSucceeded = "MigrationCompleted"

	MigrationFailed       = "MigrationFailed"
	ReasonMigrationFailed = "MigrationError"
)
