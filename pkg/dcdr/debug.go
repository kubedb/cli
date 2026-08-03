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
	"io"
	"strings"
	"time"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/util/i18n"
	"k8s.io/kubectl/pkg/util/templates"
)

// NewCmdDebug is the `dc-dr debug` group: per-symptom diagnosis walkers.
func NewCmdDebug(f cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "debug",
		Short:                 i18n.T("Diagnose DC-DR symptoms: failover not happening, switchover stuck, fenced database"),
		Long:                  templates.LongDesc(`Walk the DC-DR decision chain and report, in order, every condition that would stop the thing you are waiting for.`),
		Run:                   func(cmd *cobra.Command, args []string) {},
		DisableFlagsInUseLine: true,
		DisableAutoGenTag:     true,
	}
	cmd.AddCommand(
		newCmdDebugFailover(f),
		newCmdDebugSwitchover(f),
		newCmdDebugFence(f),
	)
	return cmd
}

type finding struct {
	ok      bool
	title   string
	detail  string
	remedy  string
	blocker bool
}

func (fd finding) print(w io.Writer) {
	mark := "OK  "
	if !fd.ok {
		mark = "FAIL"
		if !fd.blocker {
			mark = "WARN"
		}
	}
	fmt.Fprintf(w, "  [%s] %s\n", mark, fd.title)
	if fd.detail != "" {
		fmt.Fprintf(w, "         %s\n", fd.detail)
	}
	if !fd.ok && fd.remedy != "" {
		fmt.Fprintf(w, "         -> %s\n", fd.remedy)
	}
}

