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

	api "kubedb.dev/apimachinery/apis/kubedb/v1alpha2"

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
	ClusterHealthInfo() (map[string]interface{}, error)
	ClusterStatus() (string, error)
	CountData(index string) (int, error)
	CreateDBUserRole(ctx context.Context) error
	CreateIndex(index string) error
	DeleteIndex(index string) error
	GetIndicesInfo() ([]interface{}, error)
	GetClusterWriteStatus(ctx context.Context, db *api.Elasticsearch) error
	GetClusterReadStatus(ctx context.Context, db *api.Elasticsearch) error
	GetTotalDiskUsage(ctx context.Context) (string, error)
	GetDBUserRole(ctx context.Context) (error, bool)
	IndexExistsOrNot(index string) error
	NodesStats() (map[string]interface{}, error)
	PutData(index, id string, data map[string]interface{}) error
	SyncCredentialFromSecret(secret *core.Secret) error
}
