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

type redisOpts struct {
	db       *api.Redis
	dbClient *cs.Clientset

	errWriter *bytes.Buffer
}
type MasterNode struct {
	host string
	slot int64
}

var dataInsertScript = `
for i = 1, ARGV[1], 1 do
    redis.call("SET", "kubedb:{"..ARGV[2].."}-key"..i, tostring({}):sub(10))
end

return "Success!"
`

var dataDeleteScript = `
local cursor = 0
local calls = 0
local dels = 0
repeat
    local result = redis.call('SCAN', cursor, 'MATCH', ARGV[1])
    calls = calls + 1
    for _,key in ipairs(result[2]) do
        redis.call('DEL', key)
        dels = dels + 1
    end
    cursor = tonumber(result[1])
until cursor == 0

return "Success!"
`

const redisKeyPrefix = "kubedb"

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

func (opts *redisOpts) insertDataInDatabase(rows int) error {
	if opts.db.Spec.Mode == api.RedisModeCluster {
		return opts.insertDataInRedisCluster(rows)
	}

	redisCommand := []interface{}{
		"eval", dataInsertScript, "0", fmt.Sprintf("%d", rows), "hash",
	}
	output, err := opts.execCommand("", redisCommand)
	if err != nil {
		return err
	}
	if output != "Success!" {
		fmt.Printf("Error. Can not insert data in master node. Output: %s\n", output)
	}
	fmt.Printf("\n%d keys inserted in redis database %s/%s successfully\n", rows, opts.db.Namespace, opts.db.Name)
	return nil
}

func (opts *redisOpts) insertDataInRedisCluster(rows int) error {
	var slotKey map[int64]string
	err := json.Unmarshal(redisutil.ClusterSlotKeys, &slotKey)
	if err != nil {
		return err
	}
	masterNodes, err := opts.getClusterMasterNodes()
	if err != nil {
		return err
	}
	keysPerMaster := rows / len(masterNodes)
	extraKey := rows % len(masterNodes)
	for _, node := range masterNodes {
		keyHash := slotKey[node.slot]
		redisCommand := []interface{}{
			"eval", dataInsertScript, "0", fmt.Sprintf("%d", keysPerMaster+extraKey), keyHash,
		}
		extraKey = 0
		output, err := opts.execCommand(node.host, redisCommand)
		if err != nil {
			return err
		}
		if output != "Success!" {
			fmt.Printf("Error. Can not insert data in master %s. Output: %s\n", node.host, output)
		}
	}
	fmt.Printf("\n%d keys inserted in redis database %s/%s successfully\n", rows, opts.db.Namespace, opts.db.Name)
	return nil
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

			err = opts.verifyRedisData(rows)
			if err != nil {
				log.Fatalln(err)
			}
		},
	}

	rdVerifyCmd.Flags().IntVarP(&rows, "rows", "r", 10, "rows in ")

	return rdVerifyCmd
}

func (opts *redisOpts) verifyRedisData(rows int) error {
	if opts.db.Spec.Mode == api.RedisModeCluster {
		return opts.verifyDataInRedisCluster(rows)
	}
	redisCommand := []interface{}{
		"dbsize",
	}
	output, err := opts.execCommand("", redisCommand)
	if err != nil {
		return err
	}
	totalKeys, err := strconv.Atoi(output)
	fmt.Printf("\nExpected keys: %d .Redis database %s/%s contains: %d keys\n", rows, opts.db.Namespace, opts.db.Name, totalKeys)

	return nil
}

func (opts *redisOpts) verifyDataInRedisCluster(rows int) error {
	masterNodes, err := opts.getClusterMasterNodes()
	if err != nil {
		return err
	}
	var totalKeys = 0
	for _, node := range masterNodes {
		redisCommand := []interface{}{
			"dbsize",
		}
		output, err := opts.execCommand(node.host, redisCommand)
		if err != nil {
			return err
		}
		keys, err := strconv.Atoi(output)
		if err != nil {
			return err
		}
		totalKeys += keys
	}
	fmt.Printf("\nExpected keys: %d .Redis database %s/%s contains: %d keys\n", rows, opts.db.Namespace, opts.db.Name, totalKeys)
	return nil
}

