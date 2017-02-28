#!/bin/bash

set -e

set_listen_addresses() {
	sedEscapedValue="$(echo "$1" | sed 's/[\/&]/\\&/g')"
	sed -ri "s/^#?(listen_addresses\s*=\s*)\S+/\1'$sedEscapedValue'/" "$PGDATA/postgresql.conf"
}

if [ "$1" = 'postgres' ]; then

	mkdir -p "$PGDATA"
	chown -R postgres:postgres "$PGDATA"
	chmod g+s /run/postgresql
	chown -R postgres /run/postgresql

	if [ ! -f "$PGDATA/.appscode" ]; then

		rm -rf $PGDATA/*

		gosu postgres pg_basebackup -x -P -R -D $PGDATA  -h $MASTER -U rep -v

		cat <<-EOF >> "$PGDATA/postgresql.conf"
				hot_standby = on
		EOF

		sed -ri "s/krbsrvname=postgres/application_name=$HOSTPOD/" "$PGDATA/recovery.conf"

		{ echo; echo "trigger_file = '/tmp/postgresql.trigger'"; } >> "$PGDATA/recovery.conf"
		{ echo; echo "recovery_target_timeline='latest'"; } >> "$PGDATA/recovery.conf"

		echo
		echo 'PostgreSQL init process complete; ready for start up.'
		echo

		touch $PGDATA/.appscode
		chmod 700 $PGDATA
		chown -R postgres:postgres $PGDATA

	fi
	exec gosu postgres "$@"
fi
exec "$@"