package kube

import (
	"github.com/spf13/cobra"
	cmdutil "k8s.io/kubernetes/pkg/kubectl/cmd/util"
)

type Client struct {
	cmdutil.Factory
}

// New create a new Client
func GetKubeCmd(cmd *cobra.Command) *Client {
	context := cmdutil.GetFlagString(cmd, "kube-context")
	config := GetConfig(context)
	return &Client{
		Factory: cmdutil.NewFactory(config),
	}
}
