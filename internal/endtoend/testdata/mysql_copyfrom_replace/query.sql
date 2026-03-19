-- name: UpsertLocations :copyfrom
REPLACE INTO locations (id, name, address, latitude, longitude)
VALUES (?, ?, ?, ?, ?);

-- name: IgnoreLocations :copyfrom
INSERT IGNORE INTO locations (id, name, address, latitude, longitude)
VALUES (?, ?, ?, ?, ?);
