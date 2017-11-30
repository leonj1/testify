FROM ubuntu:14.04

MAINTAINER Jose Leon

RUN apt-get update
RUN apt-get install -y mysql-client

ADD localhost.crt /app/
ADD server.key /app/
ADD bootstrap.sh /
ADD testify /app/

ENTRYPOINT ["/bootstrap.sh"]

