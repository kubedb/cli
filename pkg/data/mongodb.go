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
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	api "kubedb.dev/apimachinery/apis/kubedb/v1"
	cs "kubedb.dev/apimachinery/client/clientset/versioned"

	"github.com/Masterminds/semver/v3"
	"github.com/spf13/cobra"
	shell "gomodules.xyz/go-sh"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/klog/v2"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
)

const (
	mgCAFile      = "/var/run/mongodb/tls/ca.crt"
	mgPEMFile     = "/var/run/mongodb/tls/client.pem"
	mgTempPEMFile = "/tmp/client.pem"
)

func InsertMongoDBDataCMD(f cmdutil.Factory) *cobra.Command {
	var (
		dbName string
		rows   int
	)

	insertCmd := &cobra.Command{
		Use: "mongodb",
		Aliases: []string{
			"mg",
		},
		Short:   "Insert data to mongodb",
		Long:    `Use this cmd to insert data into a mongodb database.`,
		Example: `kubectl dba data insert mg -n demo mg-rs --rows 500`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				log.Fatal("Enter mongodb object's name as an argument")
			}
			dbName = args[0]

			namespace, _, err := f.ToRawKubeConfigLoader().Namespace()
			if err != nil {
				klog.Error(err, "failed to get current namespace")
			}

			opts, err := newMongoDBOpts(f, dbName, namespace)
			if err != nil {
				log.Fatalln(err)
			}

			if rows <= 0 {
				log.Fatal("rows need to be greater than 0")
			}

			command := fmt.Sprintf("for(var i=1;i<=%d;i++){db[\"%s\"].insert({_id:\"doc\"+i,actor:\"%s\"})}", rows, KubeDBCollectionName, actor)
			_, err = opts.executeCommand(command)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("\nSuccess! %d documents inserted in MongoDB database %s/%s.\n", rows, opts.db.Namespace, opts.db.Name)
		},
	}
	insertCmd.Flags().IntVarP(&rows, "rows", "r", 100, "number of rows to insert")

	return insertCmd
}

func VerifyMongoDBDataCMD(f cmdutil.Factory) *cobra.Command {
	var (
		dbName string
		rows   int
	)

	verifyCmd := &cobra.Command{
		Use: "mongodb",
		Aliases: []string{
			"mg",
		},
		Short:   "Verify data to a mongodb resource",
		Long:    `Use this cmd to verify data existence in a mongodb object`,
		Example: `kubectl dba data verify mg -n demo mg-rs --rows 500`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				log.Fatal("Enter mongodb object's name as an argument")
			}
			dbName = args[0]

			namespace, _, err := f.ToRawKubeConfigLoader().Namespace()
			if err != nil {
				klog.Error(err, "failed to get current namespace")
			}

			opts, err := newMongoDBOpts(f, dbName, namespace)
			if err != nil {
				log.Fatal(err)
			}

			if rows <= 0 {
				log.Fatal("rows need to be greater than 0")
			}

			command := fmt.Sprintf("db.runCommand({find:\"%s\",filter:{\"actor\":\"%s\"},batchSize:10000})", KubeDBCollectionName, actor)
			command = fmt.Sprintf("JSON.stringify(%s)", command)

			out, err := opts.executeCommand(command)
			if err != nil {
				log.Fatal(err)
			}
			err = opts.verifyOutput(out, rows)
			if err != nil {
				log.Fatal(err)
			}
		},
	}

	verifyCmd.Flags().IntVarP(&rows, "rows", "r", 100, "number of rows to verify")

	return verifyCmd
}

