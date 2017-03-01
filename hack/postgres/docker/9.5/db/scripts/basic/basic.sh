#!/bin/bash

set -e

source /scripts/lib.sh
user_script=${1:-}

if [ ! -f "$PGDATA/.appscode" ]; then
  # The VERY FIRST run
	initdb
	config_postgresql 'basic'
	config_pg_hba 'basic'
	configure_pgpool
	set_password 'basic'

	if [ -n  "$user_script" ]; then
		run_user_script "$user_script"
	fi

	gosu postgres touch $PGDATA/.appscode
	touch /tmp/.done
else
	reset_owner
	set_password 'basic'
	touch /tmp/.done
fi

exec gosu postgres postgres
