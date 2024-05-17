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
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	appcat "kmodules.xyz/custom-resources/apis/appcatalog/v1alpha1"
	"os"
	"path"
	stashapi "stash.appscode.dev/apimachinery/apis/stash/v1beta1"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/yaml"
)

func writeYaml(obj client.Object, fullPath string) error {
	b, err := yaml.Marshal(obj)
	if err != nil {
		return err
	}
	return os.WriteFile(path.Join(fullPath, obj.GetName()+".yaml"), b, filePerm)
}

func isStashCRDAvailable(kc client.Client, kind string) bool {
	_, err := kc.RESTMapper().RESTMapping(schema.GroupKind{
		Group: stashapi.SchemeGroupVersion.Group,
		Kind:  kind,
	})
	return err == nil
}

func isBackupTargetMatched(ref stashapi.TargetRef, meta metav1.ObjectMeta) bool {
	fmt.Println(appcat.SchemeGroupVersion.String(), "-----", ref)
	eq := ref.Name == meta.Name && ref.APIVersion == appcat.SchemeGroupVersion.String() && ref.Kind == appcat.ResourceKindApp
	if ref.Namespace != "" {
		return eq && ref.Namespace == meta.Namespace
	}
	return eq
}
