/*
Copyright AppsCode Inc. and Contributors

Licensed under the AppsCode Community License 1.0.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://github.com/appscode/licenses/raw/1.0.0/AppsCode-Community-1.0.0.md

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package data

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"kubedb.dev/cli/pkg/data/redisutil"
	"log"
	"os"
	"strconv"
	"strings"

	api "kubedb.dev/apimachinery/apis/kubedb/v1alpha2"
	cs "kubedb.dev/apimachinery/client/clientset/versioned"
	_ "kubedb.dev/db-client-go/redis"

	"github.com/spf13/cobra"
	shell "gomodules.xyz/go-sh"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/klog/v2"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
)

func InsertRedisDataCMD(f cmdutil.Factory) *cobra.Command {
	var (
		dbName string
		rows   int
	)

	rdInsertCmd := &cobra.Command{
		Use: "redis",
		Aliases: []string{
			"rd",
		},
		Short: "Insert data to a redis object's pod",
		Long:  `Use this cmd to insert data into a redis object's primary pod.`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				log.Fatal("Enter redis object's name as an argument")
			}
			dbName = args[0]

			namespace, _, err := f.ToRawKubeConfigLoader().Namespace()
			if err != nil {
				klog.Error(err, "failed to get current namespace")
			}

			opts, err := newRedisOpts(f, dbName, namespace)
			if err != nil {
				log.Fatalln(err)
			}

			err = opts.insertDataInDatabase(rows)
			if err != nil {
				log.Fatal(err)
			}
		},
	}
	rdInsertCmd.Flags().IntVarP(&rows, "rows", "r", 10, "rows in ")

	return rdInsertCmd
}

func VerifyRedisDataCMD(f cmdutil.Factory) *cobra.Command {
	var (
		dbName string
		rows   int
	)

	rdVerifyCmd := &cobra.Command{
		Use: "redis",
		Aliases: []string{
			"rd",
		},
		Short: "Verify rows in a redis database",
		Long:  `Use this cmd to verify data in a redis object`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				log.Fatal("Enter redis object's name as an argument. Your commands will be applied on a database inside it's primary pod")
			}
			dbName = args[0]

			namespace, _, err := f.ToRawKubeConfigLoader().Namespace()
			if err != nil {
				klog.Error(err, "failed to get current namespace")
			}

			opts, err := newRedisOpts(f, dbName, namespace)
			if err != nil {
				log.Fatalln(err)
			}

			err = opts.verifyRedisKeys()
			if err != nil {
				log.Fatalln(err)
			}
		},
	}

	rdVerifyCmd.Flags().IntVarP(&rows, "rows", "r", 10, "rows in ")

	return rdVerifyCmd
}

type redisOpts struct {
	db       *api.Redis
	dbClient *cs.Clientset

	errWriter *bytes.Buffer
}

func newRedisOpts(f cmdutil.Factory, dbName, namespace string) (*redisOpts, error) {
	config, err := f.ToRESTConfig()
	if err != nil {
		return nil, err
	}

	dbClient, err := cs.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	db, err := dbClient.KubedbV1alpha2().Redises(namespace).Get(context.TODO(), dbName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	if db.Status.Phase != api.DatabasePhaseReady {
		return nil, fmt.Errorf("redis %s/%s is not ready", namespace, dbName)
	}

	return &redisOpts{
		db:        db,
		dbClient:  dbClient,
		errWriter: &bytes.Buffer{},
	}, nil
}

var script = `
for i = 1, ARGV[1], 1 do
    redis.call("SET", "{"..ARGV[2].."}-key"..i, tostring({}):sub(10))
end

return "Success!"
`

func (opts *redisOpts) verifyRedisKeys() error {
	return nil
}

func (opts *redisOpts) insertDataInDatabase(rows int) error {
	var slotKey map[int64]string
	err := json.Unmarshal(redisutil.ClusterSlotKeys, &slotKey)
	if err != nil {
		return err
	}
	masterNodes, err := opts.getMasterNodes()
	if err != nil {
		return err
	}
	keysPerMaster := fmt.Sprintf("%d", (rows+len(masterNodes)-1)/len(masterNodes))
	for _, node := range masterNodes {
		hash := slotKey[node.slot]
		redisExtraFlags := []interface{}{
			"eval", script, "0", keysPerMaster, hash,
		}
		shSession := opts.getShellCommand(node.ip, redisExtraFlags)
		out, err := shSession.Output()
		if err != nil {
			return fmt.Errorf("failed to execute command, error: %s, output: %s\n", err, out)
		}

		errOutput := opts.errWriter.String()
		if errOutput != "" {
			return fmt.Errorf("failed to execute command, stderr: %s \n output:%s", errOutput, string(out))
		}

		if strings.TrimSpace(string(out)) == "Success!" {
			fmt.Printf("Successfully insert %s keys in master %s\n", keysPerMaster, node.ip)
		} else {
			fmt.Printf("Error. Can not insert kyes in master %s. Output: %s\n", node.ip, string(out))
		}

	}
	fmt.Printf("%d keys inserted in redis database %s/%s successfully\n", rows, opts.db.Namespace, opts.db.Name)

	return nil
}

type MasterNode struct {
	ip   string
	slot int64
}

func (opts *redisOpts) getMasterNodes() ([]MasterNode, error) {
	if opts.db.Spec.Mode == api.RedisModeCluster {
		return opts.getClusterMasterNodes()
	}
	role, err := opts.getNodeRole()
	if err != nil {
		return nil, err
	}
	var curNode MasterNode
	if role == "master" {
		curNode.slot = 0
		curNode.ip = fmt.Sprintf("%s-0.%s.%s.svc", opts.db.Name, opts.db.GoverningServiceName(), opts.db.Namespace)
	} else {
		curNode.slot = 0
		master, err := opts.getMasterNode()
		if err != nil {
			return nil, err
		}
		curNode.ip = master
	}
	return []MasterNode{curNode}, nil
}

func (opts *redisOpts) getClusterMasterNodes() ([]MasterNode, error) {
	var (
		slotRange []string
		start     int
	)
	nodesConf, err := opts.getNodesConf()
	if err != nil {
		return nil, err
	}
	nodes := strings.Split(nodesConf, "\n")

	var masterNodes []MasterNode

	for _, node := range nodes {
		node = strings.TrimSpace(node)
		parts := strings.Split(strings.TrimSpace(node), " ")

		if strings.Contains(parts[2], "master") {
			var currentNode MasterNode

			for j := 8; j < len(parts); j++ {
				if parts[j][0] == '[' && parts[j][len(parts[j])-1] == ']' {
					continue
				}
				slotRange = strings.Split(parts[j], "-")
				start, _ = strconv.Atoi(slotRange[0])
				currentNode.slot = int64(start)
				break
			}
			currentNode.ip = strings.TrimSuffix(parts[1], ":6379@16379")
			masterNodes = append(masterNodes, currentNode)
		}

	}
	return masterNodes, nil
}

func (opts *redisOpts) getMasterNode() (string, error) {
	redisExtraFlags := []interface{}{
		"info",
	}
	shSession := opts.getShellCommand("", redisExtraFlags)
	out, err := shSession.Output()
	if err != nil {
		return "", fmt.Errorf("failed to execute command, error: %s, output: %s\n", err, out)
	}
	infos := strings.Split(strings.TrimSpace(string(out)), "\n")
	for _, info := range infos {
		info = strings.Trim(info, "\r")
		if strings.Contains(info, "master_host") {
			// The role line will be in this format "role:master".
			// So we are splitting it with ":" and taking the second value after splitting.
			role := strings.Split(info, ":")[1]
			return role, nil
		}
	}
	return "", fmt.Errorf("failed to get pod role")
}

func (opts *redisOpts) getNodeRole() (string, error) {
	redisExtraFlags := []interface{}{
		"info",
	}
	shSession := opts.getShellCommand("", redisExtraFlags)
	out, err := shSession.Output()
	if err != nil {
		return "", fmt.Errorf("failed to execute command, error: %s, output: %s\n", err, out)
	}
	infos := strings.Split(strings.TrimSpace(string(out)), "\n")
	for _, info := range infos {
		info = strings.Trim(info, "\r")
		if strings.Contains(info, "role") {
			// The role line will be in this format "role:master".
			// So we are splitting it with ":" and taking the second value after splitting.
			role := strings.Split(info, ":")[1]
			return role, nil
		}
	}
	return "", fmt.Errorf("failed to get pod role")
}

func (opts *redisOpts) getNodesConf() (string, error) {
	redisExtraFlags := []interface{}{
		"cluster", "nodes",
	}
	shSession := opts.getShellCommand("", redisExtraFlags)
	out, err := shSession.Output()
	if err != nil {
		return "", fmt.Errorf("failed to execute command, error: %s, output: %s\n", err, out)
	}
	return strings.TrimSpace(string(out)), nil
}

func (opts *redisOpts) getShellCommand(podIP string, redisExtraFlags []interface{}) *shell.Session {
	sh := shell.NewSession()
	sh.ShowCMD = false
	sh.Stderr = opts.errWriter

	db := opts.db
	podName := db.Name + "-0"
	redisCommand := []interface{}{
		"--", "redis-cli",
	}
	if len(podIP) != 0 {
		redisCommand = append(redisCommand, "-h", podIP)
	}
	if db.Spec.Mode == api.RedisModeCluster {
		podName = db.StatefulSetNameWithShard(0) + "-0"
	}
	kubectlCommand := []interface{}{
		"exec", "-n", db.Namespace, podName, "-c", "redis",
	}

	if db.Spec.TLS != nil {
		redisCommand = append(redisCommand,
			"--tls",
			"--cert", "/certs/client.crt",
			"--key", "/certs/client.key",
			"--cacert", "/certs/ca.crt",
		)
	}

	finalCommand := append(kubectlCommand, redisCommand...)
	if redisExtraFlags != nil {
		finalCommand = append(finalCommand, redisExtraFlags...)
	}
	return sh.Command("kubectl", finalCommand...).SetStdin(os.Stdin)
}