func DropMongoDBDataCMD(f cmdutil.Factory) *cobra.Command {
	var dbName string

	dropCmd := &cobra.Command{
		Use: "mongodb",
		Aliases: []string{
			"mg",
		},
		Short:   "Drop data from mongodb",
		Long:    `Use this cmd to drop data from a mongodb database.`,
		Example: `kubectl dba data drop mg -n demo mg-rs`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				log.Fatal("Enter mongodb object's name as an argument")
			}
			dbName = args[0]

			namespace, _, err := f.ToRawKubeConfigLoader().Namespace()
			if err != nil {
				klog.Error(err, "failed to get current namespace")
			}

			opts, err := newMongoDBOpts(f, dbName, namespace)
			if err != nil {
				log.Fatalln(err)
			}

			command := fmt.Sprintf("db[\"%s\"].drop()", KubeDBCollectionName)
			_, err = opts.executeCommand(command)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("\nSuccess! All the CLI inserted documents DELETED from MongoDB database %s/%s \n", opts.db.Namespace, opts.db.Name)
		},
	}

	return dropCmd
}

type mongoDBOpts struct {
	db       *api.MongoDB
	client   *kubernetes.Clientset
	dbClient *cs.Clientset

	errWriter  *bytes.Buffer
	username   string
	pass       string
	cliCommand string
}

