# Modifying the database schema

sqlc parses `CREATE TABLE` and `ALTER TABLE` statements in order to generate
the necessary code.

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

sqlc does not perform database migrations for you. However, sqlc is able to
differentiate between up and down migrations. sqlc ignores down migrations when
parsing SQL files.

sqlc supports parsing migrations from the following tools:

- [atlas](https://github.com/ariga/atlas)
- [dbmate](https://github.com/amacneil/dbmate)
- [golang-migrate](https://github.com/golang-migrate/migrate)
- [goose](https://github.com/pressly/goose)
- [sql-migrate](https://github.com/rubenv/sql-migrate)
- [tern](https://github.com/jackc/tern)

To enable migration parsing, specify the migration directory instead of a schema file:

```yaml
version: "2"
sql:
  - engine: "postgresql"
    queries: "query.sql"
    schema: "db/migrations"
    gen:
      go:
        package: "tutorial"
        out: "tutorial"
```

### atlas

```sql
-- Create "post" table
CREATE TABLE "public"."post" ("id" integer NOT NULL, "title" text NULL, "body" text NULL, PRIMARY KEY ("id"));
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

### golang-migrate

**Warning:**
[golang-migrate interprets](https://github.com/golang-migrate/migrate/blob/master/MIGRATIONS.md#migration-filename-format)
migration filenames numerically. However, sqlc parses migration files in
lexicographic order. If you choose to have sqlc enumerate your migration files,
make sure their numeric ordering matches their lexicographic ordering to avoid
unexpected behavior. This can be done by prepending enough zeroes to the
migration filenames.

This doesn't work as intended.

```
1_initial.up.sql
...
9_foo.up.sql
# this migration file will be parsed BEFORE 9_foo
10_bar.up.sql
```

This worked as intended.

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
