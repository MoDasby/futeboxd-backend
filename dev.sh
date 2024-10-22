#!/bin/bash

cleanup() {
	make -C services/football compose:stop
	make -C services/core compose:stop
}

trap cleanup SIGINT

pkill -f air || true

make -C services/football compose:up
make -C services/core compose:up

make -C services/core dev &
make -C services/football dev

wait

cleanup
