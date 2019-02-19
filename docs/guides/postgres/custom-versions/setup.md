---
title: Setup Custom PostgresVersions
menu:
  docs_0.9.0:
    identifier: pg-custom-versions-setup-postgres
    name: Overview
    parent: pg-custom-versions-postgres
    weight: 10
menu_name: docs_0.9.0
section_menu_id: guides
---
> New to KubeDB? Please start [here](/docs/concepts/README.md).

## Setting up Custom PostgresVersions

PostgresVersions are KubeDB crds that define the docker images KubeDB will use when deploying a postgres database. For more details about PostgresVersion crd, please visit [here](/docs/concepts/catalog/postgres.md).

## Creating a Custom Postgres Database Image for KubeDB

The best way to create a custom image is to build on top of the existing kubedb image.

```docker
FROM kubedb/postgres:10.2-v3

ENV TIMESCALEDB_VERSION 0.9.1

RUN set -ex \
    && apk add --no-cache --virtual .fetch-deps \
    ca-certificates \
    openssl \
    tar \
    && mkdir -p /build/timescaledb \
    && wget -O /timescaledb.tar.gz https://github.com/timescale/timescaledb/archive/$TIMESCALEDB_VERSION.tar.gz \
    && tar -C /build/timescaledb --strip-components 1 -zxf /timescaledb.tar.gz \
    && rm -f /timescaledb.tar.gz \
    \
    && apk add --no-cache --virtual .build-deps \
    coreutils \
    dpkg-dev dpkg \
    gcc \
    libc-dev \
    make \
    cmake \
    util-linux-dev \
    \
    && cd /build/timescaledb \
    && ./bootstrap \
    && cd build && make install \
    && cd ~ \
    \
    && apk del .fetch-deps .build-deps \
    && rm -rf /build

RUN sed -r -i "s/[#]*\s*(shared_preload_libraries)\s*=\s*'(.*)'/\1 = 'timescaledb,\2'/;s/,'/'/" /scripts/primary/postgresql.conf
```

From there, we would define a PostgresVersion that contains this new image. Lets say we tagged it as `myco/postgres:timescale-0.9.1`

```yaml
apiVersion: catalog.kubedb.com/v1alpha1
kind: PostgresVersion
metadata:
  name: timescale-0.9.1
spec:
  version: 10.2
  db:
    image: "myco/postgres:timescale-0.9.1"
  exporter:
    image: "kubedb/postgres_exporter:v0.4.6"
  tools:
    image: "kubedb/postgres-tools:10.2-v2"
```

Once we add this PostgresVersion we can use it in a new Postgres like:

```yaml
apiVersion: kubedb.com/v1alpha1
kind: Postgres
metadata:
  name: timescale-postgres
  namespace: demo
spec:
  version: "timescale-0.9.1" # points to the name of our custom PostgresVersion
  storage:
    storageClassName: "standard"
    accessModes:
    - ReadWriteOnce
    resources:
      requests:
        storage: 1Gi
```
