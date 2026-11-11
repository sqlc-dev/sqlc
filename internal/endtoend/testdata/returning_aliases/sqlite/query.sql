-- name: InsertReturningAliases :one
INSERT INTO t (name, score) VALUES (?, ?)
RETURNING id AS new_id, score * 2 AS doubled;

-- name: InsertReturningImplicitAliases :one
INSERT INTO t (name, score) VALUES (?, ?)
RETURNING id new_id, score * 2 doubled;

-- name: InsertReturningQuotedAliases :one
INSERT INTO t (name, score) VALUES (?, ?)
RETURNING id AS "newId", name AS [label];

-- name: InsertReturningSingleAlias :one
INSERT INTO t (name, score) VALUES (?, ?)
RETURNING id AS new_id;

-- name: InsertReturningStarAndAlias :one
INSERT INTO t (name, score) VALUES (?, ?)
RETURNING *, score * 2 AS doubled;

-- name: InsertReturningDuplicateAliases :one
INSERT INTO t (name, score) VALUES (?, ?)
RETURNING id AS old_id, id AS row_id;

-- name: InsertOnConflictReturningAliases :one
INSERT INTO t (id, name, score) VALUES (?, ?, ?)
ON CONFLICT (id) DO UPDATE SET score = excluded.score
RETURNING id AS new_id, score AS total;

-- name: UpdateReturningAliases :one
UPDATE t SET score = score + 1 WHERE id = ?
RETURNING id AS updated_id, score - 1 AS previous;

-- name: DeleteReturningAliases :one
DELETE FROM t WHERE id = ?
RETURNING id AS deleted_id, upper(name) AS shout;
