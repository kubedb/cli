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

// List of possible condition types for a ops request
const (
	AccessApproved            = "Approved"
	AccessDenied              = "Denied"
	DisableSharding           = "DisableSharding"
	EnableSharding            = "EnableSharding"
	Failure                   = "Failure"
	HorizontalScalingDatabase = "HorizontalScaling"
	MigratingData             = "MigratingData"
	NodeCreated               = "NodeCreated"
	NodeDeleted               = "NodeDeleted"
	NodeRestarted             = "NodeRestarted"
	PauseDatabase             = "PauseDatabase"
	Progressing               = "Progressing"
	ResumeDatabase            = "ResumeDatabase"
	ScalingDatabase           = "Scaling"
	ScalingDown               = "ScalingDown"
	ScalingUp                 = "ScalingUp"
	Successful                = "Successful"
	Updating                  = "Updating"
	Upgrading                 = "Upgrading"
	UpgradeVersion            = "UpgradeVersion"
	VerticalScalingDatabase   = "VerticalScaling"
	VotingExclusionAdded      = "VotingExclusionAdded"
	VotingExclusionDeleted    = "VotingExclusionDeleted"
	UpdateStatefulSets        = "UpdateStatefulSets"

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
	VolumeExpansion             = "VolumeExpansion"
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
)

type OpsRequestPhase string

const (
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

// +kubebuilder:validation:Enum=Upgrade;HorizontalScaling;VerticalScaling;VolumeExpansion;Restart;RotateCertificates;Reconfigure
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
	// used for RotateCertificates operation
	OpsRequestTypeRotateCertificates OpsRequestType = "RotateCertificates"
	// used for Reconfigure operation
	OpsRequestTypeReconfigure OpsRequestType = "Reconfigure"
)

type UpgradeSpec struct {
	// Specifies the target version name from catalog
	TargetVersion string `json:"targetVersion,omitempty" protobuf:"bytes,1,opt,name=targetVersion"`
}
