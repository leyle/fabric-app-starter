#!/bin/bash

export TAG=3.1.1
export NETWORK=couchdb-dev
export TIMEZONE=Asia/Shanghai
export ROOTUSER=admin
export ROOTPASSWD=passwd
export PORT=5984
export CONTAINER=couchdb-$PORT

docker-compose -f couchdb.yaml down
