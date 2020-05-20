build:
	go build ./...

test:
	go test ./...

sqlc-dev:
	go build -o ~/bin/sqlc-dev ./cmd/sqlc/

regen: sqlc-dev
	./scripts/regenerate.sh
