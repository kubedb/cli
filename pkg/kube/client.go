package kube

import (
	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
	cmdutil "k8s.io/kubernetes/pkg/kubectl/cmd/util"
)

func NewKubeFactory(cmd *cobra.Command) cmdutil.Factory {
	context := cmdutil.GetFlagString(cmd, "kube-context")
	config := getConfig(context)
	return cmdutil.NewFactory(config)
}

func getConfig(context string) clientcmd.ClientConfig {
	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	rules.DefaultClientConfig = &clientcmd.DefaultClientConfig

	overrides := &clientcmd.ConfigOverrides{ClusterDefaults: clientcmd.ClusterDefaults}

	if context != "" {
		overrides.CurrentContext = context
	}
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, overrides)
}
