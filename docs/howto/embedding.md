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

#### Nested JSON objects and arrays

`sqlc.jsonb_build_object."Name"(key, value, key, value, ...)` builds a named
Go struct from an inline JSON shape — a typed wrapper around Postgres's
`jsonb_build_object`, so keys must be string literals. This is most useful
for pulling a one-to-many relationship into a single query as a slice
field, instead of running a separate query per parent row. It's **pgx/v5
only**: `sqlc generate` fails with a clear error for any other
`sql_package`.

`"Name"` is required and must be a double-quoted identifier in that
position, not an argument — Postgres parses this as a 3-part qualified
function name (`catalog.schema.name`), and sqlc reads the struct name
directly off of it.

```sql
CREATE TABLE authors (
  id   bigserial PRIMARY KEY,
  name text NOT NULL
);

CREATE TABLE books (
  id        bigserial PRIMARY KEY,
  author_id bigint NOT NULL REFERENCES authors (id),
  title     text NOT NULL
);
```

```sql
-- name: GetAuthors :many
SELECT
  sqlc.embed(authors),
  ARRAY(
    SELECT sqlc.jsonb_build_object."Book"('id', books.id, 'title', books.title)
    FROM books
    WHERE books.author_id = authors.id
  ) AS books
FROM authors;
```

```go
type GetAuthorsRow struct {
	Author Author
	Books  []Book
}

type Book struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
}
```

No custom `Scan`/`Value` methods, no wrapper type — `Books` is a plain
`[]Book`; pgx v5 decodes the `jsonb[]` column into it directly. A lone
`sqlc.jsonb_build_object."Name"(...)` (not wrapped in `ARRAY(...)`) works
the same way and produces a plain struct field instead of a slice:

```sql
-- name: GetAuthorSummary :one
SELECT sqlc.jsonb_build_object."AuthorSummary"('name', name) AS summary FROM authors LIMIT 1;
```

```go
type GetAuthorSummaryRow struct {
	Summary AuthorSummary
}
```

Two queries that use the same explicit name reuse the same Go type, as long
as their shapes match (this works even mixing scalar and `ARRAY(...)` uses
of the same name). A shape mismatch, or a name that collides with an
existing model/enum type, fails generation with an error instead of
emitting Go code that won't compile.

##### Overriding generated names

The struct name and individual field names can be overridden via the
`rename` option:

```json
{
  "rename": {
    "Book": "BookSummary",
    "Book.id": "BookID"
  }
}
```

`"Book"` renames the type; `"Book.id"` renames just the `id` field within
it, without affecting other types that also have an `id` key. Use the plain
key (`"id"`) instead to rename that field everywhere.

##### Embedding a whole row as JSON

Listing every column by hand is tedious when you just want the whole row.
`sqlc.embed.jsonb(table)` builds a JSON object from all of a table's columns,
the same way `sqlc.embed(table)` gives you the table's model — but as a
single JSON value, so it can be nested inside `ARRAY(...)` to return a slice
of rows from one query:

```sql
-- name: GetAuthorsWithBooks :many
SELECT
  authors.id,
  ARRAY(SELECT sqlc.embed.jsonb(books) FROM books WHERE books.author_id = authors.id) AS books
FROM authors;
```

```go
type GetAuthorsWithBooksRow struct {
	ID    int64  `json:"id"`
	Books []Book `json:"books"`
}

type Book struct {
	ID       int64  `json:"id"`
	AuthorID int64  `json:"author_id"`
	Title    string `json:"title"`
}
```

The generated struct is named from the result alias — singularized for the
`ARRAY(...)` case (`AS books` → `Book`), or used as-is for a scalar
`sqlc.embed.jsonb(table) AS author` (→ `Author`). Fields, types and JSON keys
come from the table's columns, so the object always decodes cleanly. Under
the hood the call is rewritten to `to_jsonb(table)`, which you'll see in
`EXPLAIN` output and logs.

Like `sqlc.jsonb_build_object`, this is pgx/v5 only, and the generated name
must not collide with a model or another JSON type; use the `rename` option
if it does.