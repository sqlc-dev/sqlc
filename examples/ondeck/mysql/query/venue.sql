/* name: ListVenues :many */
SELECT *
FROM venue
WHERE city = ?
ORDER BY name;

/* name: DeleteVenue :exec */
DELETE FROM venue
WHERE slug = ? AND slug = ?;

/* name: GetVenue :one */
SELECT *
FROM venue
WHERE slug = ? AND city = ?;

/* name: CreateVenue :execresult */
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
    ?,
    ?,
    ?,
    NOW(),
    ?,
    ?,
    ?,
    ?
);

/* name: UpdateVenueName :exec */
UPDATE venue
SET name = ?
WHERE slug = ?;

/* name: VenueCountByCity :many */
SELECT
    city,
    count(*)
FROM venue
GROUP BY 1
ORDER BY 1;
