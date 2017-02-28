###################################################
#          Image for Ha postgres cluster          #
#         Image: appscode/postgres:9.4-ha           #
###################################################

FROM appscode/postgres:9.4-base-arc

MAINTAINER shahriar@appscode.com

ENV MW_SENDER 3
ENV WK_SEGMENTS 50

COPY postgres-master.sh /postgres-master.sh
RUN chmod 500 /postgres-master.sh
COPY postgres-slave.sh /postgres-slave.sh
RUN chmod 500 /postgres-slave.sh

EXPOSE 5432

