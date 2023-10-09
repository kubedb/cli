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
	GenericKey = "ops.kubedb.com"

	LabelOpsRequestKind = GenericKey + "/kind"
	LabelOpsRequestName = GenericKey + "/name"
)

const (
	Running    = "Running"
	Successful = "Successful"
	Failed     = "Failed"
)

// Database
const (
	DatabaseReady           = "DatabaseReady"
	PauseDatabase           = "PauseDatabase"
	DatabasePauseSucceeded  = "DatabasePauseSucceeded"
	DatabasePauseFailed     = "DatabasePauseFailed"
	ResumeDatabase          = "ResumeDatabase"
	DatabaseResumeSucceeded = "DatabaseResumeSucceeded"
	DatabaseResumeFailed    = "DatabaseResumeFailed"
	UpdateDatabase          = "UpdateDatabase"
)

// Version Update
const (
	VersionUpdate          = "VersionUpdate"
	VersionUpdateStarted   = "VersionUpdateStarted"
	VersionUpdateSucceeded = "VersionUpdateSucceeded"
	VersionUpdateFailed    = "VersionUpdateFailed"
)

// Horizontal
const (
	HorizontalScale          = "HorizontalScale"
	HorizontalScaleUp        = "HorizontalScaleUp"
	HorizontalScaleDown      = "HorizontalScaleDown"
	HorizontalScaleStarted   = "HorizontalScaleStarted"
	HorizontalScaleSucceeded = "HorizontalScaleSucceeded"
	HorizontalScaleFailed    = "HorizontalScaleFailed"
)

// Vertical
const (
	VerticalScale          = "VerticalScale"
	VerticalScaleUp        = "VerticalScaleUp"
	VerticalScaleDown      = "VerticalScaleDown"
	VerticalScaleStarted   = "VerticalScaleStarted"
	VerticalScaleSucceeded = "VerticalScaleSucceeded"
	VerticalScaleFailed    = "VerticalScaleFailed"
)

// Volume Expansion
const (
	VolumeExpansion          = "VolumeExpansion"
	VolumeExpansionSucceeded = "VolumeExpansionSucceeded"
	VolumeExpansionFailed    = "VolumeExpansionFailed"
)

// Reconfigure
const (
	Reconfigure          = "Reconfigure"
	ReconfigureSucceeded = "ReconfigureSucceeded"
	ReconfigureFailed    = "ReconfigureFailed"
)

// ReconfigureTLS
const (
	ReconfigureTLS          = "ReconfigureTLS"
	ReconfigureTLSSucceeded = "ReconfigureTLSSucceeded"
	ReconfigureTLSFailed    = "ReconfigureTLSFailed"

	RemoveTLS                  = "RemoveTLS"
	RemoveTLSSucceeded         = "RemoveTLSSucceeded"
	RemoveTLSFailed            = "RemoveTLSFailed"
	AddTLS                     = "AddTLS"
	AddTLSSucceeded            = "AddTLSSucceeded"
	AddTLSFailed               = "AddTLSFailed"
	Issuing                    = "Issuing"
	IssueCertificatesSucceeded = "IssueCertificatesSucceeded"
	CertificateSynced          = "CertificateSynced"
	IssueCertificatesFailed    = "IssueCertificatesFailed"
)

// Restart
const (
	Restart              = "Restart"
	RestartNodes         = "RestartNodes"
	RestartPods          = "RestartPods"
	RestartPodsSucceeded = "RestartPodsSucceeded"
	RestartPodsFailed    = "RestartPodsFailed"
)

// StatefulSets
const (
	UpdateStatefulSets          = "UpdateStatefulSets"
	UpdateStatefulSetsSucceeded = "UpdateStatefulSetsSucceeded"
	UpdateStatefulSetsFailed    = "UpdateStatefulSetsFailed"
	ReadyStatefulSets           = "ReadyStatefulSets"
	DeleteStatefulSets          = "DeleteStatefulSets"
)

// Stash
const (
	PauseBackupConfiguration  = "PauseBackupConfiguration"
	ResumeBackupConfiguration = "ResumeBackupConfiguration"
)

// **********************************  Database Specifics ************************************

