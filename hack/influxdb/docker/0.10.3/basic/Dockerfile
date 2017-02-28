#####################################
#      Image for influxdb base      #
#  Image: appscode/influxdb:base    #
#####################################

FROM appscode/cloud-tool:1.0

# Install InfluxDB
ENV INFLUXDB_VERSION 0.10.3-1
ENV PVROOT /var/pv
ENV DBROOT /var/pv/db


RUN apt-get update \
  && mkdir /app && wget http://s3.amazonaws.com/influxdb/influxdb_${INFLUXDB_VERSION}_amd64.deb -P /app \
  && dpkg -i /app/influxdb_${INFLUXDB_VERSION}_amd64.deb \
  && rm /app/influxdb_${INFLUXDB_VERSION}_amd64.deb

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
