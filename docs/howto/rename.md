# Renaming fields

Struct field names are generated from column names using a simple algorithm:
split the column name on underscores and capitalize the first letter of each
part.

```
account     -> Account
spotify_url -> SpotifyUrl
app_id      -> AppID
```

If you're not happy with a field's generated name, use the `rename` mapping
to pick a new name. The keys are column names and the values are the struct
field name to use.

```yaml
version: "2"
sql:
- schema: "postgresql/schema.sql"
  queries: "postgresql/query.sql"
  engine: "postgresql"
  gen:
    go: 
      package: "authors"
      out: "postgresql"
      rename:
        spotify_url: "SpotifyURL"
```

## Tables

The output structs associated with tables can also be renamed. By default, the struct name will be the singular version of the table name. For example, the `authors` table will generate an `Author` struct.

```sql
CREATE TABLE authors (
  id   BIGSERIAL PRIMARY KEY,
  name text      NOT NULL,
  bio  text
);
```

```go
package db

import (
	"database/sql"
)

type Author struct {
	ID   int64
	Name string
	Bio  sql.NullString
}
```

To rename this struct, you must use the generated struct name. In this example, that would be `author`. Use the `rename` map to change the name of this struct to `Writer` (note the uppercase `W`).

```yaml
version: '1'
packages:
- path: db
  engine: postgresql
  schema: query.sql
  queries: query.sql
rename:
  author: Writer
```

```yaml
version: "2"
sql:
  - engine: postgresql
    queries: query.sql
    schema: query.sql
overrides:
  go:
    rename:
      author: Writer
```

```go
package db

import (
	"database/sql"
)

type Writer struct {
	ID   int64
	Name string
	Bio  sql.NullString
}
```

## Limitations

Rename mappings apply to an entire package. Therefore, a column named `foo` and
a table name `foo` can't map to different rename values.