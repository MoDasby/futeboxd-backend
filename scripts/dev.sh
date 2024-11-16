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

wait 

cleanup
