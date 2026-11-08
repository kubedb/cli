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

package adopt

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/util/i18n"
	"k8s.io/kubectl/pkg/util/templates"
)

// NewCmdMilvus returns `adopt milvus`.
func NewCmdMilvus(f cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "milvus",
		Short: i18n.T("Adopt a running Distributed Milvus 2.6 in place"),
		Long: templates.LongDesc(`
			Adopt a Distributed Milvus 2.6 that runs from a Helm chart (or any set of
			Deployments and StatefulSets) in place. The metadata stays in its etcd, the
			segments, indexes and write-ahead log stay in its bucket, the root
			credentials and the client Service stay the same.

			* plan: check the cluster and show what an adoption creates (read-only)
			* run: create the Milvus and hand the roles over, one at a time
			* status: show where each role runs
			* rollback: give the cluster back to the chart's workloads

			RELEASE names the running cluster: by default its workloads are those with
			app.kubernetes.io/instance=RELEASE, and the "component" label names each
			workload's role, as in the upstream Milvus Helm chart.`),
		Run:                   func(cmd *cobra.Command, args []string) {},
		DisableFlagsInUseLine: true,
		DisableAutoGenTag:     true,
	}
	cmd.AddCommand(newCmdMilvusPlan(f), newCmdMilvusRun(f), newCmdMilvusStatus(f), newCmdMilvusRollback(f))
	return cmd
}

func addSourceFlags(cmd *cobra.Command, o *sourceOptions) {
	cmd.Flags().StringVar(&o.Selector, "selector", "", "Label selector of the cluster's workloads (default app.kubernetes.io/instance=RELEASE)")
	cmd.Flags().StringVar(&o.RoleLabel, "role-label", "component", "Pod label that names a workload's Milvus role")
	cmd.Flags().StringArrayVar(&o.Workloads, "workload", nil, "Map a workload to a role, as role[/group]=Kind/name (repeat for each workload; replaces --selector)")
}

func addPlanFlags(cmd *cobra.Command, o *planOptions) {
	addSourceFlags(cmd, &o.sourceOptions)
	cmd.Flags().StringVar(&o.Name, "name", "", "Name of the Milvus object (default: the name of the Service clients use, which KubeDB then takes over)")
	cmd.Flags().StringVar(&o.Version, "version", "", "MilvusVersion to use (default: the one matching the running image)")
	cmd.Flags().StringVar(&o.AuthSecret, "auth-secret", "", "Existing basic-auth Secret with the current root user and password")
	cmd.Flags().StringVar(&o.RootPasswordFile, "root-password-file", "", "File holding the current root password (default: common.security.defaultRootPassword of the running config)")
	cmd.Flags().StringVar(&o.EtcdAuthSecret, "etcd-auth-secret", "", "Basic-auth Secret for an etcd that uses authentication")
	cmd.Flags().StringVar(&o.TLSIssuer, "tls-issuer", "", "cert-manager issuer on the CA clients trust, as name, Issuer/name or ClusterIssuer/name; KubeDB then issues and renews the certificates (found by itself if the running certificate came from cert-manager)")
	cmd.Flags().StringVar(&o.TLSCASecret, "tls-ca-secret", "", "Secret (tls.crt, tls.key) with the CA clients trust; an Issuer is created from it and KubeDB issues and renews the certificates")
	cmd.Flags().StringVar(&o.TLSClientSecret, "tls-client-secret", "", "Existing Secret (tls.crt, tls.key, ca.crt) with the client certificate KubeDB presents when client mTLS is on")
	cmd.Flags().StringVar(&o.ObjectStorageCASecret, "object-storage-ca-secret", "", "Secret NAME[/KEY] (default key ca.crt) with the CA that signed the object storage certificate; the Milvus then reaches the object storage over TLS (spec.objectStorage.tls)")
	cmd.Flags().StringVar(&o.ObjectStorageCAConfigMap, "object-storage-ca-configmap", "", "ConfigMap NAME[/KEY] (default key ca.crt) with that CA, instead of a Secret")
	cmd.Flags().StringVar(&o.ObjectStorageAddress, "object-storage-address", "", "Object storage host when the TLS endpoint differs from the running one (needs a CA flag)")
	cmd.Flags().IntVar(&o.ObjectStoragePort, "object-storage-port", 0, "Object storage port when the TLS endpoint differs from the running one (needs a CA flag)")
	cmd.Flags().StringVar(&o.ConfigFile, "config-file", "", "The effective milvus.yaml, when it is not mounted from a ConfigMap or Secret")
	cmd.Flags().BoolVar(&o.RestrictedSCC, "restricted-scc", false, "Set spec.restrictedSCC (OpenShift)")
	cmd.Flags().BoolVar(&o.SkipCredentialCheck, "skip-credential-check", false, "Do not log in to the running cluster to check the root credentials")
}

