-- name: GetItems :many
SELECT * FROM items
WHERE id > $1  -- Mandatory parameter
sqlc.optional('Name', 'AND name = $2')
sqlc.optional('Status', 'AND status = $3')
sqlc.optional('Description', 'AND description LIKE $4');

-- name: GetItemsNoMandatory :many
SELECT * FROM items
WHERE 1=1
sqlc.optional('Name', 'AND name = $1')
sqlc.optional('Status', 'AND status = $2');
