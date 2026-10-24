-- name: ListItems :many
SELECT id, name FROM items;

-- name: ExportItems :exec
COPY (SELECT id, name FROM items) TO 'items.parquet' (FORMAT PARQUET);
