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

type MonitoringPresets struct {
	Spec MonitoringPresetsSpec `json:"spec,omitempty"`
	Form MonitoringPresetsForm `json:"form,omitempty"`
}

type MonitoringPresetsSpec struct {
	Monitoring ServiceMonitorPreset `json:"monitoring"`
}

type ServiceMonitorPreset struct {
	Agent          string               `json:"agent"`
	ServiceMonitor ServiceMonitorLabels `json:"serviceMonitor"`
}

type ServiceMonitorLabels struct {
	// +optional
	Labels map[string]string `json:"labels"`
}

type MonitoringPresetsForm struct {
	Alert AlertPreset `json:"alert"`
}

type AlertPreset struct {
	Enabled SeverityFlag `json:"enabled"`
	// +optional
	Labels map[string]string `json:"labels"`
}

// +kubebuilder:validation:Enum=none;critical;warning;info
type SeverityFlag string

const (
	SeverityFlagNone     SeverityFlag = "none"     // 0
	SeverityFlagCritical SeverityFlag = "critical" // 1
	SeverityFlagWarning  SeverityFlag = "warning"  // 2
	SeverityFlagInfo     SeverityFlag = "info"     // 3
)
