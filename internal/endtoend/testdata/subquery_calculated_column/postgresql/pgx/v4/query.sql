-- name: SubqueryCalcColumn :many
SELECT sum FROM (SELECT a + b AS sum FROM foo) AS f;
