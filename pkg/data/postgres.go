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
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"

	api "kubedb.dev/apimachinery/apis/kubedb/v1alpha2"
	cs "kubedb.dev/apimachinery/client/clientset/versioned"

	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
)

const (
	pgCaFile   = "/tls/certs/client/ca.crt"
	pgCertFile = "/tls/certs/client/client.crt"
	pgKeyFile  = "/tls/certs/client/client.key"
	rowLimit   = 100000
)

type postgresOpts struct {
	db       *api.Postgres
	config   *rest.Config
	client   *kubernetes.Clientset
	dbClient *cs.Clientset

	username string
	pass     string
}

func newPostgresOpts(f cmdutil.Factory, dbName, namespace string) (*postgresOpts, error) {
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

	db, err := dbClient.KubedbV1alpha2().Postgreses(namespace).Get(context.TODO(), dbName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	if db.Status.Phase != api.DatabasePhaseReady {
		return nil, fmt.Errorf("postgres %s/%s is not ready", namespace, dbName)
	}

	secret, err := client.CoreV1().Secrets(db.Namespace).Get(context.TODO(), db.Spec.AuthSecret.Name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return &postgresOpts{
		db:       db,
		config:   config,
		client:   client,
		dbClient: dbClient,
		username: string(secret.Data[corev1.BasicAuthUsernameKey]),
		pass:     string(secret.Data[corev1.BasicAuthPasswordKey]),
	}, nil
}

func InsertPostgresDataCMD(f cmdutil.Factory) *cobra.Command {
	var (
		dbName string
		rows   int
	)

	pgInsertCmd := &cobra.Command{
		Use: "postgres",
		Aliases: []string{
			"postgresql",
			"pgsql",
			"pg",
		},
		Short:   "Insert data to a postgres object's pod",
		Long:    `Use this cmd to insert data into a postgres object's primary pod.`,
		Example: `kubectl dba insert postgres -n demo sample-postgres --rows 500`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				log.Fatal("Enter postgres object's name as an argument. Your commands will be applied on a database inside it's primary pod")
			}
			dbName = args[0]

			namespace, _, err := f.ToRawKubeConfigLoader().Namespace()
			if err != nil {
				klog.Error(err, "failed to get current namespace")
			}

			opts, err := newPostgresOpts(f, dbName, namespace)
			if err != nil {
				log.Fatalln(err)
			}

			if rows <= 0 {
				log.Fatal("rows need to be greater than 0")
			}

			if rows <= rowLimit {
				err = opts.insertDataExecCmd(rows)
				if err != nil {
					log.Fatal(err)
				}
			} else {
				log.Printf("At most %v rows can be inserted per operation", rowLimit)
			}
		},
	}

	pgInsertCmd.Flags().IntVarP(&rows, "rows", "r", 100, "number of rows to insert")

	return pgInsertCmd
}

func (opts *postgresOpts) insertDataExecCmd(rows int) error {
	if rows <= 0 {
		return fmt.Errorf("rows need to be greater than 0")
	}
	command := fmt.Sprintf(`create table if not exists appscode_kubedb_postgres_test_table (values int not null);insert into appscode_kubedb_postgres_test_table (values) values(generate_series(1,%v))`, rows)
	out, err := opts.execCommand(command)
	if err != nil {
		return err
	}
	if strings.Contains(strings.TrimSpace(out), strconv.Itoa(rows)) {
		fmt.Printf("\nSuccess! %d keys inserted in postgres database %s/%s.\n", rows, opts.db.Namespace, opts.db.Name)
	} else {
		fmt.Printf("Error. Can not insert data properly in master %s\n. Output: %v\n", opts.db.Name, out)
	}
	return nil
}

func VerifyPostgresDataCMD(f cmdutil.Factory) *cobra.Command {
	var (
		dbName string
		rows   int
	)

	pgVerifyCmd := &cobra.Command{
		Use: "postgres",
		Aliases: []string{
			"postgresql",
			"pgsql",
			"pg",
		},
		Short:   "Verify rows in a postgres database",
		Long:    `Use this cmd to verify data in a postgres object`,
		Example: `kubectl dba verify pg -n demo sample-postgres --rows 500`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				log.Fatal("Enter postgres object's name as an argument. Your commands will be applied on a database inside it's primary pod")
			}
			dbName = args[0]

			namespace, _, err := f.ToRawKubeConfigLoader().Namespace()
			if err != nil {
				klog.Error(err, "failed to get current namespace")
			}

			opts, err := newPostgresOpts(f, dbName, namespace)
			if err != nil {
				log.Fatalln(err)
			}

			err = opts.verifyDataExecCmd(rows)
			if err != nil {
				log.Fatal(err)
			}
		},
	}
	pgVerifyCmd.Flags().IntVarP(&rows, "rows", "r", 100, "number of rows to verify")

	return pgVerifyCmd
}

