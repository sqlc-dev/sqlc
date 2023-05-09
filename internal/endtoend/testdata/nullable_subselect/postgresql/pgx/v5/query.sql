CREATE TABLE foo (a int not null, b int);

-- name: SubqueryWithWhereClause :many
SELECT a, (SELECT COUNT(a) FROM foo WHERE a > 10) as "total" FROM foo;

-- name: SubqueryWithHavingClause :many
SELECT a, (SELECT COUNT(a) FROM foo GROUP BY b HAVING COUNT(a) > 10) as "total" FROM foo;
