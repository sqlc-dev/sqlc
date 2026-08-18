-- name: GetUpdateableWishlistItemIDs :many
SELECT id
FROM wishlist_item
WHERE DATEDIFF(date_from, NOW()) >= sqlc.arg('min_days_to_date_from')
  AND DATEDIFF(date_from, NOW()) <= sqlc.arg('max_days_to_date_from')
  AND updated_at < sqlc.arg('updated_by');
