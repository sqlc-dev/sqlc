-- name: CreateUser :one
INSERT INTO users (name, bio) VALUES ($1, $2) RETURNING *;

-- name: UpdateBio :exec
UPDATE users SET bio = $1 WHERE id = $2;

-- name: UpdateNameReturning :one
UPDATE users SET name = $1 WHERE id = $2 RETURNING id, name;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;

-- name: DeleteUserReturning :one
DELETE FROM users WHERE id = $1 RETURNING id;

-- name: MergeUsers :many
MERGE INTO users AS t
USING posts AS s
ON t.id = s.user_id AND t.id = $1
WHEN MATCHED AND s.id > $2 THEN
  UPDATE SET (name, bio) = ($3, s.title)
WHEN NOT MATCHED THEN
  INSERT (name, bio) VALUES ($4, s.title)
WHEN NOT MATCHED BY SOURCE AND t.id > $5 THEN
  DELETE
RETURNING merge_action(), *;

-- name: MergeDerivedSource :many
MERGE INTO users AS t
USING (SELECT 1::int8 AS user_id, 1::int4 AS marker) AS s
ON t.id = s.user_id
WHEN NOT MATCHED BY SOURCE THEN
  DELETE
RETURNING s.marker, t.id;
