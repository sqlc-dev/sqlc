# Configuring generated structs

## Naming scheme

Structs generated from tables will attempt to use the singular form of a table
name if the table name is pluralized.

```sql
CREATE TABLE authors (
  id   SERIAL PRIMARY KEY,
  name text   NOT NULL
);
```

```go
package db

// Struct names use the singular form of table names
type Author struct {
	ID   int
	Name string
}
```

## JSON tags

```sql
CREATE TABLE authors (
  id         SERIAL    PRIMARY KEY,
  created_at timestamp NOT NULL
);
```

sqlc can generate structs with JSON tags by adding the `emit_json_tags` key to the configuration file as it shows on [configuration reference](../reference/config.md).
The JSON name for a field matches
the column name in the database.

```go
package db

import (
	"time"
)

type Author struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}
```

## More control

See the guide to [Overriding types](./overrides.md) for fine-grained control over struct field types and tags.
