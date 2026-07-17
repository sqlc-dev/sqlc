-- name: GetSystemColumns :one
SELECT xmin, cmin, xmax, cmax, ctid, tableoid FROM authors LIMIT 1;

-- name: GetSystemColumnsAliased :one
SELECT a.xmin, a.ctid FROM authors a LIMIT 1;

-- name: SelectStar :many
SELECT * FROM authors;
