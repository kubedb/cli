FROM appscode/nginx-wsgi:1.10-2.7

ENV PG_MAJOR 9.5
ENV PATH /usr/lib/postgresql/$PG_MAJOR/bin:$PATH

RUN set -x \
  & apt-key adv --keyserver ha.pool.sks-keyservers.net --recv-keys B97B0AFCAA1A47F044F244A07FCC7D46ACCC4CF8 \
  && echo 'deb http://apt.postgresql.org/pub/repos/apt/ jessie-pgdg main' $PG_MAJOR > /etc/apt/sources.list.d/pgdg.list \
  && apt-get update

RUN set -x \
  && apt-get install -y --no-install-recommends postgresql-client-$PG_MAJOR --force-yes

RUN set -x \
  && apt-get update \
  && apt-get install -y --no-install-recommends build-essential libffi-dev libssl-dev libpq-dev wget \
  && wget https://ftp.postgresql.org/pub/pgadmin3/pgadmin4/v1.1/pip/pgadmin4-1.1-py2-none-any.whl \
  && pip install ./pgadmin4-1.1-py2-none-any.whl \
  && pip install gsutil \
  && pip install awscli \
  && mkdir -p /root/.pgadmin \
  && apt-get purge -y --auto-remove build-essential wget \
  && rm -rf /root/.pip/cache /var/lib/apt/lists/* /usr/share/doc /usr/share/man /tmp/*

COPY nginx.conf  /etc/nginx/conf.d/default.conf
COPY setup.py    /usr/local/lib/python2.7/site-packages/pgadmin4/setup.py
# COPY pgadmin4.db /root/.pgadmin/pgadmin4.db

# Setup runit scripts
COPY sv /etc/sv/
RUN ln -s /etc/sv /etc/service

COPY runit.sh /runit.sh
ENTRYPOINT ["/runit.sh"]
EXPOSE 5050
