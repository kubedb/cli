/*
Copyright AppsCode Inc. and Contributors

Licensed under the AppsCode Free Trial License 1.0.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://github.com/appscode/licenses/raw/1.0.0/AppsCode-Free-Trial-1.0.0.md

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package lib

import (
	opsapi "kubedb.dev/apimachinery/apis/ops/v1alpha1"

	core "k8s.io/api/core/v1"
	nodemeta "kmodules.xyz/resource-metadata/apis/node/v1alpha1"
)

func SetNodeSelector(policy nodemeta.NodeSelectionPolicy, selector map[string]string, topology *opsapi.Topology) map[string]string {
	if topology == nil {
		return selector
	}
	if policy == nodemeta.NodeSelectionPolicyLabelSelector {
		if len(selector) == 0 {
			selector = make(map[string]string)
		}
		selector[topology.Key] = topology.Value
	}
	return selector
}

func SetToleration(policy nodemeta.NodeSelectionPolicy, tolerations []core.Toleration, topology *opsapi.Topology) []core.Toleration {
	if topology == nil {
		return tolerations
	}

	if policy == nodemeta.NodeSelectionPolicyTaint {
		if tolerations == nil {
			tolerations = make([]core.Toleration, 0)
		}
		tol := core.Toleration{
			Effect: core.TaintEffectNoSchedule,
			Key:    topology.Key,
			Value:  topology.Value,
		}
		for i, t := range tolerations {
			if t.Key == topology.Key && t.Effect == core.TaintEffectNoSchedule {
				tolerations[i] = tol
				return tolerations
			}
		}
		tolerations = append(tolerations, tol)
	}
	return tolerations
}
