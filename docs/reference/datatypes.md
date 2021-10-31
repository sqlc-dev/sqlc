# Datatypes

## Arrays

PostgreSQL [arrays](https://www.postgresql.org/docs/current/arrays.html) are
materialized as Go slices. Currently, only one-dimensional arrays are
supported.

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
`database/sql` package. 

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
