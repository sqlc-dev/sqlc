# Macros

## `sqlc.arg`

Attach a name to a parameter in a SQL query. This macro expands to an
engine-specific parameter placeholder. The name of the parameter is noted and
used during code generation.

```sql
-- name: GetAuthorByName :one
SELECT *
FROM authors
WHERE lower(name) = sqlc.arg(name);

-- >>> EXPANDS TO >>>

-- name: GetAuthorByName :one
SELECT *
FROM authors
WHERE lower(name) = ?;
```

See more examples in [Naming parameters](../howto/named_parameters.md).

## `sqlc.embed`

Embedding allows you to reuse existing model structs in more queries, resulting
in less manual serialization work. First, imagine we have the following schema
with students and test scores.

```sql
CREATE TABLE students (
  id   bigserial PRIMARY KEY,
  name text,
  age  integer
);

CREATE TABLE test_scores (
  student_id bigint,
  score integer,
  grade text
);
```

```sql
-- name: GetStudentAndScore :one
SELECT sqlc.embed(students), sqlc.embed(test_scores)
FROM students
JOIN test_scores ON test_scores.student_id = students.id
WHERE students.id = $1;

-- >>> EXPANDS TO >>>

-- name: GetStudentAndScore :one
SELECT students.*, test_scores.*
FROM students
JOIN test_scores ON test_scores.student_id = students.id
WHERE students.id = $1;
```

The Go method will return a struct with a field for the `Student` and field for
the test `TestScore` instead of each column existing on the struct.

```go
type GetStudentAndScoreRow struct {
	Student   Student
	TestScore TestScore
}

func (q *Queries) GetStudentAndScore(ctx context.Context, id int64) (GetStudentAndScoreRow, error) {
    // ...
}
```

See a full example in [Embedding structs](../howto/embedding.md).

## `sqlc.macro`

Expand a macro you defined yourself. A query annotated
[`:macro`](query-annotations.md#macro) is not a query: it generates no code,
and stands for a piece of SQL that is pasted into the queries that call it.

```sql
-- name: ActiveAuthors :macro
-- @param org bigint
SELECT * FROM authors
WHERE org_id = sqlc.arg(org) AND deleted_at IS NULL;

-- name: IsRecent :macro
-- @param ts timestamptz
-- @param days
sqlc.arg(ts) > now() - sqlc.arg(days) * interval '1 day';

-- name: ListRecentAuthors :many
SELECT a.id, a.name
FROM sqlc.macro.ActiveAuthors(sqlc.arg(org_id)) a
WHERE sqlc.macro.IsRecent(a.created_at, 7);

-- >>> EXPANDS TO >>>

-- name: ListRecentAuthors :many
SELECT a.id, a.name
FROM (SELECT * FROM authors WHERE org_id = CAST($1 AS bigint) AND deleted_at IS NULL) a
WHERE (CAST(a.created_at AS timestamptz) > now() - 7 * interval '1 day');
```

```go
func (q *Queries) ListRecentAuthors(ctx context.Context, orgID int64) ([]ListRecentAuthorsRow, error) {
	// ...
}
```

A macro's parameters are declared with `-- @param name [type]` comments, in
the order they are passed, and referenced in its body with `sqlc.arg(name)`
(or `@name`). A call passes one SQL expression for each:

- A parameter **with a type** casts its argument, as `CAST(argument AS type)`.
  The type is written the way the database's `CAST` takes it, so
  `bigint` for PostgreSQL and `signed` for MySQL. Passing `sqlc.arg(org_id)`
  to a typed parameter gives the query a parameter of that type.
- A parameter **without a type** takes its argument as written. An argument
  that is more than a single term is wrapped in parentheses, so the operators
  around the reference cannot regroup it.

The expansion is always wrapped in parentheses, so a call can stand wherever
an expression or a subquery can: in a select list, after `FROM`, `JOIN`, `IN`
or `EXISTS`, as a common table expression (`WITH a AS sqlc.macro.Active()`)
and so on.

Macros are expanded lazily, where they are called:

- A macro is never checked on its own. Its body does not have to be a
  statement, and an error in it is reported against the query that calls it.
  A macro nothing calls is never looked at.
- A macro can be called before it is defined, and from any query file of the
  same package.
- A macro can call other macros. A macro that ends up calling itself is an
  error.
- A `sqlc.arg()` in the body that does not name one of the macro's parameters
  is left alone, and becomes a parameter of every query the macro is expanded
  in.

A macro's body cannot use positional placeholders (`$1`, `?`), which would be
numbered against whichever query it lands in, and cannot hold a string that
spans lines.

Macros are available for the engines sqlc rewrites `sqlc.arg` for:
PostgreSQL, MySQL, SQLite and ClickHouse.

## `sqlc.narg`

The same as `sqlc.arg`, but always marks the parameter as nullable.

```sql
-- name: GetAuthorByName :one
SELECT *
FROM authors
WHERE lower(name) = sqlc.narg(name);

-- >>> EXPANDS TO >>>

-- name: GetAuthorByName :one
SELECT *
FROM authors
WHERE LOWER(name) = ?;
```

See more examples in [Naming parameters](../howto/named_parameters.md).

## `sqlc.slice`

For drivers that do not support passing slices to the IN operator, the
`sqlc.slice` macro generates a dynamic query at runtime with the correct
number of parameters.

```sql
/* name: SelectStudents :many */
SELECT * FROM students 
WHERE age IN (sqlc.slice("ages"))

-- >>> EXPANDS TO >>>

/* name: SelectStudents :many */
SELECT id, name, age FROM authors 
WHERE age IN (/*SLICE:ages*/?)
```

Since the `/*SLICE:ages*/` placeholder is dynamically replaced on a per-query
basis, this macro can't be used with prepared statements.

See a full example in [Passing a slice as a parameter to a
query](../howto/select.md#mysql-and-sqlite).
