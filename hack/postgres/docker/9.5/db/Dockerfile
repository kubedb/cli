#######################################
#     Base Image for postgres DB      #
#   Image: appscode/postgres:9.5-db   #
#######################################

FROM appscode/gosu:8.5

RUN groupadd -r postgres --gid=999 && useradd -r -g postgres --uid=999 postgres
RUN apt-key adv --keyserver ha.pool.sks-keyservers.net --recv-keys B97B0AFCAA1A47F044F244A07FCC7D46ACCC4CF8

ENV PG_MAJOR 9.5
ENV PGPOOL_VER 3.5.3
ENV PATH /usr/lib/postgresql/$PG_MAJOR/bin:$PATH
ENV PGROOT /var/pv
ENV PGSCRIPT /var/db-script
ENV PGDATA /var/pv/data

RUN echo 'deb http://apt.postgresql.org/pub/repos/apt/ jessie-pgdg main' $PG_MAJOR > /etc/apt/sources.list.d/pgdg.list

RUN set -x \
  && apt-get update \
  && apt-get install -y --no-install-recommends ca-certificates postgresql-common \
  && sed -ri 's/#(create_main_cluster) .*$/\1 = false/' /etc/postgresql-common/createcluster.conf \
  && apt-get install -y  --no-install-recommends postgresql-$PG_MAJOR postgresql-contrib-$PG_MAJOR

RUN set -x \
  && apt-get install -y --no-install-recommends daemontools libevent-dev python-pip python-all-dev lzop pv build-essential automake libtool

RUN set -x \
  && apt-get install -y --no-install-recommends build-essential libffi-dev libssl-dev libpq-dev wget postgresql-server-dev-all \
  && mkdir -p /opt/pgpool \
  && cd /opt/pgpool \
  && wget http://www.pgpool.net/download.php?f=pgpool-II-$PGPOOL_VER.tar.gz -O pgpool-II-$PGPOOL_VER.tar.gz \
  && tar xzvf pgpool-II-$PGPOOL_VER.tar.gz --strip-components=1 \
  && rm -rf pgpool-II-$PGPOOL_VER.tar.gz \
  && ./configure --with-openssl \
  && make \
  && make install

RUN set -x \
  && cd /opt/pgpool/src/sql/pgpool-recovery \
  && make \
  && make install \
  && cd /opt/pgpool/src/sql/pgpool_adm \
  && make \
  && make install

RUN set -x \
  && apt-get purge -y --auto-remove build-essential automake libtool git wget unzip \
  && rm -rf /var/lib/apt/lists/* /usr/share/doc /usr/share/man /tmp/*

RUN mkdir /scripts
COPY scripts /scripts/
ENTRYPOINT ["/scripts/run.sh"]
EXPOSE 5432
