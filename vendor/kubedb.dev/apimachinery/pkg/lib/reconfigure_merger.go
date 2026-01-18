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

package lib

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	opsapi "kubedb.dev/apimachinery/apis/ops/v1alpha1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/klog/v2"
	kutil "kmodules.xyz/client-go"
	kmapi "kmodules.xyz/client-go/api/v1"
	cu "kmodules.xyz/client-go/client"
	cutil "kmodules.xyz/client-go/conditions"
	meta_util "kmodules.xyz/client-go/meta"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type MergeFunc func(pendingReconfigureOps []opsapi.Accessor) (any, error)

type ConvertFunc func(u *unstructured.Unstructured) (opsapi.Accessor, error)

type ReconfigureMerger struct {
	kbClient   client.Client
	kind       string
	opsReqList []client.Object
	currentOps opsapi.Accessor
	mergeFn    MergeFunc
	convertFn  ConvertFunc
	log        interface {
		Info(msg string, keysAndValues ...any)
	}
}

func NewReconfigureMerger(kbClient client.Client, kind string, curOps client.Object, mergeFn MergeFunc, convertFn ConvertFunc, log interface {
	Info(msg string, keysAndValues ...any)
},
) (*ReconfigureMerger, error) {
	ret := &ReconfigureMerger{
		kbClient:   kbClient,
		kind:       kind,
		currentOps: curOps.(opsapi.Accessor),
		mergeFn:    mergeFn,
		convertFn:  convertFn,
		log:        log,
	}
	err := ret.populateList()
	return ret, err
}

const (
	MergedOpsSubStr     = "-rcfg-merged-"
	OriginalOpsSkipped  = "OriginalOpsSkipped"
	MergedFromOps       = "MergedFromOps"
	ConfigurationMerged = "ConfigurationMerged"
)

const (
	ContinueGeneral  = iota // 0
	MergeNeeded             // 1
	RequeueNeeded           // 2
	RequeueNotNeeded        // 3
)

func (m *ReconfigureMerger) populateList() error {
	unsList := &unstructured.UnstructuredList{}
	unsList.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "ops.kubedb.com",
		Version: "v1alpha1",
		Kind:    m.kind + "List",
	})

	err := m.kbClient.List(context.TODO(), unsList,
		client.InNamespace(m.currentOps.GetNamespace()),
	)
	if err != nil {
		return err
	}

	var (
		lst      []client.Object
		accessor opsapi.Accessor
	)
	for _, request := range unsList.Items {
		accessor, err = m.convertFn(&request)
		if err != nil {
			return err
		}
		lst = append(lst, accessor)
	}
	m.opsReqList = lst
	return nil
}

func (m *ReconfigureMerger) Run() (int, error) {
	skip, pendingReconfigureOps := m.FindPendingReconfigureOpsToMerge()
	if skip != MergeNeeded {
		return skip, nil
	}

	mergedConfig, err := m.mergeFn(pendingReconfigureOps)
	if err != nil {
		return RequeueNeeded, fmt.Errorf("failed to merge configurations from reconfigure ops requests: %w", err)
	}

	u, err := m.GetMergedOpsRequest(pendingReconfigureOps[0], mergedConfig)
	if err != nil {
		return RequeueNeeded, err
	}

	mergedOps, err := m.convertFn(u)
	if err != nil {
		return RequeueNeeded, err
	}

	klog.Infof("mmm %v %v \n", m.currentOps.GetName(), mergedOps.GetName())
	err = m.EnsureMergedOpsRequest(mergedOps, pendingReconfigureOps)
	return ContinueGeneral, err
}

