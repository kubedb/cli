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

package dashboard

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func ignoreMongoDBModeSpecificExpressions(unknown map[string]*missingOpts, db *unstructured.Unstructured) map[string]*missingOpts {
	shared := []string{
		"kubedb_com_mongodb_shard_shards",
		"kubedb_com_mongodb_shard_replicas",
		"kubedb_com_mongodb_mongos_replicas",
		"kubedb_com_mongodb_configsvr_replicas",
	}

	general := []string{
		"kubedb_com_mongodb_replicas",
	}

	has := func(shared []string, expr string) bool {
		for _, s := range shared {
			if expr == s {
				return true
			}
		}
		return false
	}

	isShardTopologySet := func(db *unstructured.Unstructured) bool {
		spec, found, err := unstructured.NestedMap(db.Object, "spec")
		if err != nil || !found {
			return false
		}

		shardTopology, found, err := unstructured.NestedMap(spec, "shardTopology")
		if err != nil || !found {
			return false
		}
		return len(shardTopology) > 0
	}

	sharded := isShardTopologySet(db)
	ret := make(map[string]*missingOpts)
	for s, o := range unknown {
		if has(shared, s) && !sharded {
			continue
		} else if has(general, s) && sharded {
			continue
		}
		ret[s] = o
	}
	return ret
}
