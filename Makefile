build:
	go build ./...

test:
	go test ./...

sqlc-dev:
	go build -o ~/bin/sqlc-dev ./cmd/sqlc/

sqlc-pg-gen:
	go build -o ~/bin/sqlc-pg-gen ./internal/tools/sqlc-pg-gen

regen: sqlc-dev
	go run ./scripts/regenerate/
