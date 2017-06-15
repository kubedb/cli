package kube

import (
	"fmt"

	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	cmdutil "k8s.io/kubernetes/pkg/kubectl/cmd/util"
)

func NewKubeFactory(cmd *cobra.Command) cmdutil.Factory {
	context := cmdutil.GetFlagString(cmd, "kube-context")
	config := configForContext(context)
	return cmdutil.NewFactory(config)
}

func NewKubeClient(cmd *cobra.Command) (kubernetes.Interface, error) {
	context := cmdutil.GetFlagString(cmd, "kube-context")
	return getKubeClient(context)
}

// getKubeClient creates a Kubernetes config and client for a given kubeconfig context.
func getKubeClient(context string) (kubernetes.Interface, error) {
	config, err := configForContext(context).ClientConfig()
	if err != nil {
		return nil, err
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("could not get Kubernetes client: %s", err)
	}
	return client, nil
}

func configForContext(context string) clientcmd.ClientConfig {
	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	rules.DefaultClientConfig = &clientcmd.DefaultClientConfig

	overrides := &clientcmd.ConfigOverrides{ClusterDefaults: clientcmd.ClusterDefaults}

	if context != "" {
		overrides.CurrentContext = context
	}
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, overrides)
}
