build:
	go build ./...

test:
	go test ./...

sqlc-dev:
	go build -o $(GOPATH)/bin/sqlc-dev ./cmd/sqlc/

regen: sqlc-dev
	./scripts/regenerate.sh
