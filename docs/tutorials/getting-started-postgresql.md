# Getting started with PostgreSQL 

This tutorial assumes that the latest version of sqlc is
[installed](../overview/install.html) and ready to use.

Create a new directory called `sqlc-tutorial` and open it up.

Initialize a new Go module named `tutorial.sql.dev/app`

```shell
go mod init tutorial.sqlc.dev/app
```

sqlc looks for either a `sqlc.yaml` or `sqlc.json` file in the current
directory. In our new directory, create a file named `sqlc.yaml` with the
following contents:

```yaml
version: "2"
sql:
  - engine: "postgresql"
    queries: "query.sql"
    schema: "schema.sql"
    gen:
      go:
        package: "tutorial"
        out: "tutorial"
```

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

Next, create a `query.sql` file with the following four queries:

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

-- name: DeleteAuthor :exec
DELETE FROM authors
WHERE id = $1;
```

If you **do not** want your SQL `UPDATE` queries to return the updated record
to the user, add this to `query.sql`:

```sql
-- name: UpdateAuthor :exec
UPDATE authors
  set name = $2,
  bio = $3
WHERE id = $1;
```

Otherwise, to return the updated record to the user, add this to `query.sql`:

```sql
-- name: UpdateAuthor :one
UPDATE authors
  set name = $2,
  bio = $3
WHERE id = $1
RETURNING *;
```

You are now ready to generate code. You shouldn't see any errors or output.

```shell
sqlc generate
```

You should now have a `tutorial` package containing three files.

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

You can use your newly generated queries in `app.go`.

```go
package main

import (
	"context"
	"database/sql"
	"log"
	"reflect"

	"tutorial.sqlc.dev/app/tutorial"

	_ "github.com/lib/pq"
)

func run() error {
	ctx := context.Background()

	db, err := sql.Open("postgres", "user=pqgotest dbname=pqgotest sslmode=verify-full")
	if err != nil {
		return err
	}

	queries := tutorial.New(db)

	// list all authors
	authors, err := queries.ListAuthors(ctx)
	if err != nil {
		return err
	}
	log.Println(authors)

	// create an author
	insertedAuthor, err := queries.CreateAuthor(ctx, tutorial.CreateAuthorParams{
		Name: "Brian Kernighan",
		Bio:  sql.NullString{String: "Co-author of The C Programming Language and The Go Programming Language", Valid: true},
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

Before the code will compile, you'll need to add the Go PostgreSQL driver.

```
go get github.com/lib/pq
go build ./...
```

sqlc generates readable, **idiomatic** Go code that you otherwise would have
had to write yourself. Take a look in the `tutorial` package to see what code
sqlc generated.