func addTimingFlags(cmd *cobra.Command, t *timing) {
	cmd.Flags().DurationVar(&t.Timeout, "timeout", 15*time.Minute, "Longest wait for one step (pods Ready or gone)")
	cmd.Flags().DurationVar(&t.Settle, "settle", 20*time.Second, "Pause after each role")
	cmd.Flags().DurationVar(&t.QueryNodeDelay, "querynode-delay", time.Minute, "Overlap: time a new QueryNode serves before an old one stops (raise it for large collections)")
}

func newCmdMilvusPlan(f cmdutil.Factory) *cobra.Command {
	var o planOptions
	var fo fileOptions
	var outDir string
	cmd := &cobra.Command{
		Use:   "plan [RELEASE]",
		Short: i18n.T("Check a running Milvus and show what adopting it creates"),
		Long: templates.LongDesc(`
			Read-only. Finds the workloads, effective configuration, image and client
			Service of the running cluster, checks everything an in-place adoption
			depends on (Woodpecker WAL, etcd root path, object storage, a matching
			MilvusVersion, the current root password), and shows the Secrets, the
			state ConfigMap and the staged Milvus that "run" would create.

			With --output-dir the objects are written there, unredacted, together with
			adopt.yaml: every option of this adoption as plan found it. Review it and
			pass it to run, status and rollback with -f instead of repeating the flags.
			Edit the objects if needed and apply them; "run" then uses what exists.`),
		Example: templates.Examples(`
			kubectl dba adopt milvus plan my-release -n milvus
			kubectl dba adopt milvus plan my-release -n milvus --output-dir ./adopt
			kubectl dba adopt milvus plan -f ./adopt/adopt.yaml`),
		DisableAutoGenTag: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			if o.Release, err = fo.resolve(cmd, args); err != nil {
				return err
			}
			ctx := context.Background()
			c, err := newCluster(f, cmd.OutOrStdout(), fo.namespace())
			if err != nil {
				return err
			}
			p, err := buildPlan(ctx, c, o)
			if err != nil {
				return err
			}
			out := cmd.OutOrStdout()
			_, _ = fmt.Fprintf(out, "Pre-flight checks for %s in namespace %s:\n", o.Release, c.ns)
			p.Findings.print(out)
			if p.Findings.hasErrors() {
				return errors.New("the adoption is blocked; fix the ERROR findings")
			}
			docs, err := p.yamlDocs(outDir == "")
			if err != nil {
				return err
			}
			if outDir != "" {
				if err := os.MkdirAll(outDir, 0o700); err != nil {
					return err
				}
				af, err := adoptionFileFor(p, o, c.ns)
				if err != nil {
					return err
				}
				docs[adoptionFileName] = string(af)
				for _, name := range sortedKeys(docs) {
					if err := os.WriteFile(filepath.Join(outDir, name), []byte(docs[name]), 0o600); err != nil {
						return err
					}
				}
				_, _ = fmt.Fprintf(out, "\nWrote %s to %s\n", strings.Join(sortedKeys(docs), ", "), outDir)
			} else {
				for _, name := range sortedKeys(docs) {
					_, _ = fmt.Fprintf(out, "\n# %s\n%s", name, docs[name])
				}
			}
			printHandoverPlan(out, p.State, methodOverlap)
			if outDir != "" {
				_, _ = fmt.Fprintf(out, "\nNext: kubectl dba adopt milvus run -f %s --yes\n", filepath.Join(outDir, adoptionFileName))
			} else {
				_, _ = fmt.Fprintf(out, "\nNext: kubectl dba adopt milvus run %s -n %s --yes\n", o.Release, c.ns)
			}
			return nil
		},
	}
	addPlanFlags(cmd, &o)
	addFileFlag(cmd, &fo)
	cmd.Flags().StringVar(&outDir, "output-dir", "", "Write the objects to this directory (Secrets unredacted, mode 0600)")
	return cmd
}

