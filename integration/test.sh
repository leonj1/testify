#!/bin/bash

export CORE_NETWORK=core_net
export PROJECT=enchilada
export container=${PROJECT}; docker stop $container; docker rm $container
export container=${PROJECT}_db; docker stop $container; docker rm $container

#docker network create ${PROJECT}_net

export HTTP_INTERNAL=443
export HTTP_EXTERNAL=8843
export DBPORT_INTERNAL=3301
export DBPORT_EXTERNAL=6665

export DOCKER_IMAGE_TAG=$(python get_docker_build_version.py)

#docker run --name ${PROJECT}_db \
#-p ${DBPORT_EXTERNAL}:${DBPORT_INTERNAL} \
#-e MYSQL_ALLOW_EMPTY_PASSWORD=yes \
#-e MYSQL_ROOT_HOST=% \
#-e CHECK_INTERVAL=15 \
#-v /Users/jose/go/src/github.com/leonj1/${PROJECT}/resources:/docker-entrypoint-initdb.d \
#-v /Users/jose/go/src/github.com/leonj1/${PROJECT}/resources/my.cnf:/etc/my.cnf \
#--net ${CORE_NETWORK} \
#-d mysql/mysql-server:latest
#
#echo "Waiting for DB to come online"
#while ! netstat -tna | grep 'LISTEN\>' | grep -q '.'${DBPORT_EXTERNAL}; do
#  sleep 5
#done

docker run -d --name ${PROJECT} \
-p ${HTTP_EXTERNAL}:${HTTP_INTERNAL} \
-e DBNAME=${PROJECT} \
-e DBUSER=root \
-e DBPASS= \
-e DBHOST=docker.for.mac.localhost \
-e DBPORT=3306 \
--net ${CORE_NETWORK} \
www.dockerhub.us/${PROJECT}:${DOCKER_IMAGE_TAG}

echo "Waiting for ${PROJECT} to come online"
while ! netstat -tna | grep 'LISTEN\>' | grep -q '.'${HTTP_EXTERNAL}; do
  sleep 5
done

sleep 10

curl http://localhost:${HTTP_EXTERNAL}/public/health

