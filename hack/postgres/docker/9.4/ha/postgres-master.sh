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

		###################################################################
		########  modify postgresql.conf &  pg_hba.conf as master #########
		cat <<-EOF >> "$PGDATA/postgresql.conf"
			listen_addresses = '*'
			wal_level = 'hot_standby'
			max_wal_senders = $MW_SENDER
			wal_keep_segments = $WK_SEGMENTS

			ssl = true
			ssl_cert_file = '/etc/ssl/certs/ssl-cert-snakeoil.pem'
			ssl_key_file = '/etc/ssl/private/ssl-cert-snakeoil.key'
		EOF

		{ echo; echo "host all all 0.0.0.0/0 md5"; } >> "$PGDATA/pg_hba.conf"
		{ echo; echo "hostssl replication rep 0.0.0.0/0 trust"; } >> "$PGDATA/pg_hba.conf"
		###################################################################

		gosu postgres pg_ctl -D "$PGDATA" -o "-c listen_addresses=''" -w start

		psql --username postgres <<-EOSQL
				ALTER USER postgres WITH SUPERUSER PASSWORD 'postgres';
		EOSQL

		psql --username postgres <<-EOSQL
				CREATE USER rep REPLICATION LOGIN ENCRYPTED PASSWORD 'change_secret';
		EOSQL

		set_listen_addresses '*'
		gosu postgres pg_ctl -D "$PGDATA" -m fast -w stop

		echo
		echo 'PostgreSQL init process complete; ready for start up.'
		echo

		touch $PGDATA/.appscode
		chown -R postgres:postgres $PGDATA/.appscode

	fi
	exec gosu postgres "$@"
fi
exec "$@"
