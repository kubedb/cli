package util

import (
	"github.com/spf13/cobra"
	"k8s.io/kubernetes/pkg/kubectl/resource"
)

func AddGetFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("all-namespaces", false, "If present, list the requested object(s) across all namespaces. Namespace in current context is ignored even if specified with --namespace.")
	cmd.Flags().Bool("show-kind", false, "If present, list the resource type for the requested object(s).")

	cmd.Flags().StringP("output", "o", "", "Output format. One of: json|yaml|wide|name.")
	cmd.Flags().BoolP("show-all", "a", false, "When printing, show all resources (default hide terminated pods.)")
	cmd.Flags().Bool("show-labels", false, "When printing, show all labels as the last column (default hide labels column)")
}

func AddCreateFlags(cmd *cobra.Command, options *resource.FilenameOptions) {
	cmd.Flags().StringSliceVarP(&options.Filenames, "filename", "f", options.Filenames, "Filename to use to create the resource")
	cmd.Flags().BoolVarP(&options.Recursive, "recursive", "R", options.Recursive, "Process the directory used in -f, --filename recursively.")
}

func AddDeleteFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("selector", "l", "", "Selector (label query) to filter on.")
	cmd.Flags().StringP("output", "o", "", "Output mode. Use \"-o name\" for shorter output (resource/name).")
}

func AddDescribeFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("all-namespaces", false, "If present, list the requested object(s) across all namespaces.")
}

func AddFilenameOptionFlags(cmd *cobra.Command, options *resource.FilenameOptions) {
	cmd.Flags().StringSliceVarP(&options.Filenames, "filename", "f", options.Filenames, "Filename to use to create the resource")
	cmd.Flags().BoolVarP(&options.Recursive, "recursive", "R", options.Recursive, "Process the directory used in -f, --filename recursively.")
}
