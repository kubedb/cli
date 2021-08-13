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

package databases

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	apiv1alpha2 "kubedb.dev/apimachinery/apis/kubedb/v1alpha2"
	cs "kubedb.dev/apimachinery/client/clientset/versioned"
	"kubedb.dev/cli/pkg/lib"

	shell "github.com/codeskyblue/go-sh"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
)

func NewRedisCMD(f cmdutil.Factory) *cobra.Command {
	var (
		dbName    string
		namespace string
		fileName  string
		command   string
		keys      []string
		argv      []string
	)

	currentNamespace, _, err := f.ToRawKubeConfigLoader().Namespace()
	if err != nil {
		klog.Error(err, "failed to get current namespace")
	}

	var rdCmd = &cobra.Command{
		Use: "redis",
		Aliases: []string{
			"rd",
		},
		Short: "Use to operate redis pods",
		Long: `Use this cmd to operate redis pods. Available sub-commands:
				apply
				connect`,
		Run: func(cmd *cobra.Command, args []string) {},
	}

	var rdConnectCmd = &cobra.Command{
		Use:   "connect",
		Short: "Connect to a redis object's pod",
		Long:  `Use this cmd to exec into a redis object's primary pod.`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				log.Fatal("Enter redis object's name as an argument")
			}
			dbName = args[0]
			opts, err := newRedisOpts(f, dbName, namespace, keys, argv)
			if err != nil {
				log.Fatalln(err)
			}

			err = opts.connect()
			if err != nil {
				log.Fatal(err)
			}

		},
	}

	var rdApplyCmd = &cobra.Command{
		Use:   "apply",
		Short: "Apply SQL commands to a redis resource",
		Long: `Use this cmd to apply SQL commands from a file to a redis object's primary pod.
				Syntax: $ kubectl dba redis apply <redis-object-name> -n <namespace> -f <fileName>`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				log.Fatal("Enter redis object's name as an argument. Your commands will be applied on a database inside it's primary pod")
			}
			dbName = args[0]

			opts, err := newRedisOpts(f, dbName, namespace, keys, argv)
			if err != nil {
				log.Fatalln(err)
			}

			if fileName == "" && command == "" {
				log.Fatal("use --file or --command to apply supported commands to a redis object's pods")
			}

			tunnel, err := lib.TunnelToDBService(opts.config, dbName, namespace, apiv1alpha2.RedisDatabasePort)
			if err != nil {
				log.Fatal("couldn't creat tunnel, error: ", err)
			}

			if command != "" {
				err = opts.applyCommand(command)
				if err != nil {
					log.Fatal(err)
				}
			}

			if fileName != "" {
				err = opts.applyFile(fileName)
				if err != nil {
					log.Fatal(err)
				}
			}

			tunnel.Close()
		},
	}

	rdCmd.AddCommand(rdConnectCmd)
	rdCmd.AddCommand(rdApplyCmd)
	rdCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", currentNamespace, "namespace of the redis object to connect to.")

	rdApplyCmd.Flags().StringVarP(&fileName, "file", "f", "", "path to lua script file")
	rdApplyCmd.Flags().StringArrayVarP(&keys, "keys", "k", []string{}, "keys to pass to the lua script, used in script as KEYS[*] and the flag can be specified multiple times with different keys. ")
	rdApplyCmd.Flags().StringArrayVarP(&argv, "argv", "a", []string{}, "args to pass to the lua script, used in script as ARGV[*] and the flag can be specified multiple times with different args. ")
	rdApplyCmd.Flags().StringVarP(&command, "command", "c", "", "single command to execute")

	return rdCmd
}

type redisOpts struct {
	db       *apiv1alpha2.Redis
	config   *rest.Config
	client   *kubernetes.Clientset
	dbClient *cs.Clientset

	errWriter *bytes.Buffer

	keys []string
	args []string
}

