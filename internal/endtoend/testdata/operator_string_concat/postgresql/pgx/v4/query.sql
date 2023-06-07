CREATE TABLE demo (txt text not null);

-- name: Test :one
select * from Demo
where txt ~~ '%' || sqlc.arg('val') || '%';

-- name: Test2 :one
select * from Demo
where txt like '%' || sqlc.arg('val') || '%';

-- name: Test3 :one
select * from Demo
where txt like concat('%', sqlc.arg('val'), '%');
