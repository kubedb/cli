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

import (
	"time"

	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var DefaultResources = core.ResourceRequirements{
	Requests: core.ResourceList{
		core.ResourceCPU:    resource.MustParse(".500"),
		core.ResourceMemory: resource.MustParse("1024Mi"),
	},
	Limits: core.ResourceList{
		core.ResourceMemory: resource.MustParse("1024Mi"),
	},
}

const (
	ES_USER_ENV     = "ELASTICSEARCH_USERNAME"
	ES_PASSWORD_ENV = "ELASTICSEARCH_PASSWORD"

	ElasticsearchDashboardPortServer            = "server"
	ElasticsearchDashboardKibanaConfigDir       = "/usr/share/kibana/config"
	ElasticsearchDashboardConfigMergeCommand    = "/usr/local/bin/dashboard-config-merger.sh"
	ElasticsearchDashboardKibanaTempConfigDir   = "/kibana/temp-config"
	ElasticsearchDashboardKibanaCustomConfigDir = "/kibana/custom-config"

	KibanaStatusEndpoint = "/api/status"
	KibanaConfigFileName = "kibana.yml"

	DashboardDeploymentAvailable           = "MinimumReplicasAvailable"
	DashboardDeploymentNotAvailable        = "MinimumReplicasNotAvailable"
	DashboardServiceReady                  = "ServiceAcceptingRequests"
	DashboardServiceNotReady               = "ServiceNotAcceptingRequests"
	DashboardAcceptingConnectionRequest    = "DashboardAcceptingConnectionRequests"
	DashboardNotAcceptingConnectionRequest = "DashboardNotAcceptingConnectionRequests"
	DashboardReadinessCheckSucceeded       = "DashboardReadinessCheckSucceeded"
	DashboardReadinessCheckFailed          = "DashboardReadinessCheckFailed"

	DashboardStateGreen  = "ServerHealthGood"
	DashboardStateYellow = "ServerHealthCritical"
	DashboardStateRed    = "ServerUnhealthy"

	DBNotFound = "DatabaseNotFound"
	DBNotReady = "DatabaseNotReady"

	ComponentDashboard                  = "dashboard"
	CaCertKey                           = "ca.crt"
	DefaultElasticsearchClientCertAlias = "archiver"
	HealthCheckInterval                 = 10 * time.Second
	GlobalHost                          = "0.0.0.0"
)

var (
	ElasticsearchDashboardPropagationPolicy      = meta.DeletePropagationForeground
	ElasticsearchDashboardDefaultPort            = (int32)(5601)
	ElasticsearchDashboardGracefulDeletionPeriod = (int64)(time.Duration(3 * time.Second))
)
