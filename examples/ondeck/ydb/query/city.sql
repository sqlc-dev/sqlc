-- name: ListCities :many
SELECT *
FROM city
ORDER BY name;

-- name: GetCity :one
SELECT *
FROM city
WHERE slug = $slug;

-- name: CreateCity :one
-- Create a new city. The slug must be unique.
-- This is the second line of the comment
-- This is the third line
INSERT INTO city (
    name,
    slug
) VALUES (
    $name,
    $slug
) RETURNING *;

-- name: UpdateCityName :exec
UPDATE city
SET name = $name
WHERE slug = $slug;

