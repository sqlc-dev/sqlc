# Goose

sqlc will ignore rollback statements when parsing
[Goose](https://github.com/pressly/goose) migrations.

```sql
-- +goose Up
CREATE TABLE post (
    id    int NOT NULL,
    title text,
    body  text,
    PRIMARY KEY(id)
);

-- +goose Down
DROP TABLE post;
```

```go
package db

type Post struct {
	ID    int
	Title sql.NullString
	Body  sql.NullString
}
```
