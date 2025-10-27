# Getting started with YDB

This tutorial assumes that the latest version of sqlc is
[installed](../overview/install.md) and ready to use.

We'll generate Go code here, but other
[language plugins](../reference/language-support.rst) are available. You'll
naturally need the Go toolchain if you want to build and run a program with the
code sqlc generates, but sqlc itself has no dependencies.

At the end, you'll push your SQL queries to [sqlc
Cloud](https://dashboard.sqlc.dev/) for further insights and analysis.

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
  - engine: "ydb"
    queries: "query.sql"
    schema: "schema.sql"
    gen:
      go:
        package: "tutorial"
        out: "tutorial"
```

## Schema and queries

sqlc needs to know your database schema and queries in order to generate code.
In the same directory, create a file named `schema.sql` with the following
content:

```sql
CREATE TABLE authors (
    id   Serial,
    name Text NOT NULL,
    bio  Text,
    PRIMARY KEY (id)
);
```

Next, create a `query.sql` file with the following five queries:

```sql
-- name: GetAuthor :one
SELECT * FROM authors
WHERE id = $id LIMIT 1;

-- name: ListAuthors :many
SELECT * FROM authors
ORDER BY name;

-- name: CreateOrUpdateAuthor :one
UPSERT INTO authors (name, bio) 
VALUES (
  $name, $bio
)
RETURNING *;

-- name: DeleteAuthor :exec
DELETE FROM authors WHERE id = $id;

-- name: DropTable :exec
DROP TABLE IF EXISTS authors;
```

Note that YDB uses named parameters (`$id`, `$name`, `$bio`) rather than
positional parameters.

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

	"github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/query"

	"tutorial.sqlc.dev/app/tutorial"
)

func ptr(s string) *string {
	return &s
}

func run() error {
	ctx := context.Background()

	// Create YDB connection
	// Replace with your actual YDB endpoint
	db, err := ydb.Open(ctx, "grpcs://localhost:2136/local")
	if err != nil {
		return err
	}
	defer db.Close(ctx)

	queries := tutorial.New(db.Query())

	// list all authors
	authors, err := queries.ListAuthors(ctx)
	if err != nil {
		return err
	}
	log.Println(authors)

	// create an author
	insertedAuthor, err := queries.CreateOrUpdateAuthor(ctx, tutorial.CreateOrUpdateAuthorParams{
		Name: "Brian Kernighan",
		Bio:  ptr("Co-author of The C Programming Language and The Go Programming Language"),
	}, query.WithIdempotent())
	if err != nil {
		return err
	}
	log.Println(insertedAuthor)

	// get the author we just inserted
	fetchedAuthor, err := queries.GetAuthor(ctx, insertedAuthor.ID)
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

Before this code will compile you'll need to fetch the relevant YDB driver:

```shell
go get github.com/ydb-platform/ydb-go-sdk/v3
go build ./...
```

The program should compile without errors. To make that possible, sqlc generates
readable, **idiomatic** Go code that you otherwise would've had to write
yourself. Take a look in `tutorial/query.sql.go`.

Of course for this program to run successfully you'll need
to compile after replacing the database connection parameters in the call to
`ydb.Open()` with the correct parameters for your database. And your
database must have the `authors` table as defined in `schema.sql`.

You should now have a working program using sqlc's generated Go source code,
and hopefully can see how you'd use sqlc in your own real-world applications.

## Query verification (Not supported for YDB yet)

[sqlc Cloud](https://dashboard.sqlc.dev) provides additional verification, catching subtle bugs. To get started, create a
[dashboard account](https://dashboard.sqlc.dev). Once you've signed in, create a
project and generate an auth token. Add your project's ID to the `cloud` block
to your sqlc.yaml.

```yaml
version: "2"
cloud:
  # Replace <PROJECT_ID> with your project ID from the sqlc Cloud dashboard
  project: "<PROJECT_ID>"
sql:
  - engine: "ydb"
    queries: "query.sql"
    schema: "schema.sql"
    gen:
      go:
        package: "tutorial"
        out: "tutorial"
```

Replace `<PROJECT_ID>` with your project ID from the sqlc Cloud dashboard. It
will look something like `01HA8SZH31HKYE9RR3N3N3TSJM`.

And finally, set the `SQLC_AUTH_TOKEN` environment variable:

```shell
export SQLC_AUTH_TOKEN="<your sqlc auth token>"
```

```shell
$ sqlc push --tag tutorial
```

In the sidebar, go to the "Queries" section to see your published queries. Run
`verify` to ensure that previously published queries continue to work against
updated database schema.

```shell
$ sqlc verify --against tutorial
```
