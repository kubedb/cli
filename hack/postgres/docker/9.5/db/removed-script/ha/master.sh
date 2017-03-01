#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

source /scripts/lib.sh
user_script=${1:-}

if [ ! -f "$PGDATA/.appscode" ]; then
	initdb
	config_postgresql 'ha'
	config_pg_hba 'ha'
	configure_pgpool
	set_password 'ha'

	if [ -n  "$user_script" ]; then
		run_user_script "$user_script"
	fi

	gosu postgres touch $PGDATA/.appscode
	touch /tmp/.done
fi

exec gosu postgres postgres

# select * from pg_stat_replication ;
