#!/usr/bin/env bash

set -o nounset
set -o pipefail

reset_owner() {
	chown -R postgres:postgres "$PGDATA"
	chmod g+s /run/postgresql
	chown -R postgres /run/postgresql
}

initdb() {
	mkdir -p $PGDATA
	rm -rf $PGDATA/*
	reset_owner
	gosu postgres initdb $PGDATA
}

init_basebackup() {
  mkdir -p "$PGDATA"
  rm -rf $PGDATA/*
  chown -R postgres:postgres "$PGDATA"
  chmod g+s /run/postgresql
  chown -R postgres /run/postgresql
  load_password
  POSTGRES_PASSWORD=""

  # Wait for postgres to start
	# ref: http://unix.stackexchange.com/a/5279
	while ! nc -q 1 $MASTER 5432 </dev/null; do echo "Waiting... Master pod is not ready yet"; sleep 5; done
	PGPASSWORD=$REPLICA_PASSWORD gosu postgres pg_basebackup -x -P -R -D $PGDATA  -h $MASTER -U rep -v

}

run_user_script() {
	# run user script to allow creating db schema
  cd "$PGSCRIPT"
  if [ -n  "$1" ]; then
    while [ ! -f "$1" ]
    do
    	echo "Waiting for $1"
      sleep 2
    done
    chown -R postgres:postgres *
    chmod -R 777 *
	gosu postgres pg_ctl -D $PGDATA  -w start
    gosu postgres "$1"
	gosu postgres pg_ctl -D $PGDATA -m fast -w stop
	cd /
  fi
}

config_pg_hba() {
	{ echo; echo 'host all all 10.0.0.0/8     password'; } >> "$PGDATA/pg_hba.conf"
	{ echo; echo 'host all all 172.16.0.0/12  password'; } >> "$PGDATA/pg_hba.conf"
	{ echo; echo 'host all all 192.168.0.0/16 password'; } >> "$PGDATA/pg_hba.conf"
	{ echo; echo 'host all all 0.0.0.0/0      md5'; }      >> "$PGDATA/pg_hba.conf"
	if [ "$1" != 'basic' ]; then
		{ echo; echo 'host replication rep 0.0.0.0/0 md5'; } >> "$PGDATA/pg_hba.conf"
	fi
}

set_listen_addresses() {
	sedEscapedValue="$(echo "$1" | sed 's/[\/&]/\\&/g')"
	sed -ri "s/^#?(listen_addresses\s*=\s*)\S+/\1'$sedEscapedValue'/" $PGDATA/postgresql.conf
}

# different from basic
load_password() {
	###### get postgres user password ######
	if [ -f '/srv/postgres/secrets/.admin' ]; then
		export $(cat /srv/postgres/secrets/.admin | xargs)
	else
		echo
		echo 'Missing environment file /srv/postgres/secrets/.admin. Using default password.'
		echo
		POSTGRES_PASSWORD=${POSTGRES_PASSWORD:-postgres}
		REPLICA_PASSWORD=${REPLICA_PASSWORD:-rep}
	fi
}

# different from basic
set_password() {
	load_password
	set_listen_addresses '*'
	gosu postgres pg_ctl -D $PGDATA  -w start

	psql --username postgres <<-EOSQL
ALTER USER postgres WITH SUPERUSER PASSWORD '$POSTGRES_PASSWORD';
EOSQL

## Check if For Basic or not
  if [ "$1" != 'basic' ]; then
		psql --username postgres <<-EOSQL
CREATE USER rep REPLICATION LOGIN ENCRYPTED PASSWORD '$REPLICA_PASSWORD';
EOSQL
  fi

	gosu postgres pg_ctl -D $PGDATA -m fast -w stop

	echo
	echo 'PostgreSQL init process complete; ready for start up.'
	echo
}

# different from basic
configure_pgpool() {
	cat >>$PGDATA/postgresql.conf <<EOL
pgpool.pg_ctl = '/usr/lib/postgresql/$PG_MAJOR/bin/pg_ctl'
EOL
	gosu postgres pg_ctl -D $PGDATA  -w start

	psql --username postgres -f /opt/pgpool/src/sql/insert_lock.sql template1
	# psql --username postgres -f /opt/pgpool/src/sql/pgpool-recovery/pgpool-recovery.sql template1
	psql --username postgres -c 'CREATE EXTENSION pgpool_recovery' template1
	psql --username postgres -c 'CREATE EXTENSION pgpool_adm' template1

	gosu postgres pg_ctl -D $PGDATA -m fast -w stop

	echo
	echo 'Pgpool-II configuration complete.'
	echo
}

set_wal-e_credential() {
		##############################################################
		################ Setting up Cloud Bucket ####################
		umask u=rwx,g=rx,o=
		mkdir -p /etc/wal-e.d/env

		if [ -f "/var/credentials/gce" ]; then
			echo 'gs://'$BUCKET'/'$DATABASE > /etc/wal-e.d/env/WALE_GS_PREFIX
			echo '/var/credentials/gce' > /etc/wal-e.d/env/GOOGLE_APPLICATION_CREDENTIALS
			echo $PROJECT > /etc/wal-e.d/env/GCLOUD_PROJECT
		elif [ -d "/var/credentials/aws" ]; then
			echo 's3://'$BUCKET'/'$DATABASE > /etc/wal-e.d/env/WALE_S3_PREFIX
			cat /var/credentials/aws/secret > /etc/wal-e.d/env/AWS_SECRET_ACCESS_KEY
			cat /var/credentials/aws/keyid > /etc/wal-e.d/env/AWS_ACCESS_KEY_ID
			cat /var/credentials/aws/region > /etc/wal-e.d/env/AWS_REGION
		else
			echo
			echo 'Missing credentials ... '
			echo
			exit 1
		fi
		chown -R root:postgres /etc/wal-e.d
		##############################################################
}

config_postgresql() {
  cat <<-EOF >> "$PGDATA/postgresql.conf"
listen_addresses = '*'
ssl = true
ssl_cert_file = '/etc/ssl/certs/ssl-cert-snakeoil.pem'
ssl_key_file = '/etc/ssl/private/ssl-cert-snakeoil.key'
EOF
  if [ "$1" = 'basic' ]; then
    return
  fi

	cat <<-EOF >> "$PGDATA/postgresql.conf"
wal_level = 'hot_standby'
max_wal_senders = $MW_SENDER
wal_keep_segments = $WK_SEGMENTS
EOF

###########################################
##  For Wal Archive in external server ####
  if [ "$1" = "cloud" ]; then
cat <<-EOF >> "$PGDATA/postgresql.conf"
archive_mode = on
archive_command = 'envdir /etc/wal-e.d/env /usr/local/bin/wal-e wal-push %p'
archive_timeout =  $ARC_TIME
EOF
  fi
}

backup-push() {
    gosu postgres envdir /etc/wal-e.d/env /usr/local/bin/wal-e backup-push $PGDATA
    retval=$?
    if [ "$retval" -ne 0 ]; then
        echo "fail"
        exit 1
    else
        echo "success"
    fi
}
