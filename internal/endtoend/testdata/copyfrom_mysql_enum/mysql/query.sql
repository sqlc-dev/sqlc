-- name: UpsertExperienceLocations :copyfrom
REPLACE INTO experience_locations (location_id, type)
VALUES (?, ?);
