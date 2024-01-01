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

	core_util "kmodules.xyz/client-go/core/v1"

	"github.com/pkg/errors"
	"gomodules.xyz/pointer"
	apps "k8s.io/api/apps/v1"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"
	kutil "kmodules.xyz/client-go"
)

func CreateOrPatchStatefulSet(ctx context.Context, c kubernetes.Interface, meta metav1.ObjectMeta, transform func(*apps.StatefulSet) *apps.StatefulSet, opts metav1.PatchOptions) (*apps.StatefulSet, kutil.VerbType, error) {
	cur, err := c.AppsV1().StatefulSets(meta.Namespace).Get(ctx, meta.Name, metav1.GetOptions{})
	if kerr.IsNotFound(err) {
		klog.V(3).Infof("Creating StatefulSet %s/%s.", meta.Namespace, meta.Name)
		out, err := c.AppsV1().StatefulSets(meta.Namespace).Create(ctx, transform(&apps.StatefulSet{
			TypeMeta: metav1.TypeMeta{
				Kind:       "StatefulSet",
				APIVersion: apps.SchemeGroupVersion.String(),
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
	return PatchStatefulSet(ctx, c, cur, transform, opts)
}

func PatchStatefulSet(ctx context.Context, c kubernetes.Interface, cur *apps.StatefulSet, transform func(*apps.StatefulSet) *apps.StatefulSet, opts metav1.PatchOptions) (*apps.StatefulSet, kutil.VerbType, error) {
	return PatchStatefulSetObject(ctx, c, cur, transform(cur.DeepCopy()), opts)
}

func PatchStatefulSetObject(ctx context.Context, c kubernetes.Interface, cur, mod *apps.StatefulSet, opts metav1.PatchOptions) (*apps.StatefulSet, kutil.VerbType, error) {
	curJson, err := json.Marshal(cur)
	if err != nil {
		return nil, kutil.VerbUnchanged, err
	}

	modJson, err := json.Marshal(mod)
	if err != nil {
		return nil, kutil.VerbUnchanged, err
	}

	patch, err := strategicpatch.CreateTwoWayMergePatch(curJson, modJson, apps.StatefulSet{})
	if err != nil {
		return nil, kutil.VerbUnchanged, err
	}
	if len(patch) == 0 || string(patch) == "{}" {
		return cur, kutil.VerbUnchanged, nil
	}
	klog.V(3).Infof("Patching StatefulSet %s/%s with %s.", cur.Namespace, cur.Name, string(patch))
	out, err := c.AppsV1().StatefulSets(cur.Namespace).Patch(ctx, cur.Name, types.StrategicMergePatchType, patch, opts)
	return out, kutil.VerbPatched, err
}

func TryUpdateStatefulSet(ctx context.Context, c kubernetes.Interface, meta metav1.ObjectMeta, transform func(*apps.StatefulSet) *apps.StatefulSet, opts metav1.UpdateOptions) (result *apps.StatefulSet, err error) {
	attempt := 0
	err = wait.PollUntilContextTimeout(ctx, kutil.RetryInterval, kutil.RetryTimeout, true, func(ctx context.Context) (bool, error) {
		attempt++
		cur, e2 := c.AppsV1().StatefulSets(meta.Namespace).Get(ctx, meta.Name, metav1.GetOptions{})
		if kerr.IsNotFound(e2) {
			return false, e2
		} else if e2 == nil {
			result, e2 = c.AppsV1().StatefulSets(cur.Namespace).Update(ctx, transform(cur.DeepCopy()), opts)
			return e2 == nil, nil
		}
		klog.Errorf("Attempt %d failed to update StatefulSet %s/%s due to %v.", attempt, cur.Namespace, cur.Name, e2)
		return false, nil
	})

	if err != nil {
		err = errors.Errorf("failed to update StatefulSet %s/%s after %d attempts due to %v", meta.Namespace, meta.Name, attempt, err)
	}
	return
}

func IsStatefulSetReady(obj *apps.StatefulSet) bool {
	replicas := int32(1)
	if obj.Spec.Replicas != nil {
		replicas = *obj.Spec.Replicas
	}
	return replicas == obj.Status.ReadyReplicas
}

func StatefulSetsAreReady(items []*apps.StatefulSet) (bool, string) {
	for _, sts := range items {
		if !IsStatefulSetReady(sts) {
			return false, fmt.Sprintf("All desired replicas are not ready. For StatefulSet: %s/%s desired replicas: %d, ready replicas: %d.", sts.Namespace, sts.Name, pointer.Int32(sts.Spec.Replicas), sts.Status.ReadyReplicas)
		}
	}
	return true, "All desired replicas are ready."
}

func WaitUntilStatefulSetReady(ctx context.Context, c kubernetes.Interface, meta metav1.ObjectMeta) error {
	return wait.PollUntilContextTimeout(ctx, kutil.RetryInterval, kutil.ReadinessTimeout, true, func(ctx context.Context) (bool, error) {
		if obj, err := c.AppsV1().StatefulSets(meta.Namespace).Get(ctx, meta.Name, metav1.GetOptions{}); err == nil {
			return IsStatefulSetReady(obj), nil
		}
		return false, nil
	})
}

func DeleteStatefulSet(ctx context.Context, c kubernetes.Interface, meta metav1.ObjectMeta) error {
	statefulSet, err := c.AppsV1().StatefulSets(meta.Namespace).Get(ctx, meta.Name, metav1.GetOptions{})
	if err != nil {
		if kerr.IsNotFound(err) {
			return nil
		} else {
			return err
		}
	}

	// Update StatefulSet
	_, _, err = PatchStatefulSet(ctx, c, statefulSet, func(in *apps.StatefulSet) *apps.StatefulSet {
		in.Spec.Replicas = pointer.Int32P(0)
		return in
	}, metav1.PatchOptions{})
	if err != nil {
		return err
	}

	err = core_util.WaitUntilPodDeletedBySelector(ctx, c, statefulSet.Namespace, statefulSet.Spec.Selector)
	if err != nil {
		return err
	}

	return c.AppsV1().StatefulSets(statefulSet.Namespace).Delete(ctx, statefulSet.Name, metav1.DeleteOptions{})
}
