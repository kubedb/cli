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
	"time"

	"kubestash.dev/apimachinery/apis"
	"kubestash.dev/apimachinery/crds"

	"kmodules.xyz/client-go/apiextensions"
)

const (
	// DefaultCBTInterval mirrors the +kubebuilder:default on CBTBackupOptions.Interval,
	// for the paths that read a spec the API server never defaulted.
	DefaultCBTInterval = apis.DefaultCBTInterval
)

func (KubeVirtArchiver) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(GroupVersion.WithResource(ResourcePluralKubeVirtArchiver))
}

// GetCBTInterval is how long the resident loop sleeps between checkpoints.
func (k *KubeVirtArchiver) GetCBTInterval() time.Duration {
	if k.Spec.CBTBackup == nil || k.Spec.CBTBackup.Interval == nil {
		return DefaultCBTInterval
	}
	return k.Spec.CBTBackup.Interval.Duration
}

// UpsertVirtualMachineStatus replaces the entry for one VirtualMachine, keeping the
// list's order stable so a `kubectl describe` does not reshuffle between reconciles.
func (s *KubeVirtArchiverStatus) UpsertVirtualMachineStatus(vm VirtualMachineArchiverStatus) {
	for i := range s.VirtualMachines {
		if s.VirtualMachines[i].Name == vm.Name && s.VirtualMachines[i].Namespace == vm.Namespace {
			s.VirtualMachines[i] = vm
			return
		}
	}
	s.VirtualMachines = append(s.VirtualMachines, vm)
}

// RemoveVirtualMachineStatus drops a VirtualMachine that no longer matches the
// selector, so a stale entry cannot keep reporting a VM the archiver has released.
func (s *KubeVirtArchiverStatus) RemoveVirtualMachineStatus(namespace, name string) {
	kept := s.VirtualMachines[:0]
	for _, vm := range s.VirtualMachines {
		if vm.Name == name && vm.Namespace == namespace {
			continue
		}
		kept = append(kept, vm)
	}
	s.VirtualMachines = kept
}

// AllVirtualMachinesReady reports whether every matched VirtualMachine can actually
// be backed up. An archiver matching nothing is not ready: it is a selector typo.
func (s *KubeVirtArchiverStatus) AllVirtualMachinesReady() bool {
	if len(s.VirtualMachines) == 0 {
		return false
	}
	for _, vm := range s.VirtualMachines {
		if !vm.ChangedBlockTrackingReady {
			return false
		}
	}
	return true
}
