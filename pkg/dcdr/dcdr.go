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

// Package dcdr implements `kubectl dba dc-dr`, the day-2 command group for KubeDB
// cross data center disaster recovery (DC-DR).
//
// The commands act on up to three different control planes, and each one documents
// which:
//
//   - The HUB cluster (where the KubeDB operator and the database CRs live): the
//     ordinary kubeconfig flags / $KUBECONFIG select it, exactly like every other
//     kubectl-dba command.
//   - The COORDINATION control plane (a separate apiserver holding the primary-DC
//     Leases): reached with a dedicated kubeconfig resolved by CoordFlags, from a
//     file, a Secret, or a ConfigMap (see AddCoordFlags).
//   - A SPOKE cluster (one data center): the pin commands create their marker
//     ConfigMaps on a specific DC's spoke, so for those the ordinary kubeconfig
//     flags must point AT that spoke.
package dcdr

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/util/i18n"
	"k8s.io/kubectl/pkg/util/templates"
)

// GroupVersionResources of the objects this command group touches. Everything is
// accessed unstructured on purpose: the DC-DR status and PlacementPolicy types are
// not part of the released apimachinery this CLI vendors, and the commands only
// need a handful of well-known paths.
var (
	PostgresGVR  = schema.GroupVersionResource{Group: "kubedb.com", Version: "v1", Resource: "postgreses"}
	PlacementGVR = schema.GroupVersionResource{Group: "apps.k8s.appscode.com", Version: "v1", Resource: "placementpolicies"}
	PgOpsGVR     = schema.GroupVersionResource{Group: "ops.kubedb.com", Version: "v1alpha1", Resource: "postgresopsrequests"}
)

// The annotation and naming contract shared with the KubeDB Postgres operator and
// the dr-controlplane service. Keep byte for byte in sync with
// postgres/pkg/dcdr/helpers.go and dr-controlplane/pkg/leases/names.go.
const (
	AnnSwitchoverTo    = "dr.kubedb.com/switchover-to"
	AnnSwitchoverAbort = "dr.kubedb.com/switchover-abort"
	AnnSwitchoverStart = "dr.kubedb.com/switchover-started"
	AnnQuiesceActive   = "dr.kubedb.com/quiesce-active"
	AnnAcceptDataLoss  = "dr.kubedb.com/accept-failover-data-loss"
	AnnFailoverGroup   = "dr.kubedb.com/failover-group"
	AnnMaxLagBytes     = "dr.kubedb.com/switchover-max-lag-bytes"

	AnnLeaseHandoffTo = "dr.open-cluster-management.io/handoff-to"
	AnnLeaseMemberDCs = "dr.open-cluster-management.io/member-dcs"

	GlobalPrimaryLease = "primary-dc"
	GroupLeasePrefix   = "primary-dc-"

	OverrideCMSuffix    = "-override"
	StandbyHoldCMSuffix = "-standby-hold"

	// DefaultCoordNamespace is where the Leases and marker ConfigMaps live, on the
	// coordination plane and on every spoke respectively.
	DefaultCoordNamespace = "dc-failover"

	// zeroRPOLagBytes mirrors the operator's switchover residual-lag tolerance.
	zeroRPOLagBytes = int64(8 * 1024)
)

// NewCmdDCDR returns the `dc-dr` command group.
func NewCmdDCDR(f cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dc-dr",
		Short: i18n.T("Cross data center DR operations: switchover, failover, pins, and diagnosis"),
		Long: templates.LongDesc(`
			Operate a KubeDB database that is distributed across data centers:
			trigger and monitor planned switchovers, accept a held failover's data
			loss, move the failover authority, pin a data center, and diagnose a
			failover that is not happening.`),
		Run:                   func(cmd *cobra.Command, args []string) {},
		DisableFlagsInUseLine: true,
		DisableAutoGenTag:     true,
	}
	cmd.AddCommand(
		NewCmdSwitchover(f),
		NewCmdStatus(f),
		NewCmdAbort(f),
		NewCmdAcceptDataLoss(f),
		NewCmdDebug(f),
		NewCmdHandoff(f),
		NewCmdPin(f, pinPrimary),
		NewCmdPin(f, pinStandby),
		NewCmdActiveDC(f),
	)
	return cmd
}

// CoordFlags resolves the kubeconfig of the coordination control plane, the
// separate apiserver that stores the primary-DC Leases. Resolution order:
//
//  1. --coord-kubeconfig: a kubeconfig FILE for the coordination plane.
//  2. --coord-kubeconfig-configmap [ns/]name: read key "kubeconfig" from that
//     ConfigMap on the CURRENT cluster.
//  3. --coord-kubeconfig-secret [ns/]name (default dc-failover/coord-kubeconfig,
//     the name the dr-controlplane chart mints): read key "kubeconfig" from that
//     Secret on the CURRENT cluster.
type CoordFlags struct {
	File      string
	Secret    string
	ConfigMap string
	// LeaseNS is the namespace holding the Leases on the coordination plane.
	LeaseNS string
}

