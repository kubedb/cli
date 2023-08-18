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
	"os"
	"strconv"
	"strings"

	api "kubedb.dev/apimachinery/apis/kubedb/v1alpha2"
	cs "kubedb.dev/apimachinery/client/clientset/versioned"
	"kubedb.dev/cli/pkg/lib"

	"github.com/spf13/cobra"
	shell "gomodules.xyz/go-sh"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"kmodules.xyz/client-go/tools/portforward"
)

const (
	pgCaFile   = "/tmp/root.crt"
	pgCertFile = "/tmp/postgresql.crt"
	pgKeyFile  = "/tmp/postgresql.key"
)

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

			tunnel, err := lib.TunnelToDBService(opts.config, dbName, namespace, api.PostgresDatabasePort)
			if err != nil {
				log.Fatal("couldn't create tunnel, error: ", err)
			}
			err = opts.insertDataExecCmd(tunnel, rows)
			if err != nil {
				log.Fatal(err)
			}
			tunnel.Close()
		},
	}

	pgInsertCmd.Flags().IntVarP(&rows, "rows", "r", 100, "number of rows to insert")

	return pgInsertCmd
}

func (opts *postgresOpts) insertDataExecCmd(tunnel *portforward.Tunnel, rows int) error {
	if rows <= 0 {
		return fmt.Errorf("rows need to be greater than 0")
	}

	command := `
	DO $$ BEGIN
		IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'appscode_kubedb_postgres_test_table') THEN
			CREATE TABLE appscode_kubedb_postgres_test_table (value int not null);
		END IF;
	END $$;
	` + "\n" +
		fmt.Sprintf("INSERT INTO appscode_kubedb_postgres_test_table (value) values (generate_series(1,%v));", rows)

	out, err := opts.executeCommand(tunnel.Local, command)
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

			tunnel, err := lib.TunnelToDBService(opts.config, dbName, namespace, api.PostgresDatabasePort)
			if err != nil {
				log.Fatal("couldn't create tunnel, error: ", err)
			}

			err = opts.verifyDataExecCmd(tunnel, rows)
			if err != nil {
				log.Fatal(err)
			}
			tunnel.Close()
		},
	}
	pgVerifyCmd.Flags().IntVarP(&rows, "rows", "r", 100, "number of rows to verify")

	return pgVerifyCmd
}

func (opts *postgresOpts) verifyDataExecCmd(tunnel *portforward.Tunnel, rows int) error {
	if rows <= 0 {
		return fmt.Errorf("rows need to be greater than 0")
	}

	command := "SELECT COUNT(*) FROM appscode_kubedb_postgres_test_table;"
	out, err := opts.executeCommand(tunnel.Local, command)
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

			tunnel, err := lib.TunnelToDBService(opts.config, dbName, namespace, api.PostgresDatabasePort)
			if err != nil {
				log.Fatal("couldn't creat tunnel, error: ", err)
			}

			err = opts.dropDataExecCmd(tunnel)
			if err != nil {
				log.Fatal(err)
			}
			tunnel.Close()
		},
	}

	return pgDropCmd
}

func (opts *postgresOpts) dropDataExecCmd(tunnel *portforward.Tunnel) error {
	command := `
	DO $$ BEGIN
		IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'appscode_kubedb_postgres_test_table') THEN
			DROP TABLE appscode_kubedb_postgres_test_table;
		END IF;
	END $$;
    `

	_, err := opts.executeCommand(tunnel.Local, command)
	if err != nil {
		return err
	}
	fmt.Printf("\nSuccess: All the CLI inserted rows DELETED from postgres database %s/%s/\n", opts.db.Namespace, opts.db.Name)

	return nil
}

