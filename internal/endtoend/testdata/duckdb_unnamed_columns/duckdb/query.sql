-- name: Cast :many
SELECT CAST(age AS BIGINT) FROM people;

-- name: Concat :many
SELECT name || '!' FROM people;

-- name: Two :many
SELECT id + 1, age * 2 FROM people;

-- name: Mixed :many
SELECT id, upper(name), age + 1 AS next_age, id * 2 FROM people;
