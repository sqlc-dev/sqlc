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
