#####################################
#      Image for influxdb base      #
#  Image: appscode/influxdb:base    #
#####################################

FROM appscode/cloud-tool:1.0

# Install InfluxDB
ENV INFLUXDB_VERSION 0.12.2-1

RUN set -x \
  && apt-get update \
  && apt-get install -y --no-install-recommends ca-certificates \
  && wget https://dl.influxdata.com/influxdb/releases/influxdb_${INFLUXDB_VERSION}_amd64.deb \
  && dpkg -i influxdb_${INFLUXDB_VERSION}_amd64.deb \
  && rm influxdb_${INFLUXDB_VERSION}_amd64.deb \
  && apt-get purge -y --auto-remove ca-certificates \
  && rm -rf /var/lib/apt/lists/* /usr/share/doc /usr/share/man /tmp/*

RUN  pip install influxdb \
  && pip install --upgrade influxdb

ENV PATH=/opt/influxdb:$PATH
COPY execute.py /execute.py
COPY utils.sh /utils.sh
RUN chmod +x /utils.sh

ENTRYPOINT ["python2.7", "execute.py"]
