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

package cmds

import (
	"kubedb.dev/cli/pkg/data"

	"github.com/spf13/cobra"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/util/i18n"
	"k8s.io/kubectl/pkg/util/templates"
)

var (
	dataLong = templates.LongDesc(`
		Insert random data or verify data in a database.
    `)
	dataExample = templates.Examples(`
	    # Insert 100 rows in mysql table
		kubectl dba data insert mysql mysql-demo -n demo --rows=100

		# Verify that postgres has at least 100 rows data
		kubectl dba data verify postgres sample-postgres -n demo --rows=100

		# Drop all the CLI inserted data from mongodb
		kubectl dba data drop mg -n demo sample-mg
		

 		Valid resource types include:
    		* elasticsearch
			* mongodb
			* mariadb
			* mysql
			* postgres
			* redis
`)
)

func NewCmdData(f cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "data",
		Short:                 i18n.T("Insert, Drop or Verify data in a database"),
		Long:                  dataLong,
		Example:               dataExample,
		Run:                   func(cmd *cobra.Command, args []string) {},
		DisableFlagsInUseLine: true,
		DisableAutoGenTag:     true,
	}

	cmd.AddCommand(InsertDataCMD(f))
	cmd.AddCommand(VerifyDataCMD(f))
	cmd.AddCommand(DropDataCMD(f))

	return cmd
}

var insertLong = templates.LongDesc(`
		Insert random data in a database.
    `)

var insertExample = templates.Examples(`
		# Insert 100 rows in mysql table
		kubectl dba data insert mysql mysql-demo -n demo --rows=100

		#Insert 100 keys in redis
		kubectl dba data insert rd sample-redis -n demo --rows=100

 		Valid resource types include:
    		* elasticsearch
			* mongodb
			* mariadb
			* mysql
			* postgres
			* redis
`)

func InsertDataCMD(f cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "insert",
		Short: i18n.T("Insert random data in a database"),
		Long:  insertLong,

		Run:                   func(cmd *cobra.Command, args []string) {},
		Example:               insertExample,
		DisableFlagsInUseLine: true,
		DisableAutoGenTag:     true,
	}

	cmd.AddCommand(data.InsertRedisDataCMD(f))

	return cmd
}

var verifyLong = templates.LongDesc(`
		Verify data in a database.
    `)

var verifyExample = templates.Examples(`
		# Verify if there is 100 rows a postgres table
		kubectl dba data verify postgres sample-pg -n demo --rows=100

		# Verify if there is 100 rows in a mongodb database
		kubectl dba data verify mongodb -n demo mg-shard --rows=100

 		Valid resource types include:
    		* elasticsearch
			* mongodb
			* mariadb
			* mysql
			* postgres
			* redis
`)

func VerifyDataCMD(f cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "verify",
		Short:                 i18n.T("Verify data in a database"),
		Long:                  verifyLong,
		Example:               verifyExample,
		Run:                   func(cmd *cobra.Command, args []string) {},
		DisableFlagsInUseLine: true,
		DisableAutoGenTag:     true,
	}

	cmd.AddCommand(data.VerifyRedisDataCMD(f))

	return cmd
}

var dropLong = templates.LongDesc(`
		Drop data in a database.
    `)

var dropExample = templates.Examples(`
		# Drop all the cli inserted data from mariadb
		kubectl dba data drop mariadb -n demo sample-maria

		# Drop all the cli inserted data from elasticsearch
		kubectl dba data drop es -n demo 

 		Valid resource types include:
    		* elasticsearch
			* mongodb
			* mariadb
			* mysql
			* postgres
			* redis
`)

func DropDataCMD(f cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "drop",
		Short:                 i18n.T("Drop data from a database"),
		Long:                  dropLong,
		Example:               dropExample,
		Run:                   func(cmd *cobra.Command, args []string) {},
		DisableFlagsInUseLine: true,
		DisableAutoGenTag:     true,
	}

	cmd.AddCommand(data.DropRedisDataCMD(f))

	return cmd
}
