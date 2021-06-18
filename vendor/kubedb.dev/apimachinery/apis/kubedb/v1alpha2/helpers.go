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

package v1alpha2

import (
	"fmt"

	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/labels"
	appslister "k8s.io/client-go/listers/apps/v1"
	apps_util "kmodules.xyz/client-go/apps/v1"
	ofst "kmodules.xyz/offshoot-api/api/v1"
)

func checkReplicas(lister appslister.StatefulSetNamespaceLister, selector labels.Selector, expectedItems int) (bool, string, error) {
	items, err := lister.List(selector)
	if err != nil {
		return false, "", err
	}
	if len(items) < expectedItems {
		return false, fmt.Sprintf("All StatefulSets are not available. Desire number of StatefulSet: %d, Available: %d", expectedItems, len(items)), nil
	}

	// return isReplicasReady, message, error
	ready, msg := apps_util.StatefulSetsAreReady(items)
	return ready, msg, nil
}

// HasServiceTemplate returns "true" if the desired serviceTemplate provided in "aliaS" is present in the serviceTemplate list.
// Otherwise, it returns "false".
func HasServiceTemplate(templates []NamedServiceTemplateSpec, alias ServiceAlias) bool {
	for i := range templates {
		if templates[i].Alias == alias {
			return true
		}
	}
	return false
}

// GetServiceTemplate returns a pointer to the desired serviceTemplate referred by "aliaS". Otherwise, it returns nil.
func GetServiceTemplate(templates []NamedServiceTemplateSpec, alias ServiceAlias) ofst.ServiceTemplateSpec {
	for i := range templates {
		c := templates[i]
		if c.Alias == alias {
			return c.ServiceTemplateSpec
		}
	}
	return ofst.ServiceTemplateSpec{}
}

func SetDefaultResourceLimits(req *core.ResourceRequirements, defaultResources core.ResourceRequirements) {
	// if request is set,
	//		- limit set:
	//			- return max(limit,request)
	// else if limit set:
	//		- return limit
	// else
	//		- return default
	calLimit := func(name core.ResourceName, defaultValue resource.Quantity) resource.Quantity {
		if r, ok := req.Requests[name]; ok {
			// l is greater than r == 1.
			if l, exist := req.Limits[name]; exist && l.Cmp(r) == 1 {
				return l
			}
			return r
		}
		if l, ok := req.Limits[name]; ok {
			return l
		}
		return defaultValue
	}
	// if request is not set,
	//		- if limit exists:
	//				- copy limit
	//		- else
	//				- set default
	// else
	// 		- return request
	// endif
	calRequest := func(name core.ResourceName, defaultValue resource.Quantity) resource.Quantity {
		if r, ok := req.Requests[name]; !ok {
			if l, exist := req.Limits[name]; exist {
				return l
			}
			return defaultValue
		} else {
			return r
		}
	}

	if req.Limits == nil {
		req.Limits = core.ResourceList{}
	}
	if req.Requests == nil {
		req.Requests = core.ResourceList{}
	}

	// Calculate the limits first
	for l := range defaultResources.Limits {
		req.Limits[l] = calLimit(l, defaultResources.Limits[l])
	}

	// Once the limit is calculated, Calculate requests
	for r := range defaultResources.Requests {
		req.Requests[r] = calRequest(r, defaultResources.Requests[r])
	}
}

// Upsert elements to string slice
func upsertStringSlice(inSlice []string, values ...string) []string {
	upsert := func(m string) {
		for _, v := range inSlice {
			if v == m {
				return
			}
		}
		inSlice = append(inSlice, m)
	}

	for _, value := range values {
		upsert(value)
	}
	return inSlice
}
