/* name: UpsertLocation :exec */
INSERT INTO locations (
    name,
    address,
    zip_code,
    latitude,
    longitude
)
VALUES (?, ?, ?, ?, ?)
ON CONFLICT(name) DO UPDATE SET 
    name = excluded.name,
    address = excluded.address,
    zip_code = excluded.zip_code,
    latitude = excluded.latitude,
    longitude = excluded.longitude;
