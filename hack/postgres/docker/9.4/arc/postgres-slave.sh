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

		##############################################################
		################ Setting up Cloud Bucket ####################
		umask u=rwx,g=rx,o=
		mkdir -p /etc/wal-e.d/env

		if [ -s "/var/credentials/gce" ]; then
			echo 'gs://'$BUCKET'/backup' > /etc/wal-e.d/env/WALE_GS_PREFIX
			echo '/var/credentials/gce' > /etc/wal-e.d/env/GOOGLE_APPLICATION_CREDENTIALS
			echo $PROJECT > /etc/wal-e.d/env/GCLOUD_PROJECT
		fi

		if [ -d "/var/credentials/aws" ]; then
			echo 's3://'$BUCKET'/backup' > /etc/wal-e.d/env/WALE_S3_PREFIX
			cat /var/credentials/aws/secret > /etc/wal-e.d/env/AWS_SECRET_ACCESS_KEY
			cat /var/credentials/aws/keyid > /etc/wal-e.d/env/AWS_ACCESS_KEY_ID
			cat /var/credentials/aws/region > /etc/wal-e.d/env/AWS_REGION
		fi
		chown -R root:postgres /etc/wal-e.d
		##############################################################

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