// Elasticsearch Constant
const (
	OrphanStatefulSetPods     = "OrphanStatefulSetPods"
	PrepareCustomConfig       = "PrepareCustomConfig"
	PrepareSecureCustomConfig = "PrepareSecureCustomConfig"
	ReconfigureSecurityAdmin  = "ReconfigureSecurityAdmin"

	HorizontalScaleMasterNode       = "HorizontalScaleMasterNode"
	HorizontalScaleDataNode         = "HorizontalScaleDataNode"
	HorizontalScaleDataHotNode      = "HorizontalScaleDataHotNode"
	HorizontalScaleDataWarmNode     = "HorizontalScaleDataWarmNode"
	HorizontalScaleDataColdNode     = "HorizontalScaleDataColdNode"
	HorizontalScaleDataFrozenNode   = "HorizontalScaleDataFrozenNode"
	HorizontalScaleDataContentNode  = "HorizontalScaleDataContentNode"
	HorizontalScaleMLNode           = "HorizontalScaleMLNode"
	HorizontalScaleTransformNode    = "HorizontalScaleTransformNode"
	HorizontalScaleCoordinatingNode = "HorizontalScaleCoordinatingNode"
	HorizontalScaleIngestNode       = "HorizontalScaleIngestNode"
	HorizontalScaleCombinedNode     = "HorizontalScaleCombinedNode"

	VolumeExpansionCombinedNode     = "VolumeExpansionCombinedNode"
	VolumeExpansionMasterNode       = "VolumeExpansionMasterNode"
	VolumeExpansionIngestNode       = "VolumeExpansionIngestNode"
	VolumeExpansionDataNode         = "VolumeExpansionDataNode"
	VolumeExpansionDataContentNode  = "VolumeExpansionDataContentNode"
	VolumeExpansionDataHotNode      = "VolumeExpansionDataHotNode"
	VolumeExpansionDataWarmNode     = "VolumeExpansionDataWarmNode"
	VolumeExpansionDataColdNode     = "VolumeExpansionDataColdNode"
	VolumeExpansionDataFrozenNode   = "VolumeExpansionDataFrozenNode"
	VolumeExpansionMLNode           = "VolumeExpansionMLNode"
	VolumeExpansionTransformNode    = "VolumeExpansionTransformNode"
	VolumeExpansionCoordinatingNode = "VolumeExpansionCoordinatingNode"
)

// Kafka Constants
const (
	ScaleUpBroker       = "ScaleUpBroker"
	ScaleUpController   = "ScaleUpController"
	ScaleUpCombined     = "ScaleUpCombined"
	ScaleDownBroker     = "ScaleDownBroker"
	ScaleDownController = "ScaleDownController"
	ScaleDownCombined   = "ScaleDownCombined"

	UpdateBrokerNodePVCs     = "UpdateBrokerNodePVCs"
	UpdateControllerNodePVCs = "UpdateControllerNodePVCs"
	UpdateCombinedNodePVCs   = "UpdateCombinedNodePVCs"
)

