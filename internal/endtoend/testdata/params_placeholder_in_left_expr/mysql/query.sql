-- name: FindByID :many
SELECT * FROM users WHERE ? = id;

-- name: FindByIDAndName :many
SELECT * FROM users WHERE ? = id AND ? = name;
