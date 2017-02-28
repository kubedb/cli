#!/bin/bash

#
# This file will create execute pg_dumpall remotely
# pg_dumpall -h service_name > dumpfile.sql
#
# Generate .boto file
#

if [ "$CLOUD" = 'GCE' ]; then
	cat <<-EOF >> "/root/.boto"
	[Credentials]
	gs_service_key_file = /var/credential/gce

	[Boto]
	https_validate_certificates = True

	[GSUtil]
	content_language = en
	default_api_version = 2
	default_project_id = $PROJECT
	EOF

	chmod +x /root/.boto

	PASS=`cat /var/credential/postgres/password`
	USER=`cat /var/credential/postgres/username`

	cd /var/dump
	gsutil cp  gs://$BUCKET/postgres/$POSTGRES/$1.sql dumpfile.sql
	PGPASSWORD=$PASS psql -U $USER -h $POSTGRES  -f dumpfile.sql postgres
fi


if [ "$CLOUD" = 'AWS' ]; then

	ACCESS_KEY=`cat /var/credential/aws/keyid`
    SECRET_KEY=`cat /var/credential/aws/secret`
    mkdir -p "/root/.aws/"
	cat <<-EOF >> "/root/.aws/credentials"
	[default]
	aws_access_key_id = ${ACCESS_KEY}
	aws_secret_access_key = ${SECRET_KEY}
	EOF

	PASS=`cat /var/credential/postgres/password`
	USER=`cat /var/credential/postgres/username`

	cd /var/dump
	aws s3 cp s3://$BUCKET/postgres/$POSTGRES/$1.sql dumpfile.sql 
	PGPASSWORD=$PASS psql -U $USER -h $POSTGRES  -f dumpfile.sql postgres
fi