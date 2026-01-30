-- name: ListVenues :many
SELECT *
FROM venue
WHERE city = $city
ORDER BY name;

-- name: DeleteVenue :exec
DELETE FROM venue
WHERE slug = $slug AND slug = $slug;

-- name: GetVenue :one
SELECT *
FROM venue
WHERE slug = $slug AND city = $city;

-- name: CreateVenue :one
INSERT INTO venue (
    slug,
    name,
    city,
    created_at,
    spotify_playlist,
    status,
    tags
) VALUES (
    $slug,
    $name,
    $city,
    CurrentUtcTimestamp(),
    $spotify_playlist,
    $status,
    $tags
) RETURNING id;

-- name: UpdateVenueName :one
UPDATE venue
SET name = $name
WHERE slug = $slug
RETURNING id;

-- name: VenueCountByCity :many
SELECT
    city,
    count(*) as count
FROM venue
GROUP BY city
ORDER BY city;

