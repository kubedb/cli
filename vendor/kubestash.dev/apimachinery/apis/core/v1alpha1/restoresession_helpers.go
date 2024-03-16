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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kmapi "kmodules.xyz/client-go/api/v1"
	"kubestash.dev/apimachinery/apis"
	"kubestash.dev/apimachinery/crds"

	"kmodules.xyz/client-go/apiextensions"
	cutil "kmodules.xyz/client-go/conditions"
	meta_util "kmodules.xyz/client-go/meta"
)

func (_ RestoreSession) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crds.MustCustomResourceDefinition(GroupVersion.WithResource(ResourcePluralRestoreSession))
}

func (rs *RestoreSession) CalculatePhase() RestorePhase {
	if cutil.IsConditionFalse(rs.Status.Conditions, TypeMetricsPushed) {
		return RestoreFailed
	}

	if cutil.IsConditionFalse(rs.Status.Conditions, TypeValidationPassed) {
		return RestoreInvalid
	}

	if cutil.IsConditionTrue(rs.Status.Conditions, TypeMetricsPushed) &&
		(cutil.IsConditionTrue(rs.Status.Conditions, TypeDeadlineExceeded) ||
			cutil.IsConditionFalse(rs.Status.Conditions, TypePreRestoreHooksExecutionSucceeded) ||
			cutil.IsConditionFalse(rs.Status.Conditions, TypePostRestoreHooksExecutionSucceeded) ||
			cutil.IsConditionFalse(rs.Status.Conditions, TypeRestoreExecutorEnsured)) {
		return RestoreFailed
	}

	componentsPhase := rs.getComponentsPhase()
	if componentsPhase == RestorePending || rs.FinalStepExecuted() {
		return componentsPhase
	}

	if componentsPhase == RestorePhaseUnknown {
		return componentsPhase
	}

	return RestoreRunning
}

func (rs *RestoreSession) AllComponentsCompleted() bool {
	phase := rs.getComponentsPhase()
	return phase == RestoreSucceeded || phase == RestoreFailed
}

func (rs *RestoreSession) FinalStepExecuted() bool {
	return cutil.HasCondition(rs.Status.Conditions, TypeMetricsPushed)
}

func (rs *RestoreSession) getComponentsPhase() RestorePhase {
	if len(rs.Status.Components) == 0 {
		return RestorePending
	}

	failedComponent := 0
	successfulComponent := 0
	unknownComponentPhase := 0

	for _, c := range rs.Status.Components {
		if c.Phase == RestoreSucceeded {
			successfulComponent++
		}
		if c.Phase == RestoreFailed {
			failedComponent++
		}
		if c.Phase == RestorePhaseUnknown {
			unknownComponentPhase++
		}
	}

	totalComponents := int(rs.Status.TotalComponents)

	if successfulComponent == totalComponents {
		return RestoreSucceeded
	}

	if successfulComponent+failedComponent+unknownComponentPhase == totalComponents {
		if failedComponent > 0 {
			return RestoreFailed
		}
		return RestorePhaseUnknown
	}

	return RestoreRunning
}

func (rs *RestoreSession) OffshootLabels() map[string]string {
	newLabels := make(map[string]string)
	newLabels[meta_util.ManagedByLabelKey] = apis.KubeStashKey
	newLabels[apis.KubeStashInvokerName] = rs.Name
	newLabels[apis.KubeStashInvokerNamespace] = rs.Namespace

	return apis.UpsertLabels(rs.Labels, newLabels)
}

func (rs *RestoreSession) GetSummary(targetRef *kmapi.TypedObjectReference) *Summary {
	errMsg := rs.getFailureMessage()
	phase := RestoreSucceeded
	if errMsg != "" {
		phase = RestoreFailed
	}

	return &Summary{
		Name:      rs.Name,
		Namespace: rs.Namespace,

		Invoker: &kmapi.TypedObjectReference{
			APIGroup:  GroupVersion.Group,
			Kind:      rs.Kind,
			Name:      rs.Name,
			Namespace: rs.Namespace,
		},

		Target: targetRef,

		Status: TargetStatus{
			Phase:    string(phase),
			Duration: rs.Status.Duration,
			Error:    errMsg,
		},
	}
}

func (rs *RestoreSession) getFailureMessage() string {
	failureFound, reason := rs.checkFailureInConditions()
	if failureFound {
		return reason
	}
	failureFound, reason = rs.checkFailureInComponents()
	if failureFound {
		return reason
	}
	return ""
}

func (rs *RestoreSession) checkFailureInConditions() (bool, string) {
	for _, condition := range rs.Status.Conditions {
		if condition.Status == metav1.ConditionFalse {
			return true, condition.Message
		}
	}

	return false, ""
}

func (rs *RestoreSession) checkFailureInComponents() (bool, string) {
	for _, comp := range rs.Status.Components {
		if comp.Phase == RestoreFailed {
			return true, comp.Error
		}
	}

	return false, ""
}

func (rs *RestoreSession) GetDataSourceNamespace() string {
	if rs.Spec.DataSource.Namespace == "" {
		return rs.Namespace
	}
	return rs.Spec.DataSource.Namespace
}
