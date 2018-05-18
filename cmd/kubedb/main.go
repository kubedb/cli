package main

import (
	"os"

	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"github.com/kubedb/cli/pkg/cmds"
)

func main() {
	cmd := cmds.NewKubedbCommand(os.Stdin, os.Stdout, os.Stderr, Version)
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
