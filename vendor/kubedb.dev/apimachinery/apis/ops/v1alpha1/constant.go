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

const (
	Retrying = "Retrying"
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

// RotateAuth
const (
	RotateAuth                     = "RotateAuth"
	UpdateCredential               = "UpdateCredential"
	BasicAuthPreviousUsernameKey   = "username.prev"
	BasicAuthPreviousPasswordKey   = "password.prev"
	BasicAuthNextUsernameKey       = "username.next"
	BasicAuthNextPasswordKey       = "password.next"
	SecretAlreadyUpdatedAnnotation = "secret-already-updated"
	AuthDataPreviousKey            = "authData.prev"
	PatchDefaultConfig             = "PatchDefaultConfig"
)

// Restart
const (
	Restart              = "Restart"
	RestartNodes         = "RestartNodes"
	RestartPods          = "RestartPods"
	RestartKeeperPods    = "RestartKeeperPods"
	RestartPodsSucceeded = "RestartPodsSucceeded"
	RestartPodsFailed    = "RestartPodsFailed"
)

// Reload
const (
	ReloadPods          = "ReloadPods"
	ReloadPodsSucceeded = "ReloadPodsSucceeded"
	ReloadPodsFailed    = "ReloadPodsFailed"
)

// StatefulSets
const (
	UpdateStatefulSets          = "UpdateStatefulSets"
	UpdateStatefulSetsSucceeded = "UpdateStatefulSetsSucceeded"
	UpdateStatefulSetsFailed    = "UpdateStatefulSetsFailed"
	ReadyStatefulSets           = "ReadyStatefulSets"
	DeleteStatefulSets          = "DeleteStatefulSets"
	OrphanStatefulSetPods       = "OrphanStatefulSetPods"
)

// PetSets
const (
	UpdatePetSets          = "UpdatePetSets"
	UpdatePetSetsSucceeded = "UpdatePetSetsSucceeded"
	UpdatePetSetsFailed    = "UpdatePetSetsFailed"
	ReadyPetSets           = "ReadyPetSets"
	DeletePetSets          = "DeletePetSets"
	OrphanPetSetPods       = "OrphanPetSetPods"
)

// Stash
const (
	PauseBackupConfiguration  = "PauseBackupConfiguration"
	ResumeBackupConfiguration = "ResumeBackupConfiguration"
)

// **********************************  Database Specifics ************************************

// Elasticsearch Constant
const (
	PrepareCustomConfig               = "PrepareCustomConfig"
	PrepareSecureCustomConfig         = "PrepareSecureCustomConfig"
	ReconfigureSecurityAdmin          = "ReconfigureSecurityAdmin"
	DisabledMasterNodeShardAllocation = "DisabledMasterNodeShardAllocation"

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
	HorizontalScaleOverseerNode     = "HorizontalScaleOverseerNode"
	HorizontalScaleCoordinatorNode  = "HorizontalScaleCoordinatorNode"

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
	VolumeExpansionOverseerNode     = "VolumeExpansionOverseerNode"
	VolumeExpansionCoordinatorNode  = "VolumeExpansionCoordinatorNode"
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

// MSSQLServer Constants
const (
	PrepareApplyConfig = "PrepareApplyConfig"
)

// Singlestore Constants
const (
	ScaleUpAggregator   = "ScaleUpAggregator"
	ScaleDownAggregator = "ScaleDownAggregator"
	ScaleUpLeaf         = "ScaleUpLeaf"
	ScaleDownLeaf       = "ScaleDownLeaf"
)

// RabbitMQ Constants
const (
	UpdateNodePVCs        = "UpdateNodePVCs"
	EnableAllFeatureFlags = "EnableAllFeatureFlags"
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

	SetHorizons    = "SetHorizons"
	RemoveHorizons = "RemoveHorizons"
)

// MySQL/MariaDB/Maxscale Constants
const (
	TempIniFilesPath             = "/tmp/kubedb-custom-ini-files"
	StopRemoteReplica            = "StopRemoteReplica"
	DBPatch                      = "DBPatch"
	StopRemoteReplicaSucceeded   = "StopRemoteReplicaSucceeded"
	DBPatchSucceeded             = "DBPatchSucceeded"
	RestartMaxscale              = "RestartMaxscale"
	RestartMaxscalePodsSucceeded = "RestartMaxscalePodsSucceeded"
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

	StartRunScript                          = "StartRunScriptWithRestart"
	KillRunScript                           = "KillRunScript"
	StickyLeader                            = "STICKYLEADER" // We want a id(sticky id) to be always leader in raft
	UpdateDataDirectory                     = "UpdateDataDirectory"
	PausePgCoordinatorBeforeUpgrade         = "PausePgCoordinatorBeforeUpdate"
	RunningLeaderSticky                     = "RunningLeaderSticky"
	EnsureStickyId                          = "EnsureStickyId"
	SetPrimaryPodNameInStatus               = "SetPrimaryPodNameInStatus"
	NonTransferableResumeAfterUpgrade       = "NonTransferableResumeAfterUpgrade"
	PausePgCoordinatorBeforeCustomRestart   = "PausePgCoordinatorBeforeCustomRestart"
	NonTransferableResumeAfterCustomRestart = "NonTransferableResumeAfterCustomRestart"
	ReadyPrimaryCheck                       = "ReadyPrimaryCheck"
	PrimaryPodName                          = "PrimaryPodName"
	OpsRequestProgressing                   = "OpsRequestProgressing"
	SetRaftKeyOpsRequestProgressing         = "SetRaftKeyOpsRequestProgressing"
	UnsetRaftKeyOpsRequestProgressing       = "UnsetRaftKeyOpsRequestProgressing"
	NotReadyReplicas                        = "NotReadyReplicas"
	RestartNotReadyStandby                  = "RestartNotReadyStandby"
	StandbyReadyCheck                       = "StandbyReadyCheck"
	PrimaryRunningCheck                     = "PrimaryRunningCheck"
	StopPostgresServer                      = "StopPostgresServer"
	SetupRecoverySettings                   = "SetupRecoverySettings"
	RunPostgresBaseBackup                   = "RunPostgresBaseBackup"
	StartPostgresRecovery                   = "StartPostgresRecovery"
	RestartNotReadyStandbyAfterBaseBackup   = "RestartNotReadyStandbyAfterBaseBackup"
	StandbyReadyCheckAfterBaseBackup        = "StandbyReadyCheckAfterBaseBackup"
	ReconnectStandbyWithRestart             = "ReconnectStandbyWithRestart"
	ReconnectStandbyWithBaseBackup          = "ReconnectStandbyWithBaseBackup"
	UpdateRaftKVStore                       = "UpdateRaftKVStore"
	CreateFailOverFile                      = "CreateFailOverFile"
	TransitionCandidateToLeader             = "TransitionCandidateToLeader"
	ForceFailOverFileName                   = "force-failover-with-lsn"
	// need reveiw about this two
	StringFalse = "false"
	StringTrue  = "true"
)

// Redis Constants
const (
	PatchedSecret                        = "patchedSecret"
	ConfigKeyRedis                       = "redis.conf"
	RedisTLSArg                          = "--tls-port 6379"
	ReplaceSentinel                      = "ReplaceSentinel"
	ScaleUpRedisReplicasInSentinelMode   = "ScaleUpRedisReplicasInSentinelMode"
	ScaleDownRedisReplicasInSentinelMode = "ScaleDownRedisReplicasInSentinelMode"

	RedisUpdateAnnounces        = "UpdateAnnounces"
	HorizontalScaleReplicasUp   = "HorizontalScaleReplicasUp"
	HorizontalScaleReplicasDown = "HorizontalScaleReplicasDown"
	HorizontalScaleSentinelUp   = "HorizontalScaleSentinelUp"
	HorizontalScaleSentinelDown = "HorizontalScaleSentinelDown"

	RedisUpdateAclSecret = "UpdateAclSecret"
)

// Druid Constants
const (
	ScaleUpCoordinators   = "ScaleUpCoordinators"
	ScaleUpOverlords      = "ScaleUpOverlords"
	ScaleUpBrokers        = "ScaleUpBrokers"
	ScaleUpHistoricals    = "ScaleUpHistoricals"
	ScaleUpMiddleManagers = "ScaleUpMiddleManagers"
	ScaleUpRouters        = "ScaleUpRouters"

	ScaleDownCoordinators   = "ScaleDownCoordinators"
	ScaleDownOverlords      = "ScaleDownOverlords"
	ScaleDownBrokers        = "ScaleDownBrokers"
	ScaleDownHistoricals    = "ScaleDownHistoricals"
	ScaleDownMiddleManagers = "ScaleDownMiddleManagers"
	ScaleDownRouters        = "ScaleDownRouters"

	UpdateMiddleManagersNodePVCs = "UpdateMiddleManagersNodePVCs"
	UpdateHistoricalsNodePVCs    = "UpdateHistoricalsNodePVCs"

	UpdateCredentialDynamically = "UpdateCredentialDynamically"
)

// SingleStore Constants
const (
	UpdateAggregatorNodePVCs = "UpdateAggregatorNodePVCs"
	UpdateLeafNodePVCs       = "UpdateLeafNodePVCs"
)

// PgBouncer Constants
const (
	UpdatePgBouncerBackendSecret = "UpdateBackendSecret"
	ConfigSecretDelete           = "ConfigSecretDeleted"
)

// Pgpool Constants
const (
	UpdateConfigSecret = "UpdateConfigSecret"
)
