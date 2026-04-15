# Datatypes

`sqlc` attempts to make reasonable default choices when mapping internal
database types to Go types. Choices for more complex types are described below.

If you're unsatisfied with the default, you can override any type using the
[overrides list](config.md#overrides) in your `sqlc` config file.

## Arrays

PostgreSQL [arrays](https://www.postgresql.org/docs/current/arrays.html) are
materialized as Go slices.

```sql
CREATE TABLE places (
  name text   not null,
  tags text[]
);
```

```go
package db

type Place struct {
	Name string
	Tags []string
}
```

## Dates and times

All date and time types are returned as `time.Time` structs. For
null time or date values, the `NullTime` type from `database/sql` is used.

The `pgx/v5` sql package uses the appropriate pgx types.

For MySQL users relying on `github.com/go-sql-driver/mysql`, ensure that
`parseTime=true` is added to your database connection string.

```sql
CREATE TABLE authors (
  id         SERIAL    PRIMARY KEY,
  created_at timestamp NOT NULL DEFAULT NOW(),
  updated_at timestamp
);
```

```go
package db

import (
	"database/sql"
	"time"
)

type Author struct {
	ID        int
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}
```

## Enums

PostgreSQL [enums](https://www.postgresql.org/docs/current/datatype-enum.html) are
mapped to an aliased string type.

```sql
CREATE TYPE status AS ENUM (
  'open',
  'closed'
);

CREATE TABLE stores (
  name   text    PRIMARY KEY,
  status status  NOT NULL
);
```

```go
package db

type Status string

const (
	StatusOpen   Status = "open"
	StatusClosed Status = "closed"
)

type Store struct {
	Name   string
	Status Status
}
```

## Null

For structs, null values are represented using the appropriate type from the
`database/sql` or `pgx` package.

```sql
CREATE TABLE authors (
  id   SERIAL PRIMARY KEY,
  name text   NOT NULL,
  bio  text
);
```

```go
package db

import (
	"database/sql"
)

type Author struct {
	ID   int
	Name string
	Bio  sql.NullString
}
```

## UUIDs

The Go standard library does not come with a `uuid` package. For UUID support,
sqlc uses the excellent `github.com/google/uuid` package. The pgx/v5 sql package uses `pgtype.UUID`.

```sql
CREATE TABLE records (
  id   uuid PRIMARY KEY
);
```

```go
package db

import (
	"github.com/google/uuid"
)

type Author struct {
	ID uuid.UUID
}
```

For MySQL, there is no native `uuid` data type. When using `UUID_TO_BIN` to store a `UUID()`, the underlying field type is `BINARY(16)` which by default sqlc would map to `sql.NullString`. To have sqlc automatically convert these fields to a `uuid.UUID` type, use an overide on the column storing the `uuid`
(see [Overriding types](../howto/overrides.md) for details).

```json
{
  "overrides": [
    {
      "column": "*.uuid",
      "go_type": "github.com/google/uuid.UUID"
    }
  ]
}
```

## JSON

By default, sqlc will generate the `[]byte`, `pgtype.JSON` or `json.RawMessage` for JSON column type.
But if you use the `pgx/v5` sql package then you can specify a struct instead of the default type
(see [Overriding types](../howto/overrides.md) for details).
The `pgx` implementation will marshal/unmarshal the struct automatically.

```go
package dto

type BookData struct {
	Genres    []string `json:"genres"`
	Title     string   `json:"title"`
	Published bool     `json:"published"`
}
```

```sql
CREATE TABLE books (
  data jsonb
);
```

```json
{
  "overrides": [
    {
      "column": "books.data",
      "go_type": {
        "import":"example.com/db",
        "package": "dto",
        "type":"BookData",
        "pointer": true
      }
    }
  ]
}
```

```go
package db

import (
	"example.com/db/dto"
)

type Book struct {
    Data *dto.BookData
}
```

## TEXT

In PostgreSQL, when you have a column with the TEXT type, sqlc will map it to a Go string by default. This default mapping applies to `TEXT` columns that are not nullable. However, for nullable `TEXT` columns, sqlc maps them to `pgtype.Text` when using the pgx/v5 driver. This distinction is crucial for developers looking to handle null values appropriately in their Go applications.

To accommodate nullable strings and map them to `*string` in Go, you can use the `emit_pointers_for_null_types` option in your sqlc configuration. This option ensures that nullable SQL columns are represented as pointer types in Go, allowing for a clear distinction between null and non-null values. Another way to do this is by passing the option `pointer: true` when you are overriding the `TEXT` datatype in your sqlc config file (see [Overriding types](../howto/overrides.md) for details).

## Geometry

### PostGIS

#### Using `github.com/twpayne/go-geos` (pgx/v5 only)

sqlc can be configured to use the [geos](https://github.com/twpayne/go-geos)
package for working with PostGIS geometry types in [GEOS](https://libgeos.org/).

There are three steps:

1. Configure sqlc to use `*github.com/twpayne/go-geos.Geom` for geometry types (see [Overriding types](../howto/overrides.md) for details).
2. Call `github.com/twpayne/pgx-geos.Register` on each
   `*github.com/jackc/pgx/v5.Conn`.
3. Annotate your SQL with `::geometry` typecasts, if needed.

```sql
-- Multipolygons in British National Grid (epsg:27700)
create table shapes(
  id serial,
  name varchar,
  geom geometry(Multipolygon, 27700)
);

-- name: GetCentroids :many
SELECT id, name, ST_Centroid(geom)::geometry FROM shapes;
```

```json
{
  "version": 2,
  "gen": {
    "go": {
      "overrides": [
        {
          "db_type": "geometry",
          "go_type": {
            "import": "github.com/twpayne/go-geos",
            "package": "geos",
            "pointer": true,
            "type": "Geom"
          },
          "nullable": true
        }
      ]
    }
  }
}
```

```go
import (
    "github.com/twpayne/go-geos"
    pgxgeos "github.com/twpayne/pgx-geos"
)

// ...

config.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
    if err := pgxgeos.Register(ctx, conn, geos.NewContext()); err != nil {
        return err
    }
    return nil
}
```


#### Using `github.com/twpayne/go-geom`

sqlc can be configured to use the [geom](https://github.com/twpayne/go-geom)
package for working with PostGIS geometry types. See [Overriding types](../howto/overrides.md) for more information.

```sql
-- Multipolygons in British National Grid (epsg:27700)
create table shapes(
  id serial,
  name varchar,
  geom geometry(Multipolygon, 27700)
);

-- name: GetShapes :many
SELECT * FROM shapes;
```

```json
{
  "version": "1",
  "packages": [
    {
      "path": "db",
      "engine": "postgresql",
      "schema": "query.sql",
      "queries": "query.sql"
    }
  ],
  "overrides": [
    {
      "db_type": "geometry",
      "go_type": "github.com/twpayne/go-geom.MultiPolygon"
    },
    {
      "db_type": "geometry",
      "go_type": "github.com/twpayne/go-geom.MultiPolygon",
      "nullable": true
    }
  ]
}
```
