dev:
	@export $(shell grep -v '^#' .env | xargs) && \
	docker compose --env-file ./.env -f ./database/compose.dev.yml up -d && \
	go mod tidy && \
	sh ./wait-to-db-ready.sh && \
	echo "Iniciando a API do football..." && \
	~/go/bin/air --build.cmd "go build -o bin/api cmd/api/main.go" --build.bin "./bin/api"

stop:
	docker compose --env-file ./.env -f ./database/compose.dev.yml stop

test:
	@go test ./...

build:
	@go build -o bin/api cmd/api/main.go

docker\:build:
	docker build -t modasby/futeboxd-football:latest .

docker\:push:
	docker push modasby/futeboxd-football:latest