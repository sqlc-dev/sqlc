-- name: UpsertLocation :exec
UPSERT INTO locations (
    id,
    name,
    address,
    zip_code,
    latitude,
    longitude
)
VALUES ($id, $name, $address, $zip_code, $latitude, $longitude);



