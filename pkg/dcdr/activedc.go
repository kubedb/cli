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
	"time"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/util/i18n"
	"k8s.io/kubectl/pkg/util/templates"
)

// NewCmdActiveDC answers "who is primary" from the authoritative source, the Lease
// on the coordination plane, resolving the scope from a database name when given.
func NewCmdActiveDC(f cmdutil.Factory) *cobra.Command {
	var cf CoordFlags
	var leaseName string
	var quiet bool
	cmd := &cobra.Command{
		Use:   "active-dc [DB_NAME] [--lease NAME]",
		Short: i18n.T("Print the data center that currently holds the primary role"),
		Long: templates.LongDesc(`
			Reads the primary-DC Lease from the coordination control plane, the
			authority for which data center is active.

			Given a database name, its failover scope is resolved first (the
			PlacementPolicy's failoverPolicy trigger, exactly as the operator
			resolves it) and the matching Lease is read. Given --lease, that Lease is
			read directly, which also works for a scope whose database is gone.

			KUBECONFIG: the hub cluster (to read the Postgres and its
			PlacementPolicy, and by default to read the coordination kubeconfig
			Secret). The coordination plane itself is reached with the --coord-*
			flags.`),
		Example: templates.Examples(`
			# By database
			kubectl dba dc-dr active-dc pg-dcdr -n demo

			# By Lease name, with an explicit coordination kubeconfig file
			kubectl dba dc-dr active-dc --lease primary-dc --coord-kubeconfig /tmp/coord.yaml

			# Scriptable: just the DC name
			kubectl dba dc-dr active-dc pg-dcdr -n demo -q`),
		DisableAutoGenTag: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if (len(args) == 0) == (leaseName == "") {
				return fmt.Errorf("give exactly one of: a database name, or --lease")
			}
			ctx := context.Background()
			out := cmd.OutOrStdout()
			scope := &Scope{LeaseName: leaseName, Source: "--lease flag"}
			var ns string
			if len(args) == 1 {
				db, dyn, dbNS, err := getDB(ctx, f, args[0])
				if err != nil {
					return err
				}
				ns = dbNS
				scope, err = ResolveScopeForDB(ctx, dyn, db)
				if err != nil {
					return err
				}
			}
			coord, err := cf.CoordClient(ctx, f)
			if err != nil {
				return err
			}
			lease, err := coord.CoordinationV1().Leases(cf.LeaseNS).Get(ctx, scope.LeaseName, metav1.GetOptions{})
			if err != nil {
				return fmt.Errorf("failed to read Lease %s/%s on the coordination plane: %w", cf.LeaseNS, scope.LeaseName, err)
			}
			holder := ""
			if lease.Spec.HolderIdentity != nil {
				holder = *lease.Spec.HolderIdentity
			}
			if quiet {
				_, _ = fmt.Fprintln(out, holder)
				return nil
			}
			_, _ = fmt.Fprintf(out, "Active DC:  %s\n", orNone(holder))
			_, _ = fmt.Fprintf(out, "Lease:      %s/%s  (scope from %s)\n", cf.LeaseNS, scope.LeaseName, scope.Source)
			if lease.Spec.RenewTime != nil {
				age := time.Since(lease.Spec.RenewTime.Time).Round(time.Second)
				dur := int32(0)
				if lease.Spec.LeaseDurationSeconds != nil {
					dur = *lease.Spec.LeaseDurationSeconds
				}
				_, _ = fmt.Fprintf(out, "Renewed:    %s ago (lease duration %ds)\n", age, dur)
				if dur > 0 && age > time.Duration(dur)*time.Second {
					_, _ = fmt.Fprintf(out, "  WARNING: the Lease is EXPIRED. Its holder stopped renewing, so a healthy Member DC may acquire it at any moment.\n")
				}
			}
			if lease.Spec.LeaseTransitions != nil {
				_, _ = fmt.Fprintf(out, "Transitions:%d\n", *lease.Spec.LeaseTransitions)
			}
			if v := lease.Annotations[AnnLeaseMemberDCs]; v != "" {
				_, _ = fmt.Fprintf(out, "Members:    %s\n", v)
			}
			if v := lease.Annotations[AnnLeaseHandoffTo]; v != "" {
				_, _ = fmt.Fprintf(out, "Handoff to: %s (a coordinated handoff is in flight)\n", v)
			}
			if v := lease.Annotations["dr.open-cluster-management.io/override-hold"]; v != "" {
				_, _ = fmt.Fprintf(out, "PINNED:     break-glass override holds this scope on %s; it cannot fail over until the override ConfigMap is removed from that DC's spoke.\n", v)
			}
			if len(args) == 1 {
				// Cross-check the CR's own view, which lags the Lease by a reconcile.
				db, _, _, err := getDB(ctx, f, args[0])
				if err == nil {
					crActive, _, _ := unstructured.NestedString(db.Object, "status", "disasterRecovery", "activeDC")
					if crActive != "" && crActive != holder {
						_, _ = fmt.Fprintf(out, "\nNOTE: the database %s/%s still reports activeDC=%s in its status; the Lease is the authority and the status trails it by a reconcile.\n", ns, args[0], crActive)
					}
				}
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&leaseName, "lease", "", "Read this Lease directly instead of resolving a database's scope")
	cmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Print only the active DC name")
	AddCoordFlags(cmd, &cf)
	return cmd
}