type postgresOpts struct {
	db       *api.Postgres
	dbImage  string
	config   *rest.Config
	client   *kubernetes.Clientset
	dbClient *cs.Clientset

	username string
	pass     string

	errWriter *bytes.Buffer
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

	dbVersion, err := dbClient.CatalogV1alpha1().PostgresVersions().Get(context.TODO(), db.Spec.Version, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	secret, err := client.CoreV1().Secrets(db.Namespace).Get(context.TODO(), db.Spec.AuthSecret.Name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return &postgresOpts{
		db:        db,
		dbImage:   dbVersion.Spec.DB.Image,
		config:    config,
		client:    client,
		dbClient:  dbClient,
		username:  string(secret.Data[corev1.BasicAuthUsernameKey]),
		pass:      string(secret.Data[corev1.BasicAuthPasswordKey]),
		errWriter: &bytes.Buffer{},
	}, nil
}

func (opts *postgresOpts) getDockerShellCommand(localPort int, dockerFlags, postgresExtraFlags []interface{}) (*shell.Session, error) {
	sh := shell.NewSession()
	sh.ShowCMD = false
	sh.Stderr = opts.errWriter

	db := opts.db
	dockerCommand := []interface{}{
		"run", "--network=host",
		"-e", fmt.Sprintf("PGPASSWORD=%s", opts.pass),
	}
	dockerCommand = append(dockerCommand, dockerFlags...)

	postgresCommand := []interface{}{
		"psql",
		"--host=127.0.0.1", fmt.Sprintf("--port=%d", localPort),
		fmt.Sprintf("--username=%s", opts.username),
	}

	if db.Spec.TLS != nil {
		secretName := db.CertificateName(api.PostgresClientCert)
		certSecret, err := opts.client.CoreV1().Secrets(db.Namespace).Get(context.TODO(), secretName, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}

		caCrt, ok := certSecret.Data[corev1.ServiceAccountRootCAKey]
		if !ok {
			return nil, fmt.Errorf("missing %s in secret %s/%s", corev1.ServiceAccountRootCAKey, certSecret.Namespace, certSecret.Name)
		}
		err = os.WriteFile(pgCaFile, caCrt, 0o644)
		if err != nil {
			return nil, err
		}

		crt, ok := certSecret.Data[corev1.TLSCertKey]
		if !ok {
			return nil, fmt.Errorf("missing %s in secret %s/%s", corev1.TLSCertKey, certSecret.Namespace, certSecret.Name)
		}
		err = os.WriteFile(pgCertFile, crt, 0o644)
		if err != nil {
			return nil, err
		}

		key, ok := certSecret.Data[corev1.TLSPrivateKeyKey]
		if !ok {
			return nil, fmt.Errorf("missing %s in secret %s/%s", corev1.TLSPrivateKeyKey, certSecret.Namespace, certSecret.Name)
		}
		err = os.WriteFile(pgKeyFile, key, 0o600)
		if err != nil {
			return nil, err
		}

		dockerCommand = append(dockerCommand,
			"-v", fmt.Sprintf("%s:%s", "/tmp/", "/root/.postgresql/"),
		)
	}

	dockerCommand = append(dockerCommand, opts.dbImage)
	finalCommand := append(dockerCommand, postgresCommand...)
	if postgresExtraFlags != nil {
		finalCommand = append(finalCommand, postgresExtraFlags...)
	}
	return sh.Command("docker", finalCommand...).SetStdin(os.Stdin), nil
}

func (opts *postgresOpts) executeCommand(localPort int, command string) (string, error) {
	postgresExtraFlags := []interface{}{
		fmt.Sprintf("--command=%s", command),
	}
	shSession, err := opts.getDockerShellCommand(localPort, nil, postgresExtraFlags)
	if err != nil {
		return "", err
	}
	out, err := shSession.Output()
	if err != nil {
		return "", fmt.Errorf("failed to execute command, error: %s, output: %s\n", err, out)
	}
	output := ""
	if string(out) != "" {
		output = string(out)
	}
	errOutput := opts.errWriter.String()
	if errOutput != "" {
		return "", fmt.Errorf("failed to execute command, stderr: %s%s", errOutput, output)
	}

	return string(out), nil
}
