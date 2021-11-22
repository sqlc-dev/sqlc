# Modifying the database schema

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

## Handling SQL migrations

sqlc will ignore rollback statements when parsing migration SQL files. The
following tools are current supported:

- [dbmate](https://github.com/amacneil/dbmate)
- [golang-migrate](https://github.com/golang-migrate/migrate)
- [goose](https://github.com/pressly/goose)
- [sql-migrate](https://github.com/rubenv/sql-migrate)
- [tern](https://github.com/jackc/tern)

### goose

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
	ID int32
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
	ID   int32
	Text string
}
```

### golang-migrate

Warning: [golang-migrate specifies](https://github.com/golang-migrate/migrate/blob/master/MIGRATIONS.md#migration-filename-format) that the version number in the migration file name is to be interpreted numerically. However, sqlc executes the migration files in **lexicographic** order. If you choose to simply enumerate your migration versions, make sure to prepend enough zeros to the version number to avoid any unexpected behavior.

Probably doesn't work as intended:
```
1_initial.up.sql
...
9_foo.up.sql
# this migration file will be executed BEFORE 9_foo
10_bar.up.sql
```
Works as was probably intended:
```
001_initial.up.sql
...
009_foo.up.sql
010_bar.up.sql
```

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

### dbmate

```sql
-- migrate:up
CREATE TABLE foo (bar INT NOT NULL);

-- migrate:down
DROP TABLE foo;
```

```go
package db

type Foo struct {
	Bar int32
}
```
