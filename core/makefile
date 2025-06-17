dev:
	@docker compose -f ./database/postgres/compose.dev.yml up -d && \
	go mod tidy && \
	echo "Iniciando a API core..." && \
	~/go/bin/air --build.cmd "go build -o bin/api cmd/api/main.go" --build.bin "./bin/api"

stop:
	@docker compose -f ./database/postgres/compose.dev.yml stop

test:
	@go test ./...

build:
	@go build -o bin/api cmd/api/main.go

docker\:build:
	docker build -t modasby/futeboxd-core:latest .

docker\:push:
	docker push modasby/futeboxd-core:latest
