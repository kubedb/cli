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
	"strconv"
	"strings"
	"time"
)

func newCatClusterManagerFunc(t Transport) CatClusterManager {
	return func(o ...func(*CatClusterManagerRequest)) (*Response, error) {
		var r = CatClusterManagerRequest{}
		for _, f := range o {
			f(&r)
		}
		return r.Do(r.ctx, t)
	}
}

// ----- API Definition -------------------------------------------------------

// CatClusterManager returns information about the cluster-manager node.
type CatClusterManager func(o ...func(*CatClusterManagerRequest)) (*Response, error)

// CatClusterManagerRequest configures the Cat Cluster Manager API request.
type CatClusterManagerRequest struct {
	Format                string
	H                     []string
	Help                  *bool
	Local                 *bool
	MasterTimeout         time.Duration
	ClusterManagerTimeout time.Duration
	S                     []string
	V                     *bool

	Pretty     bool
	Human      bool
	ErrorTrace bool
	FilterPath []string

	Header http.Header

	ctx context.Context
}

// Do executes the request and returns response or error.
func (r CatClusterManagerRequest) Do(ctx context.Context, transport Transport) (*Response, error) {
	var (
		method string
		path   strings.Builder
		params map[string]string
	)

	method = "GET"

	path.Grow(len("/_cat/cluster_manager"))
	path.WriteString("/_cat/cluster_manager")

	params = make(map[string]string)

	if r.Format != "" {
		params["format"] = r.Format
	}

	if len(r.H) > 0 {
		params["h"] = strings.Join(r.H, ",")
	}

	if r.Help != nil {
		params["help"] = strconv.FormatBool(*r.Help)
	}

	if r.Local != nil {
		params["local"] = strconv.FormatBool(*r.Local)
	}

	if r.MasterTimeout != 0 {
		params["master_timeout"] = formatDuration(r.MasterTimeout)
	}

	if r.ClusterManagerTimeout != 0 {
		params["cluster_manager_timeout"] = formatDuration(r.ClusterManagerTimeout)
	}

	if len(r.S) > 0 {
		params["s"] = strings.Join(r.S, ",")
	}

	if r.V != nil {
		params["v"] = strconv.FormatBool(*r.V)
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
func (f CatClusterManager) WithContext(v context.Context) func(*CatClusterManagerRequest) {
	return func(r *CatClusterManagerRequest) {
		r.ctx = v
	}
}

// WithFormat - a short version of the accept header, e.g. json, yaml.
func (f CatClusterManager) WithFormat(v string) func(*CatClusterManagerRequest) {
	return func(r *CatClusterManagerRequest) {
		r.Format = v
	}
}

// WithH - comma-separated list of column names to display.
func (f CatClusterManager) WithH(v ...string) func(*CatClusterManagerRequest) {
	return func(r *CatClusterManagerRequest) {
		r.H = v
	}
}

// WithHelp - return help information.
func (f CatClusterManager) WithHelp(v bool) func(*CatClusterManagerRequest) {
	return func(r *CatClusterManagerRequest) {
		r.Help = &v
	}
}

// WithLocal - return local information, do not retrieve the state from cluster-manager node (default: false).
func (f CatClusterManager) WithLocal(v bool) func(*CatClusterManagerRequest) {
	return func(r *CatClusterManagerRequest) {
		r.Local = &v
	}
}

// WithMasterTimeout - explicit operation timeout for connection to cluster-manager node.
//
// Deprecated: To promote inclusive language, use WithClusterManagerTimeout instead.
func (f CatClusterManager) WithMasterTimeout(v time.Duration) func(*CatClusterManagerRequest) {
	return func(r *CatClusterManagerRequest) {
		r.MasterTimeout = v
	}
}

// WithClusterManagerTimeout - explicit operation timeout for connection to cluster-manager node.
func (f CatClusterManager) WithClusterManagerTimeout(v time.Duration) func(*CatClusterManagerRequest) {
	return func(r *CatClusterManagerRequest) {
		r.ClusterManagerTimeout = v
	}
}

// WithS - comma-separated list of column names or column aliases to sort by.
func (f CatClusterManager) WithS(v ...string) func(*CatClusterManagerRequest) {
	return func(r *CatClusterManagerRequest) {
		r.S = v
	}
}

// WithV - verbose mode. display column headers.
func (f CatClusterManager) WithV(v bool) func(*CatClusterManagerRequest) {
	return func(r *CatClusterManagerRequest) {
		r.V = &v
	}
}

// WithPretty makes the response body pretty-printed.
func (f CatClusterManager) WithPretty() func(*CatClusterManagerRequest) {
	return func(r *CatClusterManagerRequest) {
		r.Pretty = true
	}
}

// WithHuman makes statistical values human-readable.
func (f CatClusterManager) WithHuman() func(*CatClusterManagerRequest) {
	return func(r *CatClusterManagerRequest) {
		r.Human = true
	}
}

// WithErrorTrace includes the stack trace for errors in the response body.
func (f CatClusterManager) WithErrorTrace() func(*CatClusterManagerRequest) {
	return func(r *CatClusterManagerRequest) {
		r.ErrorTrace = true
	}
}

// WithFilterPath filters the properties of the response body.
func (f CatClusterManager) WithFilterPath(v ...string) func(*CatClusterManagerRequest) {
	return func(r *CatClusterManagerRequest) {
		r.FilterPath = v
	}
}

// WithHeader adds the headers to the HTTP request.
func (f CatClusterManager) WithHeader(h map[string]string) func(*CatClusterManagerRequest) {
	return func(r *CatClusterManagerRequest) {
		if r.Header == nil {
			r.Header = make(http.Header)
		}
		for k, v := range h {
			r.Header.Add(k, v)
		}
	}
}

// WithOpaqueID adds the X-Opaque-Id header to the HTTP request.
func (f CatClusterManager) WithOpaqueID(s string) func(*CatClusterManagerRequest) {
	return func(r *CatClusterManagerRequest) {
		if r.Header == nil {
			r.Header = make(http.Header)
		}
		r.Header.Set("X-Opaque-Id", s)
	}
}
