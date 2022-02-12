#!/usr/bin/env bash
docker build -t bot .
docker run -it --rm --env-file env.list bot
