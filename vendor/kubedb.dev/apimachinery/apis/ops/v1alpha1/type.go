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

	// MongoDB Constants
	StartingBalancer            = "StartingBalancer"
	StoppingBalancer            = "StoppingBalancer"
	UpdateShardImage            = "UpdateShardImage"
	UpdateStatefulSetResources  = "UpdateStatefulSetResources"
	UpdateShardResources        = "UpdateShardResources"
	ScaleDownShard              = "ScaleDownShard"
	ScaleUpShard                = "ScaleUpShard"
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
	ReconfigureMongos           = "ReconfigureMongos"
	ReconfigureShard            = "ReconfigureShard"
	ReconfigureConfigServer     = "ReconfigureConfigServer"
	UpdateStandaloneImage       = "UpdateStandaloneImage"
	UpdateStandaloneResources   = "UpdateStandaloneResources"
	ScaleDownStandalone         = "ScaleDownStandalone"
	ScaleUpStandalone           = "ScaleUpStandalone"
	ReconfigureStandalone       = "ReconfigureStandalone"
	StandaloneVolumeExpansion   = "StandaloneVolumeExpansion"
	ReplicasetVolumeExpansion   = "ReplicasetVolumeExpansion"
	ShardVolumeExpansion        = "ShardVolumeExpansion"
	ConfigServerVolumeExpansion = "ConfigServerVolumeExpansion"
	RestartStandalone           = "RestartStandalone"
	RestartReplicaSet           = "RestartReplicaSet"
	RestartMongos               = "RestartMongos"
	RestartConfigServer         = "RestartConfigServer"
	RestartShard                = "RestartShard"

	// Elasticsearch Constant
	OrphanStatefulSetPods  = "OrphanStatefulSetPods"
	ReadyStatefulSets      = "ReadyStatefulSets"
	ScaleDownCombinedNode  = "ScaleDownCombinedNode"
	ScaleDownDataNode      = "ScaleDownDataNode"
	ScaleDownIngestNode    = "ScaleDownIngestNode"
	ScaleDownMasterNode    = "ScaleDownMasterNode"
	ScaleUpCombinedNode    = "ScaleUpCombinedNode"
	ScaleUpDataNode        = "ScaleUpDataNode"
	ScaleUpIngestNode      = "ScaleUpIngestNode"
	ScaleUpMasterNode      = "ScaleUpMasterNode"
	UpdateCombinedNodePVCs = "UpdateCombinedNodePVCs"
	UpdateDataNodePVCs     = "UpdateDataNodePVCs"
	UpdateIngestNodePVCs   = "UpdateIngestNodePVCs"
	UpdateMasterNodePVCs   = "UpdateMasterNodePVCs"
	UpdateNodeResources    = "UpdateNodeResources"

	//Redis Constants
	PatchedSecret  = "patchedSecret"
	ConfigKeyRedis = "redis.conf"
	RedisTLSArg    = "--tls-port 6379"
	DBReady        = "DBReady"
	RestartedPods  = "RestartedPods"

	//Stash Constants
	PauseBackupConfiguration  = "PauseBackupConfiguration"
	ResumeBackupConfiguration = "ResumeBackupConfiguration"
)

// +kubebuilder:validation:Enum=Pending;Progressing;Successful;WaitingForApproval;Failed;Approved;Denied
type OpsRequestPhase string

const (
	// used for ops requests that are currently in queue
	OpsRequestPhasePending OpsRequestPhase = "Pending"
	// used for ops requests that are currently Progressing
	OpsRequestPhaseProgressing OpsRequestPhase = "Progressing"
	// used for ops requests that are executed successfully
	OpsRequestPhaseSuccessful OpsRequestPhase = "Successful"
	// used for ops requests that are waiting for approval
	OpsRequestPhaseWaitingForApproval OpsRequestPhase = "WaitingForApproval"
	// used for ops requests that are failed
	OpsRequestPhaseFailed OpsRequestPhase = "Failed"
	// used for ops requests that are approved
	OpsRequestApproved OpsRequestPhase = "Approved"
	// used for ops requests that are denied
	OpsRequestDenied OpsRequestPhase = "Denied"
)

// +kubebuilder:validation:Enum=Upgrade;HorizontalScaling;VerticalScaling;VolumeExpansion;Restart;Reconfigure;ReconfigureTLS
type OpsRequestType string

const (
	// used for Upgrade operation
	OpsRequestTypeUpgrade OpsRequestType = "Upgrade"
	// used for HorizontalScaling operation
	OpsRequestTypeHorizontalScaling OpsRequestType = "HorizontalScaling"
	// used for VerticalScaling operation
	OpsRequestTypeVerticalScaling OpsRequestType = "VerticalScaling"
	// used for VolumeExpansion operation
	OpsRequestTypeVolumeExpansion OpsRequestType = "VolumeExpansion"
	// used for Restart operation
	OpsRequestTypeRestart OpsRequestType = "Restart"
	// used for Reconfigure operation
	OpsRequestTypeReconfigure OpsRequestType = "Reconfigure"
	// used for ReconfigureTLS operation
	OpsRequestTypeReconfigureTLSs OpsRequestType = "ReconfigureTLS"
)

type RestartSpec struct {
}

type TLSSpec struct {
	// TLSConfig contains updated tls configurations for client and server.
	// +optional
	kmapi.TLSConfig `json:",inline,omitempty" protobuf:"bytes,1,opt,name=tLSConfig"`

	// RotateCertificates tells operator to initiate certificate rotation
	// +optional
	RotateCertificates bool `json:"rotateCertificates,omitempty" protobuf:"varint,2,opt,name=rotateCertificates"`

	// Remove tells operator to remove TLS configuration
	// +optional
	Remove bool `json:"remove,omitempty" protobuf:"varint,3,opt,name=remove"`
}
