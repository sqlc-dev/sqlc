-- name: GetPendingSaleStatuses :many
WITH w_pending_sale_status as (
    select * from
    (values ('SAVED'), ('IDLE'), ('IN PROGRESS'), ('HELD'))
    as pending_sale_status(status)
)
SELECT status FROM w_pending_sale_status;