func printHandoverPlan(out io.Writer, st *State, method string) {
	_, _ = fmt.Fprintf(out, "\nHandover (%s), in this order:\n", method)
	for _, role := range handoverOrder(method) {
		for _, u := range st.unitsOf(role) {
			if st.isDone(u) {
				continue
			}
			how := "start KubeDB's, wait Ready, then stop the chart's"
			if method == methodStopStart {
				how = "stop the chart's, then start KubeDB's"
			} else if u.Role == roleQueryNode && u.Replicas > 1 {
				how = "one QueryNode at a time: start one, let it serve, stop one"
			}
			_, _ = fmt.Fprintf(out, "  %-24s %d replicas: %s\n", u.key(), u.Replicas, how)
		}
	}
}

type runOptions struct {
	planOptions
	timing
	Method          string
	Yes             bool
	SkipConfigCheck bool
	ReadyTimeout    time.Duration
}

func newCmdMilvusRun(f cmdutil.Factory) *cobra.Command {
	var o runOptions
	var fo fileOptions
	cmd := &cobra.Command{
		Use:   "run [RELEASE]",
		Short: i18n.T("Adopt a running Milvus: create the Milvus object and hand the roles over"),
		Long: templates.LongDesc(`
			Runs the checks of "plan", creates the Secrets, the state ConfigMap
			NAME-adopt (NAME is the Milvus) and the Milvus with every role at 0 replicas, and waits for
			the operator. Before any pod moves, it compares the config the operator
			rendered with the running one on every setting that decides where the data
			is read from, and stops on a difference.

			Then it hands the roles over one at a time:

			* --method overlap (default): start KubeDB's pods of a role, wait until they
			  are Ready, then stop the chart's. QueryNodes move one at a time. No role
			  stops serving while it moves; it needs room for one extra set of pods of
			  a role (one extra QueryNode).
			* --method stop-start: stop the chart's pods of a role, then start KubeDB's.
			  Each role is briefly unavailable; no extra capacity is needed.

			Progress is recorded in the state ConfigMap: if the command stops, run it
			again and it continues where it was. "rollback" gives the cluster back.

			Do not remove the chart with "helm uninstall" afterwards: it would delete
			the client Service that KubeDB now uses.`),
		Example: templates.Examples(`
			kubectl dba adopt milvus run my-release -n milvus --yes
			kubectl dba adopt milvus run my-release -n milvus --method stop-start --yes
			kubectl dba adopt milvus run -f ./adopt/adopt.yaml --yes`),
		DisableAutoGenTag: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			if o.Release, err = fo.resolve(cmd, args); err != nil {
				return err
			}
			if o.Method != methodOverlap && o.Method != methodStopStart {
				return fmt.Errorf("--method must be %s or %s", methodOverlap, methodStopStart)
			}
			c, err := newCluster(f, cmd.OutOrStdout(), fo.namespace())
			if err != nil {
				return err
			}
			return runAdoption(context.Background(), c, o)
		},
	}
	addPlanFlags(cmd, &o.planOptions)
	addTimingFlags(cmd, &o.timing)
	addFileFlag(cmd, &fo)
	cmd.Flags().StringVar(&o.Method, "method", methodOverlap, "Handover method: overlap or stop-start")
	cmd.Flags().BoolVar(&o.Yes, "yes", false, "Proceed; without it, run only shows what it would do")
	cmd.Flags().BoolVar(&o.SkipConfigCheck, "skip-config-check", false, "Do not stop when the rendered config differs from the running one")
	cmd.Flags().DurationVar(&o.ReadyTimeout, "ready-timeout", 15*time.Minute, "Longest wait for the operator, and for the Milvus to be Ready at the end")
	return cmd
}

