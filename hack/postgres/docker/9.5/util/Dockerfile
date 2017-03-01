#####################################
#      Image for database dump      #
#  Image: appscode/database:util  #
#####################################

FROM appscode/postgres:9.5-v3-db

RUN set -x \
	&& apt-get -y update \
	&& apt-get install -y --no-install-recommends postgresql-client-$PG_MAJOR --force-yes

RUN set -x \
    && apt-get install -y --no-install-recommends build-essential libffi-dev netcat libssl-dev \
    && pip install gsutil \
    && pip install awscli \
	&& apt-get purge -y --auto-remove build-essential \
	&& rm -rf /var/lib/apt/lists/* /usr/share/doc /usr/share/man /tmp/*

COPY execute.py /execute.py
COPY utils.sh /utils.sh
RUN chmod +x /utils.sh

ENTRYPOINT ["python2.7", "execute.py"]
