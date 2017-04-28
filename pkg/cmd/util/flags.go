package util

import "github.com/spf13/cobra"

func AddContextFlag(cmd *cobra.Command) {
	cmd.Flags().String("kube-context", "", "name of the kubeconfig context to use")
}

func AddGetFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("all-namespaces", false, "If present, list the requested object(s) across all namespaces. Namespace in current context is ignored even if specified with --namespace.")
	cmd.Flags().Bool("show-kind", false, "If present, list the resource type for the requested object(s).")

	cmd.Flags().StringP("output", "o", "", "Output format. One of: json|yaml|wide|name.")
	cmd.Flags().BoolP("show-all", "a", false, "When printing, show all resources (default hide terminated pods.)")
	cmd.Flags().Bool("show-labels", false, "When printing, show all labels as the last column (default hide labels column)")
}
