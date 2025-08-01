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
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/klog/v2"
	apiworkv1 "open-cluster-management.io/api/work/v1"
)

// PetSetListerExpansion allows custom methods to be added to
// PetSetLister.
type PetSetListerExpansion interface {
	GetPodPetSets(pod *v1.Pod) ([]*api.PetSet, error)
	GetManifestWorkPetSets(mw *apiworkv1.ManifestWork) ([]*api.PetSet, error)
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

// GetManifestWorkPetSets returns a list of PetSets that potentially match a manifestwork.
// It lists PetSets across all namespaces and matches them based on labels.
func (s *petSetLister) GetManifestWorkPetSets(mw *apiworkv1.ManifestWork) ([]*api.PetSet, error) {
	if len(mw.Labels) == 0 {
		return nil, fmt.Errorf("no PetSets found for manifestwork %s because it has no labels", mw.Name)
	}

	list, err := s.List(labels.Everything())
	if err != nil {
		return nil, err
	}

	var psList []*api.PetSet
	for _, ps := range list {
		selector, err := metav1.LabelSelectorAsSelector(ps.Spec.Selector)
		if err != nil {
			klog.Warningf("PetSet %s/%s has an invalid selector: %v", ps.Namespace, ps.Name, err)
			continue
		}

		if selector.Empty() || !selector.Matches(labels.Set(mw.Labels)) {
			continue
		}
		psList = append(psList, ps)
	}

	if len(psList) == 0 {
		return nil, fmt.Errorf("could not find any PetSet for manifestwork %v/%s in any namespace with labels: %v", mw.Namespace, mw.Name, mw.Labels)
	}

	if len(psList) > 1 {
		setNames := []string{}
		for _, s := range psList {
			setNames = append(setNames, s.Name)
		}
		utilruntime.HandleError(
			fmt.Errorf(
				"user error: more than one PetSet is selecting manifestwork with labels: %+v. Sets: %v",
				mw.Labels, setNames))
	}

	return psList, nil
}
