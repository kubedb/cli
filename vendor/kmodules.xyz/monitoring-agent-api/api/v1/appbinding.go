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

import (
	kmapi "kmodules.xyz/client-go/api/v1"
)

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
	HubUID      string `json:"hubUID,omitempty"`
	ClusterUID  string `json:"clusterUID"`
	ProjectId   string `json:"projectId,omitempty"`
	Default     bool   `json:"default"`
	IssueToken  bool   `json:"issueToken,omitempty"`
	ClientOrgID string `json:"clientOrgID,omitempty"`
}

type GrafanaContext struct {
	FolderID   *int64 `json:"folderID,omitempty"`
	Datasource string `json:"datasource,omitempty"`
}

type Prometheus struct {
	AppBindingRef   *kmapi.ObjectReference `json:"appBindingRef,omitempty"`
	*ConnectionSpec `json:",inline,omitempty"`
}

// ConnectionSpec is the spec for app
type ConnectionSpec struct {
	// ClientConfig defines how to communicate with the app.
	// Required
	ClientConfig `json:",inline"`

	// Secret is the name of the secret to create in the AppBinding's
	// namespace that will hold the credentials associated with the AppBinding.
	AuthSecret *kmapi.ObjectReference `json:"authSecret,omitempty"`

	// TLSSecret is the name of the secret that will hold
	// the client certificate and private key associated with the AppBinding.
	TLSSecret *kmapi.ObjectReference `json:"tlsSecret,omitempty"`
}

// ClientConfig contains the information to make a connection with an app
type ClientConfig struct {
	// `url` gives the location of the app, in standard URL form
	// (`[scheme://]host:port/path`). Exactly one of `url` or `service`
	// must be specified.
	// +optional
	URL string `json:"url"`

	// InsecureSkipTLSVerify disables TLS certificate verification when communicating with this app.
	// This is strongly discouraged.  You should use the CABundle instead.
	InsecureSkipTLSVerify bool `json:"insecureSkipTLSVerify,omitempty"`

	// CABundle is a PEM encoded CA bundle which will be used to validate the serving certificate of this app.
	// +optional
	CABundle []byte `json:"caBundle,omitempty"`

	// ServerName is used to verify the hostname on the returned
	// certificates unless InsecureSkipVerify is given. It is also included
	// in the client's handshake to support virtual hosting unless it is
	// an IP address.
	ServerName string `json:"serverName,omitempty"`
}
