#!/bin/bash

if [ "$service" = 'restore' ]; then
	/etc/init.d/redis-server stop
	rm -f /var/lib/redis/*
	gsutil cp gs://database-br/backup/dump.rdb /var/lib/redis
	chown redis:redis /var/lib/redis/dump.rdb
	sed -i -e "s/\(appendonly = \).*/\1\no/" /etc/redis/redis.conf
	/etc/init.d/redis-server start
	redis-cli BGREWRITEAOF
	/etc/init.d/redis-server stop
	sed -i -e "s/\(appendonly = \).*/\1\yes/" /etc/redis/redis.conf
fi

exec redis-server 
