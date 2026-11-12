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

import (
	"math"
	"time"
)

const ReasonKubeVirtNotInstalled = "KubeVirtNotInstalled"

const (
	KubeVirtGroup   = "kubevirt.io"
	KubeVirtVersion = "v1"

	KindVirtualMachine         = "VirtualMachine"
	KindVirtualMachineInstance = "VirtualMachineInstance"

	KubeVirtAddon               = "kubevirt"
	TaskKubeVirtLogical         = "logical-backup"
	TaskKubeVirtManifest        = "manifest-backup"
	TaskKubeVirtCBTVolume       = "volume-backup"
	ParamBackupMode             = "backupMode"
	ParamForceFull              = "forceFull"
	ParamCBTScratchStorage      = "cbtScratchStorageClass"
	ParamCBTScratchSize         = "cbtScratchSize"
	ParamCBTArtifact            = "cbtArtifact"
	ArchiverBackendName         = "archiver-backend"
	ArchiverFullSessionName     = "full-backup"
	ArchiverManifestSessionName = "manifest-backup"

	SuffixArchiver = "archiver"
	SuffixFull     = "full"
	SuffixManifest = "manifest"

	SuffixArchiverJob         = "cbt-archiver"
	SuffixIncrementalSnapshot = "cbt"

	ArchiverJobComponent = "kubestash-cbt-archiver"

	CmdCBTArchiver        = "cbt-archiver"
	ArchiverContainerName = "cbt-archiver"

	DefaultCBTInterval = 10 * time.Minute

	ArchiverJobBackoffLimit           int32 = math.MaxInt32
	ArchiverJobTerminationGracePeriod int64 = 120

	FinalizerCBTChain   = "kubevirtarchiver.kubestash.com/cbt-chain"
	AnnArchiverSpecHash = "kubevirtarchiver.kubestash.com/spec-hash"
	AnnArchiverStopped  = "kubevirtarchiver.kubestash.com/stopped-at"
	LabelKubeVirtVMUID  = "kubevirtarchiver.kubestash.com/vm-uid"
)

const (
	CBTStateEnabled            = "Enabled"
	CBTStatePendingRestart     = "PendingRestart"
	CBTStateInitializing       = "Initializing"
	CBTStateFeatureGateDisable = "IncrementalBackupFeatureGateDisabled"
)

const (
	TaskKubeVirtVolumeRestore   = "volume-restore"
	TaskKubeVirtManifestRestore = "manifest-restore"

	ParamRestoreCheckpoint         = "checkpoint"
	ParamRestoreManifestRepository = "manifestRepository"
	ParamRestoreTargetName         = "targetName"
	ParamRestoreTargetNamespace    = "targetNamespace"

	SuffixVolumeRestoreSession   = "volumes"
	SuffixManifestRestoreSession = "manifest"
)
