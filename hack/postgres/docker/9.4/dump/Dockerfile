#####################################
#      Image for database dump      #
#  Image: appscode/database:dump  #
#####################################

FROM appscode/base:8.6

MAINTAINER shahriar@appscode.com

ENV PG_MAJOR 9.4
ENV PG_VERSION 9.4.6-1.pgdg80+1

RUN echo 'deb http://apt.postgresql.org/pub/repos/apt/ jessie-pgdg main' $PG_MAJOR > /etc/apt/sources.list.d/pgdg.list

RUN apt-get -y update \
    && apt-get install -y --no-install-recommends \
    && apt-get install -y gcc python-dev python-setuptools libffi-dev python-pip \
    && apt-get remove --purge gsutil \
    && easy_install -U pip \
    && pip2 install gsutil \
    && pip install awscli

RUN apt-get install -y postgresql-client-$PG_MAJOR=$PG_VERSION --force-yes

########### For postgres ###########
RUN mkdir -p /data
ADD pg_dump.sh /data/pg_dump.sh
RUN chmod a+x /data/pg_dump.sh
ADD pg_restore.sh /data/pg_restore.sh
RUN chmod a+x /data/pg_restore.sh
####################################

