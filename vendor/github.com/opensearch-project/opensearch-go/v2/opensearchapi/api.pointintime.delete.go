// SPDX-License-Identifier: Apache-2.0
//
// The OpenSearch Contributors require contributions made to
// this file be licensed under the Apache-2.0 license or a
// compatible open source license.
//
// Modifications Copyright OpenSearch Contributors. See
// GitHub history for details.

// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package opensearchapi

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

func newPointInTimeDeleteFunc(t Transport) PointInTimeDelete {
	return func(o ...func(*PointInTimeDeleteRequest)) (*Response, *PointInTimeDeleteResp, error) {
		var r = PointInTimeDeleteRequest{}
		for _, f := range o {
			f(&r)
		}
		return r.Do(r.ctx, t)
	}
}

// ----- API Definition -------------------------------------------------------

// PointInTimeDelete lets you delete pits used for searching with pagination
type PointInTimeDelete func(o ...func(*PointInTimeDeleteRequest)) (*Response, *PointInTimeDeleteResp, error)

// PointInTimeDeleteRequest configures the Point In Time Delete API request.
type PointInTimeDeleteRequest struct {
	PitID []string

	Pretty     bool
	Human      bool
	ErrorTrace bool
	FilterPath []string

	Header http.Header

	ctx context.Context
}

// PointInTimeDeleteRequestBody is used to from the delete request body
type PointInTimeDeleteRequestBody struct {
	PitID []string `json:"pit_id"`
}

// PointInTimeDeleteResp is a custom type to parse the Point In Time Delete Reponse
type PointInTimeDeleteResp struct {
	Pits []struct {
		PitID      string `json:"pit_id"`
		Successful bool   `json:"successful"`
	} `json:"pits"`
}

// Do executes the request and returns response or error.
func (r PointInTimeDeleteRequest) Do(ctx context.Context, transport Transport) (*Response, *PointInTimeDeleteResp, error) {
	var (
		path   strings.Builder
		params map[string]string
		body   io.Reader

		data PointInTimeDeleteResp
	)
	method := "DELETE"

	path.Grow(len("/_search/point_in_time"))
	path.WriteString("/_search/point_in_time")

	params = make(map[string]string)

	if len(r.PitID) > 0 {
		bodyStruct := PointInTimeDeleteRequestBody{PitID: r.PitID}
		bodyJSON, err := json.Marshal(bodyStruct)
		if err != nil {
			return nil, nil, err
		}
		body = bytes.NewBuffer(bodyJSON)
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

	req, err := newRequest(method, path.String(), body)
	if err != nil {
		return nil, nil, err
	}

	if len(params) > 0 {
		q := req.URL.Query()
		for k, v := range params {
			q.Set(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}

	if body != nil {
		req.Header[headerContentType] = headerContentTypeJSON
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
		return nil, nil, err
	}

	response := Response{
		StatusCode: res.StatusCode,
		Body:       res.Body,
		Header:     res.Header,
	}

	if len(r.FilterPath) != 0 {
		return &response, nil, nil
	}

	if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
		return &response, nil, err
	}
	return &response, &data, nil
}

// WithPitID sets the Pit to delete.
func (f PointInTimeDelete) WithPitID(v ...string) func(*PointInTimeDeleteRequest) {
	return func(r *PointInTimeDeleteRequest) {
		r.PitID = v
	}
}

// WithContext sets the request context.
func (f PointInTimeDelete) WithContext(v context.Context) func(*PointInTimeDeleteRequest) {
	return func(r *PointInTimeDeleteRequest) {
		r.ctx = v
	}
}

// WithPretty makes the response body pretty-printed.
func (f PointInTimeDelete) WithPretty() func(*PointInTimeDeleteRequest) {
	return func(r *PointInTimeDeleteRequest) {
		r.Pretty = true
	}
}

// WithHuman makes statistical values human-readable.
func (f PointInTimeDelete) WithHuman() func(*PointInTimeDeleteRequest) {
	return func(r *PointInTimeDeleteRequest) {
		r.Human = true
	}
}

// WithErrorTrace includes the stack trace for errors in the response body.
func (f PointInTimeDelete) WithErrorTrace() func(*PointInTimeDeleteRequest) {
	return func(r *PointInTimeDeleteRequest) {
		r.ErrorTrace = true
	}
}

// WithFilterPath filters the properties of the response body.
func (f PointInTimeDelete) WithFilterPath(v ...string) func(*PointInTimeDeleteRequest) {
	return func(r *PointInTimeDeleteRequest) {
		r.FilterPath = v
	}
}

// WithHeader adds the headers to the HTTP request.
func (f PointInTimeDelete) WithHeader(h map[string]string) func(*PointInTimeDeleteRequest) {
	return func(r *PointInTimeDeleteRequest) {
		if r.Header == nil {
			r.Header = make(http.Header)
		}
		for k, v := range h {
			r.Header.Add(k, v)
		}
	}
}

// WithOpaqueID adds the X-Opaque-Id header to the HTTP request.
func (f PointInTimeDelete) WithOpaqueID(s string) func(*PointInTimeDeleteRequest) {
	return func(r *PointInTimeDeleteRequest) {
		if r.Header == nil {
			r.Header = make(http.Header)
		}
		r.Header.Set("X-Opaque-Id", s)
	}
}
