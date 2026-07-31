-- name: GetParents :many
SELECT
  sqlc.embed(parent_table),
  ARRAY(SELECT sqlc.jsonb_build_object."Child"('x', child_table.x, 'y', child_table.y) FROM child_table WHERE child_table.parent_id = parent_table.id) AS children
FROM parent_table;

-- name: GetParentSummary :one
SELECT sqlc.jsonb_build_object."ParentSummary"('name', name) AS obj FROM parent_table LIMIT 1;

-- name: GetChildIDs :one
SELECT ARRAY(SELECT id FROM child_table WHERE child_table.parent_id = parent_table.id) AS child_ids FROM parent_table LIMIT 1;

-- name: GetChildPoints :many
SELECT ARRAY(SELECT sqlc.jsonb_build_object."ChildPoint"('x', child_table.x, 'y', child_table.y) FROM child_table WHERE child_table.parent_id = parent_table.id) AS points FROM parent_table;

-- name: GetChildPointsAgain :many
SELECT ARRAY(SELECT sqlc.jsonb_build_object."ChildPoint"('x', child_table.x, 'y', child_table.y) FROM child_table WHERE child_table.parent_id = parent_table.id) AS points FROM parent_table WHERE parent_table.id = $1;

-- name: GetNestedParent :one
SELECT sqlc.jsonb_build_object."NestedParent"(
  'name', parent_table.name,
  'children', ARRAY(SELECT sqlc.jsonb_build_object."NestedChild"('x', child_table.x) FROM child_table WHERE child_table.parent_id = parent_table.id)
) AS obj FROM parent_table LIMIT 1;
