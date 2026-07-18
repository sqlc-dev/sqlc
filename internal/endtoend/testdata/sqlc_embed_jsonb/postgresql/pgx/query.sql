-- name: GetParent :one
SELECT sqlc.embed.jsonb(parent_table) AS parent FROM parent_table LIMIT 1;

-- name: GetParentsWithChildren :many
SELECT
  parent_table.id,
  ARRAY(SELECT sqlc.embed.jsonb(child_table) FROM child_table WHERE child_table.parent_id = parent_table.id) AS children
FROM parent_table;
