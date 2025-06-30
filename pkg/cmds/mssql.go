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

package cmds

import (
	"context"
	"fmt"
	"kubedb.dev/cli/pkg/common"
	"os"

	"github.com/spf13/cobra"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	_ "kubedb.dev/apimachinery/apis/kubedb/v1alpha2"
	"sigs.k8s.io/yaml"
)

// NewCmdMSSQL creates the parent `mssql` command
func NewCmdMSSQL(f cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mssql",
		Short: "MSSQLServer database commands",
		Long:  "Commands for managing KubeDB MSSQLServer instances.",
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}
	cmd.AddCommand(NewCmdDAGConfig(f))
	return cmd
}

// NewCmdDAGConfig creates the `kubectl dba mssql dag-config` command.
func NewCmdDAGConfig(f cmdutil.Factory) *cobra.Command {
	var (
		namespace  string
		outputDir  string
		desLong    = `Generates a YAML file with the necessary secrets for setting up a MSSQLServer Distributed Availability Group (DAG) remote replica.`
		exampleStr = `  # Generate DAG configuration secrets from MSSQLServer 'ag1' in namespace 'demo'
  kubectl dba mssql dag-config ag1 -n demo`
	)

	cmd := &cobra.Command{
		Use:     "dag-config [mssqlserver-name]",
		Short:   "Generate Distributed Availability Group configuration from a source MSSQLServer",
		Long:    desLong,
		Example: exampleStr,
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			mssqlServerName := args[0]
			// Pass the command's context for cancellation handling
			cmdutil.CheckErr(runDAGConfig(cmd.Context(), f, namespace, outputDir, mssqlServerName))
		},
	}

	cmd.Flags().StringVarP(&namespace, "namespace", "n", "default", "Namespace of the source MSSQLServer")
	cmd.Flags().StringVar(&outputDir, "output-dir", ".", "Directory where the configuration YAML file will be saved")
	return cmd
}

// runDAGConfig is now much simpler. It just orchestrates the steps.
func runDAGConfig(ctx context.Context, f cmdutil.Factory, namespace, outputDir, mssqlServerName string) error {
	fmt.Printf("🔎 Generating DAG configuration for MSSQLServer '%s' in namespace '%s'...\n", mssqlServerName, namespace)

	// Use the new common constructor to get a validated options object
	opts, err := common.NewMSSQLOpts(f, mssqlServerName, namespace)
	if err != nil {
		return err // The error from NewMSSQLOpts will be very informative
	}

	// Generate the YAML buffer using the opts object
	yamlBuffer, err := generateMSSQLDAGConfig(ctx, opts)
	if err != nil {
		return err
	}

	// Write the buffer to a file
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory '%s': %w", outputDir, err)
	}

	outputFile := fmt.Sprintf("%s/%s-dag-config.yaml", outputDir, mssqlServerName)
	if err := os.WriteFile(outputFile, yamlBuffer, 0644); err != nil {
		return fmt.Errorf("failed to write DAG config to file '%s': %w", outputFile, err)
	}

	fmt.Printf("Successfully generated DAG configuration.\n")
	fmt.Printf("Apply this file in your remote cluster: kubectl apply -f %s\n", outputFile)

	return nil
}

// generateMSSQLDAGConfig now takes the opts object and is much more robust.
func generateMSSQLDAGConfig(ctx context.Context, opts *common.MSSQLOpts) ([]byte, error) {
	// IMPROVEMENT: Get secret names directly from the CR status, not by guessing.
	secretNames := []string{
		opts.DB.DbmLoginSecretName(),
		opts.DB.MasterKeySecretName(),
		opts.DB.EndpointCertSecretName(),
	}

	var finalYAML []byte
	for _, secretName := range secretNames {
		fmt.Printf("  - Fetching secret '%s'...\n", secretName)
		// Use the client from the opts object to fetch the secret
		secret, err := opts.Client.CoreV1().Secrets(opts.DB.Namespace).Get(ctx, secretName, metav1.GetOptions{})
		if err != nil {
			// No special error handling needed; if it's not found, something is wrong.
			return nil, fmt.Errorf("failed to get required secret '%s': %w", secretName, err)
		}

		cleanedSecret := cleanupSecretForExport(secret)
		secretYAML, err := yaml.Marshal(cleanedSecret)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal secret '%s' to YAML: %w", secretName, err)
		}
		finalYAML = append(finalYAML, secretYAML...)
		finalYAML = append(finalYAML, []byte("---\n")...)
	}
	return finalYAML, nil
}

// cleanupSecretForExport creates a clean, portable version of a Secret.
func cleanupSecretForExport(secret *core.Secret) *core.Secret {
	return &core.Secret{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Secret",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      secret.Name,
			Namespace: secret.Namespace,
		},
		Data: secret.Data,
		Type: secret.Type,
	}
}
