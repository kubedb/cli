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

	v1 "kubeops.dev/petset/apis/apps/v1"

	"gomodules.xyz/pointer"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	kutil "kmodules.xyz/client-go"
)

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
