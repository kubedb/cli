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

package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"text/template"

	dbapi "kubedb.dev/apimachinery/apis/kubedb/v1"

	esv5 "github.com/elastic/go-elasticsearch/v5"
	"github.com/elastic/go-elasticsearch/v5/esapi"
	"github.com/pkg/errors"
	core "k8s.io/api/core/v1"
)

var _ ESClient = &ESClientV5{}

type ESClientV5 struct {
	client *esv5.Client
}

func (es *ESClientV5) ClusterHealthInfo() (map[string]any, error) {
	return nil, nil
}

func (es *ESClientV5) NodesStats() (map[string]any, error) {
	// Todo: need to implement for version 5
	return nil, nil
}

func (es *ESClientV5) ShardStats() ([]ShardInfo, error) {
	return nil, nil
}

// GetIndicesInfo will return the indices info of an Elasticsearch database
func (es *ESClientV5) GetIndicesInfo() ([]any, error) {
	return nil, nil
}

func (es *ESClientV5) ClusterStatus() (string, error) {
	res, err := es.client.Cluster.Health(
		es.client.Cluster.Health.WithPretty(),
	)
	if err != nil {
		return "", err
	}
	defer res.Body.Close() // nolint:errcheck

	response := make(map[string]any)
	if err2 := json.NewDecoder(res.Body).Decode(&response); err2 != nil {
		return "", errors.Wrap(err2, "failed to parse the response body")
	}
	if value, ok := response["status"]; ok {
		if strValue, ok := value.(string); ok {
			return strValue, nil
		}
		return "", errors.New("failed to convert response to string")
	}
	return "", errors.New("status is missing")
}

// kibana_system, logstash_system etc. internal users
// are not supported for versions 6.x.x and,
// kibana, logstash can be accessed using elastic superuser
// so, sysncing is not required for other builtin users
func (es *ESClientV5) SyncCredentialFromSecret(secret *core.Secret) error {
	return nil
}

func (es *ESClientV5) GetClusterWriteStatus(ctx context.Context, db *dbapi.Elasticsearch) error {
	return nil
}

func (es *ESClientV5) GetClusterReadStatus(ctx context.Context, db *dbapi.Elasticsearch) error {
	return nil
}

func (es *ESClientV5) GetTotalDiskUsage(ctx context.Context) (string, error) {
	return "", nil
}

func (es *ESClientV5) GetDBUserRole(ctx context.Context) (error, bool) {
	return errors.New("not supported in es version 5"), false
}

func (es *ESClientV5) CreateDBUserRole(ctx context.Context) error {
	return errors.New("not supported in es version 5")
}

func (es *ESClientV5) IndexExistsOrNot(index string) error {
	return errors.New("not supported in es version 5")
}

func (es *ESClientV5) CreateIndex(index string) error {
	return errors.New("not supported in es version 5")
}

func (es *ESClientV5) DeleteIndex(index string) error {
	return errors.New("not supported in es version 5")
}

func (es *ESClientV5) CountData(index string) (int, error) {
	return 0, errors.New("not supported in es version 5")
}

func (es *ESClientV5) PutData(index, id string, data map[string]any) error {
	return errors.New("not supported in es version 5")
}

func (es *ESClientV5) DisableShardAllocation() error {
	var b strings.Builder
	b.WriteString(DisableShardAllocation)
	req := esapi.ClusterPutSettingsRequest{
		Body:   strings.NewReader(b.String()),
		Pretty: true,
		Human:  true,
	}
	res, err := req.Do(context.Background(), es.client)
	if err != nil {
		return err
	}
	defer res.Body.Close() // nolint:errcheck

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("received status code: %d", res.StatusCode)
	}

	return nil
}

func (es *ESClientV5) ReEnableShardAllocation() error {
	var b strings.Builder
	b.WriteString(ReEnableShardAllocation)
	req := esapi.ClusterPutSettingsRequest{
		Body:   strings.NewReader(b.String()),
		Pretty: true,
		Human:  true,
	}
	res, err := req.Do(context.Background(), es.client)
	if err != nil {
		return err
	}
	defer res.Body.Close() // nolint:errcheck

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("received status code: %d", res.StatusCode)
	}

	return nil
}

func (es *ESClientV5) CheckVersion() (string, error) {
	req := esapi.InfoRequest{
		Pretty: true,
		Human:  true,
	}

	res, err := req.Do(context.Background(), es.client)
	if err != nil {
		return "", err
	}
	defer res.Body.Close() // nolint:errcheck

	nodeInfo := new(Info)
	if err := json.NewDecoder(res.Body).Decode(&nodeInfo); err != nil {
		return "", errors.Wrap(err, "failed to deserialize the response")
	}

	if nodeInfo.Version.Number == "" {
		return "", errors.New("elasticsearch version is empty")
	}

	return nodeInfo.Version.Number, nil
}

func (es *ESClientV5) GetClusterStatus() (string, error) {
	res, err := es.client.Cluster.Health(
		es.client.Cluster.Health.WithPretty(),
	)
	if err != nil {
		return "", err
	}
	defer res.Body.Close() // nolint:errcheck

	response := make(map[string]any)
	if err2 := json.NewDecoder(res.Body).Decode(&response); err2 != nil {
		return "", errors.Wrap(err2, "failed to parse the response body")
	}
	if value, ok := response["status"]; ok {
		return value.(string), nil
	}
	return "", errors.New("status is missing")
}

