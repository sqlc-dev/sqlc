-- name: CountOne :one
SELECT COUNT(1) FROM bar WHERE id = $id AND name <> $name LIMIT $limit;

-- name: CountTwo :one
SELECT COUNT(1) FROM bar WHERE id = $id AND name <> $name;

-- name: CountThree :one
SELECT COUNT(1) FROM bar WHERE id > $id_gt AND phone <> $phone AND name <> $name;

