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
	"fmt"
	"io"
	"net/http"
	"strings"

	api "kubedb.dev/apimachinery/apis/kubedb/v1alpha2"

	"github.com/opensearch-project/opensearch-go"
	"github.com/opensearch-project/opensearch-go/opensearchapi"
	"github.com/pkg/errors"
	core "k8s.io/api/core/v1"
	"k8s.io/klog/v2"
	kutil "kmodules.xyz/client-go"
)

var _ ESClient = &OSClientV1{}

type OSClientV1 struct {
	client *opensearch.Client
}

func (os *OSClientV1) ClusterHealthInfo() (map[string]interface{}, error) {
	res, err := os.client.Cluster.Health(
		os.client.Cluster.Health.WithPretty(),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	response := make(map[string]interface{})
	if err2 := json.NewDecoder(res.Body).Decode(&response); err2 != nil {
		return nil, errors.Wrap(err2, "failed to parse the response body")
	}
	return response, nil
}

func (os *OSClientV1) NodesStats() (map[string]interface{}, error) {
	req := opensearchapi.NodesStatsRequest{
		Pretty: true,
		Human:  true,
	}

	resp, err := req.Do(context.Background(), os.client)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	nodesStats := make(map[string]interface{})
	if err := json.NewDecoder(resp.Body).Decode(&nodesStats); err != nil {
		return nil, fmt.Errorf("failed to deserialize the response: %v", err)
	}

	return nodesStats, nil
}

// GetIndicesInfo will return the indices' info of an Elasticsearch database
func (os *OSClientV1) GetIndicesInfo() ([]interface{}, error) {
	req := opensearchapi.CatIndicesRequest{
		Bytes:  "b", // will return resource size field into byte unit
		Format: "json",
		Pretty: true,
		Human:  true,
	}
	resp, err := req.Do(context.Background(), os.client)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	indicesInfo := make([]interface{}, 0)
	if err := json.NewDecoder(resp.Body).Decode(&indicesInfo); err != nil {
		return nil, fmt.Errorf("failed to deserialize the response: %v", err)
	}

	return indicesInfo, nil
}

func (os *OSClientV1) ClusterStatus() (string, error) {
	res, err := os.client.Cluster.Health(
		os.client.Cluster.Health.WithPretty(),
	)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	response := make(map[string]interface{})
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

func (os *OSClientV1) GetClusterWriteStatus(ctx context.Context, db *api.Elasticsearch) error {
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
	res, err3 := opensearchapi.BulkRequest{
		Index:  writeRequestIndex,
		Body:   strings.NewReader(strings.Join([]string{string(index), string(body)}, "\n") + "\n"),
		Pretty: true,
	}.Do(ctx, os.client.Transport)
	if err3 != nil {
		return errors.Wrap(err3, "Failed to perform write request")
	}
	if res.IsError() {
		return errors.New(fmt.Sprintf("Failed to get response from write request with error statuscode %d", res.StatusCode))
	}

	defer func(res *opensearchapi.Response) {
		if res != nil {
			err3 = res.Body.Close()
			if err3 != nil {
				klog.Errorf("Failed to close write request response body", err3)
			}
		}
	}(res)

	responseBody := make(map[string]interface{})
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

func (os *OSClientV1) GetClusterReadStatus(ctx context.Context, db *api.Elasticsearch) error {
	// Perform a read request in writeRequestIndex/writeRequestID (kubedb-system/info) API
	// Handle error specifically if index has not been created yet
	res, err := opensearchapi.GetRequest{
		Index:      writeRequestIndex,
		DocumentID: writeRequestID,
	}.Do(ctx, os.client.Transport)
	if err != nil {
		return errors.Wrap(err, "Failed to perform read request")
	}

	defer func(res *opensearchapi.Response) {
		if res != nil {
			err = res.Body.Close()
			if err != nil {
				klog.Errorf("failed to close read request response body", err)
			}
		}
	}(res)

	if res.StatusCode == http.StatusNotFound {
		return kutil.ErrNotFound
	}
	if res.IsError() {
		return errors.New(fmt.Sprintf("Failed to get response from read request with error statuscode %d", res.StatusCode))
	}

	return nil
}

func (os *OSClientV1) GetTotalDiskUsage(ctx context.Context) (string, error) {
	// Perform a DiskUsageRequest to database to calculate store size of all the elasticsearch indices
	// primary purpose of this function is to provide operator calculated storage of interimVolumeTemplate while taking backup
	// Analyzing field disk usage is resource-intensive. To use the API, RunExpensiveTasks must be set to true. Defaults to false.
	// Get disk usage for all indices using "*" wildcard.
	flag := true
	res, err := opensearchapi.IndicesDiskUsageRequest{
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
		err = Body.Close()
		if err != nil {
			klog.Errorf("failed to close response body from Disk Usage Request", err)
		}
	}(res.Body)

	// Parse the json response to get total storage used for all index
	totalDiskUsage, err := calculateDatabaseSize(res.Body)
	if err != nil {
		return "", errors.Wrap(err, "Failed to parse json response to get disk usage")
	}

	return totalDiskUsage, nil
}

func (os *OSClientV1) SyncCredentialFromSecret(secret *core.Secret) error {
	return nil
}

func (os *OSClientV1) GetDBUserRole(ctx context.Context) (error, bool) {
	return errors.New("not supported in os version 1"), false
}

func (os *OSClientV1) CreateDBUserRole(ctx context.Context) error {
	return errors.New("not supported in os version 1")
}

func (os *OSClientV1) IndexExistsOrNot(index string) error {
	req := opensearchapi.IndicesExistsRequest{
		Index: []string{index},
	}
	res, err := req.Do(context.Background(), os.client)
	if err != nil {
		klog.Errorf(fmt.Sprintf("failed to get response while checking either index exists or not %v", err))
		return err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			klog.Errorf("failed to close response body for checking the existence of index", err)
		}
	}(res.Body)

	if res.IsError() {
		klog.Errorf(fmt.Sprintf("failed to get index with statuscode %d", res.StatusCode))
		return errors.New("index does not exist")
	}
	return nil
}

func (os *OSClientV1) CreateIndex(index string) error {
	req := opensearchapi.IndicesCreateRequest{
		Index:  index,
		Pretty: true,
		Human:  true,
	}

	res, err := req.Do(context.Background(), os.client)
	if err != nil {
		klog.Errorf("failed to apply create index request ", err)
		return err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			klog.Errorf("failed to close response body for creating index", err)
		}
	}(res.Body)

	if res.IsError() {
		klog.Errorf(fmt.Sprintf("Creating index failed with statuscode %d", res.StatusCode))
		return errors.New("failed to create index")
	}

	return nil
}

func (os *OSClientV1) DeleteIndex(index string) error {
	req := opensearchapi.IndicesDeleteRequest{
		Index: []string{index},
	}

	res, err := req.Do(context.Background(), os.client)
	if err != nil {
		klog.Errorf("failed to apply delete index request", err)
		return err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			klog.Errorf("failed to close response body for deleting index", err)
		}
	}(res.Body)

	if res.IsError() {
		klog.Errorf(fmt.Sprintf("failed to delete index with status code %d", res.StatusCode))
		return errors.New("failed to delete index")
	}

	return nil
}

