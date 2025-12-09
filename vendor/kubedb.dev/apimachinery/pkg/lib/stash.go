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

	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	kmapi "kmodules.xyz/client-go/api/v1"
	cutil "kmodules.xyz/client-go/conditions"
	core_util "kmodules.xyz/client-go/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"stash.appscode.dev/apimachinery/apis"
	stash "stash.appscode.dev/apimachinery/apis/stash/v1beta1"
	scs "stash.appscode.dev/apimachinery/client/clientset/versioned/typed/stash/v1beta1"
	scsutil "stash.appscode.dev/apimachinery/client/clientset/versioned/typed/stash/v1beta1/util"
)

func stashBackupOrRestoreRunningForDB(stashClient scs.StashV1beta1Interface, dbMeta metav1.ObjectMeta) (bool, string, error) {
	configs, err := stashClient.BackupConfigurations(dbMeta.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		if kerr.IsNotFound(err) {
			return false, "", nil
		}
		return false, "", err
	}

	var dbBackupConfig *stash.BackupConfiguration
	for _, config := range configs.Items {
		if config.Spec.Target.Ref.Name == dbMeta.Name && config.Spec.Target.Ref.Kind == apis.KindAppBinding {
			dbBackupConfig = &config
			break
		}
	}

	// skip running backup session
	if dbBackupConfig != nil {
		backupSessions, err := stashClient.BackupSessions(dbMeta.Namespace).List(context.TODO(), metav1.ListOptions{
			LabelSelector: labels.Set(dbBackupConfig.Labels).String(),
		})
		if err != nil {
			return false, "", err
		}

		for _, session := range backupSessions.Items {
			owned, _ := core_util.IsOwnedBy(&session, dbBackupConfig)
			if owned && (session.Status.Phase == stash.BackupSessionPending || session.Status.Phase == stash.BackupSessionRunning) {
				return true, fmt.Sprintf("BackupSession %s/%s is in %s Phase for this database", session.Namespace, session.Name, session.Status.Phase), nil
			}
		}
	}

	// backupconfiguration & backupsession crds can be installed with `mongodb`/`provisioner` operator.
	// so, we need some other crd to check if stash is installed or not.
	// restoresession crd is a good choice here.
	// if stash not installed (== IsNotFound error below), we just continue as if no error occurred.
	restoreSessions, err := stashClient.RestoreSessions(dbMeta.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		if kerr.IsNotFound(err) {
			return false, "", nil
		}
		return false, "", err
	}

	for _, session := range restoreSessions.Items {
		if session.Spec.Target.Ref.Name == dbMeta.Name && session.Spec.Target.Ref.Kind == apis.KindAppBinding && (session.Status.Phase == stash.RestorePending || session.Status.Phase == stash.RestoreRunning) {
			return true, fmt.Sprintf("RestoreSession %s/%s is in %s Phase for this database", session.Namespace, session.Name, session.Status.Phase), nil
		}
	}

	return false, "", nil
}

func pauseStashBackupConfiguration(stashClient scs.StashV1beta1Interface, dbMeta metav1.ObjectMeta, pausedBackups []kmapi.TypedObjectReference, opsGeneration int64) ([]kmapi.TypedObjectReference, []kmapi.Condition, error) {
	configs, err := stashClient.BackupConfigurations(dbMeta.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, nil, err
	}

	var opsConditions []kmapi.Condition
	var newPausedBackups []kmapi.TypedObjectReference
	for _, config := range configs.Items {
		if config.Spec.Target.Ref.Name == dbMeta.Name && config.Spec.Target.Ref.Kind == apis.KindAppBinding && !config.Spec.Paused {
			newBackup := kmapi.TypedObjectReference{
				APIGroup:  stash.SchemeGroupVersion.Group,
				Name:      config.Name,
				Namespace: config.Namespace,
			}
			if slices.Contains(pausedBackups, newBackup) {
				continue
			}
			_, err := scsutil.TryUpdateBackupConfiguration(context.TODO(), stashClient, config.ObjectMeta, func(configuration *stash.BackupConfiguration) *stash.BackupConfiguration {
				configuration.Spec.Paused = true
				return configuration
			}, metav1.UpdateOptions{})
			if err != nil {
				return nil, nil, err
			}
			newCondition := cutil.NewCondition(opsapi.PauseBackupConfiguration, fmt.Sprintf("BackupConfiguration %s/%s Paused", config.Namespace, config.Name), opsGeneration)
			newPausedBackups = append(newPausedBackups, newBackup)
			opsConditions = append(opsConditions, newCondition)

		}
	}
	return newPausedBackups, opsConditions, nil
}

func resumeStashBackupConfiguration(stashClient scs.StashV1beta1Interface, pausedBackups []kmapi.TypedObjectReference, opsGeneration int64) ([]kmapi.Condition, error) {
	var opsConditions []kmapi.Condition
	for _, config := range pausedBackups {
		if config.APIGroup != stash.SchemeGroupVersion.Group {
			continue
		}
		dbBackupConfig, err := stashClient.BackupConfigurations(config.Namespace).Get(context.Background(), config.Name, metav1.GetOptions{})
		if err != nil {
			return nil, client.IgnoreNotFound(err)
		}
		_, err = scsutil.TryUpdateBackupConfiguration(context.TODO(), stashClient, dbBackupConfig.ObjectMeta, func(configuration *stash.BackupConfiguration) *stash.BackupConfiguration {
			configuration.Spec.Paused = false
			return configuration
		}, metav1.UpdateOptions{})
		if err != nil {
			return nil, err
		}
		newCondition := cutil.NewCondition(opsapi.ResumeBackupConfiguration, fmt.Sprintf("BackupConfiguration %s/%s resumed", config.Namespace, config.Name), opsGeneration)
		opsConditions = append(opsConditions, newCondition)
	}
	return opsConditions, nil
}
