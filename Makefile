.PHONY: build build-endtoend test test-ci test-examples test-endtoend start psql mysqlsh proto sqlc-dev ydb test-examples-ydb gen-examples-ydb

build:
	go build ./...

install:
	go install ./...

test:
	go test ./...

test-managed:
	MYSQL_SERVER_URI="invalid" POSTGRESQL_SERVER_URI="postgres://postgres:mysecretpassword@localhost:5432/postgres" go test -v ./...

vet:
	go vet ./...

test-examples:
	go test --tags=examples ./...

ydb-examples: sqlc-dev ydb gen-examples-ydb test-examples-ydb

test-examples-ydb:
	YDB_SERVER_URI=localhost:2136 go test -v ./examples/authors/ydb/... -count=1

gen-examples-ydb:
	cd examples/authors/ && SQLCDEBUG=1 ~/bin/sqlc-dev generate && cd ../..

build-endtoend:
	cd ./internal/endtoend/testdata && go build ./...

test-ci: test-examples build-endtoend vet

sqlc-dev:
	go build -x -v -o ~/bin/sqlc-dev ./cmd/sqlc/

sqlc-pg-gen:
	go build -o ~/bin/sqlc-pg-gen ./internal/tools/sqlc-pg-gen

sqlc-gen-json:
	go build -o ~/bin/sqlc-gen-json ./cmd/sqlc-gen-json

test-json-process-plugin:
	go build -o ~/bin/test-json-process-plugin ./scripts/test-json-process-plugin/

start:
	docker compose up -d

ydb:
	docker compose up -d ydb

fmt:
	go fmt ./...

psql:
	PGPASSWORD=mysecretpassword psql --host=127.0.0.1 --port=5432 --username=postgres dinotest

mysqlsh:
	mysqlsh --sql --user root --password mysecretpassword --database dinotest 127.0.0.1:3306

proto:
	buf generate

remote-proto:
	protoc \
		--go_out=. --go_opt="Minternal/remote/gen.proto=github.com/sqlc-dev/sqlc/internal/remote" --go_opt=module=github.com/sqlc-dev/sqlc \
        --go-grpc_out=. --go-grpc_opt="Minternal/remote/gen.proto=github.com/sqlc-dev/sqlc/internal/remote" --go-grpc_opt=module=github.com/sqlc-dev/sqlc \
        internal/remote/gen.proto
