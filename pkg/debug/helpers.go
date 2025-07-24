/*
Copyright AppsCode Inc. and Contributors

Licensed under the AppsCode Community License 1.0.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://github.com/appscode/licenses/raw/1.0.0/AppsCode-Community-1.0.0.md

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package debug

import (
	"bytes"
	shell "gomodules.xyz/go-sh"
	"k8s.io/klog/v2"
	"os"
	"path"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/yaml"
)

func runHelmCommands(name, ns, fullPath string) error {
	sh := shell.NewSession()
	sh.ShowCMD = true

	buf := new(bytes.Buffer)
	command := []interface{}{
		"ls", "-n", ns,
	}

	err := sh.Command("helm", command...).WriteStdout(path.Join(fullPath, "info.txt"))
	//err := sess.Run()
	if err != nil {
		klog.ErrorS(err, "Failed to run command", "command", command)
	}

	err = os.WriteFile(path.Join(fullPath, "info.txt"), buf.Bytes(), filePerm)
	if err != nil {
		klog.ErrorS(err, "Failed to write info.txt")
	}

	command = []interface{}{
		"get", "values", "-n", ns, name,
	}
	sess := sh.Command("helm", command...).SetStdin(buf)
	err = sess.Run()
	if err != nil {
		klog.ErrorS(err, "Failed to run command", "command", command)
	}

	return os.WriteFile(path.Join(fullPath, "values.yaml"), buf.Bytes(), filePerm)
}

func writeYaml(obj client.Object, fullPath string) error {
	b, err := yaml.Marshal(obj)
	if err != nil {
		return err
	}
	return os.WriteFile(path.Join(fullPath, obj.GetName()+".yaml"), b, filePerm)
}
