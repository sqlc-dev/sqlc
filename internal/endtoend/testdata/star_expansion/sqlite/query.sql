CREATE TABLE foo (a text, b text);

-- name: StarExpansion :many
SELECT *, *, foo.* FROM foo;

-- name: StarQuotedExpansion :many
SELECT "t".* FROM foo "t";
