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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"text/template"

	dbapi "kubedb.dev/apimachinery/apis/kubedb/v1"

	esv6 "github.com/elastic/go-elasticsearch/v6"
	"github.com/elastic/go-elasticsearch/v6/esapi"
	"github.com/pkg/errors"
	core "k8s.io/api/core/v1"
	"k8s.io/klog/v2"
	kutil "kmodules.xyz/client-go"
)

var _ ESClient = &ESClientV6{}

type ESClientV6 struct {
	client *esv6.Client
}

func (es *ESClientV6) ClusterHealthInfo() (map[string]any, error) {
	res, err := es.client.Cluster.Health(
		es.client.Cluster.Health.WithPretty(),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close() // nolint:errcheck

	response := make(map[string]any)
	if err2 := json.NewDecoder(res.Body).Decode(&response); err2 != nil {
		return nil, errors.Wrap(err2, "failed to parse the response body")
	}
	return response, nil
}

func (es *ESClientV6) NodesStats() (map[string]any, error) {
	// todo: need to implement for version 6
	return nil, nil
}

func (es *ESClientV6) ShardStats() ([]ShardInfo, error) {
	req := esapi.CatShardsRequest{
		Bytes:  "b",
		Format: "json",
		Pretty: true,
		Human:  true,
		H:      []string{"index", "shard", "prirep", "state", "unassigned.reason"},
	}

	resp, err := req.Do(context.Background(), es.client)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() // nolint:errcheck

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading body:", err)
	}

	var shardStats []ShardInfo
	err = json.Unmarshal(body, &shardStats)
	if err != nil {
		return nil, err
	}
	return shardStats, nil
}

// GetIndicesInfo will return the indices info of an Elasticsearch database
func (es *ESClientV6) GetIndicesInfo() ([]any, error) {
	req := esapi.CatIndicesRequest{
		Bytes:  "b", // will return resource size field into byte unit
		Format: "json",
		Pretty: true,
		Human:  true,
	}

	resp, err := req.Do(context.Background(), es.client)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() // nolint:errcheck

	indicesInfo := make([]any, 0)
	if err := json.NewDecoder(resp.Body).Decode(&indicesInfo); err != nil {
		return nil, fmt.Errorf("failed to deserialize the response: %v", err)
	}

	return indicesInfo, nil
}

