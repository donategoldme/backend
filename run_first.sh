#!/usr/bin/env bash
docker stop dgm_postgres dgm_tarantool dgm_centrifugo
docker rm dgm_postgres dgm_tarantool dgm_centrifugo
docker run -d --env-file env -v pgdata:/var/lib/postgresql/9.5/data/ -p 5432:5432 --name dgm_postgres postgres:9.5
docker run -d --env-file env --name dgm_tarantool tarantool/tarantool:1.7.3
docker run -d -p 8000:8000 --env-file env --name dgm_centrifugo centrifugo/centrifugo:1.7.1

