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

package common

import (
	"bytes"
	"context"
	"fmt"

	dbapi "kubedb.dev/apimachinery/apis/kubedb/v1"
	cs "kubedb.dev/apimachinery/client/clientset/versioned"

	cm "github.com/cert-manager/cert-manager/pkg/client/clientset/versioned"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	as "kmodules.xyz/custom-resources/client/clientset/versioned"
)

type MySQLOpts struct {
	DB                *dbapi.MySQL
	DBImage           string
	Config            *rest.Config
	Client            *kubernetes.Clientset
	DBClient          *cs.Clientset
	AppcatClient      *as.Clientset
	CertManagerClient *cm.Clientset
	Username          string
	Pass              string

	ErrWriter *bytes.Buffer
}

func NewMySQLOpts(f cmdutil.Factory, dbName, namespace string) (*MySQLOpts, error) {
	config, err := f.ToRESTConfig()
	if err != nil {
		return nil, err
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	dbClient, err := cs.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	appCatClient, err := as.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	certmanagerClient, err := cm.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	db, err := dbClient.KubedbV1().MySQLs(namespace).Get(context.TODO(), dbName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	if db.Status.Phase != dbapi.DatabasePhaseReady {
		return nil, fmt.Errorf("MySQL %s/%s is not ready", namespace, dbName)
	}

	dbVersion, err := dbClient.CatalogV1alpha1().MySQLVersions().Get(context.TODO(), db.Spec.Version, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	secret, err := client.CoreV1().Secrets(db.Namespace).Get(context.TODO(), db.Spec.AuthSecret.Name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return &MySQLOpts{
		DB:                db,
		DBImage:           dbVersion.Spec.DB.Image,
		Config:            config,
		Client:            client,
		DBClient:          dbClient,
		AppcatClient:      appCatClient,
		CertManagerClient: certmanagerClient,
		Username:          string(secret.Data[corev1.BasicAuthUsernameKey]),
		Pass:              string(secret.Data[corev1.BasicAuthPasswordKey]),
		ErrWriter:         &bytes.Buffer{},
	}, nil
}
