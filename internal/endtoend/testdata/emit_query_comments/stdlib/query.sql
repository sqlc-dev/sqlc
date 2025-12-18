-- name: ListBar :many
-- Lists all bars
SELECT id FROM (
  SELECT * FROM bar
) bar;

-- name: RemoveBar :exec
DELETE FROM bar WHERE id = $1;
