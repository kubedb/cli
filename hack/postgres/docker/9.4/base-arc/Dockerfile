################################################
#     Base Image for postgres DB with WAL-E    #
#      Image: appscode/postgres:9.4-base-arc   #
################################################

FROM appscode/postgres:9.4-base

MAINTAINER shahriar@appscode.com

RUN apt-get -y update \
	&& apt-get install -y lzop pv daemontools python-setuptools python-pip python-dev build-essential git \
	&& easy_install --upgrade pip \
	&& pip install git+https://github.com/wal-e/wal-e.git \
	&& pip install requests --upgrade \
	&& easy_install --upgrade pip \
	&& pip install six --upgrade \
	&& pip install pyasn1 --upgrade

RUN pip show wal-e