func DropRedisDataCMD(f cmdutil.Factory) *cobra.Command {
	var (
		dbName string
	)

	rdVerifyCmd := &cobra.Command{
		Use: "redis",
		Aliases: []string{
			"rd",
		},
		Short: "Delete data from redis database",
		Long:  `Use this cmd to delete inserted data in a redis object`,
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

			err = opts.dropRedisData()
			if err != nil {
				log.Fatalln(err)
			}
		},
	}

	return rdVerifyCmd
}

func (opts *redisOpts) dropRedisData() error {
	if opts.db.Spec.Mode == api.RedisModeCluster {
		return opts.dropRedisClusterData()
	}
	redisCommand := []interface{}{
		"eval", dataDeleteScript, "0", fmt.Sprintf("%s*", redisKeyPrefix),
	}
	output, err := opts.execCommand("", redisCommand)
	if err != nil {
		return err
	}
	if output != "Success!" {
		fmt.Printf("Error. Can not insert data in master node. Output: %s\n", output)
	}
	fmt.Printf("\nAll the CLI inserted keys DELETED drom redis database %s/%s successfully\n", opts.db.Namespace, opts.db.Name)
	return nil
}

func (opts *redisOpts) dropRedisClusterData() error {
	masterNodes, err := opts.getClusterMasterNodes()
	if err != nil {
		return err
	}

	for _, node := range masterNodes {
		redisCommand := []interface{}{
			"eval", dataDeleteScript, "0", fmt.Sprintf("%s*", redisKeyPrefix),
		}
		output, err := opts.execCommand(node.host, redisCommand)
		if err != nil {
			return err
		}
		if output != "Success!" {
			fmt.Printf("Error. Can not insert data in master %s. Output: %s\n", node.host, output)
		}
	}
	fmt.Printf("\nAll the CLI inserted keys DELETED from redis database %s/%s successfully\n", opts.db.Namespace, opts.db.Name)
	return nil
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

func (opts *redisOpts) execCommand(host string, redisCommand []interface{}) (string, error) {
	shSession := opts.getShellCommand(host, redisCommand)
	out, err := shSession.Output()
	if err != nil {
		return "", fmt.Errorf("failed to execute redisCommand, error: %s, output: %s\n", err, out)
	}

	errOutput := opts.errWriter.String()
	if errOutput != "" {
		return "", fmt.Errorf("failed to execute redisCommand, stderr: %s \n output:%s", errOutput, string(out))
	}
	return strings.TrimSpace(string(out)), nil
}

func (opts *redisOpts) getClusterMasterNodes() ([]MasterNode, error) {
	var (
		slotRange []string
		start     int
	)
	nodesConf, err := opts.getClusterNodesConf()
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
			currentNode.host = strings.TrimSuffix(parts[1], ":6379@16379")
			masterNodes = append(masterNodes, currentNode)
		}

	}
	return masterNodes, nil
}

func (opts *redisOpts) getClusterNodesConf() (string, error) {
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

func (opts *redisOpts) getShellCommand(podIP string, redisCommmand []interface{}) *shell.Session {
	sh := shell.NewSession()
	sh.ShowCMD = false
	sh.Stderr = opts.errWriter

	db := opts.db
	redisBaseCommand := []interface{}{
		"--", "redis-cli",
	}
	if len(podIP) != 0 {
		redisBaseCommand = append(redisBaseCommand, "-h", podIP)
	}
	svcName := fmt.Sprintf("svc/%s", db.Name)
	kubectlCommand := []interface{}{
		"exec", "-n", db.Namespace, svcName, "-c", "redis",
	}

	if db.Spec.TLS != nil {
		redisBaseCommand = append(redisBaseCommand,
			"--tls",
			"--cert", "/certs/client.crt",
			"--key", "/certs/client.key",
			"--cacert", "/certs/ca.crt",
		)
	}

	finalCommand := append(kubectlCommand, redisBaseCommand...)
	if redisCommmand != nil {
		finalCommand = append(finalCommand, redisCommmand...)
	}
	return sh.Command("kubectl", finalCommand...).SetStdin(os.Stdin)
}
