-- name: Foo :one
SELECT * FROM register_account('a', 'b');

-- name: GetAccount :one
SELECT * FROM get_account($1, $2);
