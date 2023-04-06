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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kmapi "kmodules.xyz/client-go/api/v1"
)

// List of possible condition types for a ops request
const (
	AccessApproved               = "Approved"
	AccessDenied                 = "Denied"
	DisableSharding              = "DisableSharding"
	EnableSharding               = "EnableSharding"
	Failed                       = "Failed"
	HorizontalScalingDatabase    = "HorizontalScaling"
	MigratingData                = "MigratingData"
	NodeCreated                  = "NodeCreated"
	NodeDeleted                  = "NodeDeleted"
	NodeRestarted                = "NodeRestarted"
	PauseDatabase                = "PauseDatabase"
	Progressing                  = "Progressing"
	ResumeDatabase               = "ResumeDatabase"
	ScalingDatabase              = "Scaling"
	ScalingDown                  = "ScalingDown"
	ScalingUp                    = "ScalingUp"
	Successful                   = "Successful"
	Running                      = "Running"
	Updating                     = "Updating"
	Upgrading                    = "Upgrading"
	UpgradeVersion               = "UpgradeVersion"
	VerticalScalingDatabase      = "VerticalScaling"
	VotingExclusionAdded         = "VotingExclusionAdded"
	VotingExclusionDeleted       = "VotingExclusionDeleted"
	UpdateStatefulSets           = "UpdateStatefulSets"
	VolumeExpansion              = "VolumeExpansion"
	Reconfigure                  = "Reconfigure"
	UpgradeNodes                 = "UpgradeNodes"
	RestartNodes                 = "RestartNodes"
	TLSRemoved                   = "TLSRemoved"
	TLSAdded                     = "TLSAdded"
	TLSChanged                   = "TLSChanged"
	IssuingConditionUpdated      = "IssuingConditionUpdated"
	CertificateIssuingSuccessful = "CertificateIssuingSuccessful"
	TLSEnabling                  = "TLSEnabling"
	Restart                      = "Restart"
	RestartStatefulSet           = "RestartStatefulSet"
	CertificateSynced            = "CertificateSynced"
	Reconciled                   = "Reconciled"
	RestartStatefulSetPods       = "RestartStatefulSetPods"

	// MongoDB Constants
	StartingBalancer            = "StartingBalancer"
	StoppingBalancer            = "StoppingBalancer"
	UpdateShardImage            = "UpdateShardImage"
	UpdateStatefulSetResources  = "UpdateStatefulSetResources"
	UpdateShardResources        = "UpdateShardResources"
	UpdateArbiterResources      = "UpdateArbiterResources"
	UpdateHiddenResources       = "UpdateHiddenResources"
	ScaleDownShard              = "ScaleDownShard"
	ScaleUpShard                = "ScaleUpShard"
	ScaleDownHidden             = "ScaleDownHidden"
	ScaleUpHidden               = "ScaleUpHidden"
	UpdateReplicaSetImage       = "UpdateReplicaSetImage"
	UpdateConfigServerImage     = "UpdateConfigServerImage"
	UpdateMongosImage           = "UpdateMongosImage"
	UpdateReplicaSetResources   = "UpdateReplicaSetResources"
	UpdateConfigServerResources = "UpdateConfigServerResources"
	UpdateMongosResources       = "UpdateMongosResources"
	FlushRouterConfig           = "FlushRouterConfig"
	ScaleDownReplicaSet         = "ScaleDownReplicaSet"
	ScaleUpReplicaSet           = "ScaleUpReplicaSet"
	ScaleUpShardReplicas        = "ScaleUpShardReplicas"
	ScaleDownShardReplicas      = "ScaleDownShardReplicas"
	ScaleDownConfigServer       = "ScaleDownConfigServer "
	ScaleUpConfigServer         = "ScaleUpConfigServer "
	ScaleMongos                 = "ScaleMongos"
	ReconfigureReplicaset       = "ReconfigureReplicaset"
	ReconfigureStandalone       = "ReconfigureStandalone"
	ReconfigureMongos           = "ReconfigureMongos"
	ReconfigureShard            = "ReconfigureShard"
	ReconfigureConfigServer     = "ReconfigureConfigServer"
	ReconfigureArbiter          = "ReconfigureArbiter"
	ReconfigureHidden           = "ReconfigureHidden"
	UpdateStandaloneImage       = "UpdateStandaloneImage"
	UpdateStandaloneResources   = "UpdateStandaloneResources"
	ScaleDownStandalone         = "ScaleDownStandalone"
	ScaleUpStandalone           = "ScaleUpStandalone"
	StandaloneVolumeExpansion   = "StandaloneVolumeExpansion"
	ReplicasetVolumeExpansion   = "ReplicasetVolumeExpansion"
	ShardVolumeExpansion        = "ShardVolumeExpansion"
	HiddenVolumeExpansion       = "HiddenVolumeExpansion"
	ConfigServerVolumeExpansion = "ConfigServerVolumeExpansion"
	RestartStandalone           = "RestartStandalone"
	RestartReplicaSet           = "RestartReplicaSet"
	RestartMongos               = "RestartMongos"
	RestartConfigServer         = "RestartConfigServer"
	RestartShard                = "RestartShard"
	RestartArbiter              = "RestartArbiter"
	RestartHidden               = "RestartHidden"
	DeleteStatefulSets          = "DeleteStatefulSets"
	DatabaseReady               = "DatabaseReady"

	// Elasticsearch Constant
	OrphanStatefulSetPods      = "OrphanStatefulSetPods"
	ReadyStatefulSets          = "ReadyStatefulSets"
	ScaleMasterNode            = "ScaleMasterNode"
	ScaleDataNode              = "ScaleDataNode"
	ScaleDataHotNode           = "ScaleDataHotNode"
	ScaleDataWarmNode          = "ScaleDataWarmNode"
	ScaleDataColdNode          = "ScaleDataColdNode"
	ScaleDataFrozenNode        = "ScaleDataFrozenNode"
	ScaleDataContentNode       = "ScaleDataContentNode"
	ScaleMLNode                = "ScaleMLNode"
	ScaleTransformNode         = "ScaleTransformNode"
	ScaleCoordinatingNode      = "ScaleCoordinatingNode"
	ScaleIngestNode            = "ScaleIngestNode"
	ScaleCombinedNode          = "ScaleCombinedNode"
	UpdateCombinedNodePVCs     = "UpdateCombinedNodePVCs"
	UpdateMasterNodePVCs       = "UpdateMasterNodePVCs"
	UpdateIngestNodePVCs       = "UpdateIngestNodePVCs"
	UpdateDataNodePVCs         = "UpdateDataNodePVCs"
	UpdateDataContentNodePVCs  = "UpdateDataContentNodePVCs"
	UpdateDataHotNodePVCs      = "UpdateDataHotNodePVCs"
	UpdateDataWarmNodePVCs     = "UpdateDataWarmNodePVCs"
	UpdateDataColdNodePVCs     = "UpdateDataColdNodePVCs"
	UpdateDataFrozenNodePVCs   = "UpdateDataFrozenNodePVCs"
	UpdateMLNodePVCs           = "UpdateMLNodePVCs"
	UpdateTransformNodePVCs    = "UpdateTransformNodePVCs"
	UpdateCoordinatingNodePVCs = "UpdateCoordinatingNodePVCs"
	UpdateElasticsearchCR      = "UpdateElasticsearchCR"

	UpdateNodeResources                = "UpdateNodeResources"
	UpdateMasterStatefulSetResources   = "UpdateMasterStatefulSetResources"
	UpdateDataStatefulSetResources     = "UpdateDataStatefulSetResources"
	UpdateIngestStatefulSetResources   = "UpdateIngestStatefulSetResources"
	UpdateCombinedStatefulSetResources = "UpdateCombinedStatefulSetResources"
	UpdateMasterNodeResources          = "UpdateMasterNodeResources"
	UpdateDataNodeResources            = "UpdateDataNodeResources"
	UpdateIngestNodeResources          = "UpdateIngestNodeResources"
	UpdateCombinedNodeResources        = "UpdateCombinedNodeResources"
	PrepareCustomConfig                = "PrepareCustomConfig"
	PrepareSecureCustomConfig          = "PrepareSecureCustomConfig"
	ReconfigureSecurityAdmin           = "ReconfigureSecurityAdmin"

	// Redis Constants
	PatchedSecret                        = "patchedSecret"
	ConfigKeyRedis                       = "redis.conf"
	RedisTLSArg                          = "--tls-port 6379"
	DBReady                              = "DBReady"
	RestartedPods                        = "RestartedPods"
	ScaleUpReplicas                      = "ScaleUpReplicas"
	ScaleDownReplicas                    = "ScaleDownReplicas"
	ScaleUpSentinel                      = "ScaleUpSentinel"
	ScaleDownSentinel                    = "ScaleDownSentinel"
	UpdateRedisImage                     = "UpdateRedisImage"
	RestartPodWithResources              = "RestartedPodsWithResources"
	ReplaceSentinel                      = "ReplaceSentinel"
	ScaleUpRedisReplicasInSentinelMode   = "ScaleUpRedisReplicasInSentinelMode"
	ScaleDownRedisReplicasInSentinelMode = "ScaleDownRedisReplicasInSentinelMode"

	// Stash Constants
	PauseBackupConfiguration  = "PauseBackupConfiguration"
	ResumeBackupConfiguration = "ResumeBackupConfiguration"
	// Postgres Constants
	UpdatePrimaryPodImage = "UpdatePrimaryImage"
	UpdateStandbyPodImage = "UpdateStandbyPodImage"
	// PausePgCoordinator is used when need to pause postgres failover with pg coordinator.
	// This is useful when we don't want failover for a certain period.
	PausePgCoordinator = "PausePgCoordinator"
	// ResumePgCoordinator is used when need to resume postgres failover with pg coordinator.
	// This is set when we are done with all the process necessary to do failover again.
	ResumePgCoordinator = "ResumePgCoordinator"
	// DataDirectoryInitialized condition is used in major upgrade ops request.
	// In major upgrade we need to initialized new directory wit initDB to run pg_upgrade.
	DataDirectoryInitialized = "DataDirectoryInitialized"
	// PgUpgraded is set when pg_upgrade command ran successfully.
	// This is used in major upgrade.
	PgUpgraded = "PgUpgraded"
	// ReplacedDataDirectory condition is used in major upgrade. After pg_upgrade we need to replace old data directory with new one.
	// after replace data directory successfully set this condition true.
	ReplacedDataDirectory                   = "ReplacedDataDirectory"
	PgCoordinatorStatusResumeDefaultPrimary = "ResumeDefaultPrimary"
	PostgresPrimaryPodReady                 = "PostgresPrimaryPodReady"
	RestartPrimaryPods                      = "RestartPrimaryPods"
	RestartStandbyPods                      = "RestartStandbyPods"
	// TransferLeaderShipToFirstNode is set when we need to set the pod-0 as primary.
	// This condition is set after pod-0 restart process done.
	TransferLeaderShipToFirstNode = "TransferPrimaryRoleToDefault"
	// TransferLeaderShipToFirstNodeBeforeCoordinatorPaused is set when we need to set the pod-0 as primary Before pgcoordinator paused
	// This is the initial step where we need to set pod-0 as primary. the condition is set before the pod-0 restart process.
	TransferLeaderShipToFirstNodeBeforeCoordinatorPaused = "TransferLeaderShipToFirstNodeBeforeCoordinatorPaused"
	// CopiedOldBinaries condition is used when we are done copying old postgres binary.
	// This is needed when we are doing major upgrade.
	CopiedOldBinaries      = "CopiedOldBinaries"
	UpdateStatefulSetImage = "UpdateStatefulSetImage"
	// ResumePrimaryPgCoordinator condition is set when we have set pg-coordinator status to NonTranferableResume this is useful when primary need to run after restart.
	ResumePrimaryPgCoordinator = "NonTransferableResumePgCoordinator"

	ReconfigurePrimaryPod  = "ReconfigurePrimaryPod"
	ReconfigureStandbyPods = "ReconfigureStandbyPods"

	// MySQL/MariaDB Constants
	TempIniFilesPath = "/tmp/kubedb-custom-ini-files"
)

