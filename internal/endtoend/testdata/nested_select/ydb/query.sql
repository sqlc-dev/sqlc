-- name: NestedSelect :one
SELECT latest.id, t.count
FROM (
    SELECT id, MAX(update_time) AS update_time
    FROM test
    WHERE test.id IN sqlc.slice(ids)
        AND test.update_time >= $start_time
    GROUP BY id
) latest
INNER JOIN test t USING (id, update_time);

