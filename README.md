# sqlc: A SQL Compiler

> And lo, the Great One looked down upon the people and proclaimed:
>
>   "SQL is actually pretty great"

sqlc generates **fully type-safe idiomatic Go code** from SQL. Here's how it
works:

1. You write SQL queries
1. You run sqlc to generate Go code that presents type-safe interfaces to those
   queries
1. You write application code that calls the methods sqlc generated

Seriously, it's that easy. You don't have to write any boilerplate SQL querying
code ever again.

## Preventing Errors

But sqlc doesn't just make you more productive by generating boilerplate for
you. sqlc **also prevents entire classes of common errors in SQL code**. Have
you ever:

- Mixed up the order of the arguments when invoking the query so they didn't
  match up with the SQL text
- Updated the name of a column in one query but not another
- Mistyped the name of a column in a query
- Changed the number of arguments in a query but forgot to pass the additional
  values
- Changed the type of a column but forgot to change the type in your code?

All of these errors are *impossible* with sqlc. Wait, what? How?

sqlc parses all of your queries and the DDL (e.g. `CREATE TABLE`)
statements during the code generation processes so that it knows the names and
types of every column in your tables and every expression in your queries.  If
any of them do not match, sqlc *will fail to compile your queries*, preventing
entire classes of runtime problems at compile time.

Likewise, the methods that sqlc generates for you have a strict arity and
correct Go type definitions that match your columns. So if you change a query's
arguments or a column's type but don't update your code, it will fail to
compile.

## Getting Started
Okay, enough hype, let's see it in action.

First you pass the following SQL to `sqlc generate`:

```sql
CREATE TABLE authors (
  id   BIGSERIAL PRIMARY KEY,
  name text      NOT NULL,
  bio  text
);

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

And then in your application code you'd write:

```go

// list all authors
authors, err := db.ListAuthors(ctx)
if err != nil {
    return err
}
fmt.Println(authors)

// create an author
insertedAuthor, err := db.CreateAuthor(ctx, db.CreateAuthorParams{
        Name: "Brian Kernighan",
        Bio:  sql.NullString{String: "Co-author of The C Programming Language and The Go Programming Language", Valid: true},
})
if err != nil {
        return err
}
fmt.Println(insertedAuthor)

// get the author we just inserted
fetchedAuthor, err := db.GetAuthor(ctx, insertedAuthor.ID)
if err != nil {
        return err
}
// prints true
fmt.Println(reflect.DeepEqual(insertedAuthor, fetchedAuthor))
```

To make that possible, sqlc generates readable, **idiomatic** Go code that you
otherwise would have had to write yourself. Take a look:

```go
package db

import (
	"context"
	"database/sql"
)

type Author struct {
	ID   int64
	Name string
	Bio  sql.NullString
}

const createAuthor = `-- name: CreateAuthor :one
INSERT INTO authors (
  name, bio
) VALUES (
  $1, $2
)
RETURNING id, name, bio
`

type CreateAuthorParams struct {
	Name string
	Bio  sql.NullString
}

func (q *Queries) CreateAuthor(ctx context.Context, arg CreateAuthorParams) (Author, error) {
	row := q.db.QueryRowContext(ctx, createAuthor, arg.Name, arg.Bio)
	var i Author
	err := row.Scan(&i.ID, &i.Name, &i.Bio)
	return i, err
}

const deleteAuthor = `-- name: DeleteAuthor :exec
DELETE FROM authors
WHERE id = $1
`

func (q *Queries) DeleteAuthor(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteAuthor, id)
	return err
}

const getAuthor = `-- name: GetAuthor :one
SELECT id, name, bio FROM authors
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetAuthor(ctx context.Context, id int64) (Author, error) {
	row := q.db.QueryRowContext(ctx, getAuthor, id)
	var i Author
	err := row.Scan(&i.ID, &i.Name, &i.Bio)
	return i, err
}

const listAuthors = `-- name: ListAuthors :many
SELECT id, name, bio FROM authors
ORDER BY name
`

func (q *Queries) ListAuthors(ctx context.Context) ([]Author, error) {
	rows, err := q.db.QueryContext(ctx, listAuthors)
	 if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Author
	for rows.Next() {
		var i Author
		if err := rows.Scan(&i.ID, &i.Name, &i.Bio); err != nil {
			 return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

func New(db DBTX) *Queries {
	return &Queries{db: db}
}

type Queries struct {
	db DBTX
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db: tx,
	}
}
```

