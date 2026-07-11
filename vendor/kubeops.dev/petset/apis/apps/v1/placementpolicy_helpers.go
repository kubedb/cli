/*
Copyright AppsCode Inc. and Contributors.

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

import "fmt"

// Validate checks the DC/DR failover configuration of a ClusterSpreadConstraint.
// It is safe to call when FailoverPolicy is nil (returns nil). Wire it from the
// PlacementPolicy validating webhook.
func (c *ClusterSpreadConstraint) Validate() error {
	if c == nil || c.FailoverPolicy == nil {
		return nil
	}
	members, arbiters, witnesses := 0, 0, 0
	for i := range c.DistributionRules {
		r := &c.DistributionRules[i]
		switch r.Role {
		case "", DCRoleMember:
			members++
			if len(r.ReplicaIndices) == 0 {
				return fmt.Errorf("distributionRule for cluster %q has role Member but no replicaIndices", r.ClusterName)
			}
		case DCRoleArbiter:
			arbiters++
			if len(r.ReplicaIndices) != 0 {
				return fmt.Errorf("distributionRule for cluster %q has role Arbiter but carries replicaIndices (an arbiter holds no data)", r.ClusterName)
			}
		case DCRoleWitness:
			witnesses++
			if len(r.ReplicaIndices) == 0 {
				return fmt.Errorf("distributionRule for cluster %q has role Witness but no replicaIndices (a witness is data bearing)", r.ClusterName)
			}
		default:
			return fmt.Errorf("distributionRule for cluster %q has unknown role %q", r.ClusterName, r.Role)
		}
	}
	if members < 2 {
		return fmt.Errorf("DC/DR requires at least two Member data centers, found %d", members)
	}

	fp := c.FailoverPolicy
	switch fp.Trigger.Scope {
	case FailoverScopeGlobal:
		if fp.Trigger.Group != "" {
			return fmt.Errorf("trigger.group must be empty when scope is Global")
		}
	case FailoverScopeGroup:
		if fp.Trigger.Group == "" {
			return fmt.Errorf("trigger.group is required when scope is Group")
		}
	default:
		return fmt.Errorf("trigger.scope must be Global or Group, got %q", fp.Trigger.Scope)
	}

	switch fp.Mode {
	case "":
		// derived, nothing to check
	case FailoverModeTwoDC:
		if members != 2 || (arbiters+witnesses) < 1 {
			return fmt.Errorf("mode TwoDC requires exactly two Members and at least one Arbiter or Witness, found members=%d arbiters=%d witnesses=%d", members, arbiters, witnesses)
		}
	case FailoverModeThreeDC:
		if members < 3 {
			return fmt.Errorf("mode ThreeDC requires at least three Members, found %d", members)
		}
	default:
		return fmt.Errorf("mode must be TwoDC or ThreeDC, got %q", fp.Mode)
	}
	return nil
}
