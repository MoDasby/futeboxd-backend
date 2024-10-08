FROM golang:1.22.5 AS build

WORKDIR /app

COPY ./go.mod ./go.sum ./
RUN go mod download

COPY services/users/ .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o api ./cmd/api/main.go

FROM alpine:latest  

WORKDIR /root/

COPY --from=build /app/api ./api
COPY services/users/.env .env
COPY services/users/entrypoint.sh ./entrypoint.sh

RUN chmod +x ./entrypoint.sh

ENTRYPOINT [ "sh", "./entrypoint" ]