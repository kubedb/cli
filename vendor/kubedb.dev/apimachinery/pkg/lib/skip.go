/*
Copyright AppsCode Inc. and Contributors

Licensed under the AppsCode Free Trial License 1.0.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://github.com/appscode/licenses/raw/1.0.0/AppsCode-Free-Trial-1.0.0.md

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
	"k8s.io/klog/v2"
	kmapi "kmodules.xyz/client-go/api/v1"
	cu "kmodules.xyz/client-go/client"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type skipper struct {
	kbClient   client.Client
	kind       string
	opsReqList []client.Object
	currentOps client.Object
}

func NewSkipper(kbClient client.Client, kind string, curOps client.Object, opsReqList []client.Object) skipper {
	return skipper{
		kbClient:   kbClient,
		kind:       kind,
		opsReqList: opsReqList,
		currentOps: curOps,
	}
}

/*
The idea behind skipping opsRequests :
1) If there are multiple opsReqs of same `Type` (like 3 'VerticalScaling') in Pending state,
   We don't want to reconcile them one-by-one. Because only reconciling the last one is enough, & that is user's desired spec.
   we are setting opsReq Phase `Skipped` in this situation, for all, except the last one.
2) If there are multiple opsReqs of different `Type` (like 3 'VerticalScaling', 2 `UpdateVersion`) in Pending state,
   After skipping in the previous step, there will be exactly one opsReq of each type in Pending
   And now, as they are different types, We want to reconcile the oldest one first.
3) Reconfigure ops requests are excluded from this skipper logic and are handled by ReconfigureMerger instead.
*/

func (s *skipper) SkipOpsReq() (bool, error) {
	// if it is `Progressing`, let it go on
	if s.currentOps.(opsapi.Accessor).GetStatus().Phase == opsapi.OpsRequestPhaseProgressing {
		return false, nil
	}

	// Skip the skipper logic for Reconfigure ops requests - they are handled by ReconfigureMerger
	currentOpsType := fmt.Sprintf("%v", s.currentOps.(opsapi.Accessor).GetRequestType())
	if isInExcludeList(currentOpsType) {
		return false, nil
	}

	// req Phase is must be `Pending`, as we already returned for other phases
	// opsMap will not be empty, as at least 'req' will be there
	opsMap := make(map[string][]client.Object)
	for _, o := range s.opsReqList {
		r := o.(opsapi.Accessor)
		// all the opsReqs that refer the same db as 'req' & currently in `Pending` or "" Phase , add them in the map
		// if Phase is `Progressing`, continue to work with that
		// Do nothing if phase is `Successful`, `Failed` or `Skipped`
		if r.GetDBRefName() != s.currentOps.(opsapi.Accessor).GetDBRefName() {
			continue
		}
		if r.GetStatus().Phase == opsapi.OpsRequestPhaseProgressing {
			return true, nil
		}
		// r.Status.Phase can be "", if it has not been reconciled yet.
		if r.GetStatus().Phase == opsapi.OpsRequestPhasePending || r.GetStatus().Phase == "" {
			opsType := fmt.Sprintf("%v", r.GetRequestType()) // r.GetRequestType().(string) causes panic
			if isInExcludeList(opsType) {
				continue
			}
			// populate the map
			opsMap[opsType] = append(opsMap[opsType], o)
		}
	}

	oldestOps, err := s.getOldestOpsAfterSkipping(opsMap)
	if err != nil {
		return true, err
	}

	return oldestOps.Name != s.currentOps.GetName(), nil
}

var ExcludedFromSkipperLogic = []string{"Reconfigure"} // Skip Reconfigure ops requests - they are handled by ReconfigureMerger

func isInExcludeList(s string) bool {
	for _, ex := range ExcludedFromSkipperLogic {
		if s == ex {
			return true
		}
	}
	return false
}

func (s *skipper) getOldestOpsAfterSkipping(opsMap map[string][]client.Object) (metav1.ObjectMeta, error) {
	oldestOps := metav1.ObjectMeta{
		CreationTimestamp: metav1.Now(), // set dummy time for later comparison
	}

	// sort the map values with
	// 1) creation time (the oldest first)
	// 2) opsReq name length (the lexicographically bigger first)
	for typ, details := range opsMap {
		sort.Slice(details, func(i, j int) bool {
			iMeta := details[i].(opsapi.Accessor).GetObjectMeta()
			jMeta := details[j].(opsapi.Accessor).GetObjectMeta()
			if !iMeta.CreationTimestamp.Equal(&jMeta.CreationTimestamp) {
				return iMeta.CreationTimestamp.Before(&jMeta.CreationTimestamp)
			}
			return strings.Compare(iMeta.Name, jMeta.Name) > 0
		})

		last := details[len(details)-1].(opsapi.Accessor).GetObjectMeta()
		msg := fmt.Sprintf("skipped as %v/%v is the most updated %v of type %v at %v",
			last.Namespace, last.Name, s.kind, typ, time.Now().String())

		// all the opsReqs in 'details' are of same OpsType & for same db.
		// so make all the opsReqPhase `Skipped` other than the last one
		for i := 0; i < len(details)-1; i++ {
			if err := s.markAsSkipped(details[i], msg); err != nil {
				return oldestOps, err
			}
		}

		// just working with the last one as others are marked as 'Skipped'
		if last.CreationTimestamp.Before(&oldestOps.CreationTimestamp) {
			oldestOps = last
		}
	}
	return oldestOps, nil
}

func (s *skipper) markAsSkipped(details client.Object, msg string) error {
	_, err := cu.PatchStatus(context.TODO(), s.kbClient, details, func(obj client.Object) client.Object {
		accessor := obj.(opsapi.Accessor)
		in := accessor.GetStatus()

		in.Phase = opsapi.OpsRequestPhaseSkipped
		in.Conditions = append(in.Conditions, kmapi.Condition{
			Type:               kmapi.ConditionType(opsapi.OpsRequestPhaseSkipped),
			Status:             metav1.ConditionTrue,
			ObservedGeneration: in.ObservedGeneration,
			LastTransitionTime: metav1.Now(),
			Reason:             string(opsapi.OpsRequestPhaseSkipped),
			Message:            msg,
		})

		accessor.SetStatus(in)
		return obj
	})
	if err != nil {
		klog.Error(err, fmt.Sprintf("failed to update opsRequest %v/%v status", details.GetNamespace(), details.GetName()))
		return err
	}
	klog.Infof("%v/%v %s", details.GetNamespace(), details.GetName(), msg)
	return nil
}
