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

package v1

import (
	"fmt"

	api "kubeops.dev/petset/apis/apps/v1"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

// PetSetListerExpansion allows custom methods to be added to
// PetSetLister.
type PetSetListerExpansion interface {
	GetPodPetSets(pod *v1.Pod) ([]*api.PetSet, error)
}

// PetSetNamespaceListerExpansion allows custom methods to be added to
// PetSetNamespaceLister.
type PetSetNamespaceListerExpansion interface{}

// GetPodPetSets returns a list of PetSets that potentially match a pod.
// Only the one specified in the Pod's ControllerRef will actually manage it.
// Returns an error only if no matching PetSets are found.
func (s *petSetLister) GetPodPetSets(pod *v1.Pod) ([]*api.PetSet, error) {
	var selector labels.Selector
	var ps *api.PetSet

	if len(pod.Labels) == 0 {
		return nil, fmt.Errorf("no PetSets found for pod %v because it has no labels", pod.Name)
	}

	list, err := s.PetSets(pod.Namespace).List(labels.Everything())
	if err != nil {
		return nil, err
	}

	var psList []*api.PetSet
	for i := range list {
		ps = list[i]
		if ps.Namespace != pod.Namespace {
			continue
		}
		selector, err = metav1.LabelSelectorAsSelector(ps.Spec.Selector)
		if err != nil {
			// This object has an invalid selector, it does not match the pod
			continue
		}

		// If a PetSet with a nil or empty selector creeps in, it should match nothing, not everything.
		if selector.Empty() || !selector.Matches(labels.Set(pod.Labels)) {
			continue
		}
		psList = append(psList, ps)
	}

	if len(psList) == 0 {
		return nil, fmt.Errorf("could not find PetSet for pod %s in namespace %s with labels: %v", pod.Name, pod.Namespace, pod.Labels)
	}

	return psList, nil
}
