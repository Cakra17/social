BINARY_NAME=main

run:
	go run cmd/main.go

build.win:
	env GOOS=windows GOARCH=amd64 go build -o build/$(BINARY_NAME) cmd/main.go

build.macos:
	env GOOS=darwin GOARCH=amd64 go build -o build/$(BINARY_NAME) cmd/main.go

build.linux:
	env GOOS=linux GOARCH=amd64 go build -o build/$(BINARY_NAME) cmd/main.go

docker.start:
	docker compose up -d

docker.stop:
	docker compose down

deps:
	go mod tidy

create-migrate:
	migrate create -ext sql -seq -dir db/migrations create_new_table

clean:
	rm -rf	build/
