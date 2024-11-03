.PHONY: dev test build

dev:
	@export $(shell grep -v '^#' .env | xargs) && \
	go mod tidy && \
	sh ./wait-to-db-ready.sh && \
	echo "Iniciando a API do core..." && \
	~/go/bin/air --build.cmd "go build -o bin/api cmd/api/main.go" --build.bin "./bin/api"

test:
	@go test ./... 

build:
	@go build -o bin/api cmd/api/main.go

docker\:build:
	docker build -t futeboxd-core:latest .