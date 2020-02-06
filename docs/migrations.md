# Migrations

sqlc will ignore rollback statements when parsing migration SQL files. The following tools are current supported:

- [goose](https://github.com/pressly/goose)
- [sql-migrate](https://github.com/rubenv/sql-migrate)
- [tern](https://github.com/jackc/tern)
- [golang-migrate](https://github.com/golang-migrate/migrate)

## goose

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

### sql-migrate

```sql
-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE people (id int);


-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE people;
```

```go
package db

type People struct {
	ID    int32
}
```

### tern

```sql
CREATE TABLE comment (id int NOT NULL, text text NOT NULL);
---- create above / drop below ----
DROP TABLE comment;
```

```go
package db

type Comment struct {
	ID     int32 
	Text   string 
}
```

## golang-migrate

In `20060102.up.sql`:

```sql
CREATE TABLE post (
    id    int NOT NULL,
    title text,
    body  text,
    PRIMARY KEY(id)
);
```

In `20060102.down.sql`:
```sql
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
