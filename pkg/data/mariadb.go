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
	"errors"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"

	api "kubedb.dev/apimachinery/apis/kubedb/v1alpha2"
	cs "kubedb.dev/apimachinery/client/clientset/versioned"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
)

func InsertMariaDBDataCMD(f cmdutil.Factory) *cobra.Command {
	var (
		dbName string
		rows   int
	)

	mdInsertCmd := &cobra.Command{
		Use: "mariadb",
		Aliases: []string{
			"md",
		},
		Short:   " Insert data to mariadb",
		Long:    `Use this cmd to insert data into a mariadb database.`,
		Example: `kubectl dba data insert mariadb -n demo sample-mariadb --rows 1000`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				log.Fatal("Enter mariadb object's name as an argument")
			}
			dbName = args[0]

			namespace, _, err := f.ToRawKubeConfigLoader().Namespace()
			if err != nil {
				klog.Error(err, "failed to get current namespace")
			}

			opts, err := newMariaDBOpts(f, dbName, namespace)
			if err != nil {
				log.Fatalln(err)
			}

			if rows <= 0 {
				log.Fatal("Inserted rows must be greater than 0")
			}

			if rows > 100000 {
				log.Fatal("Inserted rows must be less than or equal 100000")
			}

			err = opts.insertDataExecCmd(rows)
			if err != nil {
				log.Fatal(err)
			}
		},
	}

	mdInsertCmd.Flags().IntVarP(&rows, "rows", "r", 100, "number of rows to insert")

	return mdInsertCmd
}

func (opts *mariadbOpts) insertDataExecCmd(rows int) error {
	command := fmt.Sprintf(`
		CREATE DATABASE IF NOT EXISTS %v;
		USE %v;
		CREATE TABLE IF NOT EXISTS %v (id VARCHAR(255) PRIMARY KEY);
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
				WHILE i < 15 DO
					SET result = CONCAT(result, SUBSTRING(characters, FLOOR(RAND() * 81) + 1, 1));
					SET i = i + 1;
				END WHILE;
				INSERT INTO kubedb_table (id) VALUES (result ); 
				SET counter = counter + 1;
			END WHILE;
		END //
		DELIMITER ;
		CALL insert_data(`+fmt.Sprintf("%v", rows)+`); 
	`, KubeDBDatabaseName, KubeDBDatabaseName, KubeDBTableName)

	_, err := opts.executeCommand(command)
	if err != nil {
		return err
	}

	fmt.Printf("\nSuccess! %d keys inserted in MariaDB database %s/%s.\n", rows, opts.db.Namespace, opts.db.Name)
	return nil
}

func VerifyMariaDBDataCMD(f cmdutil.Factory) *cobra.Command {
	var (
		dbName string
		rows   int
	)

	mdVerifyCmd := &cobra.Command{
		Use: "mariadb",
		Aliases: []string{
			"md",
		},
		Short:   "Verify rows to a mariadb resource",
		Long:    `Use this cmd to verify data in a mariadb object`,
		Example: `kubectl dba data verify mariadb -n demo sample-mariadb --rows 1000`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				log.Fatal("Enter mariadb object's name as an argument.")
			}
			dbName = args[0]

			namespace, _, err := f.ToRawKubeConfigLoader().Namespace()
			if err != nil {
				klog.Error(err, "failed to get current namespace")
			}

			opts, err := newMariaDBOpts(f, dbName, namespace)
			if err != nil {
				log.Fatalln(err)
			}

			err = opts.verifyDataExecCmd(rows)
			if err != nil {
				log.Fatal(err)
			}
		},
	}

	mdVerifyCmd.Flags().IntVarP(&rows, "rows", "r", 100, "number of rows to verify")

	return mdVerifyCmd
}

func (opts *mariadbOpts) verifyDataExecCmd(rows int) error {
	if rows <= 0 {
		return fmt.Errorf("rows need to be greater than 0")
	}

	command := fmt.Sprintf(` 
        CREATE DATABASE IF NOT EXISTS %v;
        USE %v;
        CREATE TABLE IF NOT EXISTS kubedb_table (id VARCHAR(255) PRIMARY KEY);
        SELECT COUNT(*) FROM kubedb_table; 
    `, KubeDBDatabaseName, KubeDBDatabaseName)

	o, err := opts.executeCommand(command)
	if err != nil {
		return err
	}

	out := string(o)
	output := strings.Split(out, "\n")

	totalKeys, err := strconv.Atoi(strings.TrimSpace(output[1]))
	if err != nil {
		return err
	}
	if totalKeys >= rows {
		fmt.Printf("\nSuccess! MariaDB database %s/%s contains: %d keys\n", opts.db.Namespace, opts.db.Name, totalKeys)
	} else {
		fmt.Printf("\nError! Expected keys: %d . MariaDB database %s/%s contains: %d keys\n", rows, opts.db.Namespace, opts.db.Name, totalKeys)
	}
	return nil
}