func newMongoDBOpts(f cmdutil.Factory, dbName, namespace string) (*mongoDBOpts, error) {
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

	db, err := dbClient.KubedbV1().MongoDBs(namespace).Get(context.TODO(), dbName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	if db.Status.Phase != api.DatabasePhaseReady {
		return nil, fmt.Errorf("mongodb %s/%s is not ready", namespace, dbName)
	}

	secret, err := client.CoreV1().Secrets(db.Namespace).Get(context.TODO(), db.Spec.AuthSecret.Name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	dbVersion, err := dbClient.CatalogV1alpha1().MongoDBVersions().Get(context.TODO(), db.Spec.Version, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	getCLICommand := func(dbImage string) (string, error) {
		parts := strings.Split(dbImage, ":")
		if len(parts) < 2 {
			return "", fmt.Errorf("%s mongo db version is invalid", dbVersion.Spec.DB.Image)
		}
		vCur := semver.MustParse(parts[1])
		v6, err := semver.NewVersion("6.0.0")
		if err != nil {
			return "", err
		}
		cli := ""
		if vCur.LessThan(v6) {
			cli = "mongo"
		} else {
			cli = "mongosh"
		}
		return cli, nil
	}
	cli, err := getCLICommand(dbVersion.Spec.DB.Image)
	if err != nil {
		return nil, err
	}

	return &mongoDBOpts{
		db:         db,
		client:     client,
		dbClient:   dbClient,
		errWriter:  &bytes.Buffer{},
		username:   string(secret.Data[corev1.BasicAuthUsernameKey]),
		pass:       string(secret.Data[corev1.BasicAuthPasswordKey]),
		cliCommand: cli,
	}, nil
}

func (opts *mongoDBOpts) verifyOutput(out []byte, rows int) error {
	type doc struct {
		DocID string `json:"_id"`
		Actor string `json:"actor"`
	}
	type firstBatch []doc
	type cursor struct {
		FirstBatch firstBatch `json:"firstBatch"`
	}
	type document struct {
		Cursor cursor `json:"cursor"`
	}

	var docs document
	err := json.Unmarshal(out, &docs)
	if err != nil {
		return err
	}

	// first sort the ids
	list := make([]string, 0)
	for i := 0; i < len(docs.Cursor.FirstBatch); i++ {
		d := docs.Cursor.FirstBatch[i].DocID
		list = append(list, d)
	}
	sort.Slice(list, func(i, j int) bool {
		numI, _ := strconv.Atoi(list[i][3:])
		numJ, _ := strconv.Atoi(list[j][3:])
		return numI < numJ // Compare based on the numerical part
	})

	// then, match with the targets
	lim := len(docs.Cursor.FirstBatch)
	if lim > rows {
		lim = rows
	}
	matched := 0
	for i := 0; i < lim; i++ {
		target := fmt.Sprintf("doc%d", i+1)
		if list[i] != target {
			break
		}
		matched = matched + 1
	}
	if matched == rows {
		fmt.Printf("\nSuccess! MongoDB database %s/%s contains: %d keys\n", opts.db.Namespace, opts.db.Name, rows)
	} else {
		fmt.Printf("\nError! Expected keys: %d .MongoDB database %s/%s contains: %d keys\n", rows, opts.db.Namespace, opts.db.Name, matched)
	}
	return nil
}

func (opts *mongoDBOpts) executeCommand(command string) ([]byte, error) {
	shSession, err := opts.getShellCommand(command)
	if err != nil {
		return nil, err
	}
	output, err := shSession.Output()
	if err != nil {
		return nil, err
	}

	errOutput := opts.errWriter.String()
	if errOutput != "" {
		return nil, fmt.Errorf("failed to execute command, stderr: %s", errOutput)
	}
	return output, err
}

func (opts *mongoDBOpts) getShellCommand(command string) (*shell.Session, error) {
	sh := shell.NewSession()
	sh.ShowCMD = false
	sh.Stderr = opts.errWriter

	db := opts.db
	svcName := fmt.Sprintf("svc/%s", db.Name)
	kubectlCommand := []interface{}{
		"exec", "-n", db.Namespace, svcName, "-c", "mongodb", "--",
	}

	mgCommand := []interface{}{
		opts.cliCommand,
	}

	if db.Spec.TLS != nil {
		c, err := opts.handleTLS()
		if err != nil {
			return nil, err
		}
		mgCommand = append(mgCommand, c...)
	} else {
		mgCommand = append(mgCommand,
			KubeDBDatabaseName, "--quiet",
			fmt.Sprintf("--username=%s", opts.username),
			fmt.Sprintf("--password=%s", opts.pass),
			"--authenticationDatabase=admin",
		)
	}

	mgCommand = append(mgCommand, "--eval", command)
	finalCommand := append(kubectlCommand, mgCommand...)

	return sh.Command("kubectl", finalCommand...), nil
}

func (opts *mongoDBOpts) handleTLS() ([]interface{}, error) {
	db := opts.db

	getTLSUser := func(path string) (string, error) {
		data, err := shell.Command("openssl", "x509", "-in", path, "-inform", "PEM", "-subject", "-nameopt", "RFC2253", "-noout").Output()
		if err != nil {
			return "", err
		}

		user := strings.TrimPrefix(string(data), "subject=")
		return strings.TrimSpace(user), nil
	}

	secretName := db.GetCertSecretName(api.MongoDBClientCert, "")
	certSecret, err := opts.client.CoreV1().Secrets(db.Namespace).Get(context.TODO(), secretName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	crt, ok := certSecret.Data[corev1.TLSCertKey]
	if !ok {
		return nil, fmt.Errorf("missing %s in secret %s/%s", corev1.TLSCertKey, certSecret.Namespace, certSecret.Name)
	}

	key, ok := certSecret.Data[corev1.TLSPrivateKeyKey]
	if !ok {
		return nil, fmt.Errorf("missing %s in secret %s/%s", corev1.TLSPrivateKeyKey, certSecret.Namespace, certSecret.Name)
	}

	pem := append(crt[:], []byte("\n")...)
	pem = append(pem, key...)
	err = os.WriteFile(mgTempPEMFile, pem, 0o644)
	if err != nil {
		return nil, err
	}

	sub, err := getTLSUser(mgTempPEMFile)
	if err != nil {
		return nil, err
	}

	mgCommand := []interface{}{
		KubeDBDatabaseName, "--quiet",
		"--tls",
		fmt.Sprintf("--tlsCAFile=%v", mgCAFile),
		fmt.Sprintf("--tlsCertificateKeyFile=%v", mgPEMFile),
		"--authenticationMechanism", "MONGODB-X509",
		"--authenticationDatabase", "$external",
		"-u", sub,
	}
	return mgCommand, nil
}
