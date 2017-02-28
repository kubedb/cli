FROM appscode/base:8.6

# Install MongoDB.
RUN apt-get update
RUN apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv EA312927
RUN echo "deb http://repo.mongodb.org/apt/debian wheezy/mongodb-org/3.2 main" | tee /etc/apt/sources.list.d/mongodb-org-3.2.list
RUN apt-get update
RUN apt-get install -y mongodb-org
RUN rm -rf /var/lib/apt/lists/* /usr/share/doc /usr/share/man /tmp/*

# Define mountable directories
VOLUME ["/data/db"]

# Define working directory
WORKDIR /data

# Define default command
CMD ["mongod"]

# Expose ports
EXPOSE 27017
EXPOSE 28017
