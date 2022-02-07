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

// +kubebuilder:validation:Enum=Provisioning;Ready;Critical;NotReady
type DashboardPhase string

const (
	// used for Dashboards that are currently provisioning
	DashboardPhaseProvisioning DashboardPhase = "Provisioning"
	// used for Dashboards that are currently ReplicaReady, AcceptingConnection and Ready
	DashboardPhaseReady DashboardPhase = "Ready"
	// used for Dashboards that can connect, ReplicaReady == false || Ready == false (eg, ES yellow)
	DashboardPhaseCritical DashboardPhase = "Critical"
	// used for Dashboards that can't connect
	DashboardPhaseNotReady DashboardPhase = "NotReady"
)

// +kubebuilder:validation:Enum=DeploymentReconciled;ServiceReconciled;DashboardProvisioned;ServerAcceptingConnection;ServerHealthy
type DashboardConditionType string

const (
	DashboardConditionDeploymentReconciled DashboardConditionType = "DeploymentReconciled"
	DashboardConditionServiceReconciled    DashboardConditionType = "ServiceReconciled"
	DashboardConditionProvisioned          DashboardConditionType = "DashboardProvisioned"
	DashboardConditionAcceptingConnection  DashboardConditionType = "ServerAcceptingConnection"
	DashboardConditionServerHealthy        DashboardConditionType = "ServerHealthy"
)

const (
	DashboardDeploymentAvailable           = "MinimumReplicasAvailable"
	DashboardDeploymentNotAvailable        = "MinimumReplicasNotAvailable"
	DashboardServiceReady                  = "ServiceAcceptingRequests"
	DashboardServiceNotReady               = "ServiceNotAcceptingRequests"
	DashboardAcceptingConnectionRequest    = "DashboardAcceptingConnectionRequests"
	DashboardNotAcceptingConnectionRequest = "DashboardNotAcceptingConnectionRequests"
	DashboardReadinessCheckSucceeded       = "DashboardReadinessCheckSucceeded"
	DashboardReadinessCheckFailed          = "DashboardReadinessCheckFailed"
	DashboardStateGreen                    = "ServerHealthGood"
	DashboardStateYellow                   = "ServerHealthCritical"
	DashboardStateRed                      = "ServerUnhealthy"
)

// +kubebuilder:validation:Enum=Available;OK;Warning;Error
type DashboardStatus string

const (
	Available     DashboardStatus = "Available"
	StatusOK      DashboardStatus = "OK"
	StatusWarning DashboardStatus = "Warning"
	StatusError   DashboardStatus = "Error"
)

// +kubebuilder:validation:Enum=ca;database-client;kibana-server;dashboard-config
type ElasticsearchDashboardCertificateAlias string

const (
	ElasticsearchDashboardCACert           ElasticsearchDashboardCertificateAlias = "ca"
	ElasticsearchDatabaseClient            ElasticsearchDashboardCertificateAlias = "database-client"
	ElasticsearchDashboardKibanaServerCert ElasticsearchDashboardCertificateAlias = "kibana-server"
	ElasticsearchDashboardConfig           ElasticsearchDashboardCertificateAlias = "dashboard-config"
)

// +kubebuilder:validation:Enum=primary;stats
type ServiceAlias string

const (
	PrimaryServiceAlias ServiceAlias = "primary"
	StatsServiceAlias   ServiceAlias = "stats"
)

// +kubebuilder:validation:Enum=green;yellow;red
type DashboardServerState string

const (
	StateGreen  DashboardServerState = "green"
	StateYellow DashboardServerState = "yellow"
	StateRed    DashboardServerState = "red"
)
