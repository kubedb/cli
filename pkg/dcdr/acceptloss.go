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
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/util/i18n"
	"k8s.io/kubectl/pkg/util/templates"
)

// NewCmdAcceptDataLoss releases a cross-DC failover held by the RPO budget. Acts
// on the HUB cluster.
func NewCmdAcceptDataLoss(f cmdutil.Factory) *cobra.Command {
	var yes bool
	cmd := &cobra.Command{
		Use:   "accept-data-loss DB_NAME --yes",
		Short: i18n.T("Release a failover held by the RPO budget, explicitly accepting the data loss"),
		Long: templates.LongDesc(`
			When the surviving data center lags more than
			spec.replication.bestEffortCrossDCLagBytesForFailover (or its lag cannot
			be measured), the promotion is HELD: the un-replicated WAL of the lost
			data center is unrecoverable, so choosing between an outage and a loss
			larger than the budget belongs to a human. This command records that
			decision by setting dr.kubedb.com/accept-failover-data-loss=true; both
			promotion paths (the hub gate and the coordinator's data-plane gate)
			honor it within seconds, and the operator removes the annotation
			automatically once the failover it authorized lands, so it cannot linger
			and approve a later, unrelated loss.

			Where the hold is visible before you decide:
			status.disasterRecovery.protectionMessage (the measured lag),
			condition DCDRPromotionStalled, and dc-dr status.

			KUBECONFIG: the hub cluster.`),
		Example: templates.Examples(`
			# See what would be lost first
			kubectl dba dc-dr status pg-dcdr -n demo

			# Accept it
			kubectl dba dc-dr accept-data-loss pg-dcdr -n demo --yes`),
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
			msg, _, _ := unstructured.NestedString(db.Object, "status", "disasterRecovery", "protectionMessage")
			protected, protectedFound, _ := unstructured.NestedBool(db.Object, "status", "disasterRecovery", "protected")
			if msg != "" {
				cmd.Printf("Current protection verdict: %s\n", msg)
			}
			if protectedFound && protected {
				cmd.Printf("NOTE: the database currently reads protected=true; nothing seems held. Setting the annotation is still safe: it is only honored by a promotion that the budget is actively refusing, and it is auto-removed after use.\n")
			}
			if !yes {
				return fmt.Errorf("this authorizes losing committed data beyond the configured budget; re-run with --yes to confirm")
			}
			v := "true"
			if err := annotateDB(ctx, dyn, ns, args[0], AnnAcceptDataLoss, &v); err != nil {
				return err
			}
			cmd.Printf("Data-loss acceptance recorded on %s/%s. The held promotion proceeds within seconds; the annotation is removed automatically once the failover lands.\n", ns, args[0])
			cmd.Printf("Monitor:  kubectl dba dc-dr status %s -n %s\n", args[0], ns)
			return nil
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "Confirm accepting data loss beyond the configured RPO budget")
	return cmd
}
