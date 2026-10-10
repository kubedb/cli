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

package remote_replica

import (
	"testing"

	"kubedb.dev/apimachinery/apis/kubedb"

	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func testPod(name, role string, phase core.PodPhase, containers ...string) core.Pod {
	p := core.Pod{
		ObjectMeta: metav1.ObjectMeta{Name: name, Labels: map[string]string{kubedb.LabelRole: role}},
		Status:     core.PodStatus{Phase: phase},
	}
	for _, c := range containers {
		p.Status.ContainerStatuses = append(p.Status.ContainerStatuses, core.ContainerStatus{
			Name:  c,
			State: core.ContainerState{Running: &core.ContainerStateRunning{}},
		})
	}
	return p
}

func TestSelectPod(t *testing.T) {
	primary := testPod("pg-0", "primary", core.PodRunning, "postgres", "pg-coordinator")
	standby := testPod("pg-1", "standby", core.PodRunning, "postgres", "pg-coordinator")
	// The arbiter carries the same offshoot labels but runs no postgres
	// container, so a catalog query cannot be exec-ed into it.
	arbiter := testPod("pg-arbiter-0", "arbiter", core.PodRunning, "pg-coordinator")
	pending := testPod("pg-2", "primary", core.PodPending, "postgres")

	cases := []struct {
		name        string
		pods        []core.Pod
		primaryOnly bool
		want        string
		wantErr     bool
	}{
		{name: "prefers the primary", pods: []core.Pod{standby, primary}, want: "pg-0"},
		{name: "prefers the primary for DDL too", pods: []core.Pod{standby, primary}, primaryOnly: true, want: "pg-0"},
		// A remote replica acting as the source of a chained replica has no
		// primary at all; verification must still find a pod to query.
		{name: "falls back to a standby", pods: []core.Pod{arbiter, standby}, want: "pg-1"},
		{name: "refuses DDL without a primary", pods: []core.Pod{standby}, primaryOnly: true, wantErr: true},
		{name: "skips pods without the container", pods: []core.Pod{arbiter}, wantErr: true},
		// An empty result used to be swallowed and reported as success, which
		// produced a config referencing a user that was never created.
		{name: "reports a pod that is not running", pods: []core.Pod{pending}, wantErr: true},
		{name: "reports an empty list", pods: nil, wantErr: true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := selectPod(tc.pods, PostgresContainerName, tc.primaryOnly)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected an error, got pod %q", got.Name)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.Name != tc.want {
				t.Errorf("picked %q, want %q", got.Name, tc.want)
			}
		})
	}
}
