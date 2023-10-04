-- name: PetsByName :many
SELECT * FROM pet WHERE name LIKE ?;
