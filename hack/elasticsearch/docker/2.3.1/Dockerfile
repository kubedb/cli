############################################
#      Image for elasticsearch basic       #
#  Image: appscode/elasticsearch:basic   #
############################################

FROM appscode/py-java:1.0

## elasticsearch version
ENV ES_PKG_NAME elasticsearch-2.3.1

## gosu version
ENV GOSU_VERSION 1.7

## es user to run elasticsearch
RUN groupadd -r es --gid=999 && useradd -r -g es --uid=999 es

##### Install ElasticSearch ######
RUN set -x \
  && apt-get update \
  && apt-get install -y --no-install-recommends ca-certificates \
  && wget https://download.elasticsearch.org/elasticsearch/elasticsearch/$ES_PKG_NAME.tar.gz \
  && tar xvzf $ES_PKG_NAME.tar.gz \
  && rm -f $ES_PKG_NAME.tar.gz \
  && mv /$ES_PKG_NAME /elasticsearch \
##################################
########## install gosu ##########
  && wget -O /usr/local/bin/gosu "https://github.com/tianon/gosu/releases/download/$GOSU_VERSION/gosu-$(dpkg --print-architecture)" \
  && wget -O /usr/local/bin/gosu.asc "https://github.com/tianon/gosu/releases/download/$GOSU_VERSION/gosu-$(dpkg --print-architecture).asc" \
  && export GNUPGHOME="$(mktemp -d)" \
  && gpg --keyserver ha.pool.sks-keyservers.net --recv-keys B42F6819007F00F88E364FD4036A9C25BF357DD4 \
  && gpg --batch --verify /usr/local/bin/gosu.asc /usr/local/bin/gosu \
  && rm -r "$GNUPGHOME" /usr/local/bin/gosu.asc \
  && chmod +x /usr/local/bin/gosu \
  && gosu nobody true
##################################

########## Install cloud tools ##########
RUN apt-get update \
	&& apt-get install -y python-pip \
	&& pip install awscli \
	&& pip2 install gsutil \
#########################################
###### Elasticsearch python plugin ######
	&& pip install elasticsearch
#########################################

###### Remove unnecessary things ######
# Run set -x \
#   && apt-get purge -y --auto-remove ca-certificates \
#	  && rm -rf /var/lib/apt/lists/* /usr/share/doc /usr/share/man /tmp/*
#######################################

###### install cloud-aws & cloud-gce plugin ######
RUN set -x \
  && /elasticsearch/bin/plugin install cloud-aws url file:target/releases/elasticsearch-cloud-aws-2.7.1-SNAPSHOT.zip
##################################################

###### runit service ######
COPY sv /etc/sv
RUN ln -s /etc/sv /etc/service
RUN chmod +x /etc/service/*/run
###########################

###### Copy elasticsearch_discovery ######
COPY elasticsearch_discovery /elasticsearch_discovery
##########################################

## Python script to backup & restore
ADD execute.py /
ADD utils.sh /

## elasticsearch config file
COPY elasticsearch.yml /elasticsearch/config/elasticsearch.yml

VOLUME ["/var/pv"]
VOLUME ["/mount/snapshots"]
RUN chown -R es:es /mount/snapshots


EXPOSE 9200
EXPOSE 9300
ENTRYPOINT ["/runit.sh"]
