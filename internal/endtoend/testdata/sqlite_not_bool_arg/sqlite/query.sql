-- name: Example :many
select 1 where not @argname;

-- name: Example2 :many
select 1 where cast(@argname as boolean);
