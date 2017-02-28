FROM appscode/nginx-php:1.10-5

ENV PG_MAJOR 9.5
ENV PGPOOL_VER 3.5.3

RUN set -x \
  && apt-get update \
  && apt-get install -y --no-install-recommends ca-certificates wget \
  && wget -O - http://apt.postgresql.org/pub/repos/apt/ACCC4CF8.asc | apt-key add - \
  && echo 'deb http://apt.postgresql.org/pub/repos/apt/ jessie-pgdg main' $PG_MAJOR > /etc/apt/sources.list.d/pgdg.list

RUN set -x \
  && apt-get install -y --no-install-recommends build-essential libffi-dev libssl-dev libpq-dev wget \
  && mkdir -p /opt/pgpool \
  && cd /opt/pgpool \
  && wget http://www.pgpool.net/download.php?f=pgpool-II-$PGPOOL_VER.tar.gz -O pgpool-II-$PGPOOL_VER.tar.gz \
  && tar xzvf pgpool-II-$PGPOOL_VER.tar.gz --strip-components=1 \
  && rm -rf pgpool-II-$PGPOOL_VER.tar.gz \
  && ./configure --with-openssl \
  && make \
  && make install \
  && mkdir -p /var/run/pgpool

RUN set -x \
  && mkdir -p /var/www \
  && rm -rf /var/www/* \
  && cd /var/www \
  && wget http://www.pgpool.net/download.php?f=pgpoolAdmin-$PGPOOL_VER.tar.gz -O pgpoolAdmin-$PGPOOL_VER.tar.gz \
  && tar xzvf pgpoolAdmin-$PGPOOL_VER.tar.gz --strip-components=1 \
  && rm -rf pgpoolAdmin-$PGPOOL_VER.tar.gz \
  && sed -i "s/'ja')/'en')/" /var/www/conf/pgmgt.conf.php \
  && chown -R www-data:www-data /var/www

COPY pgpool.conf /usr/local/etc/pgpool.conf
COPY nginx.conf  /etc/nginx/conf.d/default.conf

RUN set -x \
  && apt-get purge -y --auto-remove build-essential wget \
  && rm -rf /var/lib/apt/lists/* /usr/share/doc /usr/share/man /tmp/*

# Setup runit scripts
COPY sv /etc/sv
RUN ln -s /etc/sv /etc/service

COPY runit.sh /runit.sh
ENTRYPOINT ["/runit.sh"]
EXPOSE 80
EXPOSE 5432
