#####################################
#      Image for influxdb base      #
#  Image: appscode/influxdb:base    #
#####################################

FROM appscode/cloud-tool:1.0

# Install InfluxDB
ENV INFLUXDB_VERSION 1.0.0
ENV PVROOT /var/pv
ENV DBROOT /var/pv/db

RUN set -x \
  && apt-get update \
  && apt-get install -y --no-install-recommends ca-certificates \
  && wget https://dl.influxdata.com/influxdb/releases/influxdb_${INFLUXDB_VERSION}_amd64.deb \
  && dpkg -i influxdb_${INFLUXDB_VERSION}_amd64.deb \
  && rm influxdb_${INFLUXDB_VERSION}_amd64.deb \
  && apt-get purge -y --auto-remove ca-certificates \
  && rm -rf /var/lib/apt/lists/* /usr/share/doc /usr/share/man /tmp/*

ENV PATH=/opt/influxdb:$PATH
ENV CONFIG_FILE /etc/influxdb.toml
ADD influxdb.conf $CONFIG_FILE

RUN mkdir -p /srv/influxdb
COPY admin.txt /srv/influxdb/secrets/.admin

RUN mkdir -p ${DBROOT}

COPY sv /etc/sv
RUN ln -s /etc/sv /etc/service
RUN chmod +x /etc/service/*/run
ENTRYPOINT ["/runit.sh"]

# admin, http, udp
EXPOSE 8083 8086 8088
