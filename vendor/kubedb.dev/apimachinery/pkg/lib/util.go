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
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"kubedb.dev/apimachinery/apis/kubedb"

	"github.com/go-logr/logr"
	vsecretapi "go.virtual-secrets.dev/apimachinery/apis/virtual/v1alpha1"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"
	kmapi "kmodules.xyz/client-go/api/v1"
	clientutil "kmodules.xyz/client-go/client"
	"kmodules.xyz/client-go/policy"
	psapi "kubeops.dev/petset/apis/apps/v1"
	kubestash_api "kubestash.dev/apimachinery/apis/core/v1alpha1"
	ocmapi "open-cluster-management.io/api/work/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	stash_api "stash.appscode.dev/apimachinery/apis/stash/v1beta1"
)

func EvictPod(kClient kubernetes.Interface, podMeta metav1.ObjectMeta) (types.UID, error) {
	pod, err := kClient.CoreV1().Pods(podMeta.Namespace).Get(context.TODO(), podMeta.Name, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	return pod.UID, policy.EvictPod(context.TODO(), kClient, types.NamespacedName{
		Namespace: podMeta.Namespace,
		Name:      podMeta.Name,
	}, &metav1.DeleteOptions{})
}

// Evict Pod restart the pod by deleting the manifestwork resouces.
// podMeta.namespace contain the origin namespace of pod rather than the distributed manifest namespace.
func EvictDistributedPod(kbClient client.Client, pp, podName string) (types.UID, error) {
	namespace, err := GetDistributedPodNamespace(kbClient, pp, podName)
	if err != nil {
		return "", err
	}
	mw := &ocmapi.ManifestWork{}
	err = kbClient.Get(context.TODO(), types.NamespacedName{Name: podName, Namespace: namespace}, mw, &client.GetOptions{})
	if err != nil {
		return "", err
	}
	deletionPolicy := metav1.DeletePropagationOrphan
	return mw.UID, kbClient.Delete(context.TODO(), mw, &client.DeleteOptions{
		PropagationPolicy: &deletionPolicy,
	})
}

// get the manifestwork namespace for podManifestWork
func GetDistributedPodNamespace(kbClient client.Client, ppName, podName string) (string, error) {
	pp := psapi.PlacementPolicy{}
	err := kbClient.Get(context.TODO(), types.NamespacedName{
		Name: ppName,
	}, &pp)
	if err != nil {
		klog.Errorln(err)
		return "", err
	}
	splitString := strings.Split(podName, "-")
	ordinal, _ := strconv.Atoi(splitString[len(splitString)-1])
	for i := 0; i < len(pp.Spec.ClusterSpreadConstraint.DistributionRules); i++ {
		for _, ordinalIndex := range pp.Spec.ClusterSpreadConstraint.DistributionRules[i].ReplicaIndices {
			if ordinalIndex == int32(ordinal) {
				return pp.Spec.ClusterSpreadConstraint.DistributionRules[i].ClusterName, nil
			}
		}
	}
	return "", fmt.Errorf("no cluster found for the given ordinal %v", ordinal)
}

// extractPodFromManifestWork return the pod manifest. podMeta contain the original namespace of the pods.
func ExtractPodFromManifestWork(kbClient client.Client, pp, podName string) (*core.Pod, error) {
	namespace, err := GetDistributedPodNamespace(kbClient, pp, podName)
	if err != nil {
		return nil, err
	}
	mw := &ocmapi.ManifestWork{}
	err = kbClient.Get(context.TODO(), types.NamespacedName{Name: podName, Namespace: namespace}, mw, &client.GetOptions{})
	if err != nil {
		klog.Error(err, "Failed to get ManifestWork", "ManifestWork", podName)
		return nil, err
	}
	pod := &core.Pod{}
	manifest := mw.Spec.Workload.Manifests[0]
	unstructuredObj := make(map[string]any)
	err = json.Unmarshal(manifest.Raw, &unstructuredObj)
	if err != nil {
		return nil, err
	}
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(unstructuredObj, pod)
	if err != nil {
		return nil, err
	}
	return pod, nil
}

func ExtractPVCFromManifestWork(kbClient client.Client, pp, podName string) (*core.PersistentVolumeClaim, error) {
	namespace, err := GetDistributedPodNamespace(kbClient, pp, podName)
	if err != nil {
		return nil, err
	}
	pvcName := fmt.Sprintf("%s-%s", kubedb.DefaultVolumeClaimTemplateName, podName)
	mw := &ocmapi.ManifestWork{}
	err = kbClient.Get(context.TODO(), types.NamespacedName{Name: pvcName, Namespace: namespace}, mw, &client.GetOptions{})
	if err != nil {
		klog.Error(err, "Failed to get ManifestWork", "ManifestWork", podName)
		return nil, err
	}
	pvc := &core.PersistentVolumeClaim{}
	manifest := mw.Spec.Workload.Manifests[0]
	unstructuredObj := make(map[string]any)
	err = json.Unmarshal(manifest.Raw, &unstructuredObj)
	if err != nil {
		return nil, err
	}
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(unstructuredObj, pvc)
	if err != nil {
		return nil, err
	}
	for _, manifestStatus := range mw.Status.ResourceStatus.Manifests {
		feedback := manifestStatus.StatusFeedbacks
		for _, value := range feedback.Values {
			switch value.Name {
			case "Capacity":
				if value.Value.String != nil {
					quantity, err := resource.ParseQuantity(*value.Value.String)
					if err != nil {
						return nil, err
					}
					if pvc.Status.Capacity == nil {
						pvc.Status.Capacity = make(core.ResourceList)
					}
					pvc.Status.Capacity[core.ResourceStorage] = quantity
				}
			}
		}
	}
	return pvc, nil
}

func UpdatePVCManifestWork(kbClient client.Client, reqResource core.VolumeResourceRequirements, pp, podName string) error {
	namespace, err := GetDistributedPodNamespace(kbClient, pp, podName)
	if err != nil {
		return err
	}
	pvcName := fmt.Sprintf("%s-%s", kubedb.DefaultVolumeClaimTemplateName, podName)
	mw := &ocmapi.ManifestWork{}
	err = kbClient.Get(context.TODO(), types.NamespacedName{Name: pvcName, Namespace: namespace}, mw, &client.GetOptions{})
	if err != nil {
		klog.Error(err, "Failed to get ManifestWork", "ManifestWork", podName)
		return err
	}
	pvc := &core.PersistentVolumeClaim{}
	manifest := mw.Spec.Workload.Manifests[0]
	unstructuredObj := make(map[string]any)
	err = json.Unmarshal(manifest.Raw, &unstructuredObj)
	if err != nil {
		return err
	}
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(unstructuredObj, pvc)
	if err != nil {
		return err
	}
	pvc.Spec.Resources = reqResource

	return UpdateManifestWork(kbClient, mw.Name, mw.Namespace, pvc)
}

func ObjectToUnstructured(obj runtime.Object) (*unstructured.Unstructured, error) {
	unstructuredMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		return nil, fmt.Errorf("failed to convert object to unstructured: %w", err)
	}
	return &unstructured.Unstructured{Object: unstructuredMap}, nil
}

