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
	"kmodules.xyz/client-go/apiextensions"
	cutil "kmodules.xyz/client-go/conditions"
	"kmodules.xyz/client-go/meta"
	"kubestash.dev/apimachinery/apis"
	"kubestash.dev/apimachinery/crds"
	"path/filepath"
	"regexp"
	"strings"
)

func (_ Snapshot) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(GroupVersion.WithResource(ResourcePluralSnapshot))
}

func (s *Snapshot) CalculatePhase() SnapshotPhase {
	if cutil.IsConditionFalse(s.Status.Conditions, TypeSnapshotMetadataUploaded) ||
		cutil.IsConditionFalse(s.Status.Conditions, TypeRecentSnapshotListUpdated) ||
		cutil.IsConditionTrue(s.Status.Conditions, TypeBackupIncomplete) {
		return SnapshotFailed
	}
	if s.GetComponentsPhase() == SnapshotPending {
		return SnapshotPending
	}
	if cutil.HasCondition(s.Status.Conditions, TypeSnapshotMetadataUploaded) {
		return s.GetComponentsPhase()
	}
	return SnapshotRunning
}

func (s *Snapshot) GetComponentsPhase() SnapshotPhase {
	if len(s.Status.Components) == 0 {
		return SnapshotPending
	}

	failedComponent := 0
	successfulComponent := 0

	for _, c := range s.Status.Components {
		if c.Phase == ComponentPhaseSucceeded {
			successfulComponent++
		}
		if c.Phase == ComponentPhaseFailed {
			failedComponent++
		}
	}

	totalComponents := int(s.Status.TotalComponents)

	if successfulComponent == totalComponents {
		return SnapshotSucceeded
	}

	if successfulComponent+failedComponent == totalComponents {
		return SnapshotFailed
	}

	return SnapshotRunning
}

func (s *Snapshot) IsCompleted() bool {
	return s.Status.Phase == SnapshotSucceeded || s.Status.Phase == SnapshotFailed
}

func (s *Snapshot) GetIntegrity() *bool {
	if s.Status.Components == nil {
		return nil
	}

	result, hasResticComp := true, false
	for _, component := range s.Status.Components {
		if component.ResticStats != nil &&
			component.Integrity == nil {
			return nil
		}

		if component.Integrity == nil {
			continue
		}

		hasResticComp = true
		result = result && *component.Integrity
	}

	if hasResticComp {
		return &result
	}
	return nil
}

func (s *Snapshot) GetTotalBackupSizeInBytes() (uint64, error) {
	if s.Status.Components == nil {
		return 0, fmt.Errorf("no component found for snapshot %s/%s", s.Namespace, s.Name)
	}

	var totalSizeInByte uint64
	for componentName, component := range s.Status.Components {
		for _, stats := range component.ResticStats {
			if stats.Size == "" {
				return 0, fmt.Errorf("resticStats size of component %s is empty for the snapshot %s/%s", componentName, s.Namespace, s.Name)
			}

			sizeWithUnit := strings.Split(component.Size, " ")
			if len(sizeWithUnit) < 2 {
				return 0, fmt.Errorf("resticStats size of component %s is invalid for the snapshot %s/%s", componentName, s.Namespace, s.Name)
			}

			sizeInByte, err := ConvertSizeToByte(sizeWithUnit)
			if err != nil {
				return 0, err
			}
			totalSizeInByte += sizeInByte
		}
	}
	return totalSizeInByte, nil
}

func (s *Snapshot) GetSize() string {
	if s.Status.Components == nil {
		return ""
	}

	var totalSizeInByte uint64
	hasResticComp := false
	for _, component := range s.Status.Components {
		if component.ResticStats != nil &&
			component.Size == "" {
			return ""
		}

		if component.Size == "" {
			continue
		}

		sizeWithUnit := strings.Split(component.Size, " ")
		if len(sizeWithUnit) < 2 {
			return ""
		}

		sizeInByte, err := ConvertSizeToByte(sizeWithUnit)
		if err != nil {
			return ""
		}
		hasResticComp = true
		totalSizeInByte += sizeInByte
	}
	if hasResticComp {
		return FormatBytes(totalSizeInByte)
	}

	return ""
}

func GenerateSnapshotName(repoName, backupSession string) string {
	backupSessionRegex := regexp.MustCompile("(.*)-([0-9]+)$")
	subMatches := backupSessionRegex.FindStringSubmatch(backupSession)
	return meta.ValidNameWithPrefixNSuffix(repoName, subMatches[1], subMatches[2])
}

func (s *Snapshot) OffshootLabels() map[string]string {
	newLabels := make(map[string]string)
	newLabels[apis.KubeStashInvokerKind] = ResourceKindSnapshot
	newLabels[meta.ManagedByLabelKey] = apis.KubeStashKey
	newLabels[apis.KubeStashInvokerName] = s.Name
	newLabels[apis.KubeStashInvokerNamespace] = s.Namespace
	return apis.UpsertLabels(s.Labels, newLabels)
}

func (s *Snapshot) GetComponentPath(componentName string) string {
	return filepath.Join(apis.DirRepository, s.Spec.Version, s.Spec.Session, componentName)
}
