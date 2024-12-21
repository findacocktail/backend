
SHORT_HASH=$(shell git rev-parse --short HEAD)
IMAGE_NAME=cocktails

default: run

build: bin/backend

bin/backend:
	go build -o bin/backend main.go

run:
	go run main.go

image:
	docker build . -t $(IMAGE_NAME):$(SHORT_HASH)

test:
	go test -coverprofile=c.out ./...
	go tool cover -func=c.out
