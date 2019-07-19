> 🚨 
> DinoSQL is **very new** and is under rapid development. The command line
> interface and generated code will change. There are currently no guarantees
> around stability or compatibility.
> 🚨 

# Dino: A SQL Compiler

> And lo, the Great One looked down upon the people and proclaimed:
> 
>   "SQL is actually pretty great"

DinoSQL *generates idiomatic Go code from SQL*. Save yourself the pain of
writing boilerplate `database/sql` code.

A quick `dinosql generate` turns the following SQL:

```sql
CREATE TABLE authors (
  id   SERIAL PRIMARY KEY,
  name text   NOT NULL,
  bio  text
);

-- name: GetAuthor :one
SELECT * FROM authors
WHERE id = $1;

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

Into Go you'd have to write yourself:

```go
package db

import (
  "context"
  "database/sql"
)

type Author struct {
  ID   int
  Name string
  Bio  sql.NullString
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

func (q *Queries) DeleteAuthor(ctx context.Context, id int) error {
  _, err := q.db.ExecContext(ctx, deleteAuthor, id)
  return err
}

const getAuthor = `-- name: GetAuthor :one
SELECT id, name, bio FROM authors
WHERE id = $1
`

func (q *Queries) GetAuthor(ctx context.Context, id int) (Author, error) {
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
```

## Limitations

DinoSQL currently only supports PostgreSQL. There are no plans to add support
for other databases.

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
  - [Transactions](./examples/transactions.md)
  - [Prepared queries](./examples/prepared_query.md)
- PostgreSQL Types
  - [enum](./examples/enums.md)
  - [timestamp](./examples/time.md)
  - [uuid](./examples/uuid.md)
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
  dinosql [command]

Available Commands:
  compile     Statically check SQL for syntax and type errors
  generate    Generate Go code from SQL
  help        Help about any command
  init        Create an empty dinosql.json settings file
  version     Print the DinoSQL version number

Flags:
  -h, --help   help for dinosql

Use "dinosql [command] --help" for more information about a command.
```

## Settings

The `dinosql` tool is configured via a `dinosql.json` file. This file must be
in the directory where the `dinosql` command is run.

```json
{
  "package": "db",
  "emit_json_tags": true,
  "emit_prepared_queries": false,
  "out": "internal/db/db.go",
  "queries": "./sql/query/",
  "schema": "./sql/schema/"
}
```

- `package`:
  - The package name to use for the generated code
- `emit_json_tags`:
  - If true, add JSON tags to generated structs. Defaults to `false`.
- `emit_prepared_queries`:
  - If true, include support for prepared queries. Defaults to `false`.
- `out`:
  - Filename for generated code
- `queries`:
  - Directory of SQL queries stored in `.sql` files
- `schema`:
  - Directory of SQL migrations, stored in `.sql` files

## Acknowledgements

DinoSQL was inspired by [PugSQL](https://pugsql.org/) and
[HugSQL](https://www.hugsql.org/).