// MongoDB Constants
const (
	StartingBalancer  = "StartingBalancer"
	StoppingBalancer  = "StoppingBalancer"
	FlushRouterConfig = "FlushRouterConfig"

	UpdateStandaloneImage   = "UpdateStandaloneImage"
	UpdateShardImage        = "UpdateShardImage"
	UpdateReplicaSetImage   = "UpdateReplicaSetImage"
	UpdateConfigServerImage = "UpdateConfigServerImage"
	UpdateMongosImage       = "UpdateMongosImage"

	HorizontalScaleStandaloneUp      = "HorizontalScaleStandaloneUp"
	HorizontalScaleStandaloneDown    = "HorizontalScaleStandaloneDown"
	HorizontalScaleReplicaSetUp      = "HorizontalScaleReplicaSetUp"
	HorizontalScaleReplicaSetDown    = "HorizontalScaleReplicaSetDown"
	HorizontalScaleMongos            = "HorizontalScaleMongos"
	HorizontalScaleConfigServerUp    = "HorizontalScaleConfigServerUp"
	HorizontalScaleConfigServerDown  = "HorizontalScaleConfigServerDown"
	HorizontalScaleShardReplicasUp   = "HorizontalScaleShardReplicasUp"
	HorizontalScaleShardReplicasDown = "HorizontalScaleShardReplicasDown"
	HorizontalScaleShardUp           = "HorizontalScaleShardUp"
	HorizontalScaleShardDown         = "HorizontalScaleShardDown"
	HorizontalScaleArbiterUp         = "HorizontalScaleArbiterUp"
	HorizontalScaleArbiterDown       = "HorizontalScaleArbiterDown"
	HorizontalScaleHiddenUp          = "HorizontalScaleHiddenUp"
	HorizontalScaleHiddenDown        = "HorizontalScaleHiddenDown"

	VerticalScaleStandalone   = "VerticalScaleStandalone"
	VerticalScaleReplicaSet   = "VerticalScaleReplicaSet"
	VerticalScaleMongos       = "VerticalScaleMongos"
	VerticalScaleConfigServer = "VerticalScaleConfigServer"
	VerticalScaleShard        = "VerticalScaleShard"
	VerticalScaleArbiter      = "VerticalScaleArbiter"
	VerticalScaleHidden       = "VerticalScaleHidden"

	VolumeExpansionStandalone   = "VolumeExpansionStandalone"
	VolumeExpansionReplicaSet   = "VolumeExpansionReplicaSet"
	VolumeExpansionMongos       = "VolumeExpansionMongos"
	VolumeExpansionConfigServer = "VolumeExpansionConfigServer"
	VolumeExpansionShard        = "VolumeExpansionShard"
	VolumeExpansionHidden       = "VolumeExpansionHidden"

	ReconfigureStandalone   = "ReconfigureStandalone"
	ReconfigureReplicaset   = "ReconfigureReplicaset"
	ReconfigureMongos       = "ReconfigureMongos"
	ReconfigureConfigServer = "ReconfigureConfigServer"
	ReconfigureShard        = "ReconfigureShard"
	ReconfigureArbiter      = "ReconfigureArbiter"
	ReconfigureHidden       = "ReconfigureHidden"

	RestartStandalone   = "RestartStandalone"
	RestartReplicaSet   = "RestartReplicaSet"
	RestartMongos       = "RestartMongos"
	RestartConfigServer = "RestartConfigServer"
	RestartShard        = "RestartShard"
	RestartArbiter      = "RestartArbiter"
	RestartHidden       = "RestartHidden"
)

// MySQL/MariaDB Constants
const (
	TempIniFilesPath = "/tmp/kubedb-custom-ini-files"
)

// Postgres Constants
const (
	PausePgCoordinator                                   = "PausePgCoordinator"
	ResumePgCoordinator                                  = "ResumePgCoordinator"
	DataDirectoryInitialized                             = "DataDirectoryInitialized"
	ReplacedDataDirectory                                = "ReplacedDataDirectory"
	TransferLeaderShipToFirstNode                        = "TransferPrimaryRoleToDefault"
	TransferLeaderShipToFirstNodeBeforeCoordinatorPaused = "TransferLeaderShipToFirstNodeBeforeCoordinatorPaused"
	CopiedOldBinaries                                    = "CopiedOldBinaries"
	ResumePrimaryPgCoordinator                           = "NonTransferableResumePgCoordinator"

	UpdatePrimaryImage = "UpdatePrimaryImage"
	UpdateStandbyImage = "UpdateStandbyImage"

	RestartPrimary   = "RestartPrimary"
	RestartSecondary = "RestartSecondary"
)

// Redis Constants
const (
	PatchedSecret                        = "patchedSecret"
	ConfigKeyRedis                       = "redis.conf"
	RedisTLSArg                          = "--tls-port 6379"
	ReplaceSentinel                      = "ReplaceSentinel"
	ScaleUpRedisReplicasInSentinelMode   = "ScaleUpRedisReplicasInSentinelMode"
	ScaleDownRedisReplicasInSentinelMode = "ScaleDownRedisReplicasInSentinelMode"

	HorizontalScaleReplicasUp   = "HorizontalScaleReplicasUp"
	HorizontalScaleReplicasDown = "HorizontalScaleReplicasDown"
	HorizontalScaleSentinelUp   = "HorizontalScaleSentinelUp"
	HorizontalScaleSentinelDown = "HorizontalScaleSentinelDown"
)
