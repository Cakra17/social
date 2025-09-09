BINARY_NAME=main

run:
	go run cmd/main.go

build:
	go build -o $(BINARY_NAME) cmd/main.go

deps:
	go mod tidy