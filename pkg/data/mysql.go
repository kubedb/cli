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

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/cobra"
	shell "gomodules.xyz/go-sh"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"kmodules.xyz/client-go/tools/portforward"
)

const (
	caFile   = "/tmp/ca.crt"
	certFile = "/tmp/client.crt"
	keyFile  = "/tmp/client.key"
)

func InsertMySQLDataCMD(f cmdutil.Factory) *cobra.Command {
	var (
		dbName string
		rows   int
	)

	myInsertCmd := &cobra.Command{
		Use: "mysql",
		Aliases: []string{
			"my",
		},
		Short:   "Connect to a mysql object",
		Long:    `Use this cmd to exec into a mysql object's primary pod.`,
		Example: `kubectl dba insert mysql -n demo sample-mysql --rows 1000`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				log.Fatal("Enter mysql object's name as an argument")
			}
			dbName = args[0]

			namespace, _, err := f.ToRawKubeConfigLoader().Namespace()
			if err != nil {
				klog.Error(err, "failed to get current namespace")
			}

			opts, err := newMySQLOpts(f, dbName, namespace)
			if err != nil {
				log.Fatalln(err)
			}

			if rows <= 0 {
				log.Fatal("Inserted rows must be greater than 0")
			}

			tunnel, err := lib.TunnelToDBService(opts.config, dbName, namespace, api.MySQLDatabasePort)
			if err != nil {
				log.Fatal("couldn't create tunnel, error: ", err)
			}

			defer tunnel.Close()

			err = opts.insertDataExecCmd(tunnel, rows)
			if err != nil {
				log.Fatal(err)
			}
		},
	}

	myInsertCmd.Flags().IntVarP(&rows, "rows", "r", 10, "rows in ")

	return myInsertCmd
}

func (opts *mysqlOpts) insertDataExecCmd(tunnel *portforward.Tunnel, rows int) error {
	command := `
		USE mysql;
		CREATE TABLE IF NOT EXISTS kubedb_table (id VARCHAR(255) PRIMARY KEY);
		DROP PROCEDURE IF EXISTS insert_data;
		DELIMITER //
		CREATE PROCEDURE insert_data(max_value INT)
		BEGIN
			DECLARE counter INT DEFAULT 1;
			DECLARE characters VARCHAR(82) DEFAULT 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789@#$%^&*()!';
			DECLARE result VARCHAR(255) DEFAULT '';
			DECLARE i INT DEFAULT 0;
			WHILE counter <= max_value DO
				SET result = '';
				SET i = 0;
				WHILE i < 10 DO
					SET result = CONCAT(result, SUBSTRING(characters, FLOOR(RAND() * 81) + 1, 1));
					SET i = i + 1;
				END WHILE;
				INSERT INTO kubedb_table (id) VALUES (result ); 
				SET counter = counter + 1;
			END WHILE;
		END //
		DELIMITER ;
		CALL insert_data(` + fmt.Sprintf("%v", rows) + `); 
	`

	_, err := opts.executeCommand(tunnel.Local, command)
	if err != nil {
		return err
	}

	fmt.Printf("\nSuccess! %d keys inserted in mysql database %s/%s.\n", rows, opts.db.Namespace, opts.db.Name)
	return nil
}

func VerifyMySQLDataCMD(f cmdutil.Factory) *cobra.Command {
	var (
		dbName string
		rows   int
	)

	myVerifyCmd := &cobra.Command{
		Use: "mysql",
		Aliases: []string{
			"my",
		},
		Short:   "Verify rows in a MySQL database",
		Long:    `Use this cmd to verify data in a mysql object`,
		Example: `kubectl dba verify mysql -n demo sample-mysql --rows 1000`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				log.Fatal("Enter mysql object's name as an argument.")
			}
			dbName = args[0]

			namespace, _, err := f.ToRawKubeConfigLoader().Namespace()
			if err != nil {
				klog.Error(err, "failed to get current namespace")
			}

			opts, err := newMySQLOpts(f, dbName, namespace)
			if err != nil {
				log.Fatalln(err)
			}

			tunnel, err := lib.TunnelToDBService(opts.config, dbName, namespace, api.MySQLDatabasePort)
			if err != nil {
				log.Fatal("couldn't create tunnel, error: ", err)
			}
			defer tunnel.Close()

			err = opts.verifyDataExecCmd(tunnel, rows)
			if err != nil {
				log.Fatal(err)
			}
		},
	}

	myVerifyCmd.Flags().IntVarP(&rows, "rows", "r", 10, "rows in ")

	return myVerifyCmd
}

