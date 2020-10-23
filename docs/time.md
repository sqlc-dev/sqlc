# Time

```sql
CREATE TABLE authors (
  id         SERIAL    PRIMARY KEY,
  created_at timestamp NOT NULL DEFAULT NOW(),
  updated_at timestamp
);
```

All PostgreSQL time and date types are returned as `time.Time` structs. For
null time or date values, the `NullTime` type from `database/sql` is used.

```go
package db

import (
	"time"
	"database/sql"
)

type Author struct {
	ID        int
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}
```

