-- name: InsertReturningAlias :one
INSERT INTO t (name, score) VALUES (?, ?)
RETURNING id AS new_id, score * 2 AS doubled;

-- name: InsertReturningAliasNoAs :one
INSERT INTO t (name, score) VALUES (?, ?)
RETURNING id new_id, score * 2 doubled;

-- name: InsertReturningQuotedAlias :one
INSERT INTO t (name, score) VALUES (?, ?)
RETURNING id AS "newId", name AS [label];

-- name: InsertReturningSingleAlias :one
INSERT INTO t (name, score) VALUES (?, ?)
RETURNING id AS new_id;

-- name: InsertReturningStarAndAlias :one
INSERT INTO t (name, score) VALUES (?, ?)
RETURNING *, score * 2 AS doubled;

-- name: InsertReturningSameColumnTwice :one
INSERT INTO t (name, score) VALUES (?, ?)
RETURNING id AS old_id, id AS row_id;

-- name: UpsertReturningAlias :one
INSERT INTO t (name, score) VALUES (?, ?)
ON CONFLICT (name) DO UPDATE SET score = score + excluded.score
RETURNING id AS new_id, score AS total;

-- name: UpdateReturningAlias :many
UPDATE t SET score = ? WHERE name = ?
RETURNING id AS updated_id, score - 1 AS previous;

-- name: DeleteReturningAlias :many
DELETE FROM t WHERE score < ?
RETURNING id AS deleted_id, upper(name) AS shout;
