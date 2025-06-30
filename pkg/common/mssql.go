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
	"context"
	"fmt"

	// Import the correct MSSQLServer API version

	dbapi "kubedb.dev/apimachinery/apis/kubedb/v1alpha2"
	cs "kubedb.dev/apimachinery/client/clientset/versioned"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
)

// MSSQLOpts holds clients and the fetched MSSQLServer object for a command.
type MSSQLOpts struct {
	DB       *dbapi.MSSQLServer
	Config   *rest.Config
	Client   *kubernetes.Clientset
	DBClient *cs.Clientset
}

// NewMSSQLOpts creates a new MSSQLOpts instance, fetches the MSSQLServer CR,
// and performs initial validation.
func NewMSSQLOpts(f cmdutil.Factory, dbName, namespace string) (*MSSQLOpts, error) {
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

	// Fetch the source MSSQLServer custom resource
	mssql, err := dbClient.KubedbV1alpha2().MSSQLServers(namespace).Get(context.TODO(), dbName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	// IMPORTANT VALIDATION: Check if the database is in a state
	// where it has generated the necessary DAG secrets.
	if mssql.Status.Phase != dbapi.DatabasePhaseReady {
		return nil, fmt.Errorf("source MSSQLServer %s/%s is not ready (current phase: %s)", namespace, dbName, mssql.Status.Phase)
	}

	return &MSSQLOpts{
		DB:       mssql,
		Config:   config,
		Client:   client,
		DBClient: dbClient,
	}, nil
}
