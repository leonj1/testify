#!/bin/bash

export DBNAME=${DBNAME:=testify}
export DBUSER=${DBUSER:=root}
export DBPASS=${DBPASS:=}
export DBHOST=${DBHOST:=localhost}
export DBPORT=${DBPORT:=3306}
export HTTPPORT=${HTTPPORT:=443}

cd /app
/app/testify \
    -user=${DBUSER} \
    -pass=${DBPASS} \
    -db-name=${DBNAME} \
    -db-host=${DBHOST} \
    -db-port=${DBPORT} \
    -http-port=${HTTPPORT}

