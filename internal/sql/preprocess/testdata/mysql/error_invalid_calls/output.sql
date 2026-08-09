-- name: WrongFunc :one
SELECT id FROM users WHERE id = sqlc.argh(target_id);

-- name: TooManyArgs :one
SELECT id FROM users WHERE id = sqlc.arg('foo', 'bar');

-- name: TooFewArgs :one
SELECT id FROM users WHERE id = sqlc.arg();

-- name: InvalidArgPlaceholder :one
SELECT id FROM users WHERE id = sqlc.arg(?);
