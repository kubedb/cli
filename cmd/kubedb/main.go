package main

import (
	"os"

	"github.com/k8sdb/cli/pkg"
)

func main() {
	c := pkg.NewKubedbCommand(os.Stdin, os.Stdout, os.Stderr, Version)
	if err := c.Execute(); err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
