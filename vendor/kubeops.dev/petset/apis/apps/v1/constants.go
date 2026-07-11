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

const (
	ManifestWorkRoleLabel        = "open-cluster-management.io/role"
	ManifestWorkClusterNameLabel = "open-cluster-management.io/cluster-name"
	RolePod                      = "pod"
	RolePVC                      = "pvc"
	DeletionPolicyAnnotation     = GroupName + "/deletion-policy"
	DeletionPolicyOrphan         = "Orphan"
	// PetSetNameLabel identifies the PetSet that owns a ManifestWork. Ownership must
	// not be resolved by selector alone: two PetSets of the same application (for
	// example a data set and its arbiter) can have subset-overlapping selectors, and
	// label-subset matching then attributes each other's ManifestWorks to both.
	PetSetNameLabel = GroupName + "/petset-name"
)
