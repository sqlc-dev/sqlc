# Retrieving rows

To generate a database access method, annotate a query with a specific comment.

```sql
CREATE TABLE authors (
  id         SERIAL PRIMARY KEY,
  bio        text   NOT NULL,
  birth_year int    NOT NULL
);


-- name: GetAuthor :one
SELECT * FROM authors
WHERE id = $1;

-- name: ListAuthors :many
SELECT * FROM authors
ORDER BY id;
```

A few new pieces of code are generated beyond the `Author` struct. An interface
for the underlying database is generated. The `*sql.DB` and `*sql.Tx` types
satisfy this interface.

The database access methods are added to a `Queries` struct, which is created
using the `New` method.

Note that the `*` in our query has been replaced with explicit column names.
This change ensures that the query will never return unexpected data.

Our query was annotated with `:one`, meaning that it should only return a
single row. We scan the data from that one into a `Author` struct.

Since the get query has a single parameter, the `GetAuthor` method takes a single
`int` as an argument.

Since the list query has no parameters, the `ListAuthors` method accepts no
arguments.


```go
package db

import (
	"context"
	"database/sql"
)

type Author struct {
	ID        int
	Bio       string
	BirthYear int
}

type DBTX interface {
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

func New(db DBTX) *Queries {
	return &Queries{db: db}
}

type Queries struct {
	db DBTX
}

const getAuthor = `-- name: GetAuthor :one
SELECT id, bio, birth_year FROM authors
WHERE id = $1
`

func (q *Queries) GetAuthor(ctx context.Context, id int) (Author, error) {
	row := q.db.QueryRowContext(ctx, getAuthor, id)
	var i Author
	err := row.Scan(&i.ID, &i.Bio, &i.BirthYear)
	return i, err
}

const listAuthors = `-- name: ListAuthors :many
SELECT id, bio, birth_year FROM authors
ORDER BY id
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
		if err := rows.Scan(&i.ID, &i.Bio, &i.BirthYear); err != nil {
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

## Selecting columns

```sql
CREATE TABLE authors (
  id         SERIAL PRIMARY KEY,
  bio        text   NOT NULL,
  birth_year int    NOT NULL
);

-- name: GetBioForAuthor :one
SELECT bio FROM authors
WHERE id = $1;

-- name: GetInfoForAuthor :one
SELECT bio, birth_year FROM authors
WHERE id = $1;
```

When selecting a single column, only that value that returned. The `GetBioForAuthor`
method takes a single `int` as an argument and returns a `string` and an
`error`.

When selecting multiple columns, a row record (method-specific struct) is
returned. In this case, `GetInfoForAuthor` returns a struct with two fields:
`Bio` and `BirthYear`.

If a query result has no row records, a zero value and an `ErrNoRows` error are
returned instead of a zero value and `nil`. For instance, when the `GetBioForAuthor`
result has no rows, it will return `""` and `ErrNoRows`.

```go
package db

import (
	"context"
	"database/sql"
)

type DBTX interface {
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

func New(db DBTX) *Queries {
	return &Queries{db: db}
}

type Queries struct {
	db DBTX
}

const getBioForAuthor = `-- name: GetBioForAuthor :one
SELECT bio FROM authors
WHERE id = $1
`

func (q *Queries) GetBioForAuthor(ctx context.Context, id int) (string, error) {
	row := q.db.QueryRowContext(ctx, getBioForAuthor, id)
	var i string
	err := row.Scan(&i)
	return i, err
}

const getInfoForAuthor = `-- name: GetInfoForAuthor :one
SELECT bio, birth_year FROM authors
WHERE id = $1
`

type GetInfoForAuthorRow struct {
	Bio       string
	BirthYear int
}

func (q *Queries) GetInfoForAuthor(ctx context.Context, id int) (GetInfoForAuthorRow, error) {
	row := q.db.QueryRowContext(ctx, getInfoForAuthor, id)
	var i GetInfoForAuthorRow
	err := row.Scan(&i.Bio, &i.BirthYear)
	return i, err
}
```

## Passing a slice as a parameter to a query

### PostgreSQL

In PostgreSQL,
[ANY](https://www.postgresql.org/docs/current/functions-comparisons.html#id-1.5.8.28.16)
allows you to check if a value exists in an array expression. Queries using ANY
with a single parameter will generate method signatures with slices as
arguments. Use the postgres data types, eg: int, varchar, etc.

```sql
CREATE TABLE authors (
  id         SERIAL PRIMARY KEY,
  bio        text   NOT NULL,
  birth_year int    NOT NULL
);

-- name: ListAuthorsByIDs :many
SELECT * FROM authors
WHERE id = ANY($1::int[]);
```

The above SQL will generate the following code:

```go
package db

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
)

type Author struct {
	ID        int
	Bio       string
	BirthYear int
}

type DBTX interface {
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

func New(db DBTX) *Queries {
	return &Queries{db: db}
}

type Queries struct {
	db DBTX
}

const listAuthors = `-- name: ListAuthorsByIDs :many
SELECT id, bio, birth_year FROM authors
WHERE id = ANY($1::int[])
`

func (q *Queries) ListAuthorsByIDs(ctx context.Context, ids []int) ([]Author, error) {
	rows, err := q.db.QueryContext(ctx, listAuthors, pq.Array(ids))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Author
	for rows.Next() {
		var i Author
		if err := rows.Scan(&i.ID, &i.Bio, &i.BirthYear); err != nil {
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

### MySQL and SQLite

MySQL and SQLite differ from PostgreSQL in that placeholders must be generated based on
the number of elements in the slice you pass in. Though trivial it is still
something of a nuisance. The passed in slice must not be nil or empty or an
error will be returned (ie not a panic). The placeholder insertion location is
marked by the meta-function `sqlc.slice()` (which is similar to `sqlc.arg()`
that you see documented under [Naming parameters](named_parameters.md)).

To rephrase, the `sqlc.slice('param')` behaves identically to `sqlc.arg()` it
terms of how it maps the explicit argument to the function signature, eg:

  * `sqlc.slice('ids')` maps to `ids []GoType` in the function signature
  * `sqlc.slice(cust_ids)` maps to `custIds []GoType` in the function signature
    (like `sqlc.arg()`, the parameter does not have to be quoted)

This feature is not compatible with `emit_prepared_queries` statement found in the
[Configuration file](../reference/config.md).

```sql
CREATE TABLE authors (
  id         SERIAL PRIMARY KEY,
  bio        text   NOT NULL,
  birth_year int    NOT NULL
);

-- name: ListAuthorsByIDs :many
SELECT * FROM authors
WHERE id IN (sqlc.slice('ids'));
```

The above SQL will generate the following code:

```go
package db

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type Author struct {
	ID        int
	Bio       string
	BirthYear int
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

const listAuthorsByIDs = `-- name: ListAuthorsByIDs :many
SELECT id, bio, birth_year FROM authors
WHERE id IN (/*SLICE:ids*/?)
`

func (q *Queries) ListAuthorsByIDs(ctx context.Context, ids []int64) ([]Author, error) {
	sql := listAuthorsByIDs
	var queryParams []interface{}
	if len(ids) == 0 {
		return nil, fmt.Errorf("slice ids must have at least one element")
	}
	for _, v := range ids {
		queryParams = append(queryParams, v)
	}
	sql = strings.Replace(sql, "/*SLICE:ids*/?", strings.Repeat(",?", len(ids))[1:], 1)
	rows, err := q.db.QueryContext(ctx, sql, queryParams...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Author
	for rows.Next() {
		var i Author
		if err := rows.Scan(&i.ID, &i.Bio, &i.BirthYear); err != nil {
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
