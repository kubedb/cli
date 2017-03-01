#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

source /scripts/lib.sh
user_script=${1:-}

if [ ! -f "$PGDATA/.appscode" ]; then

	##############################################################
	initdb
	set_wal-e_credential
	config_postgresql 'cloud'
	config_pg_hba 'cloud'
	configure_pgpool
	set_password 'cloud'

	########## Backup base ##########
	gosu postgres pg_ctl -D $PGDATA  -w start
	backup-push
	gosu postgres pg_ctl -D $PGDATA -m fast -w stop
	#################################

	if [ -n  "$user_script" ]; then
		run_user_script "$user_script"
	fi

	gosu postgres touch $PGDATA/.appscode
	touch /tmp/.done
fi
exec gosu postgres postgres