// AddCoordFlags registers the coordination-plane kubeconfig flags on cmd.
func AddCoordFlags(cmd *cobra.Command, cf *CoordFlags) {
	cmd.Flags().StringVar(&cf.File, "coord-kubeconfig", "", "Path to a kubeconfig file for the coordination control plane (overrides the secret/configmap sources)")
	cmd.Flags().StringVar(&cf.Secret, "coord-kubeconfig-secret", DefaultCoordNamespace+"/coord-kubeconfig", "Secret ([namespace/]name, key \"kubeconfig\") on the current cluster holding the coordination-plane kubeconfig")
	cmd.Flags().StringVar(&cf.ConfigMap, "coord-kubeconfig-configmap", "", "ConfigMap ([namespace/]name, key \"kubeconfig\") on the current cluster holding the coordination-plane kubeconfig")
	cmd.Flags().StringVar(&cf.LeaseNS, "coord-namespace", DefaultCoordNamespace, "Namespace on the coordination plane that holds the primary-DC Leases")
}

func splitNSName(s, defaultNS string) (ns, name string) {
	if ns, name, found := strings.Cut(s, "/"); found {
		return ns, name
	}
	return defaultNS, s
}

// CoordClient builds a client for the coordination control plane per the CoordFlags
// resolution order. The Secret/ConfigMap sources are read through f, the CURRENT
// cluster (normally the hub).
func (cf *CoordFlags) CoordClient(ctx context.Context, f cmdutil.Factory) (kubernetes.Interface, error) {
	var kubeconfigBytes []byte
	switch {
	case cf.File != "":
		b, err := os.ReadFile(cf.File)
		if err != nil {
			return nil, fmt.Errorf("failed to read --coord-kubeconfig file: %w", err)
		}
		kubeconfigBytes = b
	case cf.ConfigMap != "":
		cur, err := f.KubernetesClientSet()
		if err != nil {
			return nil, err
		}
		ns, name := splitNSName(cf.ConfigMap, DefaultCoordNamespace)
		cm, err := cur.CoreV1().ConfigMaps(ns).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to read coordination kubeconfig ConfigMap %s/%s from the current cluster: %w", ns, name, err)
		}
		kubeconfigBytes = []byte(cm.Data["kubeconfig"])
	default:
		cur, err := f.KubernetesClientSet()
		if err != nil {
			return nil, err
		}
		ns, name := splitNSName(cf.Secret, DefaultCoordNamespace)
		sec, err := cur.CoreV1().Secrets(ns).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to read coordination kubeconfig Secret %s/%s from the current cluster (is the current kubeconfig pointing at the hub?): %w", ns, name, err)
		}
		kubeconfigBytes = sec.Data["kubeconfig"]
	}
	if len(kubeconfigBytes) == 0 {
		return nil, fmt.Errorf("resolved an empty coordination-plane kubeconfig (expected key \"kubeconfig\")")
	}
	cfg, err := clientcmd.RESTConfigFromKubeConfig(kubeconfigBytes)
	if err != nil {
		return nil, fmt.Errorf("coordination-plane kubeconfig is not valid: %w", err)
	}
	return kubernetes.NewForConfig(cfg)
}

// Scope is a database's failover scope resolved to its Lease name.
type Scope struct {
	LeaseName string
	// Source explains where the scope came from, for human output.
	Source string
	// MemberDCs are the data-bearing Member DCs from the PlacementPolicy, empty
	// when the policy was not resolvable.
	MemberDCs []string
}

// ResolveScopeForDB mirrors the operator's scopeForDB: the PlacementPolicy's
// failoverPolicy.trigger is the source of truth (Group with a name, else Global);
// the dr.kubedb.com/failover-group annotation is consulted only when the policy
// carries no failoverPolicy; otherwise the scope is Global.
func ResolveScopeForDB(ctx context.Context, dyn dynamic.Interface, db *unstructured.Unstructured) (*Scope, error) {
	ppName, _, _ := unstructured.NestedString(db.Object, "spec", "podTemplate", "spec", "podPlacementPolicy", "name")
	if ppName != "" {
		pp, err := dyn.Resource(PlacementGVR).Get(ctx, ppName, metav1.GetOptions{})
		if err == nil {
			s := &Scope{MemberDCs: memberDCsFromPP(pp)}
			trigger, found, _ := unstructured.NestedMap(pp.Object, "spec", "clusterSpreadConstraint", "failoverPolicy", "trigger")
			if found {
				scope, _ := trigger["scope"].(string)
				group, _ := trigger["group"].(string)
				if scope == "Group" && group != "" {
					s.LeaseName = GroupLeasePrefix + group
					s.Source = fmt.Sprintf("PlacementPolicy %s failoverPolicy trigger (Group %q)", ppName, group)
					return s, nil
				}
				s.LeaseName = GlobalPrimaryLease
				s.Source = fmt.Sprintf("PlacementPolicy %s failoverPolicy trigger (Global)", ppName)
				return s, nil
			}
			// Policy exists but registers no failoverPolicy: back-compat annotation.
			if g := db.GetAnnotations()[AnnFailoverGroup]; g != "" {
				s.LeaseName = GroupLeasePrefix + g
				s.Source = fmt.Sprintf("annotation %s (PlacementPolicy %s has no failoverPolicy)", AnnFailoverGroup, ppName)
				return s, nil
			}
			s.LeaseName = GlobalPrimaryLease
			s.Source = fmt.Sprintf("default Global (PlacementPolicy %s has no failoverPolicy); WARNING: this scope may not be registered, so no Lease may exist and protection may not be armed", ppName)
			return s, nil
		}
		// The policy is referenced but unreadable: fall through to the annotation,
		// but say so.
		if g := db.GetAnnotations()[AnnFailoverGroup]; g != "" {
			return &Scope{LeaseName: GroupLeasePrefix + g, Source: fmt.Sprintf("annotation %s (PlacementPolicy %s unreadable: %v)", AnnFailoverGroup, ppName, err)}, nil
		}
		return &Scope{LeaseName: GlobalPrimaryLease, Source: fmt.Sprintf("default Global (PlacementPolicy %s unreadable: %v)", ppName, err)}, nil
	}
	if g := db.GetAnnotations()[AnnFailoverGroup]; g != "" {
		return &Scope{LeaseName: GroupLeasePrefix + g, Source: "annotation " + AnnFailoverGroup}, nil
	}
	return &Scope{LeaseName: GlobalPrimaryLease, Source: "default Global (no PlacementPolicy set)"}, nil
}

