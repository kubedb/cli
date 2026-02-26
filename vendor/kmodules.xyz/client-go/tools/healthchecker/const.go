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

type HealthCheckFailureLabel string

const (
	HealthCheckClientFailure                HealthCheckFailureLabel = "ClientFailure"
	HealthCheckPingFailure                  HealthCheckFailureLabel = "PingFailure"
	HealthCheckWriteFailure                 HealthCheckFailureLabel = "WriteFailure"
	HealthCheckReadFailure                  HealthCheckFailureLabel = "ReadFailure"
	HealthCheckPrimaryFailure               HealthCheckFailureLabel = "PrimaryFailure"
	HealthCheckSecondaryFailure             HealthCheckFailureLabel = "SecondaryFailure"
	HealthCheckSecondaryUnusualLocked       HealthCheckFailureLabel = "SecondaryUnusualLocked"
	HealthCheckSecondaryLockCheckingFailure HealthCheckFailureLabel = "SecondaryLockCheckingFailure"
	HealthCheckKubernetesClientFailure      HealthCheckFailureLabel = "KubernetesClientFailure"

	// replica
	HealthCheckReplicaFailure HealthCheckFailureLabel = "ReplicaFailure"

	// MariaDB Constants
	HealthCheckClusterFailure HealthCheckFailureLabel = "ClusterFailure"

	// Redis Constants
	HealthCheckClusterSlotFailure   HealthCheckFailureLabel = "ClusterSlotFailure"
	HealthCheckNodesNotReadyFailure HealthCheckFailureLabel = "NodesNotReadyFailure"

	// Write Check Constants
	KubeDBSystemDatabase  = "kubedb_system"
	KubeDBWriteCheckTable = "kubedb_write_check"
)
