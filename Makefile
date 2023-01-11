.PHONY: build build-endtoend test test-ci test-examples test-endtoend regen start psql mysqlsh proto

build:
	go build ./...

install:
	go install ./...

test:
	go test ./...

vet:
	go vet ./...

test-examples:
	go test --tags=examples ./...

build-endtoend:
	cd ./internal/endtoend/testdata && go build ./...

test-ci: test-examples build-endtoend vet

regen: sqlc-dev sqlc-gen-json
	go run ./scripts/regenerate/

sqlc-dev:
	go build -o ~/bin/sqlc-dev ./cmd/sqlc/

sqlc-pg-gen:
	go build -o ~/bin/sqlc-pg-gen ./internal/tools/sqlc-pg-gen

sqlc-gen-json:
	go build -o ~/bin/sqlc-gen-json ./cmd/sqlc-gen-json

start:
	docker-compose up -d

fmt:
	go fmt ./...

psql:
	PGPASSWORD=mysecretpassword psql --host=127.0.0.1 --port=5432 --username=postgres dinotest

mysqlsh:
	mysqlsh --sql --user root --password mysecretpassword --database dinotest 127.0.0.1:3306

proto:
	buf generate
