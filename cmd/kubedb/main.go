package main

import (
	"os"

	"github.com/k8sdb/cli/pkg/cmd"
)

func main() {
	c := cmd.NewKubedbCommand(os.Stdin, os.Stdout, os.Stderr, Version)
	if err := c.Execute(); err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
