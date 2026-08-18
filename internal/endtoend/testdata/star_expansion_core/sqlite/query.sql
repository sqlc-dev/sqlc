-- name: StarExpansion :many
SELECT *, *, foo.* FROM foo;

-- name: StarQuotedExpansion :many
SELECT "t".* FROM foo "t";

-- name: StarExpansionJoin :many
SELECT * FROM foo, bar;

-- name: StarExpansionSubquery :many
SELECT * FROM (SELECT * FROM bar) sub;

-- name: StarExpansionCTE :many
WITH t AS (SELECT * FROM bar) SELECT * FROM t;

-- name: StarExpansionReturning :one
INSERT INTO bar (a, c) VALUES (?, ?) RETURNING *;

-- name: StarExpansionAliasedSubquery :many
SELECT (SELECT * FROM baz LIMIT 1) AS x FROM foo GROUP BY x;
