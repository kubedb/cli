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

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

const (
	ResourceKindDruidTasks = "DruidTasks"
	ResourceDruidTasks     = "druidtasks"
)

type DruidTasksSpec struct {
	CollectedAt    *metav1.Time     `json:"collectedAt,omitempty"`
	Summary        DruidTaskSummary `json:"summary"`
	Tasks          []DruidTask      `json:"tasks,omitempty"`
	TasksTruncated bool             `json:"tasksTruncated,omitempty"`
}

type DruidTaskSummary struct {
	TotalTasks      *int64 `json:"totalTasks,omitempty"`
	RunningTasks    *int64 `json:"runningTasks,omitempty"`
	WaitingTasks    *int64 `json:"waitingTasks,omitempty"`
	PendingTasks    *int64 `json:"pendingTasks,omitempty"`
	FailedTasks     *int64 `json:"failedTasks,omitempty"`
	SuccessfulTasks *int64 `json:"successfulTasks,omitempty"`
}

type DruidTask struct {
	ID                    string       `json:"id"`
	GroupID               string       `json:"groupID,omitempty"`
	Type                  string       `json:"type,omitempty"`
	Datasource            string       `json:"datasource,omitempty"`
	Status                string       `json:"status,omitempty"`
	RunnerStatus          string       `json:"runnerStatus,omitempty"`
	CreatedAt             *metav1.Time `json:"createdAt,omitempty"`
	QueueInsertionTime    *metav1.Time `json:"queueInsertionTime,omitempty"`
	DurationMillis        *int64       `json:"durationMillis,omitempty"`
	Location              string       `json:"location,omitempty"`
	Host                  string       `json:"host,omitempty"`
	PlaintextPort         *int32       `json:"plaintextPort,omitempty"`
	TLSPort               *int32       `json:"tlsPort,omitempty"`
	ErrorMessage          string       `json:"errorMessage,omitempty"`
	ErrorMessageTruncated bool         `json:"errorMessageTruncated,omitempty"`
}

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type DruidTasks struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec DruidTasksSpec `json:"spec,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type DruidTasksList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DruidTasks `json:"items"`
}
