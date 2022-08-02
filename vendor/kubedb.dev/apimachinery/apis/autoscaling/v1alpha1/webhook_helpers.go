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
	dbapi "kubedb.dev/apimachinery/apis/kubedb/v1alpha2"
	opsapi "kubedb.dev/apimachinery/apis/ops/v1alpha1"

	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func setDefaultStorageValues(storageSpec *StorageAutoscalerSpec) {
	if storageSpec == nil {
		return
	}
	if storageSpec.Trigger == "" {
		storageSpec.Trigger = AutoscalerTriggerOff
	}
	if storageSpec.ExpansionMode == nil {
		mode := opsapi.VolumeExpansionModeOnline
		storageSpec.ExpansionMode = &mode
	}
	if storageSpec.ScalingThreshold == 0 {
		storageSpec.ScalingThreshold = DefaultStorageScalingThreshold
	}
	if storageSpec.UsageThreshold == 0 {
		storageSpec.UsageThreshold = DefaultStorageUsageThreshold
	}
}

func setDefaultComputeValues(computeSpec *ComputeAutoscalerSpec) {
	if computeSpec == nil {
		return
	}
	if computeSpec.Trigger == "" {
		computeSpec.Trigger = AutoscalerTriggerOff
	}
	if computeSpec.ControlledResources == nil {
		computeSpec.ControlledResources = []core.ResourceName{core.ResourceCPU, core.ResourceMemory}
	}
	if computeSpec.ContainerControlledValues == nil {
		reqAndLim := ContainerControlledValuesRequestsAndLimits
		computeSpec.ContainerControlledValues = &reqAndLim
	}
	if computeSpec.ResourceDiffPercentage == 0 {
		computeSpec.ResourceDiffPercentage = DefaultResourceDiffPercentage
	}
	if computeSpec.PodLifeTimeThreshold.Duration == 0 {
		computeSpec.PodLifeTimeThreshold = metav1.Duration{Duration: DefaultPodLifeTimeThreshold}
	}
}

func setInMemoryDefaults(computeSpec *ComputeAutoscalerSpec, storageEngine dbapi.StorageEngine) {
	if computeSpec == nil || storageEngine != dbapi.StorageEngineInMemory {
		return
	}
	if computeSpec.InMemoryStorage == nil {
		// assigning a dummy pointer to set the defaults
		computeSpec.InMemoryStorage = &ComputeInMemoryStorageSpec{}
	}
	if computeSpec.InMemoryStorage.UsageThresholdPercentage == 0 {
		computeSpec.InMemoryStorage.UsageThresholdPercentage = DefaultInMemoryStorageUsageThresholdPercentage
	}
	if computeSpec.InMemoryStorage.ScalingFactorPercentage == 0 {
		computeSpec.InMemoryStorage.ScalingFactorPercentage = DefaultInMemoryStorageScalingFactorPercentage
	}
}
