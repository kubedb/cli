/*
Copyright AppsCode Inc. and Contributors

Licensed under the PolyForm Noncommercial License 1.0.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://github.com/appscode/licenses/raw/1.0.0/PolyForm-Noncommercial-1.0.0.md

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package describer

import (
	"context"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/duration"
	"k8s.io/kubectl/pkg/describe"
	appcat "kmodules.xyz/custom-resources/apis/appcatalog/v1alpha1"
	stashV1beta1 "stash.appscode.dev/apimachinery/apis/stash/v1beta1"
	stash "stash.appscode.dev/apimachinery/client/clientset/versioned"
)

const (
	KindAppBinding string = "AppBinding"
)

type backupInvokerInfo struct {
	name              string
	kind              string
	schedule          string
	task              string
	repository        string
	bucket            string
	creationTimestamp metav1.Time
}

func showBackups(stash stash.Interface, ab *appcat.AppBinding, w describe.PrefixWriter) error {
	w.Write(LEVEL_0, "\n")
	w.Write(LEVEL_0, "Backup:\n")
	var invokers []backupInvokerInfo
	// There could be two types of backup invokers.
	// 1. BackupConfiguration
	// 2. BackupBatch

	// Get BackupConfiguration type invokers
	bcInvokers, err := getBackupConfigurationTypeInvokers(stash, ab)
	if err != nil {
		return err
	}
	invokers = append(invokers, bcInvokers...)

	// Get BackupBatch type invokers
	bbInvokers, err := getBackupBatchTypeInvokers(stash, ab)
	if err != nil {
		return err
	}
	invokers = append(invokers, bbInvokers...)

	if len(invokers) == 0 {
		w.Write(LEVEL_1, "No backup has been configured.\n")
		return nil
	}
	// Print the backup invokers table
	w.Write(LEVEL_1, "Backup Invokers:\n")
	w.Write(LEVEL_2, "Name\tKind\tSchedule\tTask\tRepository\tBucket\tAge\n")
	w.Write(LEVEL_2, "----\t----\t--------\t----\t----------\t------\t---\n")
	for _, invk := range invokers {
		age := duration.HumanDuration(time.Since(invk.creationTimestamp.Time))
		w.Write(LEVEL_2, "%s\t%s\t%s\t%s\t%s\t%s\t%s\n", invk.name, invk.kind, invk.schedule, invk.task, invk.repository, invk.bucket, age)
	}

	// Get the BackupSessions for the above invokers
	backupSessions, err := getBackupSessions(stash, ab.Namespace, invokers)
	if err != nil {
		return err
	}
	// Print recent backup table
	if len(backupSessions) != 0 {
		w.Write(LEVEL_1, "Recent Backups:\n")
		w.Write(LEVEL_2, "Name\tInvoker-kind\tInvoker-name\tPhase\tAge\n")
		w.Write(LEVEL_2, "----\t------------\t------------\t-----\t---\n")
		for _, bs := range backupSessions {
			age := duration.HumanDuration(time.Since(bs.CreationTimestamp.Time))
			w.Write(LEVEL_2, "%s\t%s\t%s\t%s\t%s\n", bs.Name, bs.Spec.Invoker.Kind, bs.Spec.Invoker.Name, bs.Status.Phase, age)
		}
	}
	return nil
}

func getBackupConfigurationTypeInvokers(stash stash.Interface, ab *appcat.AppBinding) ([]backupInvokerInfo, error) {
	var bcInvokers []backupInvokerInfo
	backupConfigurations, err := stash.StashV1beta1().BackupConfigurations(ab.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	// Identify those BackupConfigurations that has this AppBinding as target
	for _, bc := range backupConfigurations.Items {
		if bc.Spec.Target != nil &&
			bc.Spec.Target.Ref.Kind == KindAppBinding &&
			bc.Spec.Target.Ref.Name == ab.Name {
			invoker := backupInvokerInfo{
				name:              bc.Name,
				kind:              bc.Kind,
				schedule:          bc.Spec.Schedule,
				task:              bc.Spec.Task.Name,
				repository:        bc.Spec.Repository.Name,
				creationTimestamp: bc.CreationTimestamp,
			}
			bucket, err := getBucket(stash, bc.Spec.Repository.Name, bc.Namespace)
			if err != nil {
				return nil, err
			}
			invoker.bucket = bucket

			bcInvokers = append(bcInvokers, invoker)
		}
	}
	return bcInvokers, nil
}

func getBackupBatchTypeInvokers(stash stash.Interface, ab *appcat.AppBinding) ([]backupInvokerInfo, error) {
	var bbInvokers []backupInvokerInfo
	backupBatches, err := stash.StashV1beta1().BackupBatches(ab.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, bb := range backupBatches.Items {
		for _, m := range bb.Spec.Members {
			if m.Target != nil &&
				m.Target.Ref.Kind == KindAppBinding &&
				m.Target.Ref.Name == ab.Name {
				invoker := backupInvokerInfo{
					name:              bb.Name,
					kind:              bb.Kind,
					schedule:          bb.Spec.Schedule,
					task:              m.Task.Name,
					repository:        bb.Spec.Repository.Name,
					creationTimestamp: bb.CreationTimestamp,
				}
				bucket, err := getBucket(stash, bb.Spec.Repository.Name, bb.Namespace)
				if err != nil {
					return nil, err
				}
				invoker.bucket = bucket

				bbInvokers = append(bbInvokers, invoker)
			}
		}
	}
	return bbInvokers, nil
}

func getBucket(stash stash.Interface, repoName, repoNamespace string) (string, error) {
	// Get the respective Repository
	repo, err := stash.StashV1alpha1().Repositories(repoNamespace).Get(context.TODO(), repoName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	return repo.Spec.Backend.Container()
}

func getBackupSessions(stash stash.Interface, namespace string, invokers []backupInvokerInfo) ([]stashV1beta1.BackupSession, error) {
	var backupSessions []stashV1beta1.BackupSession

	bsList, err := stash.StashV1beta1().BackupSessions(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for i, bs := range bsList.Items {
		if ownByInvoker(bs, invokers) {
			backupSessions = append(backupSessions, bsList.Items[i])
		}
	}
	return backupSessions, nil
}

func ownByInvoker(bs stashV1beta1.BackupSession, invokers []backupInvokerInfo) bool {
	for i := range invokers {
		if invokers[i].kind == bs.Spec.Invoker.Kind &&
			invokers[i].name == bs.Spec.Invoker.Name {
			return true
		}
	}
	return false
}
