-- name: AllTypes :many
SELECT * FROM things;

-- name: StarColumns :many
SELECT id, name, tag, amount, tags, labels, matrix, kind, created, updated, price, status, attrs, pos, ip, uid, fixed, flag FROM things;
