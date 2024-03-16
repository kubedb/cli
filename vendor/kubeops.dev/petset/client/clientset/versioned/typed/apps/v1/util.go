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

package v1

import (
	"context"
	"fmt"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"gomodules.xyz/pointer"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"
	kutil "kmodules.xyz/client-go"
	core_util "kmodules.xyz/client-go/core/v1"
	"kubeops.dev/petset/apis/apps/v1"
)

var json = jsoniter.ConfigFastest

func CreateOrPatchPetSet(ctx context.Context, psc AppsV1Interface, meta metav1.ObjectMeta, transform func(*v1.PetSet) *v1.PetSet, opts metav1.PatchOptions) (*v1.PetSet, kutil.VerbType, error) {
	cur, err := psc.PetSets(meta.Namespace).Get(ctx, meta.Name, metav1.GetOptions{})
	if kerr.IsNotFound(err) {
		klog.V(3).Infof("Creating PetSet %s/%s.", meta.Namespace, meta.Name)
		out, err := psc.PetSets(meta.Namespace).Create(ctx, transform(&v1.PetSet{
			TypeMeta: metav1.TypeMeta{
				Kind:       "PetSet",
				APIVersion: v1.SchemeGroupVersion.String(),
			},
			ObjectMeta: meta,
		}), metav1.CreateOptions{
			DryRun:       opts.DryRun,
			FieldManager: opts.FieldManager,
		})
		return out, kutil.VerbCreated, err
	} else if err != nil {
		return nil, kutil.VerbUnchanged, err
	}
	return PatchPetSet(ctx, psc, cur, transform, opts)
}

func PatchPetSet(ctx context.Context, psc AppsV1Interface, cur *v1.PetSet, transform func(*v1.PetSet) *v1.PetSet, opts metav1.PatchOptions) (*v1.PetSet, kutil.VerbType, error) {
	return PatchPetSetObject(ctx, psc, cur, transform(cur.DeepCopy()), opts)
}

func PatchPetSetObject(ctx context.Context, psc AppsV1Interface, cur, mod *v1.PetSet, opts metav1.PatchOptions) (*v1.PetSet, kutil.VerbType, error) {
	curJson, err := json.Marshal(cur)
	if err != nil {
		return nil, kutil.VerbUnchanged, err
	}

	modJson, err := json.Marshal(mod)
	if err != nil {
		return nil, kutil.VerbUnchanged, err
	}

	patch, err := strategicpatch.CreateTwoWayMergePatch(curJson, modJson, v1.PetSet{})
	if err != nil {
		return nil, kutil.VerbUnchanged, err
	}
	if len(patch) == 0 || string(patch) == "{}" {
		return cur, kutil.VerbUnchanged, nil
	}
	klog.V(3).Infof("Patching PetSet %s/%s with %s.", cur.Namespace, cur.Name, string(patch))
	out, err := psc.PetSets(cur.Namespace).Patch(ctx, cur.Name, types.MergePatchType, patch, opts)
	return out, kutil.VerbPatched, err
}

func TryUpdatePetSet(ctx context.Context, psc AppsV1Interface, meta metav1.ObjectMeta, transform func(*v1.PetSet) *v1.PetSet, opts metav1.UpdateOptions) (result *v1.PetSet, err error) {
	attempt := 0
	err = wait.PollUntilContextTimeout(ctx, kutil.RetryInterval, kutil.RetryTimeout, true, func(ctx context.Context) (bool, error) {
		attempt++
		cur, e2 := psc.PetSets(meta.Namespace).Get(ctx, meta.Name, metav1.GetOptions{})
		if kerr.IsNotFound(e2) {
			return false, e2
		} else if e2 == nil {
			result, e2 = psc.PetSets(cur.Namespace).Update(ctx, transform(cur.DeepCopy()), opts)
			return e2 == nil, nil
		}
		klog.Errorf("Attempt %d failed to update PetSet %s/%s due to %v.", attempt, cur.Namespace, cur.Name, e2)
		return false, nil
	})

	if err != nil {
		err = errors.Errorf("failed to update PetSet %s/%s after %d attempts due to %v", meta.Namespace, meta.Name, attempt, err)
	}
	return
}

func IsPetSetReady(obj *v1.PetSet) bool {
	replicas := int32(1)
	if obj.Spec.Replicas != nil {
		replicas = *obj.Spec.Replicas
	}
	return replicas == obj.Status.ReadyReplicas
}

func PetSetsAreReady(items []*v1.PetSet) (bool, string) {
	for _, ps := range items {
		if !IsPetSetReady(ps) {
			return false, fmt.Sprintf("All desired replicas are not ready. For PetSet: %s/%s desired replicas: %d, ready replicas: %d.", ps.Namespace, ps.Name, pointer.Int32(ps.Spec.Replicas), ps.Status.ReadyReplicas)
		}
	}
	return true, "All desired replicas are ready."
}

func WaitUntilPetSetReady(ctx context.Context, psc AppsV1Interface, meta metav1.ObjectMeta) error {
	return wait.PollUntilContextTimeout(ctx, kutil.RetryInterval, kutil.ReadinessTimeout, true, func(ctx context.Context) (bool, error) {
		if obj, err := psc.PetSets(meta.Namespace).Get(ctx, meta.Name, metav1.GetOptions{}); err == nil {
			return IsPetSetReady(obj), nil
		}
		return false, nil
	})
}

func DeletePetSet(ctx context.Context, c kubernetes.Interface, psc AppsV1Interface, meta metav1.ObjectMeta) error {
	petSet, err := psc.PetSets(meta.Namespace).Get(ctx, meta.Name, metav1.GetOptions{})
	if err != nil {
		if kerr.IsNotFound(err) {
			return nil
		} else {
			return err
		}
	}

	// Update PetSet
	_, _, err = PatchPetSet(ctx, psc, petSet, func(in *v1.PetSet) *v1.PetSet {
		in.Spec.Replicas = pointer.Int32P(0)
		return in
	}, metav1.PatchOptions{})
	if err != nil {
		return err
	}

	err = core_util.WaitUntilPodDeletedBySelector(ctx, c, petSet.Namespace, petSet.Spec.Selector)
	if err != nil {
		return err
	}

	return psc.PetSets(petSet.Namespace).Delete(ctx, petSet.Name, metav1.DeleteOptions{})
}
