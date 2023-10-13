-- name: WrongFunc :one
select id, first_name from users where id = sqlc.argh(target_id);

-- name: TooManyArgs :one
select id, first_name from users where id = sqlc.arg('foo', 'bar');

-- name: TooFewArgs :one
select id, first_name from users where id = sqlc.arg();

-- name: InvalidArgFunc :one
select id, first_name from users where id = sqlc.arg(sqlc.arg(target_id));

-- name: InvalidArgPlaceholder :one
select id, first_name from users where id = sqlc.arg($1);
