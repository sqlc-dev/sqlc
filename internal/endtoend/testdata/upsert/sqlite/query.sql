-- https://github.com/kyleconroy/sqlc/issues/1728

CREATE TABLE IF NOT EXISTS locations (
    id              INTEGER PRIMARY KEY,
    name            TEXT    NOT NULL,
    address         TEXT    NOT NULL,
    zip_code        INT     NOT NULL,
    latitude        REAL    NOT NULL,
    longitude       REAL    NOT NULL,
    UNIQUE(name)
);

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
