# Getting started with PostgreSQL

This tutorial assumes that the latest version of sqlc is
[installed](../overview/install.md) and ready to use.

We'll generate Go code here, but other
[language plugins](../reference/language-support.rst) are available. You'll
naturally need the Go toolchain if you want to build and run a program with the
code sqlc generates, but sqlc itself has no dependencies.

We'll also rely on sqlc's [managed databases](../howto/managed-databases.md),
which require a sqlc Cloud project and auth token. You can get those from
the [sqlc Cloud dashboard](https://dashboard.sqlc.dev/). Managed databases are
an optional feature that improves sqlc's query analysis in many cases, but you
can turn it off simply by removing the `cloud` and `database` sections of your
configuration.

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
cloud:
  project: "<PROJECT_ID>"
sql:
  - engine: "postgresql"
    queries: "query.sql"
    schema: "schema.sql"
    database:
      managed: true
    gen:
      go:
        package: "tutorial"
        out: "tutorial"
        sql_package: "pgx/v5"
```

And finally, set the `SQLC_AUTH_TOKEN` environment variable:

```shell
export SQLC_AUTH_TOKEN="<your sqlc auth token>"
```

## Schema and queries

sqlc needs to know your database schema and queries in order to generate code.
In the same directory, create a file named `schema.sql` with the following
content:

```sql
CREATE TABLE authors (
  id   BIGSERIAL PRIMARY KEY,
  name text      NOT NULL,
  bio  text
);
```

Next, create a `query.sql` file with the following five queries:

```sql
-- name: GetAuthor :one
SELECT * FROM authors
WHERE id = $1 LIMIT 1;

-- name: ListAuthors :many
SELECT * FROM authors
ORDER BY name;

-- name: CreateAuthor :one
INSERT INTO authors (
  name, bio
) VALUES (
  $1, $2
)
RETURNING *;

-- name: UpdateAuthor :exec
UPDATE authors
  set name = $2,
  bio = $3
WHERE id = $1;

-- name: DeleteAuthor :exec
DELETE FROM authors
WHERE id = $1;
```

If you prefer, you can alter the `UpdateAuthor` query to return the updated
record:

```sql
-- name: UpdateAuthor :one
UPDATE authors
  set name = $2,
  bio = $3
WHERE id = $1
RETURNING *;
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
	"reflect"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"tutorial.sqlc.dev/app/tutorial"
)

func run() error {
	ctx := context.Background()

	conn, err := pgx.Connect(ctx, "user=pqgotest dbname=pqgotest sslmode=verify-full")
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	queries := tutorial.New(conn)

	// list all authors
	authors, err := queries.ListAuthors(ctx)
	if err != nil {
		return err
	}
	log.Println(authors)

	// create an author
	insertedAuthor, err := queries.CreateAuthor(ctx, tutorial.CreateAuthorParams{
		Name: "Brian Kernighan",
		Bio:  pgtype.Text{String: "Co-author of The C Programming Language and The Go Programming Language", Valid: true},
	})
	if err != nil {
		return err
	}
	log.Println(insertedAuthor)

	// get the author we just inserted
	fetchedAuthor, err := queries.GetAuthor(ctx, insertedAuthor.ID)
	if err != nil {
		return err
	}

	// prints true
	log.Println(reflect.DeepEqual(insertedAuthor, fetchedAuthor))
	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
```

Before this code will compile you'll need to fetch the relevant PostgreSQL
driver. You can use `lib/pq` with the standard library's `database/sql`
package, but in this tutorial we've used `pgx/v5`:

```shell
go get github.com/jackc/pgx/v5
go build ./...
```

The program should compile without errors. To make that possible, sqlc generates
readable, **idiomatic** Go code that you otherwise would've had to write
yourself. Take a look in `tutorial/query.sql.go`.

Of course for this program to run successfully you'll need
to compile after replacing the database connection parameters in the call to
`pgx.Connect()` with the correct parameters for your database. And your
database must have the `authors` table as defined in `schema.sql`.

You should now have a working program using sqlc's generated Go source code,
and hopefully can see how you'd use sqlc in your own real-world applications.
