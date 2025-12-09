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

//nolint:unused
package policy

import (
	"context"

	du "kmodules.xyz/client-go/discovery"
	v1 "kmodules.xyz/client-go/policy/v1"
	"kmodules.xyz/client-go/policy/v1beta1"

	"gomodules.xyz/sync"
	policyv1 "k8s.io/api/policy/v1"
	policyv1beta1 "k8s.io/api/policy/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/kubernetes"
	kutil "kmodules.xyz/client-go"
)

const kindPodDisruptionBudget = "PodDisruptionBudget"

var (
	oncePDB  sync.Once
	usePDBV1 bool
)

func detectPDBVersion(c discovery.DiscoveryInterface) {
	oncePDB.Do(func() error {
		ok, err := du.HasGVK(c, policyv1.SchemeGroupVersion.String(), kindPodDisruptionBudget)
		usePDBV1 = ok
		return err
	})
}

func CreateOrPatchPodDisruptionBudget(ctx context.Context, c kubernetes.Interface, meta metav1.ObjectMeta, transform func(*policyv1.PodDisruptionBudget) *policyv1.PodDisruptionBudget, opts metav1.PatchOptions) (*policyv1.PodDisruptionBudget, kutil.VerbType, error) {
	detectPDBVersion(c.Discovery())
	if usePDBV1 {
		return v1.CreateOrPatchPodDisruptionBudget(ctx, c, meta, transform, opts)
	}

	p, vt, err := v1beta1.CreateOrPatchPodDisruptionBudget(
		ctx,
		c,
		meta,
		func(in *policyv1beta1.PodDisruptionBudget) *policyv1beta1.PodDisruptionBudget {
			out := convert_spec_v1_to_v1beta1(transform(convert_spec_v1beta1_to_v1(in)))
			out.Status = in.Status
			return out
		},
		opts,
	)
	if err != nil {
		return nil, kutil.VerbUnchanged, err
	}
	return convert_v1beta1_to_v1(p), vt, nil
}

func CreatePodDisruptionBudget(ctx context.Context, c kubernetes.Interface, in *policyv1.PodDisruptionBudget) (*policyv1.PodDisruptionBudget, error) {
	detectPDBVersion(c.Discovery())
	if usePDBV1 {
		return c.PolicyV1().PodDisruptionBudgets(in.Namespace).Create(ctx, in, metav1.CreateOptions{})
	}
	result, err := c.PolicyV1beta1().PodDisruptionBudgets(in.Namespace).Create(ctx, convert_spec_v1_to_v1beta1(in), metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}
	return convert_v1beta1_to_v1(result), nil
}

