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

package debug

import (
	"context"
	"fmt"
	"log"

	dbapi "kubedb.dev/apimachinery/apis/kubedb/v1"
	olddbapi "kubedb.dev/apimachinery/apis/kubedb/v1alpha2"

	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type dbEntry struct {
	use     string
	aliases []string
	short   string
	gvk     schema.GroupVersionKind
}

var databases = []dbEntry{
	// v1
	{
		use: dbapi.ResourceSingularElasticsearch, aliases: []string{dbapi.ResourceCodeElasticsearch, dbapi.ResourcePluralElasticsearch}, short: "Debug helper for Elasticsearch database",
		gvk: schema.GroupVersionKind{Group: dbapi.SchemeGroupVersion.Group, Version: dbapi.SchemeGroupVersion.Version, Kind: dbapi.ResourceKindElasticsearch},
	},
	{
		use: dbapi.ResourceSingularKafka, aliases: []string{dbapi.ResourceCodeKafka, dbapi.ResourcePluralKafka}, short: "Debug helper for Kafka database",
		gvk: schema.GroupVersionKind{Group: dbapi.SchemeGroupVersion.Group, Version: dbapi.SchemeGroupVersion.Version, Kind: dbapi.ResourceKindKafka},
	},
	{
		use: dbapi.ResourceSingularMariaDB, aliases: []string{dbapi.ResourceCodeMariaDB, dbapi.ResourcePluralMariaDB}, short: "Debug helper for MariaDB database",
		gvk: schema.GroupVersionKind{Group: dbapi.SchemeGroupVersion.Group, Version: dbapi.SchemeGroupVersion.Version, Kind: dbapi.ResourceKindMariaDB},
	},
	{
		use: dbapi.ResourceSingularMemcached, aliases: []string{dbapi.ResourceCodeMemcached, dbapi.ResourcePluralMemcached}, short: "Debug helper for Memcached database",
		gvk: schema.GroupVersionKind{Group: dbapi.SchemeGroupVersion.Group, Version: dbapi.SchemeGroupVersion.Version, Kind: dbapi.ResourceKindMemcached},
	},
	{
		use: dbapi.ResourceSingularMongoDB, aliases: []string{dbapi.ResourceCodeMongoDB, dbapi.ResourcePluralMongoDB}, short: "Debug helper for MongoDB database",
		gvk: schema.GroupVersionKind{Group: dbapi.SchemeGroupVersion.Group, Version: dbapi.SchemeGroupVersion.Version, Kind: dbapi.ResourceKindMongoDB},
	},
	{
		use: dbapi.ResourceSingularMySQL, aliases: []string{dbapi.ResourceCodeMySQL, dbapi.ResourcePluralMySQL}, short: "Debug helper for MySQL database",
		gvk: schema.GroupVersionKind{Group: dbapi.SchemeGroupVersion.Group, Version: dbapi.SchemeGroupVersion.Version, Kind: dbapi.ResourceKindMySQL},
	},
	{
		use: dbapi.ResourceSingularPerconaXtraDB, aliases: []string{dbapi.ResourceCodePerconaXtraDB, dbapi.ResourcePluralPerconaXtraDB}, short: "Debug helper for PerconaXtraDB database",
		gvk: schema.GroupVersionKind{Group: dbapi.SchemeGroupVersion.Group, Version: dbapi.SchemeGroupVersion.Version, Kind: dbapi.ResourceKindPerconaXtraDB},
	},
	{
		use: dbapi.ResourceSingularPgBouncer, aliases: []string{dbapi.ResourceCodePgBouncer, dbapi.ResourcePluralPgBouncer}, short: "Debug helper for PgBouncer database",
		gvk: schema.GroupVersionKind{Group: dbapi.SchemeGroupVersion.Group, Version: dbapi.SchemeGroupVersion.Version, Kind: dbapi.ResourceKindPgBouncer},
	},
	{
		use: dbapi.ResourceSingularPostgres, aliases: []string{dbapi.ResourceCodePostgres, dbapi.ResourcePluralPostgres}, short: "Debug helper for Postgres database",
		gvk: schema.GroupVersionKind{Group: dbapi.SchemeGroupVersion.Group, Version: dbapi.SchemeGroupVersion.Version, Kind: dbapi.ResourceKindPostgres},
	},
	{
		use: dbapi.ResourceSingularProxySQL, aliases: []string{dbapi.ResourceCodeProxySQL, dbapi.ResourcePluralProxySQL}, short: "Debug helper for ProxySQL database",
		gvk: schema.GroupVersionKind{Group: dbapi.SchemeGroupVersion.Group, Version: dbapi.SchemeGroupVersion.Version, Kind: dbapi.ResourceKindProxySQL},
	},
	{
		use: dbapi.ResourceSingularRedis, aliases: []string{dbapi.ResourceCodeRedis, dbapi.ResourcePluralRedis}, short: "Debug helper for Redis database",
		gvk: schema.GroupVersionKind{Group: dbapi.SchemeGroupVersion.Group, Version: dbapi.SchemeGroupVersion.Version, Kind: dbapi.ResourceKindRedis},
	},
	{
		use: dbapi.ResourceSingularRedisSentinel, aliases: []string{dbapi.ResourceCodeRedisSentinel, dbapi.ResourcePluralRedisSentinel}, short: "Debug helper for RedisSentinel database",
		gvk: schema.GroupVersionKind{Group: dbapi.SchemeGroupVersion.Group, Version: dbapi.SchemeGroupVersion.Version, Kind: dbapi.ResourceKindRedisSentinel},
	},
	// v1alpha2
	{
		use: olddbapi.ResourceSingularAerospike, aliases: []string{olddbapi.ResourceCodeAerospike, olddbapi.ResourcePluralAerospike}, short: "Debug helper for Aerospike database",
		gvk: schema.GroupVersionKind{Group: olddbapi.SchemeGroupVersion.Group, Version: olddbapi.SchemeGroupVersion.Version, Kind: olddbapi.ResourceKindAerospike},
	},
	{
		use: olddbapi.ResourceSingularCassandra, aliases: []string{olddbapi.ResourceCodeCassandra, olddbapi.ResourcePluralCassandra}, short: "Debug helper for Cassandra database",
		gvk: schema.GroupVersionKind{Group: olddbapi.SchemeGroupVersion.Group, Version: olddbapi.SchemeGroupVersion.Version, Kind: olddbapi.ResourceKindCassandra},
	},
	{
		use: olddbapi.ResourceSingularClickHouse, aliases: []string{olddbapi.ResourceCodeClickHouse, olddbapi.ResourcePluralClickHouse}, short: "Debug helper for ClickHouse database",
		gvk: schema.GroupVersionKind{Group: olddbapi.SchemeGroupVersion.Group, Version: olddbapi.SchemeGroupVersion.Version, Kind: olddbapi.ResourceKindClickHouse},
	},
	{
		use: olddbapi.ResourceSingularDB2, aliases: []string{olddbapi.ResourceCodeDB2, olddbapi.ResourcePluralDB2}, short: "Debug helper for DB2 database",
		gvk: schema.GroupVersionKind{Group: olddbapi.SchemeGroupVersion.Group, Version: olddbapi.SchemeGroupVersion.Version, Kind: olddbapi.ResourceKindDB2},
	},
	{
		use: olddbapi.ResourceSingularDocumentDB, aliases: []string{olddbapi.ResourceCodeDocumentDB, olddbapi.ResourcePluralDocumentDB}, short: "Debug helper for DocumentDB database",
		gvk: schema.GroupVersionKind{Group: olddbapi.SchemeGroupVersion.Group, Version: olddbapi.SchemeGroupVersion.Version, Kind: olddbapi.ResourceKindDocumentDB},
	},
	{
		use: olddbapi.ResourceSingularDruid, aliases: []string{olddbapi.ResourceCodeDruid, olddbapi.ResourcePluralDruid}, short: "Debug helper for Druid database",
		gvk: schema.GroupVersionKind{Group: olddbapi.SchemeGroupVersion.Group, Version: olddbapi.SchemeGroupVersion.Version, Kind: olddbapi.ResourceKindDruid},
	},
	{
		use: olddbapi.ResourceSingularFerretDB, aliases: []string{olddbapi.ResourceCodeFerretDB, olddbapi.ResourcePluralFerretDB}, short: "Debug helper for FerretDB database",
		gvk: schema.GroupVersionKind{Group: olddbapi.SchemeGroupVersion.Group, Version: olddbapi.SchemeGroupVersion.Version, Kind: olddbapi.ResourceKindFerretDB},
	},
	{
		use: olddbapi.ResourceSingularHanaDB, aliases: []string{olddbapi.ResourceCodeHanaDB, olddbapi.ResourcePluralHanaDB}, short: "Debug helper for HanaDB database",
		gvk: schema.GroupVersionKind{Group: olddbapi.SchemeGroupVersion.Group, Version: olddbapi.SchemeGroupVersion.Version, Kind: olddbapi.ResourceKindHanaDB},
	},
	{
		use: olddbapi.ResourceSingularHazelcast, aliases: []string{olddbapi.ResourceCodeHazelcast, olddbapi.ResourcePluralHazelcast}, short: "Debug helper for Hazelcast database",
		gvk: schema.GroupVersionKind{Group: olddbapi.SchemeGroupVersion.Group, Version: olddbapi.SchemeGroupVersion.Version, Kind: olddbapi.ResourceKindHazelcast},
	},
	{
		use: olddbapi.ResourceSingularIgnite, aliases: []string{olddbapi.ResourceCodeIgnite, olddbapi.ResourcePluralIgnite}, short: "Debug helper for Ignite database",
		gvk: schema.GroupVersionKind{Group: olddbapi.SchemeGroupVersion.Group, Version: olddbapi.SchemeGroupVersion.Version, Kind: olddbapi.ResourceKindIgnite},
	},
	{
		use: olddbapi.ResourceSingularMilvus, aliases: []string{olddbapi.ResourceCodeMilvus, olddbapi.ResourcePluralMilvus}, short: "Debug helper for Milvus database",
		gvk: schema.GroupVersionKind{Group: olddbapi.SchemeGroupVersion.Group, Version: olddbapi.SchemeGroupVersion.Version, Kind: olddbapi.ResourceKindMilvus},
	},
	{
		use: olddbapi.ResourceSingularMSSQLServer, aliases: []string{olddbapi.ResourceCodeMSSQLServer, olddbapi.ResourcePluralMSSQLServer}, short: "Debug helper for MSSQLServer database",
		gvk: schema.GroupVersionKind{Group: olddbapi.SchemeGroupVersion.Group, Version: olddbapi.SchemeGroupVersion.Version, Kind: olddbapi.ResourceKindMSSQLServer},
	},
	{
		use: olddbapi.ResourceSingularNeo4j, aliases: []string{olddbapi.ResourceCodeNeo4j, olddbapi.ResourcePluralNeo4j}, short: "Debug helper for Neo4j database",
		gvk: schema.GroupVersionKind{Group: olddbapi.SchemeGroupVersion.Group, Version: olddbapi.SchemeGroupVersion.Version, Kind: olddbapi.ResourceKindNeo4j},
	},
	{
		use: olddbapi.ResourceSingularOracle, aliases: []string{olddbapi.ResourceCodeOracle, olddbapi.ResourcePluralOracle}, short: "Debug helper for Oracle database",
		gvk: schema.GroupVersionKind{Group: olddbapi.SchemeGroupVersion.Group, Version: olddbapi.SchemeGroupVersion.Version, Kind: olddbapi.ResourceKindOracle},
	},
	{
		use: olddbapi.ResourceSingularPgpool, aliases: []string{olddbapi.ResourceCodePgpool, olddbapi.ResourcePluralPgpool}, short: "Debug helper for Pgpool database",
		gvk: schema.GroupVersionKind{Group: olddbapi.SchemeGroupVersion.Group, Version: olddbapi.SchemeGroupVersion.Version, Kind: olddbapi.ResourceKindPgpool},
	},
	{
		use: olddbapi.ResourceSingularQdrant, aliases: []string{olddbapi.ResourceCodeQdrant, olddbapi.ResourcePluralQdrant}, short: "Debug helper for Qdrant database",
		gvk: schema.GroupVersionKind{Group: olddbapi.SchemeGroupVersion.Group, Version: olddbapi.SchemeGroupVersion.Version, Kind: olddbapi.ResourceKindQdrant},
	},
	{
		use: olddbapi.ResourceSingularRabbitmq, aliases: []string{olddbapi.ResourceCodeRabbitmq, olddbapi.ResourcePluralRabbitmq}, short: "Debug helper for RabbitMQ database",
		gvk: schema.GroupVersionKind{Group: olddbapi.SchemeGroupVersion.Group, Version: olddbapi.SchemeGroupVersion.Version, Kind: olddbapi.ResourceKindRabbitmq},
	},
	{
		use: olddbapi.ResourceSingularSinglestore, aliases: []string{olddbapi.ResourceCodeSinglestore, olddbapi.ResourcePluralSinglestore}, short: "Debug helper for Singlestore database",
		gvk: schema.GroupVersionKind{Group: olddbapi.SchemeGroupVersion.Group, Version: olddbapi.SchemeGroupVersion.Version, Kind: olddbapi.ResourceKindSinglestore},
	},
	{
		use: olddbapi.ResourceSingularSolr, aliases: []string{olddbapi.ResourceCodeSolr, olddbapi.ResourcePluralSolr}, short: "Debug helper for Solr database",
		gvk: schema.GroupVersionKind{Group: olddbapi.SchemeGroupVersion.Group, Version: olddbapi.SchemeGroupVersion.Version, Kind: olddbapi.ResourceKindSolr},
	},
	{
		use: olddbapi.ResourceSingularWeaviate, aliases: []string{olddbapi.ResourceCodeWeaviate, olddbapi.ResourcePluralWeaviate}, short: "Debug helper for Weaviate database",
		gvk: schema.GroupVersionKind{Group: olddbapi.SchemeGroupVersion.Group, Version: olddbapi.SchemeGroupVersion.Version, Kind: olddbapi.ResourceKindWeaviate},
	},
	{
		use: olddbapi.ResourceSingularZooKeeper, aliases: []string{olddbapi.ResourceCodeZooKeeper, olddbapi.ResourcePluralZooKeeper}, short: "Debug helper for ZooKeeper database",
		gvk: schema.GroupVersionKind{Group: olddbapi.SchemeGroupVersion.Group, Version: olddbapi.SchemeGroupVersion.Version, Kind: olddbapi.ResourceKindZooKeeper},
	},
}

func AllCMDs(f cmdutil.Factory) []*cobra.Command {
	cmds := make([]*cobra.Command, 0, len(databases))
	for _, entry := range databases {
		cmds = append(cmds, buildDebugCMD(f, entry))
	}
	return cmds
}

func buildDebugCMD(f cmdutil.Factory, entry dbEntry) *cobra.Command {
	var operatorNamespace string

	cmd := &cobra.Command{
		Use:     entry.use,
		Aliases: entry.aliases,
		Short:   entry.short,
		Example: fmt.Sprintf(`kubectl dba debug %s -n demo sample-%s --operator-namespace kubedb`, entry.use, entry.use),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				log.Fatalf("Enter %s object's name as an argument", entry.use)
			}
			dbName := args[0]

			namespace, _, err := f.ToRawKubeConfigLoader().Namespace()
			if err != nil {
				klog.Error(err, "failed to get current namespace")
			}

			opts, err := newDBOpts(f, entry.gvk, dbName, namespace, operatorNamespace)
			if err != nil {
				log.Fatalln(err)
			}

			obj, err := opts.kc.Scheme().New(entry.gvk)
			if err != nil {
				log.Fatalln(err)
			}
			db, ok := obj.(client.Object)
			if !ok {
				log.Fatalf("%T does not implement client.Object", obj)
			}
			if err = opts.kc.Get(context.TODO(), types.NamespacedName{Name: dbName, Namespace: namespace}, db); err != nil {
				log.Fatalln(err)
			}
			if err = opts.run(db); err != nil {
				log.Fatalln(err)
			}
		},
	}
	cmd.Flags().StringVarP(&operatorNamespace, "operator-namespace", "o", "kubedb", "the namespace where the kubedb operator is installed")
	return cmd
}
