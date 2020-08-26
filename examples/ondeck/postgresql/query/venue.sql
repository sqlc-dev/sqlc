-- name: ListVenues :many
SELECT *
FROM venue
WHERE city = $1
ORDER BY name;

-- name: DeleteVenue :exec
DELETE FROM venue
WHERE slug = $1 AND slug = $1;

-- name: GetVenue :one
SELECT *
FROM venue
WHERE slug = $1 AND city = $2;

-- name: CreateVenue :one
INSERT INTO venue (
    slug,
    name,
    city,
    created_at,
    spotify_playlist,
    status,
    statuses,
    tags
) VALUES (
    $1,
    $2,
    $3,
    NOW(),
    $4,
    $5,
    $6,
    $7
) RETURNING id;

-- name: UpdateVenueName :one
UPDATE venue
SET name = $2
WHERE slug = $1
RETURNING id;

-- name: VenueCountByCity :many
SELECT
    city,
    count(*)
FROM venue
GROUP BY 1
ORDER BY 1;
