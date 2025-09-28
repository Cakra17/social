BINARY_NAME=main

run:
	go run cmd/main.go

build-win:
	go build -o build/$(BINARY_NAME) cmd/main.go

deps:
	go mod tidy

create-migrate:
	migrate create -ext sql -seq -dir db/migrations create_new_table

clean:
	rm -rf	build/$(BINARY_NAME).exe