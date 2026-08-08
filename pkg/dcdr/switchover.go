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
	"slices"

	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/util/i18n"
	"k8s.io/kubectl/pkg/util/templates"
)

const mergePatchType = types.MergePatchType

// NewCmdSwitchover triggers a planned, zero-RPO switchover by annotating the
// database. Acts on the HUB cluster (ordinary kubeconfig flags).
func NewCmdSwitchover(f cmdutil.Factory) *cobra.Command {
	var to string
	cmd := &cobra.Command{
		Use:   "switchover DB_NAME --to DC",
		Short: i18n.T("Trigger a planned zero-RPO switchover of a distributed database to another data center"),
		Long: templates.LongDesc(`
			Sets the dr.kubedb.com/switchover-to annotation on the Postgres. The hub
			operator then quiesces the active primary (write-locked), waits for the
			target data center to catch up to the frozen LSN, hands off the primary-DC
			Lease, and clears the annotation. Requires the active primary to be up and
			accepting connections: the safety gates measure by dialing it and fail
			closed, so a dead primary cannot be switched away from (use the failover
			path instead: dc-dr handoff, and dc-dr accept-data-loss if the RPO budget
			holds it).

			KUBECONFIG: the hub cluster (where the Postgres CR lives).`),
		Example: templates.Examples(`
			# Move demo/pg-dcdr to data center dc-a, with zero data loss
			kubectl dba dc-dr switchover pg-dcdr -n demo --to dc-a

			# Watch the progress (one-shot, run repeatedly)
			kubectl dba dc-dr status pg-dcdr -n demo`),
		DisableAutoGenTag: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("exactly one database name is required")
			}
			if to == "" {
				return fmt.Errorf("--to is required (the target data center)")
			}
			ctx := context.Background()
			db, dyn, ns, err := getDB(ctx, f, args[0])
			if err != nil {
				return err
			}
			requireDistributed(printfErr(cmd), db)

			scope, err := ResolveScopeForDB(ctx, dyn, db)
			if err != nil {
				return err
			}
			if len(scope.MemberDCs) > 0 && !slices.Contains(scope.MemberDCs, to) {
				return fmt.Errorf("%q is not a Member data center of this database (members: %v); an Arbiter or Witness DC can never become primary", to, scope.MemberDCs)
			}
			if active, _, _ := unstructured.NestedString(db.Object, "status", "disasterRecovery", "activeDC"); active == to {
				cmd.Printf("No-op: %s is already the active data center of %s/%s.\n", to, ns, args[0])
				return nil
			}
			if err := annotateDB(ctx, dyn, ns, args[0], AnnSwitchoverTo, &to); err != nil {
				return err
			}
			cmd.Printf("Switchover of %s/%s to %q requested (scope %s, from %s).\n", ns, args[0], to, scope.LeaseName, scope.Source)
			cmd.Printf("The operator will quiesce, wait for catch-up, and hand off; zero committed rows are lost.\n")
			cmd.Printf("Monitor:  kubectl dba dc-dr status %s -n %s\n", args[0], ns)
			cmd.Printf("Abort:    kubectl dba dc-dr abort %s -n %s\n", args[0], ns)
			if scope.LeaseName == GlobalPrimaryLease {
				cmd.Printf("NOTE: this database follows the GLOBAL scope; every database in that scope switches with it.\n")
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&to, "to", "", "Target data center (must be a Member DC of the database's PlacementPolicy)")
	return cmd
}

func printfErr(cmd *cobra.Command) func(string, ...any) {
	return func(format string, a ...any) { _, _ = fmt.Fprintf(cmd.ErrOrStderr(), format, a...) }
}