func memberDCsFromPP(pp *unstructured.Unstructured) []string {
	rules, _, _ := unstructured.NestedSlice(pp.Object, "spec", "clusterSpreadConstraint", "distributionRules")
	var members []string
	for _, r := range rules {
		rm, ok := r.(map[string]any)
		if !ok {
			continue
		}
		name, _ := rm["clusterName"].(string)
		indices, _ := rm["replicaIndices"].([]any)
		role, _ := rm["role"].(string)
		// A data-bearing Member: explicit role Member, or (older policies) any rule
		// with replicaIndices. Arbiter/Witness DCs never hold the primary.
		if name == "" {
			continue
		}
		if role == "Member" || (role == "" && len(indices) > 0) {
			members = append(members, name)
		}
	}
	return members
}

// getDB fetches the Postgres CR unstructured, resolving the namespace from the
// factory's kubeconfig flags.
func getDB(ctx context.Context, f cmdutil.Factory, name string) (*unstructured.Unstructured, dynamic.Interface, string, error) {
	ns, _, err := f.ToRawKubeConfigLoader().Namespace()
	if err != nil {
		return nil, nil, "", err
	}
	cfg, err := f.ToRESTConfig()
	if err != nil {
		return nil, nil, "", err
	}
	dyn, err := dynamic.NewForConfig(cfg)
	if err != nil {
		return nil, nil, "", err
	}
	db, err := dyn.Resource(PostgresGVR).Namespace(ns).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, nil, "", fmt.Errorf("failed to get Postgres %s/%s: %w", ns, name, err)
	}
	return db, dyn, ns, nil
}

// annotateDB merge-patches one annotation onto the Postgres CR (nil value removes).
func annotateDB(ctx context.Context, dyn dynamic.Interface, ns, name, key string, value *string) error {
	var v any
	if value != nil {
		v = *value
	}
	patch := fmt.Sprintf(`{"metadata":{"annotations":{%q:%s}}}`, key, jsonValue(v))
	_, err := dyn.Resource(PostgresGVR).Namespace(ns).Patch(ctx, name, mergePatchType, []byte(patch), metav1.PatchOptions{})
	return err
}

func jsonValue(v any) string {
	if v == nil {
		return "null"
	}
	return fmt.Sprintf("%q", v)
}

// requireDistributed prints a loud warning when the database is not a DC-DR
// distributed one; the annotations are harmless there but will do nothing.
func requireDistributed(out func(string, ...any), db *unstructured.Unstructured) {
	distributed, _, _ := unstructured.NestedBool(db.Object, "spec", "distributed")
	if !distributed {
		out("WARNING: %s/%s has spec.distributed=false; DC-DR commands have no effect on it\n", db.GetNamespace(), db.GetName())
	}
}

// markerConfigMap creates or deletes a human-owned marker ConfigMap (break-glass
// override or standby-hold) in the coordination namespace of the CURRENT cluster,
// which for these markers must be the target DC's spoke.
func markerConfigMap(ctx context.Context, f cmdutil.Factory, ns, name string, remove bool) (string, error) {
	cs, err := f.KubernetesClientSet()
	if err != nil {
		return "", err
	}
	if remove {
		if err := cs.CoreV1().ConfigMaps(ns).Delete(ctx, name, metav1.DeleteOptions{}); err != nil {
			return "", err
		}
		return "deleted", nil
	}
	_, err = cs.CoreV1().ConfigMaps(ns).Create(ctx, &corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: ns}}, metav1.CreateOptions{})
	if err != nil {
		return "", err
	}
	return "created", nil
}
