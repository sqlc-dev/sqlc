-- name: GetSystemColumns :one
SELECT
  tableoid, xmin, cmin, xmax, cmax, ctid
FROM test;
