# Datatypes

## Arrays

PostgreSQL [arrays](https://www.postgresql.org/docs/current/arrays.html) are
materialized as Go slices. Currently, the `pgx/v5` sql package only supports multidimensional arrays.

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

## Dates and Time

All PostgreSQL time and date types are returned as `time.Time` structs. For
null time or date values, the `NullTime` type from `database/sql` is used.
The `pgx/v5` sql package uses the appropriate pgx types.

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

PostgreSQL [enums](https://www.postgresql.org/docs/current/arrays.html) are
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
sqlc uses the excellent `github.com/google/uuid` package.

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

For MySQL, there is no native `uuid` data type. When using `UUID_TO_BIN` to store a `UUID()`, the underlying field type is `BINARY(16)` which by default sqlc would interpret this to `sql.NullString`. To have sqlc automatically convert these fields to a `uuid.UUID` type, use an overide on the column storing the `uuid`.
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
But if you use the `pgx/v5` sql package then you can specify a some struct instead of default type.
The `pgx` implementation will marshall/unmarshall the struct automatically.

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
        "import":"example/db",
        "package": "dto",
        "type":"BookData"
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
