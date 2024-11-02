#! /bin/bash

CORE_ENV_PATH="./services/core/.env"
FOOTBALL_ENV_PATH="./services/football/.env"

for ENV_PATH in "$CORE_ENV_PATH" "$FOOTBALL_ENV_PATH"; do
    if [ ! -e "$ENV_PATH" ]; then
        echo "arquivo env não encontrado: $ENV_PATH"
        exit 1
    fi
done

kill_air() {
    pkill -f air || true
}

cleanup() {
    docker compose --env-file "$CORE_ENV_PATH" -f ./infra/compose.dev.yml stop core-db 
    docker compose --env-file "$FOOTBALL_ENV_PATH" -f ./infra/compose.dev.yml stop football-db 

    rm -f services/football/entrypoint.sh
    rm -f services/football/wait-to-db-ready.sh

    rm -f services/core/entrypoint.sh
    rm -f services/core/wait-to-db-ready.sh

    kill_air
}

trap cleanup SIGINT

kill_air

set -e

docker compose --env-file "$CORE_ENV_PATH" -f ./infra/compose.dev.yml up core-db -d
docker compose --env-file "$FOOTBALL_ENV_PATH" -f ./infra/compose.dev.yml up football-db -d

cp infra/entrypoint.sh services/core/
cp infra/wait-to-db-ready.sh services/core/

cp infra/entrypoint.sh services/football/
cp infra/wait-to-db-ready.sh services/football/

make -C services/core dev &
make -C services/football dev &

wait 

cleanup
