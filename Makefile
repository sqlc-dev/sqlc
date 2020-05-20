build:
	go build ./...

test:
	SQLC_EXPERIMENTAL_PARSER=on go test ./...

sqlc-dev:
	go build -o ~/bin/sqlc-dev ./cmd/sqlc/

regen: sqlc-dev
	./scripts/regenerate.sh
