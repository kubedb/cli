/*
Copyright AppsCode Inc. and Contributors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package elasticsearch

import (
	"context"
	"encoding/json"

	dbapi "kubedb.dev/apimachinery/apis/kubedb/v1"

	core "k8s.io/api/core/v1"
)

const (
	writeRequestIndex = "kubedb-system"
	writeRequestID    = "info"
	writeRequestType  = "_doc"
	CustomRoleName    = "readWriteAnyDatabase"
	ApplicationKibana = "kibana-.kibana"
)

const (
	PrivilegeCreateSnapshot = "create_snapshot"
	PrivilegeManage         = "manage"
	PrivilegeManageILM      = "manage_ilm"
	PrivilegeManageRoleup   = "manage_rollup"
	PrivilegeMonitor        = "monitor"
	PrivilegeManageCCR      = "manage_ccr"
	PrivilegeRead           = "read"
	PrivilegeWrite          = "write"
	PrivilegeCreateIndex    = "create_index"
	PrivilegeIndexAny       = "*"
)

const (
	DisableShardAllocation = `
{
  "persistent": {
    "cluster.routing.allocation.enable": "primaries"
  }
}`
	ReEnableShardAllocation = `
{
  "persistent": {
    "cluster.routing.allocation.enable": null
  }
}`
	ExcludeNodeAllocation = `
{
  "transient" : {
    "cluster.routing.allocation.exclude._name" : "{{.}}"
  }
}`
	DeleteNodeAllocationExclusion = `
{
  "transient" : {
    "cluster.routing.allocation.exclude._name" : null
  }
}`
	VotingExclusionUrl = "/_cluster/voting_config_exclusions"
)

type Info struct {
	Name        string  `json:"name"`
	ClusterName string  `json:"cluster_name"`
	ClusterUuid string  `json:"cluster_uuid"`
	Version     Version `json:"version"`
}
type Version struct {
	Number string `json:"number"`
}

// NodesStats is used for calculation how much storage is in use in specific data nodes
type NodesStats struct {
	Nodes map[string]struct {
		Indices struct {
			Store struct {
				SizeInBytes int64 `json:"size_in_bytes"`
			} `json:"store"`
		} `json:"indices"`
	} `json:"nodes"`
}

type NodeInfo struct {
	Nodes NodeSummary `json:"_nodes"`
}

type NodeSummary struct {
	Total      json.Number `json:"total"`
	Successful json.Number `json:"successful"`
	Failed     json.Number `json:"failed"`
}

type IndexSetting struct {
	Index Index `json:"index"`
}

type Index struct {
	NumberOfReplicas string `json:"number_of_replicas"`
}

type IndexDistribution struct {
	Index string `json:"index"`
	Node  string `json:"node"`
}

// By default searchGuard/openDistro maintain a copy this index's shard in every data node.
// -->: "auto_expand_replicas" : "0-all"
var IgnorableIndexList = []string{
	"searchguard",
	".opendistro_security",
}

// IsIgnorableIndex returns true if the index shards are ignorable.
func IsIgnorableIndex(index string) bool {
	for _, i := range IgnorableIndexList {
		if i == index {
			return true
		}
	}
	return false
}

type DBPrivileges struct {
	Names                  []string `json:"names"`
	Privileges             []string `json:"privileges"`
	AllowRestrictedIndices bool     `json:"allow_restricted_indices"`
}

type ApplicationPrivileges struct {
	Application string   `json:"application"`
	Privileges  []string `json:"privileges"`
	Resources   []string `json:"resources"`
}

type TransientMetaPrivileges struct {
	Enabled bool `json:"enabled"`
}

type UserRoleReq struct {
	Cluster           []string                `json:"cluster"`
	Indices           []DBPrivileges          `json:"indices"`
	Applications      []ApplicationPrivileges `json:"applications"`
	RunAs             []string                `json:"run_as"`
	TransientMetaData TransientMetaPrivileges `json:"transient_metadata"`
}

type WriteRequestIndex struct {
	Index WriteRequestIndexBody `json:"index"`
}

type WriteRequestIndexBody struct {
	ID   string `json:"_id"`
	Type string `json:"_type,omitempty"`
}

type ESClient interface {
	ClusterHealthInfo() (map[string]any, error)
	ClusterStatus() (string, error)
	CountData(index string) (int, error)
	CreateDBUserRole(ctx context.Context) error
	CreateIndex(index string) error
	DeleteIndex(index string) error
	GetIndicesInfo() ([]any, error)
	GetClusterWriteStatus(ctx context.Context, db *dbapi.Elasticsearch) error
	GetClusterReadStatus(ctx context.Context, db *dbapi.Elasticsearch) error
	GetTotalDiskUsage(ctx context.Context) (string, error)
	GetDBUserRole(ctx context.Context) (error, bool)
	IndexExistsOrNot(index string) error
	NodesStats() (map[string]any, error)
	ShardStats() ([]ShardInfo, error)
	PutData(index, id string, data map[string]any) error
	SyncCredentialFromSecret(secret *core.Secret) error
	DisableShardAllocation() error
	ReEnableShardAllocation() error
	CheckVersion() (string, error)
	GetClusterStatus() (string, error)
	CountIndex() (int, error)
	GetData(_index, _type, _id string) (map[string]any, error)
	CountNodes() (int64, error)
	AddVotingConfigExclusions(nodes []string) error
	DeleteVotingConfigExclusions() error
	ExcludeNodeAllocation(nodes []string) error
	DeleteNodeAllocationExclusion() error
	GetUsedDataNodes() ([]string, error)
	AssignedShardsSize(node string) (int64, error)
	EnableUpgradeModeML() error
	DisableUpgradeModeML() error
}
