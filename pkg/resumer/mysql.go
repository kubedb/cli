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

	api "kubedb.dev/apimachinery/apis/kubedb/v1alpha2"
	cs "kubedb.dev/apimachinery/client/clientset/versioned/typed/kubedb/v1alpha2"
	dbutil "kubedb.dev/apimachinery/client/clientset/versioned/typed/kubedb/v1alpha2/util"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/rest"
	condutil "kmodules.xyz/client-go/conditions"
	scs "stash.appscode.dev/apimachinery/client/clientset/versioned/typed/stash/v1beta1"
)

type MySQLResumer struct {
	dbClient    cs.KubedbV1alpha2Interface
	stashClient scs.StashV1beta1Interface
	onlyDb      bool
	onlyBackup  bool
}

func NewMySQLResumer(clientConfig *rest.Config, onlyDb, onlyBackup bool) (*MySQLResumer, error) {
	dbClient, err := cs.NewForConfig(clientConfig)
	if err != nil {
		return nil, err
	}
	stashClient, err := scs.NewForConfig(clientConfig)
	if err != nil {
		return nil, err
	}

	return &MySQLResumer{
		dbClient:    dbClient,
		stashClient: stashClient,
		onlyDb:      onlyDb,
		onlyBackup:  onlyBackup,
	}, nil
}

func (e *MySQLResumer) Resume(name, namespace string) (bool, error) {
	db, err := e.dbClient.MySQLs(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return false, err
	}

	resumeAll := !(e.onlyBackup || e.onlyDb)

	if e.onlyDb || resumeAll {
		_, err = dbutil.UpdateMySQLStatus(context.TODO(), e.dbClient, db.ObjectMeta, func(status *api.MySQLStatus) (types.UID, *api.MySQLStatus) {
			status.Conditions = condutil.RemoveCondition(status.Conditions, api.DatabasePaused)
			return db.UID, status
		}, metav1.UpdateOptions{})
		if err != nil {
			return false, err
		}
	}

	backupConfigFound := false
	if e.onlyBackup || resumeAll {
		backupConfigFound, err = ResumeBackupConfiguration(e.stashClient, db.ObjectMeta)
		if err != nil {
			return false, err
		}
	}

	return backupConfigFound, wait.PollUntilContextTimeout(context.Background(), ResumeInterval, ResumeTimeout, true, func(ctx context.Context) (done bool, err error) {
		db, err = e.dbClient.MySQLs(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return false, err
		}

		if db.ObjectMeta.Generation == db.Status.ObservedGeneration {
			return true, nil
		}

		return false, nil
	})
}