func (opts *mysqlOpts) verifyDataExecCmd(tunnel *portforward.Tunnel, rows int) error {
	if rows <= 0 {
		return fmt.Errorf("rows need to be greater than 0")
	}

	command := ` 
		USE mysql;
		SELECT COUNT(*) FROM kubedb_table; 
	`
	out, err := opts.executeCommand(tunnel.Local, command)
	if err != nil {
		return err
	}

	output := strings.Split(out, "\n")

	totalKeys, err := strconv.Atoi(strings.TrimPrefix(output[1], " "))
	if err != nil {
		return err
	}
	if totalKeys >= rows {
		fmt.Printf("\nSuccess! MySQL database %s/%s contains: %d keys\n", opts.db.Namespace, opts.db.Name, totalKeys)
	} else {
		fmt.Printf("\nError! Expected keys: %d . MySQL database %s/%s contains: %d keys\n", rows, opts.db.Namespace, opts.db.Name, totalKeys)
	}
	return nil
}

func DropMySQLDataCMD(f cmdutil.Factory) *cobra.Command {
	var dbName string

	myDropCmd := &cobra.Command{
		Use: "mysql",
		Aliases: []string{
			"my",
		},
		Short:   "Verify rows in a MySQL database",
		Long:    `Use this cmd to verify data in a mysql object`,
		Example: `kubectl dba drop mysql -n demo sample-mysql`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				log.Fatal("Enter mysql object's name as an argument.")
			}
			dbName = args[0]

			namespace, _, err := f.ToRawKubeConfigLoader().Namespace()
			if err != nil {
				klog.Error(err, "failed to get current namespace")
			}

			opts, err := newMySQLOpts(f, dbName, namespace)
			if err != nil {
				log.Fatalln(err)
			}

			tunnel, err := lib.TunnelToDBService(opts.config, dbName, namespace, api.MySQLDatabasePort)
			if err != nil {
				log.Fatal("couldn't create tunnel, error: ", err)
			}
			defer tunnel.Close()

			err = opts.dropDataExecCmd(tunnel)
			if err != nil {
				log.Fatal(err)
			}
		},
	}

	return myDropCmd
}

func (opts *mysqlOpts) dropDataExecCmd(tunnel *portforward.Tunnel) error {
	command := ` 
		USE mysql;
		DROP TABLE IF EXISTS kubedb_table;
	`
	_, err := opts.executeCommand(tunnel.Local, command)
	if err != nil {
		return err
	}
	fmt.Printf("\nSuccess: All the CLI inserted rows DELETED from MySQL database %s/%s.\n", opts.db.Namespace, opts.db.Name)

	return nil
}

type mysqlOpts struct {
	db        *api.MySQL
	dbImage   string
	config    *rest.Config
	client    *kubernetes.Clientset
	username  string
	pass      string
	errWriter *bytes.Buffer
}