// newCmdDebugFailover answers "the active DC is gone / my database is down and
// nothing failed over".
func newCmdDebugFailover(f cmdutil.Factory) *cobra.Command {
	var cf CoordFlags
	cmd := &cobra.Command{
		Use:   "failover DB_NAME",
		Short: i18n.T("Diagnose why a cross-DC failover is not happening"),
		Long: templates.LongDesc(`
			Walks every gate between "something is wrong" and "another data center is
			primary", reporting which one is holding:

			  1. is the scope registered at all (no PlacementPolicy failoverPolicy
			     means no Lease, so nothing can ever move);
			  2. does the Lease exist, who holds it, and is it still being renewed
			     (a renewed Lease means the holder's agent is alive: by design NO
			     database-level condition, client errors, QPS, lag, or a crashed
			     postgres, ever moves it);
			  3. is a break-glass pin or a standby-hold blocking the move;
			  4. is the RPO budget holding the promotion (the accept remedy);
			  5. are there stale or failed ForceFailOver ops, or a tripped retry cap,
			     that make the hub skip evaluation.

			KUBECONFIG: the hub cluster; the coordination plane via --coord-*.`),
		Example: templates.Examples(`
			kubectl dba dc-dr debug failover pg-dcdr -n demo`),
		DisableAutoGenTag: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("exactly one database name is required")
			}
			ctx := context.Background()
			out := cmd.OutOrStdout()
			db, dyn, ns, err := getDB(ctx, f, args[0])
			if err != nil {
				return err
			}
			scope, err := ResolveScopeForDB(ctx, dyn, db)
			if err != nil {
				return err
			}
			fmt.Fprintf(out, "Diagnosing failover for %s/%s\n\n", ns, args[0])

			var findings []finding

			// 1. scope registration
			distributed, _, _ := unstructured.NestedBool(db.Object, "spec", "distributed")
			dcdrOn := db.GetAnnotations()["dr.kubedb.com/enabled"] == "true"
			switch {
			case !distributed:
				findings = append(findings, finding{title: "database is DC-DR distributed", detail: "spec.distributed is false", remedy: "DC-DR does not apply to this database", blocker: true})
			case !dcdrOn:
				findings = append(findings, finding{title: "DC-DR is armed on the database", detail: "annotation dr.kubedb.com/enabled is not \"true\"", remedy: "set dr.kubedb.com/enabled=true; without it the per-DC substrate and the coordinator fence are not configured", blocker: true})
			default:
				findings = append(findings, finding{ok: true, title: "database is DC-DR distributed and armed"})
			}
			if strings.Contains(scope.Source, "WARNING") {
				findings = append(findings, finding{title: "failover scope is registered", detail: scope.Source, remedy: "register the scope in the PlacementPolicy (clusterSpreadConstraint.failoverPolicy.trigger); an unregistered scope has no Lease and no protection", blocker: true})
			} else {
				findings = append(findings, finding{ok: true, title: fmt.Sprintf("failover scope resolves to %s", scope.LeaseName), detail: scope.Source})
			}

			// 2. the Lease itself
			coord, cerr := cf.CoordClient(ctx, f)
			if cerr != nil {
				findings = append(findings, finding{title: "coordination plane reachable", detail: cerr.Error(), remedy: "pass --coord-kubeconfig / --coord-kubeconfig-secret; without the Lease no diagnosis of the authority is possible", blocker: true})
				printFindings(out, findings)
				return nil
			}
			lease, lerr := coord.CoordinationV1().Leases(cf.LeaseNS).Get(ctx, scope.LeaseName, metav1.GetOptions{})
			if lerr != nil {
				findings = append(findings, finding{title: "primary-DC Lease exists", detail: lerr.Error(), remedy: "the topology controller creates it from the PlacementPolicy; check the dr-controlplane topology deployment", blocker: true})
				printFindings(out, findings)
				return nil
			}
			holder := ""
			if lease.Spec.HolderIdentity != nil {
				holder = *lease.Spec.HolderIdentity
			}
			dur := int32(45)
			if lease.Spec.LeaseDurationSeconds != nil {
				dur = *lease.Spec.LeaseDurationSeconds
			}
			var age time.Duration
			if lease.Spec.RenewTime != nil {
				age = time.Since(lease.Spec.RenewTime.Time).Round(time.Second)
			}
			switch {
			case holder == "":
				findings = append(findings, finding{ok: true, title: "Lease is unheld", detail: "no data center currently holds it; the first healthy Member to contend will acquire it"})
			case age <= time.Duration(dur)*time.Second:
				findings = append(findings, finding{
					title:  fmt.Sprintf("holder %q is renewing normally, so the authority will NOT move on its own", holder),
					detail: fmt.Sprintf("the Lease was renewed %s ago, inside its %ds duration: that data center's agent is alive and healthy", age, dur),
					remedy: "this is by design: no database-level condition (client errors, QPS, lag, a crashed postgres) ever moves the Lease. If the DATABASE is down but the DC is alive, either let the DC's own raft promote a local peer, or move the scope deliberately: kubectl dba dc-dr handoff " + args[0] + " -n " + ns + " --to <other-dc> --yes",
				})
			default:
				findings = append(findings, finding{ok: true, title: fmt.Sprintf("Lease is EXPIRED (holder %q last renewed %s ago, duration %ds)", holder, age, dur), detail: "a healthy Member DC should acquire it within one retry tick"})
			}

			// 3. pins
			if pin := lease.Annotations["dr.open-cluster-management.io/override-hold"]; pin != "" {
				findings = append(findings, finding{title: "no break-glass pin is blocking the move", detail: fmt.Sprintf("scope is PINNED to %q", pin), remedy: "remove that DC's override ConfigMap: kubectl dba dc-dr pin-primary --scope " + scope.LeaseName + " --remove --yes (run against that DC's spoke)", blocker: true})
			} else {
				findings = append(findings, finding{ok: true, title: "no break-glass pin on the Lease"})
			}
			if ho := lease.Annotations[AnnLeaseHandoffTo]; ho != "" {
				findings = append(findings, finding{title: "no handoff is stuck in flight", detail: fmt.Sprintf("handoff-to=%q is still set, so the target has not acquired yet", ho), remedy: "if the target is standby-held it will never take it: kubectl dba dc-dr pin-standby --scope " + scope.LeaseName + " --remove --yes (against that DC's spoke)"})
			}
			findings = append(findings, checkStandbyHolds(ctx, coord, cf.LeaseNS, scope, holder)...)

			// 4. the RPO budget
			protected, protSet, _ := unstructured.NestedBool(db.Object, "status", "disasterRecovery", "protected")
			protMsg, _, _ := unstructured.NestedString(db.Object, "status", "disasterRecovery", "protectionMessage")
			if protSet && !protected {
				findings = append(findings, finding{
					title:  "promotion is not held by the RPO budget",
					detail: fmt.Sprintf("protected=false: %s", protMsg),
					remedy: "if this is a real failover and the loss is acceptable: kubectl dba dc-dr accept-data-loss " + args[0] + " -n " + ns + " --yes",
				})
			} else if protSet {
				findings = append(findings, finding{ok: true, title: "protection is confirmed (RPO budget satisfied)"})
			}
			if db.GetAnnotations()[AnnAcceptDataLoss] != "" {
				findings = append(findings, finding{ok: true, title: "a data-loss acceptance is currently set", detail: "the budget is bypassed for the promotion it authorizes; the operator removes it once that failover lands"})
			}

			// 5. ops objects and conditions
			findings = append(findings, checkFailoverOps(ctx, f, ns, args[0])...)
			findings = append(findings, checkConditions(db)...)

			printFindings(out, findings)
			fmt.Fprintf(out, "\nAlso useful:\n")
			fmt.Fprintf(out, "  kubectl dba dc-dr status %s -n %s\n", args[0], ns)
			fmt.Fprintf(out, "  kubectl dba dc-dr active-dc %s -n %s\n", args[0], ns)
			return nil
		},
	}
	AddCoordFlags(cmd, &cf)
	return cmd
}

