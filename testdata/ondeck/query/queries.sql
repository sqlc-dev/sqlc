-- name: ListCities :many
SELECT *
FROM city
ORDER BY name;

-- name: GetCity :one
SELECT *
FROM city
WHERE slug = $1;

-- name: ListVenues :many
SELECT *
FROM venue
WHERE city = $1
ORDER BY name;

-- name: DeleteVenue :exec
DELETE FROM venue
WHERE slug = $1;

-- name: GetVenue :one
SELECT *
FROM venue
WHERE slug = $1 AND city = $2;

-- name: CreateCity :one
INSERT INTO city (
    name,
    slug
) VALUES (
    $1,
    $2
) RETURNING *;

-- name: CreateVenue :one
INSERT INTO venue (
    name,
    slug,
    created_at,
    spotify_playlist,
    city
) VALUES (
    $1,
    $2,
    NOW(),
    $3,
    $4
) RETURNING id;

-- name: UpdateCityName :exec
UPDATE city
SET name = $2
WHERE slug = $1;

-- name: UpdateVenueName :one
UPDATE venue
SET name = $2
WHERE slug = $1
RETURNING id;
