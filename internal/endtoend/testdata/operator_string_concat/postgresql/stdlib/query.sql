-- name: Test2 :one
select * from Demo
where txt like '%' || sqlc.arg('val')::text || '%';

-- name: Test3 :one
select * from Demo
where txt like concat('%', sqlc.arg('val')::text, '%');
