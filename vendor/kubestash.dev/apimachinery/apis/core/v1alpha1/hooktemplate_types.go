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
	"kubestash.dev/apimachinery/apis"

	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ofst "kmodules.xyz/offshoot-api/api/v1"
	prober "kmodules.xyz/prober/api/v1"
)

const (
	ResourceKindHookTemplate     = "HookTemplate"
	ResourceSingularHookTemplate = "hooktemplate"
	ResourcePluralHookTemplate   = "hooktemplates"
)

// +k8s:openapi-gen=true
// +kubebuilder:object:root=true
// +kubebuilder:resource:path=hooktemplates,singular=hooktemplate,categories={kubestash,appscode,all}
// +kubebuilder:printcolumn:name="Executor",type="string",JSONPath=".spec.executor.type"
// +kubebuilder:printcolumn:name="Timeout",type="string",JSONPath=".spec.timeout"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// HookTemplate defines a template for some action that will be executed before or/and after backup/restore process.
// For example, there could be a HookTemplate that pause an application before backup and another HookTemplate
// that resume the application after backup.
// This is a namespaced CRD. However, you can use it from other namespaces. You can control which
// namespaces are allowed to use it using the `usagePolicy` section.
type HookTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec HookTemplateSpec `json:"spec,omitempty"`
}

// HookTemplateSpec defines the template for the operation that will be performed by this hook
type HookTemplateSpec struct {
	// UsagePolicy specifies a policy of how this HookTemplate will be used. For example,
	// you can use `allowedNamespaces` policy to restrict the usage of this HookTemplate to particular namespaces.
	//
	// This field is optional. If you don't provide the usagePolicy, then it can be used only from the current namespace.
	// +optional
	UsagePolicy *apis.UsagePolicy `json:"usagePolicy,omitempty"`

	// Params defines a list of parameters that is used by the HookTemplate to execute its logic.
	// +optional
	Params []apis.ParameterDefinition `json:"params,omitempty"`

	// Action specifies the operation that is performed by this HookTemplate
	// Valid values are:
	// - "exec": Execute command in a shell
	// - "httpGet": Do an HTTP GET request
	// - "httpPost": Do an HTTP POST request
	// - "tcpSocket": Check if a TCP socket open or not
	Action *prober.Handler `json:"action,omitempty"`

	// Executor specifies the entity where the hook will be executed.
	Executor *HookExecutor `json:"executor,omitempty"`
}

// HookExecutor specifies the entity specification which is responsible for executing the hook
type HookExecutor struct {
	// Type indicate the types of entity that will execute the hook.
	// Valid values are:
	// - "Function": KubeStash will create a job with the provided information in `function` section. The job will execute the hook.
	// - "Pod": KubeStash will select the pod that matches the selector provided in `pod` section. This pod(s) will execute the hook.
	// - "Operator": KubeStash operator itself will execute the hook.
	Type HookExecutorType `json:"type,omitempty"`

	// Function specifies the function information which will be used to create the hook executor job.
	// +optional
	Function *FunctionHookExecutorSpec `json:"function,omitempty"`

	// Pod specifies the criteria to use to select the hook executor pods
	// +optional
	Pod *PodHookExecutorSpec `json:"pod,omitempty"`
}

// HookExecutorType specifies the type of entity that will execute the hook
// +kubebuilder:validation:Enum=Function;Pod;Operator
type HookExecutorType string

const (
	HookExecutorFunction HookExecutorType = "Function"
	HookExecutorPod      HookExecutorType = "Pod"
	HookExecutorOperator HookExecutorType = "Operator"
)

// FunctionHookExecutorSpec defines function and its parameters that will be used to create hook executor job
type FunctionHookExecutorSpec struct {
	// Name indicate the name of the Function that contains the container definition for executing the hook logic
	Name string `json:"name,omitempty"`

	// EnvVariables specifies a list of environment variables that will be passed to the executor container
	// +optional
	EnvVariables []core.EnvVar `json:"env,omitempty"`

	// VolumeMounts specifies the volumes mounts for the executor container
	// +optional
	VolumeMounts []core.VolumeMount `json:"volumeMounts,omitempty"`

	// Volumes specifies the volumes that will be mounted in the executor container
	// +optional
	Volumes []ofst.Volume `json:"volumes,omitempty"`
}

// PodHookExecutorSpec specifies the criteria that will be used to select the pod which will execute the hook
type PodHookExecutorSpec struct {
	// Selector specifies list of key value pair that will be used as label selector to select the desired pods.
	// You can use comma to separate multiple labels (i.e. "app=my-app,env=prod")
	Selector string `json:"selector,omitempty"`

	// Owner specifies a template for owner reference that will be used to filter the selected pods.
	// +optional
	Owner *metav1.OwnerReference `json:"owner,omitempty"`

	// Strategy specifies what should be the behavior when multiple pods are selected
	// Valid values are:
	// - "ExecuteOnOne": Execute hook on only one of the selected pods. This is default behavior
	// - "ExecuteOnAll": Execute hook on all the selected pods.
	// +kubebuilder:default=ExecuteOnOne
	Strategy PodHookExecutionStrategy `json:"strategy,omitempty"`
}

// PodHookExecutionStrategy specifies the strategy to follow when multiple pods are selected for hook execution
// +kubebuilder:validation:Enum=ExecuteOnOne;ExecuteOnAll
type PodHookExecutionStrategy string

const (
	ExecuteOnOne PodHookExecutionStrategy = "ExecuteOnOne"
	ExecuteOnAll PodHookExecutionStrategy = "ExecuteOnAll"
)

//+kubebuilder:object:root=true

// HookTemplateList contains a list of HookTemplate
type HookTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HookTemplate `json:"items"`
}

func init() {
	SchemeBuilder.Register(&HookTemplate{}, &HookTemplateList{})
}
