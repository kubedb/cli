################################################
#       Image for singleton postgres DB        #
#      Image: appscode/postgres:9.4-st           #
################################################

FROM appscode/postgres:9.4-base

MAINTAINER shahriar@appscode.com

COPY postgres-master.sh /postgres-start.sh
RUN chmod 500 /postgres-start.sh

EXPOSE 5432

ENTRYPOINT ["/postgres-start.sh"]

CMD ["postgres"]
