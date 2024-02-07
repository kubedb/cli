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
	"fmt"

	coreapi "kubedb.dev/apimachinery/apis/archiver/v1alpha1"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kmapi "kmodules.xyz/client-go/api/v1"
	kmc "kmodules.xyz/client-go/client"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"stash.appscode.dev/apimachinery/apis"
	stash "stash.appscode.dev/apimachinery/apis/stash/v1beta1"
	scs "stash.appscode.dev/apimachinery/client/clientset/versioned/typed/stash/v1beta1"
	scsutil "stash.appscode.dev/apimachinery/client/clientset/versioned/typed/stash/v1beta1/util"
)

func NewUncachedClient() (client.Client, error) {
	cfg, err := ctrl.GetConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get Kubernetes config. Reason: %w", err)
	}
	return kmc.NewUncachedClient(
		cfg,
		coreapi.AddToScheme,
	)
}

func PauseBackupConfiguration(stashClient scs.StashV1beta1Interface, dbMeta metav1.ObjectMeta) (bool, error) {
	configs, err := stashClient.BackupConfigurations(dbMeta.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return false, nil
		}
		return false, err
	}

	var dbBackupConfig *stash.BackupConfiguration
	for _, config := range configs.Items {
		if config.Spec.Target.Ref.Name == dbMeta.Name && config.Spec.Target.Ref.Kind == apis.KindAppBinding {
			dbBackupConfig = &config
			break
		}
	}

	if dbBackupConfig != nil && !dbBackupConfig.Spec.Paused {
		_, err := scsutil.TryUpdateBackupConfiguration(context.TODO(), stashClient, dbBackupConfig.ObjectMeta, func(configuration *stash.BackupConfiguration) *stash.BackupConfiguration {
			configuration.Spec.Paused = true
			return configuration
		}, metav1.UpdateOptions{})
		if err != nil {
			return false, err
		}
	}
	return dbBackupConfig != nil, nil
}

func PauseMySQLArchiver(value bool, name string, namespace string) error {
	var klient client.Client
	klient, err := NewUncachedClient()
	if err != nil {
		return err
	}
	archiver, err := getMysqlArchiver(klient, kmapi.ObjectReference{
		Name:      name,
		Namespace: namespace,
	})
	if err != nil {
		return err
	}
	_, err = kmc.CreateOrPatch(
		context.Background(),
		klient,
		archiver,
		func(obj client.Object, createOp bool) client.Object {
			in := obj.(*coreapi.MySQLArchiver)
			in.Spec.Pause = value
			return in
		},
	)
	return err
}

func getMysqlArchiver(klient client.Client, ref kmapi.ObjectReference) (*coreapi.MySQLArchiver, error) {
	archiver := &coreapi.MySQLArchiver{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ref.Name,
			Namespace: ref.Namespace,
		},
	}
	if err := klient.Get(context.Background(), client.ObjectKeyFromObject(archiver), archiver); err != nil {
		return nil, err
	}
	return archiver, nil
}

func PausePostgresArchiver(value bool, name string, namespace string) error {
	var klient client.Client
	klient, err := NewUncachedClient()
	if err != nil {
		return err
	}
	archiver, err := getPostgresArchiver(klient, kmapi.ObjectReference{
		Name:      name,
		Namespace: namespace,
	})
	if err != nil {
		return err
	}
	_, err = kmc.CreateOrPatch(
		context.Background(),
		klient,
		archiver,
		func(obj client.Object, createOp bool) client.Object {
			in := obj.(*coreapi.PostgresArchiver)
			in.Spec.Pause = value
			return in
		},
	)
	return err
}

func getPostgresArchiver(klient client.Client, ref kmapi.ObjectReference) (*coreapi.PostgresArchiver, error) {
	archiver := &coreapi.PostgresArchiver{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ref.Name,
			Namespace: ref.Namespace,
		},
	}
	if err := klient.Get(context.Background(), client.ObjectKeyFromObject(archiver), archiver); err != nil {
		return nil, err
	}
	return archiver, nil
}

func PauseMongoDBArchiver(value bool, name string, namespace string) error {
	var klient client.Client
	klient, err := NewUncachedClient()
	if err != nil {
		return err
	}
	archiver, err := getMongoDBArchiver(klient, kmapi.ObjectReference{
		Name:      name,
		Namespace: namespace,
	})
	if err != nil {
		return err
	}
	_, err = kmc.CreateOrPatch(
		context.Background(),
		klient,
		archiver,
		func(obj client.Object, createOp bool) client.Object {
			in := obj.(*coreapi.MongoDBArchiver)
			in.Spec.Pause = value
			return in
		},
	)
	return err
}

func getMongoDBArchiver(klient client.Client, ref kmapi.ObjectReference) (*coreapi.MongoDBArchiver, error) {
	archiver := &coreapi.MongoDBArchiver{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ref.Name,
			Namespace: ref.Namespace,
		},
	}
	if err := klient.Get(context.Background(), client.ObjectKeyFromObject(archiver), archiver); err != nil {
		return nil, err
	}
	return archiver, nil
}