func UpdateManifestWork(kbClient client.Client, manifestName, namespace string, objects ...client.Object) error {
	if len(objects) == 0 {
		return nil
	}

	var unstructuredObjects []*unstructured.Unstructured
	for _, object := range objects {
		obj, err := ObjectToUnstructured(object)
		if err != nil {
			return err
		}
		unstructuredObjects = append(unstructuredObjects, obj)
	}

	klog.V(5).Infof("Updating ManifestWork %s/%s", namespace, manifestName)

	mw := &ocmapi.ManifestWork{
		ObjectMeta: metav1.ObjectMeta{
			Name:      manifestName,
			Namespace: namespace,
		},
	}

	err := kbClient.Get(context.TODO(), types.NamespacedName{Name: mw.Name, Namespace: mw.Namespace}, mw, &client.GetOptions{})
	if err != nil {
		return err
	}
	manifest := []ocmapi.Manifest{}
	for _, unstructuredObject := range unstructuredObjects {
		manifest = append(manifest, ocmapi.Manifest{
			RawExtension: runtime.RawExtension{Object: unstructuredObject},
		})
	}
	mw.Spec.Workload.Manifests = manifest
	_, err = clientutil.CreateOrPatch(context.TODO(), kbClient, mw, func(obj client.Object, createOp bool) client.Object {
		in := obj.(*ocmapi.ManifestWork)
		in.Spec = mw.Spec
		return in
	})
	if err != nil {
		return fmt.Errorf("failed to create/patch manifestwork %s in namespace %s: %w", manifestName, namespace, err)
	}

	return nil
}

