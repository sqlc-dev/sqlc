-- name: SetBoth :exec
UPDATE t SET (name, score) = (?, ?) WHERE id = ?;

-- name: SetBothReversed :exec
UPDATE t SET (score, name) = (?, ?) WHERE id = ?;

-- name: SetBothAndSingle :exec
UPDATE t SET (name, score) = (?, ?), created = ? WHERE id = ?;

-- name: UpsertSetBoth :exec
INSERT INTO t (id, name, score) VALUES (?, ?, ?)
ON CONFLICT (id) DO UPDATE SET (name, score) = (?, ?);
