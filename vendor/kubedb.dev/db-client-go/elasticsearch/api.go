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
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	catalog "kubedb.dev/apimachinery/apis/catalog/v1alpha1"
	dbapi "kubedb.dev/apimachinery/apis/kubedb/v1"
	cs "kubedb.dev/apimachinery/client/clientset/versioned"

	"github.com/Masterminds/semver/v3"
	esv6 "github.com/elastic/go-elasticsearch/v6"
	esv7 "github.com/elastic/go-elasticsearch/v7"
	esv8 "github.com/elastic/go-elasticsearch/v8"
	esv9 "github.com/elastic/go-elasticsearch/v9"
	"github.com/go-logr/logr"
	opensearchv1 "github.com/opensearch-project/opensearch-go"
	opensearchv2 "github.com/opensearch-project/opensearch-go/v2"
	"github.com/pkg/errors"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
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

var log logr.Logger

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

// GetElasticClient Returns Elasticsearch Client, Health Check Status Code, Error
func GetElasticClient(kc kubernetes.Interface, dc cs.Interface, db *dbapi.Elasticsearch, url string) (ESClient, int, error) {
	var username, password string
	if !db.Spec.DisableSecurity && db.Spec.AuthSecret != nil {
		secret, err := kc.CoreV1().Secrets(db.Namespace).Get(context.TODO(), db.Spec.AuthSecret.Name, metav1.GetOptions{})
		if err != nil {
			log.Error(err, fmt.Sprintf("Failed to get secret: %s for Elasticsearch: %s/%s", db.Spec.AuthSecret.Name, db.Namespace, db.Name))
			return nil, -1, errors.Wrap(err, "failed to get the secret")
		}

		if value, ok := secret.Data[core.BasicAuthUsernameKey]; ok {
			username = string(value)
		} else {
			log.Error(errors.New("username is missing"), fmt.Sprintf("Failed for secret: %s/%s, username is missing", secret.Namespace, secret.Name))
			return nil, -1, errors.New("username is missing")
		}

		if value, ok := secret.Data[core.BasicAuthPasswordKey]; ok {
			password = string(value)
		} else {
			log.Error(errors.New("password is missing"), fmt.Sprintf("Failed for secret: %s/%s, password is missing", secret.Namespace, secret.Name))
			return nil, -1, errors.New("password is missing")
		}
	}

	// Get original Elasticsearch version, since the client is version specific
	esVersion, err := dc.CatalogV1alpha1().ElasticsearchVersions().Get(context.TODO(), db.Spec.Version, metav1.GetOptions{})
	if err != nil {
		return nil, -1, errors.Wrap(err, "failed to get elasticsearchVersion")
	}

	version, err := semver.NewVersion(esVersion.Spec.Version)
	if err != nil {
		return nil, -1, errors.Wrap(err, "failed to parse version")
	}

	switch esVersion.Spec.AuthPlugin {
	case catalog.ElasticsearchAuthPluginXpack, catalog.ElasticsearchAuthPluginSearchGuard, catalog.ElasticsearchAuthPluginOpenDistro:
		switch {
		case version.Major() == 6:
			client, err := esv6.NewClient(esv6.Config{
				Addresses:         []string{url},
				Username:          username,
				Password:          password,
				EnableDebugLogger: true,
				DisableRetry:      true,
				Transport: &http.Transport{
					IdleConnTimeout: 3 * time.Second,
					DialContext: (&net.Dialer{
						Timeout: 30 * time.Second,
					}).DialContext,
					TLSClientConfig: &tls.Config{
						InsecureSkipVerify: true,
						MaxVersion:         tls.VersionTLS12,
					},
				},
			})
			if err != nil {
				log.Error(err, fmt.Sprintf("Failed to create HTTP client for Elasticsearch: %s/%s", db.Namespace, db.Name))
				return nil, -1, err
			}
			// do a manual health check to test client
			res, err := client.Cluster.Health(
				client.Cluster.Health.WithPretty(),
			)
			if err != nil {
				return nil, -1, errors.Wrap(err, "failed to perform health check")
			}
			defer res.Body.Close() // nolint:errcheck

			if res.IsError() {
				return &ESClientV6{client: client}, res.StatusCode, fmt.Errorf("health check failed with status code: %d", res.StatusCode)
			}
			return &ESClientV6{client: client}, res.StatusCode, nil

		case version.Major() == 7:
			client, err := esv7.NewClient(esv7.Config{
				Addresses:         []string{url},
				Username:          username,
				Password:          password,
				EnableDebugLogger: true,
				DisableRetry:      true,
				Transport: &http.Transport{
					IdleConnTimeout: 3 * time.Second,
					DialContext: (&net.Dialer{
						Timeout: 30 * time.Second,
					}).DialContext,
					TLSClientConfig: &tls.Config{
						InsecureSkipVerify: true,
						MaxVersion:         tls.VersionTLS12,
					},
				},
			})
			if err != nil {
				log.Error(err, fmt.Sprintf("Failed to create HTTP client for Elasticsearch: %s/%s", db.Namespace, db.Name))
				return nil, -1, err
			}
			// do a manual health check to test client
			res, err := client.Cluster.Health(
				client.Cluster.Health.WithPretty(),
			)
			if err != nil {
				return nil, -1, errors.Wrap(err, "failed to perform health check")
			}
			defer res.Body.Close() // nolint:errcheck

			if res.IsError() {
				return &ESClientV7{client: client}, res.StatusCode, fmt.Errorf("health check failed with status code: %d", res.StatusCode)
			}
			return &ESClientV7{client: client}, res.StatusCode, nil

		case version.Major() == 8:
			client, err := esv8.NewClient(esv8.Config{
				Addresses:         []string{url},
				Username:          username,
				Password:          password,
				EnableDebugLogger: true,
				DisableRetry:      true,
				Transport: &http.Transport{
					IdleConnTimeout: 3 * time.Second,
					DialContext: (&net.Dialer{
						Timeout: 30 * time.Second,
					}).DialContext,
					TLSClientConfig: &tls.Config{
						InsecureSkipVerify: true,
						MaxVersion:         tls.VersionTLS12,
					},
				},
			})
			if err != nil {
				log.Error(err, fmt.Sprintf("Failed to create HTTP client for Elasticsearch: %s/%s", db.Namespace, db.Name))
				return nil, -1, err
			}
			// do a manual health check to test client
			res, err := client.Cluster.Health(
				client.Cluster.Health.WithPretty(),
			)
			if err != nil {
				return nil, -1, errors.Wrap(err, "failed to perform health check")
			}
			defer res.Body.Close() // nolint:errcheck

			if res.IsError() {
				return &ESClientV8{client: client}, res.StatusCode, fmt.Errorf("health check failed with status code: %d", res.StatusCode)
			}
			return &ESClientV8{client: client}, res.StatusCode, nil

		case version.Major() == 9:
			client, err := esv9.NewClient(esv9.Config{
				Addresses:         []string{url},
				Username:          username,
				Password:          password,
				EnableDebugLogger: true,
				DisableRetry:      true,
				Transport: &http.Transport{
					IdleConnTimeout: 3 * time.Second,
					DialContext: (&net.Dialer{
						Timeout: 30 * time.Second,
					}).DialContext,
					TLSClientConfig: &tls.Config{
						InsecureSkipVerify: true,
						MaxVersion:         tls.VersionTLS12,
					},
				},
			})
			if err != nil {
				log.Error(err, fmt.Sprintf("Failed to create HTTP client for Elasticsearch: %s/%s", db.Namespace, db.Name))
				return nil, -1, err
			}
			// do a manual health check to test client
			res, err := client.Cluster.Health(
				client.Cluster.Health.WithPretty(),
			)
			if err != nil {
				return nil, -1, errors.Wrap(err, "failed to perform health check")
			}
			defer res.Body.Close() // nolint:errcheck

			if res.IsError() {
				return &ESClientV9{client: client}, res.StatusCode, fmt.Errorf("health check failed with status code: %d", res.StatusCode)
			}
			return &ESClientV9{client: client}, res.StatusCode, nil
		}

	case catalog.ElasticsearchAuthPluginOpenSearch:
		switch {
		case version.Major() == 1:
			client, err := opensearchv1.NewClient(opensearchv1.Config{
				Addresses:         []string{url},
				Username:          username,
				Password:          password,
				EnableDebugLogger: true,
				DisableRetry:      true,
				Transport: &http.Transport{
					IdleConnTimeout: 3 * time.Second,
					DialContext: (&net.Dialer{
						Timeout: 30 * time.Second,
					}).DialContext,
					TLSClientConfig: &tls.Config{
						InsecureSkipVerify: true,
						MaxVersion:         tls.VersionTLS12,
					},
				},
			})
			if err != nil {
				log.Error(err, fmt.Sprintf("Failed to create HTTP client for Elasticsearch: %s/%s", db.Namespace, db.Name))
				return nil, -1, err
			}
			// do a manual health check to test client
			res, err := client.Cluster.Health(
				client.Cluster.Health.WithPretty(),
			)
			if err != nil {
				return nil, -1, errors.Wrap(err, "failed to perform health check")
			}
			defer res.Body.Close() // nolint:errcheck

			if res.IsError() {
				return &OSClientV1{client: client}, res.StatusCode, fmt.Errorf("health check failed with status code: %d", res.StatusCode)
			}
			return &OSClientV1{client: client}, res.StatusCode, nil
		case version.Major() == 2:
			client, err := opensearchv2.NewClient(opensearchv2.Config{
				Addresses:         []string{url},
				Username:          username,
				Password:          password,
				EnableDebugLogger: true,
				DisableRetry:      true,
				Transport: &http.Transport{
					IdleConnTimeout: 3 * time.Second,
					DialContext: (&net.Dialer{
						Timeout: 30 * time.Second,
					}).DialContext,
					TLSClientConfig: &tls.Config{
						InsecureSkipVerify: true,
						MaxVersion:         tls.VersionTLS12,
					},
				},
			})
			if err != nil {
				log.Error(err, fmt.Sprintf("Failed to create HTTP client for Elasticsearch: %s/%s", db.Namespace, db.Name))
				return nil, -1, err
			}
			// do a manual health check to test client
			res, err := client.Cluster.Health(
				client.Cluster.Health.WithPretty(),
			)
			if err != nil {
				return nil, -1, errors.Wrap(err, "Failed to perform health check")
			}
			defer res.Body.Close() // nolint:errcheck

			if res.IsError() {
				return &OSClientV2{client: client}, res.StatusCode, fmt.Errorf("health check failed with status code: %d", res.StatusCode)
			}
			return &OSClientV2{client: client}, res.StatusCode, nil
		}
	}

	return nil, -1, fmt.Errorf("unknown database version: %s", esVersion.Spec.Version)
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
