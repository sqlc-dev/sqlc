-- name: CreateMemories :many
INSERT INTO memories (vampire_id, memory, victim)
SELECT
    unnest(@vampires::uuid[], @memories::text[], @victims::text[])
RETURNING
    *;

-- name: GetMemories :many
SELECT vampires.id::uuid, vampires.memory::text, vampires.victim::text FROM unnest(@vampires::uuid[], @memories::text[], @victims::text[]) AS vampires (id, memory, victim);
