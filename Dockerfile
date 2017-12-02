FROM alpine:3.7

MAINTAINER Jose Leon

RUN apk update && \
    apk add -y mysql-client bash

ADD bootstrap.sh /
ADD testify /app/

ENTRYPOINT ["/bootstrap.sh"]

