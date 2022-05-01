#!/bin/bash

docker volume create fakebook-mysql-data

docker run \
	--detach \
	--name fakebook-mysql \
	--env MYSQL_ALLOW_EMPTY_PASSWORD=true \
	--volume fakebook-mysql-data:/var/lib/mysql \
	--publish 8080:8080 \
	mysql
