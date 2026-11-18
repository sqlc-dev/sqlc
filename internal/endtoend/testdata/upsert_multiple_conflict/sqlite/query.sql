-- name: Upsert :exec
INSERT INTO t (id, name, score) VALUES (?, ?, ?)
ON CONFLICT (id) DO UPDATE SET score = ?
ON CONFLICT (name) DO UPDATE SET created = ?;

-- name: UpsertWhere :exec
INSERT INTO t (id, name, score) VALUES (?, ?, ?)
ON CONFLICT (id) DO UPDATE SET score = ? WHERE score < ?
ON CONFLICT (name) DO UPDATE SET created = ? WHERE score < ?;

-- name: UpsertNothingFirst :exec
INSERT INTO t (id, name, score) VALUES (?, ?, ?)
ON CONFLICT (id) DO NOTHING
ON CONFLICT (name) DO UPDATE SET score = ?;

-- name: UpsertReturning :one
INSERT INTO t (id, name, score) VALUES (?, ?, ?)
ON CONFLICT (id) DO UPDATE SET score = ?
ON CONFLICT (name) DO UPDATE SET created = ?
RETURNING id, created;
