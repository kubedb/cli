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

package policy

import (
	"context"
	"strconv"

	"gomodules.xyz/sync"
	policyv1 "k8s.io/api/policy/v1"
	policyv1beta1 "k8s.io/api/policy/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/kubernetes"
)

var (
	onceEviction  sync.Once
	useEvictionV1 bool
)

func detectEvictionVersion(c discovery.DiscoveryInterface) {
	onceEviction.Do(func() error {
		// Note: policy/v1 Eviction is available in v1.22+. Use policy/v1beta1 with prior releases.
		// ref: https://kubernetes.io/docs/concepts/scheduling-eviction/api-eviction/#calling-the-eviction-api
		info, err := c.ServerVersion()
		if err != nil {
			return err
		}
		major, err := strconv.Atoi(info.Major)
		if err != nil {
			return err
		}
		minor, err := strconv.Atoi(info.Minor)
		if err != nil {
			return err
		}
		useEvictionV1 = major > 1 || (major == 1 && minor >= 22)
		return err
	})
}

func EvictPod(ctx context.Context, c kubernetes.Interface, meta types.NamespacedName, opts *metav1.DeleteOptions) error {
	detectEvictionVersion(c.Discovery())
	if useEvictionV1 {
		return c.CoreV1().Pods(meta.Namespace).EvictV1(ctx, &policyv1.Eviction{
			ObjectMeta: metav1.ObjectMeta{
				Name:      meta.Name,
				Namespace: meta.Namespace,
			},
			DeleteOptions: opts,
		})
	}
	return c.CoreV1().Pods(meta.Namespace).EvictV1beta1(ctx, &policyv1beta1.Eviction{
		ObjectMeta: metav1.ObjectMeta{
			Name:      meta.Name,
			Namespace: meta.Namespace,
		},
		DeleteOptions: opts,
	})
}
