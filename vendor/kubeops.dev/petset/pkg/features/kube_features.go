/*
Copyright 2017 The Kubernetes Authors.

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

package features

import (
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/component-base/featuregate"
)

const (
	// owner: @damemi
	// alpha: v1.21
	// beta: v1.22
	//
	// Enables scaling down replicas via logarithmic comparison of creation/ready timestamps
	LogarithmicScaleDown featuregate.Feature = "LogarithmicScaleDown"

	// owner: @krmayankk
	// alpha: v1.24
	//
	// Enables maxUnavailable for PetSet
	MaxUnavailablePetSet featuregate.Feature = "MaxUnavailablePetSet"

	// owner: @mattcary
	// alpha: v1.22
	// beta: v1.27
	//
	// Enables policies controlling deletion of PVCs created by a PetSet.
	PetSetAutoDeletePVC featuregate.Feature = "PetSetAutoDeletePVC"

	// owner: @psch
	// alpha: v1.26
	// beta: v1.27
	//
	// Enables a PetSet to start from an arbitrary non zero ordinal
	PetSetStartOrdinal featuregate.Feature = "PetSetStartOrdinal"

	// owner: @ahg-g
	// alpha: v1.21
	// beta: v1.22
	//
	// Enables controlling pod ranking on replicaset scale-down.
	PodDeletionCost featuregate.Feature = "PodDeletionCost"

	// owner: @danielvegamyhre
	// kep: https://kep.k8s.io/4017
	// beta: v1.28
	//
	// Set pod completion index as a pod label for Indexed Jobs.
	PodIndexLabel featuregate.Feature = "PodIndexLabel"
)

func init() {
	runtime.Must(DefaultMutableFeatureGate.Add(defaultPetSetFeatureGates))
}

// defaultPetSetFeatureGates consists of all known Kubernetes-specific feature keys.
// To add a new feature, define a key for it above and add it here. The features will be
// available throughout Kubernetes binaries.
//
// Entries are separated from each other with blank lines to avoid sweeping gofmt changes
// when adding or removing one entry.
var defaultPetSetFeatureGates = map[featuregate.Feature]featuregate.FeatureSpec{
	LogarithmicScaleDown: {Default: true, PreRelease: featuregate.Beta},

	MaxUnavailablePetSet: {Default: false, PreRelease: featuregate.Alpha},

	PetSetAutoDeletePVC: {Default: true, PreRelease: featuregate.Beta},

	PetSetStartOrdinal: {Default: true, PreRelease: featuregate.Beta},

	PodDeletionCost: {Default: true, PreRelease: featuregate.Beta},

	PodIndexLabel: {Default: true, PreRelease: featuregate.Beta},
}
