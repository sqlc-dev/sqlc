-- name: CountOne :one
SELECT count(1) FROM bar WHERE id = sqlc.arg(id) AND name <> $1 LIMIT sqlc.arg('limit');

-- name: CountTwo :one
SELECT count(1) FROM bar WHERE id = $1 AND name <> sqlc.arg(name);

-- name: CountThree :one
SELECT count(1) FROM bar WHERE id > $2 AND phone <> sqlc.arg(phone) AND name <> $1;
