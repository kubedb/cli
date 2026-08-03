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
	"strings"
	"time"

	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/util/i18n"
	"k8s.io/kubectl/pkg/util/templates"
)

// stepState is the rendered state of one switchover step.
type stepState int

const (
	stepDone stepState = iota
	stepCurrent
	stepPending
	stepSkipped
)

func (s stepState) mark() string {
	switch s {
	case stepDone:
		return "[done]   "
	case stepCurrent:
		return "[NOW]    "
	case stepSkipped:
		return "[n/a]    "
	default:
		return "[pending]"
	}
}

// NewCmdStatus renders a one-shot picture of a database's DC-DR state and, when a
// switchover is in flight, which of its steps are complete and what happens next.
// Deliberately NOT a watch: it prints once and exits, so the user re-runs it.
func NewCmdStatus(f cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status DB_NAME",
		Short: i18n.T("Show DC-DR state and switchover progress once (re-run to see further progress)"),
		Long: templates.LongDesc(`
			Prints the database's failover scope, per data center state, protection
			verdict, and, when a planned switchover is in flight, its step-by-step
			progress: what has completed, what is happening NOW, and what remains.

			One-shot by design: it does not follow. Run it again to see the next
			state, which keeps its output readable in tickets and transcripts.

			KUBECONFIG: the hub cluster.`),
		Example: templates.Examples(`
			kubectl dba dc-dr status pg-dcdr -n demo`),
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
			scope, err := ResolveScopeForDB(ctx, dyn, db)
			if err != nil {
				return err
			}
			out := cmd.OutOrStdout()
			ann := db.GetAnnotations()
			dbPhase, _, _ := unstructured.NestedString(db.Object, "status", "phase")
			activeDC, _, _ := unstructured.NestedString(db.Object, "status", "disasterRecovery", "activeDC")
			drPhase, _, _ := unstructured.NestedString(db.Object, "status", "disasterRecovery", "phase")
			protected, protectedSet, _ := unstructured.NestedBool(db.Object, "status", "disasterRecovery", "protected")
			protMsg, _, _ := unstructured.NestedString(db.Object, "status", "disasterRecovery", "protectionMessage")

			fmt.Fprintf(out, "Database:      %s/%s (%s)\n", ns, args[0], dbPhase)
			fmt.Fprintf(out, "Failover scope: %s\n                (%s)\n", scope.LeaseName, scope.Source)
			if len(scope.MemberDCs) > 0 {
				fmt.Fprintf(out, "Member DCs:    %s\n", strings.Join(scope.MemberDCs, ", "))
			}
			fmt.Fprintf(out, "Active DC:     %s   DR phase: %s\n", orNone(activeDC), orNone(drPhase))
			if protectedSet {
				fmt.Fprintf(out, "Protected:     %v", protected)
				if protMsg != "" {
					fmt.Fprintf(out, "  (%s)", protMsg)
				}
				fmt.Fprintln(out)
			}

			// Per-DC table.
			dcs, found, _ := unstructured.NestedSlice(db.Object, "status", "disasterRecovery", "dataCenters")
			if found && len(dcs) > 0 {
				fmt.Fprintf(out, "\nData centers:\n")
				fmt.Fprintf(out, "  %-10s %-9s %-9s %-9s %-12s %s\n", "NAME", "ROLE", "WRITABLE", "HEALTHY", "LAG(BYTES)", "STREAMER")
				for _, d := range dcs {
					dm, ok := d.(map[string]any)
					if !ok {
						continue
					}
					name, _ := dm["clusterName"].(string)
					role, _ := dm["role"].(string)
					writable := boolCell(dm, "writable")
					healthy := boolCell(dm, "healthy")
					lag := "-"
					if v, ok := dm["lagBytes"]; ok {
						lag = fmt.Sprintf("%v", v)
					}
					streamer, _ := dm["crossDCStreamer"].(string)
					fmt.Fprintf(out, "  %-10s %-9s %-9s %-9s %-12s %s\n", name, orNone(role), writable, healthy, lag, orNone(streamer))
				}
			}

			// Switchover progress, only when one is in flight or was just requested.
			target := ann[AnnSwitchoverTo]
			quiesced := ann[AnnQuiesceActive] == "true"
			aborting := ann[AnnSwitchoverAbort] != ""
			started := ann[AnnSwitchoverStart]
			if target != "" || quiesced || aborting {
				fmt.Fprintf(out, "\nPlanned switchover")
				if target != "" {
					fmt.Fprintf(out, " to %q", target)
				}
				if started != "" {
					if t, perr := time.Parse(time.RFC3339, started); perr == nil {
						fmt.Fprintf(out, ", started %s ago", time.Since(t).Round(time.Second))
					}
				}
				fmt.Fprintln(out, ":")
				if aborting {
					fmt.Fprintf(out, "  ABORT requested (%s is set). The hub clears the quiesce and restores writes to %s.\n", AnnSwitchoverAbort, orNone(activeDC))
					fmt.Fprintf(out, "  Re-run this command until the switchover annotations are gone and DR phase is Steady.\n")
					return nil
				}
				renderSwitchoverSteps(out, target, activeDC, quiesced, dcs)
				fmt.Fprintf(out, "\n  Abort:  kubectl dba dc-dr abort %s -n %s\n", args[0], ns)
				fmt.Fprintf(out, "  Re-run this command to see the next step; it does not follow.\n")
				return nil
			}

			fmt.Fprintf(out, "\nNo switchover in flight.\n")
			if drPhase == "FailingOver" {
				fmt.Fprintf(out, "DR phase is FailingOver with no switchover annotation, so this is an UNPLANNED failover.\n")
				fmt.Fprintf(out, "  If it is not completing: kubectl dba dc-dr debug failover %s -n %s\n", args[0], ns)
			}
			if protectedSet && !protected {
				fmt.Fprintf(out, "Protection is NOT confirmed. If a promotion is held by the RPO budget:\n")
				fmt.Fprintf(out, "  kubectl dba dc-dr accept-data-loss %s -n %s --yes\n", args[0], ns)
			}
			fmt.Fprintf(out, "Trigger one:  kubectl dba dc-dr switchover %s -n %s --to <dc>\n", args[0], ns)
			return nil
		},
	}
	return cmd
}

