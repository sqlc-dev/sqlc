# Sharing a row type between queries

sqlc generates a struct for the rows of each query that returns more than one
column, named after the query. Queries that select the same columns therefore
return different types, even though the structs are identical.

```sql
CREATE TABLE authors (
  id   BIGSERIAL PRIMARY KEY,
  name text      NOT NULL
);

CREATE TABLE books (
  id        BIGSERIAL PRIMARY KEY,
  author_id bigint    NOT NULL REFERENCES authors (id),
  title     text      NOT NULL
);
```

```sql
-- name: GetBook :one
SELECT books.id, books.title, authors.name AS author_name
FROM books
JOIN authors ON authors.id = books.author_id
WHERE books.id = $1;

-- name: ListBooksByAuthor :many
SELECT books.id, books.title, authors.name AS author_name
FROM books
JOIN authors ON authors.id = books.author_id
WHERE books.author_id = $1;
```

```go
type GetBookRow struct {
	ID         int64
	Title      string
	AuthorName string
}

type ListBooksByAuthorRow struct {
	ID         int64
	Title      string
	AuthorName string
}
```

To have them return one type, name it with `:type <TypeName>` after the
command:

```sql
-- name: GetBook :one :type BookWithAuthor
SELECT books.id, books.title, authors.name AS author_name
FROM books
JOIN authors ON authors.id = books.author_id
WHERE books.id = $1;

-- name: ListBooksByAuthor :many :type BookWithAuthor
SELECT books.id, books.title, authors.name AS author_name
FROM books
JOIN authors ON authors.id = books.author_id
WHERE books.author_id = $1;
```

```go
type BookWithAuthor struct {
	ID         int64
	Title      string
	AuthorName string
}

func (q *Queries) GetBook(ctx context.Context, id int64) (BookWithAuthor, error) {
	// ...
}

func (q *Queries) ListBooksByAuthor(ctx context.Context, authorID int64) ([]BookWithAuthor, error) {
	// ...
}
```

The type name works with [embedded structs](embedding.md) too, and a single
query can use it just to choose the name of its row type.

## Rules

- Every query that names a type must return the same fields, in the same
  order, with the same Go types and struct tags. Otherwise `sqlc generate`
  fails and says which column differs:

  ```
  query ListBooks: :type BookWithAuthor does not match query GetBook: column 3 is AuthorID int64, want AuthorName string
  ```

- The query must return more than one column, since a query with a single
  column returns that column's value instead of a struct.
- The name must not be one sqlc already uses for a table's model. Queries
  whose columns are exactly those of a table already return its model, such as
  `Book`, without an annotation.
- The type name always wins: a query annotated with `:type` returns that type
  even when its columns match a table's model.