func GetPodDisruptionBudget(ctx context.Context, c kubernetes.Interface, meta types.NamespacedName) (*policyv1.PodDisruptionBudget, error) {
	detectPDBVersion(c.Discovery())
	if usePDBV1 {
		return c.PolicyV1().PodDisruptionBudgets(meta.Namespace).Get(ctx, meta.Name, metav1.GetOptions{})
	}
	result, err := c.PolicyV1beta1().PodDisruptionBudgets(meta.Namespace).Get(ctx, meta.Name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return convert_v1beta1_to_v1(result), nil
}

func ListPodDisruptionBudget(ctx context.Context, c kubernetes.Interface, ns string, opts metav1.ListOptions) (*policyv1.PodDisruptionBudgetList, error) {
	detectPDBVersion(c.Discovery())
	if usePDBV1 {
		return c.PolicyV1().PodDisruptionBudgets(ns).List(ctx, opts)
	}
	result, err := c.PolicyV1beta1().PodDisruptionBudgets(ns).List(ctx, opts)
	if err != nil {
		return nil, err
	}
	out := policyv1.PodDisruptionBudgetList{
		TypeMeta: result.TypeMeta,
		ListMeta: result.ListMeta,
		Items:    make([]policyv1.PodDisruptionBudget, 0, len(result.Items)),
	}
	for _, item := range result.Items {
		out.Items = append(out.Items, *convert_v1beta1_to_v1(&item))
	}
	return &out, nil
}

func DeletePodDisruptionBudget(ctx context.Context, c kubernetes.Interface, meta types.NamespacedName) error {
	detectPDBVersion(c.Discovery())
	if usePDBV1 {
		return c.PolicyV1().PodDisruptionBudgets(meta.Namespace).Delete(ctx, meta.Name, metav1.DeleteOptions{})
	}
	return c.PolicyV1beta1().PodDisruptionBudgets(meta.Namespace).Delete(ctx, meta.Name, metav1.DeleteOptions{})
}

func DeletePodDisruptionBudgets(ctx context.Context, c kubernetes.Interface, ns string, listOpts metav1.ListOptions, deleteOpts metav1.DeleteOptions) error {
	detectPDBVersion(c.Discovery())
	if usePDBV1 {
		return c.PolicyV1().PodDisruptionBudgets(ns).DeleteCollection(ctx, deleteOpts, listOpts)
	}
	return c.PolicyV1beta1().PodDisruptionBudgets(ns).DeleteCollection(ctx, deleteOpts, listOpts)
}

func convert_spec_v1beta1_to_v1(in *policyv1beta1.PodDisruptionBudget) *policyv1.PodDisruptionBudget {
	return &policyv1.PodDisruptionBudget{
		TypeMeta: metav1.TypeMeta{
			Kind:       in.Kind,
			APIVersion: policyv1.SchemeGroupVersion.String(),
		},
		ObjectMeta: in.ObjectMeta,
		Spec: policyv1.PodDisruptionBudgetSpec{
			MinAvailable:   in.Spec.MinAvailable,
			Selector:       in.Spec.Selector,
			MaxUnavailable: in.Spec.MaxUnavailable,
		},
	}
}

func convert_v1beta1_to_v1(in *policyv1beta1.PodDisruptionBudget) *policyv1.PodDisruptionBudget {
	out := convert_spec_v1beta1_to_v1(in)
	out.Status = policyv1.PodDisruptionBudgetStatus{
		ObservedGeneration: in.Status.ObservedGeneration,
		DisruptedPods:      in.Status.DisruptedPods,
		DisruptionsAllowed: in.Status.DisruptionsAllowed,
		CurrentHealthy:     in.Status.CurrentHealthy,
		DesiredHealthy:     in.Status.DesiredHealthy,
		ExpectedPods:       in.Status.ExpectedPods,
		Conditions:         in.Status.Conditions,
	}
	return out
}

func convert_spec_v1_to_v1beta1(in *policyv1.PodDisruptionBudget) *policyv1beta1.PodDisruptionBudget {
	return &policyv1beta1.PodDisruptionBudget{
		TypeMeta: metav1.TypeMeta{
			Kind:       in.Kind,
			APIVersion: policyv1beta1.SchemeGroupVersion.String(),
		},
		ObjectMeta: in.ObjectMeta,
		Spec: policyv1beta1.PodDisruptionBudgetSpec{
			MinAvailable:   in.Spec.MinAvailable,
			Selector:       in.Spec.Selector,
			MaxUnavailable: in.Spec.MaxUnavailable,
		},
	}
}

func convert_v1_to_v1beta1(in *policyv1.PodDisruptionBudget) *policyv1beta1.PodDisruptionBudget {
	out := convert_spec_v1_to_v1beta1(in)
	out.Status = policyv1beta1.PodDisruptionBudgetStatus{
		ObservedGeneration: in.Status.ObservedGeneration,
		DisruptedPods:      in.Status.DisruptedPods,
		DisruptionsAllowed: in.Status.DisruptionsAllowed,
		CurrentHealthy:     in.Status.CurrentHealthy,
		DesiredHealthy:     in.Status.DesiredHealthy,
		ExpectedPods:       in.Status.ExpectedPods,
		Conditions:         in.Status.Conditions,
	}
	return out
}
