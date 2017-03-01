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

run_user_script() {
	# run user script to allow creating db schema
  if [ -n  "$1" ]; then
    cd "$PGSCRIPT"
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

config_postgresql() {
  cat <<-EOF >> "$PGDATA/postgresql.conf"
listen_addresses = '*'
ssl = true
ssl_cert_file = '/etc/ssl/certs/ssl-cert-snakeoil.pem'
ssl_key_file = '/etc/ssl/private/ssl-cert-snakeoil.key'
EOF
}
