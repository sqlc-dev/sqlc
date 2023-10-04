-- name: GetRepro :one
select * from repro where id = ? limit 1;