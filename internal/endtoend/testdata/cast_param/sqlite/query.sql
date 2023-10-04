-- name: GetData :many
select *
from my_table
where (cast(sqlc.arg(allow_invalid) as boolean) or not invalid);