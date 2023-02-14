-- name: Demo :one
SELECT uuid_generate_v5('7c4597a0-8cfa-4c19-8da0-b8474a36440d', $1)::uuid as col1;