// checkStandbyHolds reports Member DCs whose health Lease is stale, which is the
// only cross-DC-visible hint that a DC cannot take over. The standby-hold marker
// itself is spoke-local and deliberately invisible from here.
func checkStandbyHolds(ctx context.Context, coord kubernetes.Interface, ns string, scope *Scope, holder string) []finding {
	members := scope.MemberDCs
	if len(members) == 0 {
		return nil
	}
	var out []finding
	for _, dc := range members {
		if dc == holder {
			continue
		}
		hl, err := coord.CoordinationV1().Leases(ns).Get(ctx, "dc-health-"+dc, metav1.GetOptions{})
		if err != nil {
			out = append(out, finding{title: fmt.Sprintf("candidate DC %q is reporting health", dc), detail: err.Error(), remedy: "a DC with no health Lease has no running agent; it cannot acquire the primary role"})
			continue
		}
		if hl.Spec.RenewTime != nil {
			age := time.Since(hl.Spec.RenewTime.Time).Round(time.Second)
			dur := int32(15)
			if hl.Spec.LeaseDurationSeconds != nil {
				dur = *hl.Spec.LeaseDurationSeconds
			}
			if age > 3*time.Duration(dur)*time.Second {
				out = append(out, finding{title: fmt.Sprintf("candidate DC %q is healthy", dc), detail: fmt.Sprintf("its health Lease is stale (%s old)", age), remedy: "that DC's agent is down or cannot reach the coordination plane; it cannot take over until it returns"})
				continue
			}
			out = append(out, finding{ok: true, title: fmt.Sprintf("candidate DC %q is alive (health renewed %s ago)", dc, age), detail: "if it still refuses to promote, check for a standby-hold ConfigMap on ITS spoke: kubectl -n " + ns + " get cm " + scope.LeaseName + StandbyHoldCMSuffix})
		}
	}
	return out
}

