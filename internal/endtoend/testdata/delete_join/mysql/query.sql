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