func (es *ESClientV5) CountIndex() (int, error) {
	req := esapi.IndicesGetSettingsRequest{
		Index:  []string{"_all"},
		Pretty: true,
		Human:  true,
	}

	res, err := req.Do(context.Background(), es.client)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close() // nolint:errcheck

	if res.IsError() {
		return 0, fmt.Errorf("received status code: %d", res.StatusCode)
	}

	response := make(map[string]any)
	if err2 := json.NewDecoder(res.Body).Decode(&response); err2 != nil {
		return 0, errors.Wrap(err2, "failed to parse the response body")
	}
	return len(response), nil
}

func (es *ESClientV5) GetData(_index, _type, _id string) (map[string]any, error) {
	req := esapi.GetRequest{
		Index:        _index,
		DocumentType: _type,
		DocumentID:   _id,
		Pretty:       true,
		Human:        true,
	}

	res, err := req.Do(context.Background(), es.client)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close() // nolint:errcheck

	if res.IsError() {
		return nil, fmt.Errorf("received status code: %d", res.StatusCode)
	}

	response := make(map[string]any)
	if err2 := json.NewDecoder(res.Body).Decode(&response); err2 != nil {
		return nil, errors.Wrap(err2, "failed to parse the response body")
	}

	return response, nil
}

func (es *ESClientV5) CountNodes() (int64, error) {
	req := esapi.NodesInfoRequest{
		Pretty: false,
		Human:  false,
	}

	resp, err := req.Do(context.Background(), es.client)
	if err != nil {
		return -1, err
	}
	defer resp.Body.Close() // nolint:errcheck

	nodeInfo := new(NodeInfo)
	if err := json.NewDecoder(resp.Body).Decode(&nodeInfo); err != nil {
		return -1, errors.Wrap(err, "failed to deserialize the response")
	}

	if nodeInfo.Nodes.Total == "" {
		return -1, errors.New("Node count is empty")
	}

	return nodeInfo.Nodes.Total.Int64()
}

func (es *ESClientV5) AddVotingConfigExclusions(nodes []string) error {
	url := fmt.Sprintf("/_cluster/voting_config_exclusions/%s?timeout=120s",
		strings.Join(nodes, ","),
	)

	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return err
	}

	resp, err := es.client.Perform(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close() // nolint:errcheck

	if resp.StatusCode > 299 {
		return fmt.Errorf("failed with response.StatusCode: %d", resp.StatusCode)
	}
	return nil
}

func (es *ESClientV5) DeleteVotingConfigExclusions() error {
	req, err := http.NewRequest(http.MethodDelete, VotingExclusionUrl, nil)
	if err != nil {
		return err
	}

	resp, err := es.client.Perform(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close() // nolint:errcheck

	if resp.StatusCode > 299 {
		return fmt.Errorf("failed with response.StatusCode: %d", resp.StatusCode)
	}
	return nil
}

func (es *ESClientV5) ExcludeNodeAllocation(nodes []string) error {
	list := strings.Join(nodes, ",")
	var body bytes.Buffer
	t, err := template.New("").Parse(ExcludeNodeAllocation)
	if err != nil {
		return errors.Wrap(err, "failed to parse the template")
	}

	if err := t.Execute(&body, list); err != nil {
		return err
	}

	req := esapi.ClusterPutSettingsRequest{
		Body:   bytes.NewReader(body.Bytes()),
		Pretty: true,
		Human:  true,
	}
	resp, err := req.Do(context.Background(), es.client)
	if err != nil {
		return err
	}
	defer resp.Body.Close() // nolint:errcheck

	if resp.IsError() {
		return fmt.Errorf("received status code: %d", resp.StatusCode)
	}

	return nil
}

func (es *ESClientV5) DeleteNodeAllocationExclusion() error {
	var b strings.Builder
	b.WriteString(DeleteNodeAllocationExclusion)
	req := esapi.ClusterPutSettingsRequest{
		Body:   strings.NewReader(b.String()),
		Pretty: true,
		Human:  true,
	}
	resp, err := req.Do(context.Background(), es.client)
	if err != nil {
		return err
	}
	defer resp.Body.Close() // nolint:errcheck

	if resp.IsError() {
		return fmt.Errorf("received status code: %d", resp.StatusCode)
	}

	return nil
}

func (es *ESClientV5) GetUsedDataNodes() ([]string, error) {
	req := esapi.CatShardsRequest{
		Pretty: true,
		Human:  true,
		Format: "json",
		H:      []string{"index,node"},
	}

	resp, err := req.Do(context.Background(), es.client)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() // nolint:errcheck

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var list []IndexDistribution
	err = json.Unmarshal(data, &list)
	if err != nil {
		return nil, err
	}

	var nodes []string
	// Skip the ignorable shard,
	// Because every node has a copy of this.
	// We can skip it while scaling down the data node.
	for _, value := range list {
		if !IsIgnorableIndex(value.Index) {
			nodes = append(nodes, value.Node)
		}
	}
	return nodes, nil
}

// AssignedShardsSize returns the assigned shards size of a given node
func (es *ESClientV5) AssignedShardsSize(node string) (int64, error) {
	req := esapi.NodesStatsRequest{
		NodeID: []string{node},
		Pretty: true,
		Human:  true,
	}

	resp, err := req.Do(context.Background(), es.client)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close() // nolint:errcheck

	response := new(NodesStats)
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return 0, err
	}

	for _, value := range response.Nodes {
		return value.Indices.Store.SizeInBytes, nil
	}
	return 0, errors.New("empty response body")
}

// EnableUpgradeModeML	enables upgrade modes for ML nodes.
//
//	Elasticsearch v6.x doesn't have ML nodes. Return nil.
func (es *ESClientV5) EnableUpgradeModeML() error {
	return nil
}

// DisableUpgradeModeML	disables upgrade modes for ML nodes.
//
//	Elasticsearch v6.x doesn't have ML nodes. Return nil.
func (es *ESClientV5) DisableUpgradeModeML() error {
	return nil
}
