############################################
#      Image for elasticsearch dump        #
#  Image: appscode/elasticsearch:dump    #
############################################

FROM appscode/cloud-tool:1.0

RUN set -x \
  && apt-get update \
  && curl -sL https://deb.nodesource.com/setup | bash - \
  && apt-get install -y nodejs npm \
  && npm install npm -g \
  && npm install -g n \
  && n stable \
  && npm install elasticdump@2.4.2 -g \
  && ln -s /usr/bin/nodejs /usr/bin/node \
  && pip install --upgrade pip \
  && pip install elasticsearch

COPY execute.py /execute.py
COPY utils.sh /utils.sh
RUN chmod +x /utils.sh

ENTRYPOINT ["python2.7", "execute.py"]
