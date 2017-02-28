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
		gosu postgres initdb $PGDATA
		cp $PGDATA/postgresql.conf /tmp/postgresql.conf
		cp $PGDATA/pg_hba.conf /tmp/pg_hba.conf
		rm -rf $PGDATA/*

		##############################################################
		################ Setting up Cloud Bucket ####################
		umask u=rwx,g=rx,o=
		mkdir -p /etc/wal-e.d/env

		if [ -f "/var/credentials/gce" ]; then
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

		###################################################################
		########  modify postgresql.conf &  pg_hba.conf as master #########
		cat <<-EOF >> "$PGDATA/postgresql.conf"
			listen_addresses = '*'
			wal_level = 'hot_standby'
			max_wal_senders = $MW_SENDER
			wal_keep_segments = $WK_SEGMENTS

			###########################################
			##  For Wal Archive in external server ####
			archive_mode = on
			archive_command = 'envdir /etc/wal-e.d/env /usr/local/bin/wal-e wal-push %p'
			archive_timeout =  $ARC_TIME
			###########################################

			ssl = true
			ssl_cert_file = '/etc/ssl/certs/ssl-cert-snakeoil.pem'
			ssl_key_file = '/etc/ssl/private/ssl-cert-snakeoil.key'
		EOF

		{ echo; echo "host all all 0.0.0.0/0 md5"; } >> "$PGDATA/pg_hba.conf"
		{ echo; echo "hostssl replication rep 0.0.0.0/0 trust"; } >> "$PGDATA/pg_hba.conf"
		###################################################################

		## To make it master ##
		touch /tmp/postgresql.trigger

		chown -R postgres:postgres "$PGDATA"

		cat <<-EOF >> "$PGDATA/base_backup.sh"
			#!/bin/bash
			gosu postgres envdir /etc/wal-e.d/env /usr/local/bin/wal-e backup-push $PGDATA
		EOF
		chown -R postgres:postgres $PGDATA/base_backup.sh

		
		#################################
		########## Backup base ##########
		gosu postgres pg_ctl -D "$PGDATA" -o "-c listen_addresses=''" -w start
		/bin/bash $PGDATA/base_backup.sh
		gosu postgres pg_ctl -D "$PGDATA" -m fast -w stop
		#################################

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
