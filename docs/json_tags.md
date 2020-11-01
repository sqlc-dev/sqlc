# JSON Struct tags

```sql
CREATE TABLE authors (
  id         SERIAL    PRIMARY KEY,
  created_at timestamp NOT NULL
);
```

sqlc can generate structs with JSON tags. The JSON name for a field matches
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

you can also rename the json fields by setting `rename_json_tags` to `true`.

```sqlc.json
{
{
  "version": "1",
  "packages": [
    {
      "engine": "postgresql",
      ...
      "emit_json_tags": true,
      "rename_json_tags": true
    }
  ],
  "rename": {
    "created_at": "createdAt"
  }
}
```

will generate:

```go
package db

import (
	"time"
)

type Author struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
}
```