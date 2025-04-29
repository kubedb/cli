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
	"encoding/json"
)

type Restriction struct {
	VersionConstraint string `json:"versionConstraint"`
	// +optional
	Distributions []string `json:"distributions,omitempty"`
}

type Restrictions []Restriction

type LicenseRestrictions map[string]Restrictions

// LicenseRestrictionsV1 Deprecated, use LicenseRestrictions
type LicenseRestrictionsV1 map[string]Restriction

func (lr *LicenseRestrictions) UnmarshalJSON(data []byte) error {
	// First, try to unmarshal directly into the target type.
	if err := json.Unmarshal(data, (*map[string]Restrictions)(lr)); err == nil {
		return nil
	}

	// If direct unmarshaling fails, try unmarshaling to v1 type.
	var temp LicenseRestrictionsV1
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	*lr = make(map[string]Restrictions, len(temp))
	for key, restriction := range temp {
		(*lr)[key] = Restrictions{restriction}
	}
	return nil
}

func (lr LicenseRestrictions) ToV1() LicenseRestrictionsV1 {
	out := make(LicenseRestrictionsV1, len(lr))
	for key, restrictions := range lr {
		if len(restrictions) > 0 {
			out[key] = restrictions[0]
		}
	}
	return out
}
