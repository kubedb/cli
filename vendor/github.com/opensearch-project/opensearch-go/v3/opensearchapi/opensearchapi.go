// SPDX-License-Identifier: Apache-2.0
//
// The OpenSearch Contributors require contributions made to
// this file be licensed under the Apache-2.0 license or a
// compatible open source license.

package opensearchapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/opensearch-project/opensearch-go/v3"
)

// Config represents the client configuration
type Config struct {
	Client opensearch.Config
}

// Client represents the opensearchapi Client summarizing all API calls
type Client struct {
	Client            *opensearch.Client
	Cat               catClient
	Cluster           clusterClient
	Dangling          danglingClient
	Document          documentClient
	Indices           indicesClient
	Nodes             nodesClient
	Script            scriptClient
	ComponentTemplate componentTemplateClient
	IndexTemplate     indexTemplateClient
	// Deprecated: uses legacy API (/_template), correct API is /_index_template, use IndexTemplate instead
	Template    templateClient
	DataStream  dataStreamClient
	PointInTime pointInTimeClient
	Ingest      ingestClient
	Tasks       tasksClient
	Scroll      scrollClient
	Snapshot    snapshotClient
}

// clientInit inits the Client with all sub clients
func clientInit(rootClient *opensearch.Client) *Client {
	client := &Client{
		Client: rootClient,
	}
	client.Cat = catClient{apiClient: client}
	client.Indices = indicesClient{
		apiClient: client,
		Alias:     aliasClient{apiClient: client},
		Mapping:   mappingClient{apiClient: client},
		Settings:  settingsClient{apiClient: client},
	}
	client.Nodes = nodesClient{apiClient: client}
	client.Cluster = clusterClient{apiClient: client}
	client.Dangling = danglingClient{apiClient: client}
	client.Script = scriptClient{apiClient: client}
	client.Document = documentClient{apiClient: client}
	client.ComponentTemplate = componentTemplateClient{apiClient: client}
	client.IndexTemplate = indexTemplateClient{apiClient: client}
	client.Template = templateClient{apiClient: client}
	client.DataStream = dataStreamClient{apiClient: client}
	client.PointInTime = pointInTimeClient{apiClient: client}
	client.Ingest = ingestClient{apiClient: client}
	client.Tasks = tasksClient{apiClient: client}
	client.Scroll = scrollClient{apiClient: client}
	client.Snapshot = snapshotClient{
		apiClient:  client,
		Repository: repositoryClient{apiClient: client},
	}

	return client
}

// NewClient returns a opensearchapi client
func NewClient(config Config) (*Client, error) {
	rootClient, err := opensearch.NewClient(config.Client)
	if err != nil {
		return nil, err
	}

	return clientInit(rootClient), nil
}

// NewDefaultClient returns a opensearchapi client using defaults
func NewDefaultClient() (*Client, error) {
	rootClient, err := opensearch.NewClient(opensearch.Config{
		Addresses: []string{"http://localhost:9200"},
	})
	if err != nil {
		return nil, err
	}

	return clientInit(rootClient), nil
}

// do calls the opensearch.Client.Do() and checks the response for openseach api errors
func (c *Client) do(ctx context.Context, req opensearch.Request, dataPointer any) (*opensearch.Response, error) {
	resp, err := c.Client.Do(ctx, req, dataPointer)
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		if dataPointer != nil {
			return resp, parseError(resp)
		} else {
			return resp, fmt.Errorf("status: %s", resp.Status())
		}
	}

	return resp, nil
}

// praseError tries to parse the opensearch api error into an custom error
func parseError(resp *opensearch.Response) error {
	if resp.Body == nil {
		return fmt.Errorf("%w, status: %s", ErrUnexpectedEmptyBody, resp.Status())
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrReadBody, err)
	}

	if resp.StatusCode == http.StatusMethodNotAllowed {
		var apiError StringError
		if err = json.Unmarshal(body, &apiError); err != nil {
			return fmt.Errorf("%w: %w", ErrJSONUnmarshalBody, err)
		}
		return apiError
	}

	// ToDo: Parse 404 errors separate as they are not in one standard format
	// https://github.com/opensearch-project/OpenSearch/issues/9988

	var apiError Error
	if err = json.Unmarshal(body, &apiError); err != nil {
		return fmt.Errorf("%w: %w", ErrJSONUnmarshalBody, err)
	}

	return apiError
}

// formatDuration converts duration to a string in the format
// accepted by Opensearch.
func formatDuration(d time.Duration) string {
	if d < time.Millisecond {
		return strconv.FormatInt(int64(d), 10) + "nanos"
	}

	return strconv.FormatInt(int64(d)/int64(time.Millisecond), 10) + "ms"
}

// ToPointer converts any value to a pointer, mainly used for request parameters
func ToPointer[V any](value V) *V {
	return &value
}

// ResponseShards is a sub type of api repsonses containing information about shards
type ResponseShards struct {
	Total      int                     `json:"total"`
	Successful int                     `json:"successful"`
	Failed     int                     `json:"failed"`
	Failures   []ResponseShardsFailure `json:"failures"`
	Skipped    int                     `json:"skipped"`
}

// ResponseShardsFailure is a sub type of ReponseShards containing information about a failed shard
type ResponseShardsFailure struct {
	Shard  int `json:"shard"`
	Index  any `json:"index"`
	Reason struct {
		Type   string `json:"type"`
		Reason string `json:"reason"`
	} `json:"reason"`
}
