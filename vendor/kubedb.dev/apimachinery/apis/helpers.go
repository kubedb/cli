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

package apis

import (
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"kmodules.xyz/client-go/apiextensions"
)

const (
	Finalizer = "kubedb.com"
)

type ResourceInfo interface {
	ResourceFQN() string
	ResourceShortCode() string
	ResourceKind() string
	ResourceSingular() string
	ResourcePlural() string
	CustomResourceDefinition() *apiextensions.CustomResourceDefinition
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
			if l, exist := req.Limits[name]; exist && l.Cmp(r) > 0 {
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
	calRequest := func(name core.ResourceName, defaultValue resource.Quantity, originalLimit resource.Quantity) resource.Quantity {
		if r, ok := req.Requests[name]; ok {
			return r
		}
		if originalLimit.Value() > 0 {
			// If original Limits existed, use them for Requests
			return originalLimit
		}
		return defaultValue
	}

	if req.Limits == nil {
		req.Limits = core.ResourceList{}
	}
	if req.Requests == nil {
		req.Requests = core.ResourceList{}
	}

	// Store the original Limits to differentiate between newly set and existing values
	originalLimits := make(map[core.ResourceName]resource.Quantity)
	for l := range req.Limits {
		originalLimits[l] = req.Limits[l]
	}

	// Calculate limits first
	for l := range defaultResources.Limits {
		req.Limits[l] = calLimit(l, defaultResources.Limits[l])
	}

	// Calculate requests after limits
	for r := range defaultResources.Requests {
		req.Requests[r] = calRequest(r, defaultResources.Requests[r], originalLimits[r])
	}
}
