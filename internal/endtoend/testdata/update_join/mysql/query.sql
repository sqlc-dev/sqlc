-- name: UpdateJoin :exec
UPDATE  join_table as jt
        JOIN primary_table as pt
            ON jt.primary_table_id = pt.id
SET     jt.is_active = ?
WHERE   jt.id = ?
        AND pt.user_id = ?;

-- name: UpdateLeftJoin :exec
UPDATE  join_table as jt
        LEFT JOIN primary_table as pt
            ON jt.primary_table_id = pt.id
SET     jt.is_active = ?
WHERE   jt.id = ?
        AND pt.user_id = ?;

-- name: UpdateRightJoin :exec
UPDATE  join_table as jt
        RIGHT JOIN primary_table as pt
            ON jt.primary_table_id = pt.id
SET     jt.is_active = ?
WHERE   jt.id = ?
        AND pt.user_id = ?;