## Examples

Your favorite PostgreSQL / Go features are supported:
- SQL
  - [Query annotations](./docs/annotations.md)
  - [Transactions](./docs/transactions.md)
  - [Prepared queries](./docs/prepared_query.md)
  - [Named parameters](./docs/named_parameters.md)
  - [SELECT](./docs/query_one.md)
  - [NULL](./docs/null.md)
  - [COUNT](./docs/query_count.md)
  - [INSERT](./docs/insert.md)
  - [UPDATE](./docs/update.md)
  - [DELETE](./docs/delete.md)
  - [RETURNING](./docs/returning.md)
  - [ANY](./docs/any.md)
- PostgreSQL Types
  - [Arrays](./docs/arrays.md)
  - [Enums](./docs/enums.md)
  - [Timestamps](./docs/time.md)
  - [UUIDs](./docs/uuid.md)
- DDL
  - [CREATE TABLE](./docs/table.md)
  - [ALTER TABLE](./docs/alter_table.md)
- Go
  - [JSON struct tags](./docs/json_tags.md)
  - [Migration tools](./docs/migrations.md)

A full, end-to-end example can be found in the sample
[`ondeck`](./examples/ondeck) package.

## Usage

```
Usage:
  sqlc [command]

Available Commands:
  compile     Statically check SQL for syntax and type errors
  generate    Generate Go code from SQL
  help        Help about any command
  init        Create an empty sqlc.yaml settings file
  version     Print the sqlc version number

Flags:
  -h, --help   help for sqlc

Use "sqlc [command] --help" for more information about a command.
```

## Settings

The `sqlc` tool is configured via a `sqlc.yaml` file. This file must be
in the directory where the `sqlc` command is run.

```yaml
version: "1"
packages:
  - name: "db"
    path: "internal/db"
    queries: "./sql/query/"
    schema: "./sql/schema/"
    engine: "postgresql"
    emit_json_tags: true
    emit_prepared_queries: true
    emit_interface: false
    emit_exact_table_names: false
    emit_empty_slices: false
```

Each package document has the following keys:
- `name`:
  - The package name to use for the generated code. Defaults to `path` basename
- `path`:
  - Output directory for generated code
- `queries`:
  - Directory of SQL queries or path to single SQL file
- `schema`:
  - Directory of SQL migrations or path to single SQL file
- `engine`:
  - Either `postgresql` or `mysql`. Defaults to `postgresql`. MySQL support is experimental
- `emit_json_tags`:
  - If true, add JSON tags to generated structs. Defaults to `false`.
- `emit_prepared_queries`:
  - If true, include support for prepared queries. Defaults to `false`.
- `emit_interface`:
  - If true, output a `Querier` interface in the generated package. Defaults to `false`.
- `emit_exact_table_names`:
  - If true, struct names will mirror table names. Otherwise, sqlc attempts to singularize plural table names. Defaults to `false`.
- `emit_empty_slices`:
  - If true, slices returned by `:many` queries will be empty instead of `nil`. Defaults to `false`.

### Type Overrides

The default mapping of PostgreSQL types to Go types only uses packages outside
the standard library when it must.

For example, the `uuid` PostgreSQL type is mapped to `github.com/google/uuid`.
If a different Go package for UUIDs is required, specify the package in the
`overrides` array. In this case, I'm going to use the `github.com/gofrs/uuid`
instead.

```yaml
version: "1"
packages: [...]
overrides:
  - go_type: "github.com/gofrs/uuid.UUID"
    db_type: "uuid"
```

