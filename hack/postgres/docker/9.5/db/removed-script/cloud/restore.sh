#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

source /scripts/lib.sh

if [ ! -f "$PGDATA/.appscode" ]; then
	initdb
	cp $PGDATA/postgresql.conf /tmp/postgresql.conf
	cp $PGDATA/pg_hba.conf /tmp/pg_hba.conf
	rm -rf $PGDATA/*

	set_wal-e_credential

	envdir /etc/wal-e.d/env /usr/local/bin/wal-e backup-list

	###################################################################
	### recovery.conf file will restore WAL file and make it master ###
	gosu postgres touch $PGDATA/recovery.conf
	{ echo; echo "restore_command  = 'envdir /etc/wal-e.d/env wal-e wal-fetch "%f" "%p"'"; } >> "$PGDATA/recovery.conf"
	{ echo; echo "trigger_file = '/tmp/postgresql.trigger'"; } >> "$PGDATA/recovery.conf"
	###################################################################

	## To restore base, works like pg_basebackup ##
	gosu postgres envdir /etc/wal-e.d/env wal-e backup-fetch $PGDATA LATEST

	cp /tmp/postgresql.conf $PGDATA/postgresql.conf
	cp /tmp/pg_hba.conf $PGDATA/pg_hba.conf

	config_postgresql 'cloud'
	config_pg_hba 'cloud'
	###################################################################

	## To make it master ##
	touch /tmp/postgresql.trigger

	chown -R postgres:postgres "$PGDATA"

	########## Backup base ##########
	gosu postgres pg_ctl -D $PGDATA  -w start
	backup-push
	gosu postgres pg_ctl -D $PGDATA -m fast -w stop
	#################################

	echo
	echo 'PostgreSQL init process complete; ready for start up.'
	echo

	touch $PGDATA/.appscode
	chmod 700 $PGDATA
	chown -R postgres:postgres $PGDATA
	touch /tmp/.done

fi
exec gosu postgres postgres
