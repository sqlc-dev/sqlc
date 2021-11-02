# Query annotations

sqlc requires each query to have a small comment indicating the name and
command. The format of this comment is as follows:

```sql
-- name: <name> <command>
```

## `:exec`

The generated method will return the error from
[ExecContext](https://golang.org/pkg/database/sql/#DB.ExecContext).

```sql
-- name: DeleteAuthor :exec
DELETE FROM authors
WHERE id = $1;
```

```go
func (q *Queries) DeleteAuthor(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteAuthor, id)
	return err
}
```

## `:execresult`

The generated method will return the [sql.Result](https://golang.org/pkg/database/sql/#Result) returned by
[ExecContext](https://golang.org/pkg/database/sql/#DB.ExecContext).

```sql
-- name: DeleteAllAuthors :execresult
DELETE FROM authors;
```

```go
func (q *Queries) DeleteAllAuthors(ctx context.Context) (sql.Result, error) {
	return q.db.ExecContext(ctx, deleteAllAuthors)
}
```

## `:execrows`

The generated method will return the number of affected rows from the
[result](https://golang.org/pkg/database/sql/#Result) returned by
[ExecContext](https://golang.org/pkg/database/sql/#DB.ExecContext).

```sql
-- name: DeleteAllAuthors :execrows
DELETE FROM authors;
```

```go
func (q *Queries) DeleteAllAuthors(ctx context.Context) (int64, error) {
	_, err := q.db.ExecContext(ctx, deleteAllAuthors)
	// ...
}
```

## `:many`

The generated method will return a slice of records via
[QueryContext](https://golang.org/pkg/database/sql/#DB.QueryContext).

```sql
-- name: ListAuthors :many
SELECT * FROM authors
ORDER BY name;
```

```go
func (q *Queries) ListAuthors(ctx context.Context) ([]Author, error) {
	rows, err := q.db.QueryContext(ctx, listAuthors)
	// ...
}
```

## `:one`

The generated method will return a single record via
[QueryRowContext](https://golang.org/pkg/database/sql/#DB.QueryRowContext).

```sql
-- name: GetAuthor :one
SELECT * FROM authors
WHERE id = $1 LIMIT 1;
```

```go
func (q *Queries) GetAuthor(ctx context.Context, id int64) (Author, error) {
	row := q.db.QueryRowContext(ctx, getAuthor, id)
	// ...
}
```
