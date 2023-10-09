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
	"strings"

	kutil "kmodules.xyz/client-go"

	promapi "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	core "k8s.io/api/core/v1"
)

const (
	GroupName                = "monitoring.appscode.com"
	DefaultPrometheusKey     = GroupName + "/is-default-prometheus"
	DefaultAlertmanagerKey   = GroupName + "/is-default-alertmanager"
	DefaultGrafanaKey        = GroupName + "/is-default-grafana"
	PrometheusKey            = GroupName + "/prometheus"
	PrometheusValueAuto      = "auto"
	PrometheusValueFederated = "federated"
)

// +kubebuilder:validation:Enum=prometheus.io/operator;prometheus.io;prometheus.io/builtin
type AgentType string

const (
	KeyAgent   = GroupName + "/agent"
	KeyService = GroupName + "/service"

	VendorPrometheus = "prometheus.io"

	AgentPrometheus         AgentType = VendorPrometheus
	AgentPrometheusBuiltin  AgentType = VendorPrometheus + "/builtin"
	AgentPrometheusOperator AgentType = VendorPrometheus + "/operator"

	PrometheusExporterPortNumber = 56790
	PrometheusExporterPortName   = "metrics"
)

func (at AgentType) Vendor() string {
	return strings.SplitN(string(at), "/", 2)[0]
}

type AgentSpec struct {
	Agent      AgentType       `json:"agent,omitempty"`
	Prometheus *PrometheusSpec `json:"prometheus,omitempty"`
}

type PrometheusSpec struct {
	Exporter       PrometheusExporterSpec `json:"exporter,omitempty"`
	ServiceMonitor *ServiceMonitorSpec    `json:"serviceMonitor,omitempty"`
}

type ServiceMonitorSpec struct {
	// Labels are key value pairs that is used to select Prometheus instance via ServiceMonitor labels.
	// +optional
	Labels map[string]string `json:"labels,omitempty"`

	// Interval at which metrics should be scraped
	// +optional
	Interval string `json:"interval,omitempty"`
}

type PrometheusExporterSpec struct {
	// Port number for the exporter side car.
	// +optional
	// +kubebuilder:default=56790
	Port int32 `json:"port,omitempty"`

	// Arguments to the entrypoint.
	// The docker image's CMD is used if this is not provided.
	// Variable references $(VAR_NAME) are expanded using the container's environment. If a variable
	// cannot be resolved, the reference in the input string will be unchanged. The $(VAR_NAME) syntax
	// can be escaped with a double $$, ie: $$(VAR_NAME). Escaped references will never be expanded,
	// regardless of whether the variable exists or not.
	// Cannot be updated.
	// More info: https://kubernetes.io/docs/tasks/inject-data-application/define-command-argument-container/#running-a-command-in-a-shell
	// +optional
	Args []string `json:"args,omitempty"`

	// List of environment variables to set in the container.
	// Cannot be updated.
	// +optional
	// +patchMergeKey=name
	// +patchStrategy=merge
	Env []core.EnvVar `json:"env,omitempty" patchStrategy:"merge" patchMergeKey:"name"`

	// Compute Resources required by exporter container.
	// Cannot be updated.
	// More info: https://kubernetes.io/docs/concepts/configuration/manage-compute-resources-container/
	// +optional
	Resources core.ResourceRequirements `json:"resources,omitempty"`

	// Security options the pod should run with.
	// More info: https://kubernetes.io/docs/concepts/policy/security-context/
	// More info: https://kubernetes.io/docs/tasks/configure-pod-container/security-context/
	// +optional
	SecurityContext *core.SecurityContext `json:"securityContext,omitempty"`
}

type Agent interface {
	GetType() AgentType
	CreateOrUpdate(sp StatsAccessor, spec *AgentSpec) (kutil.VerbType, error)
	Delete(sp StatsAccessor) (kutil.VerbType, error)
}

type StatsAccessor interface {
	GetNamespace() string
	ServiceName() string
	ServiceMonitorName() string
	ServiceMonitorAdditionalLabels() map[string]string
	Path() string
	// Scheme is used to determine url scheme /metrics endpoint
	Scheme() string
	TLSConfig() *promapi.TLSConfig
}
