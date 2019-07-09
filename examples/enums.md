# Enums

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
