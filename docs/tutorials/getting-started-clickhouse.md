# Getting started with ClickHouse

This tutorial assumes that the latest version of sqlc is
[installed](../overview/install.md) and ready to use.

We'll generate Go code here, but other
[language plugins](../reference/language-support.rst) are available. You'll
naturally need the Go toolchain if you want to build and run a program with the
code sqlc generates, but sqlc itself has no dependencies.

## Setting up

Create a new directory called `sqlc-tutorial` and open it up.

Initialize a new Go module named `tutorial.sqlc.dev/app`:

```shell
go mod init tutorial.sqlc.dev/app
```

sqlc looks for either a `sqlc.(yaml|yml)` or `sqlc.json` file in the current
directory. In our new directory, create a file named `sqlc.yaml` with the
following contents:

```yaml
version: "2"
sql:
  - engine: "clickhouse"
    queries: "query.sql"
    schema: "schema.sql"
    gen:
      go:
        package: "tutorial"
        out: "tutorial"
        sql_package: "clickhouse/v2"
```

## Schema and queries

sqlc needs to know your database schema and queries in order to generate code.
In the same directory, create a file named `schema.sql` with the following
content:

```sql
CREATE TABLE authors (
  id   UInt64,
  name String,
  bio  String
) ENGINE = Memory;
```

Next, create a `query.sql` file with the following four queries:

```sql
-- name: GetAuthor :one
SELECT * FROM authors
WHERE id = ? LIMIT 1;

-- name: ListAuthors :many
SELECT * FROM authors
ORDER BY name;

-- name: CreateAuthor :exec
INSERT INTO authors (
  id, name, bio
) VALUES (
  ?, ?, ?
);

-- name: DeleteAuthor :exec
DELETE FROM authors
WHERE id = ?;
```

## Generating code

You are now ready to generate code. You shouldn't see any output when you run
the `generate` subcommand, unless something goes wrong:

```shell
sqlc generate
```

You should now have a `tutorial` subdirectory with three files containing Go
source code. These files comprise a Go package named `tutorial`:

```
├── go.mod
├── query.sql
├── schema.sql
├── sqlc.yaml
└── tutorial
    ├── db.go
    ├── models.go
    └── query.sql.go
```

## Using generated code

You can use your newly-generated `tutorial` package from any Go program.
Create a file named `tutorial.go` and add the following contents:

```go
package main

import (
	"context"
	"log"

	"github.com/ClickHouse/clickhouse-go/v2"

	"tutorial.sqlc.dev/app/tutorial"
)

func run() error {
	ctx := context.Background()

	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{"localhost:9000"},
		Auth: clickhouse.Auth{
			Database: "default",
			Username: "default",
			Password: "",
		},
	})
	if err != nil {
		return err
	}
	defer conn.Close()

	queries := tutorial.New(conn)

	// list all authors
	authors, err := queries.ListAuthors(ctx)
	if err != nil {
		return err
	}
	log.Println(authors)

	// create an author
	err = queries.CreateAuthor(ctx, tutorial.CreateAuthorParams{
		ID:   1,
		Name: "Brian Kernighan",
		Bio:  "Co-author of The C Programming Language and The Go Programming Language",
	})
	if err != nil {
		return err
	}

	// get the author we just inserted
	fetchedAuthor, err := queries.GetAuthor(ctx, 1)
	if err != nil {
		return err
	}
	log.Println(fetchedAuthor)

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
```

Before this code will compile you'll need to fetch the relevant ClickHouse driver:

```shell
go get github.com/ClickHouse/clickhouse-go/v2
go build ./...
```

The program should compile without errors. To make that possible, sqlc generates
readable, **idiomatic** Go code that you otherwise would've had to write
yourself. Take a look in `tutorial/query.sql.go`.

Of course for this program to run successfully you'll need to compile after
replacing the database connection parameters in the call to `clickhouse.Open()`
with the correct parameters for your database. And your database must have the
`authors` table as defined in `schema.sql`.

You should now have a working program using sqlc's generated Go source code,
and hopefully can see how you'd use sqlc in your own real-world applications.