/*
FindPendingReconfigureOpsToMerge
- Check if OriginalOpsSkipped condition is present.
- If any Reconfigure is Progressing (other than current), requeue current request
- Collect pendingReconfigureOps
- Avoid parallel processing by checking if current ops is already part of an existing merge
- markAsSkipped is the current ops is already merged into another ops.
*/
func (m *ReconfigureMerger) FindPendingReconfigureOpsToMerge() (int, []opsapi.Accessor) {
	// Only process Reconfigure type ops requests
	if m.currentOps.GetRequestType() != opsapi.Reconfigure {
		return ContinueGeneral, nil
	}

	// If it is already Progressing, let it continue
	if m.currentOps.GetStatus().Phase == opsapi.OpsRequestPhaseProgressing {
		return ContinueGeneral, nil
	}

	// If current ops is a merged ops request, check if original ops have been skipped
	// Merged ops requests have names like: <dbname>-rcfg-merged-<timestamp>
	meta := m.currentOps.GetObjectMeta()
	if strings.Contains(meta.GetName(), MergedOpsSubStr) {
		// Check if the "OriginalOpsSkipped" condition is true
		if !m.areOriginalOpsSkipped(m.currentOps.GetStatus()) {
			m.log.Info(fmt.Sprintf("Merged ops request %s/%s waiting for original ops to be skipped, requeuing",
				meta.GetName(), meta.GetName()))
			return RequeueNeeded, nil // Requeue until original ops are skipped
		}
		// Original ops are skipped, proceed with reconciliation
		m.log.Info(fmt.Sprintf("Merged ops request %s/%s ready to proceed, original ops have been skipped",
			meta.GetNamespace(), meta.GetName()))
		return ContinueGeneral, nil
	}

	// Collect all pending Reconfigure ops requests for the same database
	var pendingReconfigureOps []opsapi.Accessor
	for _, o := range m.opsReqList {
		req := o.(opsapi.Accessor)

		// Only consider ops for the same database
		if req.GetDBRefName() != m.currentOps.GetDBRefName() {
			continue
		}

		// Only Reconfigure type
		if req.GetRequestType() != opsapi.Reconfigure {
			continue
		}

		// If any Reconfigure is Progressing (other than current), requeue current request
		if req.GetStatus().Phase == opsapi.OpsRequestPhaseProgressing {
			m.log.Info(fmt.Sprintf("Reconfigure ops request %s/%s is already progressing for database %s, requeuing current request %s/%s",
				req.GetObjectMeta().Namespace, req.GetObjectMeta().Name, req.GetDBRefName(), meta.GetNamespace(), meta.GetName()))
			return RequeueNeeded, nil
		}

		// Collect only Pending or empty ("") phase ops requests
		// Ignore Successful/Failed/Skipped
		if req.GetStatus().Phase == opsapi.OpsRequestPhasePending || req.GetStatus().Phase == "" {
			pendingReconfigureOps = append(pendingReconfigureOps, req)
		}
	}

	// If only one or none, no need to merge
	if len(pendingReconfigureOps) <= 1 {
		return ContinueGeneral, nil
	}

	// Check if current ops is already part of an existing merge
	// This prevents parallel processing from creating duplicate merges
	if alreadyMerged, mergedOpsName := isAlreadyMerged(m.currentOps, pendingReconfigureOps); alreadyMerged {
		m.log.Info(fmt.Sprintf("Ops request %s/%s is already merged into %s, skipping",
			meta.GetNamespace(), meta.GetName(), mergedOpsName))

		// Mark current ops as Skipped since it's already merged
		if err := m.markAsSkippedForMergedOps(m.currentOps, metav1.ObjectMeta{
			Name:      mergedOpsName,
			Namespace: meta.GetNamespace(),
		}); err != nil {
			klog.Errorf("failed to mark already merged ops request %s/%s as skipped: %v", meta.GetNamespace(), meta.GetName(), err)
			return RequeueNotNeeded, nil
		}
		return RequeueNotNeeded, nil // Skip this ops as it's already part of a merge
	}

	m.log.Info(fmt.Sprintf("Found %d pending reconfigure ops requests to merge", len(pendingReconfigureOps)))

	// Sort by creation timestamp (oldest first), then by name (lexicographically smaller first)
	sort.Slice(pendingReconfigureOps, func(i, j int) bool {
		x := pendingReconfigureOps[i].GetObjectMeta().CreationTimestamp
		y := pendingReconfigureOps[j].GetObjectMeta()
		if !x.Equal(&y.CreationTimestamp) {
			return x.Before(&y.CreationTimestamp)
		}
		return strings.Compare(pendingReconfigureOps[i].GetObjectMeta().Name, pendingReconfigureOps[j].GetObjectMeta().Name) < 0
	})

	return MergeNeeded, pendingReconfigureOps
}

func (m *ReconfigureMerger) areOriginalOpsSkipped(opsStatus opsapi.OpsRequestStatus) bool {
	return cutil.IsConditionTrue(opsStatus.Conditions, OriginalOpsSkipped)
}

