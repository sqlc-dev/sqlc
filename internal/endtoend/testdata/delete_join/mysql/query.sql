-- name: DeleteJoin :exec
DELETE jt.*,
pt.*
FROM
        join_table as jt
        JOIN primary_table as pt ON jt.primary_table_id = pt.id
WHERE
        jt.id = ?
        AND pt.user_id = ?;

-- name: DeleteLeftJoin :exec
DELETE jt.*,
pt.*
FROM
        join_table as jt
        LEFT JOIN primary_table as pt ON jt.primary_table_id = pt.id
WHERE
        jt.id = ?
        AND pt.user_id = ?;

-- name: DeleteRightJoin :exec
DELETE jt.*,
pt.*
FROM
        join_table as jt
        RIGHT JOIN primary_table as pt ON jt.primary_table_id = pt.id
WHERE
        jt.id = ?
        AND pt.user_id = ?;

-- name: DeleteJoinWithSubquery :exec
DELETE pt
FROM primary_table pt
JOIN (SELECT 1 as id) jt ON pt.id = jt.id;
