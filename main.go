package main

import (
	"os"

	"github.com/k8sdb/kubedb/pkg/cmd"
)

func main() {
	cmd := cmd.NewKubedbCommand(os.Stdin, os.Stdout, os.Stderr)
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
