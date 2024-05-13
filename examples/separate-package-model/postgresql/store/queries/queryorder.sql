-- name: GetOrder :one
SELECT * FROM "content"."order"
WHERE id = $1 LIMIT 1;

-- name: ListOrders :many
SELECT * FROM "content"."order"
ORDER BY number;

-- name: CreateOrder :one
INSERT INTO "content"."order" (
          number, user_id, created_at
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: DeleteOrder :exec
DELETE FROM "content"."order"
WHERE id = $1;
