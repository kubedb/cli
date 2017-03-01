#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

source /scripts/lib.sh

if [ ! -f "$PGDATA/.appscode" ]; then

	init_basebackup

	cat <<-EOF >> "$PGDATA/postgresql.conf"
hot_standby = on
EOF

	sed -ri "s/krbsrvname=postgres/application_name=$HOSTNAME/" "$PGDATA/recovery.conf"

	{ echo; echo "trigger_file = '/tmp/postgresql.trigger'"; } >> "$PGDATA/recovery.conf"
	{ echo; echo "recovery_target_timeline='latest'"; } >> "$PGDATA/recovery.conf"

	echo
	echo 'PostgreSQL init process complete; ready for start up.'
	echo

	touch $PGDATA/.appscode
	chmod 700 $PGDATA
	chown -R postgres:postgres $PGDATA
	touch /tmp/.done
fi

exec gosu postgres postgres
