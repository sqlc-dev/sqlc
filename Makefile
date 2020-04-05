build:
	go build --tags=exp ./...

test:
	go test --tags=exp ./...

sqlc-dev:
	go build -o ~/bin/sqlc-dev --tags=exp ./cmd/sqlc/

regen:
	./scripts/regenerate.sh
