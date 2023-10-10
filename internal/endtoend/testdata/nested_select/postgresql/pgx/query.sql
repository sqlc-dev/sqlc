-- name: NestedSelect :one
SELECT latest.id, t.count
FROM (
  SELECT id, max(update_time) AS update_time
  FROM test
  WHERE id = ANY (sqlc.arg('IDs')::bigint[])
    -- ERROR HERE on update_time
    AND update_time >= sqlc.arg('StartTime')
  GROUP BY id
) latest
INNER JOIN test t USING (id, update_time);
