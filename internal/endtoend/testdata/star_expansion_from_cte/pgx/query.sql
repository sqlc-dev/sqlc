CREATE TABLE foo (a text, b text);
CREATE TABLE bar (c text, d text);
-- name: StarExpansionCTE :many
WITH cte AS (SELECT * FROM foo) SELECT * FROM cte;
