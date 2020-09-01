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
	PausingDatabase           = "PausingDatabase"
	PausedDatabase            = "PausedDatabase"
	Progressing               = "Progressing"
	ResumeDatabase            = "ResumeDatabase"
	ResumingDatabase          = "ResumingDatabase"
	ResumedDatabase           = "ResumedDatabase"
	ScalingDatabase           = "Scaling"
	ScalingDown               = "ScalingDown"
	ScalingUp                 = "ScalingUp"
	Successful                = "Successful"
	Updating                  = "Updating"
	UpgradedVersion           = "UpgradedVersion"
	UpgradingVersion          = "UpgradingVersion"
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
)

type AutoscalerPhase string

const (
	// used for ops requests that are currently Progressing
	AutoscalerPhaseProgressing AutoscalerPhase = "Progressing"
	// used for ops requests that are executed successfully
	AutoscalerPhaseSuccessful AutoscalerPhase = "Successful"
	// used for ops requests that are waiting for approval
	AutoscalerPhaseWaitingForApproval AutoscalerPhase = "WaitingForApproval"
	// used for ops requests that are failed
	AutoscalerPhaseFailed AutoscalerPhase = "Failed"
	// used for ops requests that are approved
	AutoscalerApproved AutoscalerPhase = "Approved"
	// used for ops requests that are denied
	AutoscalerDenied AutoscalerPhase = "Denied"
)

// +kubebuilder:validation:Enum=Upgrade;HorizontalScaling;VerticalScaling;VolumeExpansion;Restart;RotateCertificates
type AutoscalerType string

const (
	// used for Upgrade operation
	AutoscalerTypeUpgrade AutoscalerType = "Upgrade"
	// used for HorizontalScaling operation
	AutoscalerTypeHorizontalScaling AutoscalerType = "HorizontalScaling"
	// used for VerticalScaling operation
	AutoscalerTypeVerticalScaling AutoscalerType = "VerticalScaling"
	// used for VolumeExpansion operation
	AutoscalerTypeVolumeExpansion AutoscalerType = "VolumeExpansion"
	// used for Restart operation
	AutoscalerTypeRestart AutoscalerType = "Restart"
	// used for RotateCertificates operation
	AutoscalerTypeRotateCertificates AutoscalerType = "RotateCertificates"
)

type UpgradeSpec struct {
	// Specifies the target version name from catalog
	TargetVersion string `json:"targetVersion,omitempty" protobuf:"bytes,1,opt,name=targetVersion"`
}
