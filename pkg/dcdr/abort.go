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

package dcdr

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/util/i18n"
	"k8s.io/kubectl/pkg/util/templates"
)

// NewCmdAbort aborts an in-flight planned switchover. Acts on the HUB cluster.
func NewCmdAbort(f cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "abort DB_NAME",
		Short: i18n.T("Abort an in-flight planned switchover and restore writes to the current active DC"),
		Long: templates.LongDesc(`
			Sets the dr.kubedb.com/switchover-abort annotation. Its PRESENCE aborts:
			the hub clears the quiesce so the original active data center resumes
			accepting writes, and removes every switchover annotation once done.

			Do NOT abort by deleting the switchover-to annotation: in a scope shared
			by several databases the hub re-propagates it to every sibling each pass,
			so a bare removal silently reappears. This explicit abort signal is
			propagated and honored scope-wide. A switchover that cannot complete also
			auto-aborts on its own after the switchover timeout (default 10m,
			dr.kubedb.com/switchover-timeout to override).

			KUBECONFIG: the hub cluster.`),
		Example: templates.Examples(`
			kubectl dba dc-dr abort pg-dcdr -n demo`),
		DisableAutoGenTag: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("exactly one database name is required")
			}
			ctx := context.Background()
			db, dyn, ns, err := getDB(ctx, f, args[0])
			if err != nil {
				return err
			}
			if db.GetAnnotations()[AnnSwitchoverTo] == "" && db.GetAnnotations()[AnnQuiesceActive] == "" {
				cmd.Printf("No switchover appears to be in flight on %s/%s (no %s or %s annotation); setting the abort anyway is harmless.\n", ns, args[0], AnnSwitchoverTo, AnnQuiesceActive)
			}
			v := "true"
			if err := annotateDB(ctx, dyn, ns, args[0], AnnSwitchoverAbort, &v); err != nil {
				return err
			}
			cmd.Printf("Abort requested for %s/%s. The hub restores writes to the current active DC and clears the switchover annotations; verify with:\n", ns, args[0])
			cmd.Printf("  kubectl dba dc-dr status %s -n %s\n", args[0], ns)
			return nil
		},
	}
	return cmd
}
