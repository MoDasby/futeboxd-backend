#! /bin/bash

kill_air() {
    pkill -f air || true
}

core() {
    make -C core $1
}

football() {
    make -C football $1
}

football-data-updater() {
    if [ "$1" = "dev" ]; then
        cd football-data-updater && npm run dev
    fi
}

cleanup() {
    core stop
    football stop

    kill_air
}

trap cleanup SIGINT

kill_air

set -e

core dev &
football dev &
football-data-updater dev &

wait 

cleanup
