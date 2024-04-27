// SPDX-License-Identifier: Apache-2.0
//
// The OpenSearch Contributors require contributions made to
// this file be licensed under the Apache-2.0 license or a
// compatible open source license.
//
// Modifications Copyright OpenSearch Contributors. See
// GitHub history for details.

package opensearchapi

import (
	"context"
	"net/http"
	"strings"
	"time"
)

func newIndicesGetDataStreamStatsFunc(t Transport) IndicesGetDataStreamStats {
	return func(o ...func(*IndicesGetDataStreamStatsRequest)) (*Response, error) {
		var r = IndicesGetDataStreamStatsRequest{}
		for _, f := range o {
			f(&r)
		}
		return r.Do(r.ctx, t)
	}
}

// ----- API Definition -------------------------------------------------------

// IndicesGetDataStreamStats returns a more insights about the data stream.
type IndicesGetDataStreamStats func(o ...func(*IndicesGetDataStreamStatsRequest)) (*Response, error)

// IndicesGetDataStreamStatsRequest configures the Indices Get Data Stream Stats API request.
type IndicesGetDataStreamStatsRequest struct {
	Name string

	ClusterManagerTimeout time.Duration

	Pretty     bool
	Human      bool
	ErrorTrace bool
	FilterPath []string

	Header http.Header

	ctx context.Context
}

// Do execute the request and returns response or error.
func (r IndicesGetDataStreamStatsRequest) Do(ctx context.Context, transport Transport) (*Response, error) {
	var (
		method string
		path   strings.Builder
		params map[string]string
	)

	method = "GET"

	path.Grow(1 + len("_data_stream") + 1 + len(r.Name) + 1 + len("_stats"))
	path.WriteString("/_data_stream/")
	path.WriteString(r.Name)
	path.WriteString("/_stats")

	params = make(map[string]string)

	if r.ClusterManagerTimeout != 0 {
		params["cluster_manager_timeout"] = formatDuration(r.ClusterManagerTimeout)
	}

	if r.Pretty {
		params["pretty"] = "true"
	}

	if r.Human {
		params["human"] = "true"
	}

	if r.ErrorTrace {
		params["error_trace"] = "true"
	}

	if len(r.FilterPath) > 0 {
		params["filter_path"] = strings.Join(r.FilterPath, ",")
	}

	req, err := newRequest(method, path.String(), nil)
	if err != nil {
		return nil, err
	}

	if len(params) > 0 {
		q := req.URL.Query()
		for k, v := range params {
			q.Set(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}

	if len(r.Header) > 0 {
		if len(req.Header) == 0 {
			req.Header = r.Header
		} else {
			for k, vv := range r.Header {
				for _, v := range vv {
					req.Header.Add(k, v)
				}
			}
		}
	}

	if ctx != nil {
		req = req.WithContext(ctx)
	}

	res, err := transport.Perform(req)
	if err != nil {
		return nil, err
	}

	response := Response{
		StatusCode: res.StatusCode,
		Body:       res.Body,
		Header:     res.Header,
	}

	return &response, nil
}

// WithContext sets the request context.
func (f IndicesGetDataStreamStats) WithContext(v context.Context) func(*IndicesGetDataStreamStatsRequest) {
	return func(r *IndicesGetDataStreamStatsRequest) {
		r.ctx = v
	}
}

// WithName - the comma separated names of the index templates.
func (f IndicesGetDataStreamStats) WithName(v string) func(*IndicesGetDataStreamStatsRequest) {
	return func(r *IndicesGetDataStreamStatsRequest) {
		r.Name = v
	}
}

// WithClusterManagerTimeout - explicit operation timeout for connection to cluster-manager node.
func (f IndicesGetDataStreamStats) WithClusterManagerTimeout(v time.Duration) func(*IndicesGetDataStreamStatsRequest) {
	return func(r *IndicesGetDataStreamStatsRequest) {
		r.ClusterManagerTimeout = v
	}
}

// WithPretty makes the response body pretty-printed.
func (f IndicesGetDataStreamStats) WithPretty() func(*IndicesGetDataStreamStatsRequest) {
	return func(r *IndicesGetDataStreamStatsRequest) {
		r.Pretty = true
	}
}

// WithHuman makes statistical values human-readable.
func (f IndicesGetDataStreamStats) WithHuman() func(*IndicesGetDataStreamStatsRequest) {
	return func(r *IndicesGetDataStreamStatsRequest) {
		r.Human = true
	}
}

// WithErrorTrace includes the stack trace for errors in the response body.
func (f IndicesGetDataStreamStats) WithErrorTrace() func(*IndicesGetDataStreamStatsRequest) {
	return func(r *IndicesGetDataStreamStatsRequest) {
		r.ErrorTrace = true
	}
}

// WithFilterPath filters the properties of the response body.
func (f IndicesGetDataStreamStats) WithFilterPath(v ...string) func(*IndicesGetDataStreamStatsRequest) {
	return func(r *IndicesGetDataStreamStatsRequest) {
		r.FilterPath = v
	}
}

// WithHeader adds the headers to the HTTP request.
func (f IndicesGetDataStreamStats) WithHeader(h map[string]string) func(*IndicesGetDataStreamStatsRequest) {
	return func(r *IndicesGetDataStreamStatsRequest) {
		if r.Header == nil {
			r.Header = make(http.Header)
		}
		for k, v := range h {
			r.Header.Add(k, v)
		}
	}
}

// WithOpaqueID adds the X-Opaque-Id header to the HTTP request.
func (f IndicesGetDataStreamStats) WithOpaqueID(s string) func(*IndicesGetDataStreamStatsRequest) {
	return func(r *IndicesGetDataStreamStatsRequest) {
		if r.Header == nil {
			r.Header = make(http.Header)
		}
		r.Header.Set("X-Opaque-Id", s)
	}
}
