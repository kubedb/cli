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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kmapi "kmodules.xyz/client-go/api/v1"
	rscoreapi "kmodules.xyz/resource-metadata/apis/core/v1alpha1"
)

const (
	ResourceKindDatabaseSummary = "DatabaseSummary"
	ResourceDatabaseSummary     = "databasesummary"
	ResourceDatabaseSummaries   = "databasesummaries"
)

// DatabaseSummary is a request/response payload, not a stored Kubernetes
// object -- it carries no metav1.ObjectMeta, matching how upstream
// Kubernetes shapes the same kind of type (e.g.
// k8s.io/api/admission/v1.AdmissionReview). Its generated client is built
// on kmodules.xyz/client-go/gentype's create-only Client[T runtime.Object]
// instead of k8s.io/client-go/gentype's Client[T], which requires
// metav1.Object -- see that package's doc comment. This only works because
// DatabaseSummary is create-only (+genclient:onlyVerbs=create); it carries
// no +kubebuilder:object:root=true/+kubebuilder:resource markers since
// controller-gen needs ObjectMeta to treat a type as a CRD root, and
// ui-server serves this type via API aggregation, not a CRD.
// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=get,list,update,delete,watch
// +genclient:onlyVerbs=create
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type DatabaseSummary struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	Request *DatabaseSummaryRequest `json:"request,omitempty"`
	// +optional
	Response *DatabaseSummaryResponse `json:"response,omitempty"`
}

type DatabaseSummaryRequest struct {
	Source kmapi.ObjectInfo `json:"source"`
}

type DatabaseSummaryResponse struct {
	Cluster    AggregatedStats               `json:"cluster,omitempty"`
	Namespaces map[string]DatabaseIndexEntry `json:"namespaces,omitempty"`
	Kinds      map[string]DatabaseIndexEntry `json:"kinds,omitempty"`
	Stats      []StatForSpecificDB           `json:"stats,omitempty"`
}

type DatabaseIndexEntry struct {
	AggregatedStats `json:",inline"`
	Indices         []int `json:"indices,omitempty"`
}

type AggregatedStats struct {
	Total             int                            `json:"total,omitempty"`
	Resources         rscoreapi.ResourceRequirements `json:"resources,omitempty"`
	GatewayExposed    int                            `json:"gatewayExposed,omitempty"`
	BackupEnabled     int                            `json:"backupEnabled,omitempty"`
	MonitoringEnabled int                            `json:"monitoringEnabled,omitempty"`
}

type StatForSpecificDB struct {
	Name                 string                         `json:"name,omitempty"`
	Namespace            string                         `json:"namespace,omitempty"`
	Kind                 string                         `json:"kind,omitempty"`
	Resources            rscoreapi.ResourceRequirements `json:"resources,omitempty"`
	GatewayExposed       bool                           `json:"gatewayExposed"`
	BackupEnabled        bool                           `json:"backupEnabled"`
	MonitoringEnabled    bool                           `json:"monitoringEnabled"`
	LastSuccessfulBackup *metav1.Time                   `json:"lastSuccessfulBackup,omitempty"`
	LastBackupPhase      string                         `json:"lastBackupPhase,omitempty"`
}
