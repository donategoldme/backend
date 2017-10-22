#!/usr/bin/env bash
docker stop dgm_chats
docker rm dgm_chats
docker build -t dgm_chats chats/
docker run -it --env-file env --name dgm_chats --link dgm_centrifugo:centrifugo \
--link dgm_postgres:postgres dgm_chats