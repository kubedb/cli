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
	"os"
	"path"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/yaml"
)

type Resource interface {
	OffshootSelectors() map[string]string
}

type ResourceExtra interface {
	OffshootSelectors(extraSelectors ...map[string]string) map[string]string
}

func getOffshootSelectorsIfAny(obj client.Object) (map[string]string, bool) {
	if r, ok := obj.(Resource); ok {
		return r.OffshootSelectors(), true
	}
	if r, ok := obj.(ResourceExtra); ok {
		return r.OffshootSelectors(), true
	}
	return nil, false
}

func (opts *dbOpts) run(db client.Object) error {
	opts.db.OwnerReferences = db.GetOwnerReferences()
	if err := writeYaml(db, getDir(db.GetName())); err != nil {
		return err
	}
	selectors, ok := getOffshootSelectorsIfAny(db)
	if !ok {
		return fmt.Errorf("%T does not implement OffshootSelectors", db)
	}
	opts.selectors = selectors
	klog.Infof("db selectors: %v;\nDebug info has been generated in '%v' folder", opts.selectors, db.GetName())
	return opts.collectALl()
}

func IsOwnByGitOpsObject(owners []metav1.OwnerReference, name, kind string) bool {
	for _, owner := range owners {
		if owner.Kind == kind && owner.Name == name && owner.APIVersion == "gitops.kubedb.com/v1alpha1" {
			return true
		}
	}
	return false
}

func writeYaml(obj client.Object, fullPath string) error {
	b, err := yaml.Marshal(obj)
	if err != nil {
		return err
	}
	return os.WriteFile(path.Join(fullPath, obj.GetName()+".yaml"), b, filePerm)
}

func getDir(dbName string) string {
	pwd, _ := os.Getwd()
	dir := path.Join(pwd, dbName)
	return dir
}

type OpsRequest struct {
	Spec OpsRequestSpec `json:"spec,omitempty" yaml:"spec,omitempty"`
}

type OpsRequestSpec struct {
	DatabaseRef DatabaseRef `json:"databaseRef,omitempty" yaml:"databaseRef,omitempty"`
}

type Autoscaler struct {
	Spec AutoscalerSpec `json:"spec,omitempty" yaml:"spec,omitempty"`
}

type AutoscalerSpec struct {
	DatabaseRef DatabaseRef `json:"databaseRef,omitempty" yaml:"databaseRef,omitempty"`
}
type DatabaseRef struct {
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
}
