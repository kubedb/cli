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

package v1beta1

import (
	"context"
	"encoding/json"

	api "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	cs "github.com/cert-manager/cert-manager/pkg/client/clientset/versioned/typed/certmanager/v1"
	jsonpatch "github.com/evanphx/json-patch"
	"github.com/pkg/errors"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/klog/v2"
	kutil "kmodules.xyz/client-go"
)

func CreateOrPatchClusterIssuer(ctx context.Context, c cs.CertmanagerV1Interface, meta metav1.ObjectMeta, transform func(*api.ClusterIssuer) *api.ClusterIssuer, opts metav1.PatchOptions) (*api.ClusterIssuer, kutil.VerbType, error) {
	cur, err := c.ClusterIssuers().Get(ctx, meta.Name, metav1.GetOptions{})
	if kerr.IsNotFound(err) {
		klog.V(3).Infof("Creating ClusterIssuer %s", meta.Name)
		out, err := c.ClusterIssuers().Create(ctx, transform(&api.ClusterIssuer{
			TypeMeta: metav1.TypeMeta{
				Kind:       "ClusterIssuer",
				APIVersion: api.SchemeGroupVersion.String(),
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
	return PatchClusterIssuer(ctx, c, cur, transform, opts)
}

func PatchClusterIssuer(ctx context.Context, c cs.CertmanagerV1Interface, cur *api.ClusterIssuer, transform func(*api.ClusterIssuer) *api.ClusterIssuer, opts metav1.PatchOptions) (*api.ClusterIssuer, kutil.VerbType, error) {
	return PatchClusterIssuerObject(ctx, c, cur, transform(cur.DeepCopy()), opts)
}

func PatchClusterIssuerObject(ctx context.Context, c cs.CertmanagerV1Interface, cur, mod *api.ClusterIssuer, opts metav1.PatchOptions) (*api.ClusterIssuer, kutil.VerbType, error) {
	curJson, err := json.Marshal(cur)
	if err != nil {
		return nil, kutil.VerbUnchanged, err
	}

	modJson, err := json.Marshal(mod)
	if err != nil {
		return nil, kutil.VerbUnchanged, err
	}

	patch, err := jsonpatch.CreateMergePatch(curJson, modJson)
	if err != nil {
		return nil, kutil.VerbUnchanged, err
	}
	if len(patch) == 0 || string(patch) == "{}" {
		return cur, kutil.VerbUnchanged, nil
	}
	klog.V(3).Infof("Patching ClusterIssuer %s with %s.", cur.Name, string(patch))
	out, err := c.ClusterIssuers().Patch(ctx, cur.Name, types.MergePatchType, patch, opts)
	return out, kutil.VerbPatched, err
}

func TryUpdateClusterIssuer(ctx context.Context, c cs.CertmanagerV1Interface, meta metav1.ObjectMeta, transform func(*api.ClusterIssuer) *api.ClusterIssuer, opts metav1.UpdateOptions) (result *api.ClusterIssuer, err error) {
	attempt := 0
	err = wait.PollImmediate(kutil.RetryInterval, kutil.RetryTimeout, func() (bool, error) {
		attempt++
		cur, e2 := c.ClusterIssuers().Get(ctx, meta.Name, metav1.GetOptions{})
		if kerr.IsNotFound(e2) {
			return false, e2
		} else if e2 == nil {
			result, e2 = c.ClusterIssuers().Update(ctx, transform(cur.DeepCopy()), opts)
			return e2 == nil, nil
		}
		klog.Errorf("Attempt %d failed to update ClusterIssuer %s due to %v.", attempt, cur.Name, e2)
		return false, nil
	})

	if err != nil {
		err = errors.Errorf("failed to update ClusterIssuer %s after %d attempts due to %v", meta.Name, attempt, err)
	}
	return
}
