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

package data

import (
	"context"
	"fmt"
	"log"

	api "kubedb.dev/apimachinery/apis/kubedb/v1alpha2"
	cs "kubedb.dev/apimachinery/client/clientset/versioned"
	"kubedb.dev/apimachinery/pkg/factory"
	"kubedb.dev/cli/pkg/lib"
	es_clientgo "kubedb.dev/db-client-go/elasticsearch"

	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"kmodules.xyz/client-go/tools/portforward"
)

const (
	idx = "kubedb-elasticsearch-index"
)

type elasticsearchOpts struct {
	db       *api.Elasticsearch
	client   *kubernetes.Clientset
	esClient *es_clientgo.Client
	username string
	pass     string
}

func InsertElasticsearchDataCMD(f cmdutil.Factory) *cobra.Command {
	var (
		dbName string
		rows   int
	)

	esInsertCmd := &cobra.Command{
		Use: "elasticsearch",
		Aliases: []string{
			"es",
		},
		Short:   "Insert data to a elasticsearch database",
		Long:    `Use this cmd to insert data into a elasticsearch database`,
		Example: `kubectl dba insert -n demo es es-quickstart --rows 1000`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				log.Fatal("Enter elasticsearch object's name as an argument")
			}
			dbName = args[0]

			namespace, _, err := f.ToRawKubeConfigLoader().Namespace()
			if err != nil {
				klog.Error(err, "failed to get current namespace")
			}

			config, err := f.ToRESTConfig()
			if err != nil {
				log.Fatal("couldn't create config, error: ", err)
			}

			tunnel, err := lib.TunnelToDBService(config, dbName, namespace, api.ElasticsearchRestPort)
			if err != nil {
				log.Fatal("couldn't create tunnel, error: ", err)
			}
			defer tunnel.Close()

			opts, err := newElasticsearchOpts(config, dbName, namespace, tunnel)
			if err != nil {
				log.Fatalln(err)
			}

			if rows <= 0 {
				log.Fatal("rows need to be greater than 0")
			}

			err = opts.insertDataInDatabase(rows)
			if err != nil {
				log.Fatal(err)
			}
		},
	}

	esInsertCmd.Flags().IntVarP(&rows, "rows", "r", 100, "number of rows to insert")

	return esInsertCmd
}

func VerifyElasticsearchDataCMD(f cmdutil.Factory) *cobra.Command {
	var (
		dbName string
		rows   int
	)

	esVerifyCmd := &cobra.Command{
		Use: "elasticsearch",
		Aliases: []string{
			"es",
		},
		Short:   "Verify rows in a elasticsearch database",
		Long:    `Use this cmd to verify data in a elasticsearch object`,
		Example: `kubectl dba verify -n demo es es-quickstart  --rows 1000`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				log.Fatal("Enter elasticsearch object's name as an argument")
			}
			dbName = args[0]

			config, err := f.ToRESTConfig()
			if err != nil {
				log.Fatal("couldn't create config, error: ", err)
			}

			namespace, _, err := f.ToRawKubeConfigLoader().Namespace()
			if err != nil {
				klog.Error(err, "failed to get current namespace")
			}

			tunnel, err := lib.TunnelToDBService(config, dbName, namespace, api.ElasticsearchRestPort)
			if err != nil {
				log.Fatal("couldn't create tunnel, error: ", err)
			}
			defer tunnel.Close()

			opts, err := newElasticsearchOpts(config, dbName, namespace, tunnel)
			if err != nil {
				log.Fatalln(err)
			}

			if rows <= 0 {
				log.Fatal("rows need to be greater than 0")
			}

			err = opts.verifyElasticsearchData(rows)
			if err != nil {
				log.Fatal(err)
			}
		},
	}

	esVerifyCmd.Flags().IntVarP(&rows, "rows", "r", 100, "number of rows to verify")

	return esVerifyCmd
}

