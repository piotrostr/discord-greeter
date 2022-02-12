#!/usr/bin/env bash

docker kill $(docker ps -q)
docker run -it -d --rm --env-file env.list bot
