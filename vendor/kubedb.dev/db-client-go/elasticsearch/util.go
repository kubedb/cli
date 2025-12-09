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
	"encoding/json"
	"fmt"
	"io"

	"k8s.io/klog/v2"
)

var (
	diskUsageRequestIndex     = "*"
	diskUsageRequestWildcards = "all"
	diskUsageRequestKey       = "store_size_in_bytes"
	diskUsageSafetyFactor     = 0.20 // Add 20%
	diskUsageDefaultMi        = 1024 // 1024 Mi
)

func calculateDatabaseSize(body io.ReadCloser) (string, error) {
	var totalDiskUsageInBytes float64
	resMap := make(map[string]any)
	if err := json.NewDecoder(body).Decode(&resMap); err != nil {
		klog.Errorf("failed to deserialize the response body for disk usage request: %v", err)
		return "", err
	}

	// Parse the deserialized json response to find out storage of each index
	for _, val := range resMap {
		if valMap, ok := val.(map[string]any); ok {
			for key, field := range valMap {
				if key == diskUsageRequestKey {
					storeSizeInByes := field.(float64)
					totalDiskUsageInBytes += storeSizeInByes
				}
			}
		}
	}

	// Add extra 20% percent of extra storage for safety & taking metadata into account.
	// convert bytes to Mib
	totalDiskUsageInMi := int(totalDiskUsageInBytes * (1 + diskUsageSafetyFactor) / (1024 * 1024))
	if totalDiskUsageInMi < diskUsageDefaultMi {
		totalDiskUsageInMi = diskUsageDefaultMi
	}
	return fmt.Sprintf("%dMi", totalDiskUsageInMi), nil
}
