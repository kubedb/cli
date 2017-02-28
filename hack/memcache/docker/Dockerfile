FROM appscode/gosu:8.5

# add our user and group first to make sure their IDs get assigned consistently, regardless of whatever dependencies get added
RUN groupadd -r memcache && useradd -r -g memcache memcache

RUN apt-get update && apt-get install -y --no-install-recommends \
                libevent-2.0-5 \
        && rm -rf /var/lib/apt/lists/* /usr/share/doc /usr/share/man /tmp/*

ENV MEMCACHED_VERSION 1.4.25
ENV MEMCACHED_SHA1 7fd0ba9283c61204f196638ecf2e9295688b2314

RUN buildDeps='curl gcc libc6-dev libevent-dev make perl' \
        && set -x \
        && apt-get update && apt-get install -y $buildDeps --no-install-recommends \
        && rm -rf /var/lib/apt/lists/* /usr/share/doc /usr/share/man /tmp/* \
        && curl -SL "http://memcached.org/files/memcached-$MEMCACHED_VERSION.tar.gz" -o memcached.tar.gz \
        && echo "$MEMCACHED_SHA1  memcached.tar.gz" | sha1sum -c - \
        && mkdir -p /usr/src/memcached \
        && tar -xzf memcached.tar.gz -C /usr/src/memcached --strip-components=1 \
        && rm memcached.tar.gz \
        && cd /usr/src/memcached \
        && ./configure \
        && make \
        && make install \
        && cd / && rm -rf /usr/src/memcached \
        && apt-get purge -y --auto-remove $buildDeps

COPY sv /etc/sv/
RUN ln -s /etc/sv /etc/service

ENTRYPOINT ["/runit.sh"]
EXPOSE 11211
