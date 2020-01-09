# Schema modification

sqlc understands `ALTER TABLE` statements when parsing SQL.

```sql
CREATE TABLE authors (
  id          SERIAL PRIMARY KEY,
  birth_year  int    NOT NULL
);

ALTER TABLE authors ADD COLUMN bio text NOT NULL;
ALTER TABLE authors DROP COLUMN birth_year;
ALTER TABLE authors RENAME TO writers;
```

```go
package db

type Writer struct {
	ID  int
	Bio string
}
```
