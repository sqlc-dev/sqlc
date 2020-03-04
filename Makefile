build:
	go build ./...

test:
	go test --tags=exp ./...

sqlc-dev:
	go build -o ~/bin/sqlc-dev --tags=exp ./cmd/sqlc/

regen:
	cd internal/endtoend && ./regenerate.sh
