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
	"context"
	"fmt"
	"log"
	"os"
	"path"

	gitops "kubedb.dev/apimachinery/apis/gitops/v1alpha1"
	opsapi "kubedb.dev/apimachinery/apis/ops/v1alpha1"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type GitOpsStatus struct {
	GitOps gitops.GitOpsStatus `json:"gitops,omitempty" yaml:"gitops,omitempty"`
}

type GitOps struct {
	Status GitOpsStatus `json:"status,omitempty" yaml:"status,omitempty"`
}

type Ops struct {
	Status *opsapi.OpsRequestStatus `json:"status,omitempty" yaml:"status,omitempty"`
}

type dbInfo struct {
	kind      string
	name      string
	namespace string
}

type gitOpsOpts struct {
	kc client.Client
	db dbInfo

	dir     string
	summary []string
}

func (g *gitOpsOpts) collectGitOpsYamls() error {
	err := g.collectGitOpsDatabase()
	if err != nil {
		return err
	}

	fmt.Println("Summary:")
	for _, line := range g.summary {
		fmt.Println("- ", line)
	}
	fmt.Println("--------------- Done ---------------")
	return nil
}

func newGitOpsOpts(kc client.Client, name, namespace, kind, dir string) (*gitOpsOpts, error) {
	err := os.MkdirAll(dir, dirPerm)
	if err != nil {
		return nil, err
	}
	opts := &gitOpsOpts{
		kc: kc,
		db: dbInfo{
			kind:      kind,
			name:      name,
			namespace: namespace,
		},
		dir:     dir,
		summary: make([]string, 0),
	}
	return opts, nil
}

func (g *gitOpsOpts) collectGitOpsDatabase() error {
	var uns unstructured.Unstructured
	uns.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   gitops.SchemeGroupVersion.Group,
		Version: gitops.SchemeGroupVersion.Version,
		Kind:    g.db.kind,
	})
	err := g.kc.Get(context.Background(), types.NamespacedName{
		Namespace: g.db.namespace,
		Name:      g.db.name,
	}, &uns)
	if err != nil {
		log.Fatalf("failed to get gitops database obj: %v", err)
		return err
	}

	var gitOpsObj GitOps
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(uns.Object, &gitOpsObj)
	if err != nil {
		log.Fatalf("failed to convert unstructured to gitops obj: %v", err)
		return err
	}

	statuses := []string{
		string(gitops.ChangeRequestStatusInCurrent), string(gitops.ChangeRequestStatusPending),
		string(gitops.ChangeRequestStatusInProgress), string(gitops.ChangeRequestStatusFailed),
	}
	statusIdx := 0
	for _, info := range gitOpsObj.Status.GitOps.GitOpsInfo {
		for i := range statuses {
			if string(info.ChangeRequestStatus) == statuses[i] {
				if i > statusIdx {
					statusIdx = i
				}
			}
		}
	}

	g.summary = append(g.summary, fmt.Sprintf("GitOps Database Status for: %s/%s is %s", g.db.namespace, g.db.name, statuses[statusIdx]))

	if err := g.collectOpsRequests(gitOpsObj.Status); err != nil {
		return err
	}

	return writeYaml(&uns, g.dir)
}

func (g *gitOpsOpts) collectOpsRequests(gitOpsStatus GitOpsStatus) error {
	opsYamlDir := path.Join(g.dir, "ops")
	err := os.MkdirAll(opsYamlDir, dirPerm)
	if err != nil {
		return err
	}
	for _, info := range gitOpsStatus.GitOps.GitOpsInfo {
		for _, op := range info.Operations {
			var uns unstructured.Unstructured
			uns.SetGroupVersionKind(schema.GroupVersionKind{
				Group:   opsapi.SchemeGroupVersion.Group,
				Version: opsapi.SchemeGroupVersion.Version,
				Kind:    g.db.kind + "OpsRequest",
			})
			err := g.kc.Get(context.Background(), types.NamespacedName{
				Namespace: g.db.namespace,
				Name:      op.Name,
			}, &uns)
			if err != nil {
				log.Fatalf("failed to get opsrequest: %v", err)
				return err
			}
			err = writeYaml(&uns, opsYamlDir)
			if err != nil {
				return err
			}
			var ops Ops
			err = runtime.DefaultUnstructuredConverter.FromUnstructured(uns.Object, &ops)
			if err != nil {
				log.Fatalf("failed to convert unstructured to opsrequest obj: %v", err)
				return err
			}
			if ops.Status.Phase == opsapi.OpsRequestPhaseFailed {
				for _, cond := range ops.Status.Conditions {
					if cond.Reason == opsapi.Failed {
						g.summary = append(g.summary, fmt.Sprintf("RequestName %s: %s", op.Name, cond.Message))
					}
				}
			}
		}
	}

	return nil
}
