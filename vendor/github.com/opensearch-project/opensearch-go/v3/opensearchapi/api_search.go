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
	"strings"

	"github.com/opensearch-project/opensearch-go/v3"
)

// Search executes a /_search request with the optional SearchReq
func (c Client) Search(ctx context.Context, req *SearchReq) (*SearchResp, error) {
	if req == nil {
		req = &SearchReq{}
	}

	var (
		data SearchResp
		err  error
	)
	if data.response, err = c.do(ctx, req, &data); err != nil {
		return &data, err
	}

	return &data, nil
}

// SearchReq represents possible options for the /_search request
type SearchReq struct {
	Indices []string
	Body    io.Reader

	Header http.Header
	Params SearchParams
}

// GetRequest returns the *http.Request that gets executed by the client
func (r SearchReq) GetRequest() (*http.Request, error) {
	var path string
	if len(r.Indices) > 0 {
		path = fmt.Sprintf("/%s/_search", strings.Join(r.Indices, ","))
	} else {
		path = "_search"
	}

	return opensearch.BuildRequest(
		"POST",
		path,
		r.Body,
		r.Params.get(),
		r.Header,
	)
}

// SearchResp represents the returned struct of the /_search response
type SearchResp struct {
	Took    int            `json:"took"`
	Timeout bool           `json:"timed_out"`
	Shards  ResponseShards `json:"_shards"`
	Hits    struct {
		Total struct {
			Value    int    `json:"value"`
			Relation string `json:"relation"`
		} `json:"total"`
		MaxScore float32     `json:"max_score"`
		Hits     []SearchHit `json:"hits"`
	} `json:"hits"`
	Errors       bool            `json:"errors"`
	Aggregations json.RawMessage `json:"aggregations"`
	ScrollID     *string         `json:"_scroll_id,omitempty"`
	response     *opensearch.Response
}

// Inspect returns the Inspect type containing the raw *opensearch.Reponse
func (r SearchResp) Inspect() Inspect {
	return Inspect{Response: r.response}
}

// SearchHit is a sub type of SearchResp containing information of the search hit with an unparsed Source field
type SearchHit struct {
	Index  string          `json:"_index"`
	ID     string          `json:"_id"`
	Score  float32         `json:"_score"`
	Source json.RawMessage `json:"_source"`
	Type   string          `json:"_type"` // Deprecated field
	Sort   []any           `json:"sort"`
}
