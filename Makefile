.PHONY: build build-endtoend test test-ci test-examples test-endtoend regen start psql mysqlsh

build:
	go build ./...

test:
	go test ./...

test-examples:
	go test --tags=examples ./...

build-endtoend:
	cd ./internal/endtoend/testdata && go build ./...

test-ci: test-examples build-endtoend

regen: sqlc-dev
	go run ./scripts/regenerate/

sqlc-dev:
	go build -o ~/bin/sqlc-dev ./cmd/sqlc/

sqlc-pg-gen:
	go build -o ~/bin/sqlc-pg-gen ./internal/tools/sqlc-pg-gen

start:
	docker-compose up -d

fmt:
	go fmt ./...

psql:
	PGPASSWORD=mysecretpassword psql --host=127.0.0.1 --port=5432 --username=postgres dinotest

mysqlsh:
	mysqlsh --sql --user root --password mysecretpassword --database dinotest 127.0.0.1:3306

# $ protoc --version
# libprotoc 3.19.1
# $ go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
# $ go install github.com/planetscale/vtprotobuf/cmd/protoc-gen-go-vtproto
proto: internal/plugin/codegen.pb.go internal/python/ast/ast.pb.go

internal/plugin/codegen.pb.go: protos/plugin/codegen.proto
	protoc -I ./protos \
		--go_out=. \
		--go_opt=module=github.com/kyleconroy/sqlc \
		--go-vtproto_out=. \
		--go-vtproto_opt=module=github.com/kyleconroy/sqlc,features=marshal+unmarshal+size \
		./protos/plugin/codegen.proto

internal/python/ast/ast.pb.go: protos/python/ast.proto
	protoc -I ./protos \
		--go_out=. \
		--go_opt=module=github.com/kyleconroy/sqlc \
		./protos/python/ast.proto
