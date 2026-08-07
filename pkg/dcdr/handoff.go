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
	"strings"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/util/i18n"
	"k8s.io/kubectl/pkg/util/templates"
)

// NewCmdHandoff moves the failover authority by annotating the Lease itself. Acts
// on the COORDINATION plane (and on the hub only to resolve a database's scope).
func NewCmdHandoff(f cmdutil.Factory) *cobra.Command {
	var cf CoordFlags
	var leaseName, to string
	var yes bool
	cmd := &cobra.Command{
		Use:   "handoff (DB_NAME | --lease NAME) --to DC",
		Short: i18n.T("Move the failover authority for a scope by handing off its primary-DC Lease"),
		Long: templates.LongDesc(`
			Writes dr.open-cluster-management.io/handoff-to on the scope's primary-DC
			Lease. The holding data center's agent releases the Lease once, the target
			acquires it within a retry tick, and the annotation clears itself.

			This is the scope-local FAILOVER lever, and the correct tool when the
			active data center's database is down but its data center is alive: no
			quiesce and no catch-up wait happen, so loss is bounded by the RPO budget
			rather than zero. For a healthy primary prefer "dc-dr switchover", which
			is zero-RPO.

			It moves EVERY database sharing the scope. Do NOT stop a DC's agent to
			force a failover instead: one agent serves every scope its DC holds, so
			that expires all of them together.

			KUBECONFIG: the hub cluster (to resolve a database's scope and read the
			coordination kubeconfig Secret). The Lease is written on the coordination
			plane via the --coord-* flags.`),
		Example: templates.Examples(`
			# Fail a database's scope over to dc-b
			kubectl dba dc-dr handoff pg-dcdr -n demo --to dc-b --yes

			# Move a scope by Lease name (works with no database left)
			kubectl dba dc-dr handoff --lease primary-dc-orders --to dc-a --yes`),
		DisableAutoGenTag: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if (len(args) == 0) == (leaseName == "") {
				return fmt.Errorf("give exactly one of: a database name, or --lease")
			}
			if to == "" {
				return fmt.Errorf("--to is required (the target data center)")
			}
			ctx := context.Background()
			out := cmd.OutOrStdout()
			scope := &Scope{LeaseName: leaseName, Source: "--lease flag"}
			if len(args) == 1 {
				db, dyn, _, err := getDB(ctx, f, args[0])
				if err != nil {
					return err
				}
				scope, err = ResolveScopeForDB(ctx, dyn, db)
				if err != nil {
					return err
				}
				if len(scope.MemberDCs) > 0 && !slices.Contains(scope.MemberDCs, to) {
					return fmt.Errorf("%q is not a Member data center of this database (members: %v)", to, scope.MemberDCs)
				}
			}
			coord, err := cf.CoordClient(ctx, f)
			if err != nil {
				return err
			}
			lease, err := coord.CoordinationV1().Leases(cf.LeaseNS).Get(ctx, scope.LeaseName, metav1.GetOptions{})
			if err != nil {
				return fmt.Errorf("failed to read Lease %s/%s: %w", cf.LeaseNS, scope.LeaseName, err)
			}
			holder := ""
			if lease.Spec.HolderIdentity != nil {
				holder = *lease.Spec.HolderIdentity
			}
			if holder == to {
				_, _ = fmt.Fprintf(out, "No-op: %s already holds %s.\n", to, scope.LeaseName)
				return nil
			}
			if members := lease.Annotations[AnnLeaseMemberDCs]; members != "" {
				if !slices.Contains(strings.Split(members, ","), to) {
					return fmt.Errorf("%q is not listed in the Lease's Member DCs (%s); only a Member can hold the primary role", to, members)
				}
			}
			if pin := lease.Annotations["dr.open-cluster-management.io/override-hold"]; pin != "" && pin != to {
				return fmt.Errorf("scope %s is PINNED to %q by a break-glass override; remove that DC's override ConfigMap first (kubectl dba dc-dr pin-primary --remove), or the handoff cannot complete", scope.LeaseName, pin)
			}
			if !yes {
				_, _ = fmt.Fprintf(out, "Would move %s from %s to %s. Every database in this scope fails over together, without a quiesce (loss bounded by the RPO budget, not zero).\n", scope.LeaseName, orNone(holder), to)
				return fmt.Errorf("re-run with --yes to proceed")
			}
			// Re-requesting the SAME target must still fire an agent-visible event.
			// A merge patch that writes an identical value does not bump the
			// resourceVersion, so the informers see nothing and the handoff never
			// starts: exactly what happens when a target that was standby-held is
			// released and the pending handoff has to be re-driven. Clear first, then
			// set, so there is always a real transition to observe.
			if lease.Annotations[AnnLeaseHandoffTo] == to {
				clear := fmt.Sprintf(`{"metadata":{"annotations":{%q:null}}}`, AnnLeaseHandoffTo)
				if _, err := coord.CoordinationV1().Leases(cf.LeaseNS).Patch(ctx, scope.LeaseName, types.MergePatchType, []byte(clear), metav1.PatchOptions{}); err != nil {
					return fmt.Errorf("failed to clear the stale handoff annotation on %s/%s: %w", cf.LeaseNS, scope.LeaseName, err)
				}
			}
			patch := fmt.Sprintf(`{"metadata":{"annotations":{%q:%q}}}`, AnnLeaseHandoffTo, to)
			if _, err := coord.CoordinationV1().Leases(cf.LeaseNS).Patch(ctx, scope.LeaseName, types.MergePatchType, []byte(patch), metav1.PatchOptions{}); err != nil {
				return fmt.Errorf("failed to annotate Lease %s/%s: %w", cf.LeaseNS, scope.LeaseName, err)
			}
			_, _ = fmt.Fprintf(out, "Handoff of %s requested: %s -> %s.\n", scope.LeaseName, orNone(holder), to)
			_, _ = fmt.Fprintf(out, "The holder releases within seconds and %s acquires on its next retry tick; the annotation clears itself.\n", to)
			_, _ = fmt.Fprintf(out, "Verify:  kubectl dba dc-dr active-dc --lease %s\n", scope.LeaseName)
			_, _ = fmt.Fprintf(out, "NOTE: if the target is held by a standby-hold ConfigMap it will refuse; clear it with dc-dr pin-standby --remove.\n")
			return nil
		},
	}
	cmd.Flags().StringVar(&leaseName, "lease", "", "Act on this Lease directly instead of resolving a database's scope")
	cmd.Flags().StringVar(&to, "to", "", "Target data center")
	cmd.Flags().BoolVar(&yes, "yes", false, "Confirm the handoff (it moves every database in the scope)")
	AddCoordFlags(cmd, &cf)
	return cmd
}
