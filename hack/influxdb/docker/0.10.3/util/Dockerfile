#####################################
#      Image for influxdb base      #
#  Image: appscode/influxdb:base    #
#####################################

FROM appscode/cloud-tool:1.0

# Install InfluxDB
ENV INFLUXDB_VERSION 0.10.3-1

RUN apt-get update \
  && mkdir /app && wget http://s3.amazonaws.com/influxdb/influxdb_${INFLUXDB_VERSION}_amd64.deb -P /app \
  && dpkg -i /app/influxdb_${INFLUXDB_VERSION}_amd64.deb \
  && rm /app/influxdb_${INFLUXDB_VERSION}_amd64.deb

RUN  pip install influxdb \
  && pip install --upgrade influxdb

ENV PATH=/opt/influxdb:$PATH
COPY execute.py /execute.py
COPY utils.sh /utils.sh
RUN chmod +x /utils.sh

ENTRYPOINT ["python2.7", "execute.py"]
