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

package credentials

import (
	"context"
	"fmt"

	cs "kubedb.dev/apimachinery/client/clientset/versioned/typed/kubedb/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type MySQLShowCred struct {
	client   kubernetes.Interface
	dbClient cs.KubedbV1Interface
}

func NewMySQLShowCred(clientConfig *rest.Config) (*MySQLShowCred, error) {
	dbClient, err := cs.NewForConfig(clientConfig)
	if err != nil {
		return nil, err
	}
	client, err := kubernetes.NewForConfig(clientConfig)
	if err != nil {
		return nil, err
	}

	return &MySQLShowCred{
		dbClient: dbClient,
		client:   client,
	}, nil
}

func (e *MySQLShowCred) GetCred(name, namespace string) (map[string][]byte, error) {
	db, err := e.dbClient.MySQLs(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	if db.Spec.AuthSecret == nil {
		return nil, fmt.Errorf("auth secret can't be empty")
	}

	authSecret, err := e.client.CoreV1().Secrets(db.Namespace).Get(context.TODO(), db.Spec.AuthSecret.Name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return authSecret.Data, nil
}
