-- name: InsertContact :exec
INSERT INTO contacts (
  pid,
  customername
)
VALUES (?, ?);
