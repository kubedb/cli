###################################################
#          Image for Ha postgres cluster          #
#         Image: appscode/postgres:9.4-arc          #
###################################################

FROM appscode/postgres:9.4-base-arc

MAINTAINER shahriar@appscode.com

ENV MW_SENDER 3
ENV WK_SEGMENTS 50
ENV ARC_TIME 1800

COPY postgres-master.sh /postgres-master.sh
RUN chmod 500 /postgres-master.sh
COPY postgres-slave.sh /postgres-slave.sh
RUN chmod 500 /postgres-slave.sh
COPY postgres-restore.sh /postgres-restore.sh
RUN chmod 500 /postgres-restore.sh

EXPOSE 5432