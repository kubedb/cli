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

package v1alpha1

// ReplicationModeDetector is the image for the MySQL replication mode detector
type ReplicationModeDetector struct {
	Image string `json:"image"`
}

// UpgradeConstraints specifies the constraints that need to be considered during version upgrade
type UpgradeConstraints struct {
	// List of all accepted versions for upgrade request.
	// An empty list indicates all versions are accepted except the denylist.
	Allowlist []string `json:"allowlist,omitempty"`
	// List of all rejected versions for upgrade request.
	// An empty list indicates no version is rejected.
	Denylist []string `json:"denylist,omitempty"`
}
