#!/bin/bash

export TAG=3.1.1
export NETWORK=couchdb-dev
export TIMEZONE=Asia/Shanghai
export ROOTUSER=admin
export ROOTPASSWD=passwd
export PORT=5984
export CONTAINER=couchdb-$PORT
if [ $1 = "f" ]; then
    docker-compose -f couchdb.yaml up
else
    docker-compose -f couchdb.yaml up -d
fi