func DropMariaDBDataCMD(f cmdutil.Factory) *cobra.Command {
	var dbName string

	mdDropCmd := &cobra.Command{
		Use: "mariadb",
		Aliases: []string{
			"md",
		},
		Short:   "Drop data from MariaDB",
		Long:    `Use this cmd to drop data from a mongodb`,
		Example: `kubectl dba data drop mariadb -n demo sample-mariadb`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				log.Fatal("Enter mariadb object's name as an argument.")
			}
			dbName = args[0]

			namespace, _, err := f.ToRawKubeConfigLoader().Namespace()
			if err != nil {
				klog.Error(err, "failed to get current namespace")
			}

			opts, err := newMariaDBOpts(f, dbName, namespace)
			if err != nil {
				log.Fatalln(err)
			}

			err = opts.dropDataExecCmd()
			if err != nil {
				log.Fatal(err)
			}
		},
	}

	return mdDropCmd
}

func (opts *mariadbOpts) dropDataExecCmd() error {
	command := fmt.Sprintf(` 
		USE %v;
		DROP TABLE IF EXISTS %v;
	`, KubeDBDatabaseName, KubeDBTableName)
	_, err := opts.executeCommand(command)
	if err != nil {
		return err
	}
	fmt.Printf("\nSuccess: All the CLI inserted rows DELETED from MariaDB database %s/%s.\n", opts.db.Namespace, opts.db.Name)

	return nil
}

type mariadbOpts struct {
	db        *api.MariaDB
	dbImage   string
	config    *rest.Config
	client    *kubernetes.Clientset
	username  string
	pass      string
	errWriter *bytes.Buffer
}

func newMariaDBOpts(f cmdutil.Factory, dbName, namespace string) (*mariadbOpts, error) {
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

	db, err := dbClient.KubedbV1alpha2().MariaDBs(namespace).Get(context.TODO(), dbName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	if db.Status.Phase != api.DatabasePhaseReady {
		return nil, fmt.Errorf("mariadb %s/%s is not ready", namespace, dbName)
	}

	dbVersion, err := dbClient.CatalogV1alpha1().MariaDBVersions().Get(context.TODO(), db.Spec.Version, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	secret, err := client.CoreV1().Secrets(db.Namespace).Get(context.TODO(), db.Spec.AuthSecret.Name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return &mariadbOpts{
		db:        db,
		dbImage:   dbVersion.Spec.DB.Image,
		config:    config,
		client:    client,
		username:  string(secret.Data[corev1.BasicAuthUsernameKey]),
		pass:      string(secret.Data[corev1.BasicAuthPasswordKey]),
		errWriter: &bytes.Buffer{},
	}, nil
}

func (opts *mariadbOpts) getShellCommand(command string) (string, error) {
	db := opts.db
	cmd := ""
	user, password, err := opts.GetMariaDBAuthCredentials(db)
	if err != nil {
		return "", err
	}

	containerName := "mariadb"

	if db.Spec.TLS != nil {
		cmd = fmt.Sprintf("kubectl exec -n %s svc/%s -c %s -- mysql -u%s -p'%s' --host=%s --port=%s --ssl-ca='%v' --ssl-cert='%v' --ssl-key='%v' %s -e \"%s\"", db.Namespace, db.OffshootName(), containerName, user, password, "127.0.0.1", "3306", myCaFile, myCertFile, myKeyFile, api.ResourceSingularMySQL, command)
	} else {
		cmd = fmt.Sprintf("kubectl exec -n %s svc/%s -c %s -- mysql -u%s -p'%s' %s -e \"%s\"", db.Namespace, db.OffshootName(), containerName, user, password, api.ResourceSingularMySQL, command)
	}

	return cmd, err
}

func (opts *mariadbOpts) GetMariaDBAuthCredentials(db *api.MariaDB) (string, string, error) {
	if db.Spec.AuthSecret == nil {
		return "", "", errors.New("no database secret")
	}
	secret, err := opts.client.CoreV1().Secrets(db.Namespace).Get(context.TODO(), db.Spec.AuthSecret.Name, metav1.GetOptions{})
	if err != nil {
		return "", "", err
	}
	return string(secret.Data[corev1.BasicAuthUsernameKey]), string(secret.Data[corev1.BasicAuthPasswordKey]), nil
}

func (opts *mariadbOpts) executeCommand(command string) ([]byte, error) {
	cmd, err := opts.getShellCommand(command)
	if err != nil {
		return nil, err
	}
	output, err := opts.runCMD(cmd)
	if err != nil {
		return nil, err
	}
	return output, nil
}

func (opts *mariadbOpts) runCMD(cmd string) ([]byte, error) {
	sh := exec.Command("/bin/sh", "-c", cmd)
	stdout := bytes.NewBuffer(nil)
	stderr := bytes.NewBuffer(nil)
	sh.Stdout = stdout
	sh.Stderr = stderr
	err := sh.Run()
	out := stdout.Bytes()
	errOut := stderr.Bytes()
	errOutput := string(errOut)
	if errOutput != "" && !strings.Contains(errOutput, "NOTICE") && !strings.Contains(errOutput, "Warning") {
		return nil, fmt.Errorf("failed to execute command, stderr: %s", errOutput)
	}
	if err != nil {
		return nil, err
	}
	return out, nil
}
