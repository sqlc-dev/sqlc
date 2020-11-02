.PHONY: build test test-examples regen

build:
	go build ./...

test:
	go test ./...

test-examples:
	go test --tags=examples ./...

regen: sqlc-dev
	go run ./scripts/regenerate/

sqlc-dev:
	go build -o ~/bin/sqlc-dev ./cmd/sqlc/

sqlc-pg-gen:
	go build -o ~/bin/sqlc-pg-gen ./internal/tools/sqlc-pg-gen
