-- name: query_return_full_table :many
SELECT * FROM query_return_full_table;

-- name: query_param_enum_table :one
SELECT * FROM query_param_enum_table WHERE value = $1;

-- name: query_param_struct_enum_table :one
SELECT id FROM query_param_struct_enum_table WHERE id = $1 AND value = $2;

-- name: query_return_enum_table :one
SELECT value FROM query_return_enum_table WHERE id = $1;

-- name: query_return_struct_enum_table :one
SELECT value, another FROM query_return_struct_enum_table WHERE id = $1;

-- name: query_sqlc_embed_table :one
SELECT sqlc.embed(query_sqlc_embed_table) FROM query_sqlc_embed_table WHERE id = $1;
