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

func (agent *AgentSpec) SetDefaults() {
	if agent == nil {
		return
	}

	if agent.Agent.Vendor() == VendorPrometheus {
		if agent.Prometheus == nil {
			agent.Prometheus = &PrometheusSpec{}
		}
		if agent.Prometheus.Exporter.Port == 0 {
			agent.Prometheus.Exporter.Port = PrometheusExporterPortNumber
		}
	}
}

func IsKnownAgentType(at AgentType) bool {
	switch at {
	case AgentPrometheus,
		AgentPrometheusOperator,
		AgentPrometheusBuiltin:
		return true
	}
	return false
}
