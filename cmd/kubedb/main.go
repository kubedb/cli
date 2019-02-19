package main

import (
	"os"
	"github.com/kubedb/cli/pkg/cmds"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

func main() {
	cmd := cmds.NewKubedbCommand(os.Stdin, os.Stdout, os.Stderr)
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
