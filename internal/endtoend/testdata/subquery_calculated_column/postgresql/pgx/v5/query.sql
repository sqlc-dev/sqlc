CREATE TABLE foo (a int, b int);

-- name: SubqueryCalcColumn :many
SELECT sum FROM (SELECT a + b AS sum FROM foo) AS f;