// renderSwitchoverSteps prints the operator's switchover sequence with each step's
// state derived from the SAME status fields the operator's own gates read, so the
// output cannot drift from the real decision:
//
//	1 target validated (healthy, lag known and within the switchover budget)
//	2 quiesce requested   (annotation set)
//	3 quiesce in effect    (active DC's writable flipped false)
//	4 target caught up     (target lagBytes <= 8 KiB)
//	5 Lease handed off     (activeDC == target)
//	6 old DC demoted, annotations cleared, phase Steady
func renderSwitchoverSteps(out interface{ Write([]byte) (int, error) }, target, activeDC string, quiesced bool, dcs []any) {
	var targetLag any
	targetHealthy, targetKnown := false, false
	activeWritable, activeWritableKnown := true, false
	for _, d := range dcs {
		dm, ok := d.(map[string]any)
		if !ok {
			continue
		}
		name, _ := dm["clusterName"].(string)
		if name == target {
			targetKnown = true
			targetHealthy, _ = dm["healthy"].(bool)
			targetLag = dm["lagBytes"]
		}
		if name == activeDC {
			if w, ok := dm["writable"].(bool); ok {
				activeWritable, activeWritableKnown = w, true
			}
		}
	}
	handedOff := target != "" && activeDC == target

	step1 := stepPending
	switch {
	case !targetKnown:
		step1 = stepCurrent
	case targetHealthy && targetLag != nil:
		step1 = stepDone
	default:
		step1 = stepCurrent
	}
	step2 := stepPending
	if quiesced {
		step2 = stepDone
	}
	step3 := stepPending
	switch {
	case handedOff:
		step3 = stepDone
	case quiesced && activeWritableKnown && !activeWritable:
		step3 = stepDone
	case quiesced:
		step3 = stepCurrent
	}
	step4 := stepPending
	caughtUp := false
	if lag, ok := toInt64(targetLag); ok && lag <= zeroRPOLagBytes {
		caughtUp = true
	}
	switch {
	case handedOff:
		step4 = stepDone
	case step3 == stepDone && caughtUp:
		step4 = stepDone
	case step3 == stepDone:
		step4 = stepCurrent
	}
	step5 := stepPending
	if handedOff {
		step5 = stepDone
	} else if step4 == stepDone {
		step5 = stepCurrent
	}
	step6 := stepPending
	if handedOff {
		step6 = stepCurrent
	}

	fmt.Fprintf(out, "  %s 1. target %q validated: healthy and lag known, within the switchover budget\n", step1.mark(), target)
	fmt.Fprintf(out, "  %s 2. quiesce requested on the active DC (%s)\n", step2.mark(), orNone(activeDC))
	fmt.Fprintf(out, "  %s 3. quiesce IN EFFECT: active primary write-locked, its LSN frozen\n", step3.mark())
	lagText := "unknown"
	if targetLag != nil {
		lagText = fmt.Sprintf("%v bytes", targetLag)
	}
	fmt.Fprintf(out, "  %s 4. target caught up to the frozen LSN (now %s, needs <= %d)\n", step4.mark(), lagText, zeroRPOLagBytes)
	fmt.Fprintf(out, "  %s 5. primary-DC Lease handed off to %q\n", step5.mark(), target)
	fmt.Fprintf(out, "  %s 6. old DC demoted to standby, annotations cleared, DR phase back to Steady\n", step6.mark())

	fmt.Fprintf(out, "\n  NEXT: ")
	switch {
	case step1 == stepCurrent:
		fmt.Fprintf(out, "waiting for the target's health and lag to be observable. If it never becomes healthy the switchover cannot start.\n")
	case step3 == stepCurrent:
		fmt.Fprintf(out, "waiting for the write-lock to take hold on %s. This needs the active primary to be UP and reachable; a dead primary can never satisfy it (use the failover path).\n", orNone(activeDC))
	case step4 == stepCurrent:
		fmt.Fprintf(out, "waiting for %q to replay the last %s. This is the only step whose duration depends on your write volume.\n", target, lagText)
	case step5 == stepCurrent:
		fmt.Fprintf(out, "handing off the Lease; the target promotes within seconds.\n")
	case step6 == stepCurrent:
		fmt.Fprintf(out, "the Lease has moved to %q. The old DC self-fences and re-cascades, then the annotations clear and the phase returns to Steady.\n", target)
	default:
		fmt.Fprintf(out, "the operator picks it up on its next reconcile.\n")
	}
}

func toInt64(v any) (int64, bool) {
	switch n := v.(type) {
	case int64:
		return n, true
	case float64:
		return int64(n), true
	case int:
		return int64(n), true
	}
	return 0, false
}

func boolCell(m map[string]any, key string) string {
	if v, ok := m[key].(bool); ok {
		return fmt.Sprintf("%v", v)
	}
	return "-"
}

func orNone(s string) string {
	if s == "" {
		return "<none>"
	}
	return s
}
