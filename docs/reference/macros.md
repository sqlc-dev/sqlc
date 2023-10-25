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

See more examples in [Naming parameters](../howto/named_parameters).

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

See a full example in [Embedding structs](../howto/embedding).

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

See more examples in [Naming parameters](../howto/named_parameters).

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