func newRedisOpts(f cmdutil.Factory, dbName, namespace string, keys, args []string) (*redisOpts, error) {
	config, err := f.ToRESTConfig()
	if err != nil {
		return nil, err
	}

	client, err := kubernetes.NewForConfig(config)
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

	if db.Status.Phase != apiv1alpha2.DatabasePhaseReady {
		return nil, fmt.Errorf("redis %s/%s is not ready", namespace, dbName)
	}

	return &redisOpts{
		db:        db,
		config:    config,
		client:    client,
		dbClient:  dbClient,
		errWriter: &bytes.Buffer{},
		keys:      keys,
		args:      args,
	}, nil
}

func (opts *redisOpts) getShellCommand(dockerFlags, redisExtraFlags []interface{}) *shell.Session {
	sh := shell.NewSession()
	sh.ShowCMD = false
	sh.Stderr = opts.errWriter

	db := opts.db
	podName := db.Name + "-0"
	if db.Spec.Mode == apiv1alpha2.RedisModeCluster {
		podName = db.StatefulSetNameWithShard(0) + "-0"
	}
	kubectlCommand := []interface{}{
		"exec", "-n", db.Namespace, podName,
	}
	kubectlCommand = append(kubectlCommand, dockerFlags...)

	redisCommand := []interface{}{
		"--", "redis-cli", "-n", "0", "-c",
	}

	if db.Spec.TLS != nil {
		redisCommand = append(redisCommand,
			"--tls",
			"--cert", "/certs/client.crt",
			"--key", "/certs/client.key",
			"--cacert", "/certs/ca.crt",
		)
	}

	kubectlCommand = append(kubectlCommand, "redis")
	finalCommand := append(kubectlCommand, redisCommand...)
	if redisExtraFlags != nil {
		finalCommand = append(finalCommand, redisExtraFlags...)
	}
	return sh.Command("kubectl", finalCommand...).SetStdin(os.Stdin)
}

func (opts *redisOpts) connect() error {
	dockerFlag := []interface{}{
		"-it",
	}
	shSession := opts.getShellCommand(dockerFlag, nil)

	err := shSession.Run()
	if err != nil {
		return err
	}

	return nil
}

func (opts *redisOpts) applyCommand(command string) error {
	if len(opts.keys) != 0 || len(opts.args) != 0 {
		return fmt.Errorf("argv and keys flags are only allowed with lua files, please provide lua file with --file")
	}

	commands := strings.Split(command, " ")
	redisExtraFlags := convertToInterfaceArray(commands)

	shSession := opts.getShellCommand(nil, redisExtraFlags)

	out, err := shSession.Output()
	if err != nil {
		return fmt.Errorf("failed to apply command, error: %s, output: %s\n", err, out)
	}
	output := ""
	if string(out) != "" {
		output = ", output:\n\n" + string(out)
	}

	errOutput := opts.errWriter.String()
	if errOutput != "" {
		return fmt.Errorf("failed to apply command, stderr: %s%s", errOutput, output)
	}
	fmt.Printf("command applied successfully%s", output)

	return nil
}

func (opts *redisOpts) applyFile(fileName string) error {
	fileName, err := filepath.Abs(fileName)
	if err != nil {
		return err
	}
	fileData, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}

	redisExtraFlags := []interface{}{
		"eval", string(fileData), fmt.Sprintf("%v", len(opts.keys)),
	}
	keysIfcArray := convertToInterfaceArray(opts.keys)
	argsIfcArray := convertToInterfaceArray(opts.args)

	redisExtraFlags = append(redisExtraFlags, keysIfcArray...)
	redisExtraFlags = append(redisExtraFlags, argsIfcArray...)

	shSession := opts.getShellCommand(nil, redisExtraFlags)

	out, err := shSession.Output()
	if err != nil {
		fmt.Println(opts.errWriter.String())
		return fmt.Errorf("failed to apply file, error: %s, output: %s\n", err, out)
	}

	output := ""
	if string(out) != "" {
		output = ", output:\n\n" + string(out)
	}

	errOutput := opts.errWriter.String()
	if errOutput != "" {
		return fmt.Errorf("failed to apply file, stderr: %s%s", errOutput, output)
	}

	fmt.Printf("file %s applied successfully%s", fileName, output)

	return nil
}

func convertToInterfaceArray(strs []string) []interface{} {
	interfaceArray := make([]interface{}, len(strs))
	for i := range strs {
		interfaceArray[i] = strs[i]
	}

	return interfaceArray
}
