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

	coreapi "kubedb.dev/apimachinery/apis/archiver/v1alpha1"
	"kubedb.dev/apimachinery/apis/kubedb"
	dbapi "kubedb.dev/apimachinery/apis/kubedb/v1"
	cs "kubedb.dev/apimachinery/client/clientset/versioned/typed/kubedb/v1"
	dbutil "kubedb.dev/apimachinery/client/clientset/versioned/typed/kubedb/v1/util"
	pautil "kubedb.dev/cli/pkg/pauser"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/rest"
	kmc "kmodules.xyz/client-go/client"
	condutil "kmodules.xyz/client-go/conditions"
	"sigs.k8s.io/controller-runtime/pkg/client"
	scs "stash.appscode.dev/apimachinery/client/clientset/versioned/typed/stash/v1beta1"
)

type MySQLResumer struct {
	dbClient     cs.KubedbV1Interface
	stashClient  scs.StashV1beta1Interface
	kc           client.Client
	onlyDb       bool
	onlyBackup   bool
	onlyArchiver bool
}

func NewMySQLResumer(clientConfig *rest.Config, onlyDb, onlyBackup, onlyArchiver bool) (*MySQLResumer, error) {
	dbClient, err := cs.NewForConfig(clientConfig)
	if err != nil {
		return nil, err
	}
	stashClient, err := scs.NewForConfig(clientConfig)
	if err != nil {
		return nil, err
	}
	kc, err := kmc.NewUncachedClient(clientConfig, coreapi.AddToScheme)
	if err != nil {
		return nil, err
	}
	return &MySQLResumer{
		dbClient:     dbClient,
		stashClient:  stashClient,
		kc:           kc,
		onlyDb:       onlyDb,
		onlyBackup:   onlyBackup,
		onlyArchiver: onlyArchiver,
	}, nil
}

func (e *MySQLResumer) Resume(name, namespace string) (bool, error) {
	db, err := e.dbClient.MySQLs(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return false, err
	}

	resumeAll := !(e.onlyBackup || e.onlyDb || e.onlyArchiver)

	if (e.onlyArchiver || resumeAll) && db.Spec.Archiver != nil {
		if err := pautil.PauseOrResumeMySQLArchiver(e.kc, false, db.Spec.Archiver.Ref); err != nil {
			return false, err
		}
		if e.onlyArchiver {
			return false, nil
		}
	}

	if e.onlyDb || resumeAll {
		_, err = dbutil.UpdateMySQLStatus(context.TODO(), e.dbClient, db.ObjectMeta, func(status *dbapi.MySQLStatus) (types.UID, *dbapi.MySQLStatus) {
			status.Conditions = condutil.RemoveCondition(status.Conditions, kubedb.DatabasePaused)
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
