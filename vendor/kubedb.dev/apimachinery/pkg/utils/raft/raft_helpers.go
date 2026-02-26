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
	"net/http"
	"strconv"
	"strings"
	"time"

	"kubedb.dev/apimachinery/apis/kubedb"

	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	meta_util "kmodules.xyz/client-go/meta"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const leaderAPIRequestTimeout = DefaultDialContextTimeout

// DBRaftAddressProvider defines the minimum DB metadata required by raft helpers.
type DBRaftAddressProvider interface {
	GoverningServiceDNS(podName string) string
	OffshootName() string
}

type raftNodeInfo struct {
	NodeID *int    `json:"id" protobuf:"varint,1,opt,name=id"`
	URL    *string `json:"url,omitempty" protobuf:"bytes,2,opt,name=url"`
}

// GetCurrentLeaderID queries raft leader id from a coordinator endpoint.
func GetCurrentLeaderID(coordinatorClientPort int, governingDNS string, user, pass string) (uint64, error) {
	url := "http://" + governingDNS + ":" + strconv.Itoa(coordinatorClientPort) + "/current-primary"

	defaultLead := uint64(0)
	resp, err := DoRaftRequest(http.MethodGet, url, user, pass, nil, leaderAPIRequestTimeout)
	if err != nil {
		return defaultLead, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		return defaultLead, err
	}

	response := strings.TrimSpace(string(bodyText))
	podID, err := strconv.ParseUint(response, 10, 64)
	if err != nil {
		return defaultLead, err
	}
	if podID == 0 {
		return 0, fmt.Errorf("leader is not elected yet")
	}

	return podID, nil
}

// AddNodeToRaft requests raft membership add via coordinator /add-node endpoint.
func AddNodeToRaft(db DBRaftAddressProvider, coordinatorClientPort, coordinatorPort int, primaryPodName, podName string, nodeID int, user, pass string) (string, error) {
	primaryDNSName := db.GoverningServiceDNS(primaryPodName)
	primaryURL := "http://" + primaryDNSName + ":" + strconv.Itoa(coordinatorClientPort) + "/add-node"

	dnsName := db.GoverningServiceDNS(podName)
	url := "http://" + dnsName + ":" + strconv.Itoa(coordinatorPort)
	node := &raftNodeInfo{
		NodeID: &nodeID,
		URL:    &url,
	}

	return doRaftMembershipChange(http.MethodPost, primaryURL, node, user, pass, "add new node")
}

// RemoveNodeFromRaft requests raft membership remove via coordinator /remove-node endpoint.
func RemoveNodeFromRaft(db DBRaftAddressProvider, coordinatorClientPort int, primaryPodName string, nodeID int, user, pass string) (string, error) {
	primaryDNSName := db.GoverningServiceDNS(primaryPodName)
	primaryURL := "http://" + primaryDNSName + ":" + strconv.Itoa(coordinatorClientPort) + "/remove-node"

	node := &raftNodeInfo{
		NodeID: &nodeID,
	}

	return doRaftMembershipChange(http.MethodDelete, primaryURL, node, user, pass, "remove node")
}

// GetCurrentLeaderPodName returns current leader pod name by resolving raft leader id.
func GetCurrentLeaderPodName(db DBRaftAddressProvider, coordinatorClientPort int, podName, user, pass string) (string, error) {
	leaderID, err := GetCurrentLeaderID(coordinatorClientPort, db.GoverningServiceDNS(podName), user, pass)
	if err != nil {
		return "", fmt.Errorf("failed on get current primary from remote host: %w", err)
	}
	if leaderID < 1 {
		return "", fmt.Errorf("invalid raft leader id: %d", leaderID)
	}
	return fmt.Sprintf("%s-%d", db.OffshootName(), leaderID-1), nil
}

