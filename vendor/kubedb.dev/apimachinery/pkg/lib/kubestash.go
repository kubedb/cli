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
	"slices"

	opsapi "kubedb.dev/apimachinery/apis/ops/v1alpha1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kmapi "kmodules.xyz/client-go/api/v1"
	kmc "kmodules.xyz/client-go/client"
	cutil "kmodules.xyz/client-go/conditions"
	"kubestash.dev/apimachinery/apis"
	coreapi "kubestash.dev/apimachinery/apis/core/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func kubeStashBackupOrRestoreRunningForDB(KBClient client.Client, dbObjMeta metav1.ObjectMeta, kind string) (bool, string, error) {
	backupConfigList, err := getBackupConfigList(KBClient)
	if err != nil {
		return false, "", err
	}
	for _, config := range backupConfigList.Items {
		if config.Spec.Target == nil {
			continue
		}
		if matchesTarget(*config.Spec.Target, dbObjMeta, kind, config.Namespace) {
			// skip for running backup session
			if runningSession, err := getRunningBackupSession(KBClient, config); err != nil || runningSession != nil {
				if err != nil {
					return false, "", err
				}
				return true, fmt.Sprintf("BackupSession %s/%s is in %s Phase for this database", runningSession.Namespace, runningSession.Name, runningSession.Status.Phase), nil
			}
		}
	}

	restoreSessionList := coreapi.RestoreSessionList{}
	if err := KBClient.List(context.Background(), &restoreSessionList); err != nil {
		return false, "", err
	}

	for _, session := range restoreSessionList.Items {
		if session.Spec.Target == nil {
			continue
		}
		if matchesTarget(*session.Spec.Target, dbObjMeta, kind, session.Namespace) &&
			(session.Status.Phase == coreapi.RestorePending || session.Status.Phase == coreapi.RestoreRunning) {
			return true, fmt.Sprintf("RestoreSession %s/%s is in %s Phase for this database", session.Namespace, session.Name, session.Status.Phase), nil
		}
	}

	return false, "", nil
}

func pauseKubeStashBackupConfiguration(KBClient client.Client, dbObjMeta metav1.ObjectMeta, pausedBackups []kmapi.TypedObjectReference, kind string, opsGeneration int64) ([]kmapi.TypedObjectReference, []kmapi.Condition, error) {
	backupConfigList, err := getBackupConfigList(KBClient)
	if err != nil {
		return nil, nil, err
	}

	var opsConditions []kmapi.Condition
	var newPausedBackups []kmapi.TypedObjectReference
	for _, config := range backupConfigList.Items {
		if config.Spec.Target == nil {
			continue
		}
		if matchesTarget(*config.Spec.Target, dbObjMeta, kind, config.Namespace) && !config.Spec.Paused {
			newBackup := kmapi.TypedObjectReference{
				APIGroup:  coreapi.GroupVersion.Group,
				Name:      config.Name,
				Namespace: config.Namespace,
			}
			if slices.Contains(pausedBackups, newBackup) {
				continue
			}

			if err := modifyBackupConfiguration(KBClient, &config, true); err != nil {
				return nil, nil, err
			}
			newCondition := cutil.NewCondition(opsapi.PauseBackupConfiguration, fmt.Sprintf("BackupConfiguration %s/%s Paused", config.Namespace, config.Name), opsGeneration)
			newPausedBackups = append(newPausedBackups, newBackup)
			opsConditions = append(opsConditions, newCondition)
		}
	}
	return newPausedBackups, opsConditions, nil
}

func resumeKubeStashBackupConfiguration(KBClient client.Client, pausedBackups []kmapi.TypedObjectReference, opsGeneration int64) ([]kmapi.Condition, error) {
	var opsConditions []kmapi.Condition
	for _, config := range pausedBackups {
		if config.APIGroup != coreapi.GroupVersion.Group {
			continue
		}
		bc := &coreapi.BackupConfiguration{
			ObjectMeta: metav1.ObjectMeta{
				Name:      config.Name,
				Namespace: config.Namespace,
			},
		}
		if err := KBClient.Get(context.Background(), client.ObjectKeyFromObject(bc), bc); err != nil {
			return nil, client.IgnoreNotFound(err)
		}

		if err := modifyBackupConfiguration(KBClient, bc, false); err != nil {
			return nil, err
		}
		newCondition := cutil.NewCondition(opsapi.ResumeBackupConfiguration, fmt.Sprintf("BackupConfiguration %s/%s resumed", config.Namespace, config.Name), opsGeneration)
		opsConditions = append(opsConditions, newCondition)
	}
	return opsConditions, nil
}

func modifyBackupConfiguration(KBClient client.Client, config *coreapi.BackupConfiguration, paused bool) error {
	_, err := kmc.CreateOrPatch(
		context.Background(),
		KBClient,
		config,
		func(obj client.Object, createOp bool) client.Object {
			in := obj.(*coreapi.BackupConfiguration)
			in.Spec.Paused = paused
			return in
		},
	)

	return err
}

func getBackupConfigList(KBClient client.Client) (*coreapi.BackupConfigurationList, error) {
	backupConfigList := &coreapi.BackupConfigurationList{}
	if err := KBClient.List(context.Background(), backupConfigList); err != nil {
		return backupConfigList, err
	}

	return backupConfigList, nil
}

func getRunningBackupSession(KBClient client.Client, config coreapi.BackupConfiguration) (*coreapi.BackupSession, error) {
	backupSessionList := &coreapi.BackupSessionList{}
	opts := []client.ListOption{
		client.MatchingLabels{
			apis.KubeStashInvokerName: config.Name,
		},
	}
	if err := KBClient.List(context.Background(), backupSessionList, opts...); err != nil {
		return nil, err
	}
	for _, session := range backupSessionList.Items {
		if session.Status.Phase == coreapi.BackupSessionPending || session.Status.Phase == coreapi.BackupSessionRunning {
			return &session, nil
		}
	}
	return nil, nil
}

func matchesTarget(targetRef kmapi.TypedObjectReference, dbMeta metav1.ObjectMeta, kind, objNS string) bool {
	if targetRef.Name != dbMeta.Name || targetRef.Kind != kind {
		return false
	}
	if targetRef.Namespace != "" {
		return targetRef.Namespace == dbMeta.Namespace
	}
	return objNS == dbMeta.Namespace
}