func (es *ESClientV6) ClusterStatus() (string, error) {
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
func (es *ESClientV6) SyncCredentialFromSecret(secret *core.Secret) error {
	return nil
}

func (es *ESClientV6) GetClusterWriteStatus(ctx context.Context, db *dbapi.Elasticsearch) error {
	// Build the request index & request body
	// send the db specs as body
	indexBody := WriteRequestIndexBody{
		ID:   writeRequestID,
		Type: writeRequestType,
	}

	indexReq := WriteRequestIndex{indexBody}
	ReqBody := db.Spec

	// encode the request index & request body
	index, err1 := json.Marshal(indexReq)
	if err1 != nil {
		return errors.Wrap(err1, "Failed to encode index for performing write request")
	}
	body, err2 := json.Marshal(ReqBody)
	if err2 != nil {
		return errors.Wrap(err2, "Failed to encode request body for performing write request")
	}

	// make write request & fetch response
	// check for write request failure & error from response body
	// Bulk API Performs multiple indexing or delete operations in a single API call
	// This reduces overhead and can greatly increase indexing speed it Indexes the specified document
	// If the document exists, replaces the document and increments the version
	res, err3 := esapi.BulkRequest{
		Index:  writeRequestIndex,
		Body:   strings.NewReader(strings.Join([]string{string(index), string(body)}, "\n") + "\n"),
		Pretty: true,
	}.Do(ctx, es.client.Transport)
	if err3 != nil {
		return errors.Wrap(err3, "Failed to perform write request")
	}
	if res.IsError() {
		return fmt.Errorf("failed to get response from write request with error statuscode %d", res.StatusCode)
	}

	defer func(res *esapi.Response) {
		if res != nil {
			err3 = res.Body.Close() // nolint:errcheck
			if err3 != nil {
				klog.Errorf("Failed to close write request response body, reason: %s", err3)
			}
		}
	}(res)

	responseBody := make(map[string]any)
	if err4 := json.NewDecoder(res.Body).Decode(&responseBody); err4 != nil {
		return errors.Wrap(err4, "Failed to decode response from write request")
	}

	// Parse the responseBody to check if write operation failed after request being successful
	// `errors` field(boolean) in the json response becomes true if there's and error caused, otherwise it stays nil
	if value, ok := responseBody["errors"]; ok {
		if strValue, ok := value.(bool); ok {
			if !strValue {
				return nil
			}
			return errors.Errorf("Write request responded with error, %v", responseBody)
		}
		return errors.New("Failed to parse value for `errors` in response from write request")
	}
	return errors.New("Failed to parse key `errors` in response from write request")
}

func (es *ESClientV6) GetClusterReadStatus(ctx context.Context, db *dbapi.Elasticsearch) error {
	// Perform a read request in writeRequestIndex/writeRequestID (kubedb-system/info) API
	// Handle error specifically if index has not been created yet
	res, err := esapi.GetRequest{
		Index:      writeRequestIndex,
		DocumentID: writeRequestID,
	}.Do(ctx, es.client.Transport)
	if err != nil {
		return errors.Wrap(err, "Failed to perform read request")
	}

	defer func(res *esapi.Response) {
		if res != nil {
			err = res.Body.Close() // nolint:errcheck
			if err != nil {
				klog.Errorf("failed to close read request response body, reason: %s", err)
			}
		}
	}(res)

	if res.StatusCode == http.StatusNotFound {
		return kutil.ErrNotFound
	}
	if res.IsError() {
		return fmt.Errorf("failed to get response from write request with error statuscode %d", res.StatusCode)
	}

	return nil
}

func (es *ESClientV6) GetTotalDiskUsage(ctx context.Context) (string, error) {
	return "", nil
}

func (es *ESClientV6) GetDBUserRole(ctx context.Context) (error, bool) {
	return errors.New("not supported in es version 6"), false
}

func (es *ESClientV6) CreateDBUserRole(ctx context.Context) error {
	return errors.New("not supported in es version 6")
}

func (es *ESClientV6) IndexExistsOrNot(index string) error {
	return errors.New("not supported in es version 6")
}

func (es *ESClientV6) CreateIndex(index string) error {
	return errors.New("not supported in es version 6")
}

func (es *ESClientV6) DeleteIndex(index string) error {
	return errors.New("not supported in es version 6")
}

func (es *ESClientV6) CountData(index string) (int, error) {
	return 0, errors.New("not supported in es version 6")
}

func (es *ESClientV6) PutData(index, id string, data map[string]any) error {
	return errors.New("not supported in es version 6")
}

func (es *ESClientV6) DisableShardAllocation() error {
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

func (es *ESClientV6) ReEnableShardAllocation() error {
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

func (es *ESClientV6) CheckVersion() (string, error) {
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

func (es *ESClientV6) GetClusterStatus() (string, error) {
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

func (es *ESClientV6) CountIndex() (int, error) {
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

func (es *ESClientV6) GetData(_index, _type, _id string) (map[string]any, error) {
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

func (es *ESClientV6) CountNodes() (int64, error) {
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

func (es *ESClientV6) AddVotingConfigExclusions(nodes []string) error {
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

func (es *ESClientV6) DeleteVotingConfigExclusions() error {
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

func (es *ESClientV6) ExcludeNodeAllocation(nodes []string) error {
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

func (es *ESClientV6) DeleteNodeAllocationExclusion() error {
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

func (es *ESClientV6) GetUsedDataNodes() ([]string, error) {
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
func (es *ESClientV6) AssignedShardsSize(node string) (int64, error) {
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
func (es *ESClientV6) EnableUpgradeModeML() error {
	return nil
}

// DisableUpgradeModeML	disables upgrade modes for ML nodes.
//
//	Elasticsearch v6.x doesn't have ML nodes. Return nil.
func (es *ESClientV6) DisableUpgradeModeML() error {
	return nil
}