func (os *OSClientV1) CountData(index string) (int, error) {
	req := opensearchapi.CountRequest{
		Index: []string{index},
	}

	res, err := req.Do(context.Background(), os.client)
	if err != nil {
		return 0, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			klog.Errorf("failed to close response body for counting data", err)
		}
	}(res.Body)

	if res.IsError() {
		klog.Errorf(fmt.Sprintf("failed to count number of documents in index with statuscode %d", res.StatusCode))
		return 0, errors.New("failed to count number of documents in index")
	}

	var response map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return 0, err
	}

	count, ok := response["count"]
	if !ok {
		return 0, errors.New("failed to parse value for index count in response body")
	}

	return int(count.(float64)), nil
}

func (os *OSClientV1) PutData(index, id string, data map[string]interface{}) error {
	var b strings.Builder
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return errors.Wrap(err, "failed to Marshal data")
	}
	b.Write(dataBytes)

	req := opensearchapi.CreateRequest{
		Index:      index,
		DocumentID: id,
		Body:       strings.NewReader(b.String()),
		Pretty:     true,
		Human:      true,
	}

	res, err := req.Do(context.Background(), os.client)
	if err != nil {
		klog.Errorf("failed to put data in the index", err)
		return err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			klog.Errorf("failed to close response body for putting data in the index", err)
		}
	}(res.Body)

	if res.IsError() {
		klog.Errorf(fmt.Sprintf("failed to put data in an index with statuscode %d", res.StatusCode))
		return errors.New("failed to put data in an index")
	}
	return nil
}
