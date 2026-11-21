-- name: PublicPeople :many
SELECT * EXCLUDE (secret) FROM people;

-- name: RenumberedPeople :many
SELECT * REPLACE (id + 1 AS id) FROM people;
