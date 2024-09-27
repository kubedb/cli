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

package v1

type GrafanaConfig struct {
	URL         string        `json:"url"`
	Service     ServiceSpec   `json:"service"`
	BasicAuth   BasicAuth     `json:"basicAuth"`
	BearerToken string        `json:"bearerToken"`
	TLS         TLSConfig     `json:"tls"`
	Dashboard   DashboardSpec `json:"dashboard"`
}

type PrometheusConfig struct {
	URL         string      `json:"url"`
	Service     ServiceSpec `json:"service"`
	BasicAuth   BasicAuth   `json:"basicAuth"`
	BearerToken string      `json:"bearerToken"`
	TLS         TLSConfig   `json:"tls"`
}

type ServiceSpec struct {
	Scheme    string `json:"scheme"`
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Port      string `json:"port"`
	Path      string `json:"path"`
	Query     string `json:"query"`
}

type DashboardSpec struct {
	Datasource string `json:"datasource"`
	FolderID   int    `json:"folderID"`
}

type BasicAuth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type TLSConfig struct {
	Ca                    string `json:"ca"`
	Cert                  string `json:"cert"`
	Key                   string `json:"key"`
	ServerName            string `json:"serverName"`
	InsecureSkipTLSVerify bool   `json:"insecureSkipTLSVerify"`
}

type PrometheusContext struct {
	HubUID     string `json:"hubUID,omitempty"`
	ClusterUID string `json:"clusterUID"`
	ProjectId  string `json:"projectId,omitempty"`
	Default    bool   `json:"default"`
}

type GrafanaContext struct {
	FolderID   *int64 `json:"folderID,omitempty"`
	Datasource string `json:"datasource,omitempty"`
}
