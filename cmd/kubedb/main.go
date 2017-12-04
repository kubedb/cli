package main

import (
	"os"

	"github.com/kubedb/cli/pkg/cmds"
)

func main() {
	cmd := cmds.NewKubedbCommand(os.Stdin, os.Stdout, os.Stderr, Version)
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