func runAdoption(ctx context.Context, c *cluster, o runOptions) error {
	out := c.out
	st, _, err := c.loadState(ctx, o.Release)
	if err != nil {
		return err
	}
	var p *adoptPlan
	if st == nil || st.Phase == phasePlanned {
		o.allowExisting = st != nil
		if p, err = buildPlan(ctx, c, o.planOptions); err != nil {
			return err
		}
		_, _ = fmt.Fprintf(out, "Pre-flight checks for %s in namespace %s:\n", o.Release, c.ns)
		p.Findings.print(out)
		if p.Findings.hasErrors() {
			return errors.New("the adoption is blocked; fix the ERROR findings")
		}
		if st == nil {
			st = p.State
		}
	} else {
		_, _ = fmt.Fprintf(out, "Resuming the adoption of %s as Milvus %s (%s, %s); handed over: %s\n",
			o.Release, st.Name, st.Phase, st.Method, orNone(strings.Join(st.Done, ", ")))
		if st.Method != "" && st.Method != o.Method {
			c.logf("continuing with --method %s (started with %s)", o.Method, st.Method)
		}
	}
	if st.Phase == phaseAdopted {
		_, _ = fmt.Fprintf(out, "Milvus %s is already adopted.\n", st.Name)
		return nil
	}
	if !o.Yes {
		printHandoverPlan(out, st, o.Method)
		return errors.New("re-run with --yes to proceed")
	}

	if p != nil {
		if err := c.createObjects(ctx, p); err != nil {
			return err
		}
		if err := c.removeGoverningServices(ctx, st, p.GoverningServices); err != nil {
			return err
		}
		if _, err := c.createMilvus(ctx, p, o.ReadyTimeout); err != nil {
			return err
		}
		if !o.SkipConfigCheck {
			diffs, err := c.checkRenderedConfig(ctx, st.Name, p.ObjStore.expected(p.Source.Config))
			if err != nil {
				return err
			}
			if len(diffs) > 0 {
				_, _ = fmt.Fprintf(out, "The config the operator rendered for Milvus %s differs from the running one:\n", st.Name)
				for _, d := range diffs {
					_, _ = fmt.Fprintf(out, "  %s\n", d)
				}
				return fmt.Errorf("no pod was moved; fix Secret %s or %s and run again, or pass --skip-config-check", configSecretName(st.Name), objectStorageSecretName(st.Name))
			}
			c.logf("rendered config matches the running one (storage, etcd root path, WAL, security, ports)")
		}
		if err := c.checkIssuedCertificates(ctx, st.Name, p.TLS, o.ReadyTimeout); err != nil {
			return err
		}
	} else if _, err := c.dyn.Resource(MilvusGVR).Namespace(c.ns).Get(ctx, st.Name, metav1.GetOptions{}); kerr.IsNotFound(err) {
		return fmt.Errorf("the Milvus %s is gone while %s; run \"rollback\" to give the cluster back", st.Name, st.Phase)
	} else if err != nil {
		return err
	}

	st.Phase, st.Method = phaseHandingOver, o.Method
	if err := c.saveState(ctx, st); err != nil {
		return err
	}
	start := time.Now()
	if err := c.handOver(ctx, st, o.Method, o.timing); err != nil {
		return fmt.Errorf("%w\nThe handover stopped; run the same command to continue, or \"rollback\" to give the cluster back", err)
	}
	if err := c.waitMilvusReady(ctx, st.Name, o.ReadyTimeout); err != nil {
		return err
	}
	st.Phase = phaseAdopted
	if err := c.saveState(ctx, st); err != nil {
		return err
	}
	c.logf("Milvus %s is Ready; adopted in %s", st.Name, time.Since(start).Round(time.Second))
	_, _ = fmt.Fprintf(out, `
Next:
  - check the data from a client (collection count, known queries)
  - the chart's workloads are kept at 0 replicas; do not run "helm uninstall" (it would delete
    the client Service KubeDB now uses), and do not "helm upgrade" (it would start them again)
  - to give the cluster back: kubectl dba adopt milvus rollback %s -n %s --yes
`, o.Release, c.ns)
	return nil
}

