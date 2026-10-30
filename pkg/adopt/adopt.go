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

// Package adopt implements `kubectl dba adopt`, which brings a database that
// already runs in Kubernetes (for example from a Helm chart) under KubeDB in
// place: the data stays where it is and the workloads change hands one role at
// a time.
package adopt

import (
	"github.com/spf13/cobra"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/util/i18n"
	"k8s.io/kubectl/pkg/util/templates"
)

// NewCmdAdopt returns the `adopt` command group.
func NewCmdAdopt(f cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "adopt",
		Short: i18n.T("Bring a running database under KubeDB in place, without moving its data"),
		Long: templates.LongDesc(`
			Take over a database that already runs in Kubernetes, for example one
			installed by a Helm chart. The data stays in its current storage, the
			credentials stay the same, and clients keep their address; only the
			workloads change hands, one role at a time.`),
		Run:                   func(cmd *cobra.Command, args []string) {},
		DisableFlagsInUseLine: true,
		DisableAutoGenTag:     true,
	}
	cmd.AddCommand(NewCmdMilvus(f))
	return cmd
}