func (opts *postgresOpts) verifyDataExecCmd(rows int) error {
	if rows <= 0 {
		return fmt.Errorf("rows need to be greater than 0")
	}
	command := `SELECT COUNT(*) FROM appscode_kubedb_postgres_test_table`
	out, err := opts.execCommand(command)
	if err != nil {
		return err
	}
	output := strings.Split(out, "\n")

	found := strings.TrimSpace(output[2])
	totalRows, err := strconv.Atoi(found)
	if err != nil {
		return err
	}
	if totalRows == rows {
		fmt.Printf("\nSuccess! Postgres database %s/%s contains: %d Rows\n", opts.db.Namespace, opts.db.Name, totalRows)
	} else {
		fmt.Printf("\nError! Expected keys: %d .Postgres database %s/%s contains: %d Rows\n", rows, opts.db.Namespace, opts.db.Name, totalRows)
	}
	return nil
}

func DropPostgresDataCMD(f cmdutil.Factory) *cobra.Command {
	var dbName string

	pgDropCmd := &cobra.Command{
		Use: "postgres",
		Aliases: []string{
			"postgresql",
			"pgsql",
			"pg",
		},
		Short:   "Delete data from postgres database",
		Long:    `Use this cmd to delete inserted data in a postgres object`,
		Example: `kubectl dba drop pg -n demo sample-postgres`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				log.Fatal("Enter postgres object's name as an argument. Your commands will be applied on a database inside it's primary pod")
			}
			dbName = args[0]

			namespace, _, err := f.ToRawKubeConfigLoader().Namespace()
			if err != nil {
				klog.Error(err, "failed to get current namespace")
			}

			opts, err := newPostgresOpts(f, dbName, namespace)
			if err != nil {
				log.Fatalln(err)
			}

			err = opts.dropDataExecCmd()
			if err != nil {
				log.Fatal(err)
			}
		},
	}

	return pgDropCmd
}

func (opts *postgresOpts) dropDataExecCmd() error {
	command := `DROP TABLE if exists appscode_kubedb_postgres_test_table`
	_, err := opts.execCommand(command)
	if err != nil {
		return err
	}
	fmt.Printf("\nSuccess: All the CLI inserted rows DELETED from postgres database %s/%s \n", opts.db.Namespace, opts.db.Name)
	return nil
}

func (opts *postgresOpts) execCommand(command string) (string, error) {
	cmd := opts.getShellCommand(command)

	output, err := opts.runCMD(cmd)
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func (opts *postgresOpts) getShellCommand(command string) string {
	db := opts.db
	svcName := fmt.Sprintf("svc/%s", db.Name)

	cmd := ""
	if db.Spec.TLS != nil {
		if db.Spec.ClientAuthMode == api.ClientAuthModeCert {
			cmd = fmt.Sprintf("kubectl exec -n %s %s -c postgres -- env PGSSLMODE='%s' PGSSLROOTCERT='%s' PGSSLCERT='%s' PGSSLKEY='%s' PGPASSWORD='%s' psql -d postgres -U %s -c '%s'", db.Namespace, svcName, db.Spec.SSLMode, pgCaFile, pgCertFile, pgKeyFile, opts.pass, opts.username, command)
		} else {
			cmd = fmt.Sprintf("kubectl exec -n %s %s -c postgres -- env PGSSLMODE='%s' PGSSLROOTCERT='%s' PGPASSWORD='%s' psql -d postgres -U %s -c '%s'", db.Namespace, svcName, db.Spec.SSLMode, pgCaFile, opts.pass, opts.username, command)
		}
	} else {
		cmd = fmt.Sprintf("kubectl exec -n %s %s -c postgres -- env PGSSLMODE=%s PGPASSWORD='%s' psql -d postgres -U %s -c '%s'", db.Namespace, svcName, db.Spec.SSLMode, opts.pass, opts.username, command)
	}

	return cmd
}

func (opts *postgresOpts) runCMD(cmd string) ([]byte, error) {
	sh := exec.Command("/bin/sh", "-c", cmd)
	stdout := bytes.NewBuffer(nil)
	stderr := bytes.NewBuffer(nil)
	sh.Stdout = stdout
	sh.Stderr = stderr
	err := sh.Run()
	out := stdout.Bytes()
	errOut := stderr.Bytes()
	errOutput := string(errOut)
	if errOutput != "" && !strings.Contains(errOutput, "NOTICE") {
		return nil, fmt.Errorf("failed to execute command, stderr: %s", errOutput)
	}
	if err != nil {
		return nil, err
	}
	return out, nil
}