func newMySQLOpts(f cmdutil.Factory, dbName, namespace string) (*mysqlOpts, error) {
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

	db, err := dbClient.KubedbV1alpha2().MySQLs(namespace).Get(context.TODO(), dbName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	if db.Status.Phase != api.DatabasePhaseReady {
		return nil, fmt.Errorf("mysql %s/%s is not ready", namespace, dbName)
	}

	dbVersion, err := dbClient.CatalogV1alpha1().MySQLVersions().Get(context.TODO(), db.Spec.Version, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	secret, err := client.CoreV1().Secrets(db.Namespace).Get(context.TODO(), db.Spec.AuthSecret.Name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return &mysqlOpts{
		db:        db,
		dbImage:   dbVersion.Spec.DB.Image,
		config:    config,
		client:    client,
		username:  string(secret.Data[corev1.BasicAuthUsernameKey]),
		pass:      string(secret.Data[corev1.BasicAuthPasswordKey]),
		errWriter: &bytes.Buffer{},
	}, nil
}

func (opts *mysqlOpts) getDockerShellCommand(localPort int, dockerFlags, mysqlExtraFlags []interface{}) (*shell.Session, error) {
	sh := shell.NewSession()
	sh.ShowCMD = false

	db := opts.db
	dockerCommand := []interface{}{
		"run", "--network=host",
		"-e", fmt.Sprintf("MYSQL_PWD=%s", opts.pass),
	}
	dockerCommand = append(dockerCommand, dockerFlags...)

	mysqlCommand := []interface{}{
		"mysql",
		"--host=127.0.0.1", fmt.Sprintf("--port=%d", localPort),
		fmt.Sprintf("--user=%s", opts.username),
	}

	if db.Spec.TLS != nil {
		secretName := db.CertificateName(api.MySQLClientCert)
		certSecret, err := opts.client.CoreV1().Secrets(db.Namespace).Get(context.TODO(), secretName, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}

		caCrt, ok := certSecret.Data[corev1.ServiceAccountRootCAKey]
		if !ok {
			return nil, fmt.Errorf("missing %s in secret %s/%s", corev1.ServiceAccountRootCAKey, certSecret.Namespace, certSecret.Name)
		}
		err = os.WriteFile(caFile, caCrt, 0o644)
		if err != nil {
			return nil, err
		}

		crt, ok := certSecret.Data[corev1.TLSCertKey]
		if !ok {
			return nil, fmt.Errorf("missing %s in secret %s/%s", corev1.TLSCertKey, certSecret.Namespace, certSecret.Name)
		}
		err = os.WriteFile(certFile, crt, 0o644)
		if err != nil {
			return nil, err
		}

		key, ok := certSecret.Data[corev1.TLSPrivateKeyKey]
		if !ok {
			return nil, fmt.Errorf("missing %s in secret %s/%s", corev1.TLSPrivateKeyKey, certSecret.Namespace, certSecret.Name)
		}
		err = os.WriteFile(keyFile, key, 0o644)
		if err != nil {
			return nil, err
		}

		dockerCommand = append(dockerCommand,
			"-v", fmt.Sprintf("%s:%s", caFile, caFile),
			"-v", fmt.Sprintf("%s:%s", certFile, certFile),
			"-v", fmt.Sprintf("%s:%s", keyFile, keyFile),
		)
		mysqlCommand = append(mysqlCommand,
			fmt.Sprintf("--ssl-ca=%v", caFile),
			fmt.Sprintf("--ssl-cert=%v", certFile),
			fmt.Sprintf("--ssl-key=%v", keyFile),
		)
	}

	dockerCommand = append(dockerCommand, opts.dbImage)
	finalCommand := append(dockerCommand, mysqlCommand...)
	if mysqlExtraFlags != nil {
		finalCommand = append(finalCommand, mysqlExtraFlags...)
	}
	return sh.Command("docker", finalCommand...).SetStdin(os.Stdin), nil
}

func (opts *mysqlOpts) executeCommand(localPort int, command string) (string, error) {
	mysqlExtraFlags := []interface{}{
		"-e", command,
	}

	shSession, err := opts.getDockerShellCommand(localPort, nil, mysqlExtraFlags)
	if err != nil {
		return "", err
	}

	out, err := shSession.Output()
	if err != nil {
		return "", fmt.Errorf("failed to execute file, error: %s, output: %s\n", err, out)
	}

	output := ""
	if string(out) != "" {
		output = ", output:\n\n" + string(out)
	}

	errOutput := opts.errWriter.String()
	if errOutput != "" {
		return "", fmt.Errorf("failed to execute command, stderr: %s%s", errOutput, output)
	}

	return string(out), nil
}