type OpsRequestStatus struct {
	// Specifies the current phase of the ops request
	// +optional
	Phase OpsRequestPhase `json:"phase,omitempty"`
	// observedGeneration is the most recent generation observed for this resource. It corresponds to the
	// resource's generation, which is updated on mutation by the API Server.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
	// Conditions applied to the request, such as approval or denial.
	// +optional
	Conditions []kmapi.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:validation:Enum=Pending;Progressing;Successful;WaitingForApproval;Failed;Approved;Denied;Skipped
type OpsRequestPhase string

const (
	// used for ops requests that are currently in queue
	OpsRequestPhasePending OpsRequestPhase = "Pending"
	// used for ops requests that are currently Progressing
	OpsRequestPhaseProgressing OpsRequestPhase = "Progressing"
	// used for ops requests that are executed successfully
	OpsRequestPhaseSuccessful OpsRequestPhase = "Successful"
	// used for ops requests that are failed
	OpsRequestPhaseFailed OpsRequestPhase = "Failed"
	// used for ops requests that are skipped
	OpsRequestPhaseSkipped OpsRequestPhase = "Skipped"

	// Approval-related Phases

	// used for ops requests that are waiting for approval
	OpsRequestPhaseWaitingForApproval OpsRequestPhase = "WaitingForApproval"
	// used for ops requests that are approved
	OpsRequestApproved OpsRequestPhase = "Approved"
	// used for ops requests that are denied
	OpsRequestDenied OpsRequestPhase = "Denied"
)

// +kubebuilder:validation:Enum=Offline;Online
type VolumeExpansionMode string

const (
	// used to define a Online volume expansion mode
	VolumeExpansionModeOnline VolumeExpansionMode = "Online"
	// used to define a Offline volume expansion mode
	VolumeExpansionModeOffline VolumeExpansionMode = "Offline"
)

type RestartSpec struct{}

type Reprovision struct{}

type TLSSpec struct {
	// TLSConfig contains updated tls configurations for client and server.
	// +optional
	kmapi.TLSConfig `json:",inline,omitempty"`

	// RotateCertificates tells operator to initiate certificate rotation
	// +optional
	RotateCertificates bool `json:"rotateCertificates,omitempty"`

	// Remove tells operator to remove TLS configuration
	// +optional
	Remove bool `json:"remove,omitempty"`
}

// +kubebuilder:validation:Enum=IfReady;Always
type ApplyOption string

const (
	ApplyOptionIfReady ApplyOption = "IfReady"
	ApplyOptionAlways  ApplyOption = "Always"
)

type Accessor interface {
	GetObjectMeta() metav1.ObjectMeta
	GetDBRefName() string
	GetRequestType() any
	GetStatus() OpsRequestStatus
	SetStatus(_ OpsRequestStatus)
}
