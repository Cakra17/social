BINARY_NAME=main

run:
	go run cmd/main.go

build-win:
	go build -o build/$(BINARY_NAME) cmd/main.go

deps:
	go mod tidy

clean:
	rm -rf	build/$(BINARY_NAME).exe