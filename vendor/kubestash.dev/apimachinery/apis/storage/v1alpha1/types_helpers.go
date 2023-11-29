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
	"gomodules.xyz/x/filepath"
	core "k8s.io/api/core/v1"
)

// ToVolumeAndMount returns volumes and mounts for local backend
func (l LocalSpec) ToVolumeAndMount(storageName string) (core.Volume, core.VolumeMount) {
	vol := core.Volume{
		Name:         storageName,
		VolumeSource: *l.VolumeSource.ToAPIObject(),
	}
	mnt := core.VolumeMount{
		Name:      storageName,
		MountPath: l.MountPath,
		SubPath:   l.SubPath,
	}
	return vol, mnt
}

func (l LocalSpec) ToLocalMountPath(storageName string) (string, error) {
	_, mnt := l.ToVolumeAndMount(storageName)
	return filepath.SecureJoin("/", storageName, mnt.MountPath)
}
