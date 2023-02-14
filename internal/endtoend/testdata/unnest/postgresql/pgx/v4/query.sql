-- name: CreateMemories :many
INSERT INTO memories (vampire_id)
SELECT
    unnest(@vampire_id::uuid[]) AS vampire_id
RETURNING
    *;
