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

	// ======================= Condition Reasons ========================

	OpsRequestProgressingStarted           = "OpsRequestProgressingStarted"
	OpsRequestFailedToProgressing          = "OpsRequestFailedToProgressing"
	SuccessfullyPausedDatabase             = "SuccessfullyPausedDatabase"
	FailedToPauseDatabase                  = "FailedToPauseDatabase"
	SuccessfullyResumedDatabase            = "SuccessfullyResumedDatabase"
	FailedToResumedDatabase                = "FailedToResumedDatabase"
	DatabaseVersionUpgradingStarted        = "DatabaseVersionUpgradingStarted"
	SuccessfullyUpgradedDatabaseVersion    = "SuccessfullyUpgradedDatabaseVersion"
	FailedToUpgradeDatabaseVersion         = "FailedToUpgradeDatabaseVersion"
	HorizontalScalingStarted               = "HorizontalScalingStarted"
	SuccessfullyPerformedHorizontalScaling = "SuccessfullyPerformedHorizontalScaling"
	FailedToPerformHorizontalScaling       = "FailedToPerformHorizontalScaling"
	VerticalScalingStarted                 = "VerticalScalingStarted"
	SuccessfullyPerformedVerticalScaling   = "SuccessfullyPerformedVerticalScaling"
	FailedToPerformVerticalScaling         = "FailedToPerformVerticalScaling"
	OpsRequestProcessedSuccessfully        = "OpsRequestProcessedSuccessfully"
	SuccessfullyVolumeExpanded             = "SuccessfullyVolumeExpanded"
	FailedToVolumeExpand                   = "FailedToVolumeExpand"
	SuccessfullyDBReconfigured             = "SuccessfullyDBReconfigured"
	FailedToReconfigureDB                  = "FailedToReconfigureDB"
	SuccessfullyRestartedDBMembers         = "SuccessfullyRestartedDBMembers"
	FailToRestartDBMembers                 = "FailToRestartDBMembers"
	SuccessfullyRestatedStatefulSet        = "SuccessfullyRestatedStatefulSet"
	FailedToRestartStatefulSet             = "FailedToRestartStatefulSet"
	SuccessfullyRemovedTLSConfig           = "SuccessfullyRemovedTLSConfig"
	FailedToRemoveTLSConfig                = "FailedToRemoveTLSConfig"
	SuccessfullyAddedTLSConfig             = "SuccessfullyAddedTLSConfig"
	FailedToAddTLSConfig                   = "FailedToAddTLSConfig"
	SuccessfullyIssuedCertificates         = "SuccessfullyIssuedCertificates"
	FailedToIssueCertificates              = "FailedToIssueCertificates"
	SuccessfullyReconfiguredTLS            = "SuccessfullyReconfiguredTLS"
)
