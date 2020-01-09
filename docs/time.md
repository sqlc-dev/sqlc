# Time

```sql
CREATE TABLE authors (
  id         SERIAL    PRIMARY KEY,
  created_at timestamp NOT NULL DEFAULT NOW(),
  updated_at timestamp
);
```

All PostgreSQL time and date types are returned as `time.Time` structs. For
null time or date values, the `NullTime` type is used from the
`github.com/lib/pq` package.

```go
package db

import (
	"time"

	"github.com/lib/pq"
)

type Author struct {
	ID        int
	CreatedAt time.Time
	UpdatedAt pq.NullTime
}
```

