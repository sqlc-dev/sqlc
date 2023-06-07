# Naming parameters

sqlc tried to generate good names for positional parameters, but sometimes it
lacks enough context. The following SQL generates parameters with less than
ideal names:

```sql
-- name: UpsertAuthorName :one
UPDATE author
SET
  name = CASE WHEN $1::bool
    THEN $2::text
    ELSE name
    END
RETURNING *;
```

```go
type UpdateAuthorNameParams struct {
	Column1   bool   `json:""`
	Column2_2 string `json:"_2"`
}
```

In these cases, named parameters give you the control over field names on the
Params struct.

```sql
-- name: UpsertAuthorName :one
UPDATE author
SET
  name = CASE WHEN sqlc.arg(set_name)::bool
    THEN sqlc.arg(name)::text
    ELSE name
    END
RETURNING *;
```

```go
type UpdateAuthorNameParams struct {
	SetName bool   `json:"set_name"`
	Name    string `json:"name"`
}
```

If the `sqlc.arg()` syntax is too verbose for your taste, you can use the `@`
operator as a shortcut.

```sql
-- name: UpsertAuthorName :one
UPDATE author
SET
  name = CASE WHEN @set_name::bool
    THEN @name::text
    ELSE name
    END
RETURNING *;
```

## Nullable parameters

sqlc infers the nullability of any specified parameters, and often does exactly
what you want. If you want finer control over the nullability of your
parameters, you may use `sqlc.narg()` (**n**ullable arg) to override the default
behavior. Using `sqlc.narg` tells sqlc to ignore whatever nullability it has
inferred and generate a nullable parameter instead. There is no nullable
equivalent of the `@` syntax.

Here is an example that uses a single query to allow updating an author's
name, bio or both.

```sql
-- name: UpdateAuthor :one
UPDATE author
SET
 name = coalesce(sqlc.narg('name'), name),
 bio = coalesce(sqlc.narg('bio'), bio)
WHERE id = sqlc.arg('id')
RETURNING *;
```

The following code is generated:

```go
type UpdateAuthorParams struct {
	Name sql.NullString
	Bio  sql.NullString
	ID   int64
}
```
