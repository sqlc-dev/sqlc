# Developing sqlc

## Building

For local development, install `sqlc` under an alias. We suggest `sqlc-dev`.

```
go build -o ~/go/bin/sqlc-dev ./cmd/sqlc
```

## Running Tests

```
go test ./...
```

To run the tests in the examples folder, use the `examples` tag.

```
go test --tags=examples ./...
```

These tests require locally-running database instances. Run these databases
using [Docker Compose](https://docs.docker.com/compose/).

```
docker compose up -d
```

The tests use the following environment variables to connect to the
database

### For PostgreSQL

```
Variable     Default Value
-------------------------
PG_HOST      127.0.0.1
PG_PORT      5432
PG_USER      postgres
PG_PASSWORD  mysecretpassword
PG_DATABASE  dinotest
```

### For MySQL

```
Variable     Default Value
-------------------------
MYSQL_HOST      127.0.0.1
MYSQL_PORT      3306
MYSQL_USER      root
MYSQL_ROOT_PASSWORD  mysecretpassword
MYSQL_DATABASE  dinotest
```
