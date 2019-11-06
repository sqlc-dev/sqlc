> 🚨
>
> sqlc is **new** and under rapid development.
>
> The code it generates is correct and safe for production use, but there is
> currently no guarantee of stability or backwards-compatibility of the command
> line interface, configuration file format or generated code.
>
> 🚨

# sqlc: A SQL Compiler

> And lo, the Great One looked down upon the people and proclaimed:
>
>   "SQL is actually pretty great"

sqlc generates **fully-type safe idiomatic Go code** from SQL. Here's how it
works:

1. You write SQL queries
1. You run sqlc to generate Go code that presents type-safe interfaces to those
   queries
1. You write application code that calls the methods sqlc generated.

Seriously, it's that easy. You don't have to write any boilerplate SQL querying
code ever again.

## Preventing Errors

But sqlc doesn't just make you more productive by generating boilerplate for
you. sqlc **also prevents entire classes of common errors in SQL code**. Have
you ever:

- Mixed up the order of the arguments when invoking the query so they didn't
  match up with the SQL text
- Updated the name of a column in one query both not another
- Mistyped the name of a column in a query
- Changed the number of arguments in a query but forgot to pass the additional
  values
- Changed the type of a column but forgot to change the type in your code?

All of these errors are *impossible* with sqlc. Wait, what? How?

sqlc parses your all of your queries and the DDL (e.g. `CREATE TABLE`)
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
insertedAuthor, err := db.CreateAuthor(ctx, &db.CreateAuthorParams{
        Name: "Brian Kernighan",
        Bio: "Co-author of The C Programming Language and The Go Programming Language",
})
if err != nil {
        return err
}
fmt.Println(insertedAuthor)

// get the author we just inserted
fetchedAuthor, err = db.GetAuthor(ctx, author.ID)
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

type dbtx interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

func New(db dbtx) *Queries {
	return &Queries{db: db}
}

type Queries struct {
	db dbtx
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
  - [SELECT](./examples/query_one.md)
  - [NULL](./examples/null.md)
  - [COUNT](./examples/query_count.md)
  - [INSERT](./examples/insert.md)
  - [UPDATE](./examples/update.md)
  - [DELETE](./examples/delete.md)
  - [RETURNING](./examples/returning.md)
  - [ANY](./examples/any.md)
  - [Transactions](./examples/transactions.md)
  - [Prepared queries](./examples/prepared_query.md)
- PostgreSQL Types
  - [Arrays](./examples/arrays.md)
  - [Enums](./examples/enums.md)
  - [Timestamps](./examples/time.md)
  - [UUIDs](./examples/uuid.md)
- DDL
  - [CREATE TABLE](./examples/table.md)
  - [ALTER TABLE](./examples/alter_table.md)
- Go
  - [JSON struct tags](./examples/json_tags.md)
  - [Goose migrations](./examples/goose.md)

A full, end-to-end example can be found in the sample
[`ondeck`](./internal/dinosql/testdata/ondeck) package.

## Usage

```
Usage:
  sqlc [command]

Available Commands:
  compile     Statically check SQL for syntax and type errors
  generate    Generate Go code from SQL
  help        Help about any command
  init        Create an empty sqlc.json settings file
  version     Print the sqlc version number

Flags:
  -h, --help   help for sqlc

Use "sqlc [command] --help" for more information about a command.
```

## Settings

The `sqlc` tool is configured via a `sqlc.json` file. This file must be
in the directory where the `sqlc` command is run.

```json
{
  "version": "1",
  "packages": [
    {
      "name": "db",
      "emit_json_tags": true,
      "emit_prepared_queries": false,
      "path": "internal/db",
      "queries": "./sql/query/",
      "schema": "./sql/schema/"
    }
  ]
}
```

Each package document has the following keys:
- `name`:
  - The package name to use for the generated code. Defaults to `path` basename
- `emit_json_tags`:
  - If true, add JSON tags to generated structs. Defaults to `false`.
- `emit_prepared_queries`:
  - If true, include support for prepared queries. Defaults to `false`.
- `path`:
  - Output directory for generated code
- `queries`:
  - Directory of SQL queries or path to single SQL file
- `schema`:
  - Directory of SQL migrations or path to single SQL file

### Type Overrides

The default mapping of PostgreSQL types to Go types only uses packages outside
the standard library when it must.

For example, the `uuid` PostgreSQL type is mapped to `github.com/google/uuid`.
If a different Go package for UUIDs is required, specify the package in the
`overrides` array. In this case, I'm going to use the `github.com/gofrs/uuid`
instead.

```
{
  "version": "1",
  "packages": [...],
  "overrides": [
      {
          "go_type": "github.com/gofrs/uuid.UUID",
          "postgres_type": "uuid"
      }
  ]
}
```

Each override document has the following keys:
- `postgres_type`:
  - The PostgreSQL type to override. Find the full list of supported types in [gen.go](https://github.com/kyleconroy/sqlc/blob/master/internal/dinosql/gen.go#L438).
- `go_type`:
  - A fully qualified name to a Go type to use in the generated code.
- `null`:
  - If true, use this type when a column in nullable. Defaults to `false`.

### Per-Column Type Overrides

Sometimes you would like to override the Go type used in model or query generation for
a specific field of a table and not on a type basis as described in the previous section.

This may be configured by specifying the `column` property in the override definition. `column`
should be of the form `table.column` buy you may be even more specify by specifying `schema.table.column`
or `catalog.schema.table.column`.

```
{
  "version": "1",
  "packages": [...],
  "overrides": [
    {
      "column": "authors.id",
      "go_type": "github.com/segmentio/ksuid.KSUID"
    }
  ]
}
```

### Package Level Overrides

Overrides can be configured globally, as demonstrated in the previous sections, or they can be configured on a per-package which
scopes the override behavior to just a single package:

```
{
  "version": "1",
  "packages": [
    {
      ...
      "overrides": [...]
    }
  ],
}
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

```json
{
  "version": "1",
  "packages": [...],
  "rename": {
    "spotify_url": "SpotifyURL"
  }
}
```

## Downloads

Each commit is deployed to the [`devel` channel on Equinox](https://dl.equinox.io/sqlc/sqlc/devel):
- [Linux](https://bin.equinox.io/c/gvM95th6ps1/sqlc-devel-linux-amd64.tgz)
- [macOS](https://bin.equinox.io/c/gvM95th6ps1/sqlc-devel-darwin-amd64.zip)

## Other Database Engines

sqlc currently only supports PostgreSQL. If you'd like to support another database, we'd welcome a contribution.

## Other Language Backends

sqlc currently only generates Go code, but if you'd like to build another language backend, we'd welcome a contribution.

## Acknowledgements

sqlc was inspired by [PugSQL](https://pugsql.org/) and
[HugSQL](https://www.hugsql.org/).