func doRaftMembershipChange(method, endpoint string, node *raftNodeInfo, user, pass, action string) (string, error) {
	requestByte, err := json.Marshal(node)
	if err != nil {
		return "", err
	}
	requestBody := bytes.NewReader(requestByte)

	resp, err := DoRaftRequest(method, endpoint, user, pass, requestBody, leaderAPIRequestTimeout)
	if err != nil {
		return "", fmt.Errorf("failed to %s: %w", action, err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to %s: %w", action, err)
	}
	return string(bodyText), nil
}

func GetRaftLeaderIDWithRetries(db DBRaftAddressProvider, coordinatorClientPort int, dbPodName, user, pass string, maxTries int, retryDelay time.Duration) (int, error) {
	var lastErr error
	for tries := 1; tries <= maxTries; tries++ {
		currentLeaderID, err := GetCurrentLeaderID(coordinatorClientPort, db.GoverningServiceDNS(dbPodName), user, pass)
		if err == nil {
			return int(currentLeaderID), nil
		}
		lastErr = fmt.Errorf("failed on getting current leader: %w", err)
		time.Sleep(retryDelay)
	}
	return 0, fmt.Errorf("failed to get leader of raft cluster: %w", lastErr)
}

func GetRaftPrimaryNode(db DBRaftAddressProvider, coordinatorClientPort int, replicas int, user, pass string, maxTries int, retryDelay time.Duration) (int, error) {
	var lastErr error
	for rep := range replicas {
		podName := fmt.Sprintf("%s-%v", db.OffshootName(), rep)
		primaryPodID, err := GetRaftLeaderIDWithRetries(db, coordinatorClientPort, podName, user, pass, maxTries, retryDelay)
		if err == nil {
			return primaryPodID, nil
		}
		lastErr = err
	}
	return 0, lastErr
}

func AddRaftNodeWithRetries(db DBRaftAddressProvider, coordinatorClientPort, coordinatorPort int, primaryPodName, podName string, nodeID int, user, pass string, maxTries int, retryDelay time.Duration) error {
	var lastErr error
	for tries := 0; tries <= maxTries; tries++ {
		_, err := AddNodeToRaft(db, coordinatorClientPort, coordinatorPort, primaryPodName, podName, nodeID, user, pass)
		if err == nil {
			return nil
		}
		lastErr = err
		time.Sleep(retryDelay)
	}
	return fmt.Errorf("failed to add nodeId = %v to the raft: %w", nodeID, lastErr)
}

func RemoveRaftNodeWithRetries(db DBRaftAddressProvider, coordinatorClientPort int, primaryPodName string, nodeID int, user, pass string, maxTries int, retryDelay time.Duration) error {
	var lastErr error
	for tries := 0; tries <= maxTries; tries++ {
		_, err := RemoveNodeFromRaft(db, coordinatorClientPort, primaryPodName, nodeID, user, pass)
		if err == nil {
			return nil
		}
		lastErr = err
		time.Sleep(retryDelay)
	}
	return fmt.Errorf("failed to remove nodeId = %v from the raft: %w", nodeID, lastErr)
}

func TransferLeadershipByPodNameWithProvider(db DBRaftAddressProvider, coordinatorClientPort int, podName string, transferee int, user, pass string, timeout time.Duration) (string, error) {
	dnsName := db.GoverningServiceDNS(podName)
	endpoint := fmt.Sprintf("http://%s:%d/transfer", dnsName, coordinatorClientPort)
	return TransferLeadership(endpoint, transferee, user, pass, timeout)
}

func AddNodeAsVoterWithPodNameWithProvider(db DBRaftAddressProvider, coordinatorClientPort, coordinatorPort int, nodeID int, podName, user, pass string) (string, error) {
	primaryPodName, err := GetCurrentLeaderPodName(db, coordinatorClientPort, podName, user, pass)
	if err != nil {
		return "", fmt.Errorf("failed while trying to make node a voter: %w", err)
	}
	return AddNodeToRaft(db, coordinatorClientPort, coordinatorPort, primaryPodName, podName, nodeID, user, pass)
}

func TransferLeadershipByPodName(coordinatorClientPort int, governingDNS string, transferee int, user, pass string, timeout time.Duration) (string, error) {
	endpoint := fmt.Sprintf("http://%s:%d/transfer", governingDNS, coordinatorClientPort)
	return TransferLeadership(endpoint, transferee, user, pass, timeout)
}

func AddNodeAsVoterWithPodName(coordinatorClientPort, coordinatorPort int, governingDNS, offshootName string, nodeID int, podName, user, pass string) (string, error) {
	leaderID, err := GetCurrentLeaderID(coordinatorClientPort, governingDNS, user, pass)
	if err != nil {
		return "", fmt.Errorf("failed while trying to make node a voter: %w", err)
	}
	if leaderID < 1 {
		return "", fmt.Errorf("invalid raft leader id: %d", leaderID)
	}

	dotIndex := strings.Index(governingDNS, ".")
	if dotIndex <= 0 || dotIndex+1 >= len(governingDNS) {
		return "", fmt.Errorf("invalid governing DNS: %s", governingDNS)
	}

	primaryPodName := fmt.Sprintf("%s-%d", offshootName, leaderID-1)
	dnsSuffix := governingDNS[dotIndex:]
	primaryDNSName := primaryPodName + dnsSuffix
	primaryURL := "http://" + primaryDNSName + ":" + strconv.Itoa(coordinatorClientPort) + "/add-node"

	dnsName := podName + dnsSuffix
	url := "http://" + dnsName + ":" + strconv.Itoa(coordinatorPort)
	node := &raftNodeInfo{
		NodeID: &nodeID,
		URL:    &url,
	}

	return doRaftMembershipChange(http.MethodPost, primaryURL, node, user, pass, "add new node")
}

func GetPrimaryPods(cacheClient client.Client, pod *core.Pod, primaryRole string) (*core.PodList, error) {
	labelSelector := &metav1.LabelSelector{
		MatchExpressions: []metav1.LabelSelectorRequirement{
			{
				Key:      kubedb.LabelRole,
				Operator: metav1.LabelSelectorOpIn,
				Values:   []string{primaryRole},
			},
			{
				Key:      meta_util.InstanceLabelKey,
				Operator: metav1.LabelSelectorOpIn,
				Values:   []string{pod.Labels[meta_util.InstanceLabelKey]},
			},
			{
				Key:      meta_util.NameLabelKey,
				Operator: metav1.LabelSelectorOpIn,
				Values:   []string{pod.Labels[meta_util.NameLabelKey]},
			},
		},
	}
	pods := &core.PodList{}
	selector, err := metav1.LabelSelectorAsSelector(labelSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to build selector for primary pods: %w", err)
	}
	if err = cacheClient.List(context.TODO(), pods, &client.ListOptions{LabelSelector: selector}); err != nil {
		return nil, fmt.Errorf("failed to list primary pods: %w", err)
	}
	return pods, nil
}
