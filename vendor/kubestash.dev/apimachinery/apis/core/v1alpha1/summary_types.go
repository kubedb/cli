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
	kmapi "kmodules.xyz/client-go/api/v1"
)

type Summary struct {
	// Name of the respective BackupSession/RestoreSession
	Name string `json:"name,omitempty"`
	// Namespace of the respective invoker
	Namespace string `json:"namespace,omitempty"`

	// Invoker specifies the information about the invoker which resulted this session
	Invoker *kmapi.TypedObjectReference `json:"invoker,omitempty"`
	// Target specifies the target information that has been backed up /restored in this session
	Target *kmapi.TypedObjectReference `json:"target,omitempty"`
	// Status specifies the backup/restore status for the respective target
	Status TargetStatus `json:"status,omitempty"`
}

type TargetStatus struct {
	// Phase represents the backup/restore phase of the target
	Phase string `json:"phase,omitempty"`
	// Duration represents the amount of time it took to complete the backup/restore for this target.
	Duration string `json:"duration,omitempty"`
	// Error specifies the respective error message in case of backup/restore failure
	Error string `json:"error,omitempty"`
}
