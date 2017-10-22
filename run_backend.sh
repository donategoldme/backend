#!/usr/bin/env bash
docker rm dgm_backend
docker build -t dgm_backend .
docker run -it --name dgm_backend -v uploads:/var/uploads/ \
--link dgm_postgres:postgres --link dgm_tarantool:tarantool --link dgm_centrifugo:centrifugo \
--link dgm_chats:chats --env-file env \
dgm_backend




