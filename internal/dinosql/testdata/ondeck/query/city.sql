-- name: ListCities :many
SELECT *
FROM city
ORDER BY name;

-- name: GetCity :one
SELECT *
FROM city
WHERE slug = $1;

-- name: CreateCity :one
-- Create a new city. The slug must be unique
INSERT INTO city (
    name,
    slug
) VALUES (
    $1,
    $2
) RETURNING *;

-- name: UpdateCityName :exec
UPDATE city
SET name = $2
WHERE slug = $1;
