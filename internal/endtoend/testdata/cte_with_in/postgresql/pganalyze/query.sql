-- name: GetAll :many
SELECT * FROM L;

-- name: GetAll1 :many
with recursive cte as (
  select id, L_ID, F from T
  union all
  select c.id, c.L_ID, c.F from T as c where c.L_ID = $1
) select id, l_id, f from cte;

-- name: GetAll2 :many
with recursive cte as (
  select id, L_ID, F from T where T.ID=2
  union all
  select c.id, c.L_ID, c.F from T as c where c.L_ID = $1
) select id, l_id, f from cte;

-- name: GetAll4 :many
select id from T where L_ID in(
  with recursive L as (
    select id, L_ID, F from T where T.ID =2
    union all
    select c.id, c.L_ID, c.F from T as c where c.L_ID = $1
 ) select l_id from L
);

-- name: GetAll3 :many
select id from T where L_ID in(
  with recursive cte as (
    select id, L_ID, F from T where T.ID =2
    union all
    select c.id, c.L_ID, c.F from T as c where c.L_ID = $1
 ) select l_id from cte
);
