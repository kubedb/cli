#!/bin/bash

# env variable email and password
if [ -f '/srv/pgadmin/secrets/.env' ]; then
	export $(cat /srv/pgadmin/secrets/.env | xargs)
elif [ "$#" -ne 2 ]; then
    echo "Since no secrets file found, email and password must be passed as command line arguments"
    exit 1
else
	email=$1
	shift
	password=$1
fi
sed -i -e "s/email = ''/email = '$email'/g" /usr/local/lib/python2.7/site-packages/pgadmin4/setup.py
sed -i -e "s/p1 = ''/p1 = '$password'/g" /usr/local/lib/python2.7/site-packages/pgadmin4/setup.py

export > /etc/envvars

cat >/usr/local/lib/python2.7/site-packages/pgadmin4/config_local.py <<EOL
# -*- coding: utf-8 -*-

DEFAULT_SERVER='0.0.0.0'

# Secret key for signing CSRF data.
CSRF_SESSION_KEY = '$(date +%s | sha256sum | base64 | head -c 32 ; echo)'

# Secret key for signing cookies.
SECRET_KEY = '$(date +%s | sha256sum | base64 | head -c 32 ; echo)'

# Salt used when hashing passwords.
SECURITY_PASSWORD_SALT = '$(date +%s | sha256sum | base64 | head -c 32 ; echo)'
EOL

# echo "Starting runit..."
# exec /usr/sbin/runsvdir-start

# Location to store backup files.
mkdir -p /root/.pgadmin/storage/admin

exec python /usr/local/lib/python2.7/site-packages/pgadmin4/pgAdmin4.py
