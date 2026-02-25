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

package raft

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	DefaultDialContextTimeout = 30 * time.Second
	DefaultIdleConnTimeout    = 3 * time.Second
)

type transferLeadershipRequest struct {
	Transferee *int `json:"transferee" protobuf:"varint,1,opt,name=transferee"`
}

func GetRaftHTTPClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			IdleConnTimeout: DefaultIdleConnTimeout,
			DialContext: (&net.Dialer{
				Timeout: DefaultDialContextTimeout,
			}).DialContext,
		},
	}
}

func DoRaftRequest(method, endpoint, user, pass string, body io.Reader, timeout time.Duration) (*http.Response, error) {
	client := GetRaftHTTPClient()
	req, err := http.NewRequest(method, endpoint, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json; charset=UTF-8")
	req.SetBasicAuth(user, pass)

	if timeout <= 0 {
		timeout = DefaultDialContextTimeout
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	req = req.WithContext(ctx)

	return client.Do(req)
}

func TransferLeadership(endpoint string, transferee int, user, pass string, timeout time.Duration) (string, error) {
	transferInfo := transferLeadershipRequest{
		Transferee: &transferee,
	}

	requestByte, err := json.Marshal(transferInfo)
	if err != nil {
		return "", err
	}

	resp, err := DoRaftRequest(http.MethodPost, endpoint, user, pass, bytes.NewReader(requestByte), timeout)
	if err != nil {
		return "", fmt.Errorf("failed to transfer leadership to %d: %w", transferee, err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to transfer leadership to %d: %w", transferee, err)
	}

	return string(bodyText), nil
}

func GetCurrentPrimaryFromLocalhost(coordinatorClientPort int, user, pass string, timeout time.Duration) (int64, error) {
	endpoint := fmt.Sprintf("http://127.0.0.1:%d/current-primary", coordinatorClientPort)

	resp, err := DoRaftRequest(http.MethodGet, endpoint, user, pass, nil, timeout)
	if err != nil {
		return -1, fmt.Errorf("failed to get current primary from localhost: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		return -1, fmt.Errorf("failed to read current primary response: %w", err)
	}

	response := strings.TrimSpace(string(bodyText))
	raftID, err := strconv.ParseInt(response, 10, 64)
	if err != nil {
		return -1, err
	}

	return raftID, nil
}
