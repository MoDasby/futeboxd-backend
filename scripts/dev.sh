#!/bin/bash

cleanup() {
    echo "Stopping services..."
    make -C services/football compose:down
    make -C services/core compose:down
}

trap cleanup SIGINT

# Parar qualquer instância do air
pkill -f air || true

# Iniciar serviços usando Docker Compose
echo "Starting football service..."
make -C services/football compose:up

echo "Starting core service..."
make -C services/core compose:up

# Rodar os serviços em modo de desenvolvimento
echo "Running development for core service..."
make -C services/core dev &

echo "Running development for football service..."
make -C services/football dev

wait

cleanup
