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
	"github.com/spf13/pflag"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/util/i18n"
	"k8s.io/kubectl/pkg/util/templates"
)

type pinKind int

const (
	pinPrimary pinKind = iota
	pinStandby
)

// NewCmdPin builds either pin-primary (break-glass override: this DC stays primary
// and stays writable even with the control plane gone) or pin-standby
// (standby-hold: this DC never promotes).
//
// Both act on a SPOKE cluster: the marker ConfigMap is human-owned and lives on the
// data center it governs, which is exactly why it still works when the hub is
// unreachable. So the ordinary kubeconfig flags must point at that DC's spoke.
func NewCmdPin(f cmdutil.Factory, kind pinKind) *cobra.Command {
	var scopeName, dbName string
	var remove, yes bool
	var cf CoordFlags

	use, short, long, example := pinPrimaryTexts()
	suffix := OverrideCMSuffix
	if kind == pinStandby {
		use, short, long, example = pinStandbyTexts()
		suffix = StandbyHoldCMSuffix
	}

	cmd := &cobra.Command{
		Use:               use,
		Short:             i18n.T(short),
		Long:              templates.LongDesc(long),
		Example:           templates.Examples(example),
		DisableAutoGenTag: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if (scopeName == "") == (dbName == "") {
				return fmt.Errorf("give exactly one of --scope (a primary-DC Lease name) or --db (resolve the scope from a database)")
			}
			ctx := context.Background()
			out := cmd.OutOrStdout()
			lease := scopeName
			if dbName != "" {
				db, dyn, _, err := getDB(ctx, f, dbName)
				if err != nil {
					return fmt.Errorf("resolving --db (note: with --db the kubeconfig must reach the hub, while creating the marker needs the SPOKE; pass --scope when running against a spoke): %w", err)
				}
				s, err := ResolveScopeForDB(ctx, dyn, db)
				if err != nil {
					return err
				}
				lease = s.LeaseName
				fmt.Fprintf(out, "Resolved scope %s from database %s (%s).\n", lease, dbName, s.Source)
			}
			cmName := lease + suffix
			if !remove && !yes {
				if kind == pinPrimary {
					fmt.Fprintf(out, "Would create ConfigMap %s/%s on the CURRENT cluster.\n", cf.LeaseNS, cmName)
					fmt.Fprintf(out, "This is a STANDING BYPASS of split-brain protection: the local leader stays writable regardless of what the failover authority says, and the scope cannot fail over while it exists.\n")
				} else {
					fmt.Fprintf(out, "Would create ConfigMap %s/%s on the CURRENT cluster; this data center will never promote while it exists.\n", cf.LeaseNS, cmName)
				}
				return fmt.Errorf("re-run with --yes to proceed (or --remove --yes to clear an existing pin)")
			}
			if remove && !yes {
				fmt.Fprintf(out, "Would delete ConfigMap %s/%s on the CURRENT cluster.\n", cf.LeaseNS, cmName)
				return fmt.Errorf("re-run with --yes to proceed")
			}
			action, err := markerConfigMap(ctx, f, cf.LeaseNS, cmName, remove)
			if err != nil {
				return err
			}
			fmt.Fprintf(out, "ConfigMap %s/%s %s on the current cluster.\n", cf.LeaseNS, cmName, action)
			if kind == pinPrimary {
				if remove {
					fmt.Fprintf(out, "The break-glass pin is cleared; normal contention resumes and the override-hold annotation drops from the Lease within seconds.\n")
				} else {
					fmt.Fprintf(out, "This data center is now PINNED primary for scope %s:\n", lease)
					fmt.Fprintf(out, "  - its agent mirrors the pin onto the Lease, so no other Member contends;\n")
					fmt.Fprintf(out, "  - its coordinator keeps the local leader writable even if the marker goes stale (a control-plane outage no longer fences it).\n")
					fmt.Fprintf(out, "  IMPORTANT: only honored on the scope's LAST KNOWN HOLDER. On any other DC the agent refuses it and logs why; it never promotes a standby.\n")
					fmt.Fprintf(out, "  IMPORTANT: while pinned, that DC dying means NO failover happens. Remove the pin as soon as the emergency ends:\n")
					fmt.Fprintf(out, "    kubectl dba dc-dr pin-primary --scope %s --remove --yes\n", lease)
				}
			} else {
				if remove {
					fmt.Fprintf(out, "The standby-hold is cleared. NOTE: this DC resumes contending on the NEXT Lease event; if the Lease is idle, a pending handoff may take until the agent's informer resync. Touching the Lease (for example dc-dr handoff --to <dc>) makes it immediate.\n")
				} else {
					fmt.Fprintf(out, "This data center is now HELD as a standby for scope %s: it never contends for the Lease, never promotes, and refuses destructive cross-DC rewinds of its data.\n", lease)
					fmt.Fprintf(out, "  It is ignored while this DC is the ACTIVE one (demoting the active DC without a quiesce is unsafe): move the primary away with a switchover first.\n")
					fmt.Fprintf(out, "  Remove with:  kubectl dba dc-dr pin-standby --scope %s --remove --yes\n", lease)
				}
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&scopeName, "scope", "", "Primary-DC Lease name of the scope (for example primary-dc or primary-dc-orders)")
	cmd.Flags().StringVar(&dbName, "db", "", "Resolve the scope from this database instead (requires the kubeconfig to reach the hub)")
	cmd.Flags().BoolVar(&remove, "remove", false, "Remove the pin instead of creating it")
	cmd.Flags().BoolVar(&yes, "yes", false, "Confirm")
	// Only the marker namespace is meaningful here; the pin never touches the
	// coordination plane itself.
	cmd.Flags().StringVar(&cf.LeaseNS, "coord-namespace", DefaultCoordNamespace, "Namespace on this spoke that holds the marker ConfigMaps")
	cmd.Flags().VisitAll(func(fl *pflag.Flag) {})
	return cmd
}

func pinPrimaryTexts() (use, short, long, example string) {
	return "pin-primary (--scope LEASE | --db DB_NAME)",
		"Pin this data center as primary (break-glass override): no failover, writable through a control-plane outage",
		`Creates the human-owned break-glass override ConfigMap <scope>-override in
			the coordination namespace of the CURRENT cluster, which must be the data
			center you are pinning (its own spoke).

			Two effects, both live-proven: this DC's agent mirrors the pin onto the
			Lease so every other Member defers permanently, and this DC's coordinator
			forces its leader ACTIVE regardless of marker state, so a sustained
			coordination-plane outage no longer fences it read-only after the usual
			marker TTL plus uncertainty hold.

			Use it when the failover authority is unreachable and the surviving
			primary must keep accepting writes, or as a deliberate "never fail this
			scope over" policy. While it stands there is no split-brain protection for
			the scope, and nothing takes over if this DC dies.

			KUBECONFIG: the SPOKE of the data center being pinned.`,
		`# Keep dc-b primary for the global scope, come what may (run against dc-b)
			kubectl dba dc-dr pin-primary --scope primary-dc --yes --kubeconfig ~/.kube/dc-b.yaml

			# Clear it once the emergency is over
			kubectl dba dc-dr pin-primary --scope primary-dc --remove --yes --kubeconfig ~/.kube/dc-b.yaml`
}

func pinStandbyTexts() (use, short, long, example string) {
	return "pin-standby (--scope LEASE | --db DB_NAME)",
		"Pin this data center as a standby (standby-hold): it never promotes",
		`Creates the human-owned standby-hold ConfigMap <scope>-standby-hold in the
			coordination namespace of the CURRENT cluster, which must be the data
			center you are holding down.

			While it exists that DC never contends for the scope's primary-DC Lease,
			never promotes (it refuses even an explicit handoff naming it), and its
			coordinator refuses destructive cross-DC rewinds or re-seeds of the data
			it holds. It fails CLOSED: if the ConfigMap cannot be read the hold is
			assumed, so a flaky apiserver never silently drops the protection.

			It is deliberately ignored on the data center that is currently ACTIVE,
			because demoting the active DC without a quiesce is unsafe; move the
			primary away with a planned switchover first.

			KUBECONFIG: the SPOKE of the data center being held.`,
		`# Never let dc-a take the primary role for this scope (run against dc-a)
			kubectl dba dc-dr pin-standby --scope primary-dc --yes --kubeconfig ~/.kube/dc-a.yaml

			# Release it
			kubectl dba dc-dr pin-standby --scope primary-dc --remove --yes --kubeconfig ~/.kube/dc-a.yaml`
}
