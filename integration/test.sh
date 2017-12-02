#!/bin/bash -x

set -e

export CORE_NETWORK=core_net
export PROJECT=testify
export container=${PROJECT}; docker stop $container; docker rm $container
export container=${PROJECT}_db; docker stop $container; docker rm $container

docker network create core_net || true

export HTTP_INTERNAL=8080
export HTTP_EXTERNAL=8743
# my.cnf is source of truth
export DBPORT_INTERNAL=3201
export DBPORT_EXTERNAL=6261

export DOCKER_IMAGE_TAG=$(python get_docker_build_version.py)

docker run --name ${PROJECT}_db \
-p ${DBPORT_EXTERNAL}:${DBPORT_INTERNAL} \
-e MYSQL_ALLOW_EMPTY_PASSWORD=yes \
-e MYSQL_ROOT_HOST=% \
-v $(pwd)/resources:/docker-entrypoint-initdb.d \
-v $(pwd)/resources/my.cnf:/etc/my.cnf \
--net ${CORE_NETWORK} \
-d mysql/mysql-server:latest

echo "Waiting for DB to come online"
while ! netstat -tna | grep 'LISTEN\>' | grep -q '.'${DBPORT_EXTERNAL}; do
  sleep 5
done

docker run -d --name ${PROJECT} \
-p ${HTTP_EXTERNAL}:${HTTP_INTERNAL} \
-e DBNAME=${PROJECT} \
-e DBUSER=root \
-e DBPASS= \
-e DBHOST=${PROJECT}_db \
-e DBPORT=${DBPORT_INTERNAL} \
-e HTTPPORT=${HTTP_INTERNAL} \
--net ${CORE_NETWORK} \
www.dockerhub.us/${PROJECT}:${DOCKER_IMAGE_TAG}

echo "Waiting for ${PROJECT} to come online"
while ! netstat -tna | grep 'LISTEN\>' | grep -q '.'${HTTP_EXTERNAL}; do
  sleep 5
done

sleep 10

curl http://localhost:${HTTP_EXTERNAL}/public/health

