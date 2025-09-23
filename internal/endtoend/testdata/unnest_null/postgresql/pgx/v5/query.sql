-- name: CreateMemories :many
INSERT INTO memories (vampire_id)
SELECT
    unnest(sqlc.narg(vampire_id)::uuid[]) AS vampire_id
RETURNING
    *;

-- name: GetVampireIDs :many
SELECT vampires.id::uuid FROM unnest(sqlc.narg(vampire_id)::uuid[]) AS vampires (id);
