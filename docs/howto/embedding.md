#### Embedding structs

Embedding allows you to reuse existing model structs in more queries, resulting
in less manual serialization work. First, imagine we have the following schema
with students and test scores.

```sql
CREATE TABLE students (
  id   bigserial PRIMARY KEY,
  name text NOT NULL,
  age  integer NOT NULL
);

CREATE TABLE test_scores (
  student_id bigint NOT NULL,
  score      integer NOT NULL,
  grade      text NOT NULL
);
```

We want to select the student record and the scores they got on a test.
Here's how we'd usually do that:

```sql
-- name: ScoreAndTests :many
SELECT students.*, test_scores.*
FROM students
JOIN test_scores ON test_scores.student_id = students.id
WHERE students.id = $1;
```

When using Go, sqlc will produce a struct like this:

```go
type ScoreAndTestsRow struct {
	ID        int64
	Name      string
	Age       int32
	StudentID int64
	Score     int32
	Grade     string
}
```

With embedding, the struct will contain a model for both tables instead of a
flattened list of columns.

```sql
-- name: ScoreAndTests :many
SELECT sqlc.embed(students), sqlc.embed(test_scores)
FROM students
JOIN test_scores ON test_scores.student_id = students.id
WHERE students.id = $1;
```

```
type ScoreAndTestsRow struct {
	Student   Student
	TestScore TestScore
}
```

#### JSON objects and arrays

`sqlc.jsonb_build_object."Name"(key, value, ...)` returns a JSON-shaped
result decoded straight into a named Go struct (or `[]struct` inside
`ARRAY(...)`), useful for pulling a one-to-many relationship into a single
query instead of one query per parent row. It's **pgx/v5 only** — `sqlc
generate` fails with a clear error for any other `sql_package`. `"Name"` is
required and is read off the 3-part qualified function name Postgres parses
(`catalog.schema.name`), not an argument; keys must be string literals.

```sql
-- name: GetAuthors :many
SELECT
  authors.id,
  ARRAY(
    SELECT sqlc.jsonb_build_object."Book"('id', books.id, 'title', books.title)
    FROM books WHERE books.author_id = authors.id
  ) AS books
FROM authors;
```

```go
type GetAuthorsRow struct {
	ID    int64
	Books []Book
}

type Book struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
}
```

The call is rewritten to Postgres's `jsonb_build_object`, which you'll see in
`EXPLAIN` output. Two queries that use the same name share one Go type when
their shapes match; a shape mismatch, or a name that collides with a model or
another JSON type, fails generation. Names and fields can be overridden with
`rename`: `"Book"` renames the type, `"Book.id"` just the `id` field within
it.