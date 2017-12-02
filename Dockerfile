FROM ubuntu:16.04

MAINTAINER Jose Leon

RUN apt-get update
RUN apt-get install -y mysql-client

ADD bootstrap.sh /
ADD testify /app/

ENTRYPOINT ["/bootstrap.sh"]

