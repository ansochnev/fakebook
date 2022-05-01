#!/bin/bash

docker cp fakebook fakebook-mysql:/root
docker cp site fakebook-mysql:/root
docker cp configs/fakebook.yaml fakebook-mysql:/root

docker cp mysql/init.sql fakebook-mysql:/root
docker cp mysql/drop.sql fakebook-mysql:/root
