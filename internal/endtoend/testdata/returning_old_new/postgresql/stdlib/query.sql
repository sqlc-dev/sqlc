-- name: UpdateReturningOld :one
UPDATE users SET name = $1 WHERE id = $2
RETURNING old.name;

-- name: UpdateReturningOldNew :one
UPDATE users SET name = $1 WHERE id = $2
RETURNING old.name, new.name;

-- name: UpdateReturningOldStar :one
UPDATE users SET name = $1 WHERE id = $2
RETURNING old.*;

-- name: InsertReturningOldNew :one
INSERT INTO users (name) VALUES ($1)
RETURNING old.id, new.id;

-- name: DeleteReturningOldNew :one
DELETE FROM users WHERE id = $1
RETURNING old.name, new.name;

-- name: UpdateReturningWithAliases :one
UPDATE users SET name = $1 WHERE id = $2
RETURNING WITH (OLD AS o, NEW AS n) o.name, n.name;
