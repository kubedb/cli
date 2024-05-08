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

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kmapi "kmodules.xyz/client-go/api/v1"
	kmc "kmodules.xyz/client-go/client"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func PauseOrResumeMySQLArchiver(klient client.Client, value bool, reference kmapi.ObjectReference) error {
	name := reference.Name
	namespace := reference.Namespace
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

func PauseOrResumeMariaDBArchiver(klient client.Client, value bool, reference kmapi.ObjectReference) error {
	name := reference.Name
	namespace := reference.Namespace
	archiver, err := getMariaDBArchiver(klient, kmapi.ObjectReference{
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
			in := obj.(*coreapi.MariaDBArchiver)
			in.Spec.Pause = value
			return in
		},
	)
	return err
}

func getMariaDBArchiver(klient client.Client, ref kmapi.ObjectReference) (*coreapi.MariaDBArchiver, error) {
	archiver := &coreapi.MariaDBArchiver{
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

func PauseOrResumePostgresArchiver(klient client.Client, value bool, reference kmapi.ObjectReference) error {
	name := reference.Name
	namespace := reference.Namespace
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

func PauseOrResumeMongoDBArchiver(klient client.Client, value bool, reference kmapi.ObjectReference) error {
	name := reference.Name
	namespace := reference.Namespace
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
