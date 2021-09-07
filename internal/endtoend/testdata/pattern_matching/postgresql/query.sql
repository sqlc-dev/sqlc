CREATE TABLE pet (name text);

-- name: PetsByName :many
SELECT * FROM pet WHERE name LIKE $1;
