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

package healthchecker

import "k8s.io/klog/v2"

type HealthCard struct {
	lastFailure  HealthCheckFailureLabel
	totalFailure int32
	threshold    int32
	clientCount  int32
	key          string
}

func newHealthCard(key string, threshold int32) *HealthCard {
	return &HealthCard{
		threshold: threshold,
		key:       key,
	}
}

// SetThreshold sets the current failure threshold.
// Call this function on the start of each health check.
func (hcf *HealthCard) SetThreshold(threshold int32) {
	hcf.threshold = threshold
}

// HasFailed returns true or false based on the threshold.
// Update the health check condition if this function returns true.
func (hcf *HealthCard) HasFailed(label HealthCheckFailureLabel, err error) bool {
	if hcf.lastFailure == label {
		hcf.totalFailure++
	} else {
		hcf.totalFailure = 1
	}
	hcf.lastFailure = label
	klog.V(5).InfoS("Health check failed for database", "Key", hcf.key, "FailureType", hcf.lastFailure, "Error", err.Error(), "TotalFailure", hcf.totalFailure)
	return hcf.totalFailure >= hcf.threshold
}

// Clear is used to reset the error counter.
// Call this method after each successful health check.
func (hcf *HealthCard) Clear() {
	hcf.totalFailure = 0
	hcf.lastFailure = ""
}

// ClientCreated is used to track the client which are created on the health check.
// Call this method after a client is successfully created in the health check.
func (hcf *HealthCard) ClientCreated() {
	hcf.clientCount++
}

// ClientClosed is used to track the client which are closed on the health check.
// Call this method after a client is successfully closed in the health check.
func (hcf *HealthCard) ClientClosed() {
	hcf.clientCount--
}

// GetClientCount is used to get the current open client count.
// This should always be 0.
func (hcf *HealthCard) GetClientCount() int32 {
	return hcf.clientCount
}
