package app

import (
	"os"

	"github.com/k8sdb/kubedb/pkg/cmd"
)

func Run() error {
	cmd := cmd.NewKubedbCommand(os.Stdin, os.Stdout, os.Stderr)
	return cmd.Execute()
}
