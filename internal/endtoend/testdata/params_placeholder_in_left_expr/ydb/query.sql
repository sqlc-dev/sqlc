-- name: FindByID :many
SELECT * FROM users WHERE $id = id;

-- name: FindByIDAndName :many
SELECT * FROM users WHERE $id = id AND $id = name;
