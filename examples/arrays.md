# Arrays

```sql
CREATE TABLE places (
  name text   not null,
  tags text[]
);
```

PostgreSQL [arrays](https://www.postgresql.org/docs/current/arrays.html) are
materialized as Go slices. Currently, only one-dimensional arrays are
supported.

```go
package db

type Place struct {
	Name string
	Tags []string
}
```
