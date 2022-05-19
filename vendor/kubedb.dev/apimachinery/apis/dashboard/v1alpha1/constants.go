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
)

const (
	ES_USER_ENV     = "ELASTICSEARCH_USERNAME"
	ES_PASSWORD_ENV = "ELASTICSEARCH_PASSWORD"
	ES_USER_KEY     = "elasticsearch.username"
	ES_PASSWORD_KEY = "elasticsearch.password"
	OS_USER_KEY     = "opensearch.username"
	OS_PASSWORD_KEY = "opensearch.password"

	ElasticsearchDashboardPortServer         = "server"
	ElasticsearchDashboardConfigMergeCommand = "/usr/local/bin/dashboard-config-merger.sh"

	KibanaConfigDir       = "/usr/share/kibana/config"
	KibanaTempConfigDir   = "/kibana/temp-config"
	KibanaCustomConfigDir = "/kibana/custom-config"
	KibanaStatusEndpoint  = "/api/status"
	KibanaConfigFileName  = "kibana.yml"

	OpensearchDashboardsConfigDir       = "/usr/share/opensearch-dashboards/config"
	OpensearchDashboardsTempConfigDir   = "/opensearch-dashboards/temp-config"
	OpensearchDashboardsCustomConfigDir = "/opensearch-dashboards/custom-config"
	OpensearchDashboardsStatusEndpoint  = "/api/status"
	OpensearchDasboardsConfigFileName   = "opensearch_dashboards.yml"

	ElasticsearchHostsKey = "elasticsearch.hosts"
	ElasticsearchSSLCaKey = "elasticsearch.ssl.certificateAuthorities"

	OpensearchHostsKey        = "opensearch.hosts"
	OpensearchSSLCaKey        = "opensearch.ssl.certificateAuthorities"
	OpensearchCookieSecureKey = "opensearch_security.cookie.secure"

	DashboardServerNameKey       = "server.name"
	DashboardServerPortKey       = "server.port"
	DashboardServerHostKey       = "server.host"
	DashboardServerSSLEnabledKey = "server.ssl.enabled"
	DashboardServerSSLCertKey    = "server.ssl.certificate"
	DashboardServerSSLKey        = "server.ssl.key"
	DashboardServerSSLCaKey      = "server.ssl.certificateAuthorities"
	DashboardNodeOptionsKey      = "node.options"
	DashboardMaxOldSpaceFlag     = "--max-old-space-size"

	DashboardDeploymentAvailable           = "MinimumReplicasAvailable"
	DashboardDeploymentNotAvailable        = "MinimumReplicasNotAvailable"
	DashboardServiceReady                  = "ServiceAcceptingRequests"
	DashboardServiceNotReady               = "ServiceNotAcceptingRequests"
	DashboardAcceptingConnectionRequest    = "DashboardAcceptingConnectionRequests"
	DashboardNotAcceptingConnectionRequest = "DashboardNotAcceptingConnectionRequests"
	DashboardReadinessCheckSucceeded       = "DashboardReadinessCheckSucceeded"
	DashboardReadinessCheckFailed          = "DashboardReadinessCheckFailed"
	DashboardOnDeletion                    = "DashboardOnDeletion"

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

	ElasticsearchDashboardRESTPort     = 5601
	ElasticsearchDashboardRESTPortName = "http"
)

var (
	ElasticsearchDashboardGracefulDeletionPeriod = (int64)(time.Duration(3 * time.Second))

	DashboardsDefaultResources = core.ResourceRequirements{
		Requests: core.ResourceList{
			core.ResourceCPU:    resource.MustParse(".100"),
			core.ResourceMemory: resource.MustParse("1Gi"),
		},
		Limits: core.ResourceList{
			core.ResourceMemory: resource.MustParse("1Gi"),
		},
	}
)