// Check the Pod ready condition from manifestWork condition feedback.
func CheckManifestWorkPodReady(mw *ocmapi.ManifestWork, log *logr.Logger) bool {
	for _, manifestStatus := range mw.Status.ResourceStatus.Manifests {
		feedback := manifestStatus.StatusFeedbacks
		for _, value := range feedback.Values {
			switch value.Name {
			case "PodPhase":
				if value.Value.String != nil && *value.Value.String != "Running" {
					log.Info("Pod Phase", "Pod Phase", *value.Value.String)
					return false
				}
			case "ReadyCondition":
				if value.Value.JsonRaw != nil {
					var readyCondition core.PodCondition
					if err := json.Unmarshal([]byte(*value.Value.JsonRaw), &readyCondition); err == nil {
						if readyCondition.Type != core.PodReady || readyCondition.Status != core.ConditionTrue {
							log.Info("Pod Status", "Pod Status", readyCondition.Status)
							return false
						}
					}
				}
			}
		}
	}
	return true
}

func IsOpsTypeSupported(supportedTypes []string, curOpsType string) bool {
	for _, s := range supportedTypes {
		if s == curOpsType {
			return true
		}
	}
	return false
}

func stashOperatorExist(KBClient client.Client) bool {
	_, err1 := KBClient.RESTMapper().RESTMapping(schema.GroupKind{
		Group: stash_api.SchemeGroupVersion.Group,
		Kind:  stash_api.ResourceKindBackupSession,
	})
	_, err2 := KBClient.RESTMapper().RESTMapping(schema.GroupKind{
		Group: stash_api.SchemeGroupVersion.Group,
		Kind:  stash_api.ResourceKindRestoreSession,
	})
	return err1 == nil && err2 == nil
}

func kubeStashOperatorExist(KBClient client.Client) bool {
	_, err1 := KBClient.RESTMapper().RESTMapping(schema.GroupKind{
		Group: kubestash_api.GroupVersion.Group,
		Kind:  kubestash_api.ResourceKindBackupSession,
	})
	_, err2 := KBClient.RESTMapper().RESTMapping(schema.GroupKind{
		Group: kubestash_api.GroupVersion.Group,
		Kind:  kubestash_api.ResourceKindRestoreSession,
	})
	return err1 == nil && err2 == nil
}

func GetSecret(kbClient client.Client, namespacedName types.NamespacedName) (*core.Secret, error) {
	var secret core.Secret
	err := kbClient.Get(context.Background(), types.NamespacedName{
		Namespace: namespacedName.Namespace,
		Name:      namespacedName.Name,
	}, &secret)
	if err != nil {
		return nil, fmt.Errorf("failed to get secret %s/%s. error: %v", namespacedName.Namespace, namespacedName.Name, err)
	}
	return &secret, nil
}

func ValidateAuthSecret(secret *core.Secret) error {
	// verify if the desired key ["password", "username"] exist or not when secret is managed by the user
	err := validateAuthSecretData(secret.Data)
	if err != nil {
		err = errors.Join(err, fmt.Errorf(" for secret %s/%s", secret.Namespace, secret.Name))
	}
	return err
}

func ValidateVirtualAuthSecret(secret *vsecretapi.Secret) error {
	// verify if the desired key ["password", "username"] exist or not when secret is managed by the user
	err := validateAuthSecretData(secret.Data)
	if err != nil {
		err = errors.Join(err, fmt.Errorf(" for virtual secret %s/%s", secret.Namespace, secret.Name))
	}
	return err
}

func validateAuthSecretData(authData map[string][]byte) error {
	// verify if the desired key ["password", "username"] exist or not when secret is managed by the user
	if authData == nil {
		return fmt.Errorf("key \"%s\" & \"%s\" doesn't exists inside spec data", core.BasicAuthUsernameKey, core.BasicAuthPasswordKey)
	}
	if userName, ok := authData[core.BasicAuthUsernameKey]; !ok || string(userName) == "" {
		return fmt.Errorf("key \"%s\" doesn't exists inside spec data", core.BasicAuthUsernameKey)
	}
	if pass, ok := authData[core.BasicAuthPasswordKey]; !ok || string(pass) == "" {
		return fmt.Errorf("key \"%s\" doesn't exists inside spec data", core.BasicAuthPasswordKey)
	}
	return nil
}

func UpdateMachineProfileAnnotation(db, ops map[string]string) {
	if db == nil || ops == nil {
		return
	}
	if val, exist := ops[kmapi.AceMachineProfileKey]; exist {
		db[kmapi.AceMachineProfileKey] = val
	}
}
