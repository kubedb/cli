#!/bin/bash

redis-cli save
gsutil cp dump.rdb gs://database-br/backup
rm -f dump.rdb

