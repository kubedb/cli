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
	v1 "kmodules.xyz/monitoring-agent-api/api/v1"

	"k8s.io/apimachinery/pkg/conversion"
)

func Convert_v1alpha1_AgentSpec_To_v1_AgentSpec(in *AgentSpec, out *v1.AgentSpec, s conversion.Scope) error {
	out.Agent = v1.AgentType(in.Agent)
	if in.Prometheus != nil {
		out := &out.Prometheus
		*out = new(v1.PrometheusSpec)
		(*out).Exporter = v1.PrometheusExporterSpec{
			Port:            in.Prometheus.Port,
			Args:            in.Args,
			Env:             in.Env,
			Resources:       in.Resources,
			SecurityContext: in.SecurityContext,
		}
		(*out).ServiceMonitor = &v1.ServiceMonitorSpec{
			Labels:   in.Prometheus.Labels,
			Interval: in.Prometheus.Interval,
		}
	} else {
		out.Prometheus = nil
	}
	return nil
}

func Convert_v1_PrometheusSpec_To_v1alpha1_PrometheusSpec(in *v1.PrometheusSpec, out *PrometheusSpec, s conversion.Scope) error {
	out.Port = in.Exporter.Port
	// out.Namespace = ""
	out.Labels = in.ServiceMonitor.Labels
	out.Interval = in.ServiceMonitor.Interval
	return nil
}

func Convert_v1alpha1_PrometheusSpec_To_v1_PrometheusSpec(in *PrometheusSpec, out *v1.PrometheusSpec, s conversion.Scope) error {
	out.Exporter = v1.PrometheusExporterSpec{
		Port: in.Port,
	}
	out.ServiceMonitor = &v1.ServiceMonitorSpec{
		Labels:   in.Labels,
		Interval: in.Interval,
	}
	return nil
}
