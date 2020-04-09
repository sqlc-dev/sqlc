build:
	go build --tags=exp ./...

test:
	go test --tags=exp ./...

sqlc-dev:
	go build -o ~/bin/sqlc-dev --tags=exp ./cmd/sqlc/

regen: sqlc-dev
	./scripts/regenerate.sh
