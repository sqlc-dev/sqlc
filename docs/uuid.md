# Universally unique identifiers (UUIDs)

```sql
CREATE TABLE records (
  id   uuid PRIMARY KEY
);
```

The Go standard library does not come with a `uuid` package. For UUID support,
sqlc uses the excellent `github.com/google/uuid` package.

```go
package db

import (
	"github.com/google/uuid"
)

type Author struct {
	ID   uuid.UUID
}
```
