-- name: GetUser :one
SELECT * FROM "auth"."user"
WHERE id = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM "auth"."user"
ORDER BY login;

-- name: CreateUser :one
INSERT INTO "auth"."user" (
          login, password, created_at
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM "auth"."user"
WHERE id = $1;
