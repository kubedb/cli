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

import "time"

// Compute Autoscaler
const (
	// Ignore change priority that is smaller than 10%.
	DefaultResourceDiffPercentage = 10

	// Pods that live for at least that long can be evicted even if their
	// request is within the [MinRecommended...MaxRecommended] range.
	DefaultPodLifeTimeThreshold = time.Hour * 12

	DefaultInMemoryStorageUsageThresholdPercentage = 70
	DefaultInMemoryStorageScalingFactorPercentage  = 50
)

// Storage Autoscaler
const (
	DefaultStorageUsageThreshold   = 80
	DefaultStorageScalingThreshold = 50
)

const (
	// AutoscalerTriggerOn means the autoscaler is enabled.
	AutoscalerTriggerOn AutoscalerTrigger = "On"
	// AutoscalerTriggerOff means the autoscaler is disabled.
	AutoscalerTriggerOff AutoscalerTrigger = "Off"
)

const (
	// ContainerControlledValuesRequestsAndLimits means resource request and limits
	// are scaled automatically. The limit is scaled proportionally to the request.
	ContainerControlledValuesRequestsAndLimits ContainerControlledValues = "RequestsAndLimits"
	// ContainerControlledValuesRequestsOnly means only requested resource is autoscaled.
	ContainerControlledValuesRequestsOnly ContainerControlledValues = "RequestsOnly"
)

// List of possible condition types for an autoscaler
const (
	Failure          = "Failure"
	CreateOpsRequest = "CreateOpsRequest"
)

const (
	// AutoscalerPhaseInProgress is used when autoscaler is waiting for the initialization
	// if referred db is not found, It will also be in InProgress
	AutoscalerPhaseInProgress AutoscalerPhase = "InProgress"
	// AutoscalerPhaseCurrent is used as long as autoscaler is running properly
	AutoscalerPhaseCurrent AutoscalerPhase = "Current"
	// AutoscalerPhaseTerminating is used when an autoscaler object is being terminated
	AutoscalerPhaseTerminating AutoscalerPhase = "Terminating"
	// AutoscalerPhaseFailed is used when some unexpected error occurred
	AutoscalerPhaseFailed AutoscalerPhase = "Failed"
)
