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
		use: "elasticsearch", aliases: []string{"es", "elasticsearches"}, short: "Debug helper for Elasticsearch database",
		gvk: schema.GroupVersionKind{Group: dbapi.SchemeGroupVersion.Group, Version: dbapi.SchemeGroupVersion.Version, Kind: dbapi.ResourceKindElasticsearch},
	},
	{
		use: "kafka", aliases: []string{"kf", "kafkas"}, short: "Debug helper for Kafka database",
		gvk: schema.GroupVersionKind{Group: dbapi.SchemeGroupVersion.Group, Version: dbapi.SchemeGroupVersion.Version, Kind: dbapi.ResourceKindKafka},
	},
	{
		use: "mariadb", aliases: []string{"md", "mariadbs"}, short: "Debug helper for MariaDB database",
		gvk: schema.GroupVersionKind{Group: dbapi.SchemeGroupVersion.Group, Version: dbapi.SchemeGroupVersion.Version, Kind: dbapi.ResourceKindMariaDB},
	},
	{
		use: "memcached", aliases: []string{"mc", "memcacheds"}, short: "Debug helper for Memcached database",
		gvk: schema.GroupVersionKind{Group: dbapi.SchemeGroupVersion.Group, Version: dbapi.SchemeGroupVersion.Version, Kind: dbapi.ResourceKindMemcached},
	},
	{
		use: "mongodb", aliases: []string{"mg", "mongodbs"}, short: "Debug helper for MongoDB database",
		gvk: schema.GroupVersionKind{Group: dbapi.SchemeGroupVersion.Group, Version: dbapi.SchemeGroupVersion.Version, Kind: dbapi.ResourceKindMongoDB},
	},
	{
		use: "mysql", aliases: []string{"my", "mysqls"}, short: "Debug helper for MySQL database",
		gvk: schema.GroupVersionKind{Group: dbapi.SchemeGroupVersion.Group, Version: dbapi.SchemeGroupVersion.Version, Kind: dbapi.ResourceKindMySQL},
	},
	{
		use: "perconaxtradb", aliases: []string{"prx", "perconaxtradbs"}, short: "Debug helper for PerconaXtraDB database",
		gvk: schema.GroupVersionKind{Group: dbapi.SchemeGroupVersion.Group, Version: dbapi.SchemeGroupVersion.Version, Kind: dbapi.ResourceKindPerconaXtraDB},
	},
	{
		use: "pgbouncer", aliases: []string{"pb", "pgbouncers"}, short: "Debug helper for PgBouncer database",
		gvk: schema.GroupVersionKind{Group: dbapi.SchemeGroupVersion.Group, Version: dbapi.SchemeGroupVersion.Version, Kind: dbapi.ResourceKindPgBouncer},
	},
	{
		use: "postgres", aliases: []string{"pg", "postgreses"}, short: "Debug helper for Postgres database",
		gvk: schema.GroupVersionKind{Group: dbapi.SchemeGroupVersion.Group, Version: dbapi.SchemeGroupVersion.Version, Kind: dbapi.ResourceKindPostgres},
	},
	{
		use: "proxysql", aliases: []string{"px", "proxysqls"}, short: "Debug helper for ProxySQL database",
		gvk: schema.GroupVersionKind{Group: dbapi.SchemeGroupVersion.Group, Version: dbapi.SchemeGroupVersion.Version, Kind: dbapi.ResourceKindProxySQL},
	},
	{
		use: "redis", aliases: []string{"rd", "redises"}, short: "Debug helper for Redis database",
		gvk: schema.GroupVersionKind{Group: dbapi.SchemeGroupVersion.Group, Version: dbapi.SchemeGroupVersion.Version, Kind: dbapi.ResourceKindRedis},
	},
	// v1alpha2
	{
		use: "cassandra", aliases: []string{"cs", "cassandras"}, short: "Debug helper for Cassandra database",
		gvk: schema.GroupVersionKind{Group: olddbapi.SchemeGroupVersion.Group, Version: olddbapi.SchemeGroupVersion.Version, Kind: olddbapi.ResourceKindCassandra},
	},
	{
		use: "clickhouse", aliases: []string{"cl", "clickhouses"}, short: "Debug helper for ClickHouse database",
		gvk: schema.GroupVersionKind{Group: olddbapi.SchemeGroupVersion.Group, Version: olddbapi.SchemeGroupVersion.Version, Kind: olddbapi.ResourceKindClickHouse},
	},
	{
		use: "druid", aliases: []string{"dr", "druids"}, short: "Debug helper for Druid database",
		gvk: schema.GroupVersionKind{Group: olddbapi.SchemeGroupVersion.Group, Version: olddbapi.SchemeGroupVersion.Version, Kind: olddbapi.ResourceKindDruid},
	},
	{
		use: "ferretdb", aliases: []string{"fr", "ferretdbs"}, short: "Debug helper for FerretDB database",
		gvk: schema.GroupVersionKind{Group: olddbapi.SchemeGroupVersion.Group, Version: olddbapi.SchemeGroupVersion.Version, Kind: olddbapi.ResourceKindFerretDB},
	},
	{
		use: "hazelcast", aliases: []string{"hz", "hazelcasts"}, short: "Debug helper for Hazelcast database",
		gvk: schema.GroupVersionKind{Group: olddbapi.SchemeGroupVersion.Group, Version: olddbapi.SchemeGroupVersion.Version, Kind: olddbapi.ResourceKindHazelcast},
	},
	{
		use: "ignite", aliases: []string{"ig", "ignites"}, short: "Debug helper for Ignite database",
		gvk: schema.GroupVersionKind{Group: olddbapi.SchemeGroupVersion.Group, Version: olddbapi.SchemeGroupVersion.Version, Kind: olddbapi.ResourceKindIgnite},
	},
	{
		use: "mssqlserver", aliases: []string{"ms", "mssqlservers"}, short: "Debug helper for MSSQLServer database",
		gvk: schema.GroupVersionKind{Group: olddbapi.SchemeGroupVersion.Group, Version: olddbapi.SchemeGroupVersion.Version, Kind: olddbapi.ResourceKindMSSQLServer},
	},
	{
		use: "oracle", aliases: []string{"or", "oracles"}, short: "Debug helper for Oracle database",
		gvk: schema.GroupVersionKind{Group: olddbapi.SchemeGroupVersion.Group, Version: olddbapi.SchemeGroupVersion.Version, Kind: olddbapi.ResourceKindOracle},
	},
	{
		use: "pgpool", aliases: []string{"pp", "pgpools"}, short: "Debug helper for Pgpool database",
		gvk: schema.GroupVersionKind{Group: olddbapi.SchemeGroupVersion.Group, Version: olddbapi.SchemeGroupVersion.Version, Kind: olddbapi.ResourceKindPgpool},
	},
	{
		use: "rabbitmq", aliases: []string{"rm", "rabbitmqs"}, short: "Debug helper for Rabbitmq database",
		gvk: schema.GroupVersionKind{Group: olddbapi.SchemeGroupVersion.Group, Version: olddbapi.SchemeGroupVersion.Version, Kind: olddbapi.ResourceKindRabbitmq},
	},
	{
		use: "singlestore", aliases: []string{"sdb", "singlestores"}, short: "Debug helper for Singlestore database",
		gvk: schema.GroupVersionKind{Group: olddbapi.SchemeGroupVersion.Group, Version: olddbapi.SchemeGroupVersion.Version, Kind: olddbapi.ResourceKindSinglestore},
	},
	{
		use: "solr", aliases: []string{"sl", "solrs"}, short: "Debug helper for Solr database",
		gvk: schema.GroupVersionKind{Group: olddbapi.SchemeGroupVersion.Group, Version: olddbapi.SchemeGroupVersion.Version, Kind: olddbapi.ResourceKindSolr},
	},
	{
		use: "zookeeper", aliases: []string{"zk", "zookeepers"}, short: "Debug helper for Zookeeper database",
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
