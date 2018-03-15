package main

import (
	"log"

	logs "github.com/appscode/go/log/golog"
	_ "github.com/kubedb/apimachinery/client/clientset/versioned/scheme"
	"github.com/kubedb/elasticsearch/pkg/cmds"
)

func main() {
	logs.InitLogs()
	defer logs.FlushLogs()
	if err := cmds.NewRootCmd(Version).Execute(); err != nil {
		log.Fatal(err)
	}
}
