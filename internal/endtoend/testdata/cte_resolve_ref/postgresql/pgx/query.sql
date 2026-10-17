-- name: CTERef :one
WITH t1_ids AS (
    SELECT id FROM t1
)
SELECT * FROM t1_ids WHERE t1_ids.id = sqlc.arg('id');

-- name: CTEMultipleRefs :one
WITH t1_ids AS (
    SELECT id FROM t1 WHERE t1.id = sqlc.arg('id')
),
t2_ids AS (
    SELECT id FROM t2 WHERE t2.id = sqlc.arg('id')
),
all_ids AS (
    SELECT * FROM t1_ids
    UNION
    SELECT * FROM t2_ids
)
SELECT * FROM all_ids AS ids WHERE ids.id = sqlc.arg('id');


