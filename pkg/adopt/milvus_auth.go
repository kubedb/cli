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

package adopt

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"kmodules.xyz/client-go/tools/portforward"
)

// Milvus REST error codes for a failed login.
var milvusAuthErrorCodes = map[int]bool{1800: true, 2: true}

type credentialResult int

const (
	credentialsAccepted credentialResult = iota
	credentialsRejected
	credentialsUnverified
)

// checkCredentials logs in to a running proxy of the release with the given
// root credentials, through a port-forward to its RESTful API. The API server's
// service proxy cannot be used: it does not pass the Authorization header on.
func (c *cluster) checkCredentials(ctx context.Context, proxy Unit, restPort int, user, password string) (credentialResult, string) {
	pods, err := c.kube.CoreV1().Pods(c.ns).List(ctx, metav1.ListOptions{LabelSelector: labels.SelectorFromSet(proxy.Selector).String()})
	if err != nil {
		return credentialsUnverified, err.Error()
	}
	pod := ""
	for _, p := range pods.Items {
		if p.DeletionTimestamp == nil && podReady(p) {
			pod = p.Name
			break
		}
	}
	if pod == "" {
		return credentialsUnverified, "no Ready proxy pod to log in to"
	}
	tunnel := portforward.NewTunnel(portforward.TunnelOptions{
		Client:    c.kube.CoreV1().RESTClient(),
		Config:    c.cfg,
		Resource:  "pods",
		Name:      pod,
		Namespace: c.ns,
		Remote:    restPort,
	})
	if err := tunnel.ForwardPort(); err != nil {
		return credentialsUnverified, fmt.Sprintf("port-forward to %s: %v", pod, err)
	}
	defer tunnel.Close()

	scheme, httpc := "http", http.DefaultClient
	if c.restTLS != nil {
		scheme = "https"
		httpc = &http.Client{Transport: &http.Transport{TLSClientConfig: c.restTLS}}
	}
	url := scheme + "://127.0.0.1:" + strconv.Itoa(tunnel.Local) + "/v2/vectordb/collections/list"
	reqCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(reqCtx, http.MethodPost, url, bytes.NewBufferString("{}"))
	if err != nil {
		return credentialsUnverified, err.Error()
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+user+":"+password)
	resp, err := httpc.Do(req)
	if err != nil {
		return credentialsUnverified, fmt.Sprintf("RESTful API of %s on port %d: %v", pod, restPort, err)
	}
	defer resp.Body.Close() // nolint:errcheck
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 64<<10))
	var r struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(body, &r); err != nil {
		return credentialsUnverified, fmt.Sprintf("RESTful API of %s answered HTTP %d: %.200s", pod, resp.StatusCode, body)
	}
	switch {
	case r.Code == 0:
		return credentialsAccepted, pod
	case milvusAuthErrorCodes[r.Code]:
		return credentialsRejected, r.Message
	}
	return credentialsUnverified, fmt.Sprintf("RESTful API of %s answered code %d: %s", pod, r.Code, r.Message)
}

// restPort is where the running proxies serve the RESTful API: the gRPC port
// unless proxy.http.port moves it.
func restPort(cfg map[string]any) int {
	if p, err := strconv.Atoi(lookupString(cfg, "proxy.http.port")); err == nil && p > 0 {
		return p
	}
	if p, err := strconv.Atoi(lookupString(cfg, "proxy.port")); err == nil && p > 0 {
		return p
	}
	return 19530
}
