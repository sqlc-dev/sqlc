-- name: CreateEvent :exec
INSERT INTO events (id, payload, meta) VALUES ($1, $2, $3);

-- name: GetEvent :one
SELECT id, payload, meta FROM events WHERE id = $1;

-- name: EventKind :one
SELECT json_extract(payload, '$.kind') AS kind, json_extract_string(payload, '$.kind') AS kind_text
FROM events WHERE id = $1;

-- name: EventsByPayload :many
SELECT id FROM events WHERE payload = $1;

-- name: PayloadText :one
SELECT CAST(payload AS VARCHAR) AS payload_text FROM events WHERE id = $1;
