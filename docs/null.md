# Null values

```sql
CREATE TABLE authors (
  id   SERIAL PRIMARY KEY,
  name text   NOT NULL,
  bio  text
);
```

For structs, null values are represented using the appropriate type from the
`database/sql` package. 

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