func isAlreadyMerged(currentReq opsapi.Accessor, pendingReconfigureOps []opsapi.Accessor) (bool, string) {
	for _, req := range pendingReconfigureOps {
		// Only check ops requests that is merged ops
		if !strings.Contains(req.GetObjectMeta().Name, MergedOpsSubStr) {
			continue
		}

		// Check if current ops name is in the "MergedFromOps" condition message
		for _, cond := range req.GetStatus().Conditions {
			if cond.Type == MergedFromOps && cond.Status == metav1.ConditionTrue {
				// The message contains comma-separated list of original ops names
				opsNames := strings.Split(cond.Message, ", ")
				for _, name := range opsNames {
					if strings.TrimSpace(name) == currentReq.GetObjectMeta().Name {
						return true, req.GetObjectMeta().Name
					}
				}
			}
		}
	}
	return false, ""
}

func (m *ReconfigureMerger) markAsSkippedForMergedOps(ops opsapi.Accessor, mergedOpsMeta metav1.ObjectMeta) error {
	msg := fmt.Sprintf("Configuration merged into newly created ops request %s/%s", mergedOpsMeta.Namespace, mergedOpsMeta.Name)

	_, err := cu.PatchStatus(context.TODO(), m.kbClient, ops, func(obj client.Object) client.Object {
		ret := obj.(opsapi.Accessor)
		sts := ret.GetStatus()
		sts.Phase = opsapi.OpsRequestPhaseSkipped
		sts.ObservedGeneration = ret.GetObjectMeta().Generation
		sts.Conditions = cutil.SetCondition(sts.Conditions, kmapi.Condition{
			Type:               kmapi.ConditionType(opsapi.OpsRequestPhaseSkipped),
			Reason:             ConfigurationMerged,
			Status:             metav1.ConditionTrue,
			ObservedGeneration: ret.GetObjectMeta().Generation,
			LastTransitionTime: metav1.Now(),
			Message:            msg,
		})
		ret.SetStatus(sts)
		return ret
	})
	if err != nil {
		return fmt.Errorf("failed to mark ops request %s/%s as skipped: %w", ops.GetObjectMeta().Namespace, ops.GetObjectMeta().Name, err)
	}

	m.log.Info(fmt.Sprintf("%s/%s skipped: %s", ops.GetObjectMeta().Namespace, ops.GetObjectMeta().Name, msg))
	return nil
}

func (m *ReconfigureMerger) GetMergedOpsRequest(firstPendingOps opsapi.Accessor, config any) (*unstructured.Unstructured, error) {
	toUnstructured := func(obj any) (*unstructured.Unstructured, error) {
		m, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
		if err != nil {
			return nil, err
		}
		return &unstructured.Unstructured{Object: m}, nil
	}

	buildReconfigureOpsRequest := func(
		firstPendingOps opsapi.Accessor,
		config any,
	) (*unstructured.Unstructured, error) {
		mergedOpsName := meta_util.NameWithSuffix(firstPendingOps.GetDBRefName(), fmt.Sprintf("rcfg-merged-%d", time.Now().UnixMilli()))
		cfg, err := toUnstructured(config)
		if err != nil {
			return nil, err
		}

		obj := map[string]any{
			"apiVersion": "ops.kubedb.com/v1alpha1",
			"kind":       m.kind,
			"metadata": map[string]any{
				"name":      mergedOpsName,
				"namespace": firstPendingOps.GetObjectMeta().Namespace,
				"labels":    firstPendingOps.GetObjectMeta().Labels,
			},
			"spec": map[string]any{
				"type": opsapi.Reconfigure,
				"databaseRef": map[string]any{
					"name": firstPendingOps.GetDBRefName(),
				},
				"configuration": cfg.Object,
			},
		}

		return &unstructured.Unstructured{Object: obj}, nil
	}

	request, err := buildReconfigureOpsRequest(firstPendingOps, config)
	if err != nil {
		return nil, err
	}

	return request, nil
}

