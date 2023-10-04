-- name: StarExpansionCTE :many
WITH cte AS (SELECT * FROM foo) SELECT * FROM cte;
