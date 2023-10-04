-- name: UpdateJoin :exec
UPDATE  join_table
SET     is_active = $1
FROM    primary_table
WHERE   join_table.id = $2
        AND primary_table.user_id = $3
        AND join_table.primary_table_id = primary_table.id;