func newCmdMilvusStatus(f cmdutil.Factory) *cobra.Command {
	var fo fileOptions
	cmd := &cobra.Command{
		Use:               "status [RELEASE]",
		Short:             i18n.T("Show where each role of an adoption runs"),
		Example:           templates.Examples(`kubectl dba adopt milvus status my-release -n milvus`),
		DisableAutoGenTag: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			release, err := fo.resolve(cmd, args)
			if err != nil {
				return err
			}
			ctx := context.Background()
			c, err := newCluster(f, cmd.OutOrStdout(), fo.namespace())
			if err != nil {
				return err
			}
			st, _, err := c.loadState(ctx, release)
			if err != nil {
				return err
			}
			if st == nil {
				return fmt.Errorf("no adoption of %s in namespace %s", release, c.ns)
			}
			return c.printStatus(ctx, st)
		},
	}
	addFileFlag(cmd, &fo)
	return cmd
}

func (c *cluster) printStatus(ctx context.Context, st *State) error {
	phase := "not found"
	milvusExists := false
	if db, err := c.dyn.Resource(MilvusGVR).Namespace(c.ns).Get(ctx, st.Name, metav1.GetOptions{}); err == nil {
		milvusExists = true
		phase, _, _ = unstructured.NestedString(db.Object, "status", "phase")
	} else if !kerr.IsNotFound(err) {
		return err
	}
	_, _ = fmt.Fprintf(c.out, "Adoption of %s: %s (method %s). Milvus %s: %s\n\n", st.Release, st.Phase, orNone(st.Method), st.Name, phase)
	w := tabwriter.NewWriter(c.out, 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprintln(w, "ROLE\tCHART WORKLOAD\tCHART\tKUBEDB\tHANDED OVER")
	for _, role := range milvusRoles {
		for _, u := range st.unitsOf(role) {
			var kd side = petSetSide{c, st.Name, u}
			if milvusExists {
				kd = kubedbSide{c, st.Name, u}
			}
			_, _ = fmt.Fprintf(w, "%s\t%s %s\t%s\t%s\t%v\n", u.key(), u.Kind, u.Name,
				c.sideSummary(ctx, chartSide{c, u}), c.sideSummary(ctx, kd), st.isDone(u))
		}
	}
	return w.Flush()
}

func (c *cluster) sideSummary(ctx context.Context, s side) string {
	d, err := s.desired(ctx)
	if err != nil {
		return "?"
	}
	r, _, err := c.podCounts(ctx, s.selector())
	if err != nil {
		return "?"
	}
	return fmt.Sprintf("%d/%d Ready", r, d)
}

type rollbackOptions struct {
	sourceOptions
	timing
	Name    string
	Yes     bool
	Cleanup bool
}

func newCmdMilvusRollback(f cmdutil.Factory) *cobra.Command {
	var o rollbackOptions
	var fo fileOptions
	cmd := &cobra.Command{
		Use:   "rollback [RELEASE]",
		Short: i18n.T("Give an adopted (or partly adopted) Milvus back to the chart's workloads"),
		Long: templates.LongDesc(`
			Deletes the Milvus object with kubedb.com/milvus-orphan-on-delete, so its
			PetSets keep serving, then moves each role back to the chart's workloads,
			starting the chart's pods before it stops KubeDB's. Once the chart's
			proxies are Ready, the client Service gets back its original selector,
			labels and ports. The data is never touched.

			With --cleanup it then deletes the PetSets (at 0 replicas), the Services
			other than the client Service and the rendered config KubeDB left, and
			recreates the chart Services the adoption removed. PVCs are kept.

			For a Milvus adopted by hand (no state ConfigMap), give its --name: the
			chart's replicas are then taken from the PetSets, and the client Service
			gets the chart proxies' selector.`),
		Example: templates.Examples(`
			kubectl dba adopt milvus rollback my-release -n milvus --yes --cleanup
			kubectl dba adopt milvus rollback -f ./adopt/adopt.yaml --yes`),
		DisableAutoGenTag: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			if o.Release, err = fo.resolve(cmd, args); err != nil {
				return err
			}
			c, err := newCluster(f, cmd.OutOrStdout(), fo.namespace())
			if err != nil {
				return err
			}
			return runRollback(context.Background(), c, o)
		},
	}
	addSourceFlags(cmd, &o.sourceOptions)
	addTimingFlags(cmd, &o.timing)
	addFileFlag(cmd, &fo)
	cmd.Flags().StringVar(&o.Name, "name", "", "Name of the Milvus, for an adoption done by hand (no state ConfigMap)")
	cmd.Flags().BoolVar(&o.Yes, "yes", false, "Proceed; without it, rollback only shows what it would do")
	cmd.Flags().BoolVar(&o.Cleanup, "cleanup", false, "Delete what KubeDB left behind once the chart serves again (PVCs are kept)")
	return cmd
}