func (m *ReconfigureMerger) EnsureMergedOpsRequest(mergedOpsRequest opsapi.Accessor, pendingReconfigureOps []opsapi.Accessor) error {
	verb, err := cu.CreateOrPatch(context.TODO(), m.kbClient, mergedOpsRequest, func(obj client.Object, _ bool) client.Object {
		return mergedOpsRequest
	})
	if err != nil {
		return fmt.Errorf("failed to create merged ops request %s/%s: %w", mergedOpsRequest.GetObjectMeta().Namespace, mergedOpsRequest.GetObjectMeta().Name, err)
	}

	// Prepare the list of original ops names for the "MergedFromOps" condition
	var opsNames []string
	for _, ops := range pendingReconfigureOps {
		opsNames = append(opsNames, ops.GetObjectMeta().Name)
	}

	mergedFromOpsMessage := strings.Join(opsNames, ", ")
	if verb == kutil.VerbCreated {
		m.log.Info(fmt.Sprintf("Created merged reconfigure ops request %s/%s from %d pending requests",
			mergedOpsRequest.GetObjectMeta().Namespace, mergedOpsRequest.GetObjectMeta().Name, len(opsNames)))
	}

	if err := m.updateConditionForMergedOps(mergedOpsRequest, mergedFromOpsMessage); err != nil {
		return fmt.Errorf("failed to record merged-from condition for %s/%s: %w", mergedOpsRequest.GetObjectMeta().Namespace, mergedOpsRequest.GetObjectMeta().Name, err)
	}

	for _, op := range pendingReconfigureOps {
		err = m.markAsSkippedForMergedOps(op, mergedOpsRequest.GetObjectMeta())
		if err != nil {
			return err
		}
	}

	// Mark the merged ops request with "OriginalOpsSkipped" condition after skipping is done
	// This allows the merged ops to proceed with reconciliation
	if err := m.markOriginalOpsAsSkipped(mergedOpsRequest); err != nil {
		return fmt.Errorf("failed to mark original ops as skipped for merged ops %s/%s: %w", mergedOpsRequest.GetObjectMeta().Namespace, mergedOpsRequest.GetObjectMeta().Name, err)
	}

	m.log.Info(fmt.Sprintf("Marked merged ops request %s/%s as ready to proceed (original ops skipped)", mergedOpsRequest.GetObjectMeta().Namespace, mergedOpsRequest.GetObjectMeta().Name))
	return nil
}

func (m *ReconfigureMerger) updateConditionForMergedOps(mergedOps opsapi.Accessor, mergedFromOpsMessage string) error {
	_, err := cu.PatchStatus(context.TODO(), m.kbClient, mergedOps, func(obj client.Object) client.Object {
		ret := obj.(opsapi.Accessor)
		sts := ret.GetStatus()
		sts.ObservedGeneration = ret.GetGeneration()
		sts.Conditions = cutil.SetCondition(sts.Conditions, kmapi.Condition{
			Type:               MergedFromOps,
			Reason:             ConfigurationMerged,
			Status:             metav1.ConditionTrue,
			ObservedGeneration: ret.GetGeneration(),
			LastTransitionTime: metav1.Now(),
			Message:            mergedFromOpsMessage,
		})
		ret.SetStatus(sts)
		return ret
	})
	if err != nil {
		return fmt.Errorf("failed to set MergedFromOps condition for %s/%s: %w", mergedOps.GetObjectMeta().Namespace, mergedOps.GetObjectMeta().Name, err)
	}
	return nil
}

func (m *ReconfigureMerger) markOriginalOpsAsSkipped(mergedOps opsapi.Accessor) error {
	_, err := cu.PatchStatus(context.TODO(), m.kbClient, mergedOps, func(obj client.Object) client.Object {
		ret := obj.(opsapi.Accessor)
		sts := ret.GetStatus()
		sts.ObservedGeneration = ret.GetGeneration()
		sts.Conditions = cutil.SetCondition(sts.Conditions, kmapi.Condition{
			Type:               OriginalOpsSkipped,
			Reason:             OriginalOpsSkipped,
			Status:             metav1.ConditionTrue,
			ObservedGeneration: ret.GetGeneration(),
			LastTransitionTime: metav1.Now(),
			Message:            "All original ops requests have been marked as Skipped",
		})
		ret.SetStatus(sts)
		return ret
	})
	if err != nil {
		return fmt.Errorf("failed to mark original ops as skipped for merged ops %s/%s: %w", mergedOps.GetObjectMeta().Namespace, mergedOps.GetObjectMeta().Name, err)
	}
	return nil
}
