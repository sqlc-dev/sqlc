-- name: A :many
SELECT "a""b".* FROM authors AS "a""b";

-- name: B :many
SELECT "x]y".* FROM authors AS "x]y";

-- name: C :many
SELECT [x]]y].* FROM authors AS [x]]y];

-- name: D :many
SELECT "select".* FROM authors AS "select";