func runRollback(ctx context.Context, c *cluster, o rollbackOptions) error {
	st, cm, err := c.loadState(ctx, o.Release)
	if err != nil {
		return err
	}
	if st == nil {
		if o.Name == "" {
			return fmt.Errorf("no adoption of %s in namespace %s; for a Milvus adopted by hand, give its --name", o.Release, c.ns)
		}
		if st, err = c.stateFromManualAdoption(ctx, o); err != nil {
			return err
		}
	}
	out := c.out
	_, _ = fmt.Fprintf(out, "Rollback of Milvus %s to the workloads of %s:\n", st.Name, o.Release)
	for _, role := range rollbackOrder {
		for _, u := range st.unitsOf(role) {
			_, _ = fmt.Fprintf(out, "  %-24s -> %s %s, %d replicas\n", u.key(), u.Kind, u.Name, u.Replicas)
		}
	}
	if st.Service != nil {
		_, _ = fmt.Fprintf(out, "  Service %s -> selector %v\n", st.Service.Name, st.Service.Selector)
	}
	if !o.Yes {
		return errors.New("re-run with --yes to proceed")
	}
	start := time.Now()
	if err := c.releaseMilvus(ctx, st.Name, o.Timeout); err != nil {
		return err
	}
	if err := c.rollBack(ctx, st, o.timing); err != nil {
		if cm != nil {
			_ = c.saveState(ctx, st)
		}
		return fmt.Errorf("%w\nThe rollback stopped; run the same command to continue", err)
	}
	if o.Cleanup {
		if err := c.cleanupKubeDBObjects(ctx, st); err != nil {
			return err
		}
	}
	if cm != nil {
		if err := c.kube.CoreV1().ConfigMaps(c.ns).Delete(ctx, cm.Name, metav1.DeleteOptions{}); err != nil && !kerr.IsNotFound(err) {
			return err
		}
	}
	c.logf("the chart's workloads serve %s again; rolled back in %s", o.Release, time.Since(start).Round(time.Second))
	if !o.Cleanup {
		_, _ = fmt.Fprintf(out, "KubeDB's PetSets (0 replicas), Services and rendered config are left in place; re-run with --cleanup to delete them.\n")
	}
	_, _ = fmt.Fprintf(out, "The Secrets with label %s=%s (credentials and the carried config) are kept.\n", LabelRelease, o.Release)
	return nil
}

// stateFromManualAdoption rebuilds the state of a Milvus adopted by hand: the
// chart workloads from the release, their replicas from the KubeDB PetSets.
func (c *cluster) stateFromManualAdoption(ctx context.Context, o rollbackOptions) (*State, error) {
	src, err := discover(ctx, c.kube, c.ns, o.sourceOptions)
	if err != nil {
		return nil, err
	}
	if len(src.Units) == 0 {
		return nil, fmt.Errorf("found no workloads of %s", o.Release)
	}
	st := &State{Release: o.Release, Name: o.Name, Phase: phaseHandingOver, Units: src.Units}
	for i := range st.Units {
		u := &st.Units[i]
		n, err := petSetSide{c, o.Name, *u}.desired(ctx)
		if err != nil {
			return nil, err
		}
		if chart, err := (chartSide{c, *u}).desired(ctx); err == nil && chart > n {
			n = chart
		}
		if n == 0 {
			return nil, fmt.Errorf("cannot tell how many replicas %s %s should run: it and PetSet %s are both at 0", u.Kind, u.Name, petSetName(o.Name, *u))
		}
		u.Replicas = n
	}
	svc, err := c.kube.CoreV1().Services(c.ns).Get(ctx, o.Name, metav1.GetOptions{})
	if err == nil {
		proxy := src.firstUnit(roleProxy)
		st.Service = fallbackServiceRecord(*svc, proxy, src.Templates[proxy.key()])
	} else if !kerr.IsNotFound(err) {
		return nil, err
	}
	return st, nil
}

func orNone(s string) string {
	if s == "" {
		return "none"
	}
	return s
}

func sortedKeys(m map[string]string) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
