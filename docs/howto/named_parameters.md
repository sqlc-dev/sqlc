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