Each override document has the following keys:
- `db_type`:
  - The PostgreSQL type to override. Find the full list of supported types in [gen.go](https://github.com/kyleconroy/sqlc/blob/master/internal/dinosql/gen.go#L438).
- `go_type`:
  - A fully qualified name to a Go type to use in the generated code.
- `nullable`:
  - If true, use this type when a column is nullable. Defaults to `false`.

### Per-Column Type Overrides

Sometimes you would like to override the Go type used in model or query generation for
a specific field of a table and not on a type basis as described in the previous section.

This may be configured by specifying the `column` property in the override definition. `column`
should be of the form `table.column` buy you may be even more specify by specifying `schema.table.column`
or `catalog.schema.table.column`.

```yaml
version: "1"
packages: [...]
overrides:
  - column: "authors.id"
    go_type: "github.com/segmentio/ksuid.KSUID"
```

### Package Level Overrides

Overrides can be configured globally, as demonstrated in the previous sections, or they can be configured on a per-package which
scopes the override behavior to just a single package:

```yaml
version: "1"
packages:
  - overrides: [...]
```

### Renaming Struct Fields

Struct field names are generated from column names using a simple algorithm:
split the column name on underscores and capitalize the first letter of each
part.

```
account     -> Account
spotify_url -> SpotifyUrl
app_id      -> AppID
```

If you're not happy with a field's generated name, use the `rename` dictionary
to pick a new name. The keys are column names and the values are the struct
field name to use.

```yaml
version: "1"
packages: [...]
rename:
  spotify_url: "SpotifyURL"
```

## Installation

### macOS

```
brew install kyleconroy/sqlc/sqlc
```

### Ubuntu

```
sudo snap install sqlc
```

### go get

```
go get github.com/kyleconroy/sqlc/cmd/sqlc
```

### Docker

```
docker pull kjconroy/sqlc
```

Run `sqlc` using `docker run`:

```
docker run --rm -v $(pwd):/src -w /src kjconroy/sqlc generate
```

### Downloads

Binaries for a given release can be downloaded from the [stable channel on
Equinox](https://dl.equinox.io/sqlc/sqlc/stable) or the latest [GitHub
release](https://github.com/kyleconroy/sqlc/releases).

### Tip Releases

Each commit is deployed to the [`devel` channel on Equinox](https://dl.equinox.io/sqlc/sqlc/devel):
- [Linux](https://bin.equinox.io/c/gvM95th6ps1/sqlc-devel-linux-amd64.tgz)
- [macOS](https://bin.equinox.io/c/gvM95th6ps1/sqlc-devel-darwin-amd64.zip)

## Other Databases and Languages

sqlc currently only supports PostgreSQL / Go. MySQL and Kotlin support have
been merged, but both are marked as experimental. SQLite and TypeScript support
are planned.

| Language     | PostgreSQL        | MySQL             |
| ------------ | :---------------- | :---------------- |
| Go           | :white_check_mark: - Stable | :bug: - Beta |
| TypeScript   | :timer_clock: - Planned | :timer_clock: - Planned |
| Kotlin       | :warning: - Experimental | |

If you'd like to add another database or language, we'd welcome a contribution.

## Sponsors

sqlc development is funded by our generous sponsors.

- Companies
  - [Meter](https://meter.com)
  - [ngrok](https://ngrok.com)
  - [Weave](https://www.getweave.com/)
- Individuals
  - [Alex Besogonov](https://github.com/Cyberax)
  - [Myles McDonnell](https://github.com/myles-mcdonnell)

If you use sqlc at your company, please consider [becoming a
sponsor](https://github.com/sponsors/kyleconroy) today.

Sponsors receive priority support via the sqlc Slack organization.

## Development

### Building

For local development, install `sqlc` under an alias. We suggest `sqlc-dev`.

```
go build -o ~/go/bin/sqlc-dev ./cmd/sqlc
```

### Running Tests

To run the tests, include the `exp` tag. Without this tag, a few tests will
fail.

```
go test --tags=exp ./...
```

To run the tests in the examples folder, a running PostgreSQL instance is
required. The tests use the following environment variables to connect to the
database:

```
Variable     Default Value
-------------------------
PG_HOST      127.0.0.1
PG_PORT      5432
PG_USER      postgres
PG_PASSWORD  mysecretpassword
PG_DATABASE  dinotest
```

```
go test --tags=examples,exp ./...
```

### Regenerate expected test output

If you need to update a large number of expected test output in the
`internal/endtoend/testdata` directory, run the `regenerate.sh` script.

```
make regen
```

Note that this uses the `sqlc-dev` binary, not `sqlc` so make sure you have an
up to date `sqlc-dev` binary.

## Acknowledgements

sqlc was inspired by [PugSQL](https://pugsql.org/) and
[HugSQL](https://www.hugsql.org/).
