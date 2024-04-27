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
	"fmt"
	"sort"
	"strconv"
	"strings"

	dbapi "kubedb.dev/apimachinery/apis/kubedb/v1alpha2"

	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
	"k8s.io/utils/ptr"
)

func setDefaultStorageValues(storageSpec *StorageAutoscalerSpec) {
	if storageSpec == nil {
		return
	}
	if storageSpec.Trigger == "" {
		storageSpec.Trigger = AutoscalerTriggerOff
	}
	if storageSpec.ScalingThreshold == nil {
		storageSpec.ScalingThreshold = ptr.To[int32](DefaultStorageScalingThreshold)
	}
	setDefaultScalingRules(storageSpec)
}

type quantity struct {
	InString   string
	InQuantity resource.Quantity
	Threshold  string
}

// sort them by the appliesUpto value. The items with empty threshold
func setDefaultScalingRules(storageSpec *StorageAutoscalerSpec) {
	if storageSpec.UsageThreshold == nil {
		storageSpec.UsageThreshold = ptr.To[int32](DefaultStorageUsageThreshold)
	}
	var quantities []quantity
	var zeroQuantityThresholds []string
	for _, sr := range storageSpec.ScalingRules {
		if sr.AppliesUpto == "" {
			zeroQuantityThresholds = append(zeroQuantityThresholds, sr.Threshold)
			continue
		}
		quantities = append(quantities, quantity{
			InString:   sr.AppliesUpto,
			InQuantity: resource.MustParse(sr.AppliesUpto),
			Threshold:  sr.Threshold,
		})
	}
	sort.Slice(quantities, func(i, j int) bool {
		return quantities[i].InQuantity.Cmp(quantities[j].InQuantity) < 0
	})

	storageSpec.ScalingRules = make([]StorageScalingRule, 0)
	for _, q := range quantities {
		storageSpec.ScalingRules = append(storageSpec.ScalingRules, StorageScalingRule{
			AppliesUpto: q.InString,
			Threshold:   q.Threshold,
		})
	}
	for _, threshold := range zeroQuantityThresholds {
		storageSpec.ScalingRules = append(storageSpec.ScalingRules, StorageScalingRule{
			AppliesUpto: "",
			Threshold:   threshold,
		})
	}
	klog.Infof("scaling Rules = %v \n", storageSpec.ScalingRules)
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

func validateScalingRules(storageSpec *StorageAutoscalerSpec) error {
	var zeroQuantityThresholds []string
	for _, sr := range storageSpec.ScalingRules {
		if sr.AppliesUpto == "" {
			zeroQuantityThresholds = append(zeroQuantityThresholds, sr.Threshold)
		}
		th := sr.Threshold
		if strings.HasSuffix(th, "%") {
			if !isNum(strings.TrimSuffix(th, "%")) {
				return fmt.Errorf("%v is not a valid percentage value", th)
			}
		} else if strings.HasSuffix(th, "pc") {
			if !isNum(strings.TrimSuffix(th, "pc")) {
				return fmt.Errorf("%v is not a valid percentage value", th)
			}
		} else {
			_, err := resource.ParseQuantity(sr.Threshold)
			if err != nil {
				return fmt.Errorf("%v is not a valid quatity", sr.Threshold)
			}
		}
	}
	if len(zeroQuantityThresholds) > 1 {
		return fmt.Errorf("%v appliesUpto value are empty in %v", zeroQuantityThresholds, storageSpec.ScalingRules)
	}
	return nil
}

func isNum(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}
