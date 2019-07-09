# Struct tags

```sql
CREATE TABLE authors (
  id         SERIAL    PRIMARY KEY,
  created_at timestamp NOT NULL
);
```

When `emit_tags` is set, structs are generated with JSON struct tags. The JSON
name for fields match the column name in the database.

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
