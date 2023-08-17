-- name: Create :exec
insert into ding_depts (id, pid, title)
values ($1, $2, $3);


-- name: CountById :one
SELECT count(*)
FROM ding_depts
where id = $1;

-- name: ListByPid :many
SELECT *
FROM ding_depts
where pid = $1;
