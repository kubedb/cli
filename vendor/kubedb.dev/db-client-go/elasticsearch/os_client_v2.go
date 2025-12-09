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

	"github.com/opensearch-project/opensearch-go/opensearchapi"
	osv2 "github.com/opensearch-project/opensearch-go/v2"
	osv2api "github.com/opensearch-project/opensearch-go/v2/opensearchapi"
	"github.com/pkg/errors"
	core "k8s.io/api/core/v1"
	"k8s.io/klog/v2"
	kutil "kmodules.xyz/client-go"
)

var _ ESClient = &OSClientV2{}

type OSClientV2 struct {
	client *osv2.Client
}

func (os *OSClientV2) ClusterHealthInfo() (map[string]any, error) {
	res, err := os.client.Cluster.Health(
		os.client.Cluster.Health.WithPretty(),
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

func (os *OSClientV2) NodesStats() (map[string]any, error) {
	req := osv2api.NodesStatsRequest{
		Pretty: true,
		Human:  true,
	}

	resp, err := req.Do(context.Background(), os.client)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() // nolint:errcheck

	nodesStats := make(map[string]any)
	if err := json.NewDecoder(resp.Body).Decode(&nodesStats); err != nil {
		return nil, fmt.Errorf("failed to deserialize the response: %v", err)
	}

	return nodesStats, nil
}

func (os *OSClientV2) DisableShardAllocation() error {
	var b strings.Builder
	b.WriteString(DisableShardAllocation)
	req := opensearchapi.ClusterPutSettingsRequest{
		Body:   strings.NewReader(b.String()),
		Pretty: true,
		Human:  true,
	}
	res, err := req.Do(context.Background(), os.client)
	if err != nil {
		return err
	}
	defer res.Body.Close() // nolint:errcheck

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("received status code: %d", res.StatusCode)
	}

	return nil
}

func (es *OSClientV2) ShardStats() ([]ShardInfo, error) {
	req := osv2api.CatShardsRequest{
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

// GetIndicesInfo will return the indices' info of an Elasticsearch database
func (os *OSClientV2) GetIndicesInfo() ([]any, error) {
	req := osv2api.CatIndicesRequest{
		Bytes:  "b", // will return resource size field into byte unit
		Format: "json",
		Pretty: true,
		Human:  true,
	}
	resp, err := req.Do(context.Background(), os.client)
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

func (os *OSClientV2) ClusterStatus() (string, error) {
	res, err := os.client.Cluster.Health(
		os.client.Cluster.Health.WithPretty(),
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

func (os *OSClientV2) GetClusterWriteStatus(ctx context.Context, db *dbapi.Elasticsearch) error {
	// Build the request index & request body
	// send the db specs as body
	indexBody := WriteRequestIndexBody{
		ID: writeRequestID,
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
	res, err3 := osv2api.BulkRequest{
		Index:  writeRequestIndex,
		Body:   strings.NewReader(strings.Join([]string{string(index), string(body)}, "\n") + "\n"),
		Pretty: true,
	}.Do(ctx, os.client.Transport)
	if err3 != nil {
		return errors.Wrap(err3, "Failed to perform write request")
	}
	if res.IsError() {
		return fmt.Errorf("failed to get response from write request with error statuscode %d", res.StatusCode)
	}

	defer func(res *osv2api.Response) {
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

func (os *OSClientV2) GetClusterReadStatus(ctx context.Context, db *dbapi.Elasticsearch) error {
	// Perform a read request in writeRequestIndex/writeRequestID (kubedb-system/info) API
	// Handle error specifically if index has not been created yet
	res, err := osv2api.GetRequest{
		Index:      writeRequestIndex,
		DocumentID: writeRequestID,
	}.Do(ctx, os.client.Transport)
	if err != nil {
		return errors.Wrap(err, "Failed to perform read request")
	}

	defer func(res *osv2api.Response) {
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
		return fmt.Errorf("failed to get response from read request with error statuscode %d", res.StatusCode)
	}

	return nil
}

func (os *OSClientV2) GetTotalDiskUsage(ctx context.Context) (string, error) {
	// Perform a DiskUsageRequest to database to calculate store size of all the elasticsearch indices
	// primary purpose of this function is to provide operator calculated storage of interimVolumeTemplate while taking backup
	// Analyzing field disk usage is resource-intensive. To use the API, RunExpensiveTasks must be set to true. Defaults to false.
	// Get disk usage for all indices using "*" wildcard.
	flag := true
	res, err := osv2api.IndicesDiskUsageRequest{
		Index:             diskUsageRequestIndex,
		Pretty:            true,
		Human:             true,
		RunExpensiveTasks: &flag,
		ExpandWildcards:   diskUsageRequestWildcards,
	}.Do(ctx, os.client.Transport)
	if err != nil {
		return "", errors.Wrap(err, "Failed to perform Disk Usage Request")
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close() // nolint:errcheck
		if err != nil {
			klog.Errorf("failed to close response body from Disk Usage Request, reason: %s", err)
		}
	}(res.Body)

	// Parse the json response to get total storage used for all index
	totalDiskUsage, err := calculateDatabaseSize(res.Body)
	if err != nil {
		return "", errors.Wrap(err, "Failed to parse json response to get disk usage")
	}

	return totalDiskUsage, nil
}

func (os *OSClientV2) SyncCredentialFromSecret(secret *core.Secret) error {
	return nil
}

func (os *OSClientV2) GetDBUserRole(ctx context.Context) (error, bool) {
	return errors.New("not supported in os version 2"), false
}

func (os *OSClientV2) CreateDBUserRole(ctx context.Context) error {
	return errors.New("not supported in os version 2")
}

func (os *OSClientV2) IndexExistsOrNot(index string) error {
	req := opensearchapi.IndicesExistsRequest{
		Index: []string{index},
	}
	res, err := req.Do(context.Background(), os.client)
	if err != nil {
		klog.Errorf("failed to get response while checking either index exists or not %v", err)
		return err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close() // nolint:errcheck
		if err != nil {
			klog.Errorf("failed to close response body for checking the existence of index, reason: %s", err)
		}
	}(res.Body)

	if res.IsError() {
		klog.Errorf("index does not exist")
		return fmt.Errorf("failed to get index with statuscode %d", res.StatusCode)
	}
	return nil
}

func (os *OSClientV2) CreateIndex(index string) error {
	req := opensearchapi.IndicesCreateRequest{
		Index:  index,
		Pretty: true,
		Human:  true,
	}

	res, err := req.Do(context.Background(), os.client)
	if err != nil {
		klog.Errorf("failed to apply create index request, reason: %s", err)
		return err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close() // nolint:errcheck
		if err != nil {
			klog.Errorf("failed to close response body for creating index, reason: %s", err)
		}
	}(res.Body)

	if res.IsError() {
		klog.Errorf("creating index failed with statuscode %d", res.StatusCode)
		return errors.New("failed to create index")
	}

	return nil
}

func (os *OSClientV2) DeleteIndex(index string) error {
	req := opensearchapi.IndicesDeleteRequest{
		Index: []string{index},
	}

	res, err := req.Do(context.Background(), os.client)
	if err != nil {
		klog.Errorf("failed to apply delete index request, reason: %s", err)
		return err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close() // nolint:errcheck
		if err != nil {
			klog.Errorf("failed to close response body for deleting index, reason: %s", err)
		}
	}(res.Body)

	if res.IsError() {
		klog.Errorf("failed to delete index with status code %d", res.StatusCode)
		return errors.New("failed to delete index")
	}

	return nil
}

func (os *OSClientV2) CountData(index string) (int, error) {
	req := opensearchapi.CountRequest{
		Index: []string{index},
	}

	res, err := req.Do(context.Background(), os.client)
	if err != nil {
		return 0, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close() // nolint:errcheck
		if err != nil {
			klog.Errorf("failed to close response body for counting data, reason: %s", err)
		}
	}(res.Body)

	if res.IsError() {
		klog.Errorf("failed to count number of documents in index with statuscode %d", res.StatusCode)
		return 0, errors.New("failed to count number of documents in index")
	}

	var response map[string]any
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return 0, err
	}

	count, ok := response["count"]
	if !ok {
		return 0, errors.New("failed to parse value for index count in response body")
	}

	return int(count.(float64)), nil
}

func (os *OSClientV2) PutData(index, id string, data map[string]any) error {
	var b strings.Builder
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return errors.Wrap(err, "failed to Marshal data")
	}
	b.Write(dataBytes)

	// CreateRequest is not supported in OS V2.
	req := opensearchapi.IndexRequest{
		Index:      index,
		DocumentID: id,
		Body:       strings.NewReader(b.String()),
		Pretty:     true,
		Human:      true,
	}

	res, err := req.Do(context.Background(), os.client)
	if err != nil {
		klog.Errorf("failed to put data in the index, reason: %s", err)
		return err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close() // nolint:errcheck
		if err != nil {
			klog.Errorf("failed to close response body for putting data in the index, reason: %s", err)
		}
	}(res.Body)

	if res.IsError() {
		klog.Errorf("failed to put data in an index with statuscode %d", res.StatusCode)
		return errors.New("failed to put data in an index")
	}
	return nil
}

func (os *OSClientV2) ReEnableShardAllocation() error {
	var b strings.Builder
	b.WriteString(ReEnableShardAllocation)
	req := opensearchapi.ClusterPutSettingsRequest{
		Body:   strings.NewReader(b.String()),
		Pretty: true,
		Human:  true,
	}
	res, err := req.Do(context.Background(), os.client)
	if err != nil {
		return err
	}
	defer res.Body.Close() // nolint:errcheck

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("received status code: %d", res.StatusCode)
	}

	return nil
}

func (os *OSClientV2) CheckVersion() (string, error) {
	req := opensearchapi.InfoRequest{
		Pretty: true,
		Human:  true,
	}

	res, err := req.Do(context.Background(), os.client)
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

func (os *OSClientV2) GetClusterStatus() (string, error) {
	res, err := os.client.Cluster.Health(
		os.client.Cluster.Health.WithPretty(),
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

func (os *OSClientV2) CountIndex() (int, error) {
	req := opensearchapi.IndicesGetSettingsRequest{
		Index:  []string{"_all"},
		Pretty: true,
		Human:  true,
	}

	res, err := req.Do(context.Background(), os.client)
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

func (os *OSClientV2) GetData(_index, _type, _id string) (map[string]any, error) {
	req := opensearchapi.GetRequest{
		Index:        _index,
		DocumentType: _type,
		DocumentID:   _id,
		Pretty:       true,
		Human:        true,
	}

	res, err := req.Do(context.Background(), os.client)
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

func (os *OSClientV2) CountNodes() (int64, error) {
	req := opensearchapi.NodesInfoRequest{
		Pretty: false,
		Human:  false,
	}

	resp, err := req.Do(context.Background(), os.client)
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

func (os *OSClientV2) AddVotingConfigExclusions(nodes []string) error {
	nodeNames := strings.Join(nodes, ",")
	req := opensearchapi.ClusterPostVotingConfigExclusionsRequest{
		NodeNames: nodeNames,
	}

	res, err := req.Do(context.Background(), os.client)
	if err != nil {
		return err
	}
	defer res.Body.Close() // nolint:errcheck

	if res.IsError() {
		return fmt.Errorf("failed with response.StatusCode: %d", res.StatusCode)
	}
	return nil
}

func (os *OSClientV2) DeleteVotingConfigExclusions() error {
	req, err := http.NewRequest(http.MethodDelete, VotingExclusionUrl, nil)
	if err != nil {
		return err
	}

	resp, err := os.client.Perform(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close() // nolint:errcheck

	if resp.StatusCode > 299 {
		return fmt.Errorf("failed with response.StatusCode: %d", resp.StatusCode)
	}
	return nil
}

func (os *OSClientV2) ExcludeNodeAllocation(nodes []string) error {
	list := strings.Join(nodes, ",")
	var body bytes.Buffer
	t, err := template.New("").Parse(ExcludeNodeAllocation)
	if err != nil {
		return errors.Wrap(err, "failed to parse the template")
	}

	if err := t.Execute(&body, list); err != nil {
		return err
	}

	req := opensearchapi.ClusterPutSettingsRequest{
		Body:   bytes.NewReader(body.Bytes()),
		Pretty: true,
		Human:  true,
	}
	resp, err := req.Do(context.Background(), os.client)
	if err != nil {
		return err
	}
	defer resp.Body.Close() // nolint:errcheck

	if resp.IsError() {
		return fmt.Errorf("received status code: %d", resp.StatusCode)
	}

	return nil
}

func (os *OSClientV2) DeleteNodeAllocationExclusion() error {
	var b strings.Builder
	b.WriteString(DeleteNodeAllocationExclusion)
	req := opensearchapi.ClusterPutSettingsRequest{
		Body:   strings.NewReader(b.String()),
		Pretty: true,
		Human:  true,
	}
	resp, err := req.Do(context.Background(), os.client)
	if err != nil {
		return err
	}
	defer resp.Body.Close() // nolint:errcheck

	if resp.IsError() {
		return fmt.Errorf("received status code: %d", resp.StatusCode)
	}

	return nil
}

func (os *OSClientV2) GetUsedDataNodes() ([]string, error) {
	req := opensearchapi.CatShardsRequest{
		Pretty: true,
		Human:  true,
		Format: "json",
		H:      []string{"index,node"},
	}

	resp, err := req.Do(context.Background(), os.client)
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
func (os *OSClientV2) AssignedShardsSize(node string) (int64, error) {
	req := opensearchapi.NodesStatsRequest{
		NodeID: []string{node},
		Pretty: true,
		Human:  true,
	}

	resp, err := req.Do(context.Background(), os.client)
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
func (os *OSClientV2) EnableUpgradeModeML() error {
	return nil
}

// DisableUpgradeModeML	disables upgrade modes for ML nodes.
func (os *OSClientV2) DisableUpgradeModeML() error {
	return nil
}
