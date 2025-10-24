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

sqlc can generate structs with fields based on the alias inside the macro `sqlc.embed()` by adding the `emit_embed_alias` key to the configuration file as it shows on [configuration reference](../reference/config.md).

```sql
-- name: ListUserLink :many
SELECT
    sqlc.embed(owner),
    sqlc.embed(consumer)
FROM
    user_links
    INNER JOIN users AS owner ON owner.id = user_links.owner_id
    INNER JOIN users AS consumer ON consumer.id = user_links.consumer_id;
```

```
type ListUserLinkRow struct {
	Owner    User
	Consumer User
}
```
