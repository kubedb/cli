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

func ignoreElasticsearchModeSpecificExpressions(unknown map[string]*missingOpts, db *unstructured.Unstructured) map[string]*missingOpts {
	master := []string{
		"kubedb_com_elasticsearch_master_node_replicas",
		"kubedb_com_elasticsearch_master_node_storage_class_info",
		"kubedb_com_elasticsearch_master_node_max_unavailable",
	}
	data := []string{
		"kubedb_com_elasticsearch_data_node_replicas",
		"kubedb_com_elasticsearch_data_node_storage_class_info",
		"kubedb_com_elasticsearch_data_node_max_unavailable",
	}
	datacontent := []string{
		"kubedb_com_elasticsearch_datacontent_node_replicas",
		"kubedb_com_elasticsearch_datacontent_node_storage_class_info",
		"kubedb_com_elasticsearch_datacontent_node_max_unavailable",
	}
	datahot := []string{
		"kubedb_com_elasticsearch_datahot_node_replicas",
		"kubedb_com_elasticsearch_datahot_node_storage_class_info",
		"kubedb_com_elasticsearch_datahot_node_max_unavailable",
	}
	datawarm := []string{
		"kubedb_com_elasticsearch_datawarm_node_replicas",
		"kubedb_com_elasticsearch_datawarm_node_storage_class_info",
		"kubedb_com_elasticsearch_datawarm_node_max_unavailable",
	}
	datacold := []string{
		"kubedb_com_elasticsearch_datacold_node_replicas",
		"kubedb_com_elasticsearch_datacold_node_storage_class_info",
		"kubedb_com_elasticsearch_datacold_node_max_unavailable",
	}
	datafrozen := []string{
		"kubedb_com_elasticsearch_datafrozen_node_replicas",
		"kubedb_com_elasticsearch_datafrozen_node_storage_class_info",
		"kubedb_com_elasticsearch_datafrozen_node_max_unavailable",
	}
	ingest := []string{
		"kubedb_com_elasticsearch_ingest_node_replicas",
		"kubedb_com_elasticsearch_ingest_node_storage_class_info",
		"kubedb_com_elasticsearch_ingest_node_max_unavailable",
	}
	ml := []string{
		"kubedb_com_elasticsearch_ml_node_replicas",
		"kubedb_com_elasticsearch_ml_node_storage_class_info",
		"kubedb_com_elasticsearch_ml_node_max_unavailable",
	}
	transform := []string{
		"kubedb_com_elasticsearch_transform_node_replicas",
		"kubedb_com_elasticsearch_transform_node_storage_class_info",
		"kubedb_com_elasticsearch_transform_node_max_unavailable",
	}
	coordinating := []string{
		"kubedb_com_elasticsearch_coordinating_node_replicas",
		"kubedb_com_elasticsearch_coordinating_node_storage_class_info",
		"kubedb_com_elasticsearch_coordinating_node_max_unavailable",
	}

	has := func(shared []string, expr string) bool {
		for _, s := range shared {
			if expr == s {
				return true
			}
		}
		return false
	}

	isSet := func(db *unstructured.Unstructured, typ string) bool {
		spec, found, err := unstructured.NestedMap(db.Object, "spec")
		if err != nil || !found {
			return false
		}

		topo, found, err := unstructured.NestedMap(spec, "topology")
		if err != nil || !found {
			return false
		}

		typed, found, err := unstructured.NestedMap(topo, typ)
		if err != nil || !found {
			return false
		}
		return len(typed) > 0
	}

	ret := make(map[string]*missingOpts)
	for s, o := range unknown {
		if has(master, s) && !isSet(db, "master") {
			continue
		}
		if has(data, s) && !isSet(db, "data") {
			continue
		}
		if has(datacontent, s) && !isSet(db, "dataContent") {
			continue
		}
		if has(datahot, s) && !isSet(db, "dataHot") {
			continue
		}
		if has(datawarm, s) && !isSet(db, "dataWarm") {
			continue
		}
		if has(datacold, s) && !isSet(db, "dataCold") {
			continue
		}
		if has(datafrozen, s) && !isSet(db, "dataFrozen") {
			continue
		}
		if has(ingest, s) && !isSet(db, "ingest") {
			continue
		}
		if has(ml, s) && !isSet(db, "ml") {
			continue
		}
		if has(transform, s) && !isSet(db, "transform") {
			continue
		}
		if has(coordinating, s) && !isSet(db, "coordinating") {
			continue
		}
		ret[s] = o
	}
	return ret
}