func checkFailoverOps(ctx context.Context, f cmdutil.Factory, ns, dbName string) []finding {
	cfg, err := f.ToRESTConfig()
	if err != nil {
		return nil
	}
	dyn, err := dynamic.NewForConfig(cfg)
	if err != nil {
		return nil
	}
	list, err := dyn.Resource(PgOpsGVR).Namespace(ns).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil
	}
	var stale, failed, progressing []string
	for i := range list.Items {
		o := &list.Items[i]
		dbRef, _, _ := unstructured.NestedString(o.Object, "spec", "databaseRef", "name")
		typ, _, _ := unstructured.NestedString(o.Object, "spec", "type")
		if dbRef != dbName || typ != "ForceFailOver" {
			continue
		}
		phase, _, _ := unstructured.NestedString(o.Object, "status", "phase")
		switch phase {
		case "Failed":
			failed = append(failed, o.GetName())
		case "Skipped", "Successful":
			if time.Since(o.GetCreationTimestamp().Time) > time.Hour {
				stale = append(stale, o.GetName())
			}
		case "Progressing", "":
			progressing = append(progressing, o.GetName())
		}
	}
	var out []finding
	if len(progressing) > 0 {
		out = append(out, finding{ok: true, title: "a ForceFailOver is in progress", detail: strings.Join(progressing, ", ")})
	}
	if len(failed) > 0 {
		out = append(out, finding{title: "no failed ForceFailOver ops are blocking", detail: "failed: " + strings.Join(failed, ", "), remedy: "read their status for the real cause, then delete them; repeated failures trip the retry cap and the hub stops minting new ones"})
	}
	if len(stale) > 0 {
		out = append(out, finding{title: "no stale ForceFailOver ops are confusing the hub", detail: "old and completed: " + strings.Join(stale, ", "), remedy: "delete them; a stale op targeting the same DC makes the hub read \"already promoted\" and skip evaluation entirely"})
	}
	if len(out) == 0 {
		out = append(out, finding{ok: true, title: "no stale or failed ForceFailOver ops"})
	}
	return out
}

func checkConditions(db *unstructured.Unstructured) []finding {
	conds, _, _ := unstructured.NestedSlice(db.Object, "status", "conditions")
	var out []finding
	for _, c := range conds {
		cm, ok := c.(map[string]any)
		if !ok {
			continue
		}
		typ, _ := cm["type"].(string)
		status, _ := cm["status"].(string)
		msg, _ := cm["message"].(string)
		switch typ {
		case "ForceFailOverRetryCapReached":
			if status == "True" {
				out = append(out, finding{title: "the ForceFailOver retry cap is not tripped", detail: msg, remedy: "the hub has stopped minting new failover ops after repeated failures. Fix the underlying failure, delete the failed ops, and it resumes", blocker: true})
			}
		case "DCDRPromotionStalled":
			if status == "True" {
				out = append(out, finding{title: "no stalled promotion is reported", detail: msg, remedy: "the holder cannot promote; check the coordinator logs on that DC's leader pod"})
			}
		case "DCDRFailoverScopeShared":
			if status == "True" {
				out = append(out, finding{ok: true, title: "NOTE: this scope is SHARED", detail: msg + " (every database in it fails over together)"})
			}
		}
	}
	return out
}

func printFindings(w io.Writer, fs []finding) {
	blocking := 0
	for _, fd := range fs {
		fd.print(w)
		if !fd.ok && fd.blocker {
			blocking++
		}
	}
	fmt.Fprintf(w, "\n%d blocking condition(s) found.\n", blocking)
}

