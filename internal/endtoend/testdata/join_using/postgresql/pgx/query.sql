-- name: SelectJoinUsing :many
select t1.fk, sum(t2.fk) from t1 join t2 using (fk) group by fk;
