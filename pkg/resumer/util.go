/*
Copyright AppsCode Inc. and Contributors

Licensed under the AppsCode Community License 1.0.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://github.com/appscode/licenses/raw/1.0.0/AppsCode-Community-1.0.0.md

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package resumer

import (
	"context"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"stash.appscode.dev/apimachinery/apis"
	stash "stash.appscode.dev/apimachinery/apis/stash/v1beta1"
	scs "stash.appscode.dev/apimachinery/client/clientset/versioned/typed/stash/v1beta1"
	scsutil "stash.appscode.dev/apimachinery/client/clientset/versioned/typed/stash/v1beta1/util"
)

const (
	ResumeTimeout  = 5 * time.Minute
	ResumeInterval = 5 * time.Second
)

func ResumeBackupConfiguration(stashClient scs.StashV1beta1Interface, dbMeta metav1.ObjectMeta) (bool, error) {
	configs, err := stashClient.BackupConfigurations(dbMeta.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return false, err
	}

	var dbBackupConfig *stash.BackupConfiguration
	for _, config := range configs.Items {
		if config.Spec.Target.Ref.Name == dbMeta.Name && config.Spec.Target.Ref.Kind == apis.KindAppBinding {
			dbBackupConfig = &config
			break
		}
	}

	if dbBackupConfig != nil {
		_, err := scsutil.TryUpdateBackupConfiguration(context.TODO(), stashClient, dbBackupConfig.ObjectMeta, func(configuration *stash.BackupConfiguration) *stash.BackupConfiguration {
			configuration.Spec.Paused = false
			return configuration
		}, metav1.UpdateOptions{})
		if err != nil {
			return false, err
		}
	}

	return dbBackupConfig != nil, nil
}
