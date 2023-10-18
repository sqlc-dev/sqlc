-- name: ListItems :one
WITH
    items1 AS (SELECT 'id'::TEXT AS id, 'name'::TEXT AS name),
    items2 AS (SELECT 'id'::TEXT AS id, 'name'::TEXT AS name)
SELECT
    i1.id AS id1,
    i2.id AS id2
FROM
    items1 i1
        JOIN items1 i2 ON 1 = 1;