// newCmdDebugSwitchover explains a planned switchover that will not complete.
func newCmdDebugSwitchover(f cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "switchover DB_NAME",
		Short: i18n.T("Diagnose a planned switchover that is not completing"),
		Long: templates.LongDesc(`
			Reports which of the switchover's gates is holding and why, including the
			one that surprises people most: every gate measures by dialing the ACTIVE
			primary, so a switchover cannot proceed while that primary is unreachable.

			KUBECONFIG: the hub cluster.`),
		Example:           templates.Examples(`kubectl dba dc-dr debug switchover pg-dcdr -n demo`),
		DisableAutoGenTag: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("exactly one database name is required")
			}
			ctx := context.Background()
			out := cmd.OutOrStdout()
			db, _, ns, err := getDB(ctx, f, args[0])
			if err != nil {
				return err
			}
			ann := db.GetAnnotations()
			target := ann[AnnSwitchoverTo]
			if target == "" && ann[AnnQuiesceActive] == "" {
				fmt.Fprintf(out, "No switchover is in flight on %s/%s.\n", ns, args[0])
				fmt.Fprintf(out, "Start one:  kubectl dba dc-dr switchover %s -n %s --to <dc>\n", args[0], ns)
				return nil
			}
			var findings []finding
			if ann[AnnSwitchoverAbort] != "" {
				findings = append(findings, finding{title: "no abort is pending", detail: "the abort annotation is set; the hub is unwinding this switchover", remedy: "wait for the annotations to clear, then retry the switchover"})
			}
			if started := ann[AnnSwitchoverStart]; started != "" {
				if t, perr := time.Parse(time.RFC3339, started); perr == nil {
					el := time.Since(t).Round(time.Second)
					fd := finding{ok: true, title: fmt.Sprintf("switchover started %s ago", el)}
					if el > 10*time.Minute {
						fd = finding{title: "switchover is within its timeout", detail: fmt.Sprintf("running for %s, past the default 10m", el), remedy: "it auto-aborts and restores writes to the original active DC; watch for the annotations to clear"}
					}
					findings = append(findings, fd)
				}
			}
			activeDC, _, _ := unstructured.NestedString(db.Object, "status", "disasterRecovery", "activeDC")
			dcs, _, _ := unstructured.NestedSlice(db.Object, "status", "disasterRecovery", "dataCenters")
			var targetFound bool
			for _, d := range dcs {
				dm, ok := d.(map[string]any)
				if !ok {
					continue
				}
				name, _ := dm["clusterName"].(string)
				if name == target {
					targetFound = true
					healthy, _ := dm["healthy"].(bool)
					if !healthy {
						findings = append(findings, finding{title: fmt.Sprintf("target %q is healthy", target), detail: "its health Lease is not fresh", remedy: "the switchover refuses an unhealthy target; fix that DC's agent first", blocker: true})
					} else {
						findings = append(findings, finding{ok: true, title: fmt.Sprintf("target %q is healthy", target)})
					}
					if lag, ok := toInt64(dm["lagBytes"]); ok {
						if lag > zeroRPOLagBytes {
							findings = append(findings, finding{title: "target has caught up to the frozen LSN", detail: fmt.Sprintf("lag is %d bytes, needs <= %d", lag, zeroRPOLagBytes), remedy: "this resolves itself once the quiesce freezes the primary and the target replays; if it never shrinks, the target is not streaming"})
						} else {
							findings = append(findings, finding{ok: true, title: fmt.Sprintf("target lag is %d bytes, within the zero-RPO tolerance", lag)})
						}
					} else {
						findings = append(findings, finding{title: "target lag is known", detail: "no lagBytes reported for the target", remedy: "the hub measures lag by dialing the ACTIVE primary; an unreachable primary makes this permanently unknown and the switchover can never start. Use the failover path instead: dc-dr handoff", blocker: true})
					}
				}
				if name == activeDC {
					if w, ok := dm["writable"].(bool); ok && !w {
						findings = append(findings, finding{ok: true, title: fmt.Sprintf("quiesce is IN EFFECT on %q (write-locked)", activeDC)})
					} else if ann[AnnQuiesceActive] == "true" {
						findings = append(findings, finding{title: fmt.Sprintf("quiesce has taken effect on %q", activeDC), detail: "the quiesce was requested but the active DC still reads writable", remedy: "the write-lock is confirmed by dialing the active primary; if that primary is down or unreachable this never flips and the switchover stalls. Abort (dc-dr abort) and use dc-dr handoff instead", blocker: true})
					}
				}
			}
			if target != "" && !targetFound {
				findings = append(findings, finding{title: fmt.Sprintf("target %q is a known data center", target), detail: "it is not present in status.disasterRecovery.dataCenters", remedy: "check the target name against the PlacementPolicy's Member DCs", blocker: true})
			}
			printFindings(out, findings)
			fmt.Fprintf(out, "\n  kubectl dba dc-dr status %s -n %s\n", args[0], ns)
			fmt.Fprintf(out, "  kubectl dba dc-dr abort %s -n %s\n", args[0], ns)
			return nil
		},
	}
	return cmd
}

