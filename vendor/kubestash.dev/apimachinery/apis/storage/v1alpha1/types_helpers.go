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
	"fmt"
	"strconv"

	"gomodules.xyz/x/filepath"
	core "k8s.io/api/core/v1"
)

// ToVolumeAndMount returns volumes and mounts for local backend
func (l LocalSpec) ToVolumeAndMount(storageName string) (core.Volume, core.VolumeMount) {
	vol := core.Volume{
		Name:         storageName,
		VolumeSource: *l.ToAPIObject(),
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

func ConvertSizeToByte(sizeWithUnit []string) (uint64, error) {
	numeral, err := strconv.ParseFloat(sizeWithUnit[0], 64)
	if err != nil {
		return 0, err
	}

	switch sizeWithUnit[1] {
	case "TiB":
		return uint64(numeral * (1 << 40)), nil
	case "GiB":
		return uint64(numeral * (1 << 30)), nil
	case "MiB":
		return uint64(numeral * (1 << 20)), nil
	case "KiB":
		return uint64(numeral * (1 << 10)), nil
	case "B":
		return uint64(numeral), nil
	default:
		return 0, fmt.Errorf("no valid unit matched")
	}
}

func FormatBytes(c uint64) string {
	b := float64(c)
	switch {
	case c > 1<<40:
		return fmt.Sprintf("%.3f TiB", b/(1<<40))
	case c > 1<<30:
		return fmt.Sprintf("%.3f GiB", b/(1<<30))
	case c > 1<<20:
		return fmt.Sprintf("%.3f MiB", b/(1<<20))
	case c > 1<<10:
		return fmt.Sprintf("%.3f KiB", b/(1<<10))
	default:
		return fmt.Sprintf("%d B", c)
	}
}
