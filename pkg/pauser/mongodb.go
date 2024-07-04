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

package pauser

import (
	"context"

	coreapi "kubedb.dev/apimachinery/apis/archiver/v1alpha1"
	"kubedb.dev/apimachinery/apis/kubedb"
	dbapi "kubedb.dev/apimachinery/apis/kubedb/v1"
	cs "kubedb.dev/apimachinery/client/clientset/versioned/typed/kubedb/v1"
	dbutil "kubedb.dev/apimachinery/client/clientset/versioned/typed/kubedb/v1/util"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	kmc "kmodules.xyz/client-go/client"
	condutil "kmodules.xyz/client-go/conditions"
	"sigs.k8s.io/controller-runtime/pkg/client"
	scs "stash.appscode.dev/apimachinery/client/clientset/versioned/typed/stash/v1beta1"
)

type MongoDBPauser struct {
	dbClient     cs.KubedbV1Interface
	stashClient  scs.StashV1beta1Interface
	kc           client.Client
	onlyDb       bool
	onlyBackup   bool
	onlyArchiver bool
}

func NewMongoDBPauser(clientConfig *rest.Config, onlyDb, onlyBackup, onlyArchiver bool) (*MongoDBPauser, error) {
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

	return &MongoDBPauser{
		dbClient:     dbClient,
		stashClient:  stashClient,
		kc:           kc,
		onlyDb:       onlyDb,
		onlyBackup:   onlyBackup,
		onlyArchiver: onlyArchiver,
	}, nil
}

func (e *MongoDBPauser) Pause(name, namespace string) (bool, error) {
	db, err := e.dbClient.MongoDBs(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return false, nil
	}

	pauseAll := !(e.onlyBackup || e.onlyDb || e.onlyArchiver)
	if e.onlyArchiver || pauseAll {
		if err := PauseOrResumeMongoDBArchiver(e.kc, true, db.Spec.Archiver.Ref); err != nil {
			return false, err
		}
		if e.onlyArchiver {
			return false, nil
		}
	}

	if e.onlyDb || pauseAll {
		_, err = dbutil.UpdateMongoDBStatus(context.TODO(), e.dbClient, db.ObjectMeta, func(status *dbapi.MongoDBStatus) (types.UID, *dbapi.MongoDBStatus) {
			status.Conditions = condutil.SetCondition(status.Conditions, condutil.NewCondition(
				kubedb.DatabasePaused,
				"Paused by KubeDB CLI tool",
				db.Generation,
			))
			return db.UID, status
		}, metav1.UpdateOptions{})
		if err != nil {
			return false, nil
		}
	}

	if e.onlyBackup || pauseAll {
		return PauseBackupConfiguration(e.stashClient, db.ObjectMeta)
	}

	return false, nil
}