// newCmdDebugFence explains a database that is up but refusing writes.
func newCmdDebugFence(f cmdutil.Factory) *cobra.Command {
	var cf CoordFlags
	cmd := &cobra.Command{
		Use:   "fence DB_NAME",
		Short: i18n.T("Diagnose a database whose primary is fenced read-only"),
		Long: templates.LongDesc(`
			A DC-DR database goes read-only by design when its local marker is
			missing, stale past its TTL, or names another data center, so that at most
			one data center is ever writable. This reports which of those applies,
			whether the authority itself is healthy, and what to do when the
			coordination plane is the thing that is broken.

			KUBECONFIG: the hub cluster; the coordination plane via --coord-*.`),
		Example:           templates.Examples(`kubectl dba dc-dr debug fence pg-dcdr -n demo`),
		DisableAutoGenTag: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("exactly one database name is required")
			}
			ctx := context.Background()
			out := cmd.OutOrStdout()
			db, dyn, ns, err := getDB(ctx, f, args[0])
			if err != nil {
				return err
			}
			scope, err := ResolveScopeForDB(ctx, dyn, db)
			if err != nil {
				return err
			}
			var findings []finding
			coord, cerr := cf.CoordClient(ctx, f)
			if cerr != nil {
				findings = append(findings, finding{title: "coordination plane reachable", detail: cerr.Error(), remedy: "if the coordination plane is genuinely down, every DC fences read-only after the marker TTL plus its uncertainty hold. To keep the current primary writable through the outage, pin it: kubectl dba dc-dr pin-primary --scope " + scope.LeaseName + " --yes (run against that DC's spoke)", blocker: true})
				printFindings(out, findings)
				return nil
			}
			lease, lerr := coord.CoordinationV1().Leases(cf.LeaseNS).Get(ctx, scope.LeaseName, metav1.GetOptions{})
			if lerr != nil {
				findings = append(findings, finding{title: "primary-DC Lease exists", detail: lerr.Error(), remedy: "with no Lease there is no marker to project, so every DC fences. Check the dr-controlplane topology controller", blocker: true})
				printFindings(out, findings)
				return nil
			}
			holder := ""
			if lease.Spec.HolderIdentity != nil {
				holder = *lease.Spec.HolderIdentity
			}
			if holder == "" {
				findings = append(findings, finding{title: "the scope has an active data center", detail: "the Lease is currently unheld, so every DC's marker says nobody is active and all of them fence", remedy: "a healthy Member acquires within a retry tick; if none does, check that at least one DC's agent is running", blocker: true})
			} else {
				findings = append(findings, finding{ok: true, title: fmt.Sprintf("the authority says %q is active", holder)})
			}
			if lease.Spec.RenewTime != nil {
				age := time.Since(lease.Spec.RenewTime.Time).Round(time.Second)
				if age > 30*time.Second {
					findings = append(findings, finding{title: "the authority is being renewed", detail: fmt.Sprintf("last renewed %s ago; the marker projected onto each spoke goes stale after 30s and the fence closes fail-closed", age), remedy: "the holder's agent cannot write to the coordination plane. Fix that, or pin the primary to ride out the outage: kubectl dba dc-dr pin-primary --scope " + scope.LeaseName + " --yes", blocker: true})
				} else {
					findings = append(findings, finding{ok: true, title: fmt.Sprintf("the authority is fresh (renewed %s ago)", age)})
				}
			}
			activeDC, _, _ := unstructured.NestedString(db.Object, "status", "disasterRecovery", "activeDC")
			if activeDC != "" && holder != "" && activeDC != holder {
				findings = append(findings, finding{title: "the database agrees with the authority", detail: fmt.Sprintf("status says %q, the Lease says %q", activeDC, holder), remedy: "the status trails the Lease by a reconcile; if it persists, the hub operator is not reconciling this database"})
			}
			fmt.Fprintf(out, "Fence diagnosis for %s/%s (scope %s)\n\n", ns, args[0], scope.LeaseName)
			printFindings(out, findings)
			fmt.Fprintf(out, "\nThe marker each data center actually reads lives on its OWN spoke:\n")
			fmt.Fprintf(out, "  kubectl --kubeconfig <spoke> -n %s get cm %s -o jsonpath='{.data}'\n", cf.LeaseNS, scope.LeaseName)
			return nil
		},
	}
	AddCoordFlags(cmd, &cf)
	return cmd
}
