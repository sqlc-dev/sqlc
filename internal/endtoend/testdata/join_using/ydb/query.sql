-- name: SelectJoinUsing :many
SELECT t1.fk, SUM(t2.fk) FROM t1 JOIN t2 USING (fk) GROUP BY t1.fk;