func DropElasticsearchDataCMD(f cmdutil.Factory) *cobra.Command {
	var dbName string

	esDropCmd := &cobra.Command{
		Use: "elasticsearch",
		Aliases: []string{
			"es",
		},
		Short:   "Delete data from elasticsearch database",
		Long:    `Use this cmd to delete inserted data in a elasticsearch object`,
		Example: `kubectl dba drop -n demo es es-quickstart`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				log.Fatal("Enter elasticsearch object's name as an argument. Your commands will be applied on a database inside it's primary pod")
			}
			dbName = args[0]

			config, err := f.ToRESTConfig()
			if err != nil {
				log.Fatal("couldn't create config, error: ", err)
			}

			namespace, _, err := f.ToRawKubeConfigLoader().Namespace()
			if err != nil {
				klog.Error(err, "failed to get current namespace")
			}

			tunnel, err := lib.TunnelToDBService(config, dbName, namespace, api.ElasticsearchRestPort)
			if err != nil {
				log.Fatal("couldn't create tunnel, error: ", err)
			}
			defer tunnel.Close()

			opts, err := newElasticsearchOpts(config, dbName, namespace, tunnel)
			if err != nil {
				log.Fatalln(err)
			}

			err = opts.dropElasticsearchData()
			if err != nil {
				log.Fatal(err)
			}
		},
	}

	return esDropCmd
}

func newElasticsearchOpts(config *rest.Config, dbName, namespace string, tunnel *portforward.Tunnel) (*elasticsearchOpts, error) {
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	dbClient, err := cs.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	db, err := dbClient.KubedbV1alpha2().Elasticsearches(namespace).Get(context.TODO(), dbName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	if db.Status.Phase != api.DatabasePhaseReady {
		return nil, fmt.Errorf("elasticsearch %s/%s is not ready", namespace, dbName)
	}

	secret, err := client.CoreV1().Secrets(db.Namespace).Get(context.TODO(), db.Spec.AuthSecret.Name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	kcc, err := factory.NewUncachedClient(config)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%v://127.0.0.1:%d", db.GetConnectionScheme(), tunnel.Local)

	esClient, err := es_clientgo.NewKubeDBClientBuilder(kcc, db).WithURL(url).WithContext(context.TODO()).GetElasticClient()
	if err != nil {
		return nil, err
	}

	return &elasticsearchOpts{
		db:       db,
		client:   client,
		esClient: esClient,
		username: string(secret.Data[corev1.BasicAuthUsernameKey]),
		pass:     string(secret.Data[corev1.BasicAuthPasswordKey]),
	}, nil
}

func (opts *elasticsearchOpts) insertDataInDatabase(rows int) error {
	err := opts.esClient.IndexExistsOrNot(idx)
	if err == nil {
		fmt.Printf("Index named %s already exists", idx)
		return err
	}
	err = opts.esClient.CreateIndex(idx)
	if err != nil {
		fmt.Printf("Failed to create index %s for elasticsearch %s/%s", idx, opts.db.Namespace, opts.db.Name)
		return err
	}
	for i := 0; i < rows; i++ {
		name := fmt.Sprintf("document index %d", i+1)

		id := fmt.Sprintf("%d", i+1)
		body := map[string]interface{}{
			"name": name,
		}
		if err := opts.esClient.PutData(idx, id, body); err != nil {
			fmt.Printf("Failed to insert data in the index %s for elasticsearch %s/%s", idx, opts.db.Namespace, opts.db.Name)
			return err
		}
	}

	fmt.Printf("\n%d keys inserted in index %s of elasticsearch database %s/%s successfully\n", rows, idx, opts.db.Namespace, opts.db.Name)

	return nil
}

func (opts *elasticsearchOpts) verifyElasticsearchData(rows int) error {
	totalKeys, err := opts.esClient.CountData(idx)
	if err != nil {
		return err
	}
	if totalKeys == rows {
		fmt.Printf("\nSuccess! Elasticsearch database %s/%s contains: %d keys\n", opts.db.Namespace, opts.db.Name, totalKeys)
	} else {
		fmt.Printf("\nError! Expected keys: %d .Elasticsearch database %s/%s contains: %d keys\n", rows, opts.db.Namespace, opts.db.Name, totalKeys)
	}
	return nil
}

func (opts *elasticsearchOpts) dropElasticsearchData() error {
	err := opts.esClient.DeleteIndex(idx)
	if err != nil {
		fmt.Printf("Error. Can not drop data from Elasticsearch %s/%s\n", opts.db.Namespace, opts.db.Name)
		return err
	}
	fmt.Printf("\nSuccess: All the CLI inserted documents DELETED from elasticsearch %s/%s and deleted index %s\n", opts.db.Namespace, opts.db.Name, idx)
	return nil
}